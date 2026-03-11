package store

import (
	"database/sql"
	"time"
)

type Product struct {
	ID              string          `json:"id,omitempty"`
	BusinessID      string          `json:"business_id,omitempty"`
	CategoryID      string          `json:"category_id,omitempty"`
	SubCategoryID   string          `json:"sub_category_id,omitempty"`
	Name            string          `json:"name,omitempty"`
	Description     string          `json:"description,omitempty"`
	Quantity        float64         `json:"quantity,omitempty"`
	Unit            string          `json:"unit,omitempty"`
	Price           float64         `json:"price,omitempty"`
	MOQ             string          `json:"moq,omitempty"`
	Images          []ProductImages `json:"product_images,omitempty"`
	IsProductActive bool            `json:"is_product_active,omitempty"`
	CreatedAT       time.Time       `json:"created_at"`
	UpdatedAT       time.Time       `json:"updated_at"`
}

type CompleteProduct struct {
	ID                     string          `json:"id"`
	BusinessID             string          `json:"business_id"`
	BusinessName           string          `json:"business_name"`
	BusinessEmail          string          `json:"business_email"`
	BusinessPhone          string          `json:"business_phone"`
	Address                string          `json:"address"`
	City                   string          `json:"city"`
	State                  string          `json:"state"`
	Pincode                string          `json:"pincode"`
	CategoryID             string          `json:"category_id"`
	CategoryName           string          `json:"category_name"`
	CategoryDescription    string          `json:"category_description"`
	SubCategoryID          string          `json:"sub_category_id"`
	SubCategoryName        string          `json:"sub_category_name"`
	SubCategoryDescription string          `json:"sub_category_description"`
	ProductName            string          `json:"product_name"`
	ProductDescription     string          `json:"product_description"`
	Quantity               float64         `json:"quantity"`
	Unit                   string          `json:"unit"`
	Price                  float64         `json:"price"`
	IsProductActive        bool            `json:"is_product_active"`
	MOQ                    string          `json:"moq"`
	CreatedAT              string          `json:"created_at"`
	Images                 []ProductImages `json:"product_images,omitempty"`
}

type ProductImages struct {
	ID         string    `json:"id"`
	ImageIndex int       `json:"image_index"`
	ProductID  string    `json:"product_id"`
	Image      string    `json:"image"`
	CreatedAT  time.Time `json:"created_at"`
	UpdatedAT  time.Time `json:"updated_at"`
}

type PostgresProductStore struct {
	db *sql.DB
}

type ProductStore interface {
	CreateProduct(*Product) error
	UpdateProduct(*Product) error
	ChangeProductActivateStatus(*Product) error
	DeleteProduct(id string) error
	GetAllProducts() ([]Product, error)
	GetBusinessProducts(id string) ([]Product, error)
	GetFollowersProducts(id string) ([]Product, error)
	GetCategoryBasedProducts(id string) ([]Product, error)
	GetSubCategoryBasedProducts(id string) ([]Product, error)
	GetProductDetailsByID(id string) (*CompleteProduct, error)
}

func NewPostgresProductStore(db *sql.DB) *PostgresProductStore {
	return &PostgresProductStore{
		db: db,
	}
}

