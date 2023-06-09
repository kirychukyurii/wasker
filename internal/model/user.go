package model

import "github.com/kirychukyurii/wasker/internal/model/dto"

type User struct {
	dto.Model

	Name     string            `json:"name" db:"name"`
	Email    string            `json:"email" validate:"required,email" db:"email"`
	UserName string            `json:"username" validate:"required" db:"user_name"`
	Password string            `json:"password" validate:"required" db:"password"`
	Role     *dto.LookupEntity `json:"role" db:"role"`
}

type Users []*User

type UserQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	Name     string `query:"name"`
	UserName string `query:"user_name"`
	Email    string `query:"email"`
	Query    string `query:"query"`
}

type UserQueryResult struct {
	List       Users           `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a *User) TableName() string {
	return "auth_user"
}
