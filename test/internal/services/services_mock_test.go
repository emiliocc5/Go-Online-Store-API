package services

import (
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/models/response"
	"github.com/stretchr/testify/mock"
)

const (
	aValidProductId    = 1
	aValidClientId     = 2
	aValidCartId       = 3
	anInvalidProductId = 000
	anInvalidClientId  = 111
)

type CartRepositoryMock struct{ mock.Mock }
type ClientRepositoryMock struct{ mock.Mock }
type ProductRepositoryMock struct{ mock.Mock }

func (mock *CartRepositoryMock) GetCartByClient(clientId int) (*models.Cart, error) {
	args := mock.Called(clientId)
	return args.Get(0).(*models.Cart), args.Error(1)
}

func (mock *CartRepositoryMock) AddProductToCart(productId, clientId int) error {
	args := mock.Called(productId, clientId)
	return args.Error(0)
}

func (mock *ClientRepositoryMock) IsClientInDataBase(clientId int) bool {
	args := mock.Called(clientId)
	return args.Get(0).(bool)
}

func (mock *ProductRepositoryMock) FindProductsFromCart(cartId int) (*[]models.Product, error) {
	args := mock.Called(cartId)
	return args.Get(0).(*[]models.Product), args.Error(1)
}

func (mock *ProductRepositoryMock) FindProductById(productId int) (*models.Product, error) {
	args := mock.Called(productId)
	return args.Get(0).(*models.Product), args.Error(1)
}

func getValidListOfProducts() *[]models.Product {
	var products []models.Product

	product := getValidDbProduct()

	products = append(products, product)
	return &products
}

func getValidDbProduct() models.Product {
	return models.Product{
		Id:          aValidProductId,
		CategoryId:  1,
		Label:       "Keyboard",
		Type:        1,
		DownloadUrl: "",
		Weight:      3.5,
	}
}

func getValidResponseProduct() response.ProductResponse {
	return response.ProductResponse{
		Id:          aValidProductId,
		Label:       "Keyboard",
		Type:        1,
		DownloadUrl: "",
		Weight:      3.5,
	}
}

func getMockedValidCart() models.Cart {
	return models.Cart{
		Id:       aValidCartId,
		ClientId: aValidClientId,
	}
}
