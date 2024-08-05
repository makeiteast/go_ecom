package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"go-store/internal/entity"
	errorStatus "go-store/utils/errors"
	httphelper "go-store/utils/http"
)

// OderHandler  represent the httphandler for order
type CategoryHandler struct {
	catUC  entity.CategoryUsecase
	user   entity.UserUsecase
	srvLog *logrus.Entry
}

const (
	sucsess = "sucsess"
)

// NewOrderHandler will initialize the orders/ resources endpoint
func NewCategoryHandler(handler *gin.RouterGroup, mdw gin.HandlerFunc, uc *entity.Usecases, srvLog *logrus.Entry) {
	ch := &CategoryHandler{
		catUC:  uc.CategoryUsecase,
		user:   uc.UserUsecase,
		srvLog: srvLog,
	}
	h := handler.Group("/category")
	{
		h.POST("", ch.getCategory)
		h.POST("/:categoryId", mdw, ch.getCategoryById)
	}
	a := handler.Group("/admin")
	{
		a.POST("/category", ch.createCategory)
		a.PUT("/category/:categoryId", ch.updateCategory)
		a.DELETE("/category/:categoryId", ch.deleteCategory)
		a.POST("/category/:categoryId/option", ch.createCategoryOpt)
		a.POST("/category/optionValue", ch.createCategoryOptValue)
		a.PUT("/category/:categoryId/option", ch.updateCategoryOpt)
		a.PUT("/category/optionValue", ch.updateCategoryOptValue)
		a.DELETE("/category/:categoryId/option", ch.deleteCategoryOpt)
		a.DELETE("/category/optionValue", ch.deleteCategoryOptValue)
	}

}

