package store

import (
	"database/sql"

	"github.com/shubhangcs/agromart-server/internal/models"
)

type BannerStore interface {
	Create(b *models.Banner) error
	Update(b *models.Banner) error
	UpdateImage(id, path string) error
	Delete(id string) error
	GetActive() ([]models.Banner, error)
	GetAll() ([]models.Banner, error)
}

type PostgresBannerStore struct{ db *sql.DB }

func NewPostgresBannerStore(db *sql.DB) *PostgresBannerStore { return &PostgresBannerStore{db: db} }

const bannerCols = `id, title, image, target_url, sort_order, is_active, created_at, updated_at`

func (s *PostgresBannerStore) Create(b *models.Banner) error {
	return s.db.QueryRow(`
		INSERT INTO banners (title, target_url, sort_order, is_active)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`,
		b.Title, b.TargetURL, b.SortOrder, b.IsActive).Scan(&b.ID, &b.CreatedAT, &b.UpdatedAT)
}

func (s *PostgresBannerStore) Update(b *models.Banner) error {
	res, err := s.db.Exec(`
		UPDATE banners SET title = COALESCE($1, title), target_url = COALESCE($2, target_url),
		       sort_order = $3, is_active = $4, updated_at = NOW()
		WHERE id = $5`, b.Title, b.TargetURL, b.SortOrder, b.IsActive, b.ID)
	return rowsAffectedErr(res, err)
}

func (s *PostgresBannerStore) UpdateImage(id, path string) error {
	res, err := s.db.Exec(`UPDATE banners SET image = $1, updated_at = NOW() WHERE id = $2`, path, id)
	return rowsAffectedErr(res, err)
}

func (s *PostgresBannerStore) Delete(id string) error {
	res, err := s.db.Exec(`DELETE FROM banners WHERE id = $1`, id)
	return rowsAffectedErr(res, err)
}

func (s *PostgresBannerStore) GetActive() ([]models.Banner, error) {
	return s.list(`SELECT ` + bannerCols + ` FROM banners WHERE is_active = TRUE AND image IS NOT NULL ORDER BY sort_order, created_at DESC`)
}

func (s *PostgresBannerStore) GetAll() ([]models.Banner, error) {
	return s.list(`SELECT ` + bannerCols + ` FROM banners ORDER BY sort_order, created_at DESC`)
}

func (s *PostgresBannerStore) list(query string) ([]models.Banner, error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	banners := []models.Banner{}
	for rows.Next() {
		var b models.Banner
		if err = rows.Scan(&b.ID, &b.Title, &b.Image, &b.TargetURL, &b.SortOrder, &b.IsActive, &b.CreatedAT, &b.UpdatedAT); err != nil {
			return nil, err
		}
		banners = append(banners, b)
	}
	return banners, rows.Err()
}

func rowsAffectedErr(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
