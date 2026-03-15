package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// --- DB Models ---

type password struct {
	Hash          []byte
	plainPassword *string
}

func (p *password) Set(plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 12)
	if err != nil {
		return err
	}
	p.Hash = hash
	p.plainPassword = &plainPassword
	return nil
}

func (p *password) Matches(plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plainPassword))
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

// --- Request DTOs ---

// CreateAdminRequest is the payload for creating a new admin.
// swagger:model
type CreateAdminRequest struct {
	FirstName string  `json:"first_name" validate:"required"       example:"John"`
	LastName  *string `json:"last_name"                            example:"Doe"`
	Email     string  `json:"email"      validate:"required,email" example:"admin@example.com"`
	Phone     string  `json:"phone"      validate:"required,phone" example:"9876543210"`
	Password  string  `json:"password"   validate:"required"       example:"secret123"`
}

// UpdateUserDetailsRequest is the payload for updating user/admin profile details.
// swagger:model
type UpdateUserDetailsRequest struct {
	FirstName string  `json:"first_name"                              example:"John"`
	LastName  *string `json:"last_name"                               example:"Doe"`
	Email     string  `json:"email"      validate:"omitempty,email"   example:"user@example.com"`
	Phone     string  `json:"phone"      validate:"omitempty,phone"   example:"9876543210"`
}

// UpdatePasswordRequest is the payload for changing a password.
// swagger:model
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required" example:"oldpass123"`
	NewPassword string `json:"new_password" validate:"required" example:"newpass456"`
}

// LoginRequest is the payload for login endpoints.
// swagger:model
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required"       example:"secret123"`
}

// BlockUserRequest is the payload for blocking/unblocking a user.
// swagger:model
type BlockUserRequest struct {
	IsUserBlocked bool `json:"is_user_blocked" example:"true"`
}

// --- Response DTOs ---

// UserResponse is the safe public representation of a user.
// swagger:model
type UserResponse struct {
	ID            string    `json:"id"`
	ProfileImage  *string   `json:"profile_image"`
	FirstName     string    `json:"first_name"`
	LastName      *string   `json:"last_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	IsUserBlocked bool      `json:"is_user_blocked"`
	IsUserSeller  bool      `json:"is_user_seller"`
	CreatedAT     time.Time `json:"created_at"`
	UpdatedAT     time.Time `json:"updated_at"`
}

// AdminResponse is the safe public representation of an admin.
// swagger:model
type AdminResponse struct {
	ID           string    `json:"id"`
	ProfileImage *string   `json:"profile_image"`
	FirstName    string    `json:"first_name"`
	LastName     *string   `json:"last_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	CreatedAT    time.Time `json:"created_at"`
	UpdatedAT    time.Time `json:"updated_at"`
}