func (ch *CategoryHandler) getCategory(c *gin.Context) {
	ch.srvLog = log.WithFields(log.Fields{"func": "CategoryHandler.getCategory"})
	result, err := ch.catUC.Get(c)
	if err != nil {
		ch.srvLog.Warning("Cannot getCategory, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, result, nil)
}

func (ch *CategoryHandler) getCategoryById(c *gin.Context) {
	ch.srvLog = log.WithFields(log.Fields{"func": "CategoryHandler.getCategoryById"})

	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		httphelper.SendResponse(c, nil, errorStatus.ErrInternalServer)
		return
	}
	user := userCtx.(*entity.Users)

	categoryId := c.Param("categoryId")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		logrus.WithError(err).Warning("strconv.Atoi.categoryIdInt")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	result, err := ch.catUC.GetById(c, categoryIdInt, &user.Role)
	if err != nil {
		ch.srvLog.Warning("Cannot getCategoryById, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, result, nil)
}

func (ch *CategoryHandler) createCategory(c *gin.Context) {
	ch.srvLog = log.WithFields(log.Fields{"func": "CategoryHandler.createCategory"})

	category, err := httphelper.CategoryForm(c)
	if err != nil {
		ch.srvLog.WithError(err).Warning("httphelper.CategoryForm")
		httphelper.SendResponse(c, nil, err)
		return
	}

	err = ch.catUC.Create(c, category)
	if err != nil {
		ch.srvLog.WithError(err).Warning("ch.catUC.Create")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ch *CategoryHandler) updateCategory(c *gin.Context) {
	ch.srvLog = log.WithFields(log.Fields{"func": "CategoryHandler.updateCategory"})

	categoryId := c.Param("categoryId")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		logrus.WithError(err).Warning("strconv.Atoi.categoryIdInt")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	category, err := httphelper.CategoryForm(c)
	if err != nil {
		ch.srvLog.WithError(err).Warning("httphelper.CategoryForm")
		httphelper.SendResponse(c, nil, err)
		return
	}

	category.Id = categoryIdInt
	category.State = entity.State(c.PostForm("state"))
	_, err = entity.ParseState(string(category.State))
	if err != nil {
		ch.srvLog.WithError(err).Warning("entity.ParseState")
		httphelper.SendResponse(c, nil, err)
		return
	}

	err = ch.catUC.Update(c, category)
	if err != nil {
		ch.srvLog.WithError(err).Warning("ch.catUC.Get")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ch *CategoryHandler) deleteCategory(c *gin.Context) {
	ch.srvLog = log.WithFields(log.Fields{"func": "CategoryHandler.deleteCategory"})

	categoryId := c.Param("categoryId")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		logrus.WithError(err).Warning("strconv.Atoi.categoryIdInt")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	err = ch.catUC.Delete(c, categoryIdInt)
	if err != nil {
		ch.srvLog.Warning("Cannot getCategory, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ch *CategoryHandler) createCategoryOpt(c *gin.Context) {
	ch.srvLog = log.WithFields(log.Fields{"func": "CategoryHandler.addCategoryOpt"})
	categoryId := c.Param("categoryId")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		logrus.Warning("strconv.Atoi.categoryIdInt")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	optionName, ok := c.GetQuery("optionName")
	if !ok {
		logrus.Warning("createCategoryOpt.GetQuery.optionName")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	optionId, err := ch.catUC.CreateOpt(c, categoryIdInt, optionName)
	if err != nil {
		ch.srvLog.Warning("Cannot createCategoryOpt, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	result := map[string]int{"optionId": *optionId}

	httphelper.SendResponse(c, result, nil)
}

func (ch *CategoryHandler) createCategoryOptValue(c *gin.Context) {
	ch.srvLog = log.WithFields(log.Fields{"func": "CategoryHandler.addCategoryOpt"})

	optionIdStr, ok := c.GetQuery("optionId")
	if !ok {
		logrus.Warning("createCategoryOptValue.GetQuery.optionId")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	optionId, err := strconv.Atoi(optionIdStr)
	if err != nil {
		logrus.WithError(err).Warning("strconv.Atoi.categoryIdInt")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	optionValueName, ok := c.GetQuery("optionValueName")
	if !ok {
		logrus.WithError(err).Warning("createCategoryOptValue.GetQuery.optionValueName")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	err = ch.catUC.CreateOptValue(c, optionId, optionValueName)
	if err != nil {
		ch.srvLog.Warning("Cannot createCategoryOpt, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ch *CategoryHandler) updateCategoryOpt(c *gin.Context) {
	ch.srvLog = log.WithFields(log.Fields{"func": "CategoryHandler.updateCategoryOpt"})
	categoryId := c.Param("categoryId")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		logrus.WithError(err).Warning("strconv.Atoi.categoryIdInt")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	optionId, ok := c.GetQuery("optionId")
	if !ok {
		logrus.Warning("updateCategoryOpt.GetQuery.optionId")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}
	optionIdInt, err := strconv.Atoi(optionId)
	if err != nil {
		logrus.WithError(err).Warning("strconv.Atoi.optionIdInt")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	optionName, ok := c.GetQuery("optionName")
	if !ok {
		logrus.Warning("updateCategoryOpt.GetQuery.optionName")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	err = ch.catUC.UpdateCatOpt(c, categoryIdInt, optionIdInt, optionName)
	if err != nil {
		ch.srvLog.Warning("Cannot UpdateCatOpt, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ch *CategoryHandler) updateCategoryOptValue(c *gin.Context) {
	ch.srvLog = log.WithFields(log.Fields{"func": "CategoryHandler.updateCategoryOptValue"})

	optionValueId, ok := c.GetQuery("optionValueId")
	if !ok {
		logrus.Warning("updateCategoryOptValue.GetQuery.optionValueId")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	optionValueIdInt, err := strconv.Atoi(optionValueId)
	if err != nil {
		logrus.WithError(err).Warning("strconv.Atoi.optionValueIdInt")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	optionIdStr, ok := c.GetQuery("optionId")
	if !ok {
		logrus.Warning("updateCategoryOptValue.GetQuery.optionId")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	optionId, err := strconv.Atoi(optionIdStr)
	if err != nil {
		logrus.WithError(err).Warning("strconv.Atoi.categoryIdInt")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	optionValueName, ok := c.GetQuery("optionValueName")
	if !ok {
		logrus.WithError(err).Warning("updateCategoryOptValue.GetQuery.optionValueName")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	err = ch.catUC.UpdateCatOptValue(c, optionId, optionValueIdInt, optionValueName)
	if err != nil {
		ch.srvLog.Warning("Cannot getCategoryById, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ch *CategoryHandler) deleteCategoryOpt(c *gin.Context) {
	ch.srvLog = log.WithFields(log.Fields{"func": "CategoryHandler.deleteCategoryOpt"})
	categoryId, ok := c.GetQuery("categoryId")
	if !ok {
		logrus.Warning("deleteCategoryOpt.GetQuery.categoryId")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		logrus.WithError(err).Warning("strconv.Atoi.categoryIdInt")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	err = ch.catUC.DeleteCatOpt(c, categoryIdInt)
	if err != nil {
		ch.srvLog.Warning("Cannot getCategoryById, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ch *CategoryHandler) deleteCategoryOptValue(c *gin.Context) {
	ch.srvLog = log.WithFields(log.Fields{"func": "CategoryHandler.deleteCategoryOptValue"})

	optionValueId, ok := c.GetQuery("optionValueId")
	if !ok {
		logrus.Warning("deleteCategoryOpt.GetQuery.optionValueId")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	optionValueIdInt, err := strconv.Atoi(optionValueId)
	if err != nil {
		logrus.WithError(err).Warning("strconv.Atoi.optionValueIdInt")
		httphelper.SendResponse(c, nil, errorStatus.ErrBadReq)
		return
	}

	err = ch.catUC.DeleteCatOptValue(c, optionValueIdInt)
	if err != nil {
		ch.srvLog.Warning("Cannot DeleteCatOptValue, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}
