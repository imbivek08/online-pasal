package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/database"
	"github.com/imbivek08/hamropasal/internal/model"
)

type ReviewRepository struct {
	db *database.Database
}

func NewReviewRepository(db *database.Database) *ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

func (r *ReviewRepository) CreateReview(ctx context.Context, review *model.Review) error {
	query := `
		INSERT INTO reviews (id, product_id, user_id, order_id, rating, title, comment, is_verified_purchase, is_approved, helpful_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		review.ID,
		review.ProductID,
		review.UserID,
		review.OrderID,
		review.Rating,
		review.Title,
		review.Comment,
		review.IsVerifiedPurchase,
		review.IsApproved,
		review.HelpfulCount,
		review.CreatedAt,
		review.UpdatedAt,
	)
	return err
}

// GetByID retrieves a review by ID with user information
func (r *ReviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ReviewWithUser, error) {
	var review model.ReviewWithUser
	query := `
		SELECT r.id, r.product_id, r.user_id, r.order_id, r.rating, r.title, r.comment, 
		       r.is_verified_purchase, r.is_approved, r.helpful_count, r.created_at, r.updated_at,
		       u.first_name || ' ' || u.last_name as user_name, u.avatar_url as user_avatar
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&review.ID,
		&review.ProductID,
		&review.UserID,
		&review.OrderID,
		&review.Rating,
		&review.Title,
		&review.Comment,
		&review.IsVerifiedPurchase,
		&review.IsApproved,
		&review.HelpfulCount,
		&review.CreatedAt,
		&review.UpdatedAt,
		&review.UserName,
		&review.UserAvatar,
	)
	return &review, err
}

// GetByProductID retrieves reviews for a product with pagination
func (r *ReviewRepository) GetByProductID(ctx context.Context, productID uuid.UUID, limit, offset int) ([]model.ReviewWithUser, error) {
	query := `
		SELECT r.id, r.product_id, r.user_id, r.order_id, r.rating, r.title, r.comment, 
		       r.is_verified_purchase, r.is_approved, r.helpful_count, r.created_at, r.updated_at,
		       u.first_name || ' ' || u.last_name as user_name, u.avatar_url as user_avatar
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.product_id = $1 AND r.is_approved = true
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, productID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []model.ReviewWithUser
	for rows.Next() {
		var review model.ReviewWithUser
		err := rows.Scan(
			&review.ID,
			&review.ProductID,
			&review.UserID,
			&review.OrderID,
			&review.Rating,
			&review.Title,
			&review.Comment,
			&review.IsVerifiedPurchase,
			&review.IsApproved,
			&review.HelpfulCount,
			&review.CreatedAt,
			&review.UpdatedAt,
			&review.UserName,
			&review.UserAvatar,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// CountByProductID counts total reviews for a product
func (r *ReviewRepository) CountByProductID(ctx context.Context, productID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM reviews WHERE product_id = $1 AND is_approved = true`
	err := r.db.Pool.QueryRow(ctx, query, productID).Scan(&count)
	return count, err
}

// GetProductRatingStats gets rating statistics for a product
func (r *ReviewRepository) GetProductRatingStats(ctx context.Context, productID uuid.UUID) (*model.ProductRatingStats, error) {
	var stats model.ProductRatingStats
	stats.ProductID = productID

	query := `
		SELECT 
			COUNT(*) as total_reviews,
			COALESCE(AVG(rating), 0) as average_rating,
			COUNT(CASE WHEN rating = 5 THEN 1 END) as five_star_count,
			COUNT(CASE WHEN rating = 4 THEN 1 END) as four_star_count,
			COUNT(CASE WHEN rating = 3 THEN 1 END) as three_star_count,
			COUNT(CASE WHEN rating = 2 THEN 1 END) as two_star_count,
			COUNT(CASE WHEN rating = 1 THEN 1 END) as one_star_count
		FROM reviews
		WHERE product_id = $1 AND is_approved = true
	`
	err := r.db.Pool.QueryRow(ctx, query, productID).Scan(
		&stats.TotalReviews,
		&stats.AverageRating,
		&stats.FiveStarCount,
		&stats.FourStarCount,
		&stats.ThreeStarCount,
		&stats.TwoStarCount,
		&stats.OneStarCount,
	)
	return &stats, err
}

// UpdateReview updates a review
func (r *ReviewRepository) UpdateReview(ctx context.Context, id uuid.UUID, req *model.UpdateReviewRequest) error {
	query := `
		UPDATE reviews 
		SET rating = COALESCE($2, rating),
		    title = COALESCE($3, title),
		    comment = COALESCE($4, comment),
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, id, req.Rating, req.Title, req.Comment)
	return err
}

// DeleteReview deletes a review
func (r *ReviewRepository) DeleteReview(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reviews WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

// IncrementHelpfulCount increments the helpful count for a review
func (r *ReviewRepository) IncrementHelpfulCount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE reviews SET helpful_count = helpful_count + 1 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

// HasUserReviewedProduct checks if a user has already reviewed a product
func (r *ReviewRepository) HasUserReviewedProduct(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM reviews WHERE user_id = $1 AND product_id = $2)`
	err := r.db.Pool.QueryRow(ctx, query, userID, productID).Scan(&exists)
	return exists, err
}
