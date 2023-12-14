package postgres

import (
	"TelegramBot/storage"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}
func (s *Storage) Create(ctx context.Context, page *storage.Page) error {
	q := fmt.Sprintf("INSERT INTO %s (url,username) values (?,?)", tgstorageTable)

	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return fmt.Errorf("cant save page: %w", err)
	}

	return nil
}
func (s *Storage) Get(ctx context.Context, username string) (*storage.Page, error) {
	q := fmt.Sprintf("SELECT url FROM %s WHERE username=?  ", tgstorageTable)

	var url string

	err := s.db.QueryRowContext(ctx, q, username).Scan(&url)

	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.Page{URL: url}, nil
}
func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM tgstorage WHERE url = ? AND username = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}
func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{}

}
