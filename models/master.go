package models

import (
	"github.com/jinzhu/gorm"
)

type Master struct {
	gorm.Model
	CustID        int64  `json:"cust_id"`
	UID           int64  `json:"uid,omitempty"`
	Name          string `json:"name"`
	GSTIN         string `json:"gstin"`
	Area          string `json:"area,omitempty"`
	Address       string `json:"address,omitempty"`
	InterfaceCode string `json:"i_code,omitempty"`
	SearchText    string `gorm:"type:tsvector"`
	Group         string `json:"group"`
}
