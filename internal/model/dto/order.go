package dto

import "fmt"

type OrderDirection string

const (
	OrderByASC  OrderDirection = "ASC"
	OrderByDESC OrderDirection = "DESC"

	OrderDefaultKey = "id"
)

type OrderParam struct {
	Key       string         `json:"order_key" query:"order_key"`
	Direction OrderDirection `json:"order_direction" query:"order_direction"`
}

func (a OrderParam) ParseOrder() string {
	if a.Key == "" {
		a.Key = OrderDefaultKey
	}

	key := a.Key
	direction := "DESC"
	if a.Direction == OrderByASC {
		direction = "ASC"
	}

	return fmt.Sprintf("%s %s", key, direction)
}
