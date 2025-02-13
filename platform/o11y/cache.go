package o11y

import (
	"context"
	"errors"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/dgraph-io/ristretto/v2/z"
	"github.com/neonmei/challenge_urlshortener/platform/o11y/semconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func InstrumentCacheAsync[K z.Key, V any](c *ristretto.Cache[K, V]) error {
	meter := otel.GetMeterProvider().Meter("ristretto")

	_, err1 := meter.Int64ObservableCounter(
		semconv.CacheHit,
		metric.WithDescription("How many cache read operations found the key."),
		metric.WithUnit("{call}"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(c.Metrics.Hits()))
			return nil
		}),
	)

	_, err2 := meter.Int64ObservableCounter(
		semconv.CacheMiss,
		metric.WithDescription("How many cache read operations did not find  the key."),
		metric.WithUnit("{call}"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(c.Metrics.Misses()))
			return nil
		}),
	)

	_, err3 := meter.Int64ObservableCounter(
		semconv.CacheAdded,
		metric.WithDescription("How many cache keys were added."),
		metric.WithUnit("{call}"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(c.Metrics.KeysAdded()))
			return nil
		}),
	)

	_, err4 := meter.Int64ObservableCounter(
		semconv.CacheEvicted,
		metric.WithDescription("How many cache keys were evicted."),
		metric.WithUnit("{call}"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(c.Metrics.KeysEvicted()))
			return nil
		}),
	)

	_, err5 := meter.Int64ObservableCounter(
		semconv.CacheRejected,
		metric.WithDescription("How many cache keys were rejected by admission policy."),
		metric.WithUnit("{call}"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(c.Metrics.SetsRejected()))
			return nil
		}),
	)

	return errors.Join(err1, err2, err3, err4, err5)
}
