package services

import (
	"errors"
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

const (
	aValidProductId    = 123
	aValidClientId     = 123
	anInvalidProductId = 000
)

func Test_AddProductToCartSuccessful(t *testing.T) {
	clientMockRepository := &ClientRepositoryMock{}
	productMockRepository := &ProductRepositoryMock{}
	cartMockRepository := &CartRepositoryMock{}

	clientMockRepository.On("IsClientInDataBase", aValidClientId).Return(true)

	productMocked := getValidProduct()
	productMockRepository.On("FindProductById", aValidProductId).Return(&productMocked, nil)

	cartMockRepository.On("AddProductToCart", mock.Anything, mock.Anything).Return(nil)

	cs := services.CartServiceImpl{
		CartRepository:    cartMockRepository,
		ClientRepository:  clientMockRepository,
		ProductRepository: productMockRepository,
	}

	err := cs.AddProductToCart(aValidProductId, aValidClientId)

	assert.Nil(t, err)
}

func Test_AddProductToCartWithError(t *testing.T) {
	mockRepository := &CartRepositoryMock{}

	mockRepository.On("AddProductToCart", anInvalidProductId, mock.Anything).Return(errors.New("unable to find the product"))

	cs := services.CartServiceImpl{CartRepository: mockRepository}

	err := cs.AddProductToCart(anInvalidProductId, aValidClientId)

	assert.NotNil(t, err)
	assert.Error(t, err, "unable to find the product")
}

func Test_GetCartWithElements(t *testing.T) {
	mockRepository := &CartRepositoryMock{}

	mockRepository.On("GetCart", mock.Anything, mock.Anything).Return(getValidListOfProducts(), nil)

	cs := services.CartServiceImpl{CartRepository: mockRepository}
	resp, err := cs.GetCart(aValidClientId)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.Products))
	assert.Equal(t, getValidProduct(), resp.Products[0])
}

func Test_GetCartWithError(t *testing.T) {
	mockRepository := &CartRepositoryMock{}
	var emptyProducts []models.Product

	mockRepository.On("GetCart", mock.Anything, mock.Anything).Return(&emptyProducts, errors.New("cart not found"))

	cs := services.CartServiceImpl{CartRepository: mockRepository}
	resp, err := cs.GetCart(aValidClientId)

	assert.Error(t, err, "cart not found")
	assert.Equal(t, 0, len(resp.Products))
}
