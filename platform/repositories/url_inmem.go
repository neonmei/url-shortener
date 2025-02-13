package repositories

import (
	"context"

	"github.com/neonmei/challenge_urlshortener/domain"
	"github.com/neonmei/challenge_urlshortener/domain/validators"
)

type memoryRepo struct {
	data map[string]domain.ShortURL
}

func (d *memoryRepo) Delete(_ context.Context, urlID string) error {
	delete(d.data, urlID)
	return nil
}

func (d *memoryRepo) Save(_ context.Context, shortUrl domain.ShortURL) error {
	if err := validators.ValidateShortURL(shortUrl); err != nil {
		return err
	}

	d.data[shortUrl.ID] = shortUrl
	return nil
}

func (d *memoryRepo) Get(_ context.Context, urlID string) (*domain.ShortURL, error) {
	result, found := d.data[urlID]
	if !found {
		return nil, domain.ErrURLNotFound
	}

	return &result, nil
}

// NewMemory is an in-memory repository designed for troubleshooting and development
func NewMemory() domain.URLRepository {
	return &memoryRepo{data: map[string]domain.ShortURL{}}
}
