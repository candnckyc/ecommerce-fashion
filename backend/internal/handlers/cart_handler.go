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

type CartHandler struct {
	cartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// GetCart returns user's cart
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to get cart")
		return
	}

	utils.Success(w, cart)
}

// AddToCart adds item to cart
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.cartService.AddToCart(userID, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return updated cart
	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to get cart")
		return
	}

	utils.Success(w, cart)
}

// UpdateCartItem updates cart item quantity
func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	cartItemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid cart item ID")
		return
	}

	var req models.UpdateCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.cartService.UpdateCartItem(cartItemID, userID, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return updated cart
	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to get cart")
		return
	}

	utils.Success(w, cart)
}

// RemoveFromCart removes item from cart
func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	cartItemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid cart item ID")
		return
	}

	if err := h.cartService.RemoveFromCart(cartItemID, userID); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to remove item")
		return
	}

	// Return updated cart
	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to get cart")
		return
	}

	utils.Success(w, cart)
}

// ClearCart removes all items from cart
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.cartService.ClearCart(userID); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to clear cart")
		return
	}

	utils.Success(w, map[string]string{"message": "Cart cleared"})
}
