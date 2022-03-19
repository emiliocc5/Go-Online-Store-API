package models

type Cart struct {
	Id       int `gorm:"primarykey"`
	ClientId int
	Client   Client
}
