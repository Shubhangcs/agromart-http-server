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
	SELECT id, profile_image, first_name, last_name, email, phone, password_hash, created_at, updated_at
	FROM users
	WHERE email = $1
	`
	var user models.User
	err := us.db.QueryRow(query, email).Scan(
		&user.ID, &user.ProfileImage, &user.FirstName, &user.LastName,
		&user.Email, &user.Phone, &user.Password.Hash,
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
	SELECT id, profile_image, first_name, last_name, email, phone,
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
			&u.Email, &u.Phone,
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
	SELECT id, profile_image, first_name, last_name, email, phone,
	       is_user_seller, is_user_blocked, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	var u models.User
	err := us.db.QueryRow(query, id).Scan(
		&u.ID, &u.ProfileImage, &u.FirstName, &u.LastName,
		&u.Email, &u.Phone,
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
