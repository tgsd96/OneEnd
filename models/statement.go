package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Statement : For bank entries
type Statement struct {
	gorm.Model
	Date         *time.Time `json:"date"`
	Narration    string     `json:"narration"`
	RefNo        string     `json:"ref_no"`
	ValueDt      string     `json:"value_date"`
	Withdrawl    float32    `json:"withdrawl"`
	Deposit      float32    `json:"deposit"`
	ClosingValue float64    `json:"closing_value"`
	UserID       int64      `json:"user_id"`
}
