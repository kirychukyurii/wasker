package model

import (
	"github.com/google/uuid"
	"github.com/kirychukyurii/wasker/internal/model"
	"time"

	v1directorypb "github.com/kirychukyurii/wasker/gen/go/directory/v1"
)

const (
	AccessRead   = 4
	AccessWrite  = 2
	AccessDelete = 1
)

var (
	DefaultPermission = []struct {
		Scope    string
		Endpoint ScopeEndpoints
	}{
		{
			Scope: v1directorypb.UserService_ServiceDesc.ServiceName,
			Endpoint: ScopeEndpoints{
				{Name: v1directorypb.UserService_CreateUser_FullMethodName, Bit: AccessWrite},
				{Name: v1directorypb.UserService_ReadUser_FullMethodName, Bit: AccessRead},
				{Name: v1directorypb.UserService_UpdateUser_FullMethodName, Bit: AccessWrite},
				{Name: v1directorypb.UserService_DeleteUsers_FullMethodName, Bit: AccessDelete},
				{Name: v1directorypb.UserService_ListUsers_FullMethodName, Bit: AccessRead},
			},
		},
		{
			Scope: v1directorypb.RoleService_ServiceDesc.ServiceName,
			Endpoint: ScopeEndpoints{
				{Name: v1directorypb.RoleService_CreateRole_FullMethodName, Bit: AccessWrite},
				{Name: v1directorypb.RoleService_ReadRole_FullMethodName, Bit: AccessRead},
				{Name: v1directorypb.RoleService_UpdateRole_FullMethodName, Bit: AccessWrite},
				{Name: v1directorypb.RoleService_DeleteRoles_FullMethodName, Bit: AccessDelete},
				{Name: v1directorypb.RoleService_SearchRoles_FullMethodName, Bit: AccessRead},
			},
		},
	}
)

type (
	Role struct {
		model.Model

		Name string `db:"name"`
	}

	RolePermission struct {
		User  model.LookupEntity `db:"user"`
		Role  model.LookupEntity `db:"role"`
		Scope model.LookupEntity `db:"scope"`

		Permission string `db:"permission"`
		Acl        bool   `db:"acl"`
	}
)

type (
	Scope struct {
		Id   uint64 `db:"id"`
		Name string `db:"name"`
	}

	Scopes []*Scope

	ScopeQueryParam struct {
		Pagination model.PaginationParam
		Order      model.OrderParam
		Query      model.QueryParam
	}

	ScopeQueryResult struct {
		List       Scopes            `json:"list"`
		Pagination *model.Pagination `json:"pagination"`
	}

	ScopeEndpoint struct {
		Id   uint64 `db:"id"`
		Name string `db:"name"`
		Bit  uint8  `db:"bit"`

		Scope model.LookupEntity `db:"scope"`
	}

	ScopeEndpoints []*ScopeEndpoint

	ScopeEndpointQueryParam struct {
		Pagination model.PaginationParam
		Order      model.OrderParam
		Query      model.QueryParam
	}

	ScopeEndpointQueryResult struct {
		List       ScopeEndpoints    `json:"list"`
		Pagination *model.Pagination `json:"pagination"`
	}
)

type (
	UserLogin struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	UserSession struct {
		Id          uint64             `json:"id" db:"id"`
		User        model.LookupEntity `json:"user" db:"user"`
		NetworkIp   string             `json:"network_ip" db:"network_ip"`
		AccessToken string             `json:"access_token" db:"access_token"`
		CreatedAt   time.Time          `json:"created_at" db:"created_at"`
		ExpiresAt   time.Time          `json:"expires_at" db:"expires_at"`
	}

	UserSessionInfo struct {
		User      model.LookupEntity `json:"user" db:"user"`
		Role      model.LookupEntity `json:"role" db:"role"`
		ExpiresAt time.Time          `json:"expires_at" db:"expires_at"`
	}
)

func NewID() string {
	return uuid.New().String()
}
