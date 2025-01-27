package product

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/database"
)

func (r *ProductResolver) CreateProduct(ctx context.Context, input model.CreateProduct) (*model.Product, error) {
	inputData := database.CreateFields[model.Product](input);

	if err := r.db.Create(inputData).Error; err != nil {
		return nil, err
	}
	if err := addProductDocument(r.httpClient, *inputData); err != nil {
		return nil, err
	}
	return inputData, nil
}

func (r *ProductResolver) UpdateProduct(ctx context.Context, input model.UpdateProduct) (*model.Product, error) {	
	service := &model.Product{
		ID: input.ID,
	}
	
	if err := r.db.Model(service).Updates(input).Error; err != nil {
		return nil, err
	}
	newProduct, err := r.GetProductFromDB(input.ID)
	if err != nil {
		return nil, err
	}

	return newProduct, nil
}

func (r *ProductResolver) DeleteProduct(ctx context.Context, id string) (bool, error) {
	if err := r.DeleteProductFromDB(ctx, id); err != nil {
		return false, err
	}
	if err := deleteProductDocument(r.httpClient, id); err != nil {
		return false, err
	}
	return true, nil
}