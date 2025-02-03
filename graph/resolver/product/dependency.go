package product

import (
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type ProductResolver struct {
	db *gorm.DB
	httpClient *resty.Client

}

func InitializeProductResolver( db *gorm.DB, httpClient *resty.Client) *ProductResolver {
	return &ProductResolver{
		db: db,
		httpClient: httpClient,
	}
}

const ProductsKey = "products"
const ProductKey = "product"
