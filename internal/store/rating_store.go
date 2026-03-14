package store

import (
	"database/sql"

	"github.com/shubhangcs/agromart-server/internal/models"
)

// RatingStore defines operations for product and business ratings.
type RatingStore interface {
	RateProduct(*models.ProductRating) error
	GetAverageProductRating(productID string) (float64, error)
	GetRatingsByProductID(productID string, limit, offset int) ([]models.ProductRating, error)
	DeleteProductRating(id string) error

	RateBusiness(*models.BusinessRating) error
	GetAverageBusinessRating(id string) (float64, error)
	GetRatingsByBusinessID(id string) ([]models.BusinessRating, error)
}

// PostgresRatingStore is the Postgres-backed implementation.
type PostgresRatingStore struct {
	db *sql.DB
}

func NewPostgresRatingStore(db *sql.DB) *PostgresRatingStore {
	return &PostgresRatingStore{db: db}
}

// --- Product Ratings ---

// RateProduct upserts a rating for a product by a user.
func (s *PostgresRatingStore) RateProduct(r *models.ProductRating) error {
	query := `
	INSERT INTO product_ratings (product_id, user_id, rating)
	VALUES ($1, $2, $3)
	ON CONFLICT (product_id, user_id) DO UPDATE
	    SET rating     = EXCLUDED.rating,
	        updated_at = CURRENT_TIMESTAMP
	RETURNING id, created_at, updated_at
	`
	return s.db.QueryRow(query, r.ProductID, r.UserID, r.Rating).
		Scan(&r.ID, &r.CreatedAT, &r.UpdatedAT)
}

// GetAverageProductRating returns the average rating for the given product.
// Returns 0 when no ratings exist.
func (s *PostgresRatingStore) GetAverageProductRating(productID string) (float64, error) {
	var avg sql.NullFloat64
	err := s.db.QueryRow(
		`SELECT AVG(rating) FROM product_ratings WHERE product_id = $1`,
		productID,
	).Scan(&avg)
	if err != nil {
		return 0, err
	}
	if !avg.Valid {
		return 0, nil
	}
	return avg.Float64, nil
}

// GetRatingsByProductID returns paginated ratings for a product, including the rater's name.
func (s *PostgresRatingStore) GetRatingsByProductID(productID string, limit, offset int) ([]models.ProductRating, error) {
	query := `
	SELECT pr.id, pr.product_id, pr.user_id,
	       CONCAT(u.first_name, ' ', COALESCE(u.last_name, '')) AS user_name,
	       pr.rating, pr.created_at, pr.updated_at
	FROM product_ratings pr
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

	var ratings []models.ProductRating
	for rows.Next() {
		var pr models.ProductRating
		err = rows.Scan(
			&pr.ID, &pr.ProductID, &pr.UserID, &pr.UserName,
			&pr.Rating, &pr.CreatedAT, &pr.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, pr)
	}
	return ratings, rows.Err()
}

// DeleteProductRating deletes a product rating by its ID.
func (s *PostgresRatingStore) DeleteProductRating(id string) error {
	res, err := s.db.Exec(`DELETE FROM product_ratings WHERE id = $1`, id)
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

// --- Business Ratings ---

// RateBusiness inserts or updates a user's rating for a business (upsert).
func (s *PostgresRatingStore) RateBusiness(r *models.BusinessRating) error {
	query := `
	INSERT INTO business_ratings (business_id, user_id, rating)
	VALUES ($1, $2, $3)
	ON CONFLICT (business_id, user_id)
	DO UPDATE SET rating = EXCLUDED.rating, updated_at = CURRENT_TIMESTAMP
	`
	_, err := s.db.Exec(query, r.BusinessID, r.UserID, r.Rating)
	return err
}

func (s *PostgresRatingStore) GetAverageBusinessRating(id string) (float64, error) {
	query := `SELECT COALESCE(AVG(rating), 0)::NUMERIC(3,1) FROM business_ratings WHERE business_id = $1`
	var avg float64
	err := s.db.QueryRow(query, id).Scan(&avg)
	if err != nil {
		return 0, err
	}
	return avg, nil
}

func (s *PostgresRatingStore) GetRatingsByBusinessID(id string) ([]models.BusinessRating, error) {
	query := `
	SELECT r.id, r.business_id, r.user_id,
	       CONCAT(u.first_name, ' ', COALESCE(u.last_name, '')) AS user_name,
	       r.rating, r.created_at, r.updated_at
	FROM business_ratings r
	JOIN users u ON u.id = r.user_id
	WHERE r.business_id = $1
	ORDER BY r.created_at DESC
	`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []models.BusinessRating
	for rows.Next() {
		var rating models.BusinessRating
		err = rows.Scan(
			&rating.ID, &rating.BusinessID, &rating.UserID,
			&rating.UserName, &rating.Rating,
			&rating.CreatedAT, &rating.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}
	return ratings, rows.Err()
}

