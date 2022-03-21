package services

import (
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/repository"
	"github.com/emiliocc5/online-store-api/internal/utils"
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

func init() {
	logger = utils.GetLogger()
}

func (os *OrderServiceImpl) CreateOrder(clientId, cartId int) error {
	createOrderErr := os.OrderRepository.CreateOrder(clientId, cartId)

	if createOrderErr != nil {
		return createOrderErr
	}

	return nil
}
