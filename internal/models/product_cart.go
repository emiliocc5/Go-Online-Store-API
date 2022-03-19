package models

type ProductCart struct {
	Id        int `gorm:"primarykey"`
	ProductId int
	Product   Product //TODO Check if this its neccessary
	CartId    int
	Cart      Cart //TODO Check if this its neccessary
}
