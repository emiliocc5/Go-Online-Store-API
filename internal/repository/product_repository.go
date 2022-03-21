package repository

import (
	"errors"
	"fmt"
	"github.com/emiliocc5/online-store-api/internal/models"
	"gorm.io/gorm"
)

type (
	ProductRepository interface {
		FindProductsFromCart(cartId int) (*[]models.Product, error)
		FindProductById(productId int) (*models.Product, error)
	}
	PgProductRepository struct {
		DbClient IProductRepositoryDbClient
	}
	IProductRepositoryDbClient interface {
		Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	}
)

func (pr *PgProductRepository) FindProductsFromCart(cartId int) (*[]models.Product, error) {
	var productsCarts []models.ProductCart
	err := pr.findListOfProductsForCartId(&productsCarts, cartId)
	if err != nil {
		return nil, err
	}

	var productList []models.Product
	pr.findProductsFromProductCartsList(productsCarts, &productList)

	return &productList, nil
}

func (pr *PgProductRepository) FindProductById(productId int) (*models.Product, error) {
	product := models.Product{}
	productResult := pr.DbClient.Find(&product, "id = ?", productId)
	if productResult.Error != nil {
		logger.Errorf("unable to find the product with error: %v", productResult.Error)
		return nil, errors.New("unable to find the product in our db")
	}
	return &product, nil
}

func (pr *PgProductRepository) findListOfProductsForCartId(productCarts *[]models.ProductCart, clientCartId int) error {
	productsResult := pr.DbClient.Find(productCarts, "cart_id = ?", clientCartId)
	if productsResult.Error != nil {
		return errors.New("unable to retrieve the list of products")
	}
	return nil
}

func (pr *PgProductRepository) findProductsFromProductCartsList(productCarts []models.ProductCart, productList *[]models.Product) {
	for _, e := range productCarts {
		product, err := pr.FindProductById(e.ProductId)
		if err != nil {
			logger.Error(fmt.Sprintf("unable to get the product: %v", err.Error()))
		} else {
			*productList = append(*productList, *product)
		}
	}
}
