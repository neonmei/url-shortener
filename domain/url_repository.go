package domain

import (
	"context"
)

type URLRepository interface {
	Get(ctx context.Context, urlID string) (*ShortURL, error)
	Delete(ctx context.Context, urlID string) error
	Save(ctx context.Context, shortUrl ShortURL) error
}
