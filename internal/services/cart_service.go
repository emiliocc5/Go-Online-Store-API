package services

import (
	"errors"
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models/response"
	"github.com/emiliocc5/online-store-api/internal/repository"
	"github.com/emiliocc5/online-store-api/internal/utils"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

type CartService interface {
	AddProduct(productId, clientId int) error
	GetCart(clientId int) (response.GetCartResponse, error)
}
type CartServiceImpl struct {
	CartRepository repository.CartRepository
}

func init() {
	logger = utils.GetLogger()
}

func (c *CartServiceImpl) setRepositories() error {
	client, err := repository.GetClient()
	if err != nil {
		logger.Errorf("Error getting client: %v", err)
		return err
	}
	if c.CartRepository == nil {
		c.CartRepository = &repository.CartRepositoryImpl{
			DbClient: client,
		}
	}
	return nil
}

func (c *CartServiceImpl) AddProduct(productId, clientId int) error {
	err := c.setRepositories()
	if err != nil {
		logger.Errorf("Error setting repositories: %v", err)
		return errors.New(fmt.Sprintf("Failed trying to add product: %+v to cart to the client: %+v with error: %+v ",
			productId, clientId, err))
	}

	errAddProd := c.CartRepository.AddProductToCart(productId, clientId)
	if errAddProd != nil {
		fmt.Println(fmt.Sprintf("Failed trying to add product: %+v to cart to the client: %+v with error: %+v",
			productId, clientId, errAddProd))
		return errAddProd
	}
	return nil
}

func (c *CartServiceImpl) GetCart(clientId int) (response.GetCartResponse, error) {
	resp := response.GetCartResponse{}
	err := c.setRepositories()
	if err != nil {
		logger.Errorf("Error setting repositories: %v", err)
		return resp, errors.New(fmt.Sprintf("Failed trying to get cart to the client: %+v with error: %+v ",
			clientId, err))
	}

	prods, errGetCart := c.CartRepository.GetCart(clientId)

	if errGetCart != nil {
		return resp, errGetCart
	}

	resp.Products = *prods

	return resp, nil
}
