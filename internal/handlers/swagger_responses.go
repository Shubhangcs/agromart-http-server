package handlers

import (
	"github.com/shubhangcs/agromart-server/internal/models"
)

// PaginationInfo defines generic pagination payload
type PaginationInfo struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// UserListResponse defines a paginated user list payload
type UserListResponse struct {
	Message    string         `json:"message" example:"users fetched successfully"`
	Users      []models.User  `json:"users"`
	Pagination PaginationInfo `json:"pagination"`
}

// AdminDetailsResponse defines payload for fetched admin details
type AdminDetailsResponse struct {
	Message string       `json:"message" example:"admin details fetched successfully"`
	Admin   models.Admin `json:"admin"`
}

// UserDetailsResponse defines payload for fetched user details
type UserDetailsResponse struct {
	Message string      `json:"message" example:"user details fetched successfully"`
	User    models.User `json:"user"`
}

// BusinessListResponse defines a paginated business list payload
type BusinessListResponse struct {
	Message    string            `json:"message" example:"businesses fetched successfully"`
	Businesses []models.Business `json:"businesses"`
	Pagination PaginationInfo    `json:"pagination"`
}

// BusinessDetailsResponse defines a business detail payload
type BusinessDetailsResponse struct {
	Message string          `json:"message" example:"business details fetched successfully"`
	Details models.Business `json:"details"`
}

// CategoryListResponse defines a list payload for categories
type CategoryListResponse struct {
	Message    string            `json:"message" example:"categories fetched successfully"`
	Categories []models.Category `json:"categories"`
}

// SubCategoryListResponse defines a list payload for sub-categories
type SubCategoryListResponse struct {
	Message       string               `json:"message" example:"sub categories fetched successfully"`
	SubCategories []models.SubCategory `json:"sub_categories"`
}

// FollowerListResponse defines a paginated follower list payload
type FollowerListResponse struct {
	Message    string            `json:"message"`
	Followers  []models.Follower `json:"followers"`
	Pagination PaginationInfo    `json:"pagination"`
}

// FollowingListResponse defines a paginated following list payload
type FollowingListResponse struct {
	Message    string            `json:"message"`
	Followings []models.Follower `json:"followings"`
	Pagination PaginationInfo    `json:"pagination"`
}

// RFQListResponse defines a list payload for RFQs
type RFQListResponse struct {
	Message string       `json:"message" example:"rfqs fetched successfully"`
	RFQs    []models.RFQResponse `json:"rfqs"`
}

// ProductListResponse defines a list payload for products
type ProductListResponse struct {
	Message  string           `json:"message" example:"products fetched successfully"`
	Products []models.Product `json:"products"`
}

// ChatHistoryResponse defines payload for fetching chat messages between users
type ChatHistoryResponse struct {
	Message    string           `json:"message"`
	Messages   []models.Message `json:"messages"`
	Pagination PaginationInfo   `json:"pagination"`
}
