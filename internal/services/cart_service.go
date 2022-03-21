package services

import (
	"errors"
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
		CartRepository   repository.CartRepository
		ClientRepository repository.ClientRepository
	}
)

func init() {
	logger = utils.GetLogger()
}

//TODO crear una estructura de product repository para reducir la lógica en las clases de repositorio

func (c *CartServiceImpl) AddProduct(productId, clientId int) error {
	if !c.ClientRepository.IsClientInDataBase(clientId) {
		return errors.New(fmt.Sprintf("client with id: %v not found", clientId))
	}

	errAddProd := c.CartRepository.AddProductToCart(productId, clientId)
	if errAddProd != nil {
		logger.Error(fmt.Sprintf("Failed trying to add product: %+v to cart to the client: %+v with error: %+v",
			productId, clientId, errAddProd))
		return errAddProd
	}
	return nil
}

func (c *CartServiceImpl) GetCart(clientId int) (response.GetCartResponse, error) {
	resp := response.GetCartResponse{}

	if !c.ClientRepository.IsClientInDataBase(clientId) {
		return resp, errors.New(fmt.Sprintf("client with id: %v not found", clientId))
	}

	prods, errGetCart := c.CartRepository.GetCart(clientId)
	if errGetCart != nil {
		logger.Error(fmt.Sprintf("Failed trying to get cart for the client: %+v with error: %+v",
			clientId, errGetCart))
		return resp, errGetCart
	}

	resp.Products = parseProducts(*prods)

	return resp, nil
}

func parseProducts(dbProducts []models.Product) []response.ProductResponse {
	var responseProducts []response.ProductResponse
	for _, e := range dbProducts {
		prod := response.ProductResponse{
			Id:          e.Id,
			Label:       e.Label,
			Category:    e.Category.Label,
			Type:        e.Type,
			DownloadUrl: e.DownloadUrl,
			Weight:      e.Weight,
		}
		responseProducts = append(responseProducts, prod)
	}
	return responseProducts
}
