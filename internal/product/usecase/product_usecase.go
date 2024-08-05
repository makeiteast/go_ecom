package usecase

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"go-store/internal/dto"
	"go-store/internal/entity"
	errorStatus "go-store/utils/errors"
)

// ProductUsecase will initiate usecase of entity.ProductRepository interface
type ProductUsecase struct {
	productRepo   entity.ProductRepository
	prodRedisRepo entity.ProdRedisRepository
	optionRepo    entity.OptionRepository
}

// NewProductUsecase will create new an ProductUsecase object representation of entity.ProductUsecase interface
func NewProductUsecase(p entity.ProductRepository, r entity.ProdRedisRepository, o entity.OptionRepository) entity.ProductUsecase {
	return &ProductUsecase{
		productRepo:   p,
		prodRedisRepo: r,
		optionRepo:    o,
	}
}

func (p *ProductUsecase) GetSku(ctx context.Context, limit int, offset int, userRole *entity.UserRole, filter *dto.ProductListFilter) (result *entity.ResultSkuJSon, err error) {
	skuResult := &entity.ResultSkuJSon{}
	var skus []*entity.Sku
	var prods []*entity.Product
	ctLog := log.WithFields(log.Fields{"func": "ProductUsecase.GetSkuProducts"})

	filterMap := map[string]string{}
	if filter != nil {
		filterMap = filter.ToSqlFilterMap()
	}

	skus, prods, err = p.productRepo.GetSku(ctx, limit, offset, filterMap)
	if err != nil {
		ctLog.WithFields(log.Fields{"method": "p.productRepo.GetSkuProducts"}).Warning(err)
		err = errorStatus.ErrInternalServer
		return
	}
	var skuResp = make([]*entity.SkuJson, len(skus))

	for idx, sku := range skus {
		skuValues, err := p.optionRepo.GetSkuValue(ctx, sku.Id, userRole)
		if err != nil {
			ctLog.WithFields(log.Fields{"method": "p.productRepo.GetSkuValue"}).Warning(err)
			err = errorStatus.ErrInternalServer
			return nil, err
		}
		skuResp[idx] = mapSkuToJSON(sku, prods[idx])
		if &skuValues != nil {
			SkuValueArray := make([]int32, 0)
			for _, skuValue := range skuValues {
				SkuValueArray = append(SkuValueArray, int32(skuValue.Id))
			}
			skuResp[idx].SkuValueId = SkuValueArray
		}
	}
	if len(skus) > 0 {
		skuResult.Total = skus[0].Total
	}

	skuResult.SkuJson = skuResp
	return skuResult, nil
}

func mapSkuToJSON(s *entity.Sku, p *entity.Product) *entity.SkuJson {
	return &entity.SkuJson{
		ProductName: p.ProductName,
		Description: p.Description,
		CategoryId:  p.CategoryId,
		CreateTs:    s.CreateTs,
		SkuId:       s.Id,
		SkuCode:     s.Sku,
		SkuPrice:    s.Price,
		SkuQuantity: s.Quantity,
		SkuImage:    s.SmallImage,
		CountViewed: s.CountViewed,
	}
}

func (p *ProductUsecase) GetProductSkus(ctx context.Context, limit int, offset int, categoryId int) (result *entity.ResultProductJSon, err error) {
	ctLog := log.WithFields(log.Fields{"func": "ProductUsecase.GetProduct"})

	productResult := &entity.ResultProductJSon{}
	var skus []*entity.Sku

	products, err := p.productRepo.GetProducts(ctx, limit, offset, categoryId)
	if err == pgx.ErrNoRows {
		ctLog.Warning(err)
		return nil, errorStatus.ErrNotFound
	}
	if err != nil {
		ctLog.Warning("Cannot GetProduct, err: \n", err)
		return nil, err
	}
	var productResp = make([]*entity.ProductJson, len(products))

	for idx, product := range products {

		productResp[idx] = mapProductToJSON(product)

		if err != nil {
			ctLog.Warning("Cannot map[] products, err: \n", err)
			return
		}
		skus, err = p.productRepo.GetSkuByProductID(ctx, product.Id)
		if err != nil {
			ctLog.Warning("Cannot GetProductSku, err: \n", err)
			return productResult, err
		}
		var skuResp = make([]*entity.ProductSkuJson, len(skus))

		for idx, sku := range skus {

			skuResp[idx] = mapProductSkuToJSON(sku)
		}

		productResp[idx].ProductSkuJson = skuResp

	}
	if len(products) > 0 {
		productResult.Total = products[0].Total
	}

	productResult.ProductJson = productResp
	return productResult, nil
}

