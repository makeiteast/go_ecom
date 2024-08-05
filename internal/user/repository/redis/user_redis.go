package redis_repo

import (
	"context"
	"encoding/json"
	"fmt"
	"go-store/internal/entity"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

// Auth redis repository
type authRedisRepo struct {
	redisClient *redis.Client
}

// Auth redis repository constructor
func NewAuthRedisRepo(redisClient *redis.Client) entity.AuthRedisRepository {
	return &authRedisRepo{redisClient: redisClient}
}

// Get user by id
func (a *authRedisRepo) GetUser(ctx context.Context, username string) (*entity.Claims, error) {
	redLog := log.WithFields(log.Fields{"func": "redis.GetUser"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRedisRepo.GetUser")
	defer span.Finish()

	userBytes, err := a.redisClient.Get(ctx, username).Bytes()
	if err != nil {
		redLog.WithFields(log.Fields{"username": username}).Warning(err)
		return nil, err
	}
	creds := &entity.Claims{}
	if err = json.Unmarshal(userBytes, creds); err != nil {
		redLog.Warning(err)
		return nil, err
	}
	return creds, nil
}

// Cache user with duration in seconds
func (a *authRedisRepo) SetUserCtx(ctx context.Context, username string, seconds int, user *entity.Claims) error {
	redLog := log.WithFields(log.Fields{"func": "redis.SetUserCtx"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRedisRepo.SetUserCtx")
	defer span.Finish()

	userBytes, err := json.Marshal(user)
	if err != nil {
		redLog.WithFields(log.Fields{"username": username}).Warning(err)
		return err
	}
	if err = a.redisClient.Set(ctx, username, userBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return err
	}
	return nil
}

// Get token by UserId
func (a *authRedisRepo) GetUserToken(ctx context.Context, userId int) (*entity.Signatures, error) {
	redLog := log.WithFields(log.Fields{"func": "redis.GetUser"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRedisRepo.GetUserToken")
	defer span.Finish()

	signBytes, err := a.redisClient.Get(ctx, fmt.Sprintf("token-%d", userId)).Result()
	if err != nil {
		redLog.WithFields(log.Fields{"userId": userId}).Warning(err)
		return nil, err
	}
	signs := &entity.Signatures{}
	if err = json.Unmarshal([]byte(signBytes), signs); err != nil {
		redLog.Warning(err)
		return nil, err
	}
	return signs, nil
}

// Cache user token with duration in minutes
func (a *authRedisRepo) SetUserToken(ctx context.Context, tokens *entity.Signatures, userId int, timeExp time.Duration) error {
	redLog := log.WithFields(log.Fields{"func": "redis.SetUserToken"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRedisRepo.SetUserToken")
	defer span.Finish()

	tokenBytes, err := json.Marshal(tokens)
	if err != nil {
		redLog.WithFields(log.Fields{"userId": userId}).Warning(err)
		return err
	}

	if err = a.redisClient.Set(ctx, fmt.Sprintf("token-%d", userId), tokenBytes, timeExp).Err(); err != nil {
		return err
	}
	return nil
}

// Cache user token with duration in minutes
func (a *authRedisRepo) ExpireUserToken(ctx context.Context, userId int, timeExp time.Duration) error {
	redLog := log.WithFields(log.Fields{"func": "redis.SetUserToken"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRedisRepo.SetUserToken")
	defer span.Finish()

	if err := a.redisClient.Expire(ctx, fmt.Sprintf("token-%d", userId), timeExp).Err(); err != nil {
		redLog.WithError(err)
		return err
	}
	return nil
}

// Delte cached user token
func (a *authRedisRepo) DeleteUserToken(ctx context.Context, userId int) error {
	redLog := log.WithFields(log.Fields{"func": "redis.DeleteUserToken"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRedisRepo.DeleteUserToken")
	defer span.Finish()

	if err := a.redisClient.Del(ctx, fmt.Sprintf("token-%d", userId)).Err(); err != nil {
		redLog.WithError(err).Warning()
		return err
	}
	return nil
}
