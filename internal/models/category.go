package models

import "time"

type Category struct {
	ID            string    `json:"id"`
	CategoryImage *string   `json:"category_image"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	CreatedAT     time.Time `json:"created_at"`
	UpdatedAT     time.Time `json:"updated_at"`
}

type SubCategory struct {
	ID               string    `json:"id"`
	CategoryID       string    `json:"category_id"`
	SubCategoryImage *string   `json:"sub_category_image"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	CreatedAT        time.Time `json:"created_at"`
	UpdatedAT        time.Time `json:"updated_at"`
}
