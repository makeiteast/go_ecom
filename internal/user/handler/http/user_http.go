package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"go-store/internal/entity"
	"go-store/internal/user/dto"
	errorstatus "go-store/utils/errors"
	httphelper "go-store/utils/http"
)

// UserHandler  represent the httphandler for authentification
type UserHandler struct {
	UserUsecase entity.UserUsecase
	srvLog      *logrus.Entry
	tokenConf   *entity.TokenConf
}

// NewUserHandler will initialize the user resources endpoint
func NewUserHandler(handler *gin.RouterGroup, mdw gin.HandlerFunc, uc *entity.Usecases, srvLog *logrus.Entry, tokenConf *entity.TokenConf) {
	ah := &UserHandler{
		UserUsecase: uc.UserUsecase,
		srvLog:      srvLog,
		tokenConf:   tokenConf,
	}
	h := handler.Group("/user")
	{
		h.POST("/login", ah.loginHandler)
		h.POST("/refresh", ah.refresh)
		h.POST("/logout", ah.logoutHandler)
		h.POST("/register", ah.registerHandler)
		h.POST("/activate", ah.activateHandler)
	}
}

// Handler to Login(Sign in) and create jwt key
func (ah *UserHandler) loginHandler(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "server.LoginHandler"})
	// Get the JSON body and decode into credentials
	var creds dto.Credentials
	if err := json.NewDecoder(c.Request.Body).Decode(&creds); err != nil {
		srvLog.WithError(err).Error("format is wrong")
		httphelper.SendResponse(c, nil, errorstatus.ErrBadReq)
		return
	}

	ah.srvLog = ah.srvLog.WithFields(log.Fields{"username": creds.Username})
	// Create jwt key if credentials are correct
	result, err := ah.UserUsecase.Login(c, creds.Username, creds.Password, ah.tokenConf)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		ah.srvLog.WithError(err).Warning("User credentials invalid!")
		httphelper.SendResponse(c, nil, err)
		return
	}
	httphelper.SendResponse(c, result, nil)
}

// Handler to Logout(Sign out)
func (ah *UserHandler) refresh(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "server.refresh"})
	// Get the JSON body and decode into credentials
	tokenStr := c.Request.Header.Get("x-access-token")
	claims, err := ah.UserUsecase.ParseToken(c, tokenStr, ah.tokenConf.RefreshSecret)
	if err != nil {
		srvLog.WithError(err).Warning("ah.UserUsecase.ParseToken")
		httphelper.SendResponse(c, nil, errorstatus.ErrAuth)
		return
	}
	// Validate token
	_, err = ah.UserUsecase.ValidateToken(c, claims, claims.UID, true)
	if err != nil {
		srvLog.WithError(err).Warning("refresh token not valid")
		httphelper.SendResponse(c, nil, errorstatus.ErrBadReq)
		return
	}

	user, err := ah.UserUsecase.Refresh(c, claims, ah.tokenConf)
	if err != nil {
		ah.srvLog.Warning(err)
		httphelper.SendResponse(c, nil, errorstatus.ErrAuth)
		return
	}

	httphelper.SendResponse(c, user, nil)
}

// Handler to Login(Sign in) and create jwt key
func (ah *UserHandler) logoutHandler(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "server.logoutHandler"})
	// Get the JSON body and decode into credentials
	tokenStr := c.Request.Header.Get("x-access-token")
	claims, err := ah.UserUsecase.ParseToken(c, tokenStr, ah.tokenConf.AccessSecret)
	if err != nil {
		srvLog.WithError(err).Warning("ah.UserUsecase.ParseToken")
		httphelper.SendResponse(c, nil, errorstatus.ErrAuth)
		return
	}
	// Delete tokens from cache
	err = ah.UserUsecase.Logout(c, claims.Id)
	if err != nil {
		// If the cache is empty return an HTTP error
		srvLog.WithError(err).Warning("User logout error")
		httphelper.SendResponse(c, errorstatus.ErrInternalServer, err)
		return
	}
	httphelper.SendResponse(c, "success", nil)
}

// gin Middleware, which check jwt.Token valid then Validate it with cached and requested token
// and add to cache (Expire) time of expire
func ValidateJWT(uc entity.UserUsecase, tokenConf *entity.TokenConf) gin.HandlerFunc {
	srvLog := log.WithFields(log.Fields{"func": "server.ValidateJWT"})
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("x-access-token")
		if tokenStr != "" {
			claims, err := uc.ParseToken(c, tokenStr, tokenConf.AccessSecret)
			if err != nil {
				srvLog.WithError(err).Warning("uc.ParseToken")
				c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
				return
			}
			user, err := uc.ValidateToken(c, claims, claims.UID, false)
			if err != nil {
				srvLog.WithError(err).Warning("uc.ValidateToken")
				c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
				return
			}

			err = uc.TokenExpire(c, claims.Id, tokenConf.AutoLogoffTimeout)
			if err != nil {
				srvLog.WithError(err).Warning("uc.TokenExpire")
				c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
				return
			}
			c.Set("user", user)
		} else {
			user := &entity.Users{
				Role: entity.UserRoleGuest,
			}
			c.Set("user", user)
		}
		c.Next()
	}
}

func (ah *UserHandler) registerHandler(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "server.registerHandler"})

	userForm, sendMethod, err := httphelper.UserCreateForm(c)
	if err != nil {
		srvLog.WithError(err).Warning("httphelper.UserCreateForm")
		httphelper.SendResponse(c, nil, errorstatus.ErrBadReq)
		return
	}

	err = ah.UserUsecase.CreateUser(c, userForm, sendMethod)
	if err != nil {
		srvLog.WithError(err).Warning("ah.UserUsecase.CreateUser")
		httphelper.SendResponse(c, nil, errorstatus.ErrBadReq)
		return
	}

	httphelper.SendResponse(c, "success", nil)
}

func (ah *UserHandler) activateHandler(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "server.activateHandler"})

	username := c.Param("username")
	code := c.Param("code")

	err := ah.UserUsecase.ActivateUser(c, username, code)
	if err != nil {
		srvLog.WithError(err).Warning("ah.UserUsecase.ActivateUser")
		httphelper.SendResponse(c, nil, errorstatus.ErrBadReq)
		return
	}

	httphelper.SendResponse(c, "success", nil)
}
