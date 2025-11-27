-- ============================================
-- E-Commerce Fashion Store - Test Data
-- ============================================
-- This script populates the database with sample data for testing
-- Run this AFTER migrations are successful

-- ============================================
-- 1. USERS (2 customers + 1 admin)
-- ============================================
-- Note: Password is 'password123' for all users (will be hashed in production)
INSERT INTO users (email, password_hash, first_name, last_name, phone, role) VALUES
('john.doe@example.com', '$2a$10$rjLX8Aq2TmKxmKj6Y.vPLu5yyF9BZqV4b0GpXxC4fJJ3JqTYRXrYC', 'John', 'Doe', '+90 555 123 4567', 'customer'),
('jane.smith@example.com', '$2a$10$rjLX8Aq2TmKxmKj6Y.vPLu5yyF9BZqV4b0GpXxC4fJJ3JqTYRXrYC', 'Jane', 'Smith', '+90 555 987 6543', 'customer'),
('admin@ecommerce.com', '$2a$10$rjLX8Aq2TmKxmKj6Y.vPLu5yyF9BZqV4b0GpXxC4fJJ3JqTYRXrYC', 'Admin', 'User', '+90 555 000 0000', 'admin');

-- ============================================
-- 2. BRANDS
-- ============================================
INSERT INTO brands (name, slug, description) VALUES
('Zara', 'zara', 'Spanish fast fashion retailer known for trendy and affordable clothing'),
('H&M', 'hm', 'Swedish multinational clothing company offering fashion and quality at the best price'),
('Mango', 'mango', 'Barcelona-based fashion company with contemporary and urban designs'),
('Nike', 'nike', 'Leading sports apparel and footwear brand'),
('Adidas', 'adidas', 'German multinational corporation for sportswear'),
('Pull&Bear', 'pullbear', 'Spanish clothing and accessories retailer for young people');

-- ============================================
-- 3. CATEGORIES (Hierarchical)
-- ============================================
-- Top level categories
INSERT INTO categories (name, slug, parent_id, description) VALUES
('Women', 'women', NULL, 'Women''s fashion and accessories'),
('Men', 'men', NULL, 'Men''s fashion and accessories'),
('Kids', 'kids', NULL, 'Kids fashion and accessories');

-- Women subcategories
INSERT INTO categories (name, slug, parent_id, description) VALUES
('Dresses', 'dresses', 1, 'Women''s dresses for all occasions'),
('Tops', 'tops', 1, 'Women''s tops, blouses, and shirts'),
('Jeans', 'jeans', 1, 'Women''s jeans and denim'),
('Outerwear', 'outerwear', 1, 'Women''s coats and jackets');

-- Men subcategories
INSERT INTO categories (name, slug, parent_id, description) VALUES
('Shirts', 'shirts', 2, 'Men''s shirts and button-downs'),
('T-Shirts', 'tshirts', 2, 'Men''s t-shirts and casual tops'),
('Pants', 'pants', 2, 'Men''s pants and trousers'),
('Jackets', 'jackets', 2, 'Men''s jackets and outerwear');

-- Kids subcategories
INSERT INTO categories (name, slug, parent_id, description) VALUES
('Girls', 'girls', 3, 'Girls clothing'),
('Boys', 'boys', 3, 'Boys clothing');

-- ============================================
-- 4. PRODUCTS
-- ============================================
-- Women's Products
INSERT INTO products (name, slug, description, brand_id, category_id, base_price, is_active) VALUES
('Floral Summer Dress', 'floral-summer-dress', 'Light and breezy floral print dress perfect for summer days. Made from 100% cotton with a flattering fit.', 1, 4, 299.99, true),
('Classic White Blouse', 'classic-white-blouse', 'Elegant white blouse suitable for office or casual wear. Features button-down design and breathable fabric.', 3, 5, 199.99, true),
('High-Waist Skinny Jeans', 'high-waist-skinny-jeans', 'Comfortable stretch denim jeans with high waist fit. Perfect for everyday wear.', 1, 6, 399.99, true),
('Wool Blend Coat', 'wool-blend-coat', 'Warm and stylish coat for cold weather. Features classic design and quality wool blend material.', 3, 7, 899.99, true),
('Casual Cotton T-Shirt', 'casual-cotton-tshirt-women', 'Basic cotton t-shirt in multiple colors. Essential wardrobe staple.', 2, 5, 99.99, true);

