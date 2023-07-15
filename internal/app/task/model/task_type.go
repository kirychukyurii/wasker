package model

import (
	"github.com/kirychukyurii/wasker/internal/model"
)

type Type struct {
	model.Model
	Name string `json:"name" validate:"required" db:"name"`
}

type Types []*Type

type TypeQueryParam struct {
	model.PaginationParam
	model.OrderParam
	Name  string `query:"name"`
	Query string `query:"query"`
}

type TypeQueryResult struct {
	List       Types             `json:"list"`
	Pagination *model.Pagination `json:"pagination"`
}
