package grpc_server

import (
	"context"
	"fmt"
	"go-store/internal/entity"
	"net"

	"github.com/bnkamalesh/webgo/v6"
	log "github.com/sirupsen/logrus"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	status "google.golang.org/grpc/status"
)

type grpcServer struct {
	usecases *entity.Usecases
	UnimplementedGrpcHandlerServer
}

func NewGRPCServer(cfg *entity.Config, us *entity.Usecases) (*grpc.Server, net.Listener, error) {
	host := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	lis, err := net.Listen("tcp", host)
	if err != nil {
		return nil, nil, err
	}
	webgo.LOGHANDLER.Info("gRPC server, listening on", cfg.Host, cfg.Port)

	gsrv := grpc.NewServer()
	srv := grpcServer{
		usecases: us,
	}
	RegisterGrpcHandlerServer(gsrv, &srv)
	// Register reflection service on gRPC server.
	reflection.Register(gsrv)
	return gsrv, lis, nil
}

func (g *grpcServer) LoginHandler(ctx context.Context, loginGrpc *LoginRequest) (*LoginResponse, error) {
	glog := log.WithContext(ctx).WithFields(log.Fields{
		"grpc": "LoginHandler",
	})
	tokenConf := &entity.TokenConf{}
	loginUsecase, err := g.usecases.UserUsecase.Login(ctx, loginGrpc.Username, loginGrpc.Password, tokenConf)
	if err != nil {
		glog.WithError(err).Error("LoginHandler - error while processing g.usecases.UserUsecase.LoginHandler")
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := LoginResponse{
		Username: loginUsecase.Username,
		Role:     string(loginUsecase.Role),
	}
	return &res, nil
}
