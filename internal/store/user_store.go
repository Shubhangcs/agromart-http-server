package store

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	hash          []byte
	plainPassword *string
}

func (p *password) Set(plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 12)
	if err != nil {
		return err
	}
	p.hash = hash
	p.plainPassword = &plainPassword
	return nil
}

func (p *password) Matches(plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type Admin struct {
	ID           string    `json:"id"`
	ProfileImage *string   `json:"profile_image"`
	FirstName    string    `json:"first_name"`
	LastName     *string   `json:"last_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Password     password  `json:"-"`
	CreatedAT    time.Time `json:"created_at"`
	UpdatedAT    time.Time `json:"updated_at"`
}

type User struct {
	ID            string    `json:"id"`
	ProfileImage  *string   `json:"profile_image"`
	FirstName     string    `json:"first_name"`
	LastName      *string   `json:"last_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Password      password  `json:"-"`
	IsUserBlocked bool      `json:"is_user_blocked"`
	IsUserSeller  bool      `json:"is_user_seller"`
	CreatedAT     time.Time `json:"created_at"`
	UpdatedAT     time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateAdmin(*Admin) error
	CreateUser(*User) error
	GetAdminByEmail(string) (*Admin, error)
	GetUserByEmail(string) (*User, error)
	UpdateAdminDetails(*Admin) error
	UpdateUserDetails(*User) error
	UpdateAdminPassword(*Admin) error
	UpdateUserPassword(*User) error
	UpdateUserSellerStatus(*User) error
	DeleteAdmin(id string) error
	DeleteUser(id string) error
	GetAllUsers() ([]User, error)
	BlockUser(*User) error
	GetUserDetailsByID(id string) (*User, error)
	GetAdminDetailsByID(id string) (*Admin, error)
}

func (us *PostgresUserStore) CreateAdmin(admin *Admin) error {
	query := `
	INSERT INTO admins(first_name,last_name,email,phone,password_hash)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id, created_at, updated_at
	`
	err := us.db.QueryRow(
		query,
		admin.FirstName,
		admin.LastName,
		admin.Email,
		admin.Phone,
		string(admin.Password.hash),
	).Scan(
		&admin.ID,
		&admin.CreatedAT,
		&admin.UpdatedAT,
	)
	if err != nil {
		return err
	}
	return nil
}

func (us *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users(first_name,last_name,email,phone,password_hash)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id, created_at, updated_at
	`
	err := us.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		string(user.Password.hash),
	).Scan(
		&user.ID,
		&user.CreatedAT,
		&user.UpdatedAT,
	)

	if err != nil {
		return err
	}
	return nil
}

func (us *PostgresUserStore) GetAdminByEmail(email string) (*Admin, error) {

	var admin Admin

	query := `
	SELECT id,profile_image,first_name,last_name,email,phone,password_hash,created_at,updated_at
	FROM admins
	WHERE email=$1
	`

	err := us.db.QueryRow(query, email).Scan(
		&admin.ID,
		&admin.ProfileImage,
		&admin.FirstName,
		&admin.LastName,
		&admin.Email,
		&admin.Phone,
		&admin.Password.hash,
		&admin.CreatedAT,
		&admin.UpdatedAT,
	)

	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (us *PostgresUserStore) GetUserByEmail(email string) (*User, error) {
	var user User

	query := `
	SELECT id,profile_image,first_name,last_name,email,phone,password_hash,created_at,updated_at
	FROM users
	WHERE email=$1
	`

	err := us.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.ProfileImage,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.Password.hash,
		&user.CreatedAT,
		&user.UpdatedAT,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *PostgresUserStore) UpdateAdminDetails(admin *Admin) error {
	query := `
	UPDATE admins 
	SET first_name = COALESCE($1,first_name),
	last_name = COALESCE($2,last_name),
	email = COALESCE($3,email),
	phone = COALESCE($4,phone),
	updated_at = CURRENT_TIMESTAMP
	WHERE id=$5
	`

	res, err := us.db.Exec(query, admin.FirstName, admin.LastName, admin.Email, admin.Phone, admin.ID)
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

func (us *PostgresUserStore) UpdateUserDetails(user *User) error {
	query := `
	UPDATE users
	SET first_name = COALESCE($1, first_name),
	last_name = COALESCE($2, last_name),
	email = COALESCE($3, email),
	phone = COALESCE($4, phone)
	WHERE id = $5
	`

	res, err := us.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Phone, user.ID)
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

func (us *PostgresUserStore) DeleteAdmin(id string) error {
	query := `
	DELETE FROM admins
	WHERE id = $1
	`

	res, err := us.db.Exec(query, id)
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

func (us *PostgresUserStore) DeleteUser(id string) error {
	query := `
	DELETE FROM users
	WHERE id = $1
	`

	res, err := us.db.Exec(query, id)
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

func (us *PostgresUserStore) UpdateAdminPassword(admin *Admin) error {
	query := `
	UPDATE admins
	SET password_hash = $1
	WHERE id = $2
	`

	res, err := us.db.Exec(query, string(admin.Password.hash), admin.ID)
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

func (us *PostgresUserStore) UpdateUserPassword(user *User) error {
	query := `
	UPDATE users
	SET password_hash = $1
	WHERE id = $2
	`

	res, err := us.db.Exec(query, string(user.Password.hash), user.ID)
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

func (us *PostgresUserStore) GetAllUsers() ([]User, error) {
	query := `
	SELECT id, profile_image, first_name, last_name, email, phone, is_user_seller , is_user_blocked, created_at, updated_at
	FROM users
	`
	res, err := us.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var users []User
	for res.Next() {
		var user User
		err = res.Scan(
			&user.ID,
			&user.ProfileImage,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
			&user.IsUserBlocked,
			&user.IsUserSeller,
			&user.CreatedAT,
			&user.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return users, nil
}

func (us *PostgresUserStore) BlockUser(user *User) error {
	query := `
	UPDATE users 
	SET is_user_blocked = $1,
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $2
	`
	res, err := us.db.Exec(query, user.IsUserBlocked, user.ID)
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

func (us *PostgresUserStore) UpdateUserSellerStatus(user *User) error {
	query := `
	UPDATE users 
	SET is_user_seller = $1,
	updated_at = CURRENT_TIMESTAMP,
	WHERE id = $2
	`

	res, err := us.db.Exec(query, user.IsUserSeller, user.ID)
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

func (us *PostgresUserStore) GetUserDetailsByID(id string) (*User, error) {
	query := `
	SELECT 
		id,
		profile_image,
		first_name,
		last_name,
		email,
		phone,
		created_at,
		updated_at,
		is_user_blocked,
		is_user_seller
	FROM users
	WHERE id = $1;
	`
	var user User
	err := us.db.QueryRow(
		query,
		id,
	).Scan(
		&user.ID,
		&user.ProfileImage,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.CreatedAT,
		&user.UpdatedAT,
		&user.IsUserBlocked,
		&user.IsUserSeller,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *PostgresUserStore) GetAdminDetailsByID(id string) (*Admin, error) {
	query := `
	SELECT 
		id,
		profile_image,
		first_name,
		last_name,
		email,
		phone,
		created_at,
		updated_at
	FROM admins
	WHERE id = $1;
	`
	var user Admin
	err := us.db.QueryRow(
		query,
		id,
	).Scan(
		&user.ID,
		&user.ProfileImage,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.CreatedAT,
		&user.UpdatedAT,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
