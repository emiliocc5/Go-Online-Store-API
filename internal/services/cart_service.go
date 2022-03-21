package services

import (
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/models/response"
	"github.com/emiliocc5/online-store-api/internal/repository"
	"github.com/emiliocc5/online-store-api/internal/utils"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

type (
	CartService interface {
		AddProduct(productId, clientId int) error
		GetCart(clientId int) (response.GetCartResponse, error)
	}
	CartServiceImpl struct {
		CartRepository repository.CartRepository
	}
)

func init() {
	logger = utils.GetLogger()
}

func (c *CartServiceImpl) AddProduct(productId, clientId int) error {
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
	prods, errGetCart := c.CartRepository.GetCart(clientId)

	if errGetCart != nil {
		return resp, errGetCart
	}

	resp.Products = *prods

	return resp, nil
}

func parseProducts([]models.Product) []response.ProductResponse {
	
}
