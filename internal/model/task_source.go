package model

import "github.com/kirychukyurii/wasker/internal/model/dto"

type Source struct {
	dto.Model
	Name string `json:"name" validate:"required" db:"name"`
}

type Sources []*Source

type SourceQueryParam struct {
	dto.PaginationParam
	dto.OrderParam
	Name  string `query:"name"`
	Query string `query:"query"`
}

type SourceQueryResult struct {
	List       Types           `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}
