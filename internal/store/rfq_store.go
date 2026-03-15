package store

import (
	"database/sql"

	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

type PostgresRFQStore struct {
	db *sql.DB
}

type RFQStore interface {
	CreateRFQ(*models.RFQ) error
	ActivateRFQ(*models.RFQ) error
	UpdateRFQ(*models.RFQ) error
	DeleteRFQ(id string) error
	GetAllRFQ(filter utils.RFQFilter, limit, offset int) ([]models.RFQResponse, error)
	GetRFQByBusinessID(id string, limit, offset int) ([]models.RFQResponse, error)
}

func NewPostgresRFQStore(db *sql.DB) *PostgresRFQStore {
	return &PostgresRFQStore{db: db}
}

func (rs *PostgresRFQStore) CreateRFQ(r *models.RFQ) error {
	query := `
	INSERT INTO rfqs (business_id, category_id, sub_category_id, product_name, quantity, unit, price, is_rfq_active)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING id, created_at, updated_at
	`
	return rs.db.QueryRow(
		query,
		r.BusinessID, r.CategoryID, r.SubCategoryID,
		r.ProductName, r.Quantity, r.Unit, r.Price, r.IsRFQActive,
	).Scan(&r.ID, &r.CreatedAT, &r.UpdatedAT)
}

func (rs *PostgresRFQStore) ActivateRFQ(r *models.RFQ) error {
	res, err := rs.db.Exec(
		`UPDATE rfqs SET is_rfq_active = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
		r.IsRFQActive, r.ID,
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

func (rs *PostgresRFQStore) UpdateRFQ(r *models.RFQ) error {
	query := `
	UPDATE rfqs
	SET category_id     = COALESCE(NULLIF($1,''), category_id::text)::uuid,
	    sub_category_id = COALESCE(NULLIF($2,''), sub_category_id::text)::uuid,
	    product_name    = COALESCE(NULLIF($3,''), product_name),
	    quantity        = COALESCE(NULLIF($4,0),  quantity),
	    unit            = COALESCE(NULLIF($5,''), unit),
	    price           = COALESCE(NULLIF($6,0),  price),
	    updated_at      = CURRENT_TIMESTAMP
	WHERE id = $7
	`
	res, err := rs.db.Exec(query,
		r.CategoryID, r.SubCategoryID, r.ProductName,
		r.Quantity, r.Unit, r.Price, r.ID,
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

func (rs *PostgresRFQStore) DeleteRFQ(id string) error {
	res, err := rs.db.Exec(`DELETE FROM rfqs WHERE id = $1`, id)
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

// GetAllRFQ returns active RFQs with business and category details, paginated.
// Optional filters: product name search (?q=), city, and state.
func (rs *PostgresRFQStore) GetAllRFQ(filter utils.RFQFilter, limit, offset int) ([]models.RFQResponse, error) {
	query := `
	SELECT
		r.id, b.user_id, b.id, b.business_name, b.business_email, b.business_phone, b.address, b.city, b.state,
		c.id, c.name, c.description,
		sc.id, sc.name, sc.description,
		r.product_name, r.quantity, r.unit, r.price, r.is_rfq_active, r.created_at, r.updated_at
	FROM rfqs r
	JOIN businesses b ON b.id = r.business_id
	JOIN categories c ON c.id = r.category_id
	JOIN sub_categories sc ON sc.id = r.sub_category_id
	WHERE r.is_rfq_active = TRUE
	  AND ($1 = '' OR r.product_name ILIKE '%' || $1 || '%')
	  AND ($2 = '' OR b.city ILIKE $2)
	  AND ($3 = '' OR b.state ILIKE $3)
	ORDER BY r.created_at DESC
	LIMIT $4 OFFSET $5
	`
	rows, err := rs.db.Query(query, filter.ProductName, filter.City, filter.State, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rfqs []models.RFQResponse
	for rows.Next() {
		var rfq models.RFQResponse
		err = rows.Scan(
			&rfq.ID, &rfq.UserID, &rfq.BusinessID, &rfq.BusinessName, &rfq.BusinessEmail,
			&rfq.BusinessPhone, &rfq.Address, &rfq.City, &rfq.State,
			&rfq.CategoryID, &rfq.CategoryName, &rfq.CategoryDescription,
			&rfq.SubCategoryID, &rfq.SubCategoryName, &rfq.SubCategoryDescription,
			&rfq.ProductName, &rfq.Quantity, &rfq.Unit, &rfq.Price,
			&rfq.IsRFQActive, &rfq.CreatedAT, &rfq.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}
		rfqs = append(rfqs, rfq)
	}
	return rfqs, rows.Err()
}

// GetRFQByBusinessID returns all RFQs for a business (active and inactive), paginated.
func (rs *PostgresRFQStore) GetRFQByBusinessID(id string, limit, offset int) ([]models.RFQResponse, error) {
	query := `
	SELECT
		r.id, b.user_id, b.id, b.business_name, b.business_email, b.business_phone, b.address, b.city, b.state,
		r.product_name, r.quantity, r.unit, r.price, r.is_rfq_active, r.created_at, r.updated_at
	FROM rfqs r
	JOIN businesses b ON b.id = r.business_id
	WHERE r.business_id = $1
	ORDER BY r.created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := rs.db.Query(query, id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rfqs []models.RFQResponse
	for rows.Next() {
		var rfq models.RFQResponse
		err = rows.Scan(
			&rfq.ID, &rfq.UserID, &rfq.BusinessID, &rfq.BusinessName, &rfq.BusinessEmail,
			&rfq.BusinessPhone, &rfq.Address, &rfq.City, &rfq.State,
			&rfq.ProductName, &rfq.Quantity, &rfq.Unit, &rfq.Price,
			&rfq.IsRFQActive, &rfq.CreatedAT, &rfq.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}
		rfqs = append(rfqs, rfq)
	}
	return rfqs, rows.Err()
}
