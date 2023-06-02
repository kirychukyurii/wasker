package dto

import "time"

type Model struct {
	Id uint `json:"id" db:"id"`

	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	CreatedBy LookupEntity `json:"created_by" db:"created_by"`

	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	UpdatedBy LookupEntity `json:"updated_by" db:"updated_by"`

	DeletedAt *time.Time    `json:"-" db:"deleted_at"`
	DeletedBy *LookupEntity `json:"-" db:"deleted_by"`
}

type LookupEntity struct {
	Id   uint   `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
}
