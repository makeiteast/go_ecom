package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go-store/conf"
	"go-store/internal/entity"
	"go-store/utils/cachestore"
	"go-store/utils/database"
	v1 "go-store/utils/http"

	_cartHttp "go-store/internal/cart/handler/http"
	_catHttp "go-store/internal/category/handler/http"
	_orderHttp "go-store/internal/order/handler/http"
	_prodHttp "go-store/internal/product/handler/http"
	_userHttp "go-store/internal/user/handler/http"

	_orderGrpc "go-store/internal/order/handler/grpc"
	_prodGrpc "go-store/internal/product/handler/grpc"
	_userGrpc "go-store/internal/user/handler/grpc"

	_cartRepo "go-store/internal/cart/repository/pgsql"
	_catRepo "go-store/internal/category/repository/pgsql"
	_optionRepo "go-store/internal/option/repository/pgsql"
	_orderRepo "go-store/internal/order/repository/pgsql"
	_prodRepo "go-store/internal/product/repository/pgsql"
	_userRepo "go-store/internal/user/repository/pgsql"

	_prodRedisRepo "go-store/internal/product/repository/redis"
	_authRedisRepo "go-store/internal/user/repository/redis"

	_cartUsecase "go-store/internal/cart/usecase"
	_catUsecase "go-store/internal/category/usecase"
	_orderUsecase "go-store/internal/order/usecase"
	_prodUsecase "go-store/internal/product/usecase"
	_authUsecase "go-store/internal/user/usecase"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-access-token")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	mLog := log.WithFields(log.Fields{"func": "main"})

	configs, err := conf.NewService()
	if err != nil {
		// handle db connection error and log it
		mLog.WithError(err).Fatal("Config err")
		return
	}

	dbConf, err := configs.Database()
	if err != nil {
		mLog.WithError(err).Fatal("db Conf err")
		return
	}
	// connect to postgresql with pgxpool package
	// use connection string which was compiled earlier
	dbConn, err := database.NewPgxAccess(dbConf)
	if err != nil {
		// handle db connection error and log it
		mLog.WithError(err).Fatal("db Connection err")
		return
	}
	// defer before databse connection will closed
	mLog.Info("Database connected")

	cacheCfg, err := configs.Cache()
	if err != nil {
		// handle cache config error and log it
		mLog.WithError(err).Fatal("cache Conf err")
		return
	}

	redisClient, err := cachestore.NewRedisClient(cacheCfg)
	if err != nil {
		// handle cache connection error and log it
		mLog.WithError(err).Fatal("cache connection err")
		return
	}
	defer redisClient.Close()
	mLog.Info("Redis connected")

	// next code paragraphs at bottom compiled to create clean architecture
	// first send  database connection to authentification repository function
	// which in turn will transmit dbConn to user repository interface
	// then create variable of repository
	authPgxRepo := _userRepo.NewPgxAuthPgxRepository(dbConn)
	prodRepo := _prodRepo.NewPgxProductRepository(dbConn)
	orderRepo := _orderRepo.NewPgxOrderRepository(dbConn)
	categoryRepo := _catRepo.NewPgxCategoryRepository(dbConn)
	optionRepo := _optionRepo.NewPgxOptionRepository(dbConn)
	cartRepo := _cartRepo.NewPgxCartRepository(dbConn)

	authRedisRepo := _authRedisRepo.NewAuthRedisRepo(redisClient)
	prodRedisRepo := _prodRedisRepo.NewProdRedisRepo(redisClient)

	// Second send repository variable to usecase(Application Buseness Rule, usecase) interface
	// which contain available methods of usecase. In this way we can access to repository methods from usecase
	// then create varible of usecase
	authUsecase := _authUsecase.NewAuthUsecase(authPgxRepo, authRedisRepo)
	prodUsecase := _prodUsecase.NewProductUsecase(prodRepo, prodRedisRepo, optionRepo)
	orderUsecase := _orderUsecase.NewOrderUsecase(orderRepo)
	categoryUsecase := _catUsecase.NewCategoryUsecase(categoryRepo, optionRepo)
	cartUsecase := _cartUsecase.NewCartUsecase(cartRepo, prodRepo)

	uc := &entity.Usecases{
		UserUsecase:     authUsecase,
		ProductUsecase:  prodUsecase,
		OrderUsecase:    orderUsecase,
		CategoryUsecase: categoryUsecase,
		CartUsecase:     cartUsecase,
	}

	// HTTP Server
	httpConf, err := configs.HTTP()
	if err != nil {
		// handle http connection error and log it
		mLog.Warning("http conf err: ", err)
		return
	}

	secrets := configs.Token()

	router := gin.New()
	router.Use(CORSMiddleware())
	// Options
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Static("/static", "./static")
	middleware := _userHttp.ValidateJWT(authUsecase, secrets)
	// Routers
	h := router.Group("/api/v1")
	{
		_orderHttp.NewOrderHandler(h, middleware, uc, mLog)
		_userHttp.NewUserHandler(h, middleware, uc, mLog, secrets)
		_prodHttp.NewProductHandler(h, middleware, uc, mLog)
		_catHttp.NewCategoryHandler(h, middleware, uc, mLog)
		_cartHttp.NewCartHandler(h, middleware, uc, mLog)
	}
	// v1.NewRouter(router, mLog, uc)
	httpServer, err := v1.NewService(router, httpConf)
	if err != nil {
		// handle http connection error and log it
		mLog.Warning("http Connection err: ", err)
		return
	}

	grpcCfg := configs.GrpcConf()
	type grpcServer struct {
		usecases *entity.Usecases
		_userGrpc.UnimplementedGrpcHandlerServer
		_orderGrpc.UnimplementedOrderHandlerServer
		_prodGrpc.UnimplementedProductHandlerServer
	}
	host := fmt.Sprintf("%s:%s", grpcCfg.Host, grpcCfg.Port)
	lis, err := net.Listen("tcp", host)
	if err != nil {
		mLog.Warning("grpc Connection Create: ", err)
		return
	}
	mLog.Info("gRPC server, listening on ", host)

	gsrv := grpc.NewServer()
	srv := grpcServer{
		usecases: uc,
	}
	_userGrpc.RegisterGrpcHandlerServer(gsrv, &srv)
	_orderGrpc.RegisterOrderHandlerServer(gsrv, &srv)
	_prodGrpc.RegisterProductHandlerServer(gsrv, &srv)
	// Register reflection service on gRPC server.
	reflection.Register(gsrv)

	defer gsrv.GracefulStop()

	go func() {
		if err := gsrv.Serve(lis); err != nil {
			mLog.WithError(err).Panic("Cannot serve user grpc server")
			return
		}
	}()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-interrupt:
		mLog.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		mLog.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		mLog.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
