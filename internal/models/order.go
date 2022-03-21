package models

type Order struct {
	Id       int `gorm:"primarykey"`
	ClientId int
	Client   Client
	CartId   int
	Cart     Cart //TODO Check if this its neccessary
}
