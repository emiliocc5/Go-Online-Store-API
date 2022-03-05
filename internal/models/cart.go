package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	Id       int
	ClientId int
}
