package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Old database
	oldDB := "postgresql://idle_eticaret2:iDLe-1233a@185.26.144.214:5432/idle_eticaret2?sslmode=disable"
	// New database
	newDB := "postgresql://idle_eticaret:kL8zHtqX-1iYArLM@51.178.9.8:25600/idle_eticaret?sslmode=disable"

	fmt.Println("🔄 Starting database migration...")
	fmt.Println("FROM:", oldDB)
	fmt.Println("TO:", newDB)
	fmt.Println()

	// Connect to old database
	fmt.Println("📡 Connecting to OLD database...")
	dbOld, err := sql.Open("postgres", oldDB)
	if err != nil {
		log.Fatal("Failed to connect to old database:", err)
	}
	defer dbOld.Close()

	err = dbOld.Ping()
	if err != nil {
		log.Fatal("Failed to ping old database:", err)
	}
	fmt.Println("✓ Connected to old database")

	// Connect to new database
	fmt.Println("📡 Connecting to NEW database...")
	dbNew, err := sql.Open("postgres", newDB)
	if err != nil {
		log.Fatal("Failed to connect to new database:", err)
	}
	defer dbNew.Close()

	err = dbNew.Ping()
	if err != nil {
		log.Fatal("Failed to ping new database:", err)
	}
	fmt.Println("✓ Connected to new database")
	fmt.Println()

	// First, create tables in new database
	fmt.Println("🏗️  Creating tables in new database...")
	createTables(dbNew)
	fmt.Println()

	// Migrate data
	fmt.Println("📦 Migrating data...")
	
	migrateUsers(dbOld, dbNew)
	migrateBrands(dbOld, dbNew)
	migrateCategories(dbOld, dbNew)
	migrateProducts(dbOld, dbNew)
	migrateProductVariants(dbOld, dbNew)
	migrateProductImages(dbOld, dbNew)
	migrateAddresses(dbOld, dbNew)
	migrateOrders(dbOld, dbNew)
	migrateOrderItems(dbOld, dbNew)
	migrateCart(dbOld, dbNew)
	migrateWishlist(dbOld, dbNew)

	fmt.Println()
	fmt.Println("✅ Migration complete!")
	fmt.Println("Your new database is ready to use!")
}

