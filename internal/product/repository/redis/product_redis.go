package redis_repo

import (
	"context"
	"encoding/json"
	"go-store/internal/entity"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

// Product redis repository
type prodRedisRepo struct {
	redisClient *redis.Client
}

// Product redis repository constructor
func NewProdRedisRepo(redisClient *redis.Client) entity.ProdRedisRepository {
	return &prodRedisRepo{redisClient: redisClient}
}

// Get products by limit and offset
func (p *prodRedisRepo) GetSku(ctx context.Context, limit int, offset int, category int) ([]*entity.SkuJson, error) {
	redLog := log.WithFields(log.Fields{"func": "redis.GetSku"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "prodRedisRepo.GetSku")
	defer span.Finish()
	key := "skuprod" + strconv.Itoa(category)
	slice, err := p.redisClient.Do(ctx, "ZRANGEBYSCORE", key, offset, limit+offset-1).StringSlice()
	if err != nil || len(slice) == 0 {
		redLog.WithFields(log.Fields{"limit": limit}).Warning(err)
		return nil, err
	}
	products := make([]*entity.SkuJson, len(slice))
	for i := 0; i < len(slice); i++ {
		if err = json.Unmarshal([]byte(slice[i]), &products[i]); err != nil {
			redLog.Warning(err)
			return nil, err
		}
	}
	return products, nil
}

// Cache products with duration in seconds
func (p *prodRedisRepo) SetSkuCtx(ctx context.Context, offset int, category int, prod *entity.SkuJson) error {
	redLog := log.WithFields(log.Fields{"func": "redis.SetSkuProdCtx"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "prodRedisRepo.SetSkuProdCtx")
	defer span.Finish()

	prodBytes, err := json.Marshal(prod)
	if err != nil {
		redLog.WithFields(log.Fields{"offset": offset}).Warning(err)
		return err
	}
	key := "skuprod" + strconv.Itoa(category)
	err = p.redisClient.ZAdd(ctx, key, &redis.Z{Score: float64(offset), Member: prodBytes}).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get user by id
func (p *prodRedisRepo) GetProducts(ctx context.Context, limit int, offset int, category int) ([]*entity.ProductJson, error) {
	redLog := log.WithFields(log.Fields{"func": "redis.GetProducts"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "prodRedisRepo.GetProducts")
	defer span.Finish()
	key := "prod" + strconv.Itoa(category)
	slice, err := p.redisClient.Do(ctx, "ZRANGEBYSCORE", key, offset, limit+offset-1).StringSlice()
	if err != nil {
		redLog.WithFields(log.Fields{"limit": limit}).Warning(err)
		return nil, err
	}
	products := make([]*entity.ProductJson, len(slice))
	for i := 0; i < len(slice); i++ {
		if err = json.Unmarshal([]byte(slice[i]), &products[i]); err != nil {
			redLog.Warning(err)
			return nil, err
		}
	}
	return products, nil
}

// Cache user with duration in seconds
func (p *prodRedisRepo) SetProdCtx(ctx context.Context, offset int, category int, prod *entity.ProductJson) error {
	redLog := log.WithFields(log.Fields{"func": "redis.SetUserCtx"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "prodRedisRepo.SetProdCtx")
	defer span.Finish()

	prodBytes, err := json.Marshal(prod)
	if err != nil {
		redLog.WithFields(log.Fields{"offset": offset}).Warning(err)
		return err
	}
	key := "prod" + strconv.Itoa(category)
	err = p.redisClient.ZAdd(ctx, key, &redis.Z{Score: float64(offset), Member: prodBytes}).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get Cache count
func (p *prodRedisRepo) GetSkuCount(ctx context.Context, key string) (string, error) {
	redLog := log.WithFields(log.Fields{"func": "redis.GetProdCount"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "prodRedisRepo.SetProdCount")
	defer span.Finish()
	count, err := p.redisClient.Get(ctx, key).Result()
	if err != nil {
		redLog.Warning(err)
		return "", err
	}
	return count, nil
}

// Set Cache user with duration in seconds
func (p *prodRedisRepo) SetSkuCount(ctx context.Context, count string, key string) error {
	redLog := log.WithFields(log.Fields{"func": "redis.SetProdCount"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "prodRedisRepo.SetProdCount")
	defer span.Finish()
	_, err := p.redisClient.Set(ctx, key, count, time.Duration(120*time.Second)).Result()
	if err != nil {
		redLog.Warning(err)
		return err
	}
	return nil
}

// Get product by id
func (p *prodRedisRepo) GetProdByIDCtx(ctx context.Context, key string) (*entity.SkuJson, error) {
	redLog := log.WithFields(log.Fields{"func": "redis.GetProdByIDCtx"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "redis.GetProdByIDCtx")
	defer span.Finish()

	prodBytes, err := p.redisClient.Get(ctx, key).Bytes()
	if err == redis.Nil {
		redLog.WithFields(log.Fields{"key": key}).Warning(err)
		return nil, nil
	}
	if err != nil {
		redLog.WithFields(log.Fields{"key": key}).Warning(err)
		return nil, err
	}
	prod := &entity.SkuJson{}
	if err = json.Unmarshal(prodBytes, prod); err != nil {
		return nil, err
	}
	return prod, nil
}

// Cache product with duration in seconds
func (p *prodRedisRepo) SetProdByIDCtx(ctx context.Context, key string, prod *entity.SkuJson) error {
	redLog := log.WithFields(log.Fields{"func": "redis.SetProdByIDCtx"})
	span, ctx := opentracing.StartSpanFromContext(ctx, "redis.SetProdByIDCtx")
	defer span.Finish()

	prodBytes, err := json.Marshal(prod)
	if err != nil {
		redLog.WithFields(log.Fields{"key": key}).Warning(err)
		return err
	}
	if err = p.redisClient.Set(ctx, key, prodBytes, time.Second*time.Duration(120)).Err(); err != nil {
		return err
	}
	return nil
}
