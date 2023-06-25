package model

import (
	"github.com/google/uuid"
	"time"

	v1directorypb "github.com/kirychukyurii/wasker/gen/go/directory/v1"
	"github.com/kirychukyurii/wasker/internal/model/dto"
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
		dto.Model

		Name string `db:"name"`
	}

	RolePermission struct {
		User  dto.LookupEntity `db:"user"`
		Role  dto.LookupEntity `db:"role"`
		Scope dto.LookupEntity `db:"scope"`

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
		Pagination dto.PaginationParam
		Order      dto.OrderParam
		Query      dto.QueryParam
	}

	ScopeQueryResult struct {
		List       Scopes          `json:"list"`
		Pagination *dto.Pagination `json:"pagination"`
	}

	ScopeEndpoint struct {
		Id   uint64 `db:"id"`
		Name string `db:"name"`
		Bit  uint8  `db:"bit"`

		Scope dto.LookupEntity `db:"scope"`
	}

	ScopeEndpoints []*ScopeEndpoint

	ScopeEndpointQueryParam struct {
		Pagination dto.PaginationParam
		Order      dto.OrderParam
		Query      dto.QueryParam
	}

	ScopeEndpointQueryResult struct {
		List       ScopeEndpoints  `json:"list"`
		Pagination *dto.Pagination `json:"pagination"`
	}
)

type (
	UserLogin struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	UserSession struct {
		Id          uint64           `json:"id" db:"id"`
		User        dto.LookupEntity `json:"user" db:"user"`
		NetworkIp   string           `json:"network_ip" db:"network_ip"`
		AccessToken string           `json:"access_token" db:"access_token"`
		CreatedAt   time.Time        `json:"created_at" db:"created_at"`
		ExpiresAt   time.Time        `json:"expires_at" db:"expires_at"`
	}

	UserSessionInfo struct {
		User      dto.LookupEntity `json:"user" db:"user"`
		Role      dto.LookupEntity `json:"role" db:"role"`
		ExpiresAt time.Time        `json:"expires_at" db:"expires_at"`
	}
)

func NewID() string {
	return uuid.New().String()
}
