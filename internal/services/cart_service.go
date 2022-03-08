package services

import (
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models/response"
	"github.com/emiliocc5/online-store-api/internal/repository"
)

//TODO CHANGE TO FORMAT INPUT IN HANDLER

type CartService interface {
	AddProduct(productId, clientId int) error
	GetCart(clientId int) (response.GetCartResponse, error)
}
type CartServiceImpl struct {
	CartRepository repository.CartRepository
}

func NewCartService() *CartServiceImpl {
	return &CartServiceImpl{
		CartRepository: repository.GetCartRepository(),
	}
}

func (c *CartServiceImpl) AddProduct(productId, clientId int) error {
	errAddProd := c.CartRepository.AddProductToCart(productId, clientId)
	if errAddProd != nil {
		fmt.Println(fmt.Sprintf("Failed trying to add product: %+v to cart to the client: %+v with error: %+v, and message: %+v ",
			productId, clientId, errAddProd, errAddProd.Error()))
		return errAddProd
	}
	return nil
}

func (c *CartServiceImpl) GetCart(clientId int) (response.GetCartResponse, error) {
	resp := response.GetCartResponse{}

	prods, err := c.CartRepository.GetCart(clientId)

	if err != nil {
		return resp, err
	}

	resp.Products = *prods

	return resp, nil
}
