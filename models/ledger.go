package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Ledger struct {
	gorm.Model
	CustID   int64      `json:"cust_id"`
	Date     *time.Time `gorm:"default current_timestamp" json:"date"`
	Amount   int64      `json:"amount"`
	Type     string     `json:"type"`
	Comments string     `json:"comments"`
	UserID   int64      `json:"user_id"`
}
