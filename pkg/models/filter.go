package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SumFilter struct {
	UserID      *uuid.UUID `db:"user_id" json:"user_id"`
	ServiceName string     `db:"service_name" json:"service_name"`
	StartDate   *MonthYear `db:"start_date" json:"start_date"`
	EndDate     *MonthYear `db:"end_date" json:"end_date"`
}

type MonthYear struct {
	time.Time
}

// UnmarshalJSON проверяет формат "YYYY-MM" при приёме JSON
func (my *MonthYear) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	t, err := time.Parse("2006-01", s) // Формат YYYY-MM
	if err != nil {
		return err
	}

	my.Time = t
	return nil
}

// MarshalJSON возвращает "YYYY-MM" для JSON
func (my MonthYear) MarshalJSON() ([]byte, error) {
	return json.Marshal(my.Time.Format("2006-01"))
}

func (my *MonthYear) Scan(value interface{}) error {
	if value == nil {
		my.Time = time.Time{}
		return nil
	}

	// Check if the value is a time.Time type
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("could not scan type %T into MonthYear", value)
	}

	my.Time = t
	return nil
}
