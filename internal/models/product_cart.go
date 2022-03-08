package models

type ProductCart struct {
	Id        int `gorm:"primarykey"`
	ProductId int
	Product   Product
	CartId    int
	Cart      Cart
}
