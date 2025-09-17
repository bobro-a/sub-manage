package models

import (
	"github.com/google/uuid"
)

type Sub struct {
	ID          int64     `db:"id" json:"id"`
	ServiceName string    `db:"service_name" json:"service_name"`
	Price       int32     `db:"price" json:"price"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	StartDate   MonthYear `db:"start_date" json:"start_date"`
	EndDate     MonthYear `db:"end_date" json:"end_date"`
}
