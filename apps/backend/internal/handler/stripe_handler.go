package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/imbivek08/hamropasal/internal/middleware"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type StripeHandler struct {
	stripeService *service.StripeService
	orderService  *service.OrderService
	userService   *service.UserService
	webhookSecret string
}

func NewStripeHandler(
	stripeService *service.StripeService,
	orderService *service.OrderService,
	userService *service.UserService,
	webhookSecret string,
) *StripeHandler {
	return &StripeHandler{
		stripeService: stripeService,
		orderService:  orderService,
		userService:   userService,
		webhookSecret: webhookSecret,
	}
}

// CreateCheckoutSession creates an order and returns a Stripe Checkout URL.
// The frontend should redirect the user to this URL.
func (h *StripeHandler) CreateCheckoutSession(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Get internal user
	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	// Parse request — same shape as regular CreateOrderRequest
	var req model.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "validation failed")
	}

	// Force payment method to stripe
	req.PaymentMethod = "stripe"

	// Create the order (it will be in "pending" status)
	orderResp, err := h.orderService.CreateOrderFromCart(c.Request().Context(), user.ID, &req)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, err.Error())
	}

	// Create Stripe checkout session
	order := &model.Order{
		ID:          orderResp.ID,
		OrderNumber: orderResp.OrderNumber,
		Total:       orderResp.Total,
	}

	checkoutURL, err := h.stripeService.CreateCheckoutSession(c.Request().Context(), order, orderResp.Items)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to create payment session")
	}

	return SendSuccess(c, http.StatusCreated, "checkout session created", model.CreateCheckoutResponse{
		Order:       orderResp,
		CheckoutURL: checkoutURL,
	})
}

// HandleStripeWebhook processes Stripe webhook events.
// This endpoint must NOT have auth middleware — Stripe calls it directly.
func (h *StripeHandler) HandleStripeWebhook(c echo.Context) error {
	// Read raw body from context (saved by saveRawBody middleware)
	var body []byte
	rawBody := c.Get("raw_body")
	if rawBody != nil {
		body = rawBody.([]byte)
	} else {
		var err error
		body, err = io.ReadAll(c.Request().Body)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err, "failed to read request body")
		}
	}

	// Debug logging
	sigHeader := c.Request().Header.Get("Stripe-Signature")
	fmt.Printf("[Stripe Webhook] Body length: %d, Sig header present: %v, Webhook secret length: %d\n",
		len(body), sigHeader != "", len(h.webhookSecret))

	// Verify webhook signature
	event, err := webhook.ConstructEvent(body, sigHeader, h.webhookSecret)
	if err != nil {
		fmt.Printf("[Stripe Webhook] Signature verification FAILED: %v\n", err)
		return SendError(c, http.StatusBadRequest, err, "invalid webhook signature")
	}

	fmt.Printf("[Stripe Webhook] Event received: %s\n", event.Type)

	ctx := c.Request().Context()

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return SendError(c, http.StatusBadRequest, err, "invalid session data")
		}

		if err := h.stripeService.HandlePaymentSuccess(ctx, session.ID); err != nil {
			return SendError(c, http.StatusInternalServerError, err, "failed to process payment success")
		}

	case "checkout.session.expired":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return SendError(c, http.StatusBadRequest, err, "invalid session data")
		}

		if err := h.stripeService.HandlePaymentFailure(ctx, session.ID); err != nil {
			return SendError(c, http.StatusInternalServerError, err, "failed to process payment failure")
		}

	default:
		// Unhandled event type, just acknowledge
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// VerifySession lets the frontend check the payment status after redirect.
func (h *StripeHandler) VerifySession(c echo.Context) error {
	sessionID := c.QueryParam("session_id")
	if sessionID == "" {
		return SendError(c, http.StatusBadRequest, nil, "session_id is required")
	}

	status, err := h.stripeService.VerifySession(c.Request().Context(), sessionID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to verify session")
	}

	return SendSuccess(c, http.StatusOK, "session verified", status)
}
