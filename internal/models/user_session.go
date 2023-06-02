package models

import (
	"github.com/kirychukyurii/wasker/internal/models/dto"
	"time"
)

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserSession struct {
	Id          uint             `json:"id" db:"id"`
	User        dto.LookupEntity `json:"user" db:"user"`
	NetworkIp   string           `json:"network_ip" db:"network_ip"`
	AccessToken string           `json:"access_token" db:"access_token"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	ExpiresAt   time.Time        `json:"expires_at" db:"expires_at"`
}

type UserSessionInfo struct {
	User      dto.LookupEntity  `json:"user" db:"user"`
	Role      dto.LookupEntity  `json:"role" db:"role"`
	Scope     *ScopePermissions `json:"scope" db:"scope"`
	ExpiresAt time.Time         `json:"expires_at" db:"expires_at"`
}
