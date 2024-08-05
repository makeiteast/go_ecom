package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"go-store/internal/entity"
	errorStatus "go-store/utils/errors"
)

// CategoryUsecase will initiate usecase of entity.CategoryRepository interface
type CategoryUsecase struct {
	catRepo    entity.CategoryRepository
	optionRepo entity.OptionRepository
}

// NewCategoryUsecase will create new an CategoryUsecase object representation of entity.CategoryUsecase interface
func NewCategoryUsecase(c entity.CategoryRepository, o entity.OptionRepository) entity.CategoryUsecase {
	return &CategoryUsecase{
		catRepo:    c,
		optionRepo: o,
	}
}

func (o *CategoryUsecase) Get(ctx context.Context) (result []*entity.CategoryJson, err error) {
	srvLog := log.WithFields(log.Fields{"func": "CategoryUsecase.Get"})

	categories, err := o.catRepo.Get(ctx)
	if err != nil {
		srvLog.Warning("Cannot get category query, Err: ", err)
		err = errorStatus.ErrInternalServer
		return
	}

	categoriesMap := make([]*entity.CategoryJson, len(categories))
	for idx, category := range categories {
		categoriesMap[idx] = mapCatToJSON(category)
	}

	return categoriesMap, nil
}

func mapCatToJSON(s *entity.Category) *entity.CategoryJson {
	return &entity.CategoryJson{
		Id:     s.Id,
		Name:   s.Name,
		Parent: s.Parent,
		Image:  s.Image,
		Icon:   s.Icon,
	}
}

func (o *CategoryUsecase) GetById(ctx context.Context, categoryId int, userRole *entity.UserRole) (result *entity.SingleCategoryJson, err error) {
	srvLog := log.WithFields(log.Fields{"func": "CategoryUsecase.GetById"})

	category, err := o.catRepo.GetById(ctx, categoryId)
	if err == pgx.ErrNoRows {
		srvLog.Warning(err)
		err = errorStatus.ErrNotFound
		return nil, err
	}
	if err != nil {
		srvLog.Warning("Cannot get category query, Err: ", err)
		err = errorStatus.ErrInternalServer
		return
	}

	options, err := o.optionRepo.GetOptionByCat(ctx, category.Id, userRole)
	if err != nil {
		srvLog.Warning("Cannot get option query, Err: ", err)
		err = errorStatus.ErrInternalServer
		return
	}
	optionMap := make([]*entity.OptionJson, len(options))
	for idx, option := range options {
		optionMap[idx] = mapOptToJSON(option)
		optionValues, err := o.optionRepo.GetOptionValueByOptId(ctx, option.Id, userRole)
		if err != nil {
			srvLog.Warning("Cannot get option query, Err: ", err)
			err = errorStatus.ErrInternalServer
			return nil, err
		}
		optionMap[idx].OptionValueJson = optionValues
	}
	result = mapCatIdToJSON(category)
	result.CategoryOptions = optionMap

	return result, nil
}

func mapCatIdToJSON(s *entity.Category) *entity.SingleCategoryJson {
	return &entity.SingleCategoryJson{
		Id:     s.Id,
		Name:   s.Name,
		Parent: s.Parent,
		Image:  s.Image,
		Icon:   s.Icon,
	}
}

func mapOptToJSON(s *entity.Option) *entity.OptionJson {
	return &entity.OptionJson{
		Id:   s.Id,
		Name: s.Name,
	}
}

func mapOptValToJSON(s *entity.OptionValue) *entity.OptionValueJson {
	return &entity.OptionValueJson{
		Id:   s.Id,
		Name: s.Name,
	}
}

func (o *CategoryUsecase) Create(ctx context.Context, category *entity.Category) error {
	ctLog := log.WithFields(log.Fields{"func": "CategoryUsecase.Create"})
	category.CreateTs = time.Now().UTC()
	category.UpdateTs = time.Now().UTC()
	category.State = entity.Enabled
	category.Version = 0
	err := o.catRepo.Create(ctx, category)
	if err != nil {
		ctLog.WithError(err).Warning("o.catRepo.Create")
		return err
	}

	return nil
}

func (o *CategoryUsecase) Update(ctx context.Context, category *entity.Category) error {
	ctLog := log.WithFields(log.Fields{"func": "CategoryUsecase.Update"})
	category.UpdateTs = time.Now().UTC()
	err := o.catRepo.Update(ctx, category)
	if err != nil {
		ctLog.WithError(err).Warning("o.catRepo.Update")
		return err
	}
	return nil
}
func (o *CategoryUsecase) Delete(ctx context.Context, categoryId int) error {
	ctLog := log.WithFields(log.Fields{"func": "CategoryUsecase.Delete"})

	err := o.catRepo.Delete(ctx, categoryId)
	if err != nil {
		ctLog.WithError(err).Warning("o.catRepo.Delete")
		return err
	}
	return nil
}

