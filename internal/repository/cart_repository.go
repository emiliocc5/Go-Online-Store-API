package repository

import (
	"errors"
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models"
	"sync"
)

var (
	onceCartRepository sync.Once
	cartRepositoryImpl *CartRepositoryImpl
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
	client, err := cr.DbClient.GetClient()
	if err != nil {
		return nil, err
	}

	clientCart := new(models.Cart)
	client.Find(clientCart, "ClientId = ?", clientId)
	if clientCart == nil {
		return nil, errors.New("there is no Cart for the client")
	}

	productsCarts := new([]models.ProductCart)
	client.Find(&productsCarts, "CartId = ?", clientCart.ID)

	productsList := new([]models.Product)

	client.Find(productsList, "Id = ?", 1)

	return productsList, nil
}

func (cr *CartRepositoryImpl) AddProductToCart(productId, clientId int) error {
	client, err := cr.DbClient.GetClient()
	if err != nil {
		return err
	}

	clientCart := models.Cart{
		ClientId: uint(clientId),
	}

	clientResult := client.First(&clientCart, "client_id = ?", clientId)

	if clientResult.Error != nil {
		fmt.Println(fmt.Sprintf("Error: %v", clientResult.Error.Error()))
		createResult := client.Create(&clientCart)
		if createResult.Error != nil {
			return errors.New("cart not found! unable to create new cart")
		}
	}

	//TODO Recover the product

	productCart := models.ProductCart{
		ProductId: uint(productId),
		CartId:    clientCart.ID,
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
