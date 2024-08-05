package entity

import (
	"context"
	"time"
)

type Option struct {
	Id         int       `json:"id"`
	CategoryId int       `json:"categoryId"`
	Name       string    `json:"name"`
	CreateTs   time.Time `json:"createTs"`
	UpdateTs   time.Time `json:"updateTs"`
	State      State     `json:"state"`
	Version    int       `json:"version"`
}

type OptionValue struct {
	Id       int       `json:"id"`
	OptionId int       `json:"optionId"`
	Name     string    `json:"name"`
	CreateTs time.Time `json:"createTs"`
	UpdateTs time.Time `json:"updateTs"`
	State    State     `json:"state"`
	Version  int       `json:"version"`
}

type OptionJson struct {
	Id              int                `json:"id"`
	Name            string             `json:"name"`
	OptionValueJson []*OptionValueJson `json:"optionValues"`
}

type OptionValueJson struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type OptionRepository interface {
	GetSkuValue(ctx context.Context, skuID int, userRole *UserRole) (result []*SkuValue, err error)
	GetOption(ctx context.Context, optionID int, userRole *UserRole) (result *Option, err error)
	GetOptionByCat(ctx context.Context, categoryId int, userRole *UserRole) (result []*Option, err error)
	GetOptionValue(ctx context.Context, optionValueID int, userRole *UserRole) (result *OptionValue, err error)
	GetOptionValueByOptId(ctx context.Context, optionID int, userRole *UserRole) (result []*OptionValueJson, err error)
	GetOptionBySkuValue(ctx context.Context, skuValueID int, userRole *UserRole) (result *OptionJson, err error)
	CreateOption(ctx context.Context, option Option) (optionID *int, err error)
	CreateOptionValue(ctx context.Context, optionValue OptionValue, txId int) (optionValueID *int, err error)
	CreateSkuValue(ctx context.Context, sku string, skuValue *SkuValue, txId int) error
	UpdateOption(ctx context.Context, option Option) error
	UpdateOptionValue(ctx context.Context, optionValue OptionValue) error
	RemoveOption(ctx context.Context, optionID int) error
	RemoveOptionValue(ctx context.Context, optionValueID int) error
	RemoveSkuValue(ctx context.Context, skuValueID int) error
	NewTxId(ctx context.Context) (txId int, err error)
	TxEnd(ctx context.Context, txId int, err error) error
}
