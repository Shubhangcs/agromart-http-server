package store

import (
	"database/sql"

	"github.com/shubhangcs/agromart-server/internal/models"
)

// ReviewStore defines operations for business and product reviews.
type ReviewStore interface {
	CreateBusinessReview(*models.BusinessReview) error
	UpdateBusinessReview(*models.BusinessReview) error
	DeleteBusinessReview(id string) error
	GetBusinessReviews(businessID string, limit, offset int) ([]models.BusinessReview, error)

	CreateProductReview(*models.ProductReview) error
	UpdateProductReview(*models.ProductReview) error
	DeleteProductReview(id string) error
	GetProductReviews(productID string, limit, offset int) ([]models.ProductReview, error)
}

// PostgresReviewStore is the Postgres-backed implementation of ReviewStore.
type PostgresReviewStore struct {
	db *sql.DB
}

func NewPostgresReviewStore(db *sql.DB) *PostgresReviewStore {
	return &PostgresReviewStore{db: db}
}

// --- Business reviews ---

// CreateBusinessReview inserts a new review and populates id, created_at, updated_at.
func (s *PostgresReviewStore) CreateBusinessReview(r *models.BusinessReview) error {
	query := `
	INSERT INTO business_reviews (business_id, user_id, review)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`
	return s.db.QueryRow(query, r.BusinessID, r.UserID, r.Review).
		Scan(&r.ID, &r.CreatedAT, &r.UpdatedAT)
}

// UpdateBusinessReview updates the review text for the given review ID.
func (s *PostgresReviewStore) UpdateBusinessReview(r *models.BusinessReview) error {
	res, err := s.db.Exec(
		`UPDATE business_reviews SET review = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
		r.Review, r.ID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteBusinessReview deletes the review with the given ID.
func (s *PostgresReviewStore) DeleteBusinessReview(id string) error {
	res, err := s.db.Exec(`DELETE FROM business_reviews WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetBusinessReviews returns paginated reviews for a business, including reviewer name.
func (s *PostgresReviewStore) GetBusinessReviews(businessID string, limit, offset int) ([]models.BusinessReview, error) {
	query := `
	SELECT br.id, br.business_id, br.user_id,
	       CONCAT(u.first_name, ' ', COALESCE(u.last_name, '')) AS user_name,
	       br.review, br.created_at, br.updated_at
	FROM business_reviews br
	JOIN users u ON u.id = br.user_id
	WHERE br.business_id = $1
	ORDER BY br.created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := s.db.Query(query, businessID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.BusinessReview
	for rows.Next() {
		var rev models.BusinessReview
		err = rows.Scan(
			&rev.ID, &rev.BusinessID, &rev.UserID, &rev.UserName,
			&rev.Review, &rev.CreatedAT, &rev.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, rows.Err()
}

// --- Product reviews ---

// CreateProductReview inserts a new product review and populates id, created_at, updated_at.
func (s *PostgresReviewStore) CreateProductReview(r *models.ProductReview) error {
	query := `
	INSERT INTO product_reviews (product_id, user_id, review)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`
	return s.db.QueryRow(query, r.ProductID, r.UserID, r.Review).
		Scan(&r.ID, &r.CreatedAT, &r.UpdatedAT)
}

// UpdateProductReview updates the review text for the given review ID.
func (s *PostgresReviewStore) UpdateProductReview(r *models.ProductReview) error {
	res, err := s.db.Exec(
		`UPDATE product_reviews SET review = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
		r.Review, r.ID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteProductReview deletes the review with the given ID.
func (s *PostgresReviewStore) DeleteProductReview(id string) error {
	res, err := s.db.Exec(`DELETE FROM product_reviews WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetProductReviews returns paginated reviews for a product, including reviewer name.
func (s *PostgresReviewStore) GetProductReviews(productID string, limit, offset int) ([]models.ProductReview, error) {
	query := `
	SELECT pr.id, pr.product_id, pr.user_id,
	       CONCAT(u.first_name, ' ', COALESCE(u.last_name, '')) AS user_name,
	       pr.review, pr.created_at, pr.updated_at
	FROM product_reviews pr
	JOIN users u ON u.id = pr.user_id
	WHERE pr.product_id = $1
	ORDER BY pr.created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := s.db.Query(query, productID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.ProductReview
	for rows.Next() {
		var rev models.ProductReview
		err = rows.Scan(
			&rev.ID, &rev.ProductID, &rev.UserID, &rev.UserName,
			&rev.Review, &rev.CreatedAT, &rev.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, rows.Err()
}
