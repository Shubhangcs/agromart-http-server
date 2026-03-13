package models

import "time"

type Follower struct {
	UserID     string    `json:"user_id"`
	BusinessID string    `json:"business_id"`
	CreatedAT  time.Time `json:"created_at"`
}

type FollowerDetails struct {
	FollowerID           string    `json:"follower_id"`
	FollowerProfileImage *string   `json:"follower_profile_image"`
	FollowerName         string    `json:"follower_name"`
	FollowerEmail        string    `json:"follower_email"`
	FollowerPhone        string    `json:"follower_phone"`
	CreatedAT            time.Time `json:"created_at"`
}

type FollowingDetails struct {
	FollowingID           string  `json:"following_id"`
	FollowingProfileImage *string `json:"following_profile_image"`
	FollowingName         string  `json:"following_name"`
	FollowingPhone        string  `json:"following_phone"`
	FollowingAddress      string  `json:"following_address"`
	FollowingCity         string  `json:"following_city"`
	FollowingState        string  `json:"following_state"`
	FollowingTelegram     *string `json:"following_telegram"`
}
