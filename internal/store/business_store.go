package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/shubhangcs/agromart-server/internal/models"
)

type PostgresBusinessStore struct {
	db *sql.DB
}

func NewPostgresBusinessStore(db *sql.DB) *PostgresBusinessStore {
	return &PostgresBusinessStore{db: db}
}

type BusinessStore interface {
	CreateBusiness(*models.Business) error
	CreateSocial(*models.Social) error
	CreateLegal(*models.Legal) error
	CreateBusinessApplication(*models.BusinessApplication) error
	UpdateBusiness(*models.Business) error
	UpdateSocial(*models.Social) error
	UpdateLegal(*models.Legal) error
	DeleteBusiness(id string) error
	AcceptBusinessApplication(string) error
	RejectBusinessApplication(*models.BusinessApplication) error
	GetCompleteBusinessDetails(id string) (*models.BusinessDetails, error)
	GetBusiness(id string) (*models.Business, error)
	GetSocial(id string) (*models.Social, error)
	GetLegal(id string) (*models.Legal, error)
	GetBusinessApplication(id string) (*models.BusinessApplication, error)
	GetAllBusinesses(limit, offset int) ([]models.BusinessDetails, error)
	UpdateVerifyBusinessStatus(id string, status bool) error
	UpdateTrustBusinessStatus(id string, status bool) error
	UpdateBlockBusinessStatus(id string, status bool) error
	GetBusinessIDByUserID(id string) (*string, error)
	IsBusinessApproved(id string) (bool, error)
}

func (bs *PostgresBusinessStore) CreateBusiness(b *models.Business) error {
	query := `
	INSERT INTO businesses (
		user_id, business_profile_image, business_name, business_email,
		business_phone, address, city, state, pincode, business_type
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	RETURNING id, created_at, updated_at
	`
	return bs.db.QueryRow(
		query,
		b.UserID, b.ProfileImage, b.Name, b.Email,
		b.Phone, b.Address, b.City, b.State, b.Pincode, b.BusinessType,
	).Scan(&b.ID, &b.CreatedAT, &b.UpdatedAT)
}

