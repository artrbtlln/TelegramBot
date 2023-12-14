package storage

import (
	"TelegramBot/lib/e"
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

type Storage interface {
	Create(ctx context.Context, page *Page) error
	Get(ctx context.Context, username string) (*Page, error)
	IsExists(ctx context.Context, page *Page) (bool, error)
}
type Page struct {
	URL      string
	UserName string
}

var ErrNoSavedPages = errors.New("no saved pages")

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
