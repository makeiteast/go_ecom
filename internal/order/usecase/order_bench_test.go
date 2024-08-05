package usecase

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"go-store/internal/dto"
	"go-store/internal/entity"
	mocks "go-store/internal/order/mock"
)

func BenchmarkOrder(b *testing.B) {
	any := gomock.Any()
	ctx := context.Background()
	mockCtrl := gomock.NewController(b)
	defer mockCtrl.Finish()

	storageMock := mocks.NewMockOrderRepository(mockCtrl)

	b.Run("get order success", func(b *testing.B) {
		storageMock.EXPECT().GetOrders(ctx, any, any, any, any).Return([]entity.Order{{Id: uuid.UUID{3}}}, nil).Times(1)
		storageMock.EXPECT().GetItem(ctx, any).Return([]entity.OrderItem{{Id: 5}}, nil).Times(1)

		ordUsc := &OrderUsecase{
			orderRepo: storageMock,
		}
		claim := &entity.Users{
			Role: entity.UserRoleAdmin,
		}
		filter := &dto.OrderListFilter{}

		ordUsc.GetOrders(ctx, claim, filter, 5, 0)
	})
}
