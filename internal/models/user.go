package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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
