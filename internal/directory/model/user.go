package model

import (
	"github.com/kirychukyurii/wasker/internal/model"
)

type User struct {
	model.Model

	Name     string                      `json:"name" db:"name"`
	Email    string                      `json:"email" validate:"required,email" db:"email"`
	UserName string                      `json:"username" validate:"required" db:"user_name"`
	Password string                      `json:"password" validate:"required" db:"password"`
	Role     *model.NullableLookupEntity `json:"role" db:"role"`
}

type Users []*User

type UserQueryParam struct {
	Pagination model.PaginationParam
	Order      model.OrderParam
	Query      model.QueryParam

	UserName string `query:"user_name"`
}

type UserQueryResult struct {
	List       Users             `json:"list"`
	Pagination *model.Pagination `json:"pagination"`
}

func (a User) TableName() string {
	return "auth_user"
}

func (a UserQueryResult) GetPassword() string {
	if len(a.List) < 1 {
		return ""
	}

	return a.List[0].Password
}
