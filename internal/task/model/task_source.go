package model

import (
	"github.com/kirychukyurii/wasker/internal/model"
)

type Source struct {
	model.Model
	Name string `json:"name" validate:"required" db:"name"`
}

type Sources []*Source

type SourceQueryParam struct {
	model.PaginationParam
	model.OrderParam
	Name  string `query:"name"`
	Query string `query:"query"`
}

type SourceQueryResult struct {
	List       Types             `json:"list"`
	Pagination *model.Pagination `json:"pagination"`
}
