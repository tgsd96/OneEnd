package models

import (
	"github.com/jinzhu/gorm"
)

type Route struct {
	gorm.Model
	RouteName string `json:"route_name"`
	RouteDay  string `json:"route_day"`
}
