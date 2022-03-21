package repository

import (
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

const (
	aValidClientId     = 123
	aValidProductId    = 123
	aValidCartId       = 1
	aValidCategoryId   = 1
	aValidLabel        = "aValidLabel"
	aValidType         = 1
	aValidDownloadUrl  = ""
	aValidWeight       = 7.5
	aNotValidClientId  = 000
	aNotValidProductId = 0
)

type DbClientMock struct{ mock.Mock }

func (mock *DbClientMock) First(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	args := mock.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (mock *DbClientMock) Find(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	args := mock.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (mock *DbClientMock) Create(value interface{}) (tx *gorm.DB) {
	args := mock.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (mock *DbClientMock) FirstOrCreate(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	args := mock.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func getMockedDbObject(dest interface{}, err error) *gorm.DB {
	dbOb := gorm.DB{
		Statement: &gorm.Statement{Dest: dest},
		Error:     err,
	}
	return &dbOb
}

func getMockedClientCart() models.Cart {
	return models.Cart{
		Id:       aValidCartId,
		ClientId: aValidClientId,
		Client:   getMockedClient(),
	}
}

func getMockedClient() models.Client {
	return models.Client{
		Id:   aValidClientId,
		Name: "John Smith",
	}
}
func getMockedProductCart() models.ProductCart {
	return models.ProductCart{ProductId: aValidProductId, Product: getMockedProduct(), Cart: getMockedClientCart(), CartId: aValidCartId}
}

func getMockedProductCartsList() []models.ProductCart {
	var productCarts []models.ProductCart
	productCarts = append(productCarts, getMockedProductCart())
	return productCarts
}

func getMockedProduct() models.Product {
	return models.Product{Id: aValidProductId, CategoryId: aValidCategoryId, Label: aValidLabel, Type: aValidType, DownloadUrl: aValidDownloadUrl, Weight: aValidWeight}
}

func getMockedProductList() []models.Product {
	var products []models.Product
	return append(products, getMockedProduct())
}

func getMockedQuery(query string, parameter int) []interface{} {
	var queryArray []interface{}
	queryArray = append(queryArray, query, parameter)

	return queryArray
}
