package usecase

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"

	"go-store/internal/dto"
	"go-store/internal/entity"
	errorStatus "go-store/utils/errors"
)

// OrderUsecase will initiate usecase of entity.OrderRepository interface
type OrderUsecase struct {
	orderRepo entity.OrderRepository
}

// NewOrderUsecase will create new an OrderUsecase object representation of entity.OrderUsecase interface
func NewOrderUsecase(o entity.OrderRepository) entity.OrderUsecase {
	return &OrderUsecase{
		orderRepo: o,
	}
}

func (o *OrderUsecase) GetOrders(ctx context.Context, user *entity.Users, filter *dto.OrderListFilter, limit int, offset int) (result []*entity.OrderJson, err error) {
	var orders []*entity.Order
	srvLog := log.WithFields(log.Fields{"func": "OrderUsecase.GetOrders"})

	filterMap := map[string]string{}
	if filter != nil {
		filterMap = filter.ToSqlFilterMap()
	}
	// if not admin then filter only current users order
	if user.Role != entity.UserRoleAdmin {
		filterMap["user_id"] = strconv.Itoa(user.Id)
	}
	orders, err = o.orderRepo.GetOrders(ctx, filterMap, limit, offset, user.Id)
	if err != nil {
		srvLog.Warning("Cannot get order query, Err: ", err)
		err = errorStatus.ErrInternalServer
		return
	}

	orderResp := make([]*entity.OrderJson, len(orders))

	for idx, order := range orders {
		orderItems, err := o.orderRepo.GetItem(ctx, order.Id)
		if err != nil {
			srvLog.Warning("Cannot Get OrderItem, Err: ", err)
			err = errorStatus.ErrInternalServer
			return nil, err
		}

		var totalSum, subTotal float32
		orderResp[idx] = mapOrderToJSON(order)

		OrderItemList := make([]entity.OrderItemSkuJson, len(orderItems))
		for idx, orderItem := range orderItems {
			subTotal = float32(orderItem.Quantity) * orderItem.Price
			OrderItemList[idx] = mapOrderItemToJSON(orderItem)
			OrderItemList[idx].Subtotal = subTotal
			totalSum = subTotal + totalSum
		}
		orderResp[idx].OrderItem = OrderItemList
		orderResp[idx].TotalSum = totalSum
	}
	result = orderResp
	return
}

func mapOrderToJSON(s *entity.Order) *entity.OrderJson {
	return &entity.OrderJson{
		Id:       s.Id,
		UserId:   s.UserId,
		Address:  s.Address,
		Phone:    s.Phone,
		Comment:  s.Comment,
		Status:   s.Status,
		CreateTs: s.CreateTs,
	}
}

func mapOrderItemToJSON(s *entity.OrderItem) entity.OrderItemSkuJson {
	return entity.OrderItemSkuJson{
		ItemId:    s.Id,
		SkuName:   s.Sku.Sku,
		Quantity:  s.Quantity,
		Price:     s.Price,
		SmallName: s.Sku.SmallImage,
	}
}

func (o *OrderUsecase) GetOrderById(ctx context.Context, user *entity.Users, orderId string) (result *entity.OrderJson, err error) {

	srvLog := log.WithFields(log.Fields{"func": "OrderUsecase.GetOrderById"})

	order, err := o.orderRepo.GetOrder(ctx, orderId)
	if err != nil {
		srvLog.WithError(err).Warning("o.orderRepo.GetOrder")
		return nil, err
	}

	if !((user.Role == entity.UserRoleAdmin) || (user.Id == order.UserId)) {
		return nil, errorStatus.ErrAuth
	}

	orderItems, err := o.orderRepo.GetItem(ctx, order.Id)
	if err != nil {
		srvLog.WithError(err).Warning("o.orderRepo.GetItem")
		err = errorStatus.ErrInternalServer
		return nil, err
	}

	var totalSum, subTotal float32
	orderResp := mapOrderToJSON(order)

	OrderItemList := make([]entity.OrderItemSkuJson, len(orderItems))
	for idx, orderItem := range orderItems {
		subTotal = float32(orderItem.Quantity) * orderItem.Price
		OrderItemList[idx] = mapOrderItemToJSON(orderItem)
		OrderItemList[idx].Subtotal = subTotal
		totalSum = subTotal + totalSum
	}
	orderResp.OrderItem = OrderItemList
	orderResp.TotalSum = totalSum

	return orderResp, nil
}

func (o *OrderUsecase) CreateOrder(ctx context.Context, user *entity.Users, order *entity.Order) error {
	srvLog := log.WithFields(log.Fields{"func": "OrderUsecase.CreateOrder"})

	order.UserId = user.Id
	order.SetDefaults()

	//rollback
	var txId int
	txId, err := o.orderRepo.NewTxId(ctx)
	if err != nil {
		srvLog.WithError(err).Error("OrderUsecase - error processing o.orderRepo.NewTxId")
		return err
	}
	defer func() {
		err = o.orderRepo.TxEnd(ctx, txId, err)
		if err != nil {
			srvLog.WithError(err).Error("OrderUsecase - error processing o.orderRepo.TxEnd")
			return
		}
	}()

	orderId, err := o.orderRepo.CreateOrder(ctx, order, txId)
	if err != nil {
		srvLog.WithError(err).Warning("o.orderRepo.CreateOrder")
		return err
	}

	err = o.orderRepo.CreateOrderItem(ctx, *orderId, order.UserId, txId)
	if err != nil {
		srvLog.WithError(err).Warning("o.orderRepo.CreateOrderItem")
		return err
	}

	return nil
}

func (o *OrderUsecase) UpdateOrder(ctx context.Context, user *entity.Users, order *entity.Order) error {
	srvLog := log.WithFields(log.Fields{"func": "OrderUsecase.UpdateOrder"})

	err := o.orderRepo.UpdateOrder(ctx, order)
	if err != nil {
		srvLog.WithError(err).Warning("o.orderRepo.UpdateOrder")
		return err
	}

	return nil
}

func (o *OrderUsecase) UpdateOrderStatus(ctx context.Context, user *entity.Users, order *entity.Order) error {
	srvLog := log.WithFields(log.Fields{"func": "OrderUsecase.UpdateOrderStatus"})

	if user.Role != "ADMIN" {
		return errorStatus.ErrAuth
	}

	err := o.orderRepo.UpdateOrderStatus(ctx, order)
	if err != nil {
		srvLog.WithError(err).Warning("o.orderRepo.UpdateOrderStatus")
		return err
	}

	return nil
}

func (o *OrderUsecase) DeleteOrder(ctx context.Context, order *entity.Order) error {
	srvLog := log.WithFields(log.Fields{"func": "OrderUsecase.DeleteOrder"})

	err := o.orderRepo.DeleteOrder(ctx, order)
	if err != nil {
		srvLog.WithError(err).Warning("o.orderRepo.DeleteOrder")
		return err
	}

	return nil
}
