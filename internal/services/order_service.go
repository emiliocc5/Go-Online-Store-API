package services

import (
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/repository"
)

type (
	OrderService interface {
		CreateOrder(clientId, cartId int) error
		GetProductsFromOrder(clientId, orderId int) (*[]models.Product, error)
	}
	OrderServiceImpl struct {
		OrderRepository repository.OrderRepository
	}
)
