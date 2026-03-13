package store

import (
	"database/sql"

	"github.com/shubhangcs/agromart-server/internal/models"
)

type PostgresCategoryStore struct {
	db *sql.DB
}

type CategoryStore interface {
	CreateCategory(*models.Category) error
	CreateSubCategory(*models.SubCategory) error
	UpdateCategory(*models.Category) error
	UpdateSubCategory(*models.SubCategory) error
	DeleteCategory(id string) error
	DeleteSubCategory(id string) error
	GetCategoryByID(id string) (*models.Category, error)
	GetSubCategoryByID(id string) (*models.SubCategory, error)
	GetAllCategories() ([]models.Category, error)
	GetAllSubCategories() ([]models.SubCategory, error)
	GetSubCategoriesByCategoryID(id string) ([]models.SubCategory, error)
}

func NewPostgresCategoryStore(db *sql.DB) *PostgresCategoryStore {
	return &PostgresCategoryStore{
		db: db,
	}
}

func (cs *PostgresCategoryStore) CreateCategory(c *models.Category) error {
	query := `
	INSERT INTO categories (
		name,
		description
	) VALUES (
		$1 , $2
	)RETURNING id;
	`

	err := cs.db.QueryRow(query, c.Name, c.Description).Scan(
		&c.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (cs *PostgresCategoryStore) CreateSubCategory(sc *models.SubCategory) error {
	query := `
	INSERT INTO sub_categories (
		category_id,
		name,
		description
	) VALUES (
		$1 , $2 , $3
	)RETURNING id;
	`

	err := cs.db.QueryRow(query, sc.CategoryID, sc.Name, sc.Description).Scan(
		&sc.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (cs *PostgresCategoryStore) UpdateCategory(c *models.Category) error {
	query := `
	UPDATE categories
	SET name = COALESCE($1, name),
	description = COALESCE($2, description),
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $3;
	`

	res, err := cs.db.Exec(
		query,
		c.Name,
		c.Description,
		c.ID,
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

func (cs *PostgresCategoryStore) UpdateSubCategory(sc *models.SubCategory) error {
	query := `
	UPDATE sub_categories
	SET name = COALESCE($1, name),
	description = COALESCE($2, description),
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $3;
	`

	res, err := cs.db.Exec(
		query,
		sc.Name,
		sc.Description,
		sc.ID,
	)

	if err != nil {
		return err
	}

	rowsAffecetd, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffecetd == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (cs *PostgresCategoryStore) DeleteCategory(id string) error {
	query := `
	DELETE FROM categories
	WHERE id = $1;
	`

	res, err := cs.db.Exec(query, id)
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

func (cs *PostgresCategoryStore) DeleteSubCategory(id string) error {
	query := `
	DELETE FROM sub_categories
	WHERE id = $1;
	`

	res, err := cs.db.Exec(
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

func (cs *PostgresCategoryStore) GetCategoryByID(id string) (*models.Category, error) {
	query := `
	SELECT
		id,
		category_image,
		name,
		description,
		created_at,
		updated_at
	FROM categories
	WHERE id = $1;
	`
	var category models.Category
	err := cs.db.QueryRow(
		query,
		id,
	).Scan(
		&category.ID,
		&category.CategoryImage,
		&category.Name,
		&category.Description,
		&category.CreatedAT,
		&category.UpdatedAT,
	)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (cs *PostgresCategoryStore) GetSubCategoryByID(id string) (*models.SubCategory, error) {
	query := `
	SELECT
		id,
		category_id,
		sub_category_image,
		name,
		description,
		created_at,
		updated_at
	FROM sub_categories
	WHERE id = $1;
	`
	var subCategory models.SubCategory
	err := cs.db.QueryRow(
		query,
		id,
	).Scan(
		&subCategory.ID,
		&subCategory.CategoryID,
		&subCategory.SubCategoryImage,
		&subCategory.Name,
		&subCategory.Description,
		&subCategory.CreatedAT,
		&subCategory.UpdatedAT,
	)

	if err != nil {
		return nil, err
	}

	return &subCategory, nil
}

func (cs *PostgresCategoryStore) GetAllCategories() ([]models.Category, error) {
	query := `
	SELECT
		id,
		category_image,
		name,
		description,
		created_at,
		updated_at
	FROM categories;
	`

	res, err := cs.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var categories []models.Category
	for res.Next() {
		var c models.Category
		err = res.Scan(
			&c.ID,
			&c.CategoryImage,
			&c.Name,
			&c.Description,
			&c.CreatedAT,
			&c.UpdatedAT,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return categories, nil
}

func (cs *PostgresCategoryStore) GetAllSubCategories() ([]models.SubCategory, error) {
	query := `
	SELECT
		id,
		category_id,
		sub_category_image,
		name,
		description,
		created_at,
		updated_at
	FROM sub_categories;
	`

	res, err := cs.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var subCategories []models.SubCategory
	for res.Next() {
		var sc models.SubCategory
		err = res.Scan(
			&sc.ID,
			&sc.CategoryID,
			&sc.SubCategoryImage,
			&sc.Name,
			&sc.Description,
			&sc.CreatedAT,
			&sc.UpdatedAT,
		)

		if err != nil {
			return nil, err
		}

		subCategories = append(subCategories, sc)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return subCategories, nil
}

func (cs *PostgresCategoryStore) GetSubCategoriesByCategoryID(id string) ([]models.SubCategory, error) {
	query := `
	SELECT
		id,
		category_id,
		sub_category_image,
		name,
		description,
		created_at,
		updated_at
	FROM sub_categories
	WHERE category_id = $1;
	`

	res, err := cs.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var subCategories []models.SubCategory
	for res.Next() {
		var sc models.SubCategory
		err = res.Scan(
			&sc.ID,
			&sc.CategoryID,
			&sc.SubCategoryImage,
			&sc.Name,
			&sc.Description,
			&sc.CreatedAT,
			&sc.UpdatedAT,
		)

		if err != nil {
			return nil, err
		}

		subCategories = append(subCategories, sc)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return subCategories, nil
}
