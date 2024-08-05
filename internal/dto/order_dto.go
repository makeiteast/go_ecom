package dto

import (
	"strconv"
)

// OrderListRequest -.
type OrderListRequest struct {
	Filter *OrderListFilter `json:"filter"`
	Limit  int              `json:"limit" example:"10"`
	Offset int              `json:"offset" example:"0"`
}

type OrderListFilter struct {
	Id     *int   `json:"id,omitempty" example:"1"`
	UserId *int   `json:"userId,omitempty" example:"1"`
	Phone  string `json:"phone,omitempty"`
	Status string `json:"status,omitempty"`
}

func (o *OrderListFilter) ToSqlFilterMap() map[string]string {
	filterMap := map[string]string{}
	if o.Id != nil {
		filterMap["id"] = strconv.Itoa(*o.Id)
	}
	if o.UserId != nil {
		filterMap["user_id"] = strconv.Itoa(*o.UserId)
	}
	if len(o.Phone) > 0 {
		filterMap["phone"] = o.Phone
	}
	if len(o.Status) > 0 {
		filterMap["status"] = o.Status
	}
	return filterMap
}
