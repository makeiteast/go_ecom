package usecase

import (
	"context"

	log "github.com/sirupsen/logrus"

	"go-store/internal/entity"
	errorStatus "go-store/utils/errors"
)

// CartUsecase will initiate usecase of entity.CartRepository interface
type CartUsecase struct {
	cartRepo entity.CartRepository
	prodRepo entity.ProductRepository
}

// NewCartUsecase will create new an CartUsecase object representation of entity.CartUsecase interface
func NewCartUsecase(cart entity.CartRepository, prod entity.ProductRepository) entity.CartUsecase {
	return &CartUsecase{
		cartRepo: cart,
		prodRepo: prod,
	}
}

func (o *CartUsecase) GetCart(ctx context.Context, user *entity.Users, limit int, offset int) (result []*entity.CartJson, err error) {
	srvLog := log.WithFields(log.Fields{"func": "CartUsecase.GetOrders"})

	var userId *int

	//if user not admin then find only current users cart items
	if user.Role != entity.UserRoleAdmin {
		userId = &user.Id
	}

	carts, err := o.cartRepo.GetCarts(ctx, limit, offset, userId)
	if err != nil {
		srvLog.Warning("Cannot get cart query, Err: ", err)
		err = errorStatus.ErrInternalServer
		return
	}

	cartResp := make([]*entity.CartJson, len(carts))

	for idx, cart := range carts {

		cartResp[idx] = mapCartToJSON(cart)

	}
	result = cartResp
	return
}

func mapCartToJSON(c *entity.Cart) *entity.CartJson {
	return &entity.CartJson{
		UserId:   c.UserId,
		SkuId:    c.SkuId,
		Quantity: c.Quantity,
	}
}

func (o *CartUsecase) CreateCart(ctx context.Context, user *entity.Users, cart *entity.Cart) error {
	srvLog := log.WithFields(log.Fields{"func": "CartUsecase.GetOrderById"})

	err := o.cartRepo.CreateCart(ctx, cart)
	if err != nil {
		srvLog.WithError(err).Warning("o.cartRepo.CreateOrderItem")
		return err
	}

	return nil
}

func (o *CartUsecase) UpdateCart(ctx context.Context, user *entity.Users, order *entity.Cart) error {
	srvLog := log.WithFields(log.Fields{"func": "CartUsecase.UpdateOrder"})

	err := o.cartRepo.UpdateCart(ctx, order)
	if err != nil {
		srvLog.WithError(err).Warning("o.cartRepo.UpdateOrder")
		return err
	}

	return nil
}

func (o *CartUsecase) DeleteCart(ctx context.Context, cart *entity.Cart) error {
	srvLog := log.WithFields(log.Fields{"func": "CartUsecase.DeleteOrder"})

	err := o.cartRepo.DeleteCart(ctx, cart)
	if err != nil {
		srvLog.WithError(err).Warning("o.cartRepo.DeleteOrder")
		return err
	}

	return nil
}
