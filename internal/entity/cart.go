package entity

import (
	"context"
	"time"
)

type Cart struct {
	UserId   int       `json:"userId"`
	SkuId    int       `json:"skuId"`
	Quantity int       `json:"quantity"`
	CreateTs time.Time `json:"createTs"`
	UpdateTs time.Time `json:"updateTs"`
	State    State     `json:"state"`
	Version  int       `json:"version"`
}

type CartJson struct {
	UserId   int `json:"userId"`
	SkuId    int `json:"skuId"`
	Quantity int `json:"quantity"`
}

type CartUsecase interface {
	GetCart(ctx context.Context, user *Users, limit int, offset int) (result []*CartJson, err error)
	CreateCart(ctx context.Context, user *Users, cart *Cart) (err error)
	UpdateCart(ctx context.Context, user *Users, cart *Cart) (err error)
	DeleteCart(ctx context.Context, cart *Cart) (err error)
}

type CartRepository interface {
	GetCarts(ctx context.Context, limit int, offset int, userId *int) (result []*Cart, err error)
	CreateCart(ctx context.Context, cart *Cart) (err error)
	UpdateCart(ctx context.Context, cart *Cart) (err error)
	DeleteCart(ctx context.Context, cart *Cart) (err error)
}

func (c *Cart) SetDefaults() {
	now := NowUTC()
	c.State = Enabled
	c.CreateTs = now
	c.UpdateTs = now
	c.Version = 0
}
