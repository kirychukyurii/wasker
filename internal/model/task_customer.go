package model

import "github.com/kirychukyurii/wasker/internal/model/dto"

type Customer struct {
	dto.Model
	Name string `json:"name" validate:"required" db:"name"`
}

type Customers []*Customer

type CustomerQueryParam struct {
	dto.PaginationParam
	dto.OrderParam
	Name  string `query:"name"`
	Query string `query:"query"`
}

type CustomerQueryResult struct {
	List       Types           `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}
