package models

import "time"

// Order represents an order
type Order struct {
	ID           int         `json:"id"`
	UserID       int         `json:"user_id"`
	OrderNumber  string      `json:"order_number"`
	
	// Address snapshot (denormalized)
	ShippingAddressLine1 string `json:"shipping_address_line1"`
	ShippingAddressLine2 string `json:"shipping_address_line2"`
	ShippingCity         string `json:"shipping_city"`
	ShippingState        string `json:"shipping_state"`
	ShippingPostalCode   string `json:"shipping_postal_code"`
	ShippingCountry      string `json:"shipping_country"`
	ShippingFullName     string `json:"shipping_full_name"`
	ShippingPhone        string `json:"shipping_phone"`
	
	// Order totals
	Subtotal       float64 `json:"subtotal"`
	ShippingCost   float64 `json:"shipping_cost"`
	Tax            float64 `json:"tax"`
	Total          float64 `json:"total"`
	
	// Payment & Status
	Status              string `json:"status"` // pending, confirmed, processing, shipped, delivered, cancelled
	PaymentMethod       string `json:"payment_method"`
	PaymentStatus       string `json:"payment_status"` // pending, paid, failed, refunded
	PaymentTransactionID string `json:"payment_transaction_id,omitempty"`
	
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Related data
	Items []OrderItem `json:"items,omitempty"`
}

// OrderItem represents a single item in an order
type OrderItem struct {
	ID               int       `json:"id"`
	OrderID          int       `json:"order_id"`
	ProductVariantID int       `json:"product_variant_id"`
	
	// Product snapshot (denormalized)
	ProductName string `json:"product_name"`
	ProductSKU  string `json:"product_sku"`
	Size        string `json:"size"`
	Color       string `json:"color"`
	
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
	
	CreatedAt time.Time `json:"created_at"`
}

// CreateOrderRequest is the request to create an order
type CreateOrderRequest struct {
	AddressID     int    `json:"address_id"`
	PaymentMethod string `json:"payment_method"` // credit_card, debit_card, cash_on_delivery
	Notes         string `json:"notes"`
}

// Address represents a shipping address
type Address struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	Title        string `json:"title"` // Home, Work, etc.
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	IsDefault    bool   `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AddAddressRequest is the request to add an address
type AddAddressRequest struct {
	Title        string `json:"title"`
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	IsDefault    bool   `json:"is_default"`
}