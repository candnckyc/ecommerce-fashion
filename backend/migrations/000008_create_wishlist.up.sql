-- Create wishlist table
CREATE TABLE wishlist (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure user can only add a product once to wishlist
    UNIQUE(user_id, product_id)
);

-- Create index for faster user wishlist queries
CREATE INDEX idx_wishlist_user ON wishlist(user_id);
