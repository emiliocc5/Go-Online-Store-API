package services

import (
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models/response"
	"github.com/emiliocc5/online-store-api/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//TODO CHANGE TO FORMAT INPUT IN HANDLER

type CartService interface {
	AddProduct(context *gin.Context)
	GetCart(context *gin.Context)
}
type CartServiceImpl struct {
	CartRepository repository.CartRepository
}

func (c *CartServiceImpl) AddProduct(context *gin.Context) {
	productId := context.Param("productId")
	clientId := context.GetHeader("clientId")

	intProductId, err := strconv.Atoi(productId)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	intClientId, err1 := strconv.Atoi(clientId)
	if err1 != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	errAddProd := c.CartRepository.AddProductToCart(intProductId, intClientId)
	if errAddProd != nil {
		fmt.Println(fmt.Sprintf("Failed trying to add product: %+v to cart to the client: %+v with error: %+v, and message: %+v ",
			intProductId, intClientId, err, err.Error()))
	}
	context.Status(http.StatusOK)
}

func (c *CartServiceImpl) GetCart(context *gin.Context) {
	resp := response.GetCartResponse{}

	prods, err := c.CartRepository.GetCart(1)

	if err != nil {
		fmt.Println("Hubo un error")
	}

	resp.Products = *prods

	context.JSON(http.StatusOK, resp)
}
