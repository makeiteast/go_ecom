package entity

import (
	"errors"
	"time"
)

type Usecases struct {
	UserUsecase     UserUsecase
	OrderUsecase    OrderUsecase
	ProductUsecase  ProductUsecase
	CategoryUsecase CategoryUsecase
	CartUsecase     CartUsecase
}

type PaginationParam struct {
	Limit    int `json:"limit"`
	Offset   int `json:"offset"`
	Category int `json:"category"`
}

type State string

const (
	Enabled  State = "enabled"
	Disabled State = "disabled"
	Deleted  State = "deleted"
)

func NowUTC() time.Time {
	return time.Now().Round(time.Microsecond).UTC()
}

func ParseState(s string) (r State, err error) {
	rt := State(s)
	switch rt {
	case Enabled,
		Disabled,
		Deleted:
		r = rt
		return
	default:
		return "", ErrTypeNotMatched
	}
}

var (
	ErrUnauthorized         = errors.New("unauthorized")
	ErrPasswordHashMismatch = errors.New("password hash mismatch")
	ErrTypeNotMatched       = errors.New("type not matched")
	ErrNotFound             = errors.New("not found")
	ErrForbidden            = errors.New("forbidden")
	ErrAlreadyExists        = errors.New("already exists")
	ErrBadRequest           = errors.New("bad request")
	ErrInvalidInput         = errors.New("invalid input")
	ErrConflict             = errors.New("conflict")
	ErrTooManyRequests      = errors.New("too many requests")
	ErrInternalServerError  = errors.New("internal server error")
)
