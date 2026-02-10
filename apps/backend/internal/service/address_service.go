package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/repository"
)

type AddressService struct {
	addressRepo *repository.AddressRepository
}

func NewAddressService(addressRepo *repository.AddressRepository) *AddressService {
	return &AddressService{addressRepo: addressRepo}
}

// GetUserAddresses retrieves all addresses for a user
func (s *AddressService) GetUserAddresses(ctx context.Context, userID uuid.UUID) ([]*model.Address, error) {
	return s.addressRepo.GetByUserID(ctx, userID)
}

// GetAddressByID retrieves an address and verifies ownership
func (s *AddressService) GetAddressByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Address, error) {
	addr, err := s.addressRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("address not found")
	}

	if addr.UserID != userID {
		return nil, fmt.Errorf("address not found")
	}

	return addr, nil
}

// GetDefaultAddress retrieves the user's default address
func (s *AddressService) GetDefaultAddress(ctx context.Context, userID uuid.UUID) (*model.Address, error) {
	return s.addressRepo.GetDefaultByUserID(ctx, userID)
}

// CreateAddress creates a new address for a user
func (s *AddressService) CreateAddress(ctx context.Context, userID uuid.UUID, input *model.AddressInput) (*model.Address, error) {
	now := time.Now()

	addr := &model.Address{
		ID:           uuid.New(),
		UserID:       userID,
		FullName:     input.FullName,
		Phone:        input.Phone,
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		City:         input.City,
		State:        input.State,
		PostalCode:   input.PostalCode,
		Country:      input.Country,
		IsDefault:    input.IsDefault,
		AddressType:  "shipping",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// If this address should be default, unset other defaults first
	if addr.IsDefault {
		if err := s.addressRepo.UnsetDefaultForUser(ctx, userID); err != nil {
			return nil, fmt.Errorf("failed to update default address: %w", err)
		}
	}

	if err := s.addressRepo.Create(ctx, addr); err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	return addr, nil
}

// UpdateAddress updates an existing address (with ownership check)
func (s *AddressService) UpdateAddress(ctx context.Context, id uuid.UUID, userID uuid.UUID, input *model.AddressInput) (*model.Address, error) {
	// Verify ownership
	existing, err := s.addressRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("address not found")
	}
	if existing.UserID != userID {
		return nil, fmt.Errorf("address not found")
	}

	existing.FullName = input.FullName
	existing.Phone = input.Phone
	existing.AddressLine1 = input.AddressLine1
	existing.AddressLine2 = input.AddressLine2
	existing.City = input.City
	existing.State = input.State
	existing.PostalCode = input.PostalCode
	existing.Country = input.Country
	existing.UpdatedAt = time.Now()

	if err := s.addressRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}

	return existing, nil
}

// DeleteAddress deletes an address (with ownership check)
func (s *AddressService) DeleteAddress(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.addressRepo.Delete(ctx, id, userID)
}

// SetDefaultAddress sets an address as the user's default
func (s *AddressService) SetDefaultAddress(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Verify ownership first
	addr, err := s.addressRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("address not found")
	}
	if addr.UserID != userID {
		return fmt.Errorf("address not found")
	}

	return s.addressRepo.SetDefault(ctx, id, userID)
}
