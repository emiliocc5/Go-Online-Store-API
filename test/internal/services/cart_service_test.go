package services

import (
	"errors"
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_GivenAValidClientId_AValidProductIt_ThenAddProductToCart(t *testing.T) {
	clientMockRepository := &ClientRepositoryMock{}
	productMockRepository := &ProductRepositoryMock{}
	cartMockRepository := &CartRepositoryMock{}

	clientMockRepository.On("IsClientInDataBase", aValidClientId).Return(true)

	productMocked := getValidDbProduct()
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

func Test_GivenAnInvalidClientId_ThenClientNotFoundInDb(t *testing.T) {
	clientMockRepository := &ClientRepositoryMock{}
	productMockRepository := &ProductRepositoryMock{}
	cartMockRepository := &CartRepositoryMock{}

	clientMockRepository.On("IsClientInDataBase", anInvalidClientId).Return(false)

	productMocked := getValidDbProduct()
	productMockRepository.On("FindProductById", aValidProductId).Return(&productMocked, nil)

	cartMockRepository.On("AddProductToCart", mock.Anything, mock.Anything).Return(nil)

	cs := services.CartServiceImpl{
		CartRepository:    cartMockRepository,
		ClientRepository:  clientMockRepository,
		ProductRepository: productMockRepository,
	}

	err := cs.AddProductToCart(aValidProductId, anInvalidClientId)

	assert.EqualError(t, err, fmt.Sprintf("client with id: %v not found", anInvalidClientId))
	productMockRepository.AssertNotCalled(t, "FindProductById")
	cartMockRepository.AssertNotCalled(t, "AddProductToCart")
}

func Test_GivenAValidClientId_AnInvalidProductId_ThenProductNotFoundInDb(t *testing.T) {
	clientMockRepository := &ClientRepositoryMock{}
	productMockRepository := &ProductRepositoryMock{}
	cartMockRepository := &CartRepositoryMock{}

	clientMockRepository.On("IsClientInDataBase", aValidClientId).Return(true)

	productMockRepository.On("FindProductById", anInvalidProductId).Return(nil, errors.New("product not found"))

	cartMockRepository.On("AddProductToCart", mock.Anything, mock.Anything).Return(nil)

	cs := services.CartServiceImpl{
		CartRepository:    cartMockRepository,
		ClientRepository:  clientMockRepository,
		ProductRepository: productMockRepository,
	}

	err := cs.AddProductToCart(anInvalidProductId, aValidClientId)

	assert.EqualError(t, err, "unable to find the product in our db")

	clientMockRepository.AssertNumberOfCalls(t, "IsClientInDataBase", 1)
	productMockRepository.AssertNumberOfCalls(t, "FindProductById", 1)
	productMockRepository.AssertCalled(t, "FindProductById", anInvalidProductId)
	cartMockRepository.AssertNotCalled(t, "AddProductToCart")
}

func Test_GivenAValidClientId_AValidProductId_ThenUnableToAddProductToTheCart(t *testing.T) {
	clientMockRepository := &ClientRepositoryMock{}
	productMockRepository := &ProductRepositoryMock{}
	cartMockRepository := &CartRepositoryMock{}

	clientMockRepository.On("IsClientInDataBase", aValidClientId).Return(true)

	productMocked := getValidDbProduct()
	productMockRepository.On("FindProductById", aValidProductId).Return(&productMocked, nil)

	cartMockRepository.On("AddProductToCart", productMocked.Id, aValidClientId).Return(errors.New("could not add product"))

	cs := services.CartServiceImpl{
		CartRepository:    cartMockRepository,
		ClientRepository:  clientMockRepository,
		ProductRepository: productMockRepository,
	}

	err := cs.AddProductToCart(aValidProductId, aValidClientId)

	assert.EqualError(t, err, "unable to add product to the cart")

	clientMockRepository.AssertNumberOfCalls(t, "IsClientInDataBase", 1)
	productMockRepository.AssertNumberOfCalls(t, "FindProductById", 1)
	productMockRepository.AssertCalled(t, "FindProductById", aValidProductId)
	cartMockRepository.AssertNumberOfCalls(t, "AddProductToCart", 1)
	cartMockRepository.AssertCalled(t, "AddProductToCart", aValidClientId, aValidProductId)
}

func Test_GivenAValidClientId_ThenReturnCartWithElements(t *testing.T) {
	clientMockRepository := &ClientRepositoryMock{}
	productMockRepository := &ProductRepositoryMock{}
	cartMockRepository := &CartRepositoryMock{}

	clientMockRepository.On("IsClientInDataBase", aValidClientId).Return(true)

	validCart := getMockedValidCart()
	cartMockRepository.On("GetCartByClient", aValidClientId).Return(&validCart, nil)

	productMockRepository.On("FindProductsFromCart", aValidCartId).Return(getValidListOfProducts(), nil)

	cs := services.CartServiceImpl{
		CartRepository:    cartMockRepository,
		ClientRepository:  clientMockRepository,
		ProductRepository: productMockRepository,
	}
	resp, err := cs.GetCart(aValidClientId)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.Products))
	assert.Equal(t, getValidResponseProduct(), resp.Products[0])
	clientMockRepository.AssertNumberOfCalls(t, "IsClientInDataBase", 1)
	cartMockRepository.AssertNumberOfCalls(t, "GetCartByClient", 1)
	productMockRepository.AssertNumberOfCalls(t, "FindProductsFromCart", 1)
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
