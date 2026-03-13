package store

import (
	"database/sql"

	"github.com/shubhangcs/agromart-server/internal/models"
)

type PostgresFollowerStore struct {
	db *sql.DB
}

type FollowerStore interface {
	CreateFollower(*models.Follower) error
	RemoveFollower(*models.Follower) error
	GetFollowersCount(id string) (int, error)
	GetFollowingCount(id string) (int, error)
	GetAllFollowers(id string, limit, offset int) ([]models.FollowerDetails, error)
	GetAllFollowing(id string, limit, offset int) ([]models.FollowingDetails, error)
}

func NewPostgresFollowerStore(db *sql.DB) *PostgresFollowerStore {
	return &PostgresFollowerStore{db: db}
}

func (pfs *PostgresFollowerStore) CreateFollower(f *models.Follower) error {
	// INSERT ... ON CONFLICT DO NOTHING prevents duplicate-follow errors.
	query := `
	INSERT INTO followers (business_id, user_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	`
	_, err := pfs.db.Exec(query, f.BusinessID, f.UserID)
	return err
}

func (pfs *PostgresFollowerStore) RemoveFollower(f *models.Follower) error {
	res, err := pfs.db.Exec(
		`DELETE FROM followers WHERE user_id = $1 AND business_id = $2`,
		f.UserID, f.BusinessID,
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

func (pfs *PostgresFollowerStore) GetFollowersCount(id string) (int, error) {
	var count int
	err := pfs.db.QueryRow(`SELECT COUNT(*) FROM followers WHERE business_id = $1`, id).Scan(&count)
	return count, err
}

func (pfs *PostgresFollowerStore) GetFollowingCount(id string) (int, error) {
	var count int
	err := pfs.db.QueryRow(`SELECT COUNT(*) FROM followers WHERE user_id = $1`, id).Scan(&count)
	return count, err
}

// GetAllFollowers returns a paginated list of users following a business.
func (pfs *PostgresFollowerStore) GetAllFollowers(id string, limit, offset int) ([]models.FollowerDetails, error) {
	query := `
	SELECT u.id,
	       CONCAT(u.first_name, ' ', COALESCE(u.last_name, '')) AS name,
	       u.profile_image, u.email, u.phone, f.created_at
	FROM followers f
	JOIN users u ON u.id = f.user_id
	WHERE f.business_id = $1
	ORDER BY f.created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := pfs.db.Query(query, id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []models.FollowerDetails
	for rows.Next() {
		var f models.FollowerDetails
		err = rows.Scan(&f.FollowerID, &f.FollowerName, &f.FollowerProfileImage, &f.FollowerEmail, &f.FollowerPhone, &f.CreatedAT)
		if err != nil {
			return nil, err
		}
		followers = append(followers, f)
	}
	return followers, rows.Err()
}

// GetAllFollowing returns a paginated list of businesses a user is following.
func (pfs *PostgresFollowerStore) GetAllFollowing(id string, limit, offset int) ([]models.FollowingDetails, error) {
	query := `
	SELECT b.id, b.business_profile_image, b.business_name, b.business_phone,
	       b.address, b.city, b.state, bs.telegram
	FROM followers f
	JOIN businesses b ON b.id = f.business_id
	LEFT JOIN business_socials bs ON bs.business_id = f.business_id
	WHERE f.user_id = $1
	ORDER BY f.created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := pfs.db.Query(query, id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followings []models.FollowingDetails
	for rows.Next() {
		var f models.FollowingDetails
		err = rows.Scan(
			&f.FollowingID, &f.FollowingProfileImage, &f.FollowingName,
			&f.FollowingPhone, &f.FollowingAddress, &f.FollowingCity,
			&f.FollowingState, &f.FollowingTelegram,
		)
		if err != nil {
			return nil, err
		}
		followings = append(followings, f)
	}
	return followings, rows.Err()
}
