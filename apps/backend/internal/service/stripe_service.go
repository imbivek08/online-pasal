package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/repository"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
)

type StripeService struct {
	orderRepo   *repository.OrderRepository
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
	addressRepo *repository.AddressRepository
	frontendURL string
}

func NewStripeService(
	apiKey string,
	frontendURL string,
	orderRepo *repository.OrderRepository,
	cartRepo *repository.CartRepository,
	productRepo *repository.ProductRepository,
	addressRepo *repository.AddressRepository,
) *StripeService {
	stripe.Key = apiKey
	return &StripeService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		addressRepo: addressRepo,
		frontendURL: frontendURL,
	}
}

// CreateCheckoutSession creates a Stripe Checkout Session for the given order.
// The order must already exist in the database with status "pending".
func (s *StripeService) CreateCheckoutSession(ctx context.Context, order *model.Order, items []model.OrderItemWithDetails) (string, error) {
	var lineItems []*stripe.CheckoutSessionLineItemParams
	for _, item := range items {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("npr"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.ProductName),
				},
				// Stripe expects amount in cents
				UnitAmount: stripe.Int64(int64(item.UnitPrice * 100)),
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(fmt.Sprintf("%s/payment/success?session_id={CHECKOUT_SESSION_ID}", s.frontendURL)),
		CancelURL:          stripe.String(fmt.Sprintf("%s/payment/cancel", s.frontendURL)),
		Metadata: map[string]string{
			"order_id":     order.ID.String(),
			"order_number": order.OrderNumber,
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create stripe checkout session: %w", err)
	}

	// Save the stripe session ID on the order
	if err := s.orderRepo.UpdateStripeSessionID(ctx, order.ID, sess.ID); err != nil {
		return "", fmt.Errorf("failed to save stripe session id: %w", err)
	}

	return sess.URL, nil
}

// HandlePaymentSuccess is called by the webhook when payment succeeds.
func (s *StripeService) HandlePaymentSuccess(ctx context.Context, sessionID string) error {
	// Retrieve the session to get the order_id from metadata
	sess, err := session.Get(sessionID, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve stripe session: %w", err)
	}

	orderIDStr, ok := sess.Metadata["order_id"]
	if !ok {
		return fmt.Errorf("order_id not found in session metadata")
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return fmt.Errorf("invalid order_id in metadata: %w", err)
	}

	// Update payment status to paid
	if err := s.orderRepo.UpdatePaymentStatus(ctx, orderID, model.PaymentStatusPaid); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// Update order status to confirmed
	if err := s.orderRepo.UpdateStatusWithTimestamp(ctx, orderID, model.OrderStatusConfirmed); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

// HandlePaymentFailure is called by the webhook when payment fails or expires.
func (s *StripeService) HandlePaymentFailure(ctx context.Context, sessionID string) error {
	sess, err := session.Get(sessionID, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve stripe session: %w", err)
	}

	orderIDStr, ok := sess.Metadata["order_id"]
	if !ok {
		return fmt.Errorf("order_id not found in session metadata")
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return fmt.Errorf("invalid order_id in metadata: %w", err)
	}

	// Update payment status to failed
	if err := s.orderRepo.UpdatePaymentStatus(ctx, orderID, model.PaymentStatusFailed); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

// VerifySession retrieves a Stripe checkout session and returns basic info.
func (s *StripeService) VerifySession(ctx context.Context, sessionID string) (*model.StripeSessionStatus, error) {
	sess, err := session.Get(sessionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stripe session: %w", err)
	}

	orderIDStr := sess.Metadata["order_id"]
	orderID, _ := uuid.Parse(orderIDStr)

	return &model.StripeSessionStatus{
		SessionID:     sess.ID,
		PaymentStatus: string(sess.PaymentStatus),
		OrderID:       orderID,
		OrderNumber:   sess.Metadata["order_number"],
	}, nil
}
