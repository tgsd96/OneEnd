package models

import (
	"github.com/jinzhu/gorm"
)

type ErrorPurchases struct {
	gorm.Model
	CustID        int64  `json:"cust_id"`
	BillNo        string `json:"bill_no,omitempty"`
	Name          string `json:"name"`
	Amount        int64  `json:"amount,omitempty"`
	InterfaceCode string `json:"i_code,omitempty"`
}
