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

	product := models.Product{}
	productMockRepository.On("FindProductById", anInvalidProductId).Return(&product, errors.New("product not found"))

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

	cartMockRepository.On("AddProductToCart", aValidProductId, aValidClientId).Return(errors.New("could not add product"))

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
	cartMockRepository.AssertCalled(t, "AddProductToCart", aValidProductId, aValidClientId)
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

func Test_GivenANotValidClientId_ThenUnableToGetCart(t *testing.T) {
	clientMockRepository := &ClientRepositoryMock{}
	productMockRepository := &ProductRepositoryMock{}
	cartMockRepository := &CartRepositoryMock{}

	clientMockRepository.On("IsClientInDataBase", anInvalidClientId).Return(false)

	validCart := getMockedValidCart()
	cartMockRepository.On("GetCartByClient", anInvalidClientId).Return(&validCart, nil)

	productMockRepository.On("FindProductsFromCart", aValidCartId).Return(getValidListOfProducts(), nil)

	cs := services.CartServiceImpl{
		CartRepository:    cartMockRepository,
		ClientRepository:  clientMockRepository,
		ProductRepository: productMockRepository,
	}
	resp, err := cs.GetCart(anInvalidClientId)

	assert.EqualError(t, err, fmt.Sprintf("client with id: %v not found", anInvalidClientId))
	assert.NotNil(t, resp)
	assert.Equal(t, 0, len(resp.Products))
	clientMockRepository.AssertNumberOfCalls(t, "IsClientInDataBase", 1)
	cartMockRepository.AssertNumberOfCalls(t, "GetCartByClient", 0)
	productMockRepository.AssertNumberOfCalls(t, "FindProductsFromCart", 0)
}

func Test_GivenAValidClientId_ThenCartNotFound(t *testing.T) {
	clientMockRepository := &ClientRepositoryMock{}
	productMockRepository := &ProductRepositoryMock{}
	cartMockRepository := &CartRepositoryMock{}

	clientMockRepository.On("IsClientInDataBase", aValidClientId).Return(true)

	cart := models.Cart{}
	cartMockRepository.On("GetCartByClient", aValidClientId).Return(&cart, errors.New("cart not found"))

	productMockRepository.On("FindProductsFromCart", aValidCartId).Return(getValidListOfProducts(), nil)

	cs := services.CartServiceImpl{
		CartRepository:    cartMockRepository,
		ClientRepository:  clientMockRepository,
		ProductRepository: productMockRepository,
	}
	resp, err := cs.GetCart(aValidClientId)

	assert.EqualError(t, err, fmt.Sprintf("Failed trying to get cart for the client: %+v", aValidClientId))
	assert.NotNil(t, resp)
	assert.Equal(t, 0, len(resp.Products))
	clientMockRepository.AssertNumberOfCalls(t, "IsClientInDataBase", 1)
	cartMockRepository.AssertNumberOfCalls(t, "GetCartByClient", 1)
	productMockRepository.AssertNumberOfCalls(t, "FindProductsFromCart", 0)
}

func Test_GivenAValidClientId_CartFound_ThenListOfProductsNotFound(t *testing.T) {
	clientMockRepository := &ClientRepositoryMock{}
	productMockRepository := &ProductRepositoryMock{}
	cartMockRepository := &CartRepositoryMock{}

	clientMockRepository.On("IsClientInDataBase", aValidClientId).Return(true)

	cart := getMockedValidCart()
	cartMockRepository.On("GetCartByClient", aValidClientId).Return(&cart, nil)

	var emptyProducts []models.Product
	productMockRepository.On("FindProductsFromCart", aValidCartId).Return(&emptyProducts, errors.New("products not found"))

	cs := services.CartServiceImpl{
		CartRepository:    cartMockRepository,
		ClientRepository:  clientMockRepository,
		ProductRepository: productMockRepository,
	}
	resp, err := cs.GetCart(aValidClientId)

	assert.EqualError(t, err, fmt.Sprintf(fmt.Sprintf("unable to get the list of products from the cart: %v", cart.Id)))
	assert.NotNil(t, resp)
	assert.Equal(t, 0, len(resp.Products))
	clientMockRepository.AssertNumberOfCalls(t, "IsClientInDataBase", 1)
	cartMockRepository.AssertNumberOfCalls(t, "GetCartByClient", 1)
	productMockRepository.AssertNumberOfCalls(t, "FindProductsFromCart", 1)
}
