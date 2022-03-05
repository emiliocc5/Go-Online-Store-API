package services

import "fmt"

type IOrderService interface {
	CreateOrder()
}
type OrderService struct {
}

func (o *OrderService) CreateOrder() {
	fmt.Println("Creating order")
}