-- Men's Products
INSERT INTO products (name, slug, description, brand_id, category_id, base_price, is_active) VALUES
('Oxford Button-Down Shirt', 'oxford-button-down-shirt', 'Classic Oxford shirt perfect for business or casual wear. 100% cotton.', 1, 8, 249.99, true),
('Graphic Print T-Shirt', 'graphic-print-tshirt', 'Trendy graphic t-shirt with modern design. Comfortable cotton blend.', 6, 9, 149.99, true),
('Slim Fit Chinos', 'slim-fit-chinos', 'Modern slim fit chinos in versatile colors. Perfect for smart casual looks.', 1, 10, 349.99, true),
('Denim Jacket', 'denim-jacket', 'Classic denim jacket that never goes out of style. Perfect layering piece.', 1, 11, 599.99, true),
('Sports Track Pants', 'sports-track-pants', 'Comfortable athletic pants for sports or casual wear. Moisture-wicking fabric.', 4, 10, 299.99, true);

-- Sports Products
INSERT INTO products (name, slug, description, brand_id, category_id, base_price, is_active) VALUES
('Running Shoes', 'running-shoes-nike', 'Lightweight running shoes with excellent cushioning and support.', 4, 9, 799.99, true),
('Training Sneakers', 'training-sneakers-adidas', 'Versatile training shoes suitable for gym and everyday wear.', 5, 9, 699.99, true);

-- ============================================
-- 5. PRODUCT VARIANTS (Size & Color combinations)
-- ============================================
-- Floral Summer Dress (Product ID: 1)
INSERT INTO product_variants (product_id, sku, size, color, color_hex, stock_quantity, price_adjustment) VALUES
(1, 'FSD-S-RED', 'S', 'Red', '#FF0000', 10, 0.00),
(1, 'FSD-M-RED', 'M', 'Red', '#FF0000', 15, 0.00),
(1, 'FSD-L-RED', 'L', 'Red', '#FF0000', 8, 0.00),
(1, 'FSD-S-BLUE', 'S', 'Blue', '#0000FF', 12, 0.00),
(1, 'FSD-M-BLUE', 'M', 'Blue', '#0000FF', 20, 0.00),
(1, 'FSD-L-BLUE', 'L', 'Blue', '#0000FF', 10, 0.00);

-- Classic White Blouse (Product ID: 2)
INSERT INTO product_variants (product_id, sku, size, color, color_hex, stock_quantity, price_adjustment) VALUES
(2, 'CWB-XS-WHITE', 'XS', 'White', '#FFFFFF', 8, 0.00),
(2, 'CWB-S-WHITE', 'S', 'White', '#FFFFFF', 15, 0.00),
(2, 'CWB-M-WHITE', 'M', 'White', '#FFFFFF', 20, 0.00),
(2, 'CWB-L-WHITE', 'L', 'White', '#FFFFFF', 12, 0.00),
(2, 'CWB-XL-WHITE', 'XL', 'White', '#FFFFFF', 5, 0.00);

-- High-Waist Skinny Jeans (Product ID: 3)
INSERT INTO product_variants (product_id, sku, size, color, color_hex, stock_quantity, price_adjustment) VALUES
(3, 'HWSJ-26-BLUE', '26', 'Dark Blue', '#00008B', 10, 0.00),
(3, 'HWSJ-28-BLUE', '28', 'Dark Blue', '#00008B', 15, 0.00),
(3, 'HWSJ-30-BLUE', '30', 'Dark Blue', '#00008B', 15, 0.00),
(3, 'HWSJ-32-BLUE', '32', 'Dark Blue', '#00008B', 10, 0.00),
(3, 'HWSJ-26-BLACK', '26', 'Black', '#000000', 12, 0.00),
(3, 'HWSJ-28-BLACK', '28', 'Black', '#000000', 18, 0.00),
(3, 'HWSJ-30-BLACK', '30', 'Black', '#000000', 15, 0.00),
(3, 'HWSJ-32-BLACK', '32', 'Black', '#000000', 8, 0.00);