func mapProductToJSON(s *entity.Product) *entity.ProductJson {
	return &entity.ProductJson{
		ProductId:   s.Id,
		ProductName: s.ProductName,
		Description: s.Description,
		CategoryId:  s.CategoryId,
		CreateTs:    s.CreateTs,
	}
}

func mapProductSkuToJSON(s *entity.Sku) *entity.ProductSkuJson {
	return &entity.ProductSkuJson{
		SkuId:       s.Id,
		SkuCode:     s.Sku,
		SkuPrice:    s.Price,
		SkuQuantity: s.Quantity,
		SkuImage:    s.SmallImage,
	}
}

func (p *ProductUsecase) GetSingleProduct(ctx context.Context, skuCode string, userRole *entity.UserRole) (result *entity.SkuJson, err error) {
	// var sku domain.Sku
	var skuResp *entity.SkuJson
	ctLog := log.WithFields(log.Fields{"func": "ProductUsecase.GetSingleProduct"})
	cachedProd, err := p.prodRedisRepo.GetProdByIDCtx(ctx, skuCode)
	if err == pgx.ErrNoRows {
		ctLog.Warning(err)
		return nil, errorStatus.ErrNotFound
	}
	if err != nil {
		ctLog.Warning(err)
		return nil, err
	}
	if cachedProd != nil {
		skuResp = cachedProd
		return skuResp, nil
	}
	sku, prod, err := p.productRepo.GetSingleProduct(ctx, skuCode, 0)

	if err != nil {
		ctLog.Warning("Cannot GetSingleProduct, err: \n", err)
		return nil, err
	}

	if (sku == &entity.Sku{}) {
		ctLog.Warning("sku is null!")
		err = errorStatus.ErrNotFound
		return nil, err
	}
	skuResp = mapSingleSkuToJSON(sku, prod)
	skuValues, err := p.optionRepo.GetSkuValue(ctx, sku.Id, userRole)
	if err != nil {
		ctLog.Warning("Cannot GetSkuValue, err: \n", err)
		return nil, err
	}
	SkuValueArray := make([]int32, 0)
	for _, skuValue := range skuValues {
		SkuValueArray = append(SkuValueArray, int32(skuValue.Id))
	}
	skuResp.SkuValueId = SkuValueArray

	err = p.prodRedisRepo.SetProdByIDCtx(ctx, skuCode, skuResp)
	if err != nil {
		ctLog.Warning(err)
		return nil, errorStatus.ErrInternalServer
	}

	return skuResp, nil
}

func mapSingleSkuToJSON(s *entity.Sku, p *entity.Product) *entity.SkuJson {
	return &entity.SkuJson{
		ProductName: p.ProductName,
		Description: p.Description,
		CategoryId:  p.CategoryId,
		CreateTs:    s.CreateTs,
		SkuId:       s.Id,
		SkuCode:     s.Sku,
		SkuPrice:    s.Price,
		SkuQuantity: s.Quantity,
		SkuImage:    s.LargeImage,
	}
}

func (p *ProductUsecase) CreateProduct(ctx context.Context, prod *entity.Product) (productId *int, err error) {
	ctLog := log.WithFields(log.Fields{"func": "ProductUsecase.CreateProduct"})
	prod.State = entity.Enabled
	prodID, err := p.productRepo.CreateProduct(ctx, prod)
	if err != nil {
		ctLog.WithError(err).Warning("p.productRepo.GetSingleProduct")
		return nil, err
	}
	return prodID, nil
}

func (p *ProductUsecase) UpdateProduct(ctx context.Context, prod *entity.Product) error {
	ctLog := log.WithFields(log.Fields{"func": "ProductUsecase.UpdateProduct"})

	prod.UpdateTs = time.Now()
	err := p.productRepo.UpdateProduct(ctx, prod)
	if err != nil {
		ctLog.WithError(err).Warning("p.productRepo.UpdateProduct")
		return err
	}

	return nil
}

