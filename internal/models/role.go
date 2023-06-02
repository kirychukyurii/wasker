package models

import "github.com/kirychukyurii/wasker/internal/models/dto"

type Role struct {
	dto.Model

	Name string `db:"name"`
}

type RolePermission struct {
	User  dto.LookupEntity `db:"user"`
	Role  dto.LookupEntity `db:"role"`
	Scope dto.LookupEntity `db:"scope"`

	Permission string `db:"permission"`
	Acl        bool   `db:"acl"`
}
