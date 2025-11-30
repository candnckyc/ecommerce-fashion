package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"ecommerce-backend/internal/repository"
	//"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/services"
	"ecommerce-backend/internal/utils"
)

type AdminHandler struct {
	productService *services.ProductService
	orderService   *services.OrderService
	userRepo       *repository.UserRepository
}

func NewAdminHandler(productService *services.ProductService, orderService *services.OrderService, userRepo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{
		productService: productService,
		orderService:   orderService,
		userRepo:       userRepo,
	}
}

// GetAllOrders retrieves all orders (admin only)
func (h *AdminHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}
	utils.Success(w, orders)
}

// UpdateOrderStatus updates order status (admin only)
func (h *AdminHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.orderService.UpdateOrderStatus(orderID, req.Status); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	utils.Success(w, map[string]interface{}{
		"message":  "Order status updated",
		"order_id": orderID,
		"status":   req.Status,
	})
}

// GetStats returns basic admin statistics
func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement real stats from database
	stats := map[string]interface{}{
		"total_products": 100,
		"total_orders":   50,
		"total_revenue":  125000.00,
		"pending_orders": 12,
	}
	utils.Success(w, stats)
}

// ToggleProduct toggles product active status
func (h *AdminHandler) ToggleProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.productService.ToggleActive(productID); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to toggle product")
		return
	}

	utils.Success(w, map[string]interface{}{
		"message":    "Product status toggled",
		"product_id": productID,
	})
}

// GetAllCustomers retrieves all customers (admin only)
func (h *AdminHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to fetch customers")
		return
	}
	
	// Remove password hashes before sending
	for i := range users {
		users[i].PasswordHash = ""
	}
	
	utils.Success(w, users)
}