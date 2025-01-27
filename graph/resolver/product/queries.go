package product

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/contextUtil"
)

func (r *ProductResolver) Products(ctx context.Context) ([]*model.Product, error) {
	userId, err := contextUtil.UserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var products []*model.Product
	if err := r.db.Where("user_id = ?", userId).First(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}


func (r *ProductResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	product, err := r.GetProductFromDB(id) 
	if err != nil {
		return nil, err
	}
	return product, nil
}