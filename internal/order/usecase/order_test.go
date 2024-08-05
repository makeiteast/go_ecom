package usecase

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"go-store/internal/dto"
	"go-store/internal/entity"
	mocks "go-store/internal/order/mock"
)

func TestOrder(t *testing.T) {
	req := require.New(t)
	any := gomock.Any()
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	storageMock := mocks.NewMockOrderRepository(mockCtrl)

	t.Run("get order success", func(t *testing.T) {
		storageMock.EXPECT().GetOrders(ctx, any, any, any, any).Return([]*entity.Order{{Id: uuid.UUID{3}}}, nil).Times(1)
		storageMock.EXPECT().GetItem(ctx, any).Return([]*entity.OrderItem{{Id: 5}}, nil).Times(1)

		ordUsc := &OrderUsecase{
			orderRepo: storageMock,
		}
		claim := &entity.Users{
			Role: entity.UserRoleAdmin,
		}

		filter := &dto.OrderListFilter{}

		ordJs, err := ordUsc.GetOrders(ctx, claim, filter, 5, 0)
		req.NoError(err)
		req.Equal(uuid.UUID{3}, ordJs[0].Id)
	})

}
