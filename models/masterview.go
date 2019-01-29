package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type MasterView struct {
	gorm.Model
	CustID  int64      `json:"cust_id"`
	Name    string     `json:"name"`
	Area    string     `json:"area"`
	Ts      *time.Time `json:"ts"`
	Balance int64      `json:"balance"`
}
