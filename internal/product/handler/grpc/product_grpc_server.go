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
	UnimplementedProductHandlerServer
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
	RegisterProductHandlerServer(gsrv, &srv)
	// Register reflection service on gRPC server.
	reflection.Register(gsrv)
	return gsrv, lis, nil
}

func (g *grpcServer) Products(ctx context.Context, productGrpc *ProductRequest) (*ProductReponse, error) {
	glog := log.WithContext(ctx).WithFields(log.Fields{
		"grpc": "Products",
	})

	catId := int(productGrpc.Filter.CategoryId)
	brandId := int(productGrpc.Filter.BrandId)
	regionId := int(productGrpc.Filter.RegionId)
	start := float32(productGrpc.Filter.SkuPriceStart)
	end := float32(productGrpc.Filter.SkuPriceEnd)

	filter := &dto.ProductListFilter{
		ProductName: productGrpc.Filter.ProductName,
		Description: productGrpc.Filter.Description,
		CategoryId:  &catId,
		BrandId:     &brandId,
		RegionId:    &regionId,
		PriceStart:  &start,
		PriceEnd:    &end,
	}

	userRole := entity.UserRole("USER")
	productUsecase, err := g.usecases.ProductUsecase.GetSku(ctx, int(productGrpc.Limit), int(productGrpc.Offset), &userRole, filter)
	if err != nil {
		glog.WithError(err).Error("Products - error while processing g.usecases.ProductUsecase.GetSku")
		return nil, status.Error(codes.Internal, err.Error())
	}
	var res ProductReponse
	productSku := make([]*SkuJson, len(productUsecase.SkuJson))
	for idx, sku := range productUsecase.SkuJson {
		productSku[idx] = mapProductSku(sku)
	}
	res.Total = fmt.Sprint(productUsecase.Total)
	res.SkuJson = productSku
	return &res, nil
}

func mapProductSku(s *entity.SkuJson) *SkuJson {
	return &SkuJson{
		ProductName: s.ProductName,
		Description: s.Description,
		CategoryId:  fmt.Sprint(s.CategoryId),
		CreateTs:    s.CreateTs.String(),
		SkuId:       fmt.Sprint(s.SkuId),
		SkuCode:     fmt.Sprint(s.SkuCode),
		SkuPrice:    fmt.Sprint(s.SkuPrice),
		SkuQuantity: fmt.Sprint(s.SkuQuantity),
		SkuImage:    fmt.Sprint(s.SkuImage),
		SkuValueId:  s.SkuValueId,
	}
}
