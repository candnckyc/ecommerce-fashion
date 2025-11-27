package models

import "time"

// Product represents a product in the catalog
type Product struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description string           `json:"description"`
	BrandID     *int             `json:"brand_id"`
	Brand       *Brand           `json:"brand,omitempty"`
	CategoryID  *int             `json:"category_id"`
	Category    *Category        `json:"category,omitempty"`
	BasePrice   float64          `json:"base_price"`
	IsActive    bool             `json:"is_active"`
	Variants    []ProductVariant `json:"variants,omitempty"`
	Images      []ProductImage   `json:"images,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// ProductVariant represents a size/color combination
type ProductVariant struct {
	ID              int     `json:"id"`
	ProductID       int     `json:"product_id"`
	SKU             string  `json:"sku"`
	Size            string  `json:"size"`
	Color           string  `json:"color"`
	ColorHex        string  `json:"color_hex"`
	StockQuantity   int     `json:"stock_quantity"`
	PriceAdjustment float64 `json:"price_adjustment"`
	FinalPrice      float64 `json:"final_price"` // Calculated: base_price + price_adjustment
}

// ProductImage represents a product image
type ProductImage struct {
	ID           int    `json:"id"`
	ProductID    int    `json:"product_id"`
	ImageURL     string `json:"image_url"`
	AltText      string `json:"alt_text"`
	DisplayOrder int    `json:"display_order"`
	IsPrimary    bool   `json:"is_primary"`
}

// Brand represents a fashion brand
type Brand struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	LogoURL     string `json:"logo_url,omitempty"`
}

// Category represents a product category
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	ParentID    *int   `json:"parent_id"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

// ProductListQuery holds query parameters for product listing
type ProductListQuery struct {
	CategoryID *int    `json:"category_id"`
	BrandID    *int    `json:"brand_id"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	Size       string  `json:"size"`
	Color      string  `json:"color"`
	Search     string  `json:"search"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
}
