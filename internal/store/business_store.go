package store

import (
	"database/sql"
	"errors"
	"time"
)

type Business struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	ProfileImage       *string   `json:"profile_image"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Phone              string    `json:"phone"`
	Address            string    `json:"address"`
	City               string    `json:"city"`
	State              string    `json:"state"`
	Pincode            string    `json:"pincode"`
	BusinessType       string    `json:"business_type"`
	IsBusinessVerified bool      `json:"is_business_verified"`
	IsBusinessTrusted  bool      `json:"is_business_trusted"`
	IsBusinessApproved bool      `json:"is_business_approved"`
	CreatedAT          time.Time `json:"created_at"`
	UpdatedAT          time.Time `json:"updated_at"`
}

type Social struct {
	ID        string    `json:"id"`
	Linkedin  *string   `json:"linkedin"`
	Instagram *string   `json:"instagram"`
	Telegram  *string   `json:"telegram"`
	Youtube   *string   `json:"youtube"`
	X         *string   `json:"x"`
	Facebook  *string   `json:"facebook"`
	Website   *string   `json:"website"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

type Legal struct {
	ID           string    `json:"id"`
	Aadhaar      *string   `json:"aadhaar"`
	Pan          *string   `json:"pan"`
	ExportImport *string   `json:"export_import"`
	MSME         *string   `json:"msme"`
	Fassi        *string   `json:"fassi"`
	GST          *string   `json:"gst"`
	CreatedAT    time.Time `json:"created_at"`
	UpdatedAT    time.Time `json:"updated_at"`
}

type BusinessApplication struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	RejectReason *string   `json:"reject_reason"`
	CreatedAT    time.Time `json:"created_at"`
}

type BusinessDetails struct {
	CoreBusinessDetails        Business            `json:"business_details"`
	BusinessSocialDetails      Social              `json:"social_details"`
	BusinessLegalDetails       Legal               `json:"legal_details"`
	BusinessApplicationDetails BusinessApplication `json:"business_application"`
}

type BusinessRating struct {
	ID         string    `json:"id"`
	BusinessID string    `json:"busness_id"`
	UserID     string    `json:"user_id"`
	UserName   string    `json:"user_name,omitempty"`
	Rating     float64   `json:"rating"`
	CreatedAT  time.Time `json:"created_at"`
	UpdatedAT  time.Time `json:"updated_at"`
}

type PostgresBusinessStore struct {
	db *sql.DB
}

func NewPostgresBusinessStore(db *sql.DB) *PostgresBusinessStore {
	return &PostgresBusinessStore{
		db: db,
	}
}

type BusinessStore interface {
	CreateBusiness(*Business) error
	CreateSocial(*Social) error
	CreateLegal(*Legal) error
	CreateBusinessApplication(*BusinessApplication) error
	UpdateBusiness(*Business) error
	UpdateSocial(*Social) error
	UpdateLegal(*Legal) error
	DeleteBusiness(id string) error
	AcceptBusinessApplication(string) error
	RejectBusinessApplication(*BusinessApplication) error
	GetCompleteBusinessDetails(id string) (*BusinessDetails, error)
	GetBusiness(id string) (*Business, error)
	GetSocial(id string) (*Social, error)
	GetLegal(id string) (*Legal, error)
	GetBusinessApplication(id string) (*BusinessApplication, error)
	GetAllBusinesses() ([]BusinessDetails, error)
	UpdateVerifyBusinessStatus(id string, status bool) error
	UpdateTrustBusinessStatus(id string, status bool) error
	UpdateBlockBusinessStatus(id string, status bool) error
	GetBusinessIDByUserID(id string) (*string, error)
	IsBusinessApproved(id string) (bool, error)
	RateBusiness(*BusinessRating) error
	GetAvrageBusinessRating(id string) (float64, error)
	GetRatingsByBusinessID(id string) ([]BusinessRating, error)
}

