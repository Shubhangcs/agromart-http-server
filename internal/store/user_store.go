package store

import (
	"database/sql"

	"github.com/shubhangcs/agromart-server/internal/models"
)

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateAdmin(*models.Admin) error
	CreateUser(*models.User) error
	GetAdminByEmail(string) (*models.Admin, error)
	GetUserByEmail(string) (*models.User, error)
	UpdateAdminDetails(*models.Admin) error
	UpdateUserDetails(*models.User) error
	UpdateAdminPassword(*models.Admin) error
	UpdateUserPassword(*models.User) error
	UpdateUserSellerStatus(*models.User) error
	DeleteAdmin(id string) error
	DeleteUser(id string) error
	GetAllUsers(limit, offset int) ([]models.User, error)
	BlockUser(*models.User) error
	GetUserDetailsByID(id string) (*models.User, error)
	GetAdminDetailsByID(id string) (*models.Admin, error)
	AdminExists() (bool, error)
	GetUserByGoogleSub(sub string) (*models.User, error)
	GetUserAuthByID(id string) (*models.User, error)
	CreateGoogleUser(user *models.User) error
	LinkGoogle(userID, sub string) error
}

func (us *PostgresUserStore) CreateAdmin(admin *models.Admin) error {
	query := `
	INSERT INTO admins(first_name, last_name, email, phone, password_hash)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`
	return us.db.QueryRow(
		query,
		admin.FirstName, admin.LastName, admin.Email, admin.Phone,
		string(admin.Password.Hash),
	).Scan(&admin.ID, &admin.CreatedAT, &admin.UpdatedAT)
}

func (us *PostgresUserStore) CreateUser(user *models.User) error {
	query := `
	INSERT INTO users(first_name, last_name, email, phone, password_hash)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`
	return us.db.QueryRow(
		query,
		user.FirstName, user.LastName, user.Email, user.Phone,
		string(user.Password.Hash),
	).Scan(&user.ID, &user.CreatedAT, &user.UpdatedAT)
}

