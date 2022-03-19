package handler

import (
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/models/response"
	"github.com/stretchr/testify/mock"
)

type CartServiceMock struct{ mock.Mock }

func (mock *CartServiceMock) AddProduct(productId, clientId int) error {
	args := mock.Called(productId, clientId)
	return args.Error(0)
}

func (mock *CartServiceMock) GetCart(clientId int) (response.GetCartResponse, error) {
	args := mock.Called(clientId)
	return args.Get(0).(response.GetCartResponse), args.Error(1)
}

func getMockedValidCartResponse() response.GetCartResponse {
	return response.GetCartResponse{Products: *getValidListOfProducts()}
}

func getValidListOfProducts() *[]models.Product {
	var products []models.Product

	product := getValidProduct()

	products = append(products, product)
	return &products
}

func getValidProduct() models.Product {
	return models.Product{
		Id:          1,
		CategoryId:  1,
		Label:       "Keyboard",
		Type:        1,
		DownloadUrl: "",
		Weight:      3.5,
	}
}