func (bs *PostgresBusinessStore) CreateBusiness(b *Business) error {
	query := `
	INSERT INTO businesses (
		user_id,
		business_profile_image,
		business_name,
		business_email,
		business_phone,
		address,
		city,
		state,
		pincode,
		business_type
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	)
	RETURNING id, created_at, updated_at;
	`
	err := bs.db.QueryRow(
		query,
		b.UserID,
		b.ProfileImage,
		b.Name,
		b.Email,
		b.Phone,
		b.Address,
		b.City,
		b.State,
		b.Pincode,
		b.BusinessType,
	).Scan(
		&b.ID,
		&b.CreatedAT,
		&b.UpdatedAT,
	)

	if err != nil {
		return err
	}

	return nil
}

func (bs *PostgresBusinessStore) CreateSocial(s *Social) error {
	query := `
	INSERT INTO business_socials (
		business_id,
		linkedin,
		instagram,
		youtube,
		x,
		telegram,
		facebook,
		website
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8
	);
	`

	res, err := bs.db.Exec(
		query,
		s.ID,
		s.Linkedin,
		s.Instagram,
		s.Youtube,
		s.X,
		s.Telegram,
		s.Facebook,
		s.Website,
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

func (bs *PostgresBusinessStore) CreateLegal(l *Legal) error {
	query := `
	INSERT INTO business_legals (
		business_id,
		aadhaar,
		pan,
		export_import,
		msme,
		fassi,
		gst
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	);
	`

	res, err := bs.db.Exec(
		query,
		l.ID,
		l.Aadhaar,
		l.Pan,
		l.ExportImport,
		l.MSME,
		l.Fassi,
		l.GST,
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

func (bs *PostgresBusinessStore) CreateBusinessApplication(ba *BusinessApplication) error {
	query := `
	INSERT INTO business_applications (
		business_id,
		status
	) VALUES (
		$1, $2 
	);
	`

	res, err := bs.db.Exec(query, ba.ID, ba.Status)
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

func (bs *PostgresBusinessStore) UpdateBusiness(b *Business) error {
	query := `
	UPDATE businesses
	SET business_name = COALESCE($1, business_name),
	business_email = COALESCE($2, business_email),
	business_phone = COALESCE($3, business_phone),
	address = COALESCE($4, address),
	city = COALESCE($5, city),
	state = COALESCE($6, state),
	pincode = COALESCE($7, pincode),
	business_type = COALESCE($8, business_type),
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $9
	`

	res, err := bs.db.Exec(
		query,
		b.Name,
		b.Email,
		b.Phone,
		b.Address,
		b.City,
		b.State,
		b.Pincode,
		b.BusinessType,
		b.ID,
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

func (bs *PostgresBusinessStore) UpdateSocial(s *Social) error {
	// FIX: Added missing comma after `x = COALESCE($5, x)`
	query := `
	UPDATE business_socials
	SET linkedin = COALESCE($1, linkedin),
	instagram = COALESCE($2, instagram),
	youtube = COALESCE($3, youtube),
	telegram = COALESCE($4, telegram),
	x = COALESCE($5, x),
	facebook = COALESCE($6, facebook),
	website = COALESCE($7, website),
	updated_at = CURRENT_TIMESTAMP
	WHERE business_id = $8;
	`

	res, err := bs.db.Exec(
		query,
		s.Linkedin,
		s.Instagram,
		s.Youtube,
		s.Telegram,
		s.X,
		s.Facebook,
		s.Website,
		s.ID,
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

func (bs *PostgresBusinessStore) UpdateLegal(l *Legal) error {
	// FIX: Added missing comma after `gst = COALESCE($6, gst)`
	query := `
	UPDATE business_legals
	SET aadhaar = COALESCE($1, aadhaar),
	pan = COALESCE($2, pan),
	export_import = COALESCE($3, export_import),
	msme = COALESCE($4, msme),
	fassi = COALESCE($5, fassi),
	gst = COALESCE($6, gst),
	updated_at = CURRENT_TIMESTAMP
	WHERE business_id = $7;
	`

	res, err := bs.db.Exec(
		query,
		l.Aadhaar,
		l.Pan,
		l.ExportImport,
		l.MSME,
		l.Fassi,
		l.GST,
		l.ID,
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

func (bs *PostgresBusinessStore) DeleteBusiness(id string) error {
	query := `
	DELETE FROM businesses
	WHERE id = $1;
	`

	res, err := bs.db.Exec(query, id)
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

func (bs *PostgresBusinessStore) AcceptBusinessApplication(id string) error {
	// FIX: Table name corrected from `business_application` to `business_applications`
	query1 := `
	UPDATE business_applications
	SET status = 'ACCEPTED'
	WHERE business_id = $1;
	`

	query2 := `
	UPDATE businesses
	SET is_business_approved = TRUE
	WHERE id = $1
	`

	tx, err := bs.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(query1, id)
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

	res, err = tx.Exec(query2, id)
	if err != nil {
		return err
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// FIX: Missing tx.Commit() — without this the transaction always rolls back
	return tx.Commit()
}

func (bs *PostgresBusinessStore) RejectBusinessApplication(ba *BusinessApplication) error {
	// FIX: Table name corrected from `business_application` to `business_applications`
	query := `
	UPDATE business_applications
	SET status = 'REJECTED',
	reject_reason = $1
	WHERE business_id = $2;
	`

	res, err := bs.db.Exec(query, ba.RejectReason, ba.ID)
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

func (bs *PostgresBusinessStore) GetCompleteBusinessDetails(id string) (*BusinessDetails, error) {
	businessesQuery := `
	SELECT 
		id,
		business_profile_image,
		business_name,
		business_email,
		business_phone,
		address,
		city,
		state,
		pincode,
		business_type,
		is_business_verified,
		is_business_approved,
		is_business_trusted,
		created_at,
		updated_at
	FROM businesses
	WHERE id = $1;
	`

	socialsQuery := `
	SELECT 
		business_id,
		linkedin,
		instagram,
		youtube,
		telegram,
		x,
		facebook,
		website,
		created_at,
		updated_at
	FROM business_socials
	WHERE business_id = $1;
	`

	legalQuery := `
	SELECT 
		business_id,
		aadhaar,
		pan,
		export_import,
		msme,
		fassi,
		gst,
		created_at,
		updated_at
	FROM business_legals
	WHERE business_id = $1;
	`

	// FIX: Table name corrected from `business_application` to `business_applications`
	businessApplicationQuery := `
	SELECT 
		business_id,
		status,
		reject_reason,
		created_at
	FROM business_applications
	WHERE business_id = $1
	`
	var businessDetails BusinessDetails
	err := bs.db.QueryRow(
		businessesQuery,
		id,
	).Scan(
		&businessDetails.CoreBusinessDetails.ID,
		&businessDetails.CoreBusinessDetails.ProfileImage,
		&businessDetails.CoreBusinessDetails.Name,
		&businessDetails.CoreBusinessDetails.Email,
		&businessDetails.CoreBusinessDetails.Phone,
		&businessDetails.CoreBusinessDetails.Address,
		&businessDetails.CoreBusinessDetails.City,
		&businessDetails.CoreBusinessDetails.State,
		&businessDetails.CoreBusinessDetails.Pincode,
		&businessDetails.CoreBusinessDetails.BusinessType,
		&businessDetails.CoreBusinessDetails.IsBusinessVerified,
		&businessDetails.CoreBusinessDetails.IsBusinessApproved,
		&businessDetails.CoreBusinessDetails.IsBusinessTrusted,
		&businessDetails.CoreBusinessDetails.CreatedAT,
		&businessDetails.CoreBusinessDetails.UpdatedAT,
	)

	if err != nil {
		return nil, err
	}

	err = bs.db.QueryRow(
		socialsQuery,
		id,
	).Scan(
		&businessDetails.BusinessSocialDetails.ID,
		&businessDetails.BusinessSocialDetails.Linkedin,
		&businessDetails.BusinessSocialDetails.Instagram,
		&businessDetails.BusinessSocialDetails.Youtube,
		&businessDetails.BusinessSocialDetails.Telegram,
		&businessDetails.BusinessSocialDetails.X,
		&businessDetails.BusinessSocialDetails.Facebook,
		&businessDetails.BusinessSocialDetails.Website,
		&businessDetails.BusinessSocialDetails.CreatedAT,
		&businessDetails.BusinessSocialDetails.UpdatedAT,
	)

	if err != nil {
		return nil, err
	}

	err = bs.db.QueryRow(
		legalQuery,
		id,
	).Scan(
		&businessDetails.BusinessLegalDetails.ID,
		&businessDetails.BusinessLegalDetails.Aadhaar,
		&businessDetails.BusinessLegalDetails.Pan,
		&businessDetails.BusinessLegalDetails.ExportImport,
		&businessDetails.BusinessLegalDetails.MSME,
		&businessDetails.BusinessLegalDetails.Fassi,
		&businessDetails.BusinessLegalDetails.GST,
		&businessDetails.BusinessLegalDetails.CreatedAT,
		&businessDetails.BusinessLegalDetails.UpdatedAT,
	)

	if err != nil {
		return nil, err
	}

	err = bs.db.QueryRow(
		businessApplicationQuery,
		id,
	).Scan(
		&businessDetails.BusinessApplicationDetails.ID,
		&businessDetails.BusinessApplicationDetails.Status,
		&businessDetails.BusinessApplicationDetails.RejectReason,
		&businessDetails.BusinessApplicationDetails.CreatedAT,
	)

	if err != nil {
		return nil, err
	}

	return &businessDetails, nil
}

func (bs *PostgresBusinessStore) GetBusiness(id string) (*Business, error) {
	businessesQuery := `
	SELECT 
		id,
		business_profile_image,
		business_name,
		business_email,
		business_phone,
		address,
		city,
		state,
		pincode,
		business_type,
		is_business_verified,
		is_business_approved,
		is_business_trusted,
		created_at,
		updated_at
	FROM businesses
	WHERE id = $1;
	`
	var businessDetails Business
	err := bs.db.QueryRow(
		businessesQuery,
		id,
	).Scan(
		&businessDetails.ID,
		&businessDetails.ProfileImage,
		&businessDetails.Name,
		&businessDetails.Email,
		&businessDetails.Phone,
		&businessDetails.Address,
		&businessDetails.City,
		&businessDetails.State,
		&businessDetails.Pincode,
		&businessDetails.BusinessType,
		&businessDetails.IsBusinessVerified,
		&businessDetails.IsBusinessApproved,
		&businessDetails.IsBusinessTrusted,
		&businessDetails.CreatedAT,
		&businessDetails.UpdatedAT,
	)

	if err != nil {
		return nil, err
	}

	return &businessDetails, nil
}

func (bs *PostgresBusinessStore) GetSocial(id string) (*Social, error) {
	socialsQuery := `
	SELECT 
		business_id,
		linkedin,
		instagram,
		youtube,
		telegram,
		x,
		facebook,
		website,
		created_at,
		updated_at
	FROM business_socials
	WHERE business_id = $1;
	`
	var businessDetails Social
	err := bs.db.QueryRow(
		socialsQuery,
		id,
	).Scan(
		&businessDetails.ID,
		&businessDetails.Linkedin,
		&businessDetails.Instagram,
		&businessDetails.Youtube,
		&businessDetails.Telegram,
		&businessDetails.X,
		&businessDetails.Facebook,
		&businessDetails.Website,
		&businessDetails.CreatedAT,
		&businessDetails.UpdatedAT,
	)

	if err != nil {
		return nil, err
	}

	return &businessDetails, nil
}

func (bs *PostgresBusinessStore) GetLegal(id string) (*Legal, error) {
	legalQuery := `
	SELECT 
		business_id,
		aadhaar,
		pan,
		export_import,
		msme,
		fassi,
		gst,
		created_at,
		updated_at
	FROM business_legals
	WHERE business_id = $1;
	`
	var businessDetails Legal
	err := bs.db.QueryRow(
		legalQuery,
		id,
	).Scan(
		&businessDetails.ID,
		&businessDetails.Aadhaar,
		&businessDetails.Pan,
		&businessDetails.ExportImport,
		&businessDetails.MSME,
		&businessDetails.Fassi,
		&businessDetails.GST,
		&businessDetails.CreatedAT,
		&businessDetails.UpdatedAT,
	)

	if err != nil {
		return nil, err
	}

	return &businessDetails, nil
}

func (bs *PostgresBusinessStore) GetBusinessApplication(id string) (*BusinessApplication, error) {
	// FIX: Table name corrected from `business_application` to `business_applications`
	businessApplicationQuery := `
	SELECT 
		business_id,
		status,
		reject_reason,
		created_at
	FROM business_applications
	WHERE business_id = $1
	`
	var businessDetails BusinessApplication
	err := bs.db.QueryRow(
		businessApplicationQuery,
		id,
	).Scan(
		&businessDetails.ID,
		&businessDetails.Status,
		&businessDetails.RejectReason,
		&businessDetails.CreatedAT,
	)
	if err != nil {
		return nil, err
	}

	return &businessDetails, nil
}

func (bs *PostgresBusinessStore) GetAllBusinesses() ([]BusinessDetails, error) {
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
		ORDER BY b.created_at DESC;
	`

	rows, err := bs.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var businesses []BusinessDetails

	for rows.Next() {
		var bd BusinessDetails
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

		err := rows.Scan(
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return businesses, nil
}

func (bs *PostgresBusinessStore) UpdateVerifyBusinessStatus(id string, status bool) error {
	query := `
	UPDATE businesses
	SET is_business_verified = $1,
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $2
	`

	res, err := bs.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	// FIX: Added missing error check after RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (bs *PostgresBusinessStore) UpdateTrustBusinessStatus(id string, status bool) error {
	query := `
	UPDATE businesses
	SET is_business_trusted = $1,
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $2;
	`

	res, err := bs.db.Exec(query, status, id)
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

func (bs *PostgresBusinessStore) UpdateBlockBusinessStatus(id string, status bool) error {
	query := `
	UPDATE businesses
	SET is_business_approved = $1,
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $2;
	`

	res, err := bs.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// FIX: Was returning `err` (which is nil here) instead of sql.ErrNoRows
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (bs *PostgresBusinessStore) GetBusinessIDByUserID(id string) (*string, error) {
	query := `
	SELECT 
		id
	FROM businesses
	WHERE user_id = $1;
	`

	var businessId *string
	err := bs.db.QueryRow(query, id).Scan(
		&businessId,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return businessId, nil
}

func (bs *PostgresBusinessStore) IsBusinessApproved(id string) (bool, error) {
	query := `
	SELECT is_business_approved
	FROM businesses
	WHERE id = $1;
	`
	var isBusinessApproved bool
	err := bs.db.QueryRow(query, id).Scan(
		&isBusinessApproved,
	)

	if err != nil {
		return false, err
	}

	return isBusinessApproved, nil
}

func (bs *PostgresBusinessStore) RateBusiness(r *BusinessRating) error {
	query := `
	INSERT INTO business_ratings (
		business_id,
		user_id,
		rating
	) VALUES (
		$1, $2, $3 
	);
	`

	res, err := bs.db.Exec(query, r.BusinessID, r.UserID, r.Rating)
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

func (bs *PostgresBusinessStore) GetAvrageBusinessRating(id string) (float64, error) {
	query := `
	SELECT
		AVG(rating)::NUMERIC(1,1)
	FROM business_ratings
	WHERE business_id = $1;
	`
	var businessRating float64
	err := bs.db.QueryRow(
		query,
		id,
	).Scan(
		&businessRating,
	)

	if err != nil {
		return 0, err
	}

	return businessRating, nil
}

func (bs *PostgresBusinessStore) GetRatingsByBusinessID(id string) ([]BusinessRating, error) {
	query := `
	SELECT 
		r.id,
		r.business_id,
		r.user_id,
		u.name,
		r.rating,
		r.created_at,
		r.updated_at
	FROM business_ratings r
	JOIN users u
		ON u.id = r.user_id
	WHERE r.business_id = $1;
	`

	res, err := bs.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var ratings []BusinessRating
	for res.Next() {
		var r BusinessRating
		err = res.Scan(
			&r.ID,
			&r.BusinessID,
			&r.UserID,
			&r.UserName,
			&r.Rating,
			&r.CreatedAT,
			&r.UpdatedAT,
		)

		if err != nil {
			return nil, err
		}

		ratings = append(ratings, r)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return ratings, nil
}