func (p *ProductUsecase) DeleteProduct(ctx context.Context, productID int) error {
	ctLog := log.WithFields(log.Fields{"func": "ProductUsecase.CreateProduct"})

	err := p.productRepo.RemoveProduct(ctx, productID)
	if err == pgx.ErrNoRows {
		ctLog.Warning(err)
		return errorStatus.ErrNotFound
	}
	if err != nil {
		ctLog.WithError(err).Warning("p.productRepo.RemoveProduct")
		return err
	}
	return nil
}

func (p *ProductUsecase) CreateProductOption(ctx context.Context, optionForm *dto.ProductOptionRequest) error {
	ctLog := log.WithFields(log.Fields{"func": "ProductUsecase.CreateProduct"})

	skuValue := &entity.SkuValue{
		OptionId:      optionForm.OptionId,
		OptionValueId: optionForm.OptionValueId,
		CreateTs:      time.Now().UTC(),
		UpdateTs:      time.Now().UTC(),
		State:         entity.Enabled,
		Version:       0,
	}

	var txId int
	txId, err := p.optionRepo.NewTxId(ctx)
	if err != nil {
		ctLog.WithError(err).Warning("o.optionRepo.NewTxId")
		return err
	}
	defer func() {
		err = p.optionRepo.TxEnd(ctx, txId, err)
		if err != nil {
			ctLog.WithError(err).Warning("o.orderRepo.TxEnd")
			return
		}
	}()

	err = p.optionRepo.CreateSkuValue(ctx, optionForm.SkuId, skuValue, txId)
	if err != nil {
		ctLog.WithError(err).Warning("p.productRepo.GetSingleProduct")
		return err
	}

	return nil
}

func (p *ProductUsecase) DeleteProductOption(ctx context.Context, skuValueId int) error {
	ctLog := log.WithFields(log.Fields{"func": "ProductUsecase.DeleteProductOption"})

	err := p.optionRepo.RemoveSkuValue(ctx, skuValueId)
	if err != nil {
		ctLog.WithError(err).Warning("p.productRepo.RemoveSkuValue")
		return err
	}

	return nil
}

func (p *ProductUsecase) CreateSku(ctx context.Context, sku *entity.Sku) (err error) {
	ctLog := log.WithFields(log.Fields{"func": "ProductUsecase.CreateProduct"})
	err = p.productRepo.CreateSku(ctx, sku)
	if err != nil {
		ctLog.WithError(err).Warning("p.productRepo.GetSingleProduct")
		return err
	}
	return nil
}

func (p *ProductUsecase) UpdateSku(ctx context.Context, sku *entity.Sku) error {
	ctLog := log.WithFields(log.Fields{"func": "ProductUsecase.UpdateProduct"})

	sku.UpdateTs = time.Now()
	err := p.productRepo.UpdateSku(ctx, sku)
	if err != nil {
		ctLog.WithError(err).Warning("p.productRepo.UpdateProduct")
		return err
	}

	return nil
}

func (p *ProductUsecase) DeleteSku(ctx context.Context, skuId int) error {
	ctLog := log.WithFields(log.Fields{"func": "ProductUsecase.CreateProduct"})

	err := p.productRepo.RemoveSku(ctx, skuId)
	if err != nil {
		ctLog.WithError(err).Warning("p.productRepo.RemoveProduct")
		return err
	}
	return nil
}

func (p *ProductUsecase) GetSkuOption(ctx context.Context, skuValueId int, userRole *entity.UserRole) (result *entity.OptionJson, err error) {
	ctLog := log.WithFields(log.Fields{"func": "ProductUsecase.CreateProduct"})

	opt, err := p.optionRepo.GetOptionBySkuValue(ctx, skuValueId, userRole)
	if err == pgx.ErrNoRows {
		ctLog.Warning(err)
		return nil, errorStatus.ErrNotFound
	}
	if err != nil {
		ctLog.WithError(err).Warning("p.productRepo.RemoveProduct")
		return nil, err
	}

	return opt, nil
}
