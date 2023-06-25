package model

import (
	"github.com/kirychukyurii/wasker/internal/model"
)

type Customer struct {
	model.Model
	Name string `json:"name" validate:"required" db:"name"`
}

type Customers []*Customer

type CustomerQueryParam struct {
	model.PaginationParam
	model.OrderParam
	Name  string `query:"name"`
	Query string `query:"query"`
}

type CustomerQueryResult struct {
	List       Types             `json:"list"`
	Pagination *model.Pagination `json:"pagination"`
}
