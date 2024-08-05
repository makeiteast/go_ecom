package http

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	dto "go-store/internal/dto"
	"go-store/internal/entity"
	errorstatus "go-store/utils/errors"
	httphelper "go-store/utils/http"
)

// ProductHandler  represent the httphandler for article
type ProductHandler struct {
	PrUsecase entity.ProductUsecase
	user      entity.UserUsecase
	srvLog    *logrus.Entry
}

const (
	sucsess = "sucsess"
)

// NewProductHandler will initialize the articles/ resources endpoint
func NewProductHandler(handler *gin.RouterGroup, mdw gin.HandlerFunc, uc *entity.Usecases, srvLog *logrus.Entry) {
	ph := &ProductHandler{
		PrUsecase: uc.ProductUsecase,
		user:      uc.UserUsecase,
		srvLog:    srvLog,
	}
	h := handler.Group("/product")
	{
		h.POST("", mdw, ph.products)
		h.POST("/:skuCode", mdw, ph.singleProduct)
		h.POST("/skuvalue/:skuValueId", mdw, ph.optionBySkuValue)
	}
	a := handler.Group("/admin")
	{
		a.POST("/products", mdw, ph.productBySkus)
		a.POST("/product", mdw, ph.createProduct)
		a.PUT("/product/:productId", mdw, ph.updateProduct)
		a.DELETE("/product/:productId", mdw, ph.deleteProduct)
		a.POST("/sku/option/:sku", mdw, ph.createProductOption) //add id's of option and value.
		a.DELETE("/sku/option/:skuValueId", mdw, ph.deleteProductOption)
		a.POST("/product/:productId/sku", mdw, ph.createSku)
		a.PUT("/sku/:sku", mdw, ph.updateSku)
		a.DELETE("/sku/:skuId", mdw, ph.deleteSku)
	}

}

func (ph *ProductHandler) products(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "products"})

	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)
	var listReq dto.ProductListRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&listReq); err != nil {
		srvLog.WithError(err).Error("format is wrong")
		httphelper.SendResponse(c, nil, errorstatus.ErrBadReq)
		return
	}
	srvLog.WithFields(log.Fields{"limit": listReq.Limit, "offset": listReq.Offset, "filter": listReq.Filter})

	result, err := ph.PrUsecase.GetSku(c, listReq.Limit, listReq.Offset, &user.Role, listReq.Filter)
	if err != nil {
		srvLog.Warning("Cannot GetSkuProduct, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, result, nil)
}

func (ph *ProductHandler) productBySkus(c *gin.Context) {
	var err error

	srvLog := log.WithFields(log.Fields{"func": "productBySkus"})

	pageParam, err := httphelper.SortingParams(c)
	if err != nil {
		srvLog.Warning(err)
		httphelper.SendResponse(c, nil, errorstatus.ErrBadReq)
		return
	}

	srvLog.WithFields(log.Fields{"limit": pageParam.Limit, "offset": pageParam.Offset, "categoryId": pageParam.Category})

	result, err := ph.PrUsecase.GetProductSkus(c, pageParam.Limit, pageParam.Offset, pageParam.Category)

	if err != nil {
		srvLog.Warning("Cannot get GetProduct, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, result, nil)
}

func (ph *ProductHandler) singleProduct(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "GetSingleProductHandler"})

	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	skuCode := c.Param("skuCode")

	result, err := ph.PrUsecase.GetSingleProduct(c, skuCode, &user.Role)
	if err != nil {
		srvLog.Warning(err)
	}

	if err != nil {
		srvLog.Warning("Cannot get product, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, result, nil)
}

func (ph *ProductHandler) createProduct(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "CreateProduct"})

	// exist := httphelper.CheckIfAdmin(c)
	// if !exist {
	// 	err := errorstatus.ErrAuth
	// 	srvLog.WithError(err).Warning("Not ADMIN role. Access denied")
	// 	httphelper.SendResponse(c, nil, err)
	// 	return
	// }

	prodForm, err := httphelper.ProdCreateForm(c)
	if err != nil {
		srvLog.WithError(err).Warning("httphelper.ProdCreateForm")
		httphelper.SendResponse(c, nil, err)
		return
	}
	productId, err := ph.PrUsecase.CreateProduct(c, prodForm)
	if err != nil {
		srvLog.WithError(err).Warning("ph.PrUsecase.CreateProduct")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, productId, nil)
}

func (ph *ProductHandler) updateProduct(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "ProductRemove"})

	// exist := httphelper.CheckIfAdmin(c)
	// if !exist {
	// 	err := errorstatus.ErrAuth
	// 	srvLog.WithError(err).Warning("Not ADMIN role. Access denied")
	// 	httphelper.SendResponse(c, nil, err)
	// 	return
	// }

	productIdStr := c.Param("productId")

	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		logrus.WithError(err).Warning("utils.DeleteProduct.productId")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	prodForm, err := httphelper.ProdCreateForm(c)
	if err != nil {
		srvLog.WithError(err).Warning("httphelper.ProdCreateForm")
		httphelper.SendResponse(c, nil, err)
		return
	}

	prodForm.Id = productId
	err = ph.PrUsecase.UpdateProduct(c, prodForm)
	if err != nil {
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) deleteProduct(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "ProductRemove"})
	// exist := httphelper.CheckIfAdmin(c)
	// if !exist {
	// 	err := errorstatus.ErrAuth
	// 	srvLog.WithError(err).Warning("Not ADMIN role. Access denied")
	// 	httphelper.SendResponse(c, nil, err)
	// 	return
	// }

	productIdStr := c.Param("productId")

	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		srvLog.WithError(err).Warning("utils.DeleteProduct.productId")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	err = ph.PrUsecase.DeleteProduct(c, productId)
	if err != nil {
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) createProductOption(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "CreateProductOption"})
	// exist := httphelper.CheckIfAdmin(c)
	// if !exist {
	// 	err := errorstatus.ErrAuth
	// 	srvLog.WithError(err).Warning("Not ADMIN role. Access denied")
	// 	httphelper.SendResponse(c, nil, err)
	// 	return
	// }

	optionForm, err := httphelper.OptionForm(c)

	err = ph.PrUsecase.CreateProductOption(c, optionForm)
	if err != nil {
		srvLog.WithError(err).Warning("ph.PrUsecase.CreateProduct")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) deleteProductOption(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "ProductRemove"})
	// exist := httphelper.CheckIfAdmin(c)
	// if !exist {
	// 	err := errorstatus.ErrAuth
	// 	srvLog.WithError(err).Warning("Not ADMIN role. Access denied")
	// 	httphelper.SendResponse(c, nil, err)
	// 	return
	// }

	skuValueIdStr := c.Param("skuValueId")

	skuValueId, err := strconv.Atoi(skuValueIdStr)
	if err != nil {
		srvLog.WithError(err).Warning("utils.DeleteProductOption.skuValueIdStr")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	err = ph.PrUsecase.DeleteProductOption(c, skuValueId)
	if err != nil {
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) createSku(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "CreateProduct"})
	// exist := httphelper.CheckIfAdmin(c)
	// if !exist {
	// 	err := errorstatus.ErrAuth
	// 	srvLog.WithError(err).Warning("Not ADMIN role. Access denied")
	// 	httphelper.SendResponse(c, nil, err)
	// 	return
	// }

	skuForm, err := httphelper.SkuCreateForm(c)
	if err != nil {
		srvLog.WithError(err).Warning("httphelper.ProdCreateForm")
		httphelper.SendResponse(c, nil, err)
		return
	}

	err = ph.PrUsecase.CreateSku(c, skuForm)
	if err != nil {
		srvLog.WithError(err).Warning("ph.PrUsecase.CreateProduct")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, skuForm.Sku, nil)
}

