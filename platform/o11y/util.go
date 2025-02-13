package o11y

import (
	"context"

	"github.com/neonmei/challenge_urlshortener/platform/o11y/semconv"

	"github.com/neonmei/challenge_urlshortener/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TraceShortURL(ctx context.Context, u *domain.ShortURL) *domain.ShortURL {
	if u == nil {
		return nil
	}

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.String(semconv.UrlId, u.ID),
		attribute.String(semconv.UrlAuthor, u.CreatedBy),
		attribute.String(semconv.UrlFull, u.Upstream.String()),
		attribute.Bool(semconv.UrlEnabled, u.Enabled),
		attribute.Int64(semconv.UrlCreated, u.CreatedAt.Unix()),
	)

	return u
}
