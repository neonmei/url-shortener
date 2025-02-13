package repositories

import (
	"context"
	"errors"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/neonmei/challenge_urlshortener/domain"
	"github.com/neonmei/challenge_urlshortener/platform/o11y/semconv"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	URLCache       = ristretto.Cache[string, domain.ShortURL]
	URLCacheConfig = ristretto.Config[string, domain.ShortURL]
)

const CachedShortURLCost = 1

type cachedRepository struct {
	upstream domain.URLRepository
	cache    *URLCache
}

func (d *cachedRepository) Save(ctx context.Context, shortUrl domain.ShortURL) error {
	err := d.upstream.Save(ctx, shortUrl)
	if err != nil {
		return err
	}

	d.cache.Set(shortUrl.ID, shortUrl, CachedShortURLCost)
	d.cache.Wait()
	return nil
}

func (d *cachedRepository) Get(ctx context.Context, urlID string) (*domain.ShortURL, error) {
	cacheItem, found := d.cache.Get(urlID)
	if found {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Bool(semconv.CacheHit, true))
		return &cacheItem, nil
	}

	result, err := d.upstream.Get(ctx, urlID)
	if err != nil {
		return nil, err
	}

	trace.SpanFromContext(ctx).SetAttributes(attribute.Bool(semconv.CacheHit, false))
	d.cache.Set(urlID, *result, CachedShortURLCost)
	return result, nil
}

func (d *cachedRepository) Delete(ctx context.Context, urlID string) error {
	if _, err := d.upstream.Get(ctx, urlID); err != nil {
		// This is more of a consistency assertion
		if errors.Is(err, domain.ErrURLNotFound) {
			d.cache.Del(urlID)
		}
		return err
	}

	err := d.upstream.Delete(ctx, urlID)
	if err != nil {
		return err
	}

	d.tryNegativeCache(urlID)
	return nil
}

// tryNegativeCache if item is in cache, mark it as disabled
func (d *cachedRepository) tryNegativeCache(shortId string) {
	shortUrl, found := d.cache.Get(shortId)
	if !found {
		return
	}

	shortUrl.Enabled = false
	d.cache.Set(shortUrl.ID, shortUrl, CachedShortURLCost)
	d.cache.Wait()
}

func NewCached(repo domain.URLRepository, cache *URLCache) domain.URLRepository {
	return &cachedRepository{
		cache:    cache,
		upstream: repo,
	}
}
