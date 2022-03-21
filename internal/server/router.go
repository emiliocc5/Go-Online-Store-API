package server

import (
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/repository"
	"github.com/emiliocc5/online-store-api/internal/services"
	"github.com/emiliocc5/online-store-api/pkg/handler"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	BaseEndpoint = "/api"
	Cart         = "/cart"
	AddProduct   = "/products/:productId"
)

var (
	logger      *logrus.Logger
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
	client, err := repository.GetClient()
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting DbClient: %v", err.Error()))
	}
	cartHandler = &handler.CartHandlerImpl{
		CartService: &services.CartServiceImpl{
			CartRepository: &repository.PgCartRepository{
				DbClient: client,
			},
			ClientRepository: &repository.PgClientRepository{
				DbClient: client,
			},
			ProductRepository: &repository.PgProductRepository{
				DbClient: client,
			},
		},
	}
}
