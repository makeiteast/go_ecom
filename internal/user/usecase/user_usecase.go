package usecase

import (
	"bytes"
	"context"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"go-store/internal/entity"
	errorStatus "go-store/utils/errors"
	generator "go-store/utils/generator"
)

// UserUsecase will initiate usecase of entity.AuthPgxRepository interface
type UserUsecase struct {
	authRepo   entity.AuthPgxRepository
	redisRepo  entity.AuthRedisRepository
	brokerRepo entity.AuthBroker
}

// NewAuthUsecase will create new an UserUsecase object representation of entity.UserUsecase interface
func NewAuthUsecase(a entity.AuthPgxRepository, r entity.AuthRedisRepository) entity.UserUsecase {
	return &UserUsecase{
		authRepo:  a,
		redisRepo: r,
	}
}

// Method GetLoginUser which recieve authentification credintials and return users token, role and username
func (a *UserUsecase) Login(ctx context.Context, username string, password string, tokenConf *entity.TokenConf) (result *entity.UserJson, err error) {

	ctLog := log.WithFields(log.Fields{"func": "UserUsecase.GetLoginUser"})
	// take from database username data
	userDb, err := a.authRepo.UserByUsername(ctx, username)
	if err != nil {
		ctLog.WithError(err).Warning("a.authRepo.UserByUsername - can't get user")
		return nil, err
	}
	// recieve hashed password from database and compared with hash
	if err = bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(password)); err != nil {
		ctLog.WithError(err).Warning("bcrypt.CompareHashAndPassword - User credentials invalid")
		return nil, err
	}
	var accessToken, refreshToken *string

	accessExpirationTime := time.Now().Add(tokenConf.AccesTokenTimeout)
	refreshExpirationTime := time.Now().Add(tokenConf.RefreshTokenTimeout)

	accessToken, accSign, err := a.createToken(ctx, userDb, accessExpirationTime, tokenConf.AccessSecret)
	refreshToken, refSign, err := a.createToken(ctx, userDb, refreshExpirationTime, tokenConf.RefreshSecret)

	// Create the JWT claims, which includes the username and expiry time
	user := &entity.UserJson{
		PublicId: userDb.PublicId,
		Username: userDb.Username,
		Role:     userDb.Role,
	}
	tokens := &entity.Tokens{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}
	user.Tokens = *tokens

	signatures := &entity.Signatures{
		AccessSign:  *accSign,
		RefreshSign: *refSign,
	}
	err = a.redisRepo.SetUserToken(ctx, signatures, userDb.Id, tokenConf.AutoLogoffTimeout)
	if err != nil {
		ctLog.WithError(err).Warning("a.redisRepo.SetUserToken")
		err = errorStatus.ErrInternalServer
		return
	}

	return user, nil
}

// Method GetLoginUser which recieve authentification credintials and return users token, role and username
func (a *UserUsecase) Refresh(ctx context.Context, claims *entity.JwtClaims, tokenConf *entity.TokenConf) (result *entity.UserJson, err error) {

	ctLog := log.WithFields(log.Fields{"func": "UserUsecase.Refresh"})
	userDb, err := a.authRepo.UserByUsername(ctx, claims.Username)
	if err != nil {
		ctLog.WithError(err).Warning("a.authRepo.UserByUsername - can't get user")
		return nil, err
	}

	var accessToken, refreshToken *string

	accessExpirationTime := time.Now().Add(tokenConf.AccesTokenTimeout)
	refreshExpirationTime := time.Now().Add(tokenConf.RefreshTokenTimeout)

	accessToken, accSign, err := a.createToken(ctx, userDb, accessExpirationTime, tokenConf.AccessSecret)
	refreshToken, refSign, err := a.createToken(ctx, userDb, refreshExpirationTime, tokenConf.RefreshSecret)

	// Create the JWT claims, which includes the username and expiry time
	userJs := &entity.UserJson{
		PublicId: userDb.PublicId,
		Username: userDb.Username,
		Role:     userDb.Role,
	}
	tokens := &entity.Tokens{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}
	userJs.Tokens = *tokens

	signatures := &entity.Signatures{
		AccessSign:  *accSign,
		RefreshSign: *refSign,
	}
	err = a.redisRepo.SetUserToken(ctx, signatures, userDb.Id, tokenConf.AutoLogoffTimeout)
	if err != nil {
		ctLog.WithError(err).Warning("a.redisRepo.SetUserToken")
		return nil, err
	}

	return userJs, nil
}

func (a *UserUsecase) CheckIsAuthorized(ctx context.Context, tokenString string) (claims *entity.Claims, err error) {
	return nil, nil
}

// Method Logout which recieve userId and remove tokens with this id from cache
func (a *UserUsecase) Logout(ctx context.Context, userId int) (err error) {
	ctLog := log.WithFields(log.Fields{"func": "UserUsecase.Logout"})
	// take from database username data

	err = a.redisRepo.DeleteUserToken(ctx, userId)
	if err != nil {
		ctLog.WithError(err).Warning("a.redisRepo.DeleteUserToken")
		err = errorStatus.ErrInternalServer
		return
	}

	return nil
}

