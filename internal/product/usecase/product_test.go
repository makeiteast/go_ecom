package usecase

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"go-store/internal/dto"
	"go-store/internal/entity"
	optionMocks "go-store/internal/option/mock"
	mocks "go-store/internal/product/mock"
)

func TestProduct(t *testing.T) {
	req := require.New(t)
	any := gomock.Any()
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	productMock := mocks.NewMockProductRepository(mockCtrl)
	optionMock := optionMocks.NewMockOptionRepository(mockCtrl)
	cachedProd := mocks.NewMockProdRedisRepository(mockCtrl)
	t.Run("get product sku(admin)  success", func(t *testing.T) {
		productMock.EXPECT().GetSku(ctx, any, any, any).Return([]*entity.Sku{{Id: 44}}, nil).Times(1)
		optionMock.EXPECT().GetSkuValue(ctx, 44).Return([]*entity.SkuValue{{Id: 5, SkuId: 44}}, nil).Times(1)
		optionMock.EXPECT().GetOption(ctx, 0).Return(&entity.Option{Id: 6, CategoryId: 3}, nil)
		optionMock.EXPECT().GetOptionValue(ctx, 0).Return(&entity.OptionValue{Id: 8, OptionId: 6}, nil)
		prodUsc := &ProductUsecase{
			productRepo:   productMock,
			optionRepo:    optionMock,
			prodRedisRepo: cachedProd,
		}
		filter := &dto.ProductListFilter{}
		skuJs, err := prodUsc.GetSku(ctx, 5, 0, filter)
		req.NoError(err)
		req.Equal(44, skuJs.SkuJson[0].SkuId)
	})

	t.Run("get sku product success", func(t *testing.T) {
		productMock.EXPECT().GetProducts(ctx, any, any, any).Return([]*entity.Product{{Id: 44}}, nil).Times(1)
		productMock.EXPECT().GetSkuByProductID(ctx, any).Return([]*entity.Sku{{Id: 55}}, nil).Times(1)
		prodUsc := &ProductUsecase{
			productRepo: productMock,
		}

		skuJs, err := prodUsc.GetProductSkus(ctx, 5, 0, 3)
		req.NoError(err)
		req.Equal(44, skuJs.ProductJson[0].ProductId)
		req.Equal(55, skuJs.ProductJson[0].ProductSkuJson[0].SkuId)
	})

	t.Run("get single product success", func(t *testing.T) {

		cachedProd.EXPECT().GetProdByIDCtx(ctx, any).Return(&entity.SkuJson{SkuId: 44}, nil)

		prodUsc := &ProductUsecase{
			productRepo:   productMock,
			prodRedisRepo: cachedProd,
		}

		skuJs, err := prodUsc.GetSingleProduct(ctx, "skuCode")
		req.NoError(err)
		req.Equal(44, skuJs.SkuId)
	})
}
