package product

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)

func (r *ProductResolver) GetProductFromDB(id string) (*model.Product, error) {
	var product *model.Product
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductResolver) DeleteProductFromDB(ctx context.Context, id string) (error) {
	product := &model.Product{
		ID: id,
	}
	if err := r.db.Delete(product).Error; err != nil {
		return  err
	}
	return  nil
} 