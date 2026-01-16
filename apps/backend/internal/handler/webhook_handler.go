package handler

import (
	"encoding/json"
	"net/http"

	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/service"
	"github.com/labstack/echo/v4"
)

type WebhookHandler struct {
	userService *service.UserService
}

func NewWebhookHandler(userService *service.UserService) *WebhookHandler {
	return &WebhookHandler{
		userService: userService,
	}
}

// ClerkWebhookData represents common webhook data structure
type ClerkWebhookEvent struct {
	Type   string          `json:"type"`
	Object string          `json:"object"`
	Data   json.RawMessage `json:"data"`
}

// ClerkUserData represents user data from Clerk webhook
type ClerkUserData struct {
	ID                    string              `json:"id"`
	EmailAddresses        []ClerkEmailAddress `json:"email_addresses"`
	FirstName             *string             `json:"first_name"`
	LastName              *string             `json:"last_name"`
	Username              *string             `json:"username"`
	ImageURL              *string             `json:"image_url"`
	PrimaryEmailAddressID *string             `json:"primary_email_address_id"`
}

type ClerkEmailAddress struct {
	ID           string `json:"id"`
	EmailAddress string `json:"email_address"`
}

// HandleClerkWebhook processes Clerk webhooks
func (h *WebhookHandler) HandleClerkWebhook(c echo.Context) error {
	// Get the webhook signature from headers
	svixID := c.Request().Header.Get("svix-id")
	svixTimestamp := c.Request().Header.Get("svix-timestamp")
	svixSignature := c.Request().Header.Get("svix-signature")

	if svixID == "" || svixTimestamp == "" || svixSignature == "" {
		return SendError(c, http.StatusBadRequest, nil, "missing webhook headers")
	}

	// Bind the webhook data
	var event ClerkWebhookEvent
	if err := c.Bind(&event); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid webhook payload")
	}

	// Handle different event types
	switch event.Type {
	case "user.created":
		return h.handleUserCreated(c, event.Data)
	case "user.updated":
		return h.handleUserUpdated(c, event.Data)
	case "user.deleted":
		return h.handleUserDeleted(c, event.Data)
	default:
		// Unknown event type, just acknowledge
		return SendSuccess(c, http.StatusOK, "webhook received", nil)
	}
}

func (h *WebhookHandler) handleUserCreated(c echo.Context, data json.RawMessage) error {
	var userData ClerkUserData
	if err := json.Unmarshal(data, &userData); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid user data")
	}

	// Get primary email
	email := h.getPrimaryEmail(&userData)
	if email == "" {
		return SendError(c, http.StatusBadRequest, nil, "no email address found")
	}

	// Create user in database
	req := &model.CreateUserRequest{
		ClerkID:   userData.ID,
		Email:     email,
		Username:  userData.Username,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		AvatarURL: userData.ImageURL,
		Role:      model.RoleCustomer,
	}

	user, err := h.userService.GetOrCreateUser(c.Request().Context(), req)
	if err != nil {
		return SendInternalError(c, err)
	}

	return SendSuccess(c, http.StatusOK, "user created successfully", user.ToResponse())
}

func (h *WebhookHandler) handleUserUpdated(c echo.Context, data json.RawMessage) error {
	var userData ClerkUserData
	if err := json.Unmarshal(data, &userData); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid user data")
	}

	// Get primary email
	email := h.getPrimaryEmail(&userData)
	if email == "" {
		return SendError(c, http.StatusBadRequest, nil, "no email address found")
	}

	// Sync user data
	user, err := h.userService.SyncUserFromClerk(
		c.Request().Context(),
		userData.ID,
		email,
		userData.FirstName,
		userData.LastName,
		userData.Username,
		userData.ImageURL,
	)
	if err != nil {
		return SendInternalError(c, err)
	}

	return SendSuccess(c, http.StatusOK, "user updated successfully", user.ToResponse())
}

func (h *WebhookHandler) handleUserDeleted(c echo.Context, data json.RawMessage) error {
	var userData ClerkUserData
	if err := json.Unmarshal(data, &userData); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid user data")
	}

	// Soft delete user
	err := h.userService.DeleteUser(c.Request().Context(), userData.ID)
	if err != nil {
		return SendInternalError(c, err)
	}

	return SendSuccess(c, http.StatusOK, "user deleted successfully", nil)
}

// getPrimaryEmail extracts the primary email from Clerk user data
func (h *WebhookHandler) getPrimaryEmail(userData *ClerkUserData) string {
	if len(userData.EmailAddresses) == 0 {
		return ""
	}

	// If primary email ID is specified, find it
	if userData.PrimaryEmailAddressID != nil {
		for _, email := range userData.EmailAddresses {
			if email.ID == *userData.PrimaryEmailAddressID {
				return email.EmailAddress
			}
		}
	}

	// Otherwise, return the first email
	return userData.EmailAddresses[0].EmailAddress
}
