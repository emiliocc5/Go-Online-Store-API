package repository

import (
	"errors"
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/utils"
	"gorm.io/gorm"
)

type (
	OrderRepository interface {
		CreateOrder(clientId, cartId int) error
		GetProductsFromOrder(clientId, orderId int) (*[]models.Product, error)
	}
	PgOrderRepository struct {
		DbClient IOrderRepositoryDbClient
	}
	IOrderRepositoryDbClient interface {
		First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
		Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
		Create(value interface{}) (tx *gorm.DB)
		FirstOrCreate(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	}
)

func init() {
	logger = utils.GetLogger()
}

func (or *PgOrderRepository) CreateOrder(clientId, cartId int) error {
	if !or.isClientInDataBase(clientId) {
		return errors.New(fmt.Sprintf("client with id: %v not found", clientId))
	}

	clientCart := models.Cart{ClientId: clientId}
	err := or.findCartForClientId(&clientCart, clientId)
	if err != nil {
		return err
	}

	var productsCarts []models.ProductCart
	err = or.findListOfProductsForCartId(&productsCarts, clientCart.Id)
	if err != nil {
		return err
	}

	if len(productsCarts) == 0 {
		return errors.New("the cart has no products")
	}

	err = or.createOrderForClientId(clientCart, clientId)
	if err != nil {
		return err
	}

	return nil
}

func (or *PgOrderRepository) GetProductsFromOrder(clientId, orderId int) (*[]models.Product, error) {
	if !or.isClientInDataBase(clientId) {
		return nil, errors.New(fmt.Sprintf("client with id: %v not found", clientId))
	}

	order, err := or.findOrderForClientId(clientId, orderId)
	if err != nil {
		return nil, err
	}

	var productsCarts []models.ProductCart
	err = or.findListOfProductsForCartId(&productsCarts, order.CartId)
	if err != nil {
		return nil, err
	}

	var productList []models.Product
	or.findProductsFromProductCartsList(productsCarts, &productList)

	return &productList, nil
}

func (or *PgOrderRepository) isClientInDataBase(clientId int) bool {
	client := models.Client{}
	findClientResult := or.DbClient.Find(&client, "id = ?", clientId)
	return findClientResult.Error == nil
}

func (or *PgOrderRepository) findCartForClientId(clientCart *models.Cart, clientId int) error {
	clientResult := or.DbClient.First(clientCart, "client_id = ?", clientId)
	if clientResult.Error != nil {
		return errors.New(fmt.Sprintf("cart for client id: %v not found", clientId))
	}
	return nil
}

func (or *PgOrderRepository) createOrderForClientId(cart models.Cart, clientId int) error {
	order := models.Order{
		Cart:     cart,
		CartId:   cart.Id,
		ClientId: clientId,
	}
	createOrderResult := or.DbClient.Create(&order)
	if createOrderResult.Error != nil {
		return errors.New(fmt.Sprintf("Unable to create order for clientId: %v", clientId))
	}
	return nil
}

func (or *PgOrderRepository) findListOfProductsForCartId(productCarts *[]models.ProductCart, clientCartId int) error {
	productsResult := or.DbClient.Find(productCarts, "cart_id = ?", clientCartId)
	if productsResult.Error != nil {
		return errors.New("unable to retrieve the list of products")
	}
	return nil
}

func (or *PgOrderRepository) findProductsFromProductCartsList(productCarts []models.ProductCart, productList *[]models.Product) {
	for _, e := range productCarts {
		product := models.Product{}
		err := or.findProductById(&product, e.ProductId)
		if err != nil {
			logger.Error(fmt.Sprintf("unable to get the product: %v", err.Error()))
		} else {
			*productList = append(*productList, product)
		}
	}
}

func (or *PgOrderRepository) findProductById(product *models.Product, productId int) error {
	productResult := or.DbClient.Find(product, "id = ?", productId)
	if productResult.Error != nil {
		return productResult.Error
	}
	return nil
}

func (or *PgOrderRepository) findOrderForClientId(clientId, orderId int) (models.Order, error) {
	order := models.Order{}
	orderResult := or.DbClient.Find(&order, "client_id = ?", clientId, "order_id = ?", orderId)
	if orderResult.Error != nil {
		return order, errors.New(fmt.Sprintf("order for clientId: %v not found", clientId))
	}
	return order, nil
}