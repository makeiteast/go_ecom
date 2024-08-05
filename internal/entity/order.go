package entity

import (
	"context"
	"go-store/internal/dto"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id        uuid.UUID `json:"id"`
	UserId    int       `json:"userId"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Comment   string    `json:"comment"`
	Notes     string    `json:"notes"`
	Status    string    `json:"status"`
	CreateTs  time.Time `json:"createTs"`
	UpdateTs  time.Time `json:"updateTs"`
	State     State     `json:"state"`
	Version   int       `json:"version"`
	OrderItem OrderItem `db:"order_item"`
}

type OrderItem struct {
	Id       int       `json:"id"`
	OrderId  uuid.UUID `json:"orderId"`
	SkuId    int       `json:"skuId"`
	Quantity int       `json:"quantity"`
	Price    float32   `json:"price"`
	CreateTs time.Time `json:"createTs"`
	UpdateTs time.Time `json:"updateTs"`
	State    State     `json:"state"`
	Version  int       `json:"version"`
	Sku      `db:"sku"`
}

type OrderItemSkuJson struct {
	ItemId    int     `json:"itemId"`
	SkuName   string  `json:"sku"`
	Quantity  int     `json:"quantity"`
	Price     float32 `json:"price"`
	Subtotal  float32 `json:"subTotal"`
	SmallName string  `json:"skuImage"`
}

type OrderJson struct {
	Id        uuid.UUID          `json:"id"`
	UserId    int                `json:"user_id"`
	Address   string             `json:"address"`
	Phone     string             `json:"phone"`
	Comment   string             `json:"comment"`
	Status    string             `json:"status"`
	CreateTs  time.Time          `json:"createTs"`
	OrderItem []OrderItemSkuJson `json:"orderItems"`
	TotalSum  float32            `json:"totalSum"`
}

func (o *Order) SetDefaults() {
	o.State = Enabled
	o.CreateTs = NowUTC()
	o.UpdateTs = NowUTC()
	o.Version = 0
}

func (o *OrderItem) SetDefaults() {
	o.State = Enabled
	o.CreateTs = NowUTC()
	o.UpdateTs = NowUTC()
	o.Version = 0
}

type OrderUsecase interface {
	GetOrders(ctx context.Context, user *Users, filter *dto.OrderListFilter, limit int, offset int) (result []*OrderJson, err error)
	GetOrderById(ctx context.Context, user *Users, orderId string) (result *OrderJson, err error)
	CreateOrder(ctx context.Context, user *Users, order *Order) (err error)
	UpdateOrder(ctx context.Context, user *Users, order *Order) (err error)
	UpdateOrderStatus(ctx context.Context, user *Users, order *Order) (err error)
	DeleteOrder(ctx context.Context, order *Order) (err error)
}

type OrderRepository interface {
	GetOrders(ctx context.Context, filterMap map[string]string, limit int, offset int, id int) (result []*Order, err error)
	GetOrder(ctx context.Context, orderId string) (result *Order, err error)
	GetItem(ctx context.Context, orderId uuid.UUID) (result []*OrderItem, err error)
	CreateOrder(ctx context.Context, order *Order, txId int) (result *uuid.UUID, err error)
	CreateOrderItem(ctx context.Context, orderId uuid.UUID, userId int, txId int) (err error)
	UpdateOrder(ctx context.Context, order *Order) (err error)
	UpdateOrderStatus(ctx context.Context, order *Order) (err error)
	DeleteOrder(ctx context.Context, order *Order) (err error)
	NewTxId(ctx context.Context) (txId int, err error)
	TxEnd(ctx context.Context, txId int, err error) error
}

type OrderRedisRepository interface {
	GetOrderSuperadminQuery(ctx context.Context, limit int, offset int, dbSchema string) (result []*Order, err error)
	GetOrderQuery(ctx context.Context, limit int, offset int, id int, dbSchema string) (result []*Order, err error)
	GetItem(ctx context.Context, order_id uuid.UUID) (result []*OrderItem, err error)
}
