package model

import (
	"time"
)

type Model struct {
	Id int64 `json:"id" db:"id"`

	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	CreatedBy LookupEntity `json:"created_by" db:"created_by"`

	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	UpdatedBy LookupEntity `json:"updated_by" db:"updated_by"`

	DeletedAt *time.Time    `json:"-" db:"-"`
	DeletedBy *LookupEntity `json:"-" db:"-"`
}

type LookupEntity struct {
	Id   int64  `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
}

type NullableLookupEntity struct {
	Id   *int64  `json:"id,omitempty" db:"id"`
	Name *string `json:"name,omitempty" db:"name"`
}
