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
	if err := r.db.Where("user_id = ?", userId).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductResolver) SearchProducts(ctx context.Context, query *string) ([]*model.Product, error) {
	userId, err := contextUtil.UserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var products []*model.Product

	db := r.db // assuming you have a `*gorm.DB` instance in your resolver

	q := db.
		Where("deleted_at IS NULL").
		Where("user_id = ?", userId)

	if query != nil && *query != "" {
		// This will use the idx_products_name index if the LIKE pattern is left-anchored (e.g., "term%")
		q = q.Where("name ILIKE ?", *query+"%")
	}

	if err := q.Find(&products).Error; err != nil {
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