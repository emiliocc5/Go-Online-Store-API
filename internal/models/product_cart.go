package models

import "gorm.io/gorm"

type ProductCart struct {
	gorm.Model
	ProductId uint
	Product   Product
	CartId    uint
	Cart      Cart
}
