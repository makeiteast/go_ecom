package http

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"go-store/internal/dto"
	"go-store/internal/entity"
	errorStatus "go-store/utils/errors"
	errorstatus "go-store/utils/errors"
	httphelper "go-store/utils/http"
)

// OderHandler  represent the httphandler for cart
type CartHandler struct {
	cartUc entity.CartUsecase
	user   entity.UserUsecase
	srvLog *logrus.Entry
}

const (
	sucsess = "sucsess"
)

// NewCartHandler will initialize the cart items endpoint
func NewCartHandler(handler *gin.RouterGroup, mdw gin.HandlerFunc, uc *entity.Usecases, srvLog *logrus.Entry) {
	oh := &CartHandler{
		cartUc: uc.CartUsecase,
		user:   uc.UserUsecase,
		srvLog: srvLog,
	}
	h := handler.Group("/cart")
	{
		h.POST("", mdw, oh.getCartItems)
		h.POST("/create", mdw, oh.createCartItem)
		h.PUT("/:cartId", mdw, oh.updateCartItem)
		h.DELETE("/:cartId", mdw, oh.deleteCartItem)
	}
}

func (oh *CartHandler) getCartItems(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "CartHandler.getCartItems"})
	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	pageParam, err := httphelper.PaginationParams(c)
	if err != nil {
		srvLog.WithError(err).Warning(err)
		httphelper.SendResponse(c, nil, errorstatus.ErrBadReq)
		return
	}

	srvLog.WithFields(log.Fields{"limit": pageParam.Limit, "offset": pageParam.Offset, "func": "server.getCartItems"})

	result, err := oh.cartUc.GetCart(c, user, pageParam.Limit, pageParam.Offset)
	if err != nil {
		srvLog.WithError(err).Warning("oh.ordUsecase.GetOrders")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, result, nil)
}

func (oh *CartHandler) createCartItem(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "CartHandler.createCartItem"})
	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	var createReq dto.CartCreateRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&createReq); err != nil {
		srvLog.WithError(err).Error("format is wrong")
		httphelper.SendResponse(c, nil, errorstatus.ErrBadReq)
		return
	}

	cart := &entity.Cart{
		UserId:   user.Id,
		SkuId:    createReq.SkuId,
		Quantity: createReq.Quantity,
	}
	cart.SetDefaults()

	err := oh.cartUc.CreateCart(c, user, cart)
	if err != nil {
		oh.srvLog.WithError(err).Warning("oh.ordUsecase.CreateOrder")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)

}

func (oh *CartHandler) updateCartItem(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "CartHandler.updateOrder"})
	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	cart := &entity.Cart{}

	err := oh.cartUc.UpdateCart(c, user, cart)
	if err != nil {
		srvLog.WithError(err).Warning("oh.ordUsecase.UpdateOrder")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)

}

func (oh *CartHandler) deleteCartItem(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "CartHandler.deleteCartItem"})
	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	skuId, ok := c.GetQuery("skuId")
	if !ok {
		srvLog.Warning("createCategoryOpt.GetQuery.optionName")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	skuIdInt, err := strconv.Atoi(skuId)
	if err != nil {
		logrus.Warning("strconv.Atoi.skuId")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	cart := &entity.Cart{
		UserId: user.Id,
		SkuId:  skuIdInt,
	}

	err = oh.cartUc.DeleteCart(c, cart)
	if err != nil {
		oh.srvLog.WithError(err).Warning("oh.ordUsecase.DeleteCart")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}
