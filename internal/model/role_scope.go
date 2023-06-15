package model

import "github.com/kirychukyurii/wasker/internal/model/dto"

var (
	CreatePermission = "create"
	ReadPermission   = "read"
	UpdatePermission = "update"
	DeletePermission = "delete"
)

var (
	UserScope = "user"
	RoleScope = "role"
)

type ScopePermission struct {
	Action ScopeAction `json:"action" db:"action"`

	Name string `json:"name" db:"name"`
}

type ScopeAction struct {
	Create bool `json:"create" db:"create"`
	Read   bool `json:"read" db:"read"`
	Update bool `json:"update" db:"update"`
	Delete bool `json:"delete" db:"delete"`
}

type ScopePermissions []ScopePermission

type Scope struct {
	dto.Model
	ScopePermissions

	Name string `db:"name"`
}
