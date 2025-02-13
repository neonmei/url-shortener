package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/neonmei/challenge_urlshortener/domain"
	"github.com/stretchr/testify/assert"
)

func TestInmemRepoBasic(t *testing.T) {
	repo := NewMemory()
	ctx := context.Background()

	validItem := domain.ShortURL{
		ID:        validId,
		Upstream:  *validURL,
		CreatedBy: validAuthor,
		CreatedAt: time.Now(),
		Enabled:   true,
	}

	// REF: Save-time validation
	assert.Error(t, repo.Save(ctx, domain.ShortURL{
		ID:        validId,
		Upstream:  *URLBadSchema,
		CreatedBy: validAuthor,
		CreatedAt: time.Now(),
		Enabled:   true,
	}))

	retrieved, err := repo.Get(ctx, validId)
	assert.Nil(t, retrieved)
	assert.ErrorIs(t, err, domain.ErrURLNotFound)

	// REF: Existance
	err = repo.Save(ctx, validItem)
	assert.NoError(t, err)

	retrieved, err = repo.Get(ctx, validId)
	assert.NoError(t, err)
	assert.Equal(t, validItem, *retrieved)

	// REF: Deletion
	err = repo.Delete(ctx, validId)
	assert.NoError(t, err)

	retrieved, err = repo.Get(ctx, validId)
	assert.Nil(t, retrieved)
	assert.ErrorIs(t, err, domain.ErrURLNotFound)
}
