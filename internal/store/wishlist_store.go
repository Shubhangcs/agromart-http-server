package store

import (
	"database/sql"

	"github.com/shubhangcs/agromart-server/internal/models"
)

// WishlistStore defines the data operations for the wishlist feature.
type WishlistStore interface {
	// AddToWishlist adds a product to a user's wishlist. Duplicate adds are silently ignored.
	AddToWishlist(userID, productID string) error
	// RemoveFromWishlist removes a product from a user's wishlist.
	RemoveFromWishlist(userID, productID string) error
	// GetUserWishlist returns all wishlist items for a user with product details, paginated.
	GetUserWishlist(userID string, limit, offset int) ([]models.WishlistItemModel, error)
	// IsInWishlist reports whether a product is already in the user's wishlist.
	IsInWishlist(userID, productID string) (bool, error)
}

// PostgresWishlistStore is the Postgres implementation of WishlistStore.
type PostgresWishlistStore struct {
	db *sql.DB
}

func NewPostgresWishlistStore(db *sql.DB) *PostgresWishlistStore {
	return &PostgresWishlistStore{db: db}
}

func (ws *PostgresWishlistStore) AddToWishlist(userID, productID string) error {
	_, err := ws.db.Exec(`
		INSERT INTO wishlists (user_id, product_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, product_id) DO NOTHING
	`, userID, productID)
	return err
}

func (ws *PostgresWishlistStore) RemoveFromWishlist(userID, productID string) error {
	res, err := ws.db.Exec(`
		DELETE FROM wishlists WHERE user_id = $1 AND product_id = $2
	`, userID, productID)
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

func (ws *PostgresWishlistStore) GetUserWishlist(userID string, limit, offset int) ([]models.WishlistItemModel, error) {
	query := `
	SELECT
		w.id, w.user_id, w.product_id,
		p.name, p.description, p.price, p.unit, p.moq, p.business_id,
		w.created_at
	FROM wishlists w
	JOIN products p ON p.id = w.product_id
	WHERE w.user_id = $1
	ORDER BY w.created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := ws.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.WishlistItemModel
	for rows.Next() {
		var item models.WishlistItemModel
		if err = rows.Scan(
			&item.ID, &item.UserID, &item.ProductID,
			&item.ProductName, &item.Description, &item.Price, &item.Unit, &item.MOQ, &item.BusinessID,
			&item.CreatedAT,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (ws *PostgresWishlistStore) IsInWishlist(userID, productID string) (bool, error) {
	var exists bool
	err := ws.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM wishlists WHERE user_id = $1 AND product_id = $2)
	`, userID, productID).Scan(&exists)
	return exists, err
}