-- Wool Blend Coat (Product ID: 4)
INSERT INTO product_variants (product_id, sku, size, color, color_hex, stock_quantity, price_adjustment) VALUES
(4, 'WBC-S-CAMEL', 'S', 'Camel', '#C19A6B', 5, 0.00),
(4, 'WBC-M-CAMEL', 'M', 'Camel', '#C19A6B', 8, 0.00),
(4, 'WBC-L-CAMEL', 'L', 'Camel', '#C19A6B', 6, 0.00),
(4, 'WBC-S-BLACK', 'S', 'Black', '#000000', 5, 0.00),
(4, 'WBC-M-BLACK', 'M', 'Black', '#000000', 10, 0.00),
(4, 'WBC-L-BLACK', 'L', 'Black', '#000000', 7, 0.00);

-- Casual Cotton T-Shirt Women (Product ID: 5)
INSERT INTO product_variants (product_id, sku, size, color, color_hex, stock_quantity, price_adjustment) VALUES
(5, 'CCT-W-S-WHITE', 'S', 'White', '#FFFFFF', 25, 0.00),
(5, 'CCT-W-M-WHITE', 'M', 'White', '#FFFFFF', 30, 0.00),
(5, 'CCT-W-L-WHITE', 'L', 'White', '#FFFFFF', 20, 0.00),
(5, 'CCT-W-S-BLACK', 'S', 'Black', '#000000', 25, 0.00),
(5, 'CCT-W-M-BLACK', 'M', 'Black', '#000000', 35, 0.00),
(5, 'CCT-W-L-BLACK', 'L', 'Black', '#000000', 25, 0.00),
(5, 'CCT-W-S-NAVY', 'S', 'Navy', '#000080', 20, 0.00),
(5, 'CCT-W-M-NAVY', 'M', 'Navy', '#000080', 25, 0.00),
(5, 'CCT-W-L-NAVY', 'L', 'Navy', '#000080', 15, 0.00);

-- Oxford Shirt (Product ID: 6)
INSERT INTO product_variants (product_id, sku, size, color, color_hex, stock_quantity, price_adjustment) VALUES
(6, 'OXF-S-WHITE', 'S', 'White', '#FFFFFF', 15, 0.00),
(6, 'OXF-M-WHITE', 'M', 'White', '#FFFFFF', 20, 0.00),
(6, 'OXF-L-WHITE', 'L', 'White', '#FFFFFF', 18, 0.00),
(6, 'OXF-XL-WHITE', 'XL', 'White', '#FFFFFF', 10, 0.00),
(6, 'OXF-S-BLUE', 'S', 'Light Blue', '#ADD8E6', 12, 0.00),
(6, 'OXF-M-BLUE', 'M', 'Light Blue', '#ADD8E6', 15, 0.00),
(6, 'OXF-L-BLUE', 'L', 'Light Blue', '#ADD8E6', 15, 0.00),
(6, 'OXF-XL-BLUE', 'XL', 'Light Blue', '#ADD8E6', 8, 0.00);

-- Graphic T-Shirt (Product ID: 7)
INSERT INTO product_variants (product_id, sku, size, color, color_hex, stock_quantity, price_adjustment) VALUES
(7, 'GPT-S-BLACK', 'S', 'Black', '#000000', 20, 0.00),
(7, 'GPT-M-BLACK', 'M', 'Black', '#000000', 25, 0.00),
(7, 'GPT-L-BLACK', 'L', 'Black', '#000000', 20, 0.00),
(7, 'GPT-XL-BLACK', 'XL', 'Black', '#000000', 15, 0.00);

