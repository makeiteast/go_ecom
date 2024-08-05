package dto

// CartListRequest -.
type CartCreateRequest struct {
	SkuId    int `json:"skuId"`
	Quantity int `json:"quantity"`
}