func (a *UserUsecase) ValidateToken(ctx context.Context, claims *entity.JwtClaims, uidMdlw string, isRefresh bool) (user *entity.Users, err error) {
	ctLog := log.WithFields(log.Fields{"func": "UserUsecase.ValidateToken"})
	cachedJson, err := a.redisRepo.GetUserToken(ctx, claims.Id)
	if err != nil {
		ctLog.WithError(err).Warning("a.redisRepo.GetUserToken")
		err = errorStatus.ErrInternalServer
		return
	}
	var tokenUid string
	if isRefresh {
		tokenUid = cachedJson.RefreshSign
	} else {
		tokenUid = cachedJson.AccessSign
	}

	if err != nil || tokenUid != uidMdlw {
		ctLog.Warning("token not validate")
		err = errorStatus.ErrAuth
		return nil, err
	}

	userDb, err := a.authRepo.UserByUsername(ctx, claims.Username)
	if err != nil {
		ctLog.WithError(err).Warning("a.authRepo.UserByUsername")
		err = errorStatus.ErrInternalServer
		return nil, err
	}
	if userDb.Id == 0 {
		ctLog.WithError(err).Warning("a.authRepo.UserByUsername - User not found")
		err = errorStatus.ErrAuth
		return nil, err
	}

	user = &entity.Users{
		Id:       userDb.Id,
		PublicId: userDb.PublicId,
		Username: userDb.Username,
		Role:     userDb.Role,
	}

	return user, nil
}

func (a *UserUsecase) TokenExpire(ctx context.Context, userId int, timeExp time.Duration) (err error) {
	ctLog := log.WithFields(log.Fields{"func": "UserUsecase.TokenExpire"})
	err = a.redisRepo.ExpireUserToken(ctx, userId, timeExp)
	if err != nil {
		ctLog.WithError(err).Warning("a.redisRepo.ExpireUserToken")
		err = errorStatus.ErrInternalServer
		return
	}

	return nil
}

func (a *UserUsecase) createToken(ctx context.Context, user *entity.Users, timeExp time.Time, jwtKey []byte) (tokenStr *string, sign *string, err error) {
	// Declare the expiration time of the token
	// here, we have kept it as X minutes
	uuidS := uuid.New().String()
	claim := &entity.JwtClaims{
		Id:       user.Id,
		Username: user.Username,
		Role:     user.Role,
		UID:      uuidS,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: timeExp.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, nil, err
	}

	return &tokenString, &uuidS, nil
}

func (a *UserUsecase) ParseToken(ctx context.Context, tokenStr string, secret []byte) (claims *entity.JwtClaims, err error) {
	ctLog := log.WithFields(log.Fields{"func": "UserUsecase.ParseToken"})
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := errorStatus.ErrToken
			return nil, err
		}
		return secret, nil
	})
	if err != nil {
		ctLog.WithError(err).Warning("jwt.Parse")
		return nil, err
	}
	claimsM := token.Claims.(jwt.MapClaims)

	claims = &entity.JwtClaims{
		Id:       int(claimsM["id"].(float64)),
		UID:      claimsM["uid"].(string),
		Username: claimsM["username"].(string),
		Role:     entity.UserRole(claimsM["role"].(string)),
	}
	return claims, nil
}

func (a *UserUsecase) CreateUser(ctx context.Context, user *entity.Users, sendMethod *string) (err error) {
	ctLog := log.WithFields(log.Fields{"func": "UserUsecase.CreateUser"})

	user.PublicId = uuid.New()
	user.VerificationCode = generator.RandStringRunes(6)
	user.Role = entity.UserRoleUser

	err = a.authRepo.Create(ctx, user)
	if err != nil {
		ctLog.WithError(err).Warning("a.authRepo.Create")
		return err
	}

	tmpl := template.Must(template.ParseFiles("utils/broker/email.html"))
	buff := new(bytes.Buffer)
	if err = tmpl.Execute(buff, struct{ Code string }{user.VerificationCode}); err != nil {
		ctLog.WithError(err).Warning("tmpl.Execute")
		return err
	}

	err = a.brokerRepo.SendEmail(ctx, user.Email, buff.Bytes())
	if err != nil {
		ctLog.WithError(err).Warning("a.brokerRepo.SendEmail")
		return err
	}

	return nil
}

func (a *UserUsecase) ActivateUser(ctx context.Context, username string, verificationCode string) (err error) {
	ctLog := log.WithFields(log.Fields{"func": "UserUsecase.ActivateUser"})

	user, err := a.authRepo.UserByUsername(ctx, username)
	if err != nil {
		ctLog.WithError(err).Warning("a.authRepo.UserByUsername")
		return err
	}

	if !(user.State == entity.Disabled && user.VerificationCode == verificationCode) {
		ctLog.WithError(err).Warning("not permit!")
		return err
	}

	user.UpdateTs = time.Now()
	user.State = entity.Enabled

	err = a.authRepo.Update(ctx, user)
	if err != nil {
		ctLog.WithError(err).Warning("a.authRepo.Update")
		return err
	}

	return nil
}
