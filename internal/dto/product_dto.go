package dto

import (
	"strconv"
)

// ProductListRequest -.
type ProductListRequest struct {
	Filter *ProductListFilter `json:"filter"`
	Limit  int                `json:"limit" example:"10"`
	Offset int                `json:"offset" example:"0"`
}

type ProductListFilter struct {
	ProductName string   `json:"productName,omitempty"`
	Description string   `json:"description,omitempty"`
	CategoryId  *int     `json:"categoryId,omitempty"`
	BrandId     *int     `json:"brandId,omitempty"`
	RegionId    *int     `json:"regionId,omitempty"`
	PriceStart  *float32 `json:"priceStart,omitempty"`
	PriceEnd    *float32 `json:"priceEnd,omitempty"`
}

type ProductOptionRequest struct {
	SkuId         string `json:"skuId"`
	OptionId      int    `json:"optionId"`
	OptionValueId int    `json:"optionValueId"`
}

func (p *ProductListFilter) ToSqlFilterMap() map[string]string {
	filterMap := map[string]string{}
	if len(p.ProductName) > 0 {
		filterMap["product_name"] = p.ProductName
	}
	if len(p.Description) > 0 {
		filterMap["description"] = p.Description
	}
	if p.CategoryId != nil {
		filterMap["category_id"] = strconv.Itoa(*p.CategoryId)
	}
	if p.BrandId != nil {
		filterMap["brand_id"] = strconv.Itoa(*p.BrandId)
	}
	if p.RegionId != nil {
		filterMap["region_id"] = strconv.Itoa(*p.RegionId)
	}
	if p.PriceStart != nil {
		float64Price := float64(*p.PriceStart)
		filterMap["priceStart"] = strconv.FormatFloat(float64Price, 'f', -1, 32)
	}
	if p.PriceEnd != nil {
		float64Price := float64(*p.PriceEnd)
		filterMap["priceEnd"] = strconv.FormatFloat(float64Price, 'f', -1, 32)
	}
	return filterMap
}
