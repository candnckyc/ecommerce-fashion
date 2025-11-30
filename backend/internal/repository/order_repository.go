package repository

import (
	"database/sql"
	"ecommerce-backend/internal/models"
	"fmt"
	"time"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder creates a new order
func (r *OrderRepository) CreateOrder(order *models.Order) (int, error) {
	query := `
		INSERT INTO orders (
			user_id, order_number, 
			shipping_address_line1, shipping_address_line2, shipping_city, 
			shipping_state, shipping_postal_code, shipping_country,
			shipping_full_name, shipping_phone,
			subtotal, shipping_cost, tax, total,
			status, payment_method, payment_status, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id, created_at, updated_at
	`
	
	var id int
	var createdAt, updatedAt time.Time
	
	err := r.db.QueryRow(
		query,
		order.UserID, order.OrderNumber,
		order.ShippingAddressLine1, order.ShippingAddressLine2, order.ShippingCity,
		order.ShippingState, order.ShippingPostalCode, order.ShippingCountry,
		order.ShippingFullName, order.ShippingPhone,
		order.Subtotal, order.ShippingCost, order.Tax, order.Total,
		order.Status, order.PaymentMethod, order.PaymentStatus, order.Notes,
	).Scan(&id, &createdAt, &updatedAt)
	
	if err != nil {
		return 0, err
	}
	
	order.ID = id
	order.CreatedAt = createdAt
	order.UpdatedAt = updatedAt
	
	return id, nil
}

// CreateOrderItem creates an order item
func (r *OrderRepository) CreateOrderItem(item *models.OrderItem) error {
	query := `
		INSERT INTO order_items (
			order_id, product_variant_id,
			product_name, product_sku, size, color,
			quantity, unit_price, total_price
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`
	
	return r.db.QueryRow(
		query,
		item.OrderID, item.ProductVariantID,
		item.ProductName, item.ProductSKU, item.Size, item.Color,
		item.Quantity, item.UnitPrice, item.TotalPrice,
	).Scan(&item.ID, &item.CreatedAt)
}

// GetUserOrders retrieves all orders for a user
func (r *OrderRepository) GetUserOrders(userID int) ([]models.Order, error) {
	query := `
		SELECT id, user_id, order_number, 
		       shipping_address_line1, shipping_address_line2, shipping_city,
		       shipping_state, shipping_postal_code, shipping_country,
		       shipping_full_name, shipping_phone,
		       subtotal, shipping_cost, tax, total,
		       status, payment_method, payment_status, notes,
		       created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		err := rows.Scan(
			&o.ID, &o.UserID, &o.OrderNumber,
			&o.ShippingAddressLine1, &o.ShippingAddressLine2, &o.ShippingCity,
			&o.ShippingState, &o.ShippingPostalCode, &o.ShippingCountry,
			&o.ShippingFullName, &o.ShippingPhone,
			&o.Subtotal, &o.ShippingCost, &o.Tax, &o.Total,
			&o.Status, &o.PaymentMethod, &o.PaymentStatus, &o.Notes,
			&o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

// GetAllOrders retrieves all orders (admin only)
func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
	query := `
		SELECT id, user_id, order_number, 
		       shipping_address_line1, shipping_address_line2, shipping_city,
		       shipping_state, shipping_postal_code, shipping_country,
		       shipping_full_name, shipping_phone,
		       subtotal, shipping_cost, tax, total,
		       status, payment_method, payment_status, notes,
		       created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		err := rows.Scan(
			&o.ID, &o.UserID, &o.OrderNumber,
			&o.ShippingAddressLine1, &o.ShippingAddressLine2, &o.ShippingCity,
			&o.ShippingState, &o.ShippingPostalCode, &o.ShippingCountry,
			&o.ShippingFullName, &o.ShippingPhone,
			&o.Subtotal, &o.ShippingCost, &o.Tax, &o.Total,
			&o.Status, &o.PaymentMethod, &o.PaymentStatus, &o.Notes,
			&o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

// GetOrderByID retrieves a single order with items
func (r *OrderRepository) GetOrderByID(orderID, userID int) (*models.Order, error) {
	query := `
		SELECT id, user_id, order_number,
		       shipping_address_line1, shipping_address_line2, shipping_city,
		       shipping_state, shipping_postal_code, shipping_country,
		       shipping_full_name, shipping_phone,
		       subtotal, shipping_cost, tax, total,
		       status, payment_method, payment_status, notes,
		       created_at, updated_at
		FROM orders
		WHERE id = $1 AND user_id = $2
	`
	
	order := &models.Order{}
	err := r.db.QueryRow(query, orderID, userID).Scan(
		&order.ID, &order.UserID, &order.OrderNumber,
		&order.ShippingAddressLine1, &order.ShippingAddressLine2, &order.ShippingCity,
		&order.ShippingState, &order.ShippingPostalCode, &order.ShippingCountry,
		&order.ShippingFullName, &order.ShippingPhone,
		&order.Subtotal, &order.ShippingCost, &order.Tax, &order.Total,
		&order.Status, &order.PaymentMethod, &order.PaymentStatus, &order.Notes,
		&order.CreatedAt, &order.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	// Get order items
	items, err := r.getOrderItems(orderID)
	if err == nil {
		order.Items = items
	}
	
	return order, nil
}

// getOrderItems retrieves items for an order
func (r *OrderRepository) getOrderItems(orderID int) ([]models.OrderItem, error) {
	query := `
		SELECT id, order_id, product_variant_id,
		       product_name, product_sku, size, color,
		       quantity, unit_price, total_price, created_at
		FROM order_items
		WHERE order_id = $1
	`
	
	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	items := []models.OrderItem{}
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductVariantID,
			&item.ProductName, &item.ProductSKU, &item.Size, &item.Color,
			&item.Quantity, &item.UnitPrice, &item.TotalPrice, &item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	
	return items, nil
}

// GenerateOrderNumber generates a unique order number
func (r *OrderRepository) GenerateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano()/1000000)
}

// GetAddressByID retrieves an address
func (r *OrderRepository) GetAddressByID(addressID, userID int) (*models.Address, error) {
	query := `
		SELECT id, user_id, title, full_name, phone,
		       address_line1, address_line2, city, state, postal_code, country,
		       is_default, created_at, updated_at
		FROM addresses
		WHERE id = $1 AND user_id = $2
	`
	
	addr := &models.Address{}
	err := r.db.QueryRow(query, addressID, userID).Scan(
		&addr.ID, &addr.UserID, &addr.Title, &addr.FullName, &addr.Phone,
		&addr.AddressLine1, &addr.AddressLine2, &addr.City, &addr.State,
		&addr.PostalCode, &addr.Country, &addr.IsDefault,
		&addr.CreatedAt, &addr.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return addr, err
}

// GetUserAddresses retrieves all addresses for a user
func (r *OrderRepository) GetUserAddresses(userID int) ([]models.Address, error) {
	query := `
		SELECT id, user_id, title, full_name, phone,
		       address_line1, address_line2, city, state, postal_code, country,
		       is_default, created_at, updated_at
		FROM addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	addresses := []models.Address{}
	for rows.Next() {
		var addr models.Address
		err := rows.Scan(
			&addr.ID, &addr.UserID, &addr.Title, &addr.FullName, &addr.Phone,
			&addr.AddressLine1, &addr.AddressLine2, &addr.City, &addr.State,
			&addr.PostalCode, &addr.Country, &addr.IsDefault,
			&addr.CreatedAt, &addr.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}
	
	return addresses, nil
}

// CreateAddress creates a new address
func (r *OrderRepository) CreateAddress(addr *models.Address) (int, error) {
	query := `
		INSERT INTO addresses (
			user_id, title, full_name, phone,
			address_line1, address_line2, city, state, postal_code, country,
			is_default
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	
	var id int
	err := r.db.QueryRow(
		query,
		addr.UserID, addr.Title, addr.FullName, addr.Phone,
		addr.AddressLine1, addr.AddressLine2, addr.City, addr.State,
		addr.PostalCode, addr.Country, addr.IsDefault,
	).Scan(&id, &addr.CreatedAt, &addr.UpdatedAt)
	
	return id, err
}

// UpdateOrderStatus updates order status
func (r *OrderRepository) UpdateOrderStatus(orderID int, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, status, orderID)
	return err
}