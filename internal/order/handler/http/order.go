package http

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	dto "go-store/internal/dto"
	"go-store/internal/entity"
	errorstatus "go-store/utils/errors"
	httphelper "go-store/utils/http"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// OderHandler  represent the httphandler for order
type OrderHandler struct {
	ordUsecase entity.OrderUsecase
	user       entity.UserUsecase
	srvLog     *logrus.Entry
}

const (
	sucsess = "sucsess"
)

// NewOrderHandler will initialize the orders/ resources endpoint
func NewOrderHandler(handler *gin.RouterGroup, mdw gin.HandlerFunc, uc *entity.Usecases, srvLog *logrus.Entry) {
	oh := &OrderHandler{
		ordUsecase: uc.OrderUsecase,
		user:       uc.UserUsecase,
		srvLog:     srvLog,
	}
	h := handler.Group("/order")
	{
		h.POST("/list", mdw, oh.getOrders)
		h.POST("/:orderId", mdw, oh.getOrder)
		h.POST("", mdw, oh.createOrder)
		h.PUT("/:orderId", mdw, oh.updateOrder)
		h.DELETE("/:orderId", mdw, oh.deleteOrder)
	}
	a := handler.Group("/admin")
	{
		a.PUT("/order/:orderId", mdw, oh.updateOrderStatus)
	}
}

func (oh *OrderHandler) getOrders(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "OrderHandler.getOrders"})

	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	if user.Role == entity.UserRoleGuest {

	}

	var listReq dto.OrderListRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&listReq); err != nil {
		srvLog.WithError(err).Error("format is wrong")
		httphelper.SendResponse(c, nil, errorstatus.ErrBadReq)
		return
	}

	srvLog.WithFields(log.Fields{"limit": listReq.Limit, "offset": listReq.Offset, "func": "server.GetOrdersHandler"})

	result, err := oh.ordUsecase.GetOrders(c, user, listReq.Filter, listReq.Limit, listReq.Offset)
	if err != nil {
		srvLog.WithError(err).Warning("oh.ordUsecase.GetOrders")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, result, nil)
}

func (oh *OrderHandler) getOrder(c *gin.Context) {
	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	orderId := c.Param("orderId")

	result, err := oh.ordUsecase.GetOrderById(c, user, orderId)
	if err != nil {
		oh.srvLog.WithError(err).Warning("oh.ordUsecase.GetOrderById")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, result, nil)
}

func (oh *OrderHandler) createOrder(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "OrderHandler.getOrder"})
	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	order, err := httphelper.OrderForm(c)
	if err != nil {
		srvLog.WithError(err).Warning("createOrder.httphelper.OrderForm")
		httphelper.SendResponse(c, nil, err)
		return
	}

	err = oh.ordUsecase.CreateOrder(c, user, order)
	if err != nil {
		oh.srvLog.WithError(err).Warning("oh.ordUsecase.CreateOrder")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)

}

func (oh *OrderHandler) updateOrder(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "OrderHandler.updateOrder"})

	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	orderId := c.Param("orderId")
	orderIduuid := uuid.Must(uuid.Parse(orderId))
	srvLog.Info(orderIduuid)
	order, err := httphelper.OrderForm(c)
	if err != nil {
		srvLog.WithError(err).Warning("updateOrder.httphelper.OrderForm")
		httphelper.SendResponse(c, nil, err)
		return
	}
	order.Id = orderIduuid
	order.UserId = user.Id
	order.UpdateTs = time.Now()

	err = oh.ordUsecase.UpdateOrder(c, user, order)
	if err != nil {
		oh.srvLog.WithError(err).Warning("oh.ordUsecase.UpdateOrder")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)

}

func (oh *OrderHandler) updateOrderStatus(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "OrderHandler.updateOrderStatus"})
	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	orderId := c.Param("orderId")
	orderIduuid := uuid.Must(uuid.Parse(orderId))

	order, err := httphelper.OrderForm(c)
	if err != nil {
		srvLog.WithError(err).Warning("updateOrder.httphelper.OrderForm")
		httphelper.SendResponse(c, nil, err)
		return
	}
	order.Id = orderIduuid
	order.UserId = user.Id
	order.UpdateTs = time.Now()

	err = oh.ordUsecase.UpdateOrder(c, user, order)
	if err != nil {
		oh.srvLog.WithError(err).Warning("oh.ordUsecase.UpdateOrder")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)

}

func (oh *OrderHandler) deleteOrder(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "OrderHandler.deleteOrder"})
	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	orderId := c.Param("orderId")
	orderIduuid := uuid.Must(uuid.Parse(orderId))

	order := &entity.Order{
		Id:       orderIduuid,
		UserId:   user.Id,
		UpdateTs: time.Now(),
	}

	err := oh.ordUsecase.DeleteOrder(c, order)
	if err != nil {
		srvLog.WithError(err).Warning("oh.ordUsecase.DeleteOrder")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}
