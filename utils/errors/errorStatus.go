package errors

import (
	"net/http"
)

//Handling Errors

type APIError interface {
	// APIError returns an HTTP status code and an API-safe error message.
	APIError() (int, string)
}

type sentinelAPIError struct {
	status int
	msg    string
}

func (e sentinelAPIError) Error() string {
	return e.msg
}

func (e sentinelAPIError) APIError() (int, string) {
	return e.status, e.msg
}

var (
	ErrAuth           = &sentinelAPIError{status: http.StatusUnauthorized, msg: "not authorized"}
	ErrToken          = &sentinelAPIError{status: http.StatusUnauthorized, msg: "invalid token"}
	ErrPermission     = &sentinelAPIError{status: http.StatusNetworkAuthenticationRequired, msg: "permission denied"}
	ErrNotFound       = &sentinelAPIError{status: http.StatusNotFound, msg: "not found"}
	ErrBadReq         = &sentinelAPIError{status: http.StatusBadRequest, msg: "bad request"}
	ErrDuplicate      = &sentinelAPIError{status: http.StatusBadRequest, msg: "duplicate"}
	ErrInternalServer = &sentinelAPIError{status: http.StatusInternalServerError, msg: "internal server"}
)
