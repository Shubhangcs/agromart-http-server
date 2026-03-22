package store

import (
	"database/sql"

	"github.com/shubhangcs/agromart-server/internal/models"
)

type LeadStore interface {
	CreateLead(*models.Lead) error
	GetSentLeads(enquirerBusinessID string, limit, offset int) ([]models.LeadSentResponse, error)
	GetReceivedLeads(sellerBusinessID string, limit, offset int) ([]models.LeadReceivedResponse, error)
}

type PostgresLeadStore struct {
	db *sql.DB
}

func NewPostgresLeadStore(db *sql.DB) *PostgresLeadStore {
	return &PostgresLeadStore{db: db}
}

func (ls *PostgresLeadStore) CreateLead(l *models.Lead) error {
	query := `
	INSERT INTO leads (enquirer_business_id, enquire_to_id, product_id, enquiry_message, order_quantity, expected_price)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING lead_id, created_at
	`
	return ls.db.QueryRow(query,
		l.EnquirerBusinessID, l.EnquireToID, l.ProductID, l.EnquiryMessage, l.OrderQuantity, l.ExpectedPrice,
	).Scan(&l.LeadID, &l.CreatedAT)
}

// GetSentLeads returns leads submitted by the given business,
// enriched with the enquired product details and the complete seller business info.
func (ls *PostgresLeadStore) GetSentLeads(enquirerBusinessID string, limit, offset int) ([]models.LeadSentResponse, error) {
	query := `
	SELECT
		l.lead_id, l.product_id, p.name,
		(SELECT image FROM product_images WHERE product_id = p.id ORDER BY image_index ASC LIMIT 1),
		b.id, b.business_name, b.business_email, b.business_phone, b.business_profile_image,
		b.address, b.city, b.state, b.pincode,
		l.enquiry_message, l.order_quantity, l.expected_price, l.created_at
	FROM leads l
	JOIN products p   ON p.id  = l.product_id
	JOIN businesses b ON b.id  = l.enquire_to_id
	WHERE l.enquirer_business_id = $1
	ORDER BY l.created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := ls.db.Query(query, enquirerBusinessID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []models.LeadSentResponse
	for rows.Next() {
		var lead models.LeadSentResponse
		if err = rows.Scan(
			&lead.LeadID, &lead.ProductID, &lead.ProductName,
			&lead.ProductImage,
			&lead.SellerBusinessID, &lead.SellerBusinessName, &lead.SellerBusinessEmail, &lead.SellerBusinessPhone, &lead.SellerBusinessProfileImage,
			&lead.SellerAddress, &lead.SellerCity, &lead.SellerState, &lead.SellerPincode,
			&lead.EnquiryMessage, &lead.OrderQuantity, &lead.ExpectedPrice, &lead.CreatedAT,
		); err != nil {
			return nil, err
		}
		leads = append(leads, lead)
	}
	return leads, rows.Err()
}

// GetReceivedLeads returns leads received by the given business,
// enriched with the enquired product details and the complete enquirer business info.
func (ls *PostgresLeadStore) GetReceivedLeads(sellerBusinessID string, limit, offset int) ([]models.LeadReceivedResponse, error) {
	query := `
	SELECT
		l.lead_id, l.product_id, p.name,
		(SELECT image FROM product_images WHERE product_id = p.id ORDER BY image_index ASC LIMIT 1),
		eb.id, eb.business_name, eb.business_email, eb.business_phone, eb.business_profile_image,
		eb.address, eb.city, eb.state, eb.pincode,
		l.enquiry_message, l.order_quantity, l.expected_price, l.created_at
	FROM leads l
	JOIN products p    ON p.id  = l.product_id
	JOIN businesses eb ON eb.id = l.enquirer_business_id
	WHERE l.enquire_to_id = $1
	ORDER BY l.created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := ls.db.Query(query, sellerBusinessID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []models.LeadReceivedResponse
	for rows.Next() {
		var lead models.LeadReceivedResponse
		if err = rows.Scan(
			&lead.LeadID, &lead.ProductID, &lead.ProductName,
			&lead.ProductImage,
			&lead.EnquirerBusinessID, &lead.EnquirerBusinessName, &lead.EnquirerBusinessEmail, &lead.EnquirerBusinessPhone, &lead.EnquirerBusinessProfileImage,
			&lead.EnquirerAddress, &lead.EnquirerCity, &lead.EnquirerState, &lead.EnquirerPincode,
			&lead.EnquiryMessage, &lead.OrderQuantity, &lead.ExpectedPrice, &lead.CreatedAT,
		); err != nil {
			return nil, err
		}
		leads = append(leads, lead)
	}
	return leads, rows.Err()
}
