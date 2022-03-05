package handler

import (
	"github.com/emiliocc5/online-store-api/internal/services"
	"github.com/gin-gonic/gin"
)

type (
	CartHandler interface {
		HandleGetCart(context *gin.Context)
		HandleAddProduct(productId int)
	}
	CartHandlerImpl struct {
		CartService services.CartService //TODO change this to initialize alone
	}
)

func (ch *CartHandlerImpl) HandleGetCart(context *gin.Context) {
	ch.CartService.GetCart(context)
}

func (ch *CartHandlerImpl) HandleAddProduct(productId int) {
	ch.CartService.AddProduct(productId)
}
