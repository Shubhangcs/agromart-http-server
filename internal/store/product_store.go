package store

import (
	"database/sql"

	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

type PostgresProductStore struct {
	db *sql.DB
}

type ProductStore interface {
	CreateProduct(*models.Product) error
	UpdateProduct(*models.Product) error
	ChangeProductActivateStatus(*models.Product) error
	DeleteProduct(id string) error
	GetAllProducts(filter utils.ProductFilter, limit, offset int) ([]models.Product, error)
	GetBusinessProducts(id string, limit, offset int) ([]models.Product, error)
	GetFollowersProducts(id string, limit, offset int) ([]models.Product, error)
	GetCategoryBasedProducts(id string, limit, offset int) ([]models.Product, error)
	GetSubCategoryBasedProducts(id string, limit, offset int) ([]models.Product, error)
	GetProductDetailsByID(id string) (*models.CompleteProduct, error)
}

func NewPostgresProductStore(db *sql.DB) *PostgresProductStore {
	return &PostgresProductStore{db: db}
}

func (ps *PostgresProductStore) CreateProduct(p *models.Product) error {
	query := `
	INSERT INTO products (
		business_id, category_id, sub_category_id, name, description,
		quantity, unit, price, moq, is_product_active
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	RETURNING id
	`
	return ps.db.QueryRow(
		query,
		p.BusinessID, p.CategoryID, p.SubCategoryID, p.Name, p.Description,
		p.Quantity, p.Unit, p.Price, p.MOQ, p.IsProductActive,
	).Scan(&p.ID)
}

func (ps *PostgresProductStore) UpdateProduct(p *models.Product) error {
	query := `
	UPDATE products
	SET category_id     = COALESCE(NULLIF($1,''), category_id::text)::uuid,
	    sub_category_id = COALESCE(NULLIF($2,''), sub_category_id::text)::uuid,
	    name            = COALESCE(NULLIF($3,''), name),
	    quantity        = COALESCE(NULLIF($4,0),  quantity),
	    unit            = COALESCE(NULLIF($5,''), unit),
	    price           = COALESCE(NULLIF($6,0),  price),
	    moq             = COALESCE(NULLIF($7,''), moq),
	    description     = COALESCE(NULLIF($8,''), description),
	    updated_at      = CURRENT_TIMESTAMP
	WHERE id = $9
	`
	res, err := ps.db.Exec(query,
		p.CategoryID, p.SubCategoryID, p.Name, p.Quantity,
		p.Unit, p.Price, p.MOQ, p.Description, p.ID,
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

func (ps *PostgresProductStore) ChangeProductActivateStatus(p *models.Product) error {
	res, err := ps.db.Exec(
		`UPDATE products SET is_product_active = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
		p.IsProductActive, p.ID,
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

func (ps *PostgresProductStore) DeleteProduct(id string) error {
	res, err := ps.db.Exec(`DELETE FROM products WHERE id = $1`, id)
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

// fetchProductImages returns all images for a product in a single query.
func (ps *PostgresProductStore) fetchProductImages(productID string) ([]models.ProductImages, error) {
	rows, err := ps.db.Query(`
	SELECT id, image_index, product_id, image, created_at, updated_at
	FROM product_images
	WHERE product_id = $1
	ORDER BY image_index
	`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.ProductImages
	for rows.Next() {
		var img models.ProductImages
		err = rows.Scan(&img.ID, &img.ImageIndex, &img.ProductID, &img.Image, &img.CreatedAT, &img.UpdatedAT)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, rows.Err()
}

// GetAllProducts returns paginated active products with optional name/city/state filtering
// via a JOIN to the businesses table.
func (ps *PostgresProductStore) GetAllProducts(filter utils.ProductFilter, limit, offset int) ([]models.Product, error) {
	query := `
	SELECT p.id, b.user_id, p.name, p.description, p.quantity, p.unit, p.price, p.moq, p.created_at, p.updated_at
	FROM products p
	JOIN businesses b ON b.id = p.business_id
	WHERE p.is_product_active = TRUE
	  AND ($1 = '' OR p.name ILIKE '%' || $1 || '%')
	  AND ($2 = '' OR b.city ILIKE $2)
	  AND ($3 = '' OR b.state ILIKE $3)
	ORDER BY p.created_at DESC
	LIMIT $4 OFFSET $5
	`
	return ps.scanProducts(query, filter.Name, filter.City, filter.State, limit, offset)
}

func (ps *PostgresProductStore) GetBusinessProducts(id string, limit, offset int) ([]models.Product, error) {
	query := `
	SELECT p.id, b.user_id, p.name, p.description, p.quantity, p.unit, p.price, p.moq, p.is_product_active, p.created_at, p.updated_at
	FROM products p
	JOIN businesses b ON b.id = p.business_id
	WHERE p.business_id = $1
	ORDER BY p.created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := ps.db.Query(query, id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err = rows.Scan(
			&p.ID, &p.UserID, &p.Name, &p.Description, &p.Quantity, &p.Unit, &p.Price,
			&p.MOQ, &p.IsProductActive, &p.CreatedAT, &p.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}
		p.Images, err = ps.fetchProductImages(p.ID)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// GetFollowersProducts returns active products from businesses the given user
// follows, using a JOIN instead of two separate queries.
func (ps *PostgresProductStore) GetFollowersProducts(id string, limit, offset int) ([]models.Product, error) {
	query := `
	SELECT DISTINCT p.id, b.user_id, p.name, p.description, p.quantity, p.unit, p.price, p.moq, p.created_at, p.updated_at
	FROM products p
	JOIN businesses b ON b.id = p.business_id
	JOIN followers f ON f.business_id = p.business_id
	WHERE f.user_id = $1
	  AND p.is_product_active = TRUE
	ORDER BY p.created_at DESC
	LIMIT $2 OFFSET $3
	`
	return ps.scanProducts(query, id, limit, offset)
}

func (ps *PostgresProductStore) GetCategoryBasedProducts(id string, limit, offset int) ([]models.Product, error) {
	query := `
	SELECT p.id, b.user_id, p.name, p.description, p.quantity, p.unit, p.price, p.moq, p.created_at, p.updated_at
	FROM products p
	JOIN businesses b ON b.id = p.business_id
	WHERE p.category_id = $1
	  AND p.is_product_active = TRUE
	ORDER BY p.created_at DESC
	LIMIT $2 OFFSET $3
	`
	return ps.scanProducts(query, id, limit, offset)
}

func (ps *PostgresProductStore) GetSubCategoryBasedProducts(id string, limit, offset int) ([]models.Product, error) {
	query := `
	SELECT p.id, b.user_id, p.name, p.description, p.quantity, p.unit, p.price, p.moq, p.created_at, p.updated_at
	FROM products p
	JOIN businesses b ON b.id = p.business_id
	WHERE p.sub_category_id = $1
	  AND p.is_product_active = TRUE
	ORDER BY p.created_at DESC
	LIMIT $2 OFFSET $3
	`
	return ps.scanProducts(query, id, limit, offset)
}

// scanProducts is a helper that scans a 10-column product row (id, user_id, name, ...).
func (ps *PostgresProductStore) scanProducts(query string, args ...any) ([]models.Product, error) {
	rows, err := ps.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err = rows.Scan(
			&p.ID, &p.UserID, &p.Name, &p.Description, &p.Quantity, &p.Unit, &p.Price,
			&p.MOQ, &p.CreatedAT, &p.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}
		p.Images, err = ps.fetchProductImages(p.ID)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (ps *PostgresProductStore) GetProductDetailsByID(id string) (*models.CompleteProduct, error) {
	query := `
	SELECT
		p.id, b.user_id, p.business_id,
		b.business_name, b.business_email, b.business_phone, b.address, b.city, b.state, b.pincode,
		p.category_id, c.name, c.description,
		p.sub_category_id, s.name, s.description,
		p.name, p.description, p.quantity, p.unit, p.price, p.moq, p.is_product_active, p.created_at
	FROM products p
	JOIN businesses b ON b.id = p.business_id
	JOIN categories c ON c.id = p.category_id
	JOIN sub_categories s ON s.id = p.sub_category_id
	WHERE p.id = $1
	`
	var c models.CompleteProduct
	err := ps.db.QueryRow(query, id).Scan(
		&c.ID, &c.UserID, &c.BusinessID,
		&c.BusinessName, &c.BusinessEmail, &c.BusinessPhone, &c.Address, &c.City, &c.State, &c.Pincode,
		&c.CategoryID, &c.CategoryName, &c.CategoryDescription,
		&c.SubCategoryID, &c.SubCategoryName, &c.SubCategoryDescription,
		&c.ProductName, &c.ProductDescription, &c.Quantity, &c.Unit, &c.Price, &c.MOQ,
		&c.IsProductActive, &c.CreatedAT,
	)
	if err != nil {
		return nil, err
	}
	c.Images, err = ps.fetchProductImages(c.ID)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