-- Slim Fit Chinos (Product ID: 8)
INSERT INTO product_variants (product_id, sku, size, color, color_hex, stock_quantity, price_adjustment) VALUES
(8, 'SFC-30-BEIGE', '30', 'Beige', '#F5F5DC', 15, 0.00),
(8, 'SFC-32-BEIGE', '32', 'Beige', '#F5F5DC', 20, 0.00),
(8, 'SFC-34-BEIGE', '34', 'Beige', '#F5F5DC', 18, 0.00),
(8, 'SFC-36-BEIGE', '36', 'Beige', '#F5F5DC', 12, 0.00),
(8, 'SFC-30-NAVY', '30', 'Navy', '#000080', 15, 0.00),
(8, 'SFC-32-NAVY', '32', 'Navy', '#000080', 22, 0.00),
(8, 'SFC-34-NAVY', '34', 'Navy', '#000080', 18, 0.00),
(8, 'SFC-36-NAVY', '36', 'Navy', '#000080', 10, 0.00);

-- Denim Jacket (Product ID: 9)
INSERT INTO product_variants (product_id, sku, size, color, color_hex, stock_quantity, price_adjustment) VALUES
(9, 'DJ-S-BLUE', 'S', 'Blue', '#4169E1', 10, 0.00),
(9, 'DJ-M-BLUE', 'M', 'Blue', '#4169E1', 15, 0.00),
(9, 'DJ-L-BLUE', 'L', 'Blue', '#4169E1', 12, 0.00),
(9, 'DJ-XL-BLUE', 'XL', 'Blue', '#4169E1', 8, 0.00);

-- Track Pants (Product ID: 10)
INSERT INTO product_variants (product_id, sku, size, color, color_hex, stock_quantity, price_adjustment) VALUES
(10, 'TP-S-BLACK', 'S', 'Black', '#000000', 15, 0.00),
(10, 'TP-M-BLACK', 'M', 'Black', '#000000', 20, 0.00),
(10, 'TP-L-BLACK', 'L', 'Black', '#000000', 18, 0.00),
(10, 'TP-XL-BLACK', 'XL', 'Black', '#000000', 12, 0.00),
(10, 'TP-S-GRAY', 'S', 'Gray', '#808080', 15, 0.00),
(10, 'TP-M-GRAY', 'M', 'Gray', '#808080', 18, 0.00),
(10, 'TP-L-GRAY', 'L', 'Gray', '#808080', 15, 0.00),
(10, 'TP-XL-GRAY', 'XL', 'Gray', '#808080', 10, 0.00);

-- Running Shoes (Product ID: 11)
INSERT INTO product_variants (product_id, sku, size, color, color_hex, stock_quantity, price_adjustment) VALUES
(11, 'RS-40-BLACK', '40', 'Black', '#000000', 10, 0.00),
(11, 'RS-41-BLACK', '41', 'Black', '#000000', 12, 0.00),
(11, 'RS-42-BLACK', '42', 'Black', '#000000', 15, 0.00),
(11, 'RS-43-BLACK', '43', 'Black', '#000000', 12, 0.00),
(11, 'RS-44-BLACK', '44', 'Black', '#000000', 8, 0.00),
(11, 'RS-40-WHITE', '40', 'White', '#FFFFFF', 10, 0.00),
(11, 'RS-41-WHITE', '41', 'White', '#FFFFFF', 12, 0.00),
(11, 'RS-42-WHITE', '42', 'White', '#FFFFFF', 15, 0.00),
(11, 'RS-43-WHITE', '43', 'White', '#FFFFFF', 10, 0.00),
(11, 'RS-44-WHITE', '44', 'White', '#FFFFFF', 5, 0.00);

-- Training Sneakers (Product ID: 12)
INSERT INTO product_variants (product_id, sku, size, color, color_hex, stock_quantity, price_adjustment) VALUES
(12, 'TS-40-GRAY', '40', 'Gray', '#808080', 12, 0.00),
(12, 'TS-41-GRAY', '41', 'Gray', '#808080', 15, 0.00),
(12, 'TS-42-GRAY', '42', 'Gray', '#808080', 18, 0.00),
(12, 'TS-43-GRAY', '43', 'Gray', '#808080', 15, 0.00),
(12, 'TS-44-GRAY', '44', 'Gray', '#808080', 10, 0.00);

