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
		GetCart(clientId int) (response.GetCartResponse, error)
		AddProductToCart(productId, clientId int) error
	}
	CartServiceImpl struct {
		CartRepository    repository.CartRepository
		ClientRepository  repository.ClientRepository
		ProductRepository repository.ProductRepository
	}
)

func init() {
	logger = utils.GetLogger()
}

func (c *CartServiceImpl) GetCart(clientId int) (response.GetCartResponse, error) {
	resp := response.GetCartResponse{}

	if !c.ClientRepository.IsClientInDataBase(clientId) {
		return resp, errors.New(fmt.Sprintf("client with id: %v not found", clientId))
	}

	cart, errGetCart := c.CartRepository.GetCartByClient(clientId)
	if errGetCart != nil {
		logger.Error(fmt.Sprintf("Failed trying to get cart for the client: %+v with error: %+v",
			clientId, errGetCart))
		return resp, errGetCart
	}

	prods, errGetProds := c.ProductRepository.FindProductsFromCart(cart.Id)
	if errGetProds != nil {
		logger.Errorf("unable to get the list of products from the cart: %v, with error: %v", cart.Id, errGetProds)
		return resp, errGetProds
	}
	resp.Products = parseProducts(*prods)

	return resp, nil
}

func (c *CartServiceImpl) AddProductToCart(productId, clientId int) error {
	if !c.ClientRepository.IsClientInDataBase(clientId) {
		return errors.New(fmt.Sprintf("client with id: %v not found", clientId))
	}

	product, err := c.ProductRepository.FindProductById(productId)
	if err != nil {
		logger.Error(fmt.Sprintf("unable to get the product: %v", err.Error()))
		return errors.New("unable to find the product in our db")
	}

	errAddProd := c.CartRepository.AddProductToCart(product.Id, clientId)
	if errAddProd != nil {
		logger.Error(fmt.Sprintf("Failed trying to add product: %+v to cart to the client: %+v with error: %+v",
			productId, clientId, errAddProd))
		return errAddProd
	}
	return nil
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
