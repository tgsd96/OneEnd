package models

import (
	"tgsd96/onend/utils"

	"github.com/jinzhu/gorm"
)

type Users struct {
	gorm.Model
	Username string `gorm:"unique;not_null"`
	PwdHash  string `json:"-"`
}

func CreateUser(username string, password string) *Users {
	return &Users{
		Username: username,
		PwdHash:  utils.HashAndSalt([]byte(password)),
	}
}
