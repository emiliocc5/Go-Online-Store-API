package repository

import (
	"errors"
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/utils"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	onceCartRepository sync.Once
	cartRepositoryImpl *CartRepositoryImpl
	logger             *logrus.Logger
)

type (
	CartRepository interface {
		GetCart(clientId int) (*[]models.Product, error)
		AddProductToCart(productId, clientId int) error
	}
	CartRepositoryImpl struct {
		DbClient IDBClient
	}
)

func (cr *CartRepositoryImpl) GetCart(clientId int) (*[]models.Product, error) {
	logger = utils.GetLogger()
	client, err := cr.DbClient.GetClient()
	if err != nil {
		return nil, err
	}

	clientCart := models.Cart{
		ClientId: clientId,
	}

	clientResult := client.First(&clientCart, "client_id = ?", clientId)
	if clientResult.Error != nil {
		return nil, errors.New("cart not found")
	}

	var productsCarts []models.ProductCart
	productsResult := client.Find(&productsCarts, "cart_id = ?", clientCart.Id)
	if productsResult.Error != nil {
		return nil, errors.New("unable to retrieve the list of products")
	}

	var productsList []models.Product
	for _, e := range productsCarts {
		product := models.Product{}
		productResult := client.Find(&product, "id = ?", e.ProductId)
		if productResult.Error != nil {
			logger.Error("unable to get the product: %s", productResult.Error.Error())
		}
		productsList = append(productsList, product)
	}

	return &productsList, nil
}

func (cr *CartRepositoryImpl) AddProductToCart(productId, clientId int) error {
	client, err := cr.DbClient.GetClient()
	if err != nil {
		return err
	}

	clientCart := models.Cart{
		ClientId: clientId,
	}

	clientResult := client.First(&clientCart, "client_id = ?", clientId)

	if clientResult.Error != nil {
		fmt.Println(fmt.Sprintf("Error: %v", clientResult.Error.Error()))
		createResult := client.Create(&clientCart)
		if createResult.Error != nil {
			return errors.New("cart not found! unable to create new cart")
		}
	}

	product := models.Product{}
	productResult := client.First(&product, "id = ?", productId)
	if productResult.Error != nil {
		return errors.New("unable to find the product")
	}

	productCart := models.ProductCart{
		ProductId: product.Id,
		CartId:    clientCart.Id,
	}

	productCartResult := client.Create(&productCart)
	if productCartResult.Error != nil {
		return errors.New("unable to add product to the cart")
	}

	fmt.Println("Product added to the cart")

	return nil
}

func GetCartRepository() *CartRepositoryImpl {
	onceCartRepository.Do(func() {
		cartRepositoryImpl = &CartRepositoryImpl{
			DbClient: GetDbClient(),
		}
	})
	return cartRepositoryImpl
}
