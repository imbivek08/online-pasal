package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/repository"
)

type ShopService struct {
	shopRepo *repository.ShopRepository
	userRepo *repository.UserRepository
}

func NewShopService(shopRepo *repository.ShopRepository, userRepo *repository.UserRepository) *ShopService {
	return &ShopService{
		shopRepo: shopRepo,
		userRepo: userRepo,
	}
}

// GetMyShop retrieves the shop for the authenticated vendor
func (s *ShopService) GetMyShop(ctx context.Context, vendorID uuid.UUID) (*model.Shop, error) {
	shop, err := s.shopRepo.GetByVendorID(ctx, vendorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No shop found, not an error
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	return shop, nil
}

// CreateShop creates a new shop for a vendor
func (s *ShopService) CreateShop(ctx context.Context, vendorID uuid.UUID, req *model.CreateShopRequest) (*model.Shop, error) {
	// Check if vendor exists and has vendor role
	user, err := s.userRepo.GetByID(ctx, vendorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("vendor not found")
		}
		return nil, fmt.Errorf("failed to get vendor: %w", err)
	}

	if user.Role != "vendor" {
		return nil, errors.New("user is not a vendor")
	}

	// Check if vendor already has a shop
	existingShop, err := s.shopRepo.GetByVendorID(ctx, vendorID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing shop: %w", err)
	}
	if existingShop != nil {
		return nil, errors.New("vendor already has a shop")
	}

	// Generate slug from shop name
	slug := s.generateSlug(req.Name)

	// Ensure slug is unique
	slugExists, err := s.shopRepo.SlugExists(ctx, slug, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if slugExists {
		slug = s.generateUniqueSlug(slug)
	}

	shop := &model.Shop{
		ID:          uuid.New(),
		VendorID:    vendorID,
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		Address:     req.Address,
		City:        req.City,
		State:       req.State,
		Country:     req.Country,
		PostalCode:  req.PostalCode,
		Phone:       req.Phone,
		Email:       req.Email,
		IsActive:    true,
		IsVerified:  false,
	}

	if err := s.shopRepo.Create(ctx, shop); err != nil {
		return nil, fmt.Errorf("failed to create shop: %w", err)
	}

	return shop, nil
}

// GetShopByID retrieves a shop by ID
func (s *ShopService) GetShopByID(ctx context.Context, id uuid.UUID) (*model.Shop, error) {
	shop, err := s.shopRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("shop not found")
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}
	return shop, nil
}

// GetShopBySlug retrieves a shop by slug
func (s *ShopService) GetShopBySlug(ctx context.Context, slug string) (*model.Shop, error) {
	shop, err := s.shopRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("shop not found")
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}
	return shop, nil
}

// ListShops retrieves shops with pagination and filters
func (s *ShopService) ListShops(ctx context.Context, page, pageSize int, search string, activeOnly bool) (*model.ShopListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	shops, total, err := s.shopRepo.List(ctx, page, pageSize, search, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to list shops: %w", err)
	}

	shopResponses := make([]model.ShopResponse, len(shops))
	for i, shop := range shops {
		shopResponses[i] = shop.ToResponse()
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &model.ShopListResponse{
		Shops:      shopResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateShop updates a shop
func (s *ShopService) UpdateShop(ctx context.Context, shopID, vendorID uuid.UUID, req *model.UpdateShopRequest) (*model.Shop, error) {
	// Get existing shop
	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("shop not found")
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	// Verify ownership
	if shop.VendorID != vendorID {
		return nil, errors.New("unauthorized: you don't own this shop")
	}

	// Update fields if provided
	if req.Name != nil {
		shop.Name = *req.Name
	}
	if req.Description != nil {
		shop.Description = req.Description
	}
	if req.LogoURL != nil {
		shop.LogoURL = req.LogoURL
	}
	if req.BannerURL != nil {
		shop.BannerURL = req.BannerURL
	}
	if req.Address != nil {
		shop.Address = req.Address
	}
	if req.City != nil {
		shop.City = req.City
	}
	if req.State != nil {
		shop.State = req.State
	}
	if req.Country != nil {
		shop.Country = req.Country
	}
	if req.PostalCode != nil {
		shop.PostalCode = req.PostalCode
	}
	if req.Phone != nil {
		shop.Phone = req.Phone
	}
	if req.Email != nil {
		shop.Email = req.Email
	}

	if err := s.shopRepo.Update(ctx, shop); err != nil {
		return nil, fmt.Errorf("failed to update shop: %w", err)
	}

	return shop, nil
}

// ToggleShopStatus toggles shop's active status
func (s *ShopService) ToggleShopStatus(ctx context.Context, shopID, vendorID uuid.UUID) (*model.Shop, error) {
	// Get existing shop
	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("shop not found")
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	// Verify ownership
	if shop.VendorID != vendorID {
		return nil, errors.New("unauthorized: you don't own this shop")
	}

	// Toggle status
	newStatus := !shop.IsActive
	if err := s.shopRepo.UpdateStatus(ctx, shopID, newStatus); err != nil {
		return nil, fmt.Errorf("failed to update shop status: %w", err)
	}

	shop.IsActive = newStatus
	return shop, nil
}

// GetShopStats retrieves shop statistics
func (s *ShopService) GetShopStats(ctx context.Context, shopID, vendorID uuid.UUID) (*model.ShopWithStats, error) {
	// Get shop to verify ownership
	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("shop not found")
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	// Verify ownership
	if shop.VendorID != vendorID {
		return nil, errors.New("unauthorized: you don't own this shop")
	}

	stats, err := s.shopRepo.GetStats(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop stats: %w", err)
	}

	return stats, nil
}

// DeleteShop soft deletes a shop (admin only)
func (s *ShopService) DeleteShop(ctx context.Context, shopID uuid.UUID) error {
	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("shop not found")
		}
		return fmt.Errorf("failed to get shop: %w", err)
	}

	if err := s.shopRepo.Delete(ctx, shop.ID); err != nil {
		return fmt.Errorf("failed to delete shop: %w", err)
	}

	return nil
}

// VerifyShop verifies a shop (admin only)
func (s *ShopService) VerifyShop(ctx context.Context, shopID uuid.UUID, verified bool) (*model.Shop, error) {
	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("shop not found")
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	if err := s.shopRepo.UpdateVerification(ctx, shopID, verified); err != nil {
		return nil, fmt.Errorf("failed to verify shop: %w", err)
	}

	shop.IsVerified = verified
	return shop, nil
}

// generateSlug converts a string to a URL-friendly slug
func (s *ShopService) generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// Limit length
	if len(slug) > 100 {
		slug = slug[:100]
	}

	return slug
}

// generateUniqueSlug appends a UUID suffix to make slug unique
func (s *ShopService) generateUniqueSlug(baseSlug string) string {
	suffix := uuid.New().String()[:8]
	return fmt.Sprintf("%s-%s", baseSlug, suffix)
}
