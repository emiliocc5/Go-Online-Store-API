package repository

import (
	"errors"
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

const (
	aValidClientId    = 123
	aValidProductId   = 123
	aValidCartId      = 1
	aValidCategoryId  = 1
	aValidLabel       = "aValidLabel"
	aValidType        = 1
	aValidDownloadUrl = ""
	aValidWeight      = 7.5

	aNotValidClientId = 000
)

func Test_GetCart_thenReturnCart(t *testing.T) {
	dbClientMock := &DbClientMock{}

	var findClientByIdQuery []interface{}
	findClientByIdQuery = append(findClientByIdQuery, "id = ?", aValidClientId)

	var firsCartByClientIdQuery []interface{}
	firsCartByClientIdQuery = append(firsCartByClientIdQuery, "client_id = ?", aValidClientId)

	var findProductsByCartIdQuery []interface{}
	findProductsByCartIdQuery = append(findProductsByCartIdQuery, "cart_id = ?", aValidCartId)

	var findProductByIdQuery []interface{}
	findProductByIdQuery = append(findProductByIdQuery, "id = ?", aValidProductId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(getMockedClient(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Client)
		arg.Name = getMockedClient().Name
		arg.Id = getMockedClient().Id
	})

	clientCart := models.Cart{
		ClientId: aValidClientId,
	}

	dbClientMock.On("First", &clientCart, firsCartByClientIdQuery).Return(getMockedDbObject(getMockedClientCart(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Cart)
		arg.Id = aValidCartId
		arg.Client = getMockedClient()
		arg.ClientId = aValidClientId
	})

	var productsCarts []models.ProductCart
	dbClientMock.On("Find", &productsCarts, findProductsByCartIdQuery).Return(getMockedDbObject(getMockedProductCartsList(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*[]models.ProductCart)
		*arg = getMockedProductCartsList()
	})

	product := models.Product{}
	dbClientMock.On("Find", &product, findProductByIdQuery).Return(getMockedDbObject(getMockedProductList(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Product)
		arg.Weight = getMockedProduct().Weight
		arg.DownloadUrl = getMockedProduct().DownloadUrl
		arg.Id = getMockedProduct().Id
		arg.Label = getMockedProduct().Label
		arg.CategoryId = getMockedProduct().CategoryId
		arg.Type = getMockedProduct().Type
	})

	repo := repository.CartRepositoryImpl{DbClient: dbClientMock}

	resp, err := repo.GetCart(aValidClientId)
	assert.Nil(t, err)
	assert.Equal(t, getMockedProductList(), *resp)
}

func Test_GetCart_thenClientNotFoundInDb(t *testing.T) {
	dbClientMock := &DbClientMock{}

	var findClientByIdQuery []interface{}
	findClientByIdQuery = append(findClientByIdQuery, "id = ?", aNotValidClientId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(nil, errors.New("client not found")))

	repo := repository.CartRepositoryImpl{DbClient: dbClientMock}

	resp, err := repo.GetCart(aNotValidClientId)

	assert.Error(t, err, fmt.Sprintf("client with id: %v not found", aNotValidClientId))
	assert.Nil(t, resp)
}

func Test_GetCart_thenCartNotFoundInDb(t *testing.T) {
	dbClientMock := &DbClientMock{}

	var findClientByIdQuery []interface{}
	findClientByIdQuery = append(findClientByIdQuery, "id = ?", aValidClientId)

	var firsCartByClientIdQuery []interface{}
	firsCartByClientIdQuery = append(firsCartByClientIdQuery, "client_id = ?", aValidClientId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(getMockedClient(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Client)
		arg.Name = getMockedClient().Name
		arg.Id = getMockedClient().Id
	})
	clientCart := models.Cart{
		ClientId: aValidClientId,
	}
	dbClientMock.On("First", &clientCart, firsCartByClientIdQuery).Return(getMockedDbObject(nil, errors.New("cart not found")))

	repo := repository.CartRepositoryImpl{DbClient: dbClientMock}

	resp, err := repo.GetCart(aValidClientId)

	assert.Error(t, err, fmt.Sprintf("cart for client id: %v not found", aValidClientId))
	assert.Nil(t, resp)
}

func Test_addProductToACart_Successful(t *testing.T) {
	dbClientMock := &DbClientMock{}

	var findClientByIdQuery []interface{}
	findClientByIdQuery = append(findClientByIdQuery, "id = ?", aValidClientId)

	var firstOrCreateCartByClientIdQuery []interface{}
	firstOrCreateCartByClientIdQuery = append(firstOrCreateCartByClientIdQuery, "client_id = ?", aValidClientId)

	var findProductByIdQuery []interface{}
	findProductByIdQuery = append(findProductByIdQuery, "id = ?", aValidProductId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(getMockedClient(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Client)
		arg.Name = getMockedClient().Name
		arg.Id = getMockedClient().Id
	})

	clientCart := models.Cart{
		ClientId: aValidClientId,
	}
	dbClientMock.On("FirstOrCreate", &clientCart, firstOrCreateCartByClientIdQuery).Return(getMockedDbObject(getMockedClientCart(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Cart)
		arg.Id = getMockedClientCart().Id
		arg.ClientId = getMockedClientCart().ClientId
		arg.Client = getMockedClientCart().Client
	})

	product := models.Product{}
	dbClientMock.On("Find", &product, findProductByIdQuery).Return(getMockedDbObject(getMockedProductList(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Product)
		arg.Weight = getMockedProduct().Weight
		arg.DownloadUrl = getMockedProduct().DownloadUrl
		arg.Id = getMockedProduct().Id
		arg.Label = getMockedProduct().Label
		arg.CategoryId = getMockedProduct().CategoryId
		arg.Type = getMockedProduct().Type
	})

	productCart := models.ProductCart{
		ProductId: aValidProductId,
		CartId:    aValidCartId,
	}
	dbClientMock.On("Create", &productCart).Return(getMockedDbObject(getMockedProductCart(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.ProductCart)
		arg.ProductId = getMockedProductCart().ProductId
		arg.Cart = getMockedProductCart().Cart
		arg.Product = getMockedProductCart().Product
		arg.CartId = getMockedProductCart().CartId
	})

	repo := repository.CartRepositoryImpl{DbClient: dbClientMock}

	err := repo.AddProductToCart(aValidProductId, aValidClientId)
	assert.Nil(t, err)

}

func Test_AddProduct_thenClientNotFoundInDb(t *testing.T) {
	dbClientMock := &DbClientMock{}

	var findClientByIdQuery []interface{}
	findClientByIdQuery = append(findClientByIdQuery, "id = ?", aNotValidClientId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(nil, errors.New("client not found")))

	repo := repository.CartRepositoryImpl{DbClient: dbClientMock}

	err := repo.AddProductToCart(aValidProductId, aNotValidClientId)

	assert.Error(t, err, fmt.Sprintf("client with id: %v not found", aNotValidClientId))
}