func (ph *ProductHandler) updateSku(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "CreateProduct"})
	// exist := httphelper.CheckIfAdmin(c)
	// if !exist {
	// 	err := errorstatus.ErrAuth
	// 	srvLog.WithError(err).Warning("Not ADMIN role. Access denied")
	// 	httphelper.SendResponse(c, nil, err)
	// 	return
	// }

	skuForm, err := httphelper.SkuUpdateForm(c)
	if err != nil {
		srvLog.WithError(err).Warning("httphelper.ProdCreateForm")
		httphelper.SendResponse(c, nil, err)
		return
	}

	err = ph.PrUsecase.UpdateSku(c, skuForm)
	if err != nil {
		srvLog.WithError(err).Warning("ph.PrUsecase.CreateProduct")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) deleteSku(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "DeleteSku"})
	// exist := httphelper.CheckIfAdmin(c)
	// if !exist {
	// 	err := errorstatus.ErrAuth
	// 	srvLog.WithError(err).Warning("Not ADMIN role. Access denied")
	// 	httphelper.SendResponse(c, nil, err)
	// 	return
	// }

	skuIdStr := c.Param("skuId")
	skuId, err := strconv.Atoi(skuIdStr)
	if err != nil {
		logrus.WithError(err).Warning("utils.DeleteSku.skuId")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	err = ph.PrUsecase.DeleteSku(c, skuId)
	if err != nil {
		srvLog.WithError(err).Warning("ph.PrUsecase.DeleteSku")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) optionBySkuValue(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "OptionBySkuValue"})

	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorstatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	skuValueIdStr := c.Param("skuValueId")
	skuValueId, err := strconv.Atoi(skuValueIdStr)
	if err != nil {
		logrus.WithError(err).Warning("utils.OptionBySkuValue.skuValueId")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	option, err := ph.PrUsecase.GetSkuOption(c, skuValueId, &user.Role)
	if err != nil {
		srvLog.WithError(err).Warning("ph.PrUsecase.DeleteSku")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, option, nil)
}
