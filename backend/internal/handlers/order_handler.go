package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/services"
	"ecommerce-backend/internal/utils"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder creates a new order
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	order, err := h.orderService.CreateOrder(userID, &req)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(w, order)
}

// GetOrders retrieves user's orders
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	orders, err := h.orderService.GetUserOrders(userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}

	utils.Success(w, orders)
}

// GetOrder retrieves a single order
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.orderService.GetOrderByID(orderID, userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}
	if order == nil {
		utils.Error(w, http.StatusNotFound, "Order not found")
		return
	}

	utils.Success(w, order)
}

// GetAddresses retrieves user's addresses
func (h *OrderHandler) GetAddresses(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	addresses, err := h.orderService.GetUserAddresses(userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve addresses")
		return
	}

	utils.Success(w, addresses)
}

// CreateAddress creates a new address
func (h *OrderHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.AddAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	address, err := h.orderService.CreateAddress(userID, &req)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(w, address)
}