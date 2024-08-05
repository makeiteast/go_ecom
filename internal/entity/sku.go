package entity

import (
	"time"
)

type Sku struct {
	Id          int
	ProductId   int       `json:"productId"`
	Sku         string    `json:"skuName"`
	Price       float32   `json:"skuCode"`
	Quantity    int       `json:"quaintity"`
	LargeImage  string    `json:"largeImage"`
	SmallImage  string    `json:"smallImage"`
	ThumbImage  string    `json:"thumbImage"`
	CountViewed int       `json:"countViewed"`
	CreateTs    time.Time `json:"createTs"`
	UpdateTs    time.Time `json:"updateTs"`
	State       State     `json:"state"`
	Version     int       `json:"version"`
	Total       int       `json:"total"`
}

type SkuValue struct {
	Id            int       `json:"id"`
	SkuId         int       `json:"skuId"`
	OptionId      int       `json:"optionId"`
	OptionValueId int       `json:"optionValueId"`
	CreateTs      time.Time `json:"createTs"`
	UpdateTs      time.Time `json:"updateTs"`
	State         State     `json:"state"`
	Version       int       `json:"version"`
}

type SkuJson struct {
	ProductName string    `json:"name"`
	Description string    `json:"description"`
	CategoryId  int       `json:"categoryId"`
	CreateTs    time.Time `json:"createTs"`
	CountViewed int       `json:"countViewed"`
	SkuId       int       `json:"skuId"`
	SkuCode     string    `json:"skuCode"`
	SkuPrice    float32   `json:"skuPrice"`
	SkuQuantity int       `json:"skuQuantity"`
	SkuImage    string    `json:"skuImage"`
	SkuValueId  []int32   `json:"skuValueId"`
}

type ResultSkuJSon struct {
	Total   int        `json:"total"`
	SkuJson []*SkuJson `json:"products"`
}
