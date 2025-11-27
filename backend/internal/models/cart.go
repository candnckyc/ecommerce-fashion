package models

import "time"

// CartItem represents an item in user's cart
type CartItem struct {
	ID               int             `json:"id"`
	UserID           int             `json:"user_id"`
	ProductVariantID int             `json:"product_variant_id"`
	Quantity         int             `json:"quantity"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	
	// Populated fields (from joins)
	Product *Product         `json:"product,omitempty"`
	Variant *ProductVariant  `json:"variant,omitempty"`
}

// AddToCartRequest is the request to add item to cart
type AddToCartRequest struct {
	ProductVariantID int `json:"product_variant_id"`
	Quantity         int `json:"quantity"`
}

// UpdateCartRequest is the request to update cart item quantity
type UpdateCartRequest struct {
	Quantity int `json:"quantity"`
}

// CartResponse is the full cart with all items
type CartResponse struct {
	Items      []CartItem `json:"items"`
	TotalItems int        `json:"total_items"`
	TotalPrice float64    `json:"total_price"`
}
