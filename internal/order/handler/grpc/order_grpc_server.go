package grpc_server

import (
	"context"
	"fmt"
	"go-store/internal/dto"
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
	UnimplementedOrderHandlerServer
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
	RegisterOrderHandlerServer(gsrv, &srv)
	// Register reflection service on gRPC server.
	reflection.Register(gsrv)
	return gsrv, lis, nil
}

func (g *grpcServer) OrderList(ctx context.Context, orderGrpc *OrderRequest) (*OrderResponse, error) {
	glog := log.WithContext(ctx).WithFields(log.Fields{
		"grpc": "OrderHandler",
	})

	claim := &entity.Users{}

	filter := &dto.OrderListFilter{}

	ordersUsecase, err := g.usecases.OrderUsecase.GetOrders(ctx, claim, filter, int(orderGrpc.Limit), int(orderGrpc.Offset))
	if err != nil {
		glog.WithError(err).Error("OrderHandler - error while processing g.usecases.OrderUsecase.GetOrders")
		return nil, status.Error(codes.Internal, err.Error())
	}
	var res OrderResponse
	orderResp := make([]*OrdersJson, len(ordersUsecase))
	for idx, orders := range ordersUsecase {
		orderResp[idx] = mapOrderToJSON(orders)

		orderRespItem := make([]*OrderItemSkuJson, len(orders.OrderItem))
		for idx, ordersItem := range orders.OrderItem {
			orderRespItem[idx] = mapOrderItemToJSON(ordersItem)
		}
		orderResp[idx].OrderItems = orderRespItem
	}
	res.OrdersJson = orderResp
	return &res, nil
}

func mapOrderToJSON(s *entity.OrderJson) *OrdersJson {
	return &OrdersJson{
		Id:      s.Id.String(),
		UserId:  fmt.Sprint(s.UserId),
		Address: s.Address,
		Phone:   s.Phone,
		Comment: s.Comment,
		Status:  s.Status,
	}
}

func mapOrderItemToJSON(s entity.OrderItemSkuJson) *OrderItemSkuJson {
	return &OrderItemSkuJson{
		Item_Id:   int32(s.ItemId),
		SkuName:   s.SkuName,
		Quantity:  int32(s.Quantity),
		Price:     s.Price,
		SmallName: s.SmallName,
	}
}
