package model

import (
	"github.com/kirychukyurii/wasker/internal/model"
)

type Task struct {
	model.Model
	Name        string             `json:"name" validate:"required" db:"name"`
	Description string             `json:"description" db:"description"`
	Link        string             `json:"link" db:"link"`
	Type        model.LookupEntity `json:"type" db:"type"`
	Source      model.LookupEntity `json:"source" db:"source"`
	Customer    model.LookupEntity `json:"customer" db:"customer"`
}

type Tasks []*Task

type TaskQueryParam struct {
	model.PaginationParam
	model.OrderParam
	Name        string `query:"name"`
	Description string `query:"description"`
	Link        string `query:"link"`
	Query       string `query:"query"`
}

type TaskQueryResult struct {
	List       Tasks             `json:"list"`
	Pagination *model.Pagination `json:"pagination"`
}
