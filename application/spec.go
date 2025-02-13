package application

import (
	"context"
	"net/url"

	"github.com/neonmei/challenge_urlshortener/domain"
)

type Service interface {
	Redirect(ctx context.Context, urlID string) (string, error)
	Shorten(ctx context.Context, longURL string, author string) (*url.URL, error)
	Delete(ctx context.Context, urlID string) error
	Fetch(ctx context.Context, urlID string) (*domain.ShortURL, error)
}
