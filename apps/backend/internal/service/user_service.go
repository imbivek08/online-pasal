package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetOrCreateUser gets a user by Clerk ID, or creates one if it doesn't exist
func (s *UserService) GetOrCreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// Try to get existing user
	user, err := s.userRepo.GetByClerkID(ctx, req.ClerkID)
	if err == nil {
		// User exists, update last login
		_ = s.userRepo.UpdateLastLogin(ctx, req.ClerkID)
		return user, nil
	}

	// User doesn't exist, create new one
	user, err = s.userRepo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByClerkID retrieves a user by their Clerk ID
func (s *UserService) GetUserByClerkID(ctx context.Context, clerkID string) (*model.User, error) {
	user, err := s.userRepo.GetByClerkID(ctx, clerkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by their UUID
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(ctx context.Context, clerkID string, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.Update(ctx, clerkID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(ctx context.Context, clerkID string) error {
	err := s.userRepo.Delete(ctx, clerkID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// SyncUserFromClerk syncs user data from Clerk webhook
func (s *UserService) SyncUserFromClerk(ctx context.Context, clerkID, email string, firstName, lastName, username, avatarURL *string) (*model.User, error) {
	// Check if user exists
	existingUser, err := s.userRepo.GetByClerkID(ctx, clerkID)
	if err != nil {
		// User doesn't exist, create
		req := &model.CreateUserRequest{
			ClerkID:   clerkID,
			Email:     email,
			Username:  username,
			FirstName: firstName,
			LastName:  lastName,
			AvatarURL: avatarURL,
			Role:      model.RoleCustomer,
		}

		return s.userRepo.Create(ctx, req)
	}

	// User exists, update if needed
	updateReq := &model.UpdateUserRequest{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		AvatarURL: avatarURL,
	}

	return s.userRepo.Update(ctx, existingUser.ClerkID, updateReq)
}

// GetShopIDByVendorID retrieves the shop ID for a vendor
func (s *UserService) GetShopIDByVendorID(ctx context.Context, vendorID uuid.UUID) (uuid.UUID, error) {
	shopID, err := s.userRepo.GetShopIDByVendorID(ctx, vendorID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get shop ID: %w", err)
	}

	return shopID, nil
}

// ConvertToVendor converts a customer to vendor with business information
func (s *UserService) ConvertToVendor(ctx context.Context, userID uuid.UUID, req *model.BecomeVendorRequest) (*model.User, error) {
	// Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Validate current role
	if user.Role == model.RoleVendor {
		return nil, errors.New("already a vendor")
	}

	if user.Role == model.RoleAdmin {
		return nil, errors.New("admin cannot become vendor")
	}

	// Update user to vendor role with business info
	updateReq := &model.UpdateUserRequest{
		Phone: &req.Phone,
		Role:  model.RoleVendor,
	}

	// Update in database
	updatedUser, err := s.userRepo.Update(ctx, user.ClerkID, updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}
