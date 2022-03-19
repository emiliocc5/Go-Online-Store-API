package models

type Client struct {
	Id   int `gorm:"primarykey"`
	Name string
}
