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

func Test_GivenAValidClientId_thenReturnCart(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findClientByIdQuery := getMockedQuery("id = ?", aValidClientId)
	firstCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)
	findProductsByCartIdQuery := getMockedQuery("cart_id = ?", aValidCartId)
	findProductByIdQuery := getMockedQuery("id = ?", aValidProductId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(getMockedClient(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Client)
		arg.Name = getMockedClient().Name
		arg.Id = getMockedClient().Id
	})

	clientCart := models.Cart{
		ClientId: aValidClientId,
	}
	dbClientMock.On("First", &clientCart, firstCartByClientIdQuery).Return(getMockedDbObject(getMockedClientCart(), nil)).Run(func(args mock.Arguments) {
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

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	resp, err := repo.GetCart(aValidClientId)

	assert.Nil(t, err)
	assert.Equal(t, getMockedProductList(), *resp)
	dbClientMock.AssertNumberOfCalls(t, "First", 1)
	dbClientMock.AssertNumberOfCalls(t, "Find", 3)
}

func Test_GivenANotValidClientId_thenClientNotFoundInDb(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findClientByIdQuery := getMockedQuery("id = ?", aNotValidClientId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(nil, errors.New("client not found")))

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	resp, err := repo.GetCart(aNotValidClientId)

	assert.EqualError(t, err, fmt.Sprintf("client with id: %v not found", aNotValidClientId))
	assert.Nil(t, resp)
	dbClientMock.AssertNotCalled(t, "First", mock.AnythingOfType("models.Cart"))
	dbClientMock.AssertNotCalled(t, "Find", mock.AnythingOfType("models.ProductCart"))
	dbClientMock.AssertNotCalled(t, "Find", mock.AnythingOfType("models.Product"))
}

func Test_GivenAValidClientId_thenCartNotFoundInDb(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findClientByIdQuery := getMockedQuery("id = ?", aValidClientId)
	firstCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(getMockedClient(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Client)
		arg.Name = getMockedClient().Name
		arg.Id = getMockedClient().Id
	})
	clientCart := models.Cart{
		ClientId: aValidClientId,
	}
	dbClientMock.On("First", &clientCart, firstCartByClientIdQuery).Return(getMockedDbObject(nil, errors.New("cart not found")))

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	resp, err := repo.GetCart(aValidClientId)

	assert.EqualError(t, err, fmt.Sprintf("cart for client id: %v not found", aValidClientId))
	assert.Nil(t, resp)

	dbClientMock.AssertNumberOfCalls(t, "Find", 1)
	dbClientMock.AssertNotCalled(t, "Find", mock.AnythingOfType("models.ProductCart"))
	dbClientMock.AssertNotCalled(t, "Find", mock.AnythingOfType("models.Product"))
}

func Test_GivenAValidClientId_ThenFoundCartButNotFoundListOfProducts(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findClientByIdQuery := getMockedQuery("id = ?", aValidClientId)
	firstCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)
	findProductsByCartIdQuery := getMockedQuery("cart_id = ?", aValidCartId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(getMockedClient(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Client)
		arg.Name = getMockedClient().Name
		arg.Id = getMockedClient().Id
	})

	clientCart := models.Cart{
		ClientId: aValidClientId,
	}
	dbClientMock.On("First", &clientCart, firstCartByClientIdQuery).Return(getMockedDbObject(getMockedClientCart(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Cart)
		arg.Id = aValidCartId
		arg.Client = getMockedClient()
		arg.ClientId = aValidClientId
	})

	var productsCarts []models.ProductCart
	dbClientMock.On("Find", &productsCarts, findProductsByCartIdQuery).Return(getMockedDbObject(nil, errors.New("unable to find list of products")))

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	resp, err := repo.GetCart(aValidClientId)

	assert.EqualError(t, err, fmt.Sprintf("unable to retrieve the list of products"))
	assert.Nil(t, resp)
	dbClientMock.AssertNumberOfCalls(t, "Find", 2)
	dbClientMock.AssertNotCalled(t, "Find", mock.AnythingOfType("models.Product"))
}

