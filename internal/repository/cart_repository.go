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
		GetCartByClient(clientId int) (*models.Cart, error)
		AddProductToCart(productId, clientId int) error
	}
	PgCartRepository struct {
		DbClient ICartRepositoryDbClient
	}
	ICartRepositoryDbClient interface {
		First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
		Create(value interface{}) (tx *gorm.DB)
		FirstOrCreate(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	}
)

func init() {
	logger = utils.GetLogger()
}

func (cr *PgCartRepository) GetCartByClient(clientId int) (*models.Cart, error) {
	clientCart := models.Cart{ClientId: clientId}

	err := cr.findCartForClientId(&clientCart, clientId)
	if err != nil {
		return nil, err
	}

	return &clientCart, nil
}

func (cr *PgCartRepository) AddProductToCart(productId, clientId int) error {
	clientCart := models.Cart{ClientId: clientId}
	err := cr.findOrCreateCartForClientId(&clientCart, clientId)
	if err != nil {
		return err
	}

	productCart := models.ProductCart{
		ProductId: productId,
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

func (cr *PgCartRepository) findCartForClientId(clientCart *models.Cart, clientId int) error {
	clientResult := cr.DbClient.First(clientCart, "client_id = ?", clientId)
	if clientResult.Error != nil {
		return errors.New(fmt.Sprintf("cart for client id: %v not found", clientId))
	}
	return nil
}

func (cr *PgCartRepository) findOrCreateCartForClientId(clientCart *models.Cart, clientId int) error {
	clientResult := cr.DbClient.FirstOrCreate(clientCart, "client_id = ?", clientId)
	if clientResult.Error != nil {
		logger.Errorf("Error: %v", clientResult.Error.Error())
		return errors.New("cart not found! unable to create new cart")
	}
	return nil
}
