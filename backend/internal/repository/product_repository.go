package repository

import (
	"database/sql"
	"ecommerce-backend/internal/models"
	"fmt"
	//"strings"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll retrieves products with optional filters
func (r *ProductRepository) GetAll(query *models.ProductListQuery) ([]models.Product, error) {
	// Build dynamic SQL query
	sql := `
		SELECT DISTINCT p.id, p.name, p.slug, p.description, p.brand_id, p.category_id, 
		       p.base_price, p.is_active, p.created_at, p.updated_at
		FROM products p
		LEFT JOIN product_variants pv ON p.id = pv.product_id
		WHERE p.is_active = true
	`
	
	args := []interface{}{}
	argCount := 1

	// Apply filters
	if query.CategoryID != nil {
		sql += fmt.Sprintf(" AND p.category_id = $%d", argCount)
		args = append(args, *query.CategoryID)
		argCount++
	}
	
	if query.BrandID != nil {
		sql += fmt.Sprintf(" AND p.brand_id = $%d", argCount)
		args = append(args, *query.BrandID)
		argCount++
	}
	
	if query.MinPrice != nil {
		sql += fmt.Sprintf(" AND p.base_price >= $%d", argCount)
		args = append(args, *query.MinPrice)
		argCount++
	}
	
	if query.MaxPrice != nil {
		sql += fmt.Sprintf(" AND p.base_price <= $%d", argCount)
		args = append(args, *query.MaxPrice)
		argCount++
	}
	
	if query.Size != "" {
		sql += fmt.Sprintf(" AND pv.size = $%d", argCount)
		args = append(args, query.Size)
		argCount++
	}
	
	if query.Color != "" {
		sql += fmt.Sprintf(" AND pv.color = $%d", argCount)
		args = append(args, query.Color)
		argCount++
	}
	
	if query.Search != "" {
		sql += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.description ILIKE $%d)", argCount, argCount)
		searchTerm := "%" + query.Search + "%"
		args = append(args, searchTerm)
		argCount++
	}

	sql += " ORDER BY p.created_at DESC"

	// Pagination
	if query.Limit > 0 {
		sql += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, query.Limit)
		argCount++
		
		if query.Page > 0 {
			offset := (query.Page - 1) * query.Limit
			sql += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID, &p.Name, &p.Slug, &p.Description, &p.BrandID, &p.CategoryID,
			&p.BasePrice, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// GetByID retrieves a single product with all details
func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	product := &models.Product{}
	
	query := `
		SELECT id, name, slug, description, brand_id, category_id, 
		       base_price, is_active, created_at, updated_at
		FROM products
		WHERE id = $1 AND is_active = true
	`
	
	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Slug, &product.Description,
		&product.BrandID, &product.CategoryID, &product.BasePrice,
		&product.IsActive, &product.CreatedAt, &product.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get brand if exists
	if product.BrandID != nil {
		brand, err := r.getBrandByID(*product.BrandID)
		if err == nil {
			product.Brand = brand
		}
	}

	// Get category if exists
	if product.CategoryID != nil {
		category, err := r.getCategoryByID(*product.CategoryID)
		if err == nil {
			product.Category = category
		}
	}

	// Get variants
	variants, err := r.getVariantsByProductID(product.ID)
	if err == nil {
		product.Variants = variants
	}

	// Get images
	images, err := r.GetImagesByProductID(product.ID)
	if err == nil {
		product.Images = images
	}

	return product, nil
}

// getVariantsByProductID retrieves all variants for a product
func (r *ProductRepository) getVariantsByProductID(productID int) ([]models.ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, size, color, color_hex, stock_quantity, price_adjustment
		FROM product_variants
		WHERE product_id = $1
		ORDER BY size, color
	`
	
	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variants := []models.ProductVariant{}
	for rows.Next() {
		var v models.ProductVariant
		err := rows.Scan(
			&v.ID, &v.ProductID, &v.SKU, &v.Size, &v.Color,
			&v.ColorHex, &v.StockQuantity, &v.PriceAdjustment,
		)
		if err != nil {
			return nil, err
		}
		variants = append(variants, v)
	}

	return variants, nil
}

// getImagesByProductID retrieves all images for a product
// GetImagesByProductID retrieves images for a product
func (r *ProductRepository) GetImagesByProductID(productID int) ([]models.ProductImage, error) {
	query := `
		SELECT id, product_id, image_url, alt_text, display_order, is_primary
		FROM product_images
		WHERE product_id = $1
		ORDER BY display_order, id
	`
	
	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := []models.ProductImage{}
	for rows.Next() {
		var img models.ProductImage
		err := rows.Scan(
			&img.ID, &img.ProductID, &img.ImageURL, &img.AltText,
			&img.DisplayOrder, &img.IsPrimary,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return images, nil
}

// getBrandByID retrieves a brand
func (r *ProductRepository) getBrandByID(id int) (*models.Brand, error) {
	brand := &models.Brand{}
	query := `SELECT id, name, slug, description, logo_url FROM brands WHERE id = $1`
	
	err := r.db.QueryRow(query, id).Scan(
		&brand.ID, &brand.Name, &brand.Slug, &brand.Description, &brand.LogoURL,
	)
	if err != nil {
		return nil, err
	}
	
	return brand, nil
}

// getCategoryByID retrieves a category
func (r *ProductRepository) getCategoryByID(id int) (*models.Category, error) {
	category := &models.Category{}
	query := `SELECT id, name, slug, parent_id, description, image_url FROM categories WHERE id = $1`
	
	err := r.db.QueryRow(query, id).Scan(
		&category.ID, &category.Name, &category.Slug, &category.ParentID,
		&category.Description, &category.ImageURL,
	)
	if err != nil {
		return nil, err
	}
	
	return category, nil
}

// GetAllBrands retrieves all brands
func (r *ProductRepository) GetAllBrands() ([]models.Brand, error) {
	query := `SELECT id, name, slug, COALESCE(description, ''), COALESCE(logo_url, '') FROM brands ORDER BY name`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	brands := []models.Brand{}
	for rows.Next() {
		var b models.Brand
		err := rows.Scan(&b.ID, &b.Name, &b.Slug, &b.Description, &b.LogoURL)
		if err != nil {
			return nil, err
		}
		brands = append(brands, b)
	}

	return brands, nil
}

// GetAllCategories retrieves all categories
func (r *ProductRepository) GetAllCategories() ([]models.Category, error) {
	query := `SELECT id, name, slug, parent_id, COALESCE(description, ''), COALESCE(image_url, '') FROM categories ORDER BY name`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.ParentID, &c.Description, &c.ImageURL)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

// ToggleActive toggles product active status
func (r *ProductRepository) ToggleActive(productID int) error {
	query := `UPDATE products SET is_active = NOT is_active WHERE id = $1`
	_, err := r.db.Exec(query, productID)
	return err
}