func createTables(db *sql.DB) {
	// Create all tables (same as setup script)
	queries := []string{
		// Users
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			phone VARCHAR(20),
			role VARCHAR(20) DEFAULT 'customer',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Brands
		`CREATE TABLE IF NOT EXISTS brands (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL,
			slug VARCHAR(100) UNIQUE NOT NULL,
			description TEXT,
			logo_url VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Categories
		`CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			slug VARCHAR(100) UNIQUE NOT NULL,
			parent_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
			description TEXT,
			image_url VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Products
		`CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL,
			description TEXT,
			brand_id INTEGER REFERENCES brands(id) ON DELETE SET NULL,
			category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
			base_price DECIMAL(10, 2) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			search_vector tsvector,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Product Variants
		`CREATE TABLE IF NOT EXISTS product_variants (
			id SERIAL PRIMARY KEY,
			product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
			sku VARCHAR(100) UNIQUE NOT NULL,
			size VARCHAR(20) NOT NULL,
			color VARCHAR(50) NOT NULL,
			color_hex VARCHAR(7),
			stock_quantity INTEGER NOT NULL DEFAULT 0,
			price_adjustment DECIMAL(10, 2) DEFAULT 0.00,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT check_stock CHECK (stock_quantity >= 0)
		)`,
		
		// Product Images
		`CREATE TABLE IF NOT EXISTS product_images (
			id SERIAL PRIMARY KEY,
			product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
			image_url VARCHAR(500) NOT NULL,
			alt_text VARCHAR(255),
			display_order INTEGER DEFAULT 0,
			is_primary BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Addresses
		`CREATE TABLE IF NOT EXISTS addresses (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(50) NOT NULL,
			full_name VARCHAR(200) NOT NULL,
			phone VARCHAR(20) NOT NULL,
			address_line1 VARCHAR(255) NOT NULL,
			address_line2 VARCHAR(255),
			city VARCHAR(100) NOT NULL,
			state VARCHAR(100),
			postal_code VARCHAR(20) NOT NULL,
			country VARCHAR(100) NOT NULL DEFAULT 'Turkey',
			is_default BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Orders
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
			order_number VARCHAR(50) UNIQUE NOT NULL,
			shipping_address_line1 VARCHAR(255) NOT NULL,
			shipping_address_line2 VARCHAR(255),
			shipping_city VARCHAR(100) NOT NULL,
			shipping_state VARCHAR(100),
			shipping_postal_code VARCHAR(20) NOT NULL,
			shipping_country VARCHAR(100) NOT NULL,
			shipping_full_name VARCHAR(200) NOT NULL,
			shipping_phone VARCHAR(20) NOT NULL,
			subtotal DECIMAL(10, 2) NOT NULL,
			shipping_cost DECIMAL(10, 2) DEFAULT 0.00,
			tax DECIMAL(10, 2) DEFAULT 0.00,
			total DECIMAL(10, 2) NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			payment_method VARCHAR(50),
			payment_status VARCHAR(50) DEFAULT 'pending',
			payment_transaction_id VARCHAR(255),
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Order Items
		`CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
			product_variant_id INTEGER REFERENCES product_variants(id) ON DELETE SET NULL,
			product_name VARCHAR(255) NOT NULL,
			product_sku VARCHAR(100) NOT NULL,
			size VARCHAR(20) NOT NULL,
			color VARCHAR(50) NOT NULL,
			quantity INTEGER NOT NULL,
			unit_price DECIMAL(10, 2) NOT NULL,
			total_price DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Cart
		`CREATE TABLE IF NOT EXISTS cart (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			product_variant_id INTEGER REFERENCES product_variants(id) ON DELETE CASCADE,
			quantity INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT check_cart_quantity CHECK (quantity > 0),
			UNIQUE(user_id, product_variant_id)
		)`,
		
		// Wishlist
		`CREATE TABLE IF NOT EXISTS wishlist (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, product_id)
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Println("Warning:", err)
		}
	}
	
	fmt.Println("✓ Tables created")
}

func migrateUsers(old, new *sql.DB) {
	fmt.Println("  → Migrating users...")
	rows, err := old.Query("SELECT id, email, password_hash, first_name, last_name, phone, role, created_at, updated_at FROM users")
	if err != nil {
		log.Println("Warning: Could not read users:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var email, passwordHash, firstName, lastName, role string
		var phone sql.NullString
		var createdAt, updatedAt sql.NullTime

		rows.Scan(&id, &email, &passwordHash, &firstName, &lastName, &phone, &role, &createdAt, &updatedAt)

		_, err = new.Exec(
			"INSERT INTO users (id, email, password_hash, first_name, last_name, phone, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (email) DO NOTHING",
			id, email, passwordHash, firstName, lastName, phone, role, createdAt, updatedAt,
		)
		if err == nil {
			count++
		}
	}
	fmt.Printf("  ✓ Migrated %d users\n", count)
}

func migrateBrands(old, new *sql.DB) {
	fmt.Println("  → Migrating brands...")
	rows, err := old.Query("SELECT id, name, slug, description, logo_url, created_at, updated_at FROM brands")
	if err != nil {
		log.Println("Warning: Could not read brands:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name, slug string
		var description, logoURL sql.NullString
		var createdAt, updatedAt sql.NullTime

		rows.Scan(&id, &name, &slug, &description, &logoURL, &createdAt, &updatedAt)

		_, err = new.Exec(
			"INSERT INTO brands (id, name, slug, description, logo_url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (slug) DO NOTHING",
			id, name, slug, description, logoURL, createdAt, updatedAt,
		)
		if err == nil {
			count++
		}
	}
	fmt.Printf("  ✓ Migrated %d brands\n", count)
}

func migrateCategories(old, new *sql.DB) {
	fmt.Println("  → Migrating categories...")
	rows, err := old.Query("SELECT id, name, slug, parent_id, description, image_url, created_at, updated_at FROM categories")
	if err != nil {
		log.Println("Warning: Could not read categories:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name, slug string
		var parentID sql.NullInt64
		var description, imageURL sql.NullString
		var createdAt, updatedAt sql.NullTime

		rows.Scan(&id, &name, &slug, &parentID, &description, &imageURL, &createdAt, &updatedAt)

		_, err = new.Exec(
			"INSERT INTO categories (id, name, slug, parent_id, description, image_url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (slug) DO NOTHING",
			id, name, slug, parentID, description, imageURL, createdAt, updatedAt,
		)
		if err == nil {
			count++
		}
	}
	fmt.Printf("  ✓ Migrated %d categories\n", count)
}

func migrateProducts(old, new *sql.DB) {
	fmt.Println("  → Migrating products...")
	rows, err := old.Query("SELECT id, name, slug, description, brand_id, category_id, base_price, is_active, created_at, updated_at FROM products")
	if err != nil {
		log.Println("Warning: Could not read products:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name, slug string
		var description sql.NullString
		var brandID, categoryID sql.NullInt64
		var basePrice float64
		var isActive bool
		var createdAt, updatedAt sql.NullTime

		rows.Scan(&id, &name, &slug, &description, &brandID, &categoryID, &basePrice, &isActive, &createdAt, &updatedAt)

		_, err = new.Exec(
			"INSERT INTO products (id, name, slug, description, brand_id, category_id, base_price, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (slug) DO NOTHING",
			id, name, slug, description, brandID, categoryID, basePrice, isActive, createdAt, updatedAt,
		)
		if err == nil {
			count++
		}
	}
	fmt.Printf("  ✓ Migrated %d products\n", count)
}

func migrateProductVariants(old, new *sql.DB) {
	fmt.Println("  → Migrating product variants...")
	rows, err := old.Query("SELECT id, product_id, sku, size, color, color_hex, stock_quantity, price_adjustment, created_at, updated_at FROM product_variants")
	if err != nil {
		log.Println("Warning: Could not read product variants:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, productID, stockQuantity int
		var sku, size, color string
		var colorHex sql.NullString
		var priceAdjustment float64
		var createdAt, updatedAt sql.NullTime

		rows.Scan(&id, &productID, &sku, &size, &color, &colorHex, &stockQuantity, &priceAdjustment, &createdAt, &updatedAt)

		_, err = new.Exec(
			"INSERT INTO product_variants (id, product_id, sku, size, color, color_hex, stock_quantity, price_adjustment, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (sku) DO NOTHING",
			id, productID, sku, size, color, colorHex, stockQuantity, priceAdjustment, createdAt, updatedAt,
		)
		if err == nil {
			count++
		}
	}
	fmt.Printf("  ✓ Migrated %d product variants\n", count)
}

func migrateProductImages(old, new *sql.DB) {
	fmt.Println("  → Migrating product images...")
	rows, err := old.Query("SELECT id, product_id, image_url, alt_text, display_order, is_primary, created_at FROM product_images")
	if err != nil {
		log.Println("Warning: Could not read product images:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, productID, displayOrder int
		var imageURL string
		var altText sql.NullString
		var isPrimary bool
		var createdAt sql.NullTime

		rows.Scan(&id, &productID, &imageURL, &altText, &displayOrder, &isPrimary, &createdAt)

		_, err = new.Exec(
			"INSERT INTO product_images (id, product_id, image_url, alt_text, display_order, is_primary, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			id, productID, imageURL, altText, displayOrder, isPrimary, createdAt,
		)
		if err == nil {
			count++
		}
	}
	fmt.Printf("  ✓ Migrated %d product images\n", count)
}

func migrateAddresses(old, new *sql.DB) {
	fmt.Println("  → Migrating addresses...")
	rows, err := old.Query("SELECT id, user_id, title, full_name, phone, address_line1, address_line2, city, state, postal_code, country, is_default, created_at, updated_at FROM addresses")
	if err != nil {
		log.Println("Warning: Could not read addresses:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, userID int
		var title, fullName, phone, addressLine1, city, postalCode, country string
		var addressLine2, state sql.NullString
		var isDefault bool
		var createdAt, updatedAt sql.NullTime

		rows.Scan(&id, &userID, &title, &fullName, &phone, &addressLine1, &addressLine2, &city, &state, &postalCode, &country, &isDefault, &createdAt, &updatedAt)

		_, err = new.Exec(
			"INSERT INTO addresses (id, user_id, title, full_name, phone, address_line1, address_line2, city, state, postal_code, country, is_default, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
			id, userID, title, fullName, phone, addressLine1, addressLine2, city, state, postalCode, country, isDefault, createdAt, updatedAt,
		)
		if err == nil {
			count++
		}
	}
	fmt.Printf("  ✓ Migrated %d addresses\n", count)
}

func migrateOrders(old, new *sql.DB) {
	fmt.Println("  → Migrating orders...")
	rows, err := old.Query(`SELECT id, user_id, order_number, shipping_address_line1, shipping_address_line2, 
		shipping_city, shipping_state, shipping_postal_code, shipping_country, shipping_full_name, shipping_phone,
		subtotal, shipping_cost, tax, total, status, payment_method, payment_status, payment_transaction_id, 
		notes, created_at, updated_at FROM orders`)
	if err != nil {
		log.Println("Warning: Could not read orders:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var userID sql.NullInt64
		var orderNumber, shippingAddressLine1, shippingCity, shippingPostalCode, shippingCountry, shippingFullName, shippingPhone, status string
		var shippingAddressLine2, shippingState, paymentMethod, paymentStatus, paymentTransactionID, notes sql.NullString
		var subtotal, shippingCost, tax, total float64
		var createdAt, updatedAt sql.NullTime

		rows.Scan(&id, &userID, &orderNumber, &shippingAddressLine1, &shippingAddressLine2, 
			&shippingCity, &shippingState, &shippingPostalCode, &shippingCountry, &shippingFullName, &shippingPhone,
			&subtotal, &shippingCost, &tax, &total, &status, &paymentMethod, &paymentStatus, &paymentTransactionID,
			&notes, &createdAt, &updatedAt)

		_, err = new.Exec(
			`INSERT INTO orders (id, user_id, order_number, shipping_address_line1, shipping_address_line2,
			shipping_city, shipping_state, shipping_postal_code, shipping_country, shipping_full_name, shipping_phone,
			subtotal, shipping_cost, tax, total, status, payment_method, payment_status, payment_transaction_id,
			notes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22) ON CONFLICT (order_number) DO NOTHING`,
			id, userID, orderNumber, shippingAddressLine1, shippingAddressLine2,
			shippingCity, shippingState, shippingPostalCode, shippingCountry, shippingFullName, shippingPhone,
			subtotal, shippingCost, tax, total, status, paymentMethod, paymentStatus, paymentTransactionID,
			notes, createdAt, updatedAt,
		)
		if err == nil {
			count++
		}
	}
	fmt.Printf("  ✓ Migrated %d orders\n", count)
}

func migrateOrderItems(old, new *sql.DB) {
	fmt.Println("  → Migrating order items...")
	rows, err := old.Query("SELECT id, order_id, product_variant_id, product_name, product_sku, size, color, quantity, unit_price, total_price, created_at FROM order_items")
	if err != nil {
		log.Println("Warning: Could not read order items:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, orderID, quantity int
		var productVariantID sql.NullInt64
		var productName, productSku, size, color string
		var unitPrice, totalPrice float64
		var createdAt sql.NullTime

		rows.Scan(&id, &orderID, &productVariantID, &productName, &productSku, &size, &color, &quantity, &unitPrice, &totalPrice, &createdAt)

		_, err = new.Exec(
			"INSERT INTO order_items (id, order_id, product_variant_id, product_name, product_sku, size, color, quantity, unit_price, total_price, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
			id, orderID, productVariantID, productName, productSku, size, color, quantity, unitPrice, totalPrice, createdAt,
		)
		if err == nil {
			count++
		}
	}
	fmt.Printf("  ✓ Migrated %d order items\n", count)
}

func migrateCart(old, new *sql.DB) {
	fmt.Println("  → Migrating cart...")
	rows, err := old.Query("SELECT id, user_id, product_variant_id, quantity, created_at, updated_at FROM cart")
	if err != nil {
		log.Println("Warning: Could not read cart:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, userID, productVariantID, quantity int
		var createdAt, updatedAt sql.NullTime

		rows.Scan(&id, &userID, &productVariantID, &quantity, &createdAt, &updatedAt)

		_, err = new.Exec(
			"INSERT INTO cart (id, user_id, product_variant_id, quantity, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
			id, userID, productVariantID, quantity, createdAt, updatedAt,
		)
		if err == nil {
			count++
		}
	}
	fmt.Printf("  ✓ Migrated %d cart items\n", count)
}

func migrateWishlist(old, new *sql.DB) {
	fmt.Println("  → Migrating wishlist...")
	rows, err := old.Query("SELECT id, user_id, product_id, created_at FROM wishlist")
	if err != nil {
		log.Println("Warning: Could not read wishlist:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, userID, productID int
		var createdAt sql.NullTime

		rows.Scan(&id, &userID, &productID, &createdAt)

		_, err = new.Exec(
			"INSERT INTO wishlist (id, user_id, product_id, created_at) VALUES ($1, $2, $3, $4)",
			id, userID, productID, createdAt,
		)
		if err == nil {
			count++
		}
	}
	fmt.Printf("  ✓ Migrated %d wishlist items\n", count)
}