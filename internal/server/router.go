package server

import (
	"github.com/emiliocc5/online-store-api/internal/repository"
	"github.com/emiliocc5/online-store-api/internal/services"
	"github.com/emiliocc5/online-store-api/pkg/handler"
	"github.com/gin-gonic/gin"
)

const (
	BaseEndpoint = "/api"
	Cart         = "/cart"
	AddProduct   = "/products/:productId"
)

var (
	cartHandler handler.CartHandler
)

func ConfigureRouter(engine *gin.Engine) {
	configureCartRoutes(engine)
}

func configureCartRoutes(engine *gin.Engine) {
	engine.GET(BaseEndpoint+Cart, cartHandler.HandleGetCart)
	engine.POST(BaseEndpoint+Cart+AddProduct, cartHandler.HandleAddProduct)
}

func init() {
	cartHandler = &handler.CartHandlerImpl{
		CartService: &services.CartServiceImpl{
			CartRepository: repository.GetCartRepository(),
		},
	}
}
