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

func Test_GivenAValidClientId_thenReturnCart(t *testing.T) {
	dbClientMock := &DbClientMock{}

	firstCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)

	clientCart := models.Cart{
		ClientId: aValidClientId,
	}
	dbClientMock.On("First", &clientCart, firstCartByClientIdQuery).Return(getMockedDbObject(getMockedClientCart(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Cart)
		arg.Id = aValidCartId
		arg.Client = getMockedClient()
		arg.ClientId = aValidClientId
	})

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	resp, err := repo.GetCartByClient(aValidClientId)

	assert.Nil(t, err)
	assert.Equal(t, getMockedClientCart(), *resp)
	dbClientMock.AssertNumberOfCalls(t, "First", 1)
}

func Test_GivenAValidClientId_thenCartNotFoundInDb(t *testing.T) {
	dbClientMock := &DbClientMock{}

	firstCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)

	clientCart := models.Cart{
		ClientId: aValidClientId,
	}
	dbClientMock.On("First", &clientCart, firstCartByClientIdQuery).Return(getMockedDbObject(nil, errors.New("cart not found")))

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	resp, err := repo.GetCartByClient(aValidClientId)

	assert.Nil(t, resp)
	assert.EqualError(t, err, fmt.Sprintf("cart for client id: %v not found", aValidClientId))
	dbClientMock.AssertNumberOfCalls(t, "First", 1)
}

func Test_AddProductToACart_Successful(t *testing.T) {
	dbClientMock := &DbClientMock{}

	firstOrCreateCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)

	clientCart := models.Cart{
		ClientId: aValidClientId,
	}
	dbClientMock.On("FirstOrCreate", &clientCart, firstOrCreateCartByClientIdQuery).Return(getMockedDbObject(getMockedClientCart(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Cart)
		arg.Id = getMockedClientCart().Id
		arg.ClientId = getMockedClientCart().ClientId
		arg.Client = getMockedClientCart().Client
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

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	err := repo.AddProductToCart(aValidProductId, aValidClientId)

	assert.Nil(t, err)
	dbClientMock.AssertNumberOfCalls(t, "FirstOrCreate", 1)
	dbClientMock.AssertNumberOfCalls(t, "Create", 1)
}

func Test_GivenValidClientId_ThenUnableToFindOrCreateCartToAddProduct(t *testing.T) {
	dbClientMock := &DbClientMock{}

	firstOrCreateCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)

	clientCart := models.Cart{
		ClientId: aValidClientId,
	}
	dbClientMock.On("FirstOrCreate", &clientCart, firstOrCreateCartByClientIdQuery).Return(getMockedDbObject(nil, errors.New("cart not found and not created")))

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	err := repo.AddProductToCart(aValidProductId, aValidClientId)

	assert.EqualError(t, err, "cart not found! unable to create new cart")
	dbClientMock.AssertNumberOfCalls(t, "FirstOrCreate", 1)
	dbClientMock.AssertNumberOfCalls(t, "Create", 0)
}

func Test_GivenValidClientId_ValidProductById_ThenUnableToCreateProductCartInDb(t *testing.T) {
	dbClientMock := &DbClientMock{}

	firstOrCreateCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)

	clientCart := models.Cart{
		ClientId: aValidClientId,
	}
	dbClientMock.On("FirstOrCreate", &clientCart, firstOrCreateCartByClientIdQuery).Return(getMockedDbObject(getMockedClientCart(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Cart)
		arg.Id = getMockedClientCart().Id
		arg.ClientId = getMockedClientCart().ClientId
		arg.Client = getMockedClientCart().Client
	})

	productCart := models.ProductCart{
		ProductId: aValidProductId,
		CartId:    aValidCartId,
	}
	dbClientMock.On("Create", &productCart).Return(getMockedDbObject(nil, errors.New("unable to create association")))

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	err := repo.AddProductToCart(aValidProductId, aValidClientId)

	assert.EqualError(t, err, "unable to add product to the cart")
	dbClientMock.AssertNumberOfCalls(t, "FirstOrCreate", 1)
	dbClientMock.AssertNumberOfCalls(t, "Create", 1)
}
