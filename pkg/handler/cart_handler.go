package handler

import (
	"github.com/emiliocc5/online-store-api/internal/models/response"
	"github.com/emiliocc5/online-store-api/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type (
	CartHandler interface {
		HandleGetCart(context *gin.Context)
		HandleAddProduct(context *gin.Context)
	}
	CartHandlerImpl struct {
		CartService services.CartService
	}
)

func (ch *CartHandlerImpl) HandleGetCart(context *gin.Context) {
	resp, err := ch.CartService.GetCart(getClientIdFromContext(context))
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.JSON(http.StatusOK, resp)
}

func (ch *CartHandlerImpl) HandleAddProduct(context *gin.Context) {
	err := ch.CartService.AddProductToCart(getProductIdFromContext(context), getClientIdFromContext(context))
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.Status(http.StatusOK)
}

func getClientIdFromContext(context *gin.Context) int {
	clientId := context.GetHeader("clientId")
	intClientId, err1 := strconv.Atoi(clientId)
	if err1 != nil {
		context.AbortWithStatus(http.StatusBadRequest)
	}
	return intClientId
}

func getProductIdFromContext(context *gin.Context) int {
	productId := context.Param("productId")
	intProductId, err := strconv.Atoi(productId)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
	}
	return intProductId
}