func (us *PostgresUserStore) GetAdminByEmail(email string) (*models.Admin, error) {
	query := `
	SELECT id, profile_image, first_name, last_name, email, phone, password_hash, created_at, updated_at
	FROM admins
	WHERE email = $1
	`
	var admin models.Admin
	err := us.db.QueryRow(query, email).Scan(
		&admin.ID, &admin.ProfileImage, &admin.FirstName, &admin.LastName,
		&admin.Email, &admin.Phone, &admin.Password.Hash,
		&admin.CreatedAT, &admin.UpdatedAT,
	)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (us *PostgresUserStore) GetUserByEmail(email string) (*models.User, error) {
	query := `
	SELECT id, profile_image, first_name, last_name, email, COALESCE(phone, ''), password_hash, COALESCE(auth_provider, 'password'), google_sub, created_at, updated_at
	FROM users
	WHERE email = $1
	`
	var user models.User
	err := us.db.QueryRow(query, email).Scan(
		&user.ID, &user.ProfileImage, &user.FirstName, &user.LastName,
		&user.Email, &user.Phone, &user.Password.Hash, &user.AuthProvider, &user.GoogleSub,
		&user.CreatedAT, &user.UpdatedAT,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *PostgresUserStore) UpdateAdminDetails(admin *models.Admin) error {
	query := `
	UPDATE admins
	SET first_name  = COALESCE(NULLIF($1, ''), first_name),
	    last_name   = COALESCE(NULLIF($2, ''), last_name),
	    email       = COALESCE(NULLIF($3, ''), email),
	    phone       = COALESCE(NULLIF($4, ''), phone),
	    updated_at  = CURRENT_TIMESTAMP
	WHERE id = $5
	`
	res, err := us.db.Exec(query, admin.FirstName, admin.LastName, admin.Email, admin.Phone, admin.ID)
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

func (us *PostgresUserStore) UpdateUserDetails(user *models.User) error {
	query := `
	UPDATE users
	SET first_name  = COALESCE(NULLIF($1, ''), first_name),
	    last_name   = COALESCE(NULLIF($2, ''), last_name),
	    email       = COALESCE(NULLIF($3, ''), email),
	    phone       = COALESCE(NULLIF($4, ''), phone),
	    updated_at  = CURRENT_TIMESTAMP
	WHERE id = $5
	`
	res, err := us.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Phone, user.ID)
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

func (us *PostgresUserStore) UpdateAdminPassword(admin *models.Admin) error {
	query := `UPDATE admins SET password_hash = $1 WHERE id = $2`
	res, err := us.db.Exec(query, string(admin.Password.Hash), admin.ID)
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

func (us *PostgresUserStore) UpdateUserPassword(user *models.User) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	res, err := us.db.Exec(query, string(user.Password.Hash), user.ID)
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

func (us *PostgresUserStore) UpdateUserSellerStatus(user *models.User) error {
	query := `
	UPDATE users
	SET is_user_seller = $1,
	    updated_at     = CURRENT_TIMESTAMP
	WHERE id = $2
	`
	res, err := us.db.Exec(query, user.IsUserSeller, user.ID)
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

func (us *PostgresUserStore) DeleteAdmin(id string) error {
	res, err := us.db.Exec(`DELETE FROM admins WHERE id = $1`, id)
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

func (us *PostgresUserStore) DeleteUser(id string) error {
	res, err := us.db.Exec(`DELETE FROM users WHERE id = $1`, id)
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

// GetAllUsers returns a paginated list of users (without password hashes).
func (us *PostgresUserStore) GetAllUsers(limit, offset int) ([]models.User, error) {
	query := `
	SELECT id, profile_image, first_name, last_name, email, COALESCE(phone, ''), COALESCE(auth_provider, 'password'),
	       is_user_seller, is_user_blocked, created_at, updated_at
	FROM users
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`
	rows, err := us.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err = rows.Scan(
			&u.ID, &u.ProfileImage, &u.FirstName, &u.LastName,
			&u.Email, &u.Phone, &u.AuthProvider,
			&u.IsUserSeller, &u.IsUserBlocked, // order matches query
			&u.CreatedAT, &u.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (us *PostgresUserStore) BlockUser(user *models.User) error {
	query := `
	UPDATE users
	SET is_user_blocked = $1,
	    updated_at      = CURRENT_TIMESTAMP
	WHERE id = $2
	`
	res, err := us.db.Exec(query, user.IsUserBlocked, user.ID)
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

func (us *PostgresUserStore) GetUserDetailsByID(id string) (*models.User, error) {
	query := `
	SELECT id, profile_image, first_name, last_name, email, COALESCE(phone, ''), COALESCE(auth_provider, 'password'),
	       is_user_seller, is_user_blocked, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	var u models.User
	err := us.db.QueryRow(query, id).Scan(
		&u.ID, &u.ProfileImage, &u.FirstName, &u.LastName,
		&u.Email, &u.Phone, &u.AuthProvider,
		&u.IsUserSeller, &u.IsUserBlocked,
		&u.CreatedAT, &u.UpdatedAT,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (us *PostgresUserStore) GetAdminDetailsByID(id string) (*models.Admin, error) {
	query := `
	SELECT id, profile_image, first_name, last_name, email, phone, created_at, updated_at
	FROM admins
	WHERE id = $1
	`
	var a models.Admin
	err := us.db.QueryRow(query, id).Scan(
		&a.ID, &a.ProfileImage, &a.FirstName, &a.LastName,
		&a.Email, &a.Phone, &a.CreatedAT, &a.UpdatedAT,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// AdminExists reports whether at least one admin account exists (used to gate bootstrap).
func (us *PostgresUserStore) AdminExists() (bool, error) {
	var exists bool
	err := us.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM admins)`).Scan(&exists)
	return exists, err
}

const userAuthCols = `id, profile_image, first_name, last_name, email, COALESCE(phone, ''), password_hash, COALESCE(auth_provider, 'password'), google_sub, is_user_seller, is_user_blocked, created_at, updated_at`

func (us *PostgresUserStore) scanUserAuth(row *sql.Row) (*models.User, error) {
	var u models.User
	err := row.Scan(&u.ID, &u.ProfileImage, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.Password.Hash,
		&u.AuthProvider, &u.GoogleSub, &u.IsUserSeller, &u.IsUserBlocked, &u.CreatedAT, &u.UpdatedAT)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByGoogleSub finds a user previously linked to a Google account.
func (us *PostgresUserStore) GetUserByGoogleSub(sub string) (*models.User, error) {
	return us.scanUserAuth(us.db.QueryRow(`SELECT `+userAuthCols+` FROM users WHERE google_sub = $1`, sub))
}

// GetUserAuthByID returns the user including the password hash (for old-password verification).
func (us *PostgresUserStore) GetUserAuthByID(id string) (*models.User, error) {
	return us.scanUserAuth(us.db.QueryRow(`SELECT `+userAuthCols+` FROM users WHERE id = $1`, id))
}

// CreateGoogleUser inserts a password-less user created from a Google identity.
func (us *PostgresUserStore) CreateGoogleUser(user *models.User) error {
	return us.db.QueryRow(`
		INSERT INTO users (first_name, last_name, email, phone, password_hash, profile_image, auth_provider, google_sub)
		VALUES ($1, $2, $3, NULL, NULL, $4, 'google', $5)
		RETURNING id, created_at, updated_at`,
		user.FirstName, user.LastName, user.Email, user.ProfileImage, user.GoogleSub,
	).Scan(&user.ID, &user.CreatedAT, &user.UpdatedAT)
}

// LinkGoogle attaches a Google identity to an existing (password) account.
func (us *PostgresUserStore) LinkGoogle(userID, sub string) error {
	_, err := us.db.Exec(`UPDATE users SET google_sub = $1, updated_at = NOW() WHERE id = $2`, sub, userID)
	return err
}
