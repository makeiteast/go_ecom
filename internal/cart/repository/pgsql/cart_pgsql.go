package pgsql

import (
	"context"
	"fmt"
	"go-store/internal/entity"
	"time"

	"go-store/utils/database"

	"github.com/Masterminds/squirrel"
	log "github.com/sirupsen/logrus"
)

type PgxAccess struct {
	*database.PgxAccess
}

func NewPgxCartRepository(pgx *database.PgxAccess) entity.CartRepository {
	return &PgxAccess{pgx}
}

func (d *PgxAccess) GetCarts(ctx context.Context, limit int, offset int, userId *int) (result []*entity.Cart, err error) {
	dbLog := log.WithFields(log.Fields{"func": "db.GetOrderQuery"})
	baseQuery := d.Builder.
		Select("user_id",
			"sku_id",
			"quantity",
			"create_ts",
			"update_ts").
		From("cart").
		Where("state = 'enabled'").
		Limit(uint64(limit)).
		Offset(uint64(offset))
	if userId != nil {
		baseQuery = baseQuery.Where("cart.user_id = ?", userId)
	}
	query, args, err := baseQuery.ToSql()
	rows, err := d.Pool.Query(ctx, query, args...)
	if err != nil {
		dbLog.Warning(err)
		err = fmt.Errorf("db.GetCarts: %w", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tmp entity.Cart
		if err := rows.Scan(&tmp.UserId, &tmp.SkuId, &tmp.Quantity, &tmp.CreateTs, &tmp.UpdateTs); err != nil {
			dbLog.WithFields(log.Fields{"user_id": userId}).Warning(err)
			return nil, err
		}
		result = append(result, &tmp)

	}
	return result, nil

}

func (d *PgxAccess) CreateCart(ctx context.Context, cart *entity.Cart) (err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.CreateOrder"})
	query, args, err := d.Builder.
		Insert("cart").
		Columns("user_id",
			"sku_id",
			"quantity",
			"create_ts",
			"update_ts",
			"state",
			"version").
		Values(cart.UserId,
			cart.SkuId,
			cart.Quantity,
			cart.CreateTs,
			cart.UpdateTs,
			cart.State,
			cart.Version).
		ToSql()
	_, err = d.Pool.Exec(ctx, query, args...)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - CreateCart - Exec")
		return err
	}
	return nil
}

func (d *PgxAccess) UpdateCart(ctx context.Context, cart *entity.Cart) (err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.UpdateOrder"})
	query, args, err := d.Builder.
		Update("cart").
		SetMap(map[string]interface{}{
			"quantity":  cart.Quantity,
			"update_ts": time.Now(),
			"state":     cart.State}).
		Set("version", squirrel.Expr("version+1")).
		Where("sku_id = $1 AND user_id = $2", cart.SkuId, cart.UserId).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - UpdateCart - r.Builder - query")
		return err
	}
	_, err = d.Pool.Exec(ctx, query, args...)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - UpdateCart - Exec")
		return err
	}
	return nil
}

func (d *PgxAccess) DeleteCart(ctx context.Context, cart *entity.Cart) (err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.DeleteOrder"})
	query, args, err := d.Builder.
		Delete("cart").
		Where("sku_id = $1 AND user_id = $2", cart.SkuId, cart.UserId).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - UpdateCart - r.Builder - query")
		return err
	}
	_, err = d.Pool.Exec(ctx, query, args...)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - DeleteOrder - Exec")
		return err
	}
	return nil
}
