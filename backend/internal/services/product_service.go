package services

import (
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

// GetProducts returns a list of products with filters
func (s *ProductService) GetProducts(query *models.ProductListQuery) ([]models.Product, error) {
	// Set defaults
	if query.Limit == 0 {
		query.Limit = 20
	}
	if query.Page == 0 {
		query.Page = 1
	}

	products, err := s.productRepo.GetAll(query)
	if err != nil {
		return nil, err
	}

	// Populate images for each product
	for i := range products {
		images, err := s.productRepo.GetImagesByProductID(products[i].ID)
		if err == nil {
			products[i].Images = images
		}
	}

	return products, nil
}

// GetProductByID returns a single product with full details
func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Calculate final prices for variants
	if product != nil && len(product.Variants) > 0 {
		for i := range product.Variants {
			product.Variants[i].FinalPrice = product.BasePrice + product.Variants[i].PriceAdjustment
		}
	}

	return product, nil
}

// GetBrands returns all brands
func (s *ProductService) GetBrands() ([]models.Brand, error) {
	return s.productRepo.GetAllBrands()
}

// GetCategories returns all categories
func (s *ProductService) GetCategories() ([]models.Category, error) {
	return s.productRepo.GetAllCategories()
}