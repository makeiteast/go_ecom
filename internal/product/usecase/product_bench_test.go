package usecase

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"

	"go-store/internal/dto"
	"go-store/internal/entity"
	mocks "go-store/internal/product/mock"
)

func BenchmarkOrder(b *testing.B) {
	any := gomock.Any()
	ctx := context.Background()
	mockCtrl := gomock.NewController(b)
	defer mockCtrl.Finish()

	storageMock := mocks.NewMockProductRepository(mockCtrl)

	b.Run("get product success", func(b *testing.B) {
		storageMock.EXPECT().GetSku(ctx, any, any, any).Return([]entity.Sku{{Id: 44}}, nil).Times(1)

		prodUsc := &ProductUsecase{
			productRepo: storageMock,
		}
		filter := &dto.ProductListFilter{}
		prodUsc.GetSku(ctx, 5, 0, filter)
	})
}
