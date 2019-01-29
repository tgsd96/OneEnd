package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Purchase struct {
	gorm.Model
	CustID        int64      `json:"cust_id"`
	BillNo        string     `json:"bill_no"`
	Date          *time.Time `gorm:"default current_timestamp" json:"date"`
	Amount        int64      `gorm:"default 0" json:"amount"`
	InterfaceCode string     `json:"i_code"`
}
