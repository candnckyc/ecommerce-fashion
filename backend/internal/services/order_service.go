package services

import (
	"errors"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// CreateOrder creates an order from user's cart
func (s *OrderService) CreateOrder(userID int, req *models.CreateOrderRequest) (*models.Order, error) {
	// Validate
	if req.AddressID == 0 {
		return nil, errors.New("address is required")
	}
	if req.PaymentMethod == "" {
		return nil, errors.New("payment method is required")
	}

	// Get address
	address, err := s.orderRepo.GetAddressByID(req.AddressID, userID)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, errors.New("address not found")
	}

	// Get cart
	cartItems, err := s.cartRepo.GetUserCart(userID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Calculate totals and validate stock
	subtotal := 0.0
	orderItems := []models.OrderItem{}
	
	for _, cartItem := range cartItems {
		// Get variant
		variant, err := s.cartRepo.GetVariantByID(cartItem.ProductVariantID)
		if err != nil || variant == nil {
			continue
		}

		// Check stock
		if variant.StockQuantity < cartItem.Quantity {
			return nil, errors.New("insufficient stock for " + variant.SKU)
		}

		// Get product
		product, err := s.productRepo.GetByID(variant.ProductID)
		if err != nil || product == nil {
			continue
		}

		// Calculate price
		unitPrice := product.BasePrice + variant.PriceAdjustment
		totalPrice := unitPrice * float64(cartItem.Quantity)
		subtotal += totalPrice

		// Create order item
		orderItems = append(orderItems, models.OrderItem{
			ProductVariantID: variant.ID,
			ProductName:      product.Name,
			ProductSKU:       variant.SKU,
			Size:             variant.Size,
			Color:            variant.Color,
			Quantity:         cartItem.Quantity,
			UnitPrice:        unitPrice,
			TotalPrice:       totalPrice,
		})
	}

	// Create order
	order := &models.Order{
		UserID:               userID,
		OrderNumber:          s.orderRepo.GenerateOrderNumber(),
		ShippingAddressLine1: address.AddressLine1,
		ShippingAddressLine2: address.AddressLine2,
		ShippingCity:         address.City,
		ShippingState:        address.State,
		ShippingPostalCode:   address.PostalCode,
		ShippingCountry:      address.Country,
		ShippingFullName:     address.FullName,
		ShippingPhone:        address.Phone,
		Subtotal:             subtotal,
		ShippingCost:         0.00, // Free shipping
		Tax:                  0.00, // No tax for now
		Total:                subtotal,
		Status:               "pending",
		PaymentMethod:        req.PaymentMethod,
		PaymentStatus:        "pending",
		Notes:                req.Notes,
	}

	orderID, err := s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	// Create order items and reduce stock
	for i := range orderItems {
		orderItems[i].OrderID = orderID
		err := s.orderRepo.CreateOrderItem(&orderItems[i])
		if err != nil {
			return nil, err
		}
		
		// Reduce stock for this variant
		err = s.productRepo.ReduceStock(orderItems[i].ProductVariantID, orderItems[i].Quantity)
		if err != nil {
			return nil, err
		}
	}

	order.Items = orderItems

	// Clear cart
	s.cartRepo.ClearCart(userID)

	return order, nil
}

// GetUserOrders retrieves all orders for a user
func (s *OrderService) GetUserOrders(userID int) ([]models.Order, error) {
	return s.orderRepo.GetUserOrders(userID)
}

// GetOrderByID retrieves a single order
func (s *OrderService) GetOrderByID(orderID, userID int) (*models.Order, error) {
	return s.orderRepo.GetOrderByID(orderID, userID)
}

// GetUserAddresses retrieves all addresses
func (s *OrderService) GetUserAddresses(userID int) ([]models.Address, error) {
	return s.orderRepo.GetUserAddresses(userID)
}

// CreateAddress creates a new address
func (s *OrderService) CreateAddress(userID int, req *models.AddAddressRequest) (*models.Address, error) {
	if req.FullName == "" || req.AddressLine1 == "" || req.City == "" || req.Country == "" {
		return nil, errors.New("required fields missing")
	}

	addr := &models.Address{
		UserID:       userID,
		Title:        req.Title,
		FullName:     req.FullName,
		Phone:        req.Phone,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		Country:      req.Country,
		IsDefault:    req.IsDefault,
	}

	id, err := s.orderRepo.CreateAddress(addr)
	if err != nil {
		return nil, err
	}
	addr.ID = id

	return addr, nil
}

// GetAllOrders returns all orders (admin only)
func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAllOrders()
}

// UpdateOrderStatus updates order status (admin only)
func (s *OrderService) UpdateOrderStatus(orderID int, status string) error {
	return s.orderRepo.UpdateOrderStatus(orderID, status)
}