package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"go-store/internal/entity"
	errorStatus "go-store/utils/errors"
)

func SendResponse(c *gin.Context, data interface{}, handleErr error) {
	// w.Header().Add("Content-Type", "application/json")
	// if data == nil {
	// 	log.Info("Data in 'sendResponse' is null")
	// 	return
	// }
	var status int
	var message string
	if handleErr != nil {

		var apiErr errorStatus.APIError
		if errors.As(handleErr, &apiErr) {
			status, message = apiErr.APIError()
		} else {
			status = http.StatusInternalServerError
			message = "internal server Error"
		}
		c.JSON(status, message)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

func CheckIfAdmin(c *gin.Context) bool {
	userCtx, exists := c.Get("user")
	// This shouldn't happen, as our middleware ought to throw an error.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		return false
	}
	user := userCtx.(*entity.Users)

	if user.Role != entity.UserRoleAdmin {
		return false
	}
	return true
}
