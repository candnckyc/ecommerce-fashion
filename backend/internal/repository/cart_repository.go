package repository

import (
	"database/sql"
	"ecommerce-backend/internal/models"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

// GetUserCart gets all items in user's cart
func (r *CartRepository) GetUserCart(userID int) ([]models.CartItem, error) {
	query := `
		SELECT c.id, c.user_id, c.product_variant_id, c.quantity, c.created_at, c.updated_at
		FROM cart c
		WHERE c.user_id = $1
		ORDER BY c.created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.CartItem{}
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductVariantID,
			&item.Quantity, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// AddItem adds or updates item in cart
func (r *CartRepository) AddItem(userID, variantID, quantity int) error {
	query := `
		INSERT INTO cart (user_id, product_variant_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_variant_id)
		DO UPDATE SET quantity = cart.quantity + $3, updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(query, userID, variantID, quantity)
	return err
}

// UpdateQuantity updates item quantity
func (r *CartRepository) UpdateQuantity(cartItemID, userID, quantity int) error {
	query := `
		UPDATE cart 
		SET quantity = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND user_id = $3
	`
	_, err := r.db.Exec(query, quantity, cartItemID, userID)
	return err
}

// RemoveItem removes item from cart
func (r *CartRepository) RemoveItem(cartItemID, userID int) error {
	query := `DELETE FROM cart WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, cartItemID, userID)
	return err
}

// ClearCart removes all items from user's cart
func (r *CartRepository) ClearCart(userID int) error {
	query := `DELETE FROM cart WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

// GetVariantByID gets a product variant
func (r *CartRepository) GetVariantByID(variantID int) (*models.ProductVariant, error) {
	variant := &models.ProductVariant{}
	query := `
		SELECT id, product_id, sku, size, color, color_hex, stock_quantity, price_adjustment
		FROM product_variants
		WHERE id = $1
	`
	err := r.db.QueryRow(query, variantID).Scan(
		&variant.ID, &variant.ProductID, &variant.SKU, &variant.Size,
		&variant.Color, &variant.ColorHex, &variant.StockQuantity, &variant.PriceAdjustment,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return variant, err
}
