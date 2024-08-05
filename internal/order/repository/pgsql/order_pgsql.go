package pgsql

import (
	"context"
	"errors"
	"fmt"
	"go-store/internal/entity"
	"time"

	errorStatus "go-store/utils/errors"

	"go-store/utils/database"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type PgxAccess struct {
	*database.PgxAccess
}

func NewPgxOrderRepository(pgx *database.PgxAccess) entity.OrderRepository {
	return &PgxAccess{pgx}
}

func (r *PgxAccess) NewTxId(ctx context.Context) (int, error) {
	id, err := r.PgTxBegin(ctx)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PgxAccess) TxEnd(ctx context.Context, txId int, err error) error {
	return r.PgTxEnd(ctx, txId, err)
}

func (d *PgxAccess) GetOrders(ctx context.Context, filterMap map[string]string, limit int, offset int, id int) (result []*entity.Order, err error) {
	dbLog := log.WithFields(log.Fields{"func": "db.GetOrderQuery"})
	baseQuery := d.Builder.
		Select("id",
			"user_id",
			"address",
			"phone",
			"comment",
			"status",
			"create_ts",
			"update_ts",
			"state",
			"version",
			"notes").
		From("orders").
		Where("state = 'enabled'").
		Limit(uint64(limit)).
		Offset(uint64(offset))
	for k, v := range filterMap {
		baseQuery = baseQuery.Where(fmt.Sprintf("%s = ?", k), v)
	}
	query, args, err := baseQuery.ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - GetOrders - r.Builder - query")
		return nil, err
	}
	dbLog.Info(query)
	dbLog.Info(args...)
	rows, err := d.Pool.Query(ctx, query, args...)
	if err != nil {
		dbLog.Warning(err)
		err = fmt.Errorf("db.GetOrderQuery: %w", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tmp entity.Order
		if err := rows.Scan(&tmp.Id, &tmp.UserId, &tmp.Address, &tmp.Phone, &tmp.Comment, &tmp.Status, &tmp.CreateTs, &tmp.UpdateTs, &tmp.State, &tmp.Version, &tmp.Notes); err != nil {
			dbLog.WithFields(log.Fields{"order_id": tmp.Id}).Warning(err)
			err = fmt.Errorf("db.GetOrderQuery: %w", err)
			return nil, err
		}
		result = append(result, &tmp)

	}
	return result, nil

}

func (d *PgxAccess) GetItem(ctx context.Context, order_id uuid.UUID) (result []*entity.OrderItem, err error) {
	dbLog := log.WithFields(log.Fields{"func": "db.GetOrderQuery"})
	query, args, err := d.Builder.
		Select("order_item.id",
			"order_item.order_id",
			"order_item.quantity",
			"order_item.price",
			"order_item.create_ts",
			"order_item.update_ts",
			"order_item.state",
			"order_item.version",
			"sku.sku",
			"sku.small_name").
		From("order_item").
		InnerJoin("sku ON order_item.sku_id = sku.id").
		Where("order_item.state = 'enabled' AND order_item.order_id = $1", order_id).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - GetItem - r.Builder - query")
		return nil, err
	}
	rows, err := d.Pool.Query(ctx, query, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.New("db: No Result")
			dbLog.Warning(err)
			err = nil
			return nil, err
		}
		err = fmt.Errorf("db: %w", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tmp entity.OrderItem
		if err := rows.Scan(&tmp.Id, &tmp.OrderId, &tmp.Quantity, &tmp.Price, &tmp.CreateTs, &tmp.UpdateTs, &tmp.State, &tmp.Version, &tmp.Sku.Sku, &tmp.Sku.SmallImage); err != nil {
			dbLog.WithFields(log.Fields{"item_id": tmp.Id}).Warning(err)
			err = fmt.Errorf("db: %w", err)
			return nil, err
		}
		result = append(result, &tmp)
	}
	return result, nil

}

func (d *PgxAccess) CreateOrder(ctx context.Context, order *entity.Order, txId int) (orderId *uuid.UUID, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.CreateOrder"})

	tx, err := d.GetTxById(txId)
	if err != nil {
		dbLog.WithError(err).Errorf("OrderRepo - Update - r.GetTxById")
		return nil, err
	}

	query, args, err := d.Builder.
		Insert("orders").
		Columns(
			"id",
			"user_id",
			"address",
			"phone",
			"comment",
			"notes",
			"status",
			"create_ts",
			"update_ts",
			"state",
			"version").
		Values(
			squirrel.Expr("uuid_generate_v4()"),
			order.UserId,
			order.Address,
			order.Phone,
			order.Comment,
			order.Notes,
			order.Status,
			order.CreateTs,
			order.UpdateTs,
			order.State,
			order.Version).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - CreateOrder - r.Builder - query")
		return nil, err
	}
	row := tx.QueryRow(ctx, query, args...)
	err = row.Scan(&orderId)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - CreateOrder - QueryRow")
		return nil, err
	}
	return orderId, nil
}

