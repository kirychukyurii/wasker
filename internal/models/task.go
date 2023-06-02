package models

import "github.com/kirychukyurii/wasker/internal/models/dto"

type Task struct {
	dto.Model
	Name        string           `json:"name" validate:"required" db:"name"`
	Description string           `json:"description" db:"description"`
	Link        string           `json:"link" db:"link"`
	Type        dto.LookupEntity `json:"type" db:"type"`
	Source      dto.LookupEntity `json:"source" db:"source"`
	Customer    dto.LookupEntity `json:"customer" db:"customer"`
}

type Tasks []*Task

type TaskQueryParam struct {
	dto.PaginationParam
	dto.OrderParam
	Name        string `query:"name"`
	Description string `query:"description"`
	Link        string `query:"link"`
	Query       string `query:"query"`
}

type TaskQueryResult struct {
	List       Tasks           `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}
