package entity

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Users struct {
	Id               int
	PublicId         uuid.UUID `db:"public_id"`
	Username         string    `db:"author"`
	FullName         string    `db:"full_name"`
	Password         string    `db:"password"`
	Email            string    `db:"email"`
	PhoneNumber      string    `db:"phone_number"`
	Address          string    `db:"address"`
	Photo            string    `db:"photo"`
	Role             UserRole  `db:"user_role"`
	RegionId         int       `db:"region_id"`
	Parent           int       `db:"parent"`
	VerificationCode string    `db:"verification_code"`
	CreateTs         time.Time `json:"createTs"`
	UpdateTs         time.Time `json:"updateTs"`
	State            State     `db:"state"`
	Version          int       `db:"version"`
}

type (
	UserRole string
)

const (
	UserRoleAdmin UserRole = "ADMIN"
	UserRoleUser  UserRole = "USER"
	UserRoleGuest UserRole = "GUEST"
)

type Claims struct {
	Id       int      `json:"id"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
}

type JwtClaims struct {
	Id       int      `json:"id"`
	UID      string   `json:"uid"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
	jwt.StandardClaims
	Tokens Tokens `json:"tokens"`
}

type Tokens struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type Signatures struct {
	AccessSign  string `json:"access"`
	RefreshSign string `json:"refresh"`
}

type TokenConf struct {
	AccesTokenTimeout   time.Duration
	RefreshTokenTimeout time.Duration
	AutoLogoffTimeout   time.Duration
	AccessSecret        []byte
	RefreshSecret       []byte
}

type UserJson struct {
	PublicId uuid.UUID `json:"public_id"`
	Username string    `json:"username"`
	Role     UserRole  `json:"role"`
	Tokens   Tokens    `json:"tokens"`
}

type UserUsecase interface {
	Login(ctx context.Context, username string, password string, tokenConf *TokenConf) (userJs *UserJson, err error)
	Logout(ctx context.Context, userId int) (err error)
	Refresh(ctx context.Context, claims *JwtClaims, tokenConf *TokenConf) (userJs *UserJson, err error)
	ParseToken(ctx context.Context, tokenStr string, secret []byte) (claims *JwtClaims, err error)
	ValidateToken(ctx context.Context, claims *JwtClaims, tokenMdlw string, isRefresh bool) (user *Users, err error)
	TokenExpire(ctx context.Context, userId int, timeExp time.Duration) (err error)
	CreateUser(ctx context.Context, user *Users, sendMethod *string) (err error)
	ActivateUser(ctx context.Context, username string, verificationCode string) (err error)
}

type AuthPgxRepository interface {
	UserByUsername(ctx context.Context, username string) (creds *Users, err error)
	Create(ctx context.Context, user *Users) (err error)
	Update(ctx context.Context, user *Users) (err error)
}

type AuthRedisRepository interface {
	GetUser(ctx context.Context, username string) (creds *Claims, err error)
	SetUserCtx(ctx context.Context, username string, seconds int, creds *Claims) error
	GetUserToken(ctx context.Context, userId int) (tokens *Signatures, err error)
	SetUserToken(ctx context.Context, user *Signatures, userId int, timeExp time.Duration) error
	ExpireUserToken(ctx context.Context, userId int, timeExp time.Duration) (err error)
	DeleteUserToken(ctx context.Context, userId int) error
}

type AuthBroker interface {
	SendEmail(ctx context.Context, to string, message []byte) (err error)
}
