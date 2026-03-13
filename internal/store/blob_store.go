package store

import (
	"database/sql"

	"github.com/shubhangcs/agromart-server/internal/models"
)

type PostgresBlobStore struct {
	db *sql.DB
}

type BlobStore interface {
	UpdateAdminProfileImage(id, path string) error
	UpdateUserProfileImage(id, path string) error
	UpdateBusinessProfileImage(id, path string) error
	UpdateCategoryImage(id, path string) error
	UpdateSubCategoryImage(id, path string) error
	UpdateProductImage(*models.ProductImage) error
	DeleteProductImage(id string) error
}

func NewPostgresBlobStore(db *sql.DB) *PostgresBlobStore {
	return &PostgresBlobStore{
		db: db,
	}
}

func (bs *PostgresBlobStore) UpdateAdminProfileImage(id, path string) error {
	query := `
	UPDATE admins
	SET profile_image = $1,
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $2
	`

	res, err := bs.db.Exec(query, path, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (bs *PostgresBlobStore) UpdateUserProfileImage(id, path string) error {
	query := `
	UPDATE users
	SET profile_image = $1,
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $2
	`

	res, err := bs.db.Exec(query, path, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (bs *PostgresBlobStore) UpdateBusinessProfileImage(id, path string) error {
	query := `
	UPDATE businesses
	SET business_profile_image = $1
	WHERE id = $2;
	`

	res, err := bs.db.Exec(
		query,
		path,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (bs *PostgresBlobStore) UpdateCategoryImage(id, path string) error {
	query := `
	UPDATE categories
	SET category_image = $1,
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $2;
	`

	res, err := bs.db.Exec(
		query,
		path,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (bs *PostgresBlobStore) UpdateSubCategoryImage(id, path string) error {
	query := `
	UPDATE sub_categories
	SET sub_category_image = $1,
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $2
	`

	res, err := bs.db.Exec(
		query,
		path,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (bs *PostgresBlobStore) UpdateProductImage(p *models.ProductImage) error {
	query := `
	INSERT INTO product_images(id, product_id, image_index, image)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (product_id, image_index) DO UPDATE
	    SET image = EXCLUDED.image, updated_at = CURRENT_TIMESTAMP
	`

	res, err := bs.db.Exec(
		query,
		p.ID,
		p.ProductID,
		p.Index,
		p.Image,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (bs *PostgresBlobStore) DeleteProductImage(id string) error {
	query := `
	DELETE FROM product_images
	WHERE id = $1;
	`

	res, err := bs.db.Exec(
		query,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