func (d *PgxAccess) CreateOrderItem(ctx context.Context, orderId uuid.UUID, userId int, txId int) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.CreateOrderItem"})
	var idCount int

	tx, err := d.GetTxById(txId)
	if err != nil {
		dbLog.WithError(err).Errorf("OrderRepo - CreateOrderItem - d.GetTxById")
		return err
	}

	queryStr := fmt.Sprintf(`WITH I AS ( INSERT
		INTO "order_item"(
			order_id,
			sku_id,
			quantity,
			price,
			create_ts,
			update_ts,
			state,
			version)
		SELECT
			$2,
			cart.sku_id,
			cart.quantity,
			sku.price,
			$3,
			$4,
			$5,
			$6
		FROM "cart"
		INNER JOIN "sku" ON sku.id = cart.sku_id
		RETURNING id ), D AS(
			DELETE FROM "cart"
			WHERE cart.user_id = $1)
		SELECT id FROM I`)
	row := tx.QueryRow(ctx, queryStr, userId, orderId, time.Now(), time.Now(), entity.Enabled, 0)
	err = row.Scan(&idCount)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return err
	}
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - CreateOrderItem - Exec")
		return err
	}
	return nil
}

func (d *PgxAccess) GetOrder(ctx context.Context, orderId string) (result *entity.Order, err error) {
	order := &entity.Order{}
	dbLog := log.WithFields(log.Fields{"func": "db.GetOrder"})
	query, args, err := d.Builder.
		Select("id",
			"user_id",
			"address",
			"phone",
			"comment",
			"status",
			"notes").
		From("orders").
		Where("orders.id = $1 AND orders.state = 'enabled'", orderId).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - GetOrder - r.Builder - query")
		return nil, err
	}
	row := d.Pool.QueryRow(ctx, query, args...)
	err = row.Scan(&order.Id, &order.UserId, &order.Address, &order.Phone, &order.Comment, &order.Status, &order.Notes)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return nil, err
	}
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - GetOrder - QueryRow")
		return nil, err
	}
	return order, nil

}

func (d *PgxAccess) UpdateOrder(ctx context.Context, order *entity.Order) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.UpdateOrder"})
	queryStr := `UPDATE public."orders" AS o
	SET
	address = COALESCE(NULLIF($3, ''), o.address),
	phone = COALESCE(NULLIF($4, ''), o.phone),
	comment = COALESCE(NULLIF($5, ''), o.comment),
	notes = COALESCE(NULLIF($6, ''), o.notes),
	update_ts = COALESCE($7, o.update_ts),
	version = o.version + 1
	FROM public."users" AS u
	WHERE o.id = $1 AND o.user_id = $2 OR u.role = 'ADMIN'
	RETURNING o.id;`
	dbLog.Info(order.UserId)
	row := d.Pool.QueryRow(ctx, queryStr, order.Id, order.UserId, order.Address, order.Phone, order.Comment, order.Notes, order.UpdateTs)
	err := row.Scan(&order.Id)
	if err == pgx.ErrNoRows {
		return errorStatus.ErrNotFound
	}
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - UpdateOrder - Exec")
		return err
	}
	return nil
}

func (d *PgxAccess) UpdateOrderStatus(ctx context.Context, order *entity.Order) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.UpdateOrderStatus"})
	query, args, err := d.Builder.
		Update("orders").
		Set("status", order.Status).
		Set("update_ts", order.UpdateTs).
		Set("version", squirrel.Expr("version+1")).
		Where("orders.id = $1", order.Id).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - UpdateOrderStatus - r.Builder - query")
		return err
	}
	row := d.Pool.QueryRow(ctx, query, args...)
	err = row.Scan(&order.Id)
	if err == pgx.ErrNoRows {
		return errorStatus.ErrNotFound
	}
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - UpdateOrderStatus - Exec")
		return err
	}
	return nil
}

func (d *PgxAccess) DeleteOrder(ctx context.Context, order *entity.Order) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.DeleteOrder"})
	query, args, err := d.Builder.
		Update("orders").
		Set("state", entity.Deleted).
		Set("update_ts", order.UpdateTs).
		Set("version", squirrel.Expr("version+1")).
		Where("orders.id = $3", order.Id).
		Suffix("RETURNING \"id\"").
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - DeleteOrder - r.Builder - query")
		return err
	}
	row := d.Pool.QueryRow(ctx, query, args...)
	err = row.Scan(&order.Id)
	if err == pgx.ErrNoRows {
		return errorStatus.ErrNotFound
	}
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - DeleteOrder - Exec")
		return err
	}
	return nil
}
