package services

import (
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type CartRepositoryMock struct{ mock.Mock }

func (mock *CartRepositoryMock) GetCart(clientId int) (*[]models.Product, error) {
	args := mock.Called(clientId)
	return args.Get(0).(*[]models.Product), args.Error(1)
}

func (mock *CartRepositoryMock) AddProductToCart(productId, clientId int) error {
	args := mock.Called(productId, clientId)
	return args.Error(0)
}
