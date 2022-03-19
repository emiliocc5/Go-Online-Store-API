package models

type Category struct {
	Id   int `gorm:"primarykey"`
	Name string
}
