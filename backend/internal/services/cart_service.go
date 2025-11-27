package services

import (
	"errors"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"
)

type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewCartService(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// GetCart returns user's cart with full details
func (s *CartService) GetCart(userID int) (*models.CartResponse, error) {
	items, err := s.cartRepo.GetUserCart(userID)
	if err != nil {
		return nil, err
	}

	// Populate each item with product and variant details
	totalPrice := 0.0
	totalItems := 0

	for i := range items {
		// Get variant
		variant, err := s.cartRepo.GetVariantByID(items[i].ProductVariantID)
		if err != nil || variant == nil {
			continue
		}
		items[i].Variant = variant

		// Get product
		product, err := s.productRepo.GetByID(variant.ProductID)
		if err != nil || product == nil {
			continue
		}
		items[i].Product = product

		// Calculate price
		itemPrice := product.BasePrice + variant.PriceAdjustment
		variant.FinalPrice = itemPrice
		totalPrice += itemPrice * float64(items[i].Quantity)
		totalItems += items[i].Quantity
	}

	return &models.CartResponse{
		Items:      items,
		TotalItems: totalItems,
		TotalPrice: totalPrice,
	}, nil
}

// AddToCart adds item to cart
func (s *CartService) AddToCart(userID int, req *models.AddToCartRequest) error {
	// Validate
	if req.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// Check if variant exists and has stock
	variant, err := s.cartRepo.GetVariantByID(req.ProductVariantID)
	if err != nil {
		return err
	}
	if variant == nil {
		return errors.New("product variant not found")
	}
	if variant.StockQuantity < req.Quantity {
		return errors.New("insufficient stock")
	}

	// Add to cart
	return s.cartRepo.AddItem(userID, req.ProductVariantID, req.Quantity)
}

// UpdateCartItem updates cart item quantity
func (s *CartService) UpdateCartItem(cartItemID, userID int, req *models.UpdateCartRequest) error {
	if req.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	return s.cartRepo.UpdateQuantity(cartItemID, userID, req.Quantity)
}

// RemoveFromCart removes item from cart
func (s *CartService) RemoveFromCart(cartItemID, userID int) error {
	return s.cartRepo.RemoveItem(cartItemID, userID)
}

// ClearCart removes all items
func (s *CartService) ClearCart(userID int) error {
	return s.cartRepo.ClearCart(userID)
}
