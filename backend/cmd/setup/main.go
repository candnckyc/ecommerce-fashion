package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// New database connection
	dbURL := "postgresql://idle_eticaret:kL8zHtqX-1iYArLM@51.178.9.8:25600/idle_eticaret?sslmode=disable"

	fmt.Println("🗄️  Fresh Database Setup")
	fmt.Println("Connecting to:", dbURL)
	fmt.Println()

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("✓ Connected to database successfully!")
	fmt.Println()
	fmt.Println("Creating tables...")

	// Table 1: Users
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			phone VARCHAR(20),
			role VARCHAR(20) DEFAULT 'customer',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}
	fmt.Println("✓ Users table created")

	// Table 2: Brands
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS brands (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL,
			slug VARCHAR(100) UNIQUE NOT NULL,
			description TEXT,
			logo_url VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Failed to create brands table:", err)
	}
	fmt.Println("✓ Brands table created")

	// Table 3: Categories
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			slug VARCHAR(100) UNIQUE NOT NULL,
			parent_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
			description TEXT,
			image_url VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Failed to create categories table:", err)
	}
	fmt.Println("✓ Categories table created")

	// Table 4: Products
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
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
		)
	`)
	if err != nil {
		log.Fatal("Failed to create products table:", err)
	}
	fmt.Println("✓ Products table created")

	// Create indexes for products
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_products_search ON products USING gin(search_vector)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand_id)`)
	fmt.Println("✓ Product indexes created")

	// Create search trigger
	db.Exec(`
		CREATE OR REPLACE FUNCTION products_search_trigger() RETURNS trigger AS $$
		BEGIN
			NEW.search_vector := to_tsvector('english', 
				coalesce(NEW.name, '') || ' ' || 
				coalesce(NEW.description, '')
			);
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql
	`)
	db.Exec(`DROP TRIGGER IF EXISTS tsvector_update ON products`)
	db.Exec(`
		CREATE TRIGGER tsvector_update BEFORE INSERT OR UPDATE
		ON products FOR EACH ROW EXECUTE FUNCTION products_search_trigger()
	`)
	fmt.Println("✓ Search trigger created")

	// Table 5: Product Variants
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS product_variants (
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
		)
	`)
	if err != nil {
		log.Fatal("Failed to create product_variants table:", err)
	}
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_variants_product ON product_variants(product_id)`)
	fmt.Println("✓ Product variants table created")

	// Table 6: Product Images
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS product_images (
			id SERIAL PRIMARY KEY,
			product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
			image_url VARCHAR(500) NOT NULL,
			alt_text VARCHAR(255),
			display_order INTEGER DEFAULT 0,
			is_primary BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Failed to create product_images table:", err)
	}
	fmt.Println("✓ Product images table created")

	// Table 7: Cart
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cart (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			product_variant_id INTEGER REFERENCES product_variants(id) ON DELETE CASCADE,
			quantity INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT check_cart_quantity CHECK (quantity > 0),
			UNIQUE(user_id, product_variant_id)
		)
	`)
	if err != nil {
		log.Fatal("Failed to create cart table:", err)
	}
	fmt.Println("✓ Cart table created")

	// Table 8: Wishlist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS wishlist (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, product_id)
		)
	`)
	if err != nil {
		log.Fatal("Failed to create wishlist table:", err)
	}
	fmt.Println("✓ Wishlist table created")

	// Table 9: Addresses
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS addresses (
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
		)
	`)
	if err != nil {
		log.Fatal("Failed to create addresses table:", err)
	}
	fmt.Println("✓ Addresses table created")

	// Table 10: Orders
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
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
		)
	`)
	if err != nil {
		log.Fatal("Failed to create orders table:", err)
	}
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)`)
	fmt.Println("✓ Orders table created")

	// Table 11: Order Items
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
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
		)
	`)
	if err != nil {
		log.Fatal("Failed to create order_items table:", err)
	}
	fmt.Println("✓ Order items table created")

	// Create admin user
	_, err = db.Exec(`
		INSERT INTO users (email, password_hash, first_name, last_name, role)
		VALUES ('admin@fashion.com', '$2a$10$rGvH8xqvZ8K4fLJqE7XLHOkXqXg8vT.Xc5XZQqZVzN8mZqVLQzWYK', 'Admin', 'User', 'admin')
		ON CONFLICT (email) DO NOTHING
	`)
	if err == nil {
		fmt.Println("✓ Admin user created (email: admin@fashion.com, password: admin123)")
	}

	fmt.Println()
	fmt.Println("✅ Database setup complete!")
	fmt.Println("You can now run: go run cmd/api/main.go")
}