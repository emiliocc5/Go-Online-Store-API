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

func Test_GivenAValidCartId_ThenReturnAListOfProducts(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findProductsByCartIdQuery := getMockedQuery("cart_id = ?", aValidCartId)
	findProductByIdQuery := getMockedQuery("id = ?", aValidProductId)

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

	repo := repository.PgProductRepository{DbClient: dbClientMock}

	resp, err := repo.FindProductsFromCart(aValidCartId)

	assert.Nil(t, err)
	assert.Equal(t, getMockedProductList(), *resp)
	dbClientMock.AssertNumberOfCalls(t, "Find", 2)
}

func Test_GivenAValidCartId_ThenNotFoundListOfProducts(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findProductsByCartIdQuery := getMockedQuery("cart_id = ?", aValidCartId)

	var productsCarts []models.ProductCart
	dbClientMock.On("Find", &productsCarts, findProductsByCartIdQuery).Return(getMockedDbObject(nil, errors.New("unable to find list of products")))

	repo := repository.PgProductRepository{DbClient: dbClientMock}

	resp, err := repo.FindProductsFromCart(aValidCartId)

	assert.EqualError(t, err, fmt.Sprintf("unable to retrieve the list of products"))
	assert.Nil(t, resp)
	dbClientMock.AssertNumberOfCalls(t, "Find", 1)
	dbClientMock.AssertNotCalled(t, "Find", mock.AnythingOfType("models.Product"))
}

func Test_GivenAValidCartId_AndProductsFound_NotFoundProductsById_ThenReturnEmptyList(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findProductsByCartIdQuery := getMockedQuery("cart_id = ?", aValidCartId)
	findProductByIdQuery := getMockedQuery("id = ?", aValidProductId)

	var productsCarts []models.ProductCart
	dbClientMock.On("Find", &productsCarts, findProductsByCartIdQuery).Return(getMockedDbObject(getMockedProductCartsList(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*[]models.ProductCart)
		*arg = getMockedProductCartsList()
	})

	product := models.Product{}
	dbClientMock.On("Find", &product, findProductByIdQuery).Return(getMockedDbObject(nil, errors.New("unable to find by id")))
	repo := repository.PgProductRepository{DbClient: dbClientMock}

	resp, err := repo.FindProductsFromCart(aValidCartId)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(*resp))
	dbClientMock.AssertNumberOfCalls(t, "Find", 2)
}

func Test_GivenAValidProductId_ThenReturnAValidProduct(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findProductByIdQuery := getMockedQuery("id = ?", aValidProductId)

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

	repo := repository.PgProductRepository{DbClient: dbClientMock}

	prod, err := repo.FindProductById(aValidProductId)

	assert.Nil(t, err)
	assert.Equal(t, getMockedProduct(), *prod)
	dbClientMock.AssertNumberOfCalls(t, "Find", 1)
}

func Test_GivenANotValidProductId_ThenUnableToFindProductById(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findProductByIdQuery := getMockedQuery("id = ?", aNotValidProductId)

	product := models.Product{}
	dbClientMock.On("Find", &product, findProductByIdQuery).Return(getMockedDbObject(nil, errors.New("unable to find by id")))

	repo := repository.PgProductRepository{DbClient: dbClientMock}

	prod, err := repo.FindProductById(aNotValidProductId)

	assert.EqualError(t, err, "unable to find the product in our db")
	assert.Nil(t, prod)
	dbClientMock.AssertNumberOfCalls(t, "Find", 1)
}