func Test_GivenAValidClientId_CartAndProductsFound_ButNotFoundProductsById_ThenReturnEmptyList(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findClientByIdQuery := getMockedQuery("id = ?", aValidClientId)
	firstCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)
	findProductsByCartIdQuery := getMockedQuery("cart_id = ?", aValidCartId)
	findProductByIdQuery := getMockedQuery("id = ?", aValidProductId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(getMockedClient(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Client)
		arg.Name = getMockedClient().Name
		arg.Id = getMockedClient().Id
	})

	clientCart := models.Cart{
		ClientId: aValidClientId,
	}
	dbClientMock.On("First", &clientCart, firstCartByClientIdQuery).Return(getMockedDbObject(getMockedClientCart(), nil)).Run(func(args mock.Arguments) {
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
	dbClientMock.On("Find", &product, findProductByIdQuery).Return(getMockedDbObject(nil, errors.New("unable to find by id")))
	repo := repository.PgCartRepository{DbClient: dbClientMock}

	resp, err := repo.GetCart(aValidClientId)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(*resp))
	dbClientMock.AssertNumberOfCalls(t, "First", 1)
	dbClientMock.AssertNumberOfCalls(t, "Find", 3)
}

func Test_AddProductToACart_Successful(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findClientByIdQuery := getMockedQuery("id = ?", aValidClientId)
	firstOrCreateCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)
	findProductByIdQuery := getMockedQuery("id = ?", aValidProductId)

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

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	err := repo.AddProductToCart(aValidProductId, aValidClientId)

	assert.Nil(t, err)
	dbClientMock.AssertNumberOfCalls(t, "Find", 2)
	dbClientMock.AssertNumberOfCalls(t, "FirstOrCreate", 1)
	dbClientMock.AssertNumberOfCalls(t, "Create", 1)
}

func Test_GivenNotValidClientId_thenClientNotFoundInDbAndNotAbleToAddProduct(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findClientByIdQuery := getMockedQuery("id = ?", aNotValidClientId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(nil, errors.New("client not found")))

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	err := repo.AddProductToCart(aValidProductId, aNotValidClientId)

	assert.EqualError(t, err, fmt.Sprintf("client with id: %v not found", aNotValidClientId))
	dbClientMock.AssertNumberOfCalls(t, "Find", 1)
	dbClientMock.AssertNumberOfCalls(t, "FirstOrCreate", 0)
	dbClientMock.AssertNumberOfCalls(t, "Create", 0)
}

func Test_GivenValidClientId_ThenUnableToFindOrCreateCartToAddProduct(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findClientByIdQuery := getMockedQuery("id = ?", aValidClientId)
	firstOrCreateCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(getMockedClient(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Client)
		arg.Name = getMockedClient().Name
		arg.Id = getMockedClient().Id
	})

	clientCart := models.Cart{
		ClientId: aValidClientId,
	}
	dbClientMock.On("FirstOrCreate", &clientCart, firstOrCreateCartByClientIdQuery).Return(getMockedDbObject(nil, errors.New("cart not found and not created")))

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	err := repo.AddProductToCart(aValidProductId, aValidClientId)

	assert.EqualError(t, err, "cart not found! unable to create new cart")
	dbClientMock.AssertNumberOfCalls(t, "Find", 1)
	dbClientMock.AssertNumberOfCalls(t, "FirstOrCreate", 1)
	dbClientMock.AssertNumberOfCalls(t, "Create", 0)
}

func Test_GivenValidClientIdAndGetValidCart_ThenUnableToFindProductById(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findClientByIdQuery := getMockedQuery("id = ?", aValidClientId)
	firstOrCreateCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)
	findProductByIdQuery := getMockedQuery("id = ?", aValidProductId)

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
	dbClientMock.On("Find", &product, findProductByIdQuery).Return(getMockedDbObject(nil, errors.New("unable to find by id")))

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	err := repo.AddProductToCart(aValidProductId, aValidClientId)

	assert.EqualError(t, err, "unable to find the product in our db")
	dbClientMock.AssertNumberOfCalls(t, "Find", 2)
	dbClientMock.AssertNumberOfCalls(t, "FirstOrCreate", 1)
	dbClientMock.AssertNumberOfCalls(t, "Create", 0)
}

func Test_GivenValidClientId_FoundValidCart_FoundProductById_ThenUnableToCreateProductCartInDb(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findClientByIdQuery := getMockedQuery("id = ?", aValidClientId)
	firstOrCreateCartByClientIdQuery := getMockedQuery("client_id = ?", aValidClientId)
	findProductByIdQuery := getMockedQuery("id = ?", aValidProductId)

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
	dbClientMock.On("Create", &productCart).Return(getMockedDbObject(nil, errors.New("unable to create association")))

	repo := repository.PgCartRepository{DbClient: dbClientMock}

	err := repo.AddProductToCart(aValidProductId, aValidClientId)

	assert.EqualError(t, err, "unable to add product to the cart")
	dbClientMock.AssertNumberOfCalls(t, "Find", 2)
	dbClientMock.AssertNumberOfCalls(t, "FirstOrCreate", 1)
	dbClientMock.AssertNumberOfCalls(t, "Create", 1)
}
