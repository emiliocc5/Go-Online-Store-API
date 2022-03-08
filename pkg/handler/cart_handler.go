package handler

import (
	"github.com/emiliocc5/online-store-api/internal/services"
	"github.com/gin-gonic/gin"
)

type (
	CartHandler interface {
		HandleGetCart(context *gin.Context)
		HandleAddProduct(context *gin.Context)
	}
	CartHandlerImpl struct {
		CartService services.CartService //TODO change this to initialize alone
	}
)

//TODO FORMAT CONTEXT IN HANDLER AND PASS VALUES

func (ch *CartHandlerImpl) HandleGetCart(context *gin.Context) {
	ch.CartService.GetCart(context)
}

func (ch *CartHandlerImpl) HandleAddProduct(context *gin.Context) {
	ch.CartService.AddProduct(context)
}
