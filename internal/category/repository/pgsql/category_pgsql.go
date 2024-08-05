package pgsql

import (
	"context"
	"go-store/internal/entity"

	"github.com/Masterminds/squirrel"
	log "github.com/sirupsen/logrus"

	"go-store/utils/database"
)

type PgxAccess struct {
	*database.PgxAccess
}

func NewPgxCategoryRepository(pgx *database.PgxAccess) entity.CategoryRepository {
	return &PgxAccess{pgx}
}

func (d *PgxAccess) Get(ctx context.Context) (result []*entity.Category, err error) {
	dbLog := log.WithFields(log.Fields{"func": "CategoryRepository.Get"})
	query, args, err := d.Builder.
		Select("id",
			"name",
			"parent",
			"image",
			"icon",
			"create_ts",
			"update_ts",
			"state",
			"version").
		From("category").
		Where("state = 'enabled'").
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("CategoryRepository - Get - r.Builder - query")
		return nil, err
	}
	rows, err := d.Pool.Query(context.Background(), query, args...)
	if err != nil {
		dbLog.WithError(err).Warning("d.pool.Query")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tmp entity.Category
		if err := rows.Scan(&tmp.Id, &tmp.Name, &tmp.Parent, &tmp.Image, &tmp.Icon, &tmp.CreateTs, &tmp.UpdateTs, &tmp.State, &tmp.Version); err != nil {
			dbLog.WithError(err).Warning("rows.Scan")
			return nil, err
		}
		result = append(result, &tmp)

	}
	return result, nil

}
func (d *PgxAccess) GetById(ctx context.Context, categoryId int) (result *entity.Category, err error) {
	dbLog := log.WithFields(log.Fields{"func": "CategoryRepository.GetById"})
	category := &entity.Category{}
	query, args, err := d.Builder.
		Select("id",
			"name",
			"parent",
			"image",
			"icon",
			"create_ts",
			"update_ts",
			"state",
			"version").
		From("category").
		Where("state = 'enabled' AND id = ?", categoryId).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("CategoryRepository - GetById - r.Builder - query")
		return nil, err
	}
	row := d.Pool.QueryRow(ctx, query, args...)
	err = row.Scan(&category.Id, &category.Name, &category.Parent, &category.Image, &category.Icon, &category.CreateTs, &category.UpdateTs, &category.State, &category.Version)
	if err != nil {
		dbLog.Warning(err)
		return nil, err
	}
	return category, nil
}
func (d *PgxAccess) Create(ctx context.Context, category *entity.Category) error {
	dbLog := log.WithFields(log.Fields{"func": "CategoryRepository.Create"})
	query, args, err := d.Builder.
		Insert("category").
		Columns("name",
			"parent",
			"icon",
			"image",
			"create_ts",
			"update_ts",
			"state",
			"version").
		Values(category.Name,
			category.Parent,
			category.Icon,
			category.Image,
			category.CreateTs,
			category.UpdateTs,
			category.State,
			category.Version).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("CategoryRepository - Create - r.Builder - query")
		return err
	}
	_, err = d.Pool.Exec(ctx, query, args...)
	if err != nil {
		dbLog.Warning(err)
		return err
	}
	return nil
}

func (d *PgxAccess) Update(ctx context.Context, category *entity.Category) error {
	dbLog := log.WithFields(log.Fields{"func": "CategoryRepository.Update"})
	query, args, err := d.Builder.
		Update("category").
		SetMap(map[string]interface{}{
			"name":      category.Name,
			"parent":    category.Parent,
			"icon":      category.Icon,
			"image":     category.Image,
			"update_ts": category.UpdateTs,
			"state":     category.State}).
		Set("version", squirrel.Expr("version+1")).
		Where("category.id = $7", category.Id).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("CategoryRepository - Update - r.Builder - query")
		return err
	}
	dbLog.Info(query)
	dbLog.Info(args...)
	_, err = d.Pool.Exec(ctx, query, args...)
	if err != nil {
		dbLog.WithError(err).Errorf("CategoryRepository - Update - r.Builder - query")
		return err
	}
	return nil
}

func (d *PgxAccess) Delete(ctx context.Context, categoryId int) error {
	dbLog := log.WithFields(log.Fields{"func": "CategoryRepository.Delete"})
	query, args, err := d.Builder.
		Update("category").
		Set("state", entity.Deleted).
		Set("version", squirrel.Expr("version+1")).
		Where("category.id = $2", categoryId).
		ToSql()
	_, err = d.Pool.Exec(ctx, query, args...)
	if err != nil {
		dbLog.WithError(err).Errorf("CategoryRepository - Update - r.Builder - query")
		return err
	}
	return nil
}
