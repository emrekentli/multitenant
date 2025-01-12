package service

import (
	"context"
	"errors"
	"github.com/emrekentli/multitenant-boilerplate/src/model"
	"github.com/emrekentli/multitenant-boilerplate/src/repository"
)

type ProductService struct {
	ProductRepository *repository.ProductRepository
}

func NewProductService(productRepository *repository.ProductRepository) *ProductService {
	return &ProductService{ProductRepository: productRepository}
}

func (s *ProductService) GetAllProducts(tenantSchema string) ([]model.Product, error) {
	// Pass tenant schema to repository
	products, err := s.ProductRepository.FindAll(context.Background(), tenantSchema)
	if err != nil {
		return nil, errors.New("failed to fetch products: " + err.Error())
	}
	return products, nil
}

func (s *ProductService) CreateProduct(tenantSchema string, product *model.Product) error {
	// Pass tenant schema to repository
	err := s.ProductRepository.Create(context.Background(), tenantSchema, product)
	if err != nil {
		return errors.New("failed to create product: " + err.Error())
	}
	return nil
}

func (s *ProductService) GetProductByID(tenantSchema string, id int) (*model.Product, error) {
	// Pass tenant schema to repository
	product, err := s.ProductRepository.FindByID(context.Background(), tenantSchema, id)
	if err != nil {
		return nil, errors.New("product not found: " + err.Error())
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(tenantSchema string, id int, updatedProduct *model.Product) error {
	// Pass tenant schema to repository
	product, err := s.ProductRepository.FindByID(context.Background(), tenantSchema, id)
	if err != nil {
		return errors.New("product not found: " + err.Error())
	}

	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price

	err = s.ProductRepository.Update(context.Background(), tenantSchema, product)
	if err != nil {
		return errors.New("failed to update product: " + err.Error())
	}
	return nil
}

func (s *ProductService) DeleteProduct(tenantSchema string, id int) error {
	// Pass tenant schema to repository
	err := s.ProductRepository.Delete(context.Background(), tenantSchema, id)
	if err != nil {
		return errors.New("failed to delete product: " + err.Error())
	}
	return nil
}