func (bs *PostgresBusinessStore) CreateSocial(s *models.Social) error {
	query := `
	INSERT INTO business_socials (business_id, linkedin, instagram, youtube, x, telegram, facebook, website)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`
	res, err := bs.db.Exec(query, s.ID, s.Linkedin, s.Instagram, s.Youtube, s.X, s.Telegram, s.Facebook, s.Website)
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

func (bs *PostgresBusinessStore) CreateLegal(l *models.Legal) error {
	query := `
	INSERT INTO business_legals (business_id, aadhaar, pan, export_import, msme, fassi, gst)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	`
	res, err := bs.db.Exec(query, l.ID, l.Aadhaar, l.Pan, l.ExportImport, l.MSME, l.Fassi, l.GST)
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

func (bs *PostgresBusinessStore) CreateBusinessApplication(ba *models.BusinessApplication) error {
	query := `INSERT INTO business_applications (business_id, status) VALUES ($1,$2)`
	res, err := bs.db.Exec(query, ba.ID, ba.Status)
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

func (bs *PostgresBusinessStore) UpdateBusiness(b *models.Business) error {
	query := `
	UPDATE businesses
	SET business_name  = COALESCE(NULLIF($1,''), business_name),
	    business_email = COALESCE(NULLIF($2,''), business_email),
	    business_phone = COALESCE(NULLIF($3,''), business_phone),
	    address        = COALESCE(NULLIF($4,''), address),
	    city           = COALESCE(NULLIF($5,''), city),
	    state          = COALESCE(NULLIF($6,''), state),
	    pincode        = COALESCE(NULLIF($7,''), pincode),
	    business_type  = COALESCE(NULLIF($8,''), business_type),
	    updated_at     = CURRENT_TIMESTAMP
	WHERE id = $9
	`
	res, err := bs.db.Exec(query, b.Name, b.Email, b.Phone, b.Address, b.City, b.State, b.Pincode, b.BusinessType, b.ID)
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

func (bs *PostgresBusinessStore) UpdateSocial(s *models.Social) error {
	query := `
	UPDATE business_socials
	SET linkedin   = COALESCE($1, linkedin),
	    instagram  = COALESCE($2, instagram),
	    youtube    = COALESCE($3, youtube),
	    telegram   = COALESCE($4, telegram),
	    x          = COALESCE($5, x),
	    facebook   = COALESCE($6, facebook),
	    website    = COALESCE($7, website),
	    updated_at = CURRENT_TIMESTAMP
	WHERE business_id = $8
	`
	res, err := bs.db.Exec(query, s.Linkedin, s.Instagram, s.Youtube, s.Telegram, s.X, s.Facebook, s.Website, s.ID)
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

func (bs *PostgresBusinessStore) UpdateLegal(l *models.Legal) error {
	query := `
	UPDATE business_legals
	SET aadhaar       = COALESCE($1, aadhaar),
	    pan           = COALESCE($2, pan),
	    export_import = COALESCE($3, export_import),
	    msme          = COALESCE($4, msme),
	    fassi         = COALESCE($5, fassi),
	    gst           = COALESCE($6, gst),
	    updated_at    = CURRENT_TIMESTAMP
	WHERE business_id = $7
	`
	res, err := bs.db.Exec(query, l.Aadhaar, l.Pan, l.ExportImport, l.MSME, l.Fassi, l.GST, l.ID)
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

func (bs *PostgresBusinessStore) DeleteBusiness(id string) error {
	res, err := bs.db.Exec(`DELETE FROM businesses WHERE id = $1`, id)
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

func (bs *PostgresBusinessStore) AcceptBusinessApplication(id string) error {
	tx, err := bs.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`UPDATE business_applications SET status = 'ACCEPTED' WHERE business_id = $1`, id,
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

	res, err = tx.Exec(
		`UPDATE businesses SET is_business_approved = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, id,
	)
	if err != nil {
		return err
	}
	rows, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

func (bs *PostgresBusinessStore) RejectBusinessApplication(ba *models.BusinessApplication) error {
	query := `
	UPDATE business_applications
	SET status = 'REJECTED', reject_reason = $1
	WHERE business_id = $2
	`
	res, err := bs.db.Exec(query, ba.RejectReason, ba.ID)
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

// GetCompleteBusinessDetails fetches all business sub-records with a single
// four-way LEFT JOIN so missing social/legal/application rows don't error.
func (bs *PostgresBusinessStore) GetCompleteBusinessDetails(id string) (*models.BusinessDetails, error) {
	query := `
	SELECT
		b.id, b.business_profile_image, b.business_name, b.business_email,
		b.business_phone, b.address, b.city, b.state, b.pincode, b.business_type,
		b.is_business_verified, b.is_business_approved, b.is_business_trusted,
		b.created_at, b.updated_at,
		s.business_id, s.linkedin, s.instagram, s.youtube, s.telegram, s.x,
		s.facebook, s.website, s.created_at, s.updated_at,
		l.business_id, l.aadhaar, l.pan, l.export_import, l.msme, l.fassi,
		l.gst, l.created_at, l.updated_at,
		a.business_id, a.status, a.reject_reason, a.created_at
	FROM businesses b
	LEFT JOIN business_socials s ON s.business_id = b.id
	LEFT JOIN business_legals l ON l.business_id = b.id
	LEFT JOIN business_applications a ON a.business_id = b.id
	WHERE b.id = $1
	`
	var bd models.BusinessDetails
	var (
		socialID        *string
		legalID         *string
		applicationID   *string
		appStatus       *string
		socialCreatedAt *time.Time
		socialUpdatedAt *time.Time
		legalCreatedAt  *time.Time
		legalUpdatedAt  *time.Time
		appCreatedAt    *time.Time
	)

	err := bs.db.QueryRow(query, id).Scan(
		&bd.CoreBusinessDetails.ID,
		&bd.CoreBusinessDetails.ProfileImage,
		&bd.CoreBusinessDetails.Name,
		&bd.CoreBusinessDetails.Email,
		&bd.CoreBusinessDetails.Phone,
		&bd.CoreBusinessDetails.Address,
		&bd.CoreBusinessDetails.City,
		&bd.CoreBusinessDetails.State,
		&bd.CoreBusinessDetails.Pincode,
		&bd.CoreBusinessDetails.BusinessType,
		&bd.CoreBusinessDetails.IsBusinessVerified,
		&bd.CoreBusinessDetails.IsBusinessApproved,
		&bd.CoreBusinessDetails.IsBusinessTrusted,
		&bd.CoreBusinessDetails.CreatedAT,
		&bd.CoreBusinessDetails.UpdatedAT,
		&socialID,
		&bd.BusinessSocialDetails.Linkedin,
		&bd.BusinessSocialDetails.Instagram,
		&bd.BusinessSocialDetails.Youtube,
		&bd.BusinessSocialDetails.Telegram,
		&bd.BusinessSocialDetails.X,
		&bd.BusinessSocialDetails.Facebook,
		&bd.BusinessSocialDetails.Website,
		&socialCreatedAt,
		&socialUpdatedAt,
		&legalID,
		&bd.BusinessLegalDetails.Aadhaar,
		&bd.BusinessLegalDetails.Pan,
		&bd.BusinessLegalDetails.ExportImport,
		&bd.BusinessLegalDetails.MSME,
		&bd.BusinessLegalDetails.Fassi,
		&bd.BusinessLegalDetails.GST,
		&legalCreatedAt,
		&legalUpdatedAt,
		&applicationID,
		&appStatus,
		&bd.BusinessApplicationDetails.RejectReason,
		&appCreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if socialID != nil {
		bd.BusinessSocialDetails.ID = *socialID
	}
	if socialCreatedAt != nil {
		bd.BusinessSocialDetails.CreatedAT = *socialCreatedAt
	}
	if socialUpdatedAt != nil {
		bd.BusinessSocialDetails.UpdatedAT = *socialUpdatedAt
	}
	if legalID != nil {
		bd.BusinessLegalDetails.ID = *legalID
	}
	if legalCreatedAt != nil {
		bd.BusinessLegalDetails.CreatedAT = *legalCreatedAt
	}
	if legalUpdatedAt != nil {
		bd.BusinessLegalDetails.UpdatedAT = *legalUpdatedAt
	}
	if applicationID != nil {
		bd.BusinessApplicationDetails.ID = *applicationID
	}
	if appStatus != nil {
		bd.BusinessApplicationDetails.Status = *appStatus
	}
	if appCreatedAt != nil {
		bd.BusinessApplicationDetails.CreatedAT = *appCreatedAt
	}

	return &bd, nil
}

func (bs *PostgresBusinessStore) GetBusiness(id string) (*models.Business, error) {
	query := `
	SELECT id, business_profile_image, business_name, business_email, business_phone,
	       address, city, state, pincode, business_type,
	       is_business_verified, is_business_approved, is_business_trusted,
	       created_at, updated_at
	FROM businesses
	WHERE id = $1
	`
	var b models.Business
	err := bs.db.QueryRow(query, id).Scan(
		&b.ID, &b.ProfileImage, &b.Name, &b.Email, &b.Phone,
		&b.Address, &b.City, &b.State, &b.Pincode, &b.BusinessType,
		&b.IsBusinessVerified, &b.IsBusinessApproved, &b.IsBusinessTrusted,
		&b.CreatedAT, &b.UpdatedAT,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (bs *PostgresBusinessStore) GetSocial(id string) (*models.Social, error) {
	query := `
	SELECT business_id, linkedin, instagram, youtube, telegram, x, facebook, website, created_at, updated_at
	FROM business_socials
	WHERE business_id = $1
	`
	var s models.Social
	err := bs.db.QueryRow(query, id).Scan(
		&s.ID, &s.Linkedin, &s.Instagram, &s.Youtube, &s.Telegram,
		&s.X, &s.Facebook, &s.Website, &s.CreatedAT, &s.UpdatedAT,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (bs *PostgresBusinessStore) GetLegal(id string) (*models.Legal, error) {
	query := `
	SELECT business_id, aadhaar, pan, export_import, msme, fassi, gst, created_at, updated_at
	FROM business_legals
	WHERE business_id = $1
	`
	var l models.Legal
	err := bs.db.QueryRow(query, id).Scan(
		&l.ID, &l.Aadhaar, &l.Pan, &l.ExportImport, &l.MSME,
		&l.Fassi, &l.GST, &l.CreatedAT, &l.UpdatedAT,
	)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (bs *PostgresBusinessStore) GetBusinessApplication(id string) (*models.BusinessApplication, error) {
	query := `
	SELECT business_id, status, reject_reason, created_at
	FROM business_applications
	WHERE business_id = $1
	`
	var ba models.BusinessApplication
	err := bs.db.QueryRow(query, id).Scan(&ba.ID, &ba.Status, &ba.RejectReason, &ba.CreatedAT)
	if err != nil {
		return nil, err
	}
	return &ba, nil
}

// GetAllBusinesses returns a paginated list of businesses with their
// social, legal and application details via a single LEFT JOIN query.
func (bs *PostgresBusinessStore) GetAllBusinesses(limit, offset int) ([]models.BusinessDetails, error) {
	query := `
	SELECT
		b.id, b.user_id, b.business_profile_image, b.business_name, b.business_email,
		b.business_phone, b.address, b.city, b.state, b.pincode, b.business_type,
		b.is_business_verified, b.is_business_trusted, b.is_business_approved,
		b.created_at, b.updated_at,
		s.business_id, s.linkedin, s.instagram, s.youtube, s.telegram, s.x,
		s.facebook, s.website, s.created_at, s.updated_at,
		l.business_id, l.aadhaar, l.pan, l.export_import, l.msme, l.fassi,
		l.gst, l.created_at, l.updated_at,
		a.business_id, a.status, a.reject_reason, a.created_at
	FROM businesses b
	LEFT JOIN business_socials s ON s.business_id = b.id
	LEFT JOIN business_legals l ON l.business_id = b.id
	LEFT JOIN business_applications a ON a.business_id = b.id
	ORDER BY b.created_at DESC
	LIMIT $1 OFFSET $2
	`
	rows, err := bs.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var businesses []models.BusinessDetails
	for rows.Next() {
		var bd models.BusinessDetails
		var (
			socialID        *string
			legalID         *string
			applicationID   *string
			appStatus       *string
			socialCreatedAt *time.Time
			socialUpdatedAt *time.Time
			legalCreatedAt  *time.Time
			legalUpdatedAt  *time.Time
			appCreatedAt    *time.Time
		)

		err = rows.Scan(
			&bd.CoreBusinessDetails.ID,
			&bd.CoreBusinessDetails.UserID,
			&bd.CoreBusinessDetails.ProfileImage,
			&bd.CoreBusinessDetails.Name,
			&bd.CoreBusinessDetails.Email,
			&bd.CoreBusinessDetails.Phone,
			&bd.CoreBusinessDetails.Address,
			&bd.CoreBusinessDetails.City,
			&bd.CoreBusinessDetails.State,
			&bd.CoreBusinessDetails.Pincode,
			&bd.CoreBusinessDetails.BusinessType,
			&bd.CoreBusinessDetails.IsBusinessVerified,
			&bd.CoreBusinessDetails.IsBusinessTrusted,
			&bd.CoreBusinessDetails.IsBusinessApproved,
			&bd.CoreBusinessDetails.CreatedAT,
			&bd.CoreBusinessDetails.UpdatedAT,
			&socialID,
			&bd.BusinessSocialDetails.Linkedin,
			&bd.BusinessSocialDetails.Instagram,
			&bd.BusinessSocialDetails.Youtube,
			&bd.BusinessSocialDetails.Telegram,
			&bd.BusinessSocialDetails.X,
			&bd.BusinessSocialDetails.Facebook,
			&bd.BusinessSocialDetails.Website,
			&socialCreatedAt,
			&socialUpdatedAt,
			&legalID,
			&bd.BusinessLegalDetails.Aadhaar,
			&bd.BusinessLegalDetails.Pan,
			&bd.BusinessLegalDetails.ExportImport,
			&bd.BusinessLegalDetails.MSME,
			&bd.BusinessLegalDetails.Fassi,
			&bd.BusinessLegalDetails.GST,
			&legalCreatedAt,
			&legalUpdatedAt,
			&applicationID,
			&appStatus,
			&bd.BusinessApplicationDetails.RejectReason,
			&appCreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if socialID != nil {
			bd.BusinessSocialDetails.ID = *socialID
		}
		if socialCreatedAt != nil {
			bd.BusinessSocialDetails.CreatedAT = *socialCreatedAt
		}
		if socialUpdatedAt != nil {
			bd.BusinessSocialDetails.UpdatedAT = *socialUpdatedAt
		}
		if legalID != nil {
			bd.BusinessLegalDetails.ID = *legalID
		}
		if legalCreatedAt != nil {
			bd.BusinessLegalDetails.CreatedAT = *legalCreatedAt
		}
		if legalUpdatedAt != nil {
			bd.BusinessLegalDetails.UpdatedAT = *legalUpdatedAt
		}
		if applicationID != nil {
			bd.BusinessApplicationDetails.ID = *applicationID
		}
		if appStatus != nil {
			bd.BusinessApplicationDetails.Status = *appStatus
		}
		if appCreatedAt != nil {
			bd.BusinessApplicationDetails.CreatedAT = *appCreatedAt
		}

		businesses = append(businesses, bd)
	}
	return businesses, rows.Err()
}

func (bs *PostgresBusinessStore) UpdateVerifyBusinessStatus(id string, status bool) error {
	res, err := bs.db.Exec(
		`UPDATE businesses SET is_business_verified = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
		status, id,
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

func (bs *PostgresBusinessStore) UpdateTrustBusinessStatus(id string, status bool) error {
	res, err := bs.db.Exec(
		`UPDATE businesses SET is_business_trusted = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
		status, id,
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

func (bs *PostgresBusinessStore) UpdateBlockBusinessStatus(id string, status bool) error {
	res, err := bs.db.Exec(
		`UPDATE businesses SET is_business_approved = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
		status, id,
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

func (bs *PostgresBusinessStore) GetBusinessIDByUserID(id string) (*string, error) {
	var businessID *string
	err := bs.db.QueryRow(`SELECT id FROM businesses WHERE user_id = $1`, id).Scan(&businessID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return businessID, nil
}

func (bs *PostgresBusinessStore) IsBusinessApproved(id string) (bool, error) {
	var approved bool
	err := bs.db.QueryRow(`SELECT is_business_approved FROM businesses WHERE id = $1`, id).Scan(&approved)
	if err != nil {
		return false, err
	}
	return approved, nil
}

