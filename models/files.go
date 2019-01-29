package models

import (
	"github.com/jinzhu/gorm"
)

type File struct {
	gorm.Model
	Location string
	Original string
	Filename string
	Company  string
}
