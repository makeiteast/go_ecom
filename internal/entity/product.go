package entity

import (
	"context"
	"go-store/internal/dto"
	"time"
)

type Product struct {
	Id          int
	ProductName string    `json:"productName"`
	Description string    `json:"description"`
	CategoryId  int       `json:"categoryId"`
	BrandId     int       `json:"brandId"`
	RegionId    int       `json:"regionId"`
	CreateTs    time.Time `json:"createTs"`
	UpdateTs    time.Time `json:"updateTs"`
	State       State     `json:"state"`
	Version     int       `json:"version"`
	Total       int       `json:"total"`
	Users       `db:"users"`
}

type ResultProductJSon struct {
	Total       int            `json:"total"`
	ProductJson []*ProductJson `json:"products"`
}

type ProductJson struct {
	ProductId      int               `json:"id"`
	ProductName    string            `json:"name"`
	Description    string            `json:"description"`
	CategoryId     int               `json:"categoryId"`
	CreateTs       time.Time         `json:"createTs"`
	ProductSkuJson []*ProductSkuJson `json:"variants"`
}

type ProductSkuJson struct {
	SkuId       int     `json:"skuId"`
	SkuCode     string  `json:"skuCode"`
	SkuPrice    float32 `json:"skuPrice"`
	SkuQuantity int     `json:"skuQuantity"`
	SkuImage    string  `json:"skuImage"`
}

type ProductUsecase interface {
	GetSku(ctx context.Context, limit int, offset int, userRole *UserRole, filter *dto.ProductListFilter) (result *ResultSkuJSon, err error)
	GetProductSkus(ctx context.Context, limit int, offset int, categoryID int) (result *ResultProductJSon, err error)
	GetSingleProduct(ctx context.Context, skuCode string, userRole *UserRole) (result *SkuJson, err error)
	CreateProduct(ctx context.Context, prod *Product) (productId *int, err error)
	UpdateProduct(ctx context.Context, prod *Product) (err error)
	DeleteProduct(ctx context.Context, productID int) (err error)
	CreateSku(ctx context.Context, prod *Sku) (err error)
	UpdateSku(ctx context.Context, prod *Sku) (err error)
	DeleteSku(ctx context.Context, skuId int) (err error)
	CreateProductOption(ctx context.Context, optionForm *dto.ProductOptionRequest) (err error)
	DeleteProductOption(ctx context.Context, skuValueId int) (err error)
	GetSkuOption(ctx context.Context, skuValueId int, userRole *UserRole) (result *OptionJson, err error)
}

type ProductRepository interface {
	GetSku(ctx context.Context, limit int, offset int, filter map[string]string) (result []*Sku, products []*Product, err error)
	GetProducts(ctx context.Context, limit int, offset int, categoryID int) (result []*Product, err error)
	GetSkuByProductID(ctx context.Context, productID int) (result []*Sku, err error)
	GetSingleProduct(ctx context.Context, skuCode string, skuId int) (result *Sku, product *Product, err error)
	CreateProduct(ctx context.Context, prod *Product) (prodID *int, err error)
	UpdateProduct(ctx context.Context, prod *Product) error
	RemoveProduct(ctx context.Context, productID int) error
	CreateSku(ctx context.Context, sku *Sku) (err error)
	UpdateSku(ctx context.Context, sku *Sku) error
	RemoveSku(ctx context.Context, skuId int) error
}

type ProdRedisRepository interface {
	GetProducts(ctx context.Context, limit int, offset int, category int) (prod []*ProductJson, err error)
	SetProdCtx(ctx context.Context, offset int, category int, prod *ProductJson) error
	SetSkuCount(ctx context.Context, count string, key string) error
	GetSkuCount(ctx context.Context, key string) (string, error)
	GetSku(ctx context.Context, limit int, offset int, category int) ([]*SkuJson, error)
	SetSkuCtx(ctx context.Context, offset int, category int, prod *SkuJson) error
	GetProdByIDCtx(ctx context.Context, key string) (*SkuJson, error)
	SetProdByIDCtx(ctx context.Context, key string, user *SkuJson) error
}
