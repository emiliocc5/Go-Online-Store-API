package repository

import (
	"errors"
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	logger *logrus.Logger
)

type (
	CartRepository interface {
		GetCart(clientId int) (*[]models.Product, error)
		AddProductToCart(productId, clientId int) error
	}
	CartRepositoryImpl struct {
		DbClient ICartRepositoryDbClient
	}
	ICartRepositoryDbClient interface {
		First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
		Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
		Create(value interface{}) (tx *gorm.DB)
		FirstOrCreate(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	}
)

func init() {
	logger = utils.GetLogger()
}

func (cr *CartRepositoryImpl) GetCart(clientId int) (*[]models.Product, error) {
	if !cr.isClientInDataBase(clientId) {
		return nil, errors.New(fmt.Sprintf("client with id: %v not found", clientId))
	}
	clientCart := models.Cart{ClientId: clientId}

	err := cr.findCartForClientId(&clientCart, clientId)
	if err != nil {
		return nil, err
	}

	var productsCarts []models.ProductCart
	err = cr.findListOfProductsForCartId(&productsCarts, clientCart.Id)
	if err != nil {
		return nil, err
	}

	var productList []models.Product
	cr.findProductsFromProductCartsList(productsCarts, &productList)

	return &productList, nil
}

func (cr *CartRepositoryImpl) AddProductToCart(productId, clientId int) error {
	if !cr.isClientInDataBase(clientId) {
		return errors.New(fmt.Sprintf("client with id: %v not found", clientId))
	}

	clientCart := models.Cart{ClientId: clientId}
	err := cr.findOrCreateCartForClientId(&clientCart, clientId)
	if err != nil {
		return err
	}

	product := models.Product{}
	err = cr.findProductById(&product, productId)
	if err != nil {
		logger.Error(fmt.Sprintf("unable to get the product: %v", err.Error()))
		return errors.New("unable to find the product in our db")
	}

	productCart := models.ProductCart{
		ProductId: product.Id,
		CartId:    clientCart.Id,
	}

	productCartResult := cr.DbClient.Create(&productCart)
	if productCartResult.Error != nil {
		logger.Errorf("Error: %v", productCartResult.Error.Error())
		return errors.New("unable to add product to the cart")
	}

	logger.Info("Product added to the cart")

	return nil
}

func (cr *CartRepositoryImpl) isClientInDataBase(clientId int) bool {
	client := models.Client{}
	findClientResult := cr.DbClient.Find(&client, "id = ?", clientId)
	return findClientResult.Error == nil
}

func (cr *CartRepositoryImpl) findCartForClientId(clientCart *models.Cart, clientId int) error {
	clientResult := cr.DbClient.First(clientCart, "client_id = ?", clientId)
	if clientResult.Error != nil {
		return errors.New(fmt.Sprintf("cart for client id: %v not found", clientId))
	}
	return nil
}

func (cr *CartRepositoryImpl) findListOfProductsForCartId(productCarts *[]models.ProductCart, clientCartId int) error {
	productsResult := cr.DbClient.Find(productCarts, "cart_id = ?", clientCartId)
	if productsResult.Error != nil {
		return errors.New("unable to retrieve the list of products")
	}
	return nil
}

func (cr *CartRepositoryImpl) findProductsFromProductCartsList(productCarts []models.ProductCart, productList *[]models.Product) {
	for _, e := range productCarts {
		product := models.Product{}
		err := cr.findProductById(&product, e.ProductId)
		if err != nil {
			logger.Error(fmt.Sprintf("unable to get the product: %v", err.Error()))
		} else {
			*productList = append(*productList, product)
		}
	}
}

func (cr *CartRepositoryImpl) findProductById(product *models.Product, productId int) error {
	productResult := cr.DbClient.Find(product, "id = ?", productId)
	if productResult.Error != nil {
		return productResult.Error
	}
	return nil
}

func (cr *CartRepositoryImpl) findOrCreateCartForClientId(clientCart *models.Cart, clientId int) error {
	clientResult := cr.DbClient.FirstOrCreate(clientCart, "client_id = ?", clientId)
	if clientResult.Error != nil {
		logger.Errorf("Error: %v", clientResult.Error.Error())
		return errors.New("cart not found! unable to create new cart")
	}
	return nil
}
