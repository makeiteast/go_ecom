package pgsql

import (
	"context"
	"fmt"
	"go-store/internal/entity"
	"time"

	"go-store/utils/database"
	errorStatus "go-store/utils/errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type PgxAccess struct {
	*database.PgxAccess
}

func NewPgxOptionRepository(pgx *database.PgxAccess) entity.OptionRepository {
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

func (d *PgxAccess) GetSkuValue(ctx context.Context, skuId int, userRole *entity.UserRole) (res []*entity.SkuValue, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetSkuValue"})
	var skuValue []*entity.SkuValue
	rows, err := d.Pool.Query(ctx, GetSkuValueProducts, skuId)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - GetSkuValue - Query")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tmp := &entity.SkuValue{}
		if err := rows.Scan(&tmp.Id, &tmp.OptionId, &tmp.OptionValueId); err != nil {
			dbLog.WithError(err).Errorf("PgxAccess - GetSkuValue - Scan")
			return nil, err
		}
		skuValue = append(skuValue, tmp)
	}

	return skuValue, nil
}

func (d *PgxAccess) GetOption(ctx context.Context, optionId int, userRole *entity.UserRole) (res *entity.Option, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetOption"})
	opt := &entity.Option{}
	row := d.Pool.QueryRow(ctx, GetProductOptions, optionId)
	err = row.Scan(&opt.Id, &opt.Name)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - GetOption - QueryRow")
		return nil, err
	}
	return opt, nil
}

func (d *PgxAccess) GetOptionByCat(ctx context.Context, categoryId int, userRole *entity.UserRole) (res []*entity.Option, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetOptionByCat"})
	opt := []*entity.Option{}

	query, args, err := d.Builder.
		Select("id",
			"name").
		From("option").
		Where("option.category_id = $1 AND state = $2", categoryId, userRole).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - GetSkuByProductID - r.Builder - query")
		return nil, err
	}
	rows, err := d.Pool.Query(ctx, query, args...)
	if err != nil {
		dbLog.Warning(err)
		err = fmt.Errorf("pg.GetProductSku: %w", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tmp := &entity.Option{}
		if err := rows.Scan(&tmp.Id, &tmp.Name); err != nil {
			dbLog.WithError(err).Errorf("PgxAccess - GetOptionByCat - QueryRow")
			return nil, err
		}
		opt = append(opt, tmp)

	}

	return opt, nil
}

func (d *PgxAccess) GetOptionValue(ctx context.Context, optionValueId int, userRole *entity.UserRole) (res *entity.OptionValue, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetOptionValue"})
	optVal := &entity.OptionValue{}
	row := d.Pool.QueryRow(ctx, GetOptionValues, optionValueId)
	err = row.Scan(&optVal.Id, &optVal.Name)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - GetOptionValue - QueryRow")
		return nil, err
	}
	return optVal, nil
}

func (d *PgxAccess) GetOptionValueByOptId(ctx context.Context, optionId int, userRole *entity.UserRole) (res []*entity.OptionValueJson, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetOptionValueByOptId"})
	optVal := []*entity.OptionValueJson{}
	rows, err := d.Pool.Query(ctx, GetOptionValuesByOptId, optionId)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - GetOptionByCat - QueryRow")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tmp := &entity.OptionValueJson{}
		if err := rows.Scan(&tmp.Id, &tmp.Name); err != nil {
			dbLog.WithError(err).Errorf("PgxAccess - GetOptionByCat - QueryRow")
			return nil, err
		}
		optVal = append(optVal, tmp)

	}
	return optVal, nil
}

func (d *PgxAccess) GetOptionBySkuValue(ctx context.Context, skuValueId int, userRole *entity.UserRole) (res *entity.OptionJson, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.GetOption"})
	opt := &entity.Option{}
	optValue := &entity.OptionValue{}
	baseQuery := d.Builder.
		Select("option.id",
			"option.name",
			"option_value.id",
			"option_value.name").
		From("sku_value").
		InnerJoin("option ON sku_value.option_id = option.id").
		InnerJoin("option_value ON sku_value.option_value_id = option_value.id").
		Where("sku_value.id = $1", skuValueId)
	if *userRole != entity.UserRoleAdmin {
		baseQuery = baseQuery.Where("sku_value.state = 'enabled' AND option.state = 'enabled'")
	}

	query, args, err := baseQuery.ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - GetOptionBySkuValue - r.Builder - query")
		return nil, err
	}
	row := d.Pool.QueryRow(ctx, query, args...)
	err = row.Scan(&opt.Id, &opt.Name, &optValue.Id, &optValue.Name)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - GetOptionBySkuValue - QueryRow")
		return nil, err
	}
	option := &entity.OptionJson{
		Id:   opt.Id,
		Name: opt.Name,
	}

	optionValue := []*entity.OptionValueJson{{
		Id:   optValue.Id,
		Name: optValue.Name,
	},
	}
	option.OptionValueJson = optionValue
	return option, nil
}