func (ps *PostgresProductStore) CreateProduct(p *Product) error {
	query := `
	INSERT INTO products (
		business_id,
		category_id,
		sub_category_id,
		name,
		description,
		quantity,
		unit,
		price,
		moq,
		is_product_active
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10 
	) RETURNING id;
	`

	err := ps.db.QueryRow(
		query,
		p.BusinessID,
		p.CategoryID,
		p.SubCategoryID,
		p.Name,
		p.Description,
		p.Quantity,
		p.Unit,
		p.Price,
		p.MOQ,
		p.IsProductActive,
	).Scan(
		&p.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ps *PostgresProductStore) UpdateProduct(p *Product) error {
	query := `
	UPDATE products
	SET category_id = COALESCE($1, category_id),
	sub_category_id = COALESCE($2, sub_category_id),
	name = COALESCE($3, name),
	quantity = COALESCE($4, quantity),
	unit = COALESCE($5, unit),
	price = COALESCE($6, price),
	moq = COALESCE($7, moq),
	description = COALESCE($8, description),
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $9;
	`

	res, err := ps.db.Exec(
		query,
		p.CategoryID,
		p.SubCategoryID,
		p.Name,
		p.Quantity,
		p.Unit,
		p.Price,
		p.MOQ,
		p.Description,
		p.ID,
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

func (ps *PostgresProductStore) ChangeProductActivateStatus(p *Product) error {
	query := `
	UPDATE products
	SET is_product_active = $1,
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $2;
	`

	res, err := ps.db.Exec(
		query,
		p.IsProductActive,
		p.ID,
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

func (ps *PostgresProductStore) DeleteProduct(id string) error {
	query := `
	DELETE FROM products
	WHERE id = $1;
	`

	res, err := ps.db.Exec(query, id)
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

func (ps *PostgresProductStore) fetchProductImages(productID string) ([]ProductImages, error) {
	imageQuery := `
	SELECT 
		id,
		image_index,
		product_id,
		image,
		created_at,
		updated_at
	FROM product_images
	WHERE product_id = $1;
	`

	imgRes, err := ps.db.Query(imageQuery, productID)
	if err != nil {
		return nil, err
	}
	defer imgRes.Close()

	var images []ProductImages
	for imgRes.Next() {
		var img ProductImages
		err = imgRes.Scan(
			&img.ID,
			&img.ImageIndex,
			&img.ProductID,
			&img.Image,
			&img.CreatedAT,
			&img.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	if imgRes.Err() != nil {
		return nil, imgRes.Err()
	}

	return images, nil
}

func (ps *PostgresProductStore) GetAllProducts() ([]Product, error) {
	query := `
	SELECT 
		id,
		name,
		description,
		quantity,
		unit,
		price,
		moq,
		created_at,
		updated_at
	FROM products
	WHERE is_product_active = TRUE;
	`

	res, err := ps.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var products []Product
	for res.Next() {
		var p Product
		err = res.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Quantity,
			&p.Unit,
			&p.Price,
			&p.MOQ,
			&p.CreatedAT,
			&p.UpdatedAT,
		)

		if err != nil {
			return nil, err
		}

		images, err := ps.fetchProductImages(p.ID)
		if err != nil {
			return nil, err
		}
		p.Images = images

		products = append(products, p)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return products, nil
}

func (ps *PostgresProductStore) GetBusinessProducts(id string) ([]Product, error) {
	query := `
	SELECT 
		id,
		name,
		description,
		quantity,
		unit,
		price,
		moq,
		is_product_active,
		created_at,
		updated_at
	FROM products
	WHERE business_id = $1;
	`

	res, err := ps.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var products []Product
	for res.Next() {
		var p Product
		err = res.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Quantity,
			&p.Unit,
			&p.Price,
			&p.MOQ,
			&p.IsProductActive,
			&p.CreatedAT,
			&p.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}

		images, err := ps.fetchProductImages(p.ID)
		if err != nil {
			return nil, err
		}
		p.Images = images

		products = append(products, p)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return products, nil
}

func (ps *PostgresProductStore) GetFollowersProducts(id string) ([]Product, error) {

	businessIdsQuery := `
	SELECT
		business_id
	FROM followers
	WHERE user_id = $1;
	`

	bidres, err := ps.db.Query(
		businessIdsQuery,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer bidres.Close()

	var businessIds []string
	for bidres.Next() {
		var id string
		err = bidres.Scan(
			&id,
		)
		if err != nil {
			return nil, err
		}

		businessIds = append(businessIds, id)
	}

	if bidres.Err() != nil {
		return nil, bidres.Err()
	}

	productsQuery := `
	SELECT
		id,
		name,
		description,
		quantity,
		unit,
		price,
		moq,
		created_at,
		updated_at
	FROM products
	WHERE business_id = ANY($1)
	AND is_product_active = TRUE;
	`

	pres, err := ps.db.Query(
		productsQuery,
		businessIds,
	)
	if err != nil {
		return nil, err
	}
	defer pres.Close()

	var products []Product
	for pres.Next() {
		var p Product
		err = pres.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Quantity,
			&p.Unit,
			&p.Price,
			&p.MOQ,
			&p.CreatedAT,
			&p.UpdatedAT,
		)

		if err != nil {
			return nil, err
		}

		images, err := ps.fetchProductImages(p.ID)
		if err != nil {
			return nil, err
		}
		p.Images = images

		products = append(products, p)
	}

	if pres.Err() != nil {
		return nil, pres.Err()
	}

	return products, nil
}

func (ps *PostgresProductStore) GetCategoryBasedProducts(id string) ([]Product, error) {
	query := `
	SELECT 
		id,
		name,
		description,
		quantity,
		unit,
		price,
		moq,
		created_at,
		updated_at
	FROM products
	WHERE category_id = $1
	AND is_product_active = TRUE;
	`

	res, err := ps.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var products []Product
	for res.Next() {
		var p Product
		err = res.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Quantity,
			&p.Unit,
			&p.Price,
			&p.MOQ,
			&p.CreatedAT,
			&p.UpdatedAT,
		)

		if err != nil {
			return nil, err
		}

		images, err := ps.fetchProductImages(p.ID)
		if err != nil {
			return nil, err
		}
		p.Images = images

		products = append(products, p)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return products, nil
}

func (ps *PostgresProductStore) GetSubCategoryBasedProducts(id string) ([]Product, error) {
	query := `
	SELECT 
		id,
		name,
		description,
		quantity,
		unit,
		price,
		moq,
		created_at,
		updated_at
	FROM products
	WHERE sub_category_id = $1
	AND is_product_active = TRUE;
	`

	res, err := ps.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var products []Product
	for res.Next() {
		var p Product
		err = res.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Quantity,
			&p.Unit,
			&p.Price,
			&p.MOQ,
			&p.CreatedAT,
			&p.UpdatedAT,
		)

		if err != nil {
			return nil, err
		}

		images, err := ps.fetchProductImages(p.ID)
		if err != nil {
			return nil, err
		}
		p.Images = images

		products = append(products, p)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return products, nil
}

func (ps *PostgresProductStore) GetProductDetailsByID(id string) (*CompleteProduct, error) {
	query := `
	SELECT
		p.id,
		p.business_id,
		b.business_name,
		b.business_email,
		b.business_phone,
		b.address,
		b.city,
		b.state,
		b.pincode,
		p.category_id,
		c.name,
		c.description,
		p.sub_category_id,
		s.name,
		s.description,
		p.name,
		p.description,
		p.quantity,
		p.unit,
		p.price,
		p.moq,
		p.is_product_active,
		p.created_at
	FROM products p
	JOIN businesses b
		ON b.id = p.business_id
	JOIN categories c
		ON c.id = p.category_id
	JOIN sub_categories s
		ON s.id = p.sub_category_id
	WHERE p.id = $1;
	`
	var c CompleteProduct
	err := ps.db.QueryRow(
		query,
		id,
	).Scan(
		&c.ID,
		&c.BusinessID,
		&c.BusinessName,
		&c.BusinessEmail,
		&c.BusinessPhone,
		&c.Address,
		&c.City,
		&c.State,
		&c.Pincode,
		&c.CategoryID,
		&c.CategoryName,
		&c.CategoryDescription,
		&c.SubCategoryID,
		&c.SubCategoryName,
		&c.SubCategoryDescription,
		&c.ProductName,
		&c.ProductDescription,
		&c.Quantity,
		&c.Unit,
		&c.Price,
		&c.MOQ,
		&c.IsProductActive,
		&c.CreatedAT,
	)

	if err != nil {
		return nil, err
	}

	images, err := ps.fetchProductImages(c.ID)
	if err != nil {
		return nil, err
	}
	c.Images = images

	return &c, nil
}