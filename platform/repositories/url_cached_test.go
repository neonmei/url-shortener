package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/neonmei/challenge_urlshortener/domain"
	"github.com/stretchr/testify/assert"
)

func makeCache(t *testing.T) *URLCache {
	cache, err := ristretto.NewCache(&URLCacheConfig{
		NumCounters: 10000,
		MaxCost:     5000,
		BufferItems: 64,
	})

	assert.NoError(t, err)
	return cache
}

func TestCachedBasic(t *testing.T) {
	upstreamRepo := NewMemory()
	cachedRepo := NewCached(upstreamRepo, makeCache(t))
	ctx := context.Background()

	validItem := domain.ShortURL{
		ID:        validId,
		Upstream:  *validURL,
		CreatedBy: validAuthor,
		CreatedAt: time.Now(),
		Enabled:   true,
	}

	// REF: Missing items in upstream should propagate error
	result, err := cachedRepo.Get(ctx, validId)
	assert.ErrorIs(t, err, domain.ErrURLNotFound)
	assert.Nil(t, result)

	// REF: Save into cache, delete from upstream
	assert.NoError(t, cachedRepo.Save(ctx, validItem))
	assert.NoError(t, upstreamRepo.Delete(ctx, validId))

	// REF: cache should answer without reaching out to backend
	result, err = cachedRepo.Get(ctx, validId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, validId, result.ID)
}

func TestCachedFetch(t *testing.T) {
	upstreamRepo := NewMemory()
	cachedRepo := NewCached(upstreamRepo, makeCache(t))
	ctx := context.Background()

	validItem := domain.ShortURL{
		ID:        validId,
		Upstream:  *validURL,
		CreatedBy: validAuthor,
		CreatedAt: time.Now(),
		Enabled:   true,
	}

	// Fetch from upstream
	assert.NoError(t, upstreamRepo.Save(ctx, validItem))
	result, err := cachedRepo.Get(ctx, validId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, validId, result.ID)
}

func TestCachedDeleteBasic(t *testing.T) {
	upstreamRepo := NewMemory()
	cachedRepo := NewCached(upstreamRepo, makeCache(t))
	ctx := context.Background()

	validItem := domain.ShortURL{
		ID:        validId,
		Upstream:  *validURL,
		CreatedBy: validAuthor,
		CreatedAt: time.Now(),
		Enabled:   true,
	}
	// REF: Save into cache, perform logical deletion
	assert.NoError(t, cachedRepo.Save(ctx, validItem))
	assert.NoError(t, cachedRepo.Delete(ctx, validId))
	assert.NoError(t, upstreamRepo.Delete(ctx, validId))

	// REF: cache should answer with a negative cache
	result, err := cachedRepo.Get(ctx, validId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, validId, result.ID)
	assert.Equal(t, false, result.Enabled)
}

func TestCachedDeleteErrShouldPropagate(t *testing.T) {
	upstreamRepo := NewMemory()
	cachedRepo := NewCached(upstreamRepo, makeCache(t))
	ctx := context.Background()

	// REF: Save into cache, perform logical deletion
	assert.ErrorIs(t, domain.ErrURLNotFound, cachedRepo.Delete(ctx, validId))
}