func (d *PgxAccess) CreateOption(ctx context.Context, option entity.Option) (optionID *int, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.CreateOption"})
	row := d.Pool.QueryRow(ctx, CreateOption, option.CategoryId, option.Name, option.CreateTs, option.UpdateTs, option.State, option.Version)
	err = row.Scan(&optionID)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - CreateOption - QueryRow")
		return nil, err
	}
	return optionID, nil
}

func (d *PgxAccess) CreateOptionValue(ctx context.Context, optionValue entity.OptionValue, txId int) (optionValueID *int, err error) {
	dbLog := log.WithFields(log.Fields{"func": "pg.CreateOptionValue"})

	tx, err := d.GetTxById(txId)
	if err != nil {
		dbLog.WithError(err).Errorf("OptionRepo - CreateOptionValue - d.GetTxById")
		return nil, err
	}

	query, args, err := d.Builder.
		Insert("option_value").
		Columns(
			"option_id",
			"name",
			"create_ts",
			"update_ts",
			"state",
			"version").
		Values(
			optionValue.OptionId,
			optionValue.Name,
			optionValue.CreateTs,
			optionValue.UpdateTs,
			optionValue.State,
			optionValue.Version).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("CreateOptionValue - d.Builder - query")
		return nil, err
	}
	row := tx.QueryRow(ctx, query, args...)
	err = row.Scan(&optionValueID)
	if err != nil {
		dbLog.Warning(err)
		return nil, err
	}
	return optionValueID, nil
}

func (d *PgxAccess) CreateSkuValue(ctx context.Context, sku string, skuValue *entity.SkuValue, txId int) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.CreateSkuValue"})

	tx, err := d.GetTxById(txId)
	if err != nil {
		dbLog.WithError(err).Errorf("OptionRepo - CreateSkuValue - r.GetTxById")
		return err
	}

	skuVal := fmt.Sprintf(`INSERT INTO 
		sku_value(sku_id, 
			option_id, 
			option_value_id, 
			create_ts, 
			update_ts, 
			state, 
			version)
		SELECT id, 
			$1, 
			$2, 
			$3, 
			$4, 
			$5, 
			$6
		FROM sku
		WHERE sku.sku = $7`)
	_, err = tx.Exec(ctx, skuVal, skuValue.OptionId, skuValue.OptionValueId, skuValue.CreateTs, skuValue.UpdateTs, skuValue.State, skuValue.Version, sku)
	if err != nil {
		dbLog.WithError(err).Errorf("OptionRepo - CreateSkuValue - Exec")
		return err
	}

	return nil
}

func (d *PgxAccess) UpdateOption(ctx context.Context, option entity.Option) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.UpdateOption"})
	_, err := d.Pool.Exec(ctx, UpdateOptionName, option.Id, option.Name, option.CategoryId, option.UpdateTs)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return err
	}
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - UpdateOption - Exec")
		return err
	}
	return nil
}

func (d *PgxAccess) UpdateOptionValue(ctx context.Context, optionValue entity.OptionValue) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.UpdateOptionValue"})
	_, err := d.Pool.Exec(ctx, UpdateOptionValueName, optionValue.Id, optionValue.Name, optionValue.OptionId, optionValue.UpdateTs)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return err
	}
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - UpdateOptionValue - Exec")
		return err
	}
	return nil
}

func (d *PgxAccess) RemoveOption(ctx context.Context, optionId int) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.RemoveOption"})
	_, err := d.Pool.Exec(ctx, RemoveOption, optionId)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - RemoveOption - Exec")
		return err
	}
	return nil
}

func (d *PgxAccess) RemoveOptionValue(ctx context.Context, optionValueId int) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.RemoveOptionValue"})
	_, err := d.Pool.Exec(ctx, RemoveOptionValue, optionValueId)
	if err != nil {
		dbLog.WithError(err).Errorf("PgxAccess - RemoveOptionValue - Exec")
		return err
	}
	return nil
}
func (d *PgxAccess) RemoveSkuValue(ctx context.Context, skuValueId int) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.RemoveSkuValue"})
	query, args, err := d.Builder.
		Update("sku_value").
		SetMap(map[string]interface{}{
			"update_ts": time.Now().UTC(),
			"state":     entity.Deleted}).
		Set("version", squirrel.Expr("version+1")).
		Where("sku_value.id = $3", skuValueId).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("RemoveSkuValue - d.Builder - query")
		return err
	}

	_, err = d.Pool.Exec(ctx, query, args...)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return err
	}
	if err != nil {
		dbLog.WithError(err).Errorf("RemoveSkuValue - Exec")
		return err
	}
	return nil
}