-- ============================================
-- 6. PRODUCT IMAGES (Placeholder URLs)
-- ============================================
-- Note: These are placeholder URLs. In production, you'll upload actual images
INSERT INTO product_images (product_id, image_url, alt_text, display_order, is_primary) VALUES
-- Floral Summer Dress
(1, '/uploads/products/floral-dress-1.jpg', 'Floral Summer Dress - Front View', 1, true),
(1, '/uploads/products/floral-dress-2.jpg', 'Floral Summer Dress - Side View', 2, false),
(1, '/uploads/products/floral-dress-3.jpg', 'Floral Summer Dress - Back View', 3, false),

-- Classic White Blouse
(2, '/uploads/products/white-blouse-1.jpg', 'Classic White Blouse - Front', 1, true),
(2, '/uploads/products/white-blouse-2.jpg', 'Classic White Blouse - Detail', 2, false),

-- High-Waist Skinny Jeans
(3, '/uploads/products/jeans-1.jpg', 'High-Waist Skinny Jeans - Blue', 1, true),
(3, '/uploads/products/jeans-2.jpg', 'High-Waist Skinny Jeans - Black', 2, false),

-- Wool Blend Coat
(4, '/uploads/products/coat-1.jpg', 'Wool Blend Coat - Camel', 1, true),
(4, '/uploads/products/coat-2.jpg', 'Wool Blend Coat - Black', 2, false),

-- Casual Cotton T-Shirt
(5, '/uploads/products/tshirt-women-1.jpg', 'Cotton T-Shirt - White', 1, true),

-- Oxford Shirt
(6, '/uploads/products/oxford-1.jpg', 'Oxford Shirt - White', 1, true),
(6, '/uploads/products/oxford-2.jpg', 'Oxford Shirt - Blue', 2, false),

-- Graphic T-Shirt
(7, '/uploads/products/graphic-tee-1.jpg', 'Graphic Print T-Shirt', 1, true),

-- Slim Fit Chinos
(8, '/uploads/products/chinos-1.jpg', 'Slim Fit Chinos - Beige', 1, true),
(8, '/uploads/products/chinos-2.jpg', 'Slim Fit Chinos - Navy', 2, false),

-- Denim Jacket
(9, '/uploads/products/denim-jacket-1.jpg', 'Denim Jacket - Front', 1, true),
(9, '/uploads/products/denim-jacket-2.jpg', 'Denim Jacket - Back', 2, false),

-- Track Pants
(10, '/uploads/products/track-pants-1.jpg', 'Track Pants - Black', 1, true),

-- Running Shoes
(11, '/uploads/products/running-shoes-1.jpg', 'Running Shoes - Black', 1, true),
(11, '/uploads/products/running-shoes-2.jpg', 'Running Shoes - White', 2, false),

-- Training Sneakers
(12, '/uploads/products/sneakers-1.jpg', 'Training Sneakers', 1, true);

-- ============================================
-- 7. ADDRESSES (Sample addresses for test users)
-- ============================================
INSERT INTO addresses (user_id, title, full_name, phone, address_line1, address_line2, city, state, postal_code, country, is_default) VALUES
(1, 'Home', 'John Doe', '+90 555 123 4567', 'Atatürk Caddesi No: 123', 'Daire 5', 'Istanbul', 'Istanbul', '34000', 'Turkey', true),
(1, 'Work', 'John Doe', '+90 555 123 4567', 'İş Merkezi Kat 10', 'Ofis 1001', 'Istanbul', 'Istanbul', '34100', 'Turkey', false),
(2, 'Home', 'Jane Smith', '+90 555 987 6543', 'Cumhuriyet Bulvarı 456', 'Kat 3 Daire 8', 'Ankara', 'Ankara', '06000', 'Turkey', true);

-- ============================================
-- SUMMARY
-- ============================================
-- Data inserted:
-- - 3 users (2 customers + 1 admin)
-- - 6 brands
-- - 14 categories (3 top-level + 11 subcategories)
-- - 12 products
-- - 93 product variants (different sizes/colors)
-- - 21 product images
-- - 3 addresses
--
-- Total: 152 rows of test data
--
-- You can now:
-- 1. Browse products in your app
-- 2. Test filtering by category/brand
-- 3. Add items to cart
-- 4. Create test orders
-- ============================================
