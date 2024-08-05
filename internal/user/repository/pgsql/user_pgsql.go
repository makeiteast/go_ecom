package postgsql_repo

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"go-store/internal/entity"
	"go-store/utils/database"
	errorStatus "go-store/utils/errors"
)

type PgxAccess struct {
	*database.PgxAccess
}

func NewPgxAuthPgxRepository(pgx *database.PgxAccess) entity.AuthPgxRepository {
	return &PgxAccess{pgx}
}

func (d *PgxAccess) UserByUsername(ctx context.Context, username string) (result *entity.Users, err error) {
	dbLog := log.WithFields(log.Fields{"func": "db.GetLoginUser"})
	user := &entity.Users{}
	query, args, err := d.Builder.
		Select("users.id",
			"users.public_id",
			"users.username",
			"users.password",
			"users.role").
		From("users").
		Where("users.username = $1", username).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("SourceRepo - GetById - r.Builder")
		return nil, err
	}
	if err := d.Pool.QueryRow(context.Background(), query, args...).Scan(&user.Id, &user.PublicId, &user.Username, &user.Password, &user.Role); err != nil {
		dbLog.WithFields(log.Fields{"user_id": user.Id}).Warning(err)
		return nil, err
	}
	return user, nil

}

func (d *PgxAccess) Create(ctx context.Context, user *entity.Users) (err error) {
	dbLog := log.WithFields(log.Fields{"func": "db.GetLoginUser"})
	query, args, err := d.Builder.
		Insert("users").
		Columns("public_id",
			"username",
			"password",
			"email",
			"phone_number",
			"address",
			"photo",
			"role",
			"region_id",
			"parent",
			"create_ts",
			"update_ts",
			"state",
			"version",
			"fullname",
			"verification_code").
		Values(user.PublicId,
			user.Username,
			user.Password,
			user.Email,
			user.PhoneNumber,
			user.Address,
			user.Photo,
			user.Role,
			user.RegionId,
			user.Parent,
			user.CreateTs,
			user.UpdateTs,
			user.State,
			user.Version,
			user.FullName,
			user.VerificationCode).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("SourceRepo - GetById - r.Builder")
		return err
	}
	_, err = d.Pool.Exec(context.Background(), query, args...)

	if err != nil {
		dbLog.WithFields(log.Fields{"user_id": user.Username}).Warning(err)
		return err
	}
	return nil

}

func (d *PgxAccess) Update(ctx context.Context, user *entity.Users) error {
	dbLog := log.WithFields(log.Fields{"func": "pg.Update"})
	query, args, err := d.Builder.
		Update("users").
		SetMap(map[string]interface{}{
			"public_id":         user.PublicId,
			"username":          user.Username,
			"fullname":          user.FullName,
			"password":          user.Password,
			"email":             user.Email,
			"phone_number":      user.PhoneNumber,
			"address":           user.Address,
			"photo":             user.Photo,
			"user_role":         user.Role,
			"region_id":         user.RegionId,
			"parent":            user.Parent,
			"verification_code": user.VerificationCode,
			"updateTs":          user.UpdateTs,
			"state":             user.State}).
		Set("version", squirrel.Expr("version+1")).
		Where("users.id = $1", user.Id).
		ToSql()
	if err != nil {
		dbLog.WithError(err).Errorf("UserLogRepo - Update - r.Builder - query")
		return err
	}
	_, err = d.Pool.Exec(ctx, query, args...)
	if err == pgx.ErrNoRows {
		err = errorStatus.ErrNotFound
		return err
	}
	if err != nil {
		dbLog.Warning(err)
		return err
	}
	return nil
}
