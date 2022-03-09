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
	}
)

func init() {
	logger = utils.GetLogger()
}

func (cr *CartRepositoryImpl) GetCart(clientId int) (*[]models.Product, error) {

	clientCart := models.Cart{
		ClientId: clientId,
	}

	clientResult := cr.DbClient.First(&clientCart, "client_id = ?", clientId)
	if clientResult.Error != nil {
		return nil, errors.New("cart not found")
	}

	var productsCarts []models.ProductCart
	productsResult := cr.DbClient.Find(&productsCarts, "cart_id = ?", clientCart.Id)
	if productsResult.Error != nil {
		return nil, errors.New("unable to retrieve the list of products")
	}

	var productsList []models.Product
	for _, e := range productsCarts {
		product := models.Product{}
		productResult := cr.DbClient.Find(&product, "id = ?", e.ProductId)
		if productResult.Error != nil {
			logger.Error("unable to get the product: %s", productResult.Error.Error())
		}
		productsList = append(productsList, product)
	}

	return &productsList, nil
}

func (cr *CartRepositoryImpl) AddProductToCart(productId, clientId int) error {
	clientCart := models.Cart{
		ClientId: clientId,
	}

	clientResult := cr.DbClient.First(&clientCart, "client_id = ?", clientId)

	if clientResult.Error != nil {
		fmt.Println(fmt.Sprintf("Error: %v", clientResult.Error.Error()))
		createResult := cr.DbClient.Create(&clientCart) //TODO Hay un first or create
		if createResult.Error != nil {
			return errors.New("cart not found! unable to create new cart")
		}
	}

	product := models.Product{}
	productResult := cr.DbClient.First(&product, "id = ?", productId)
	if productResult.Error != nil {
		return errors.New("unable to find the product")
	}

	productCart := models.ProductCart{
		ProductId: product.Id,
		CartId:    clientCart.Id,
	}

	productCartResult := cr.DbClient.Create(&productCart)
	if productCartResult.Error != nil {
		return errors.New("unable to add product to the cart")
	}

	fmt.Println("Product added to the cart")

	return nil
}
