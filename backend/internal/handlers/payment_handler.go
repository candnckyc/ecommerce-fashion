package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"ecommerce-backend/internal/utils"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

type PaymentHandler struct {
	db *sql.DB
}

func NewPaymentHandler(db *sql.DB) *PaymentHandler {
	// Set Stripe API key from environment
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	
	return &PaymentHandler{db: db}
}

// CreatePaymentIntentRequest is the request body for creating a payment intent
type CreatePaymentIntentRequest struct {
	Amount   int64  `json:"amount"`    // Amount in cents (e.g., 10000 = $100.00)
	Currency string `json:"currency"`  // e.g., "usd", "try"
	OrderID  int    `json:"order_id"`  // Optional: link to order
}

// CreatePaymentIntent creates a Stripe payment intent
func (h *PaymentHandler) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentIntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate amount
	if req.Amount <= 0 {
		utils.Error(w, http.StatusBadRequest, "Amount must be greater than 0")
		return
	}

	// Default currency to USD if not provided
	if req.Currency == "" {
		req.Currency = "usd"
	}

	// Create payment intent parameters
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(req.Currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	// Add metadata if order ID is provided
	if req.OrderID > 0 {
		params.Metadata = map[string]string{
			"order_id": string(rune(req.OrderID)),
		}
	}

	// Create the payment intent
	pi, err := paymentintent.New(params)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create payment intent: "+err.Error())
		return
	}

	// Return the client secret
	utils.Success(w, map[string]string{
		"client_secret": pi.ClientSecret,
		"payment_intent_id": pi.ID,
	})
}

// ConfirmPaymentRequest is the request body for confirming a payment
type ConfirmPaymentRequest struct {
	PaymentIntentID string `json:"payment_intent_id"`
	OrderID         int    `json:"order_id"`
}

// ConfirmPayment updates order status after successful payment
func (h *PaymentHandler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	var req ConfirmPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Retrieve payment intent from Stripe to verify status
	pi, err := paymentintent.Get(req.PaymentIntentID, nil)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve payment intent")
		return
	}

	// Check if payment was successful
	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		utils.Error(w, http.StatusBadRequest, "Payment not completed")
		return
	}

	// Update order status to 'confirmed' and add payment info
	query := `
		UPDATE orders 
		SET status = 'confirmed',
		    payment_status = 'paid',
		    payment_method = 'stripe',
		    payment_transaction_id = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err = h.db.Exec(query, req.PaymentIntentID, req.OrderID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update order")
		return
	}

	utils.Success(w, map[string]interface{}{
		"message": "Payment confirmed successfully",
		"order_id": req.OrderID,
		"payment_status": "paid",
	})
}

// GetPaymentStatus retrieves the status of a payment intent
func (h *PaymentHandler) GetPaymentStatus(w http.ResponseWriter, r *http.Request) {
	paymentIntentID := r.URL.Query().Get("payment_intent_id")
	
	if paymentIntentID == "" {
		utils.Error(w, http.StatusBadRequest, "Missing payment_intent_id")
		return
	}

	// Retrieve payment intent from Stripe
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve payment intent")
		return
	}

	utils.Success(w, map[string]interface{}{
		"payment_intent_id": pi.ID,
		"status": pi.Status,
		"amount": pi.Amount,
		"currency": pi.Currency,
	})
}