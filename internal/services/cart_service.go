package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CartService interface {
	AddProduct(productId int)
	GetCart(context *gin.Context)
}
type CartServiceImpl struct {
}

func (c *CartServiceImpl) AddProduct(productId int) {
	fmt.Println(productId)
}

func (c *CartServiceImpl) GetCart(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"data": "hello world"})
}
