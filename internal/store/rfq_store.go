package store

import (
	"database/sql"
	"time"
)

type RFQ struct {
	ID            string    `json:"id"`
	BusinessID    string    `json:"business_id"`
	BusinessName  string    `json:"business_name"`
	BusinessPhone string    `json:"business_phone"`
	BusinessEmail string    `json:"business_email"`
	Address       string    `json:"address"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	CategoryID    string    `json:"category_id"`
	SubCategoryID string    `json:"sub_category_id"`
	ProductName   string    `json:"product_name"`
	Quantity      float64   `json:"quantity"`
	Unit          string    `json:"unit"`
	Price         float64   `json:"price"`
	IsRFQActive   bool      `json:"is_rfq_active"`
	CreatedAT     time.Time `json:"created_at"`
	UpdatedAT     time.Time `json:"updated_at"`
}

type PostgresRFQStore struct {
	db *sql.DB
}

type RFQStore interface {
	CreateRFQ(*RFQ) error
	ActivateRFQ(*RFQ) error
	UpdateRFQ(*RFQ) error
	DeleteRFQ(id string) error
	GetAllRFQ() ([]RFQ, error)
	GetRFQByBusinessID(id string) ([]RFQ, error)
}

func NewPostgresRFQStore(db *sql.DB) *PostgresRFQStore {
	return &PostgresRFQStore{
		db: db,
	}
}

func (rs *PostgresRFQStore) CreateRFQ(r *RFQ) error {
	query := `
	INSERT INTO rfqs (
		business_id,
		category_id,
		sub_category_id,
		product_name,
		quantity,
		unit,
		price,
		is_rfq_active
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8 
	);
	`

	res, err := rs.db.Exec(
		query,
		r.BusinessID,
		r.CategoryID,
		r.SubCategoryID,
		r.ProductName,
		r.Quantity,
		r.Unit,
		r.Price,
		r.IsRFQActive,
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

func (rs *PostgresRFQStore) ActivateRFQ(r *RFQ) error {
	query := `
	UPDATE rfqs
	SET is_rfq_active = $1,
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $2;
	`

	res, err := rs.db.Exec(query, r.IsRFQActive, r.ID)
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

func (rs *PostgresRFQStore) UpdateRFQ(r *RFQ) error {
	query := `
	UPDATE rfqs 
	SET category_id = COALESCE($1, category_id),
	sub_category_id = COALESCE($2, sub_category_id),
	product_name = COALESCE($3, product_name),
	quantity = COALESCE($4, quantity),
	unit = COALESCE($5, unit),
	price = COALESCE($6, price),
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $7;
	`

	res, err := rs.db.Exec(
		query,
		r.CategoryID,
		r.SubCategoryID,
		r.ProductName,
		r.Quantity,
		r.Unit,
		r.Price,
		r.ID,
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

func (rs *PostgresRFQStore) DeleteRFQ(id string) error {
	query := `
	DELETE FROM rfqs
	WHERE id = $1;
	`

	res, err := rs.db.Exec(query, id)
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

func (rs *PostgresRFQStore) GetAllRFQ() ([]RFQ, error) {
	query := `
	SELECT
		r.id,
		b.id,
		b.business_name,
		b.business_email,
		b.business_phone,
		b.address,
		b.city,
		b.state,
		r.product_name,
		r.quantity,
		r.unit,
		r.price,
		r.is_rfq_active,
		r.created_at,
		r.updated_at
	FROM rfqs r
	JOIN businesses b
		ON b.id = r.business_id
	WHERE r.is_rfq_active = TRUE;
	`

	res, err := rs.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rfqs []RFQ
	for res.Next() {
		var r RFQ
		err = res.Scan(
			&r.ID,
			&r.BusinessID,
			&r.BusinessName,
			&r.BusinessEmail,
			&r.BusinessPhone,
			&r.Address,
			&r.City,
			&r.State,
			&r.ProductName,
			&r.Quantity,
			&r.Unit,
			&r.Price,
			&r.IsRFQActive,
			&r.CreatedAT,
			&r.UpdatedAT,
		)

		if err != nil {
			return nil, err
		}

		rfqs = append(rfqs, r)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return rfqs, nil
}

func (rs *PostgresRFQStore) GetRFQByBusinessID(id string) ([]RFQ, error) {
	query := `
	SELECT
		r.id,
		b.id,
		b.business_name,
		b.business_email,
		b.business_phone,
		b.address,
		b.city,
		b.state,
		r.product_name,
		r.quantity,
		r.unit,
		r.price,
		r.is_rfq_active,
		r.created_at,
		r.updated_at
	FROM rfqs r
	JOIN businesses b
		ON b.id = r.business_id
	WHERE r.business_id = $1;
	`
	res, err := rs.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rfqs []RFQ
	for res.Next() {
		var r RFQ
		err = res.Scan(
			&r.ID,
			&r.BusinessID,
			&r.BusinessName,
			&r.BusinessEmail,
			&r.BusinessPhone,
			&r.Address,
			&r.City,
			&r.State,
			&r.ProductName,
			&r.Quantity,
			&r.Unit,
			&r.Price,
			&r.IsRFQActive,
			&r.CreatedAT,
			&r.UpdatedAT,
		)

		if err != nil {
			return nil, err
		}

		rfqs = append(rfqs, r)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return rfqs, nil
}
