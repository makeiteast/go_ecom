package entity

import (
	"context"
	"time"
)

type Category struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Parent   int       `json:"parent"`
	Image    string    `json:"image"`
	Icon     string    `json:"icon"`
	CreateTs time.Time `json:"createTs"`
	UpdateTs time.Time `json:"updateTs"`
	State    State     `json:"state"`
	Version  int       `json:"version"`
}

type CategoryJson struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Parent int    `json:"parent"`
	Image  string `json:"image"`
	Icon   string `json:"icon"`
}

type SingleCategoryJson struct {
	Id              int           `json:"id"`
	Name            string        `json:"name"`
	Parent          int           `json:"parent"`
	Image           string        `json:"image"`
	Icon            string        `json:"icon"`
	CategoryOptions []*OptionJson `json:"options"`
}

type CategoryUsecase interface {
	Get(ctx context.Context) (result []*CategoryJson, err error)
	GetById(ctx context.Context, categoryId int, userRole *UserRole) (result *SingleCategoryJson, err error)
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, categoryId int) error
	CreateOpt(ctx context.Context, categoryId int, optionName string) (optionId *int, err error)
	CreateOptValue(ctx context.Context, optionId int, optionValueName string) error
	UpdateCatOpt(c context.Context, categoryIdInt int, optionIdInt int, optionName string) error
	UpdateCatOptValue(c context.Context, optionId int, optionValueIdInt int, optionValueName string) error
	DeleteCatOpt(c context.Context, optionId int) error
	DeleteCatOptValue(c context.Context, optionValueId int) error
}

type CategoryRepository interface {
	Get(ctx context.Context) (result []*Category, err error)
	GetById(ctx context.Context, categoryId int) (result *Category, err error)
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, categoryId int) error
}

type CategoryRedisRepository interface {
	Get(ctx context.Context) (result []*CategoryJson, err error)
}