func (o *CategoryUsecase) CreateOpt(ctx context.Context, categoryId int, optionName string) (optionId *int, err error) {
	ctLog := log.WithFields(log.Fields{"func": "CategoryUsecase.CreateOpt"})
	option := entity.Option{
		Name:       optionName,
		CategoryId: categoryId,
		CreateTs:   time.Now().UTC(),
		UpdateTs:   time.Now().UTC(),
		State:      entity.Enabled,
		Version:    0,
	}
	optionId, err = o.optionRepo.CreateOption(ctx, option)
	if err != nil {
		pqErr := err.(*pgconn.PgError)
		if strings.Contains(pqErr.Code, "23505") {
			ctLog.WithError(err).Warning("o.optionRepo.CreateOption")
			return nil, errorStatus.ErrBadReq
		} else {
			ctLog.WithError(err).Warning("o.optionRepo.CreateOption")
			return nil, err
		}
	}

	return optionId, nil
}

func (o *CategoryUsecase) CreateOptValue(ctx context.Context, optionId int, optionValueName string) error {
	ctLog := log.WithFields(log.Fields{"func": "CategoryUsecase.CreateOptValue"})

	var txId int
	txId, err := o.optionRepo.NewTxId(ctx)
	if err != nil {
		ctLog.WithError(err).Warning("o.optionRepo.NewTxId")
		return err
	}
	defer func() {
		err = o.optionRepo.TxEnd(ctx, txId, err)
		if err != nil {
			ctLog.WithError(err).Warning("o.orderRepo.TxEnd")
			return
		}
	}()

	optionValue := entity.OptionValue{
		Name:     optionValueName,
		OptionId: optionId,
		CreateTs: time.Now().UTC(),
		UpdateTs: time.Now().UTC(),
		State:    entity.Enabled,
		Version:  0,
	}
	_, err = o.optionRepo.CreateOptionValue(ctx, optionValue, txId)
	if err != nil {
		pqErr := err.(*pgconn.PgError)
		if strings.Contains(pqErr.Code, "23503") {
			ctLog.WithError(err).Warning("o.optionRepo.CreateOptValue")
			return errorStatus.ErrBadReq
		} else {
			ctLog.WithError(err).Warning("o.optionRepo.CreateOptValue")
			return err
		}

	}

	return nil
}

func (o *CategoryUsecase) UpdateCatOpt(c context.Context, categoryIdInt int, optionIdInt int, optionName string) error {
	ctLog := log.WithFields(log.Fields{"func": "CategoryUsecase.UpdateCatOpt"})
	option := entity.Option{
		Id:         optionIdInt,
		Name:       optionName,
		CategoryId: categoryIdInt,
		UpdateTs:   time.Now().UTC(),
	}
	err := o.optionRepo.UpdateOption(c, option)
	if err != nil {
		ctLog.WithError(err).Warning("o.optionRepo.UpdateOption")
		return err
	}
	return nil
}
func (o *CategoryUsecase) UpdateCatOptValue(c context.Context, optionId int, optionValueIdInt int, optionValueName string) error {
	ctLog := log.WithFields(log.Fields{"func": "CategoryUsecase.UpdateCatOptValue"})
	optionValue := entity.OptionValue{
		Id:       optionValueIdInt,
		Name:     optionValueName,
		OptionId: optionId,
		UpdateTs: time.Now().UTC(),
	}
	err := o.optionRepo.UpdateOptionValue(c, optionValue)
	if err != nil {
		ctLog.WithError(err).Warning("o.optionRepo.UpdateOptionValue")
		return err
	}
	return nil
}
func (o *CategoryUsecase) DeleteCatOpt(c context.Context, optionId int) error {
	ctLog := log.WithFields(log.Fields{"func": "CategoryUsecase.DeleteCatOpt"})
	err := o.optionRepo.RemoveOption(c, optionId)
	if err != nil {
		ctLog.WithError(err).Warning("o.optionRepo.DeleteCatOpt")
		return err
	}
	return nil
}
func (o *CategoryUsecase) DeleteCatOptValue(c context.Context, optionValueId int) error {
	ctLog := log.WithFields(log.Fields{"func": "CategoryUsecase.DeleteCatOptValue"})
	err := o.optionRepo.RemoveOptionValue(c, optionValueId)
	if err != nil {
		ctLog.WithError(err).Warning("o.optionRepo.RemoveOptionValue")
		return err
	}
	return nil
}
