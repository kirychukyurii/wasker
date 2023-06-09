package model

import "github.com/kirychukyurii/wasker/internal/model/dto"

type Type struct {
	dto.Model
	Name string `json:"name" validate:"required" db:"name"`
}

type Types []*Type

type TypeQueryParam struct {
	dto.PaginationParam
	dto.OrderParam
	Name  string `query:"name"`
	Query string `query:"query"`
}

type TypeQueryResult struct {
	List       Types           `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}
