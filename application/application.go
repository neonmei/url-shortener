package application

import (
	"context"
	"errors"
	"math/big"
	"math/rand/v2"
	"net/url"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/neonmei/challenge_urlshortener/platform/o11y/semconv"

	"github.com/neonmei/challenge_urlshortener/domain"
	"github.com/neonmei/challenge_urlshortener/domain/validators"
	"github.com/neonmei/challenge_urlshortener/platform/config"
	"github.com/neonmei/challenge_urlshortener/platform/o11y"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type shortenerService struct {
	urlRepo      domain.URLRepository
	hitCounter   metric.Int64Counter
	serviceMeter metric.Meter
	svcURL       url.URL
	cfg          config.AppConfig
}

func (e shortenerService) Shorten(ctx context.Context, longURL string, author string) (*url.URL, error) {
	u, err := url.Parse(longURL)
	if err != nil {
		return nil, errors.Join(domain.ErrInvalidURL, err)
	}

	if err := validators.ValidateURL(u); err != nil {
		return nil, err
	}

	base62string, err := e.generateHash(ctx)
	if err != nil {
		return nil, err
	}

	newURL := domain.ShortURL{
		ID:        base62string,
		Upstream:  *u,
		CreatedBy: author,
		CreatedAt: time.Now(),
		Enabled:   true,
	}

	o11y.TraceShortURL(ctx, &newURL)
	if err := validators.ValidateShortURL(newURL); err != nil {
		return nil, err
	}

	if err := e.urlRepo.Save(ctx, newURL); err != nil {
		return nil, errors.Join(domain.ErrUnavailableRepo, err)
	}

	return e.svcURL.JoinPath(newURL.ID), nil
}

func (e shortenerService) Redirect(ctx context.Context, urlID string) (string, error) {
	urlEntry, err := e.urlRepo.Get(ctx, urlID)
	o11y.TraceShortURL(ctx, urlEntry)

	// If redirection cannot be performed because of disabled entry, turn it into an err
	if urlEntry != nil && !urlEntry.Enabled {
		err = errors.Join(domain.ErrCannotUseDisabled)
	}

	if err != nil {
		return "", err
	}

	metric.WithAttributeSet(attribute.NewSet())

	e.hitCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("url_id", urlID)),
	)

	return urlEntry.Upstream.String(), nil
}

func (e shortenerService) Delete(ctx context.Context, urlID string) error {
	return e.urlRepo.Delete(ctx, urlID)
}

func (e shortenerService) Fetch(ctx context.Context, urlID string) (*domain.ShortURL, error) {
	urlEntry, err := e.urlRepo.Get(ctx, urlID)
	if err != nil {
		return nil, err
	}

	o11y.TraceShortURL(ctx, urlEntry)
	return urlEntry, err
}

func (e shortenerService) generateHash(ctx context.Context) (string, error) {
	currentRounds := uint64(0)
	base62string := ""
	resultErr := error(nil)

	for currentRounds < e.cfg.Hasher.MaxRounds {
		bigInt := big.Int{}
		base62string = bigInt.SetUint64(rand.Uint64N(e.cfg.Hasher.RandomMaxValue)).Text(62)
		_, resultErr = e.urlRepo.Get(ctx, base62string)

		// REF: already exists
		if resultErr == nil {
			currentRounds++
			continue
		}

		// REF: does not exist
		if errors.Is(resultErr, domain.ErrURLNotFound) {
			trace.SpanFromContext(ctx).SetAttributes(
				attribute.Int64(semconv.HasherRounds, int64(currentRounds)),
				attribute.Int(semconv.HasherLength, len(base62string)),
			)
			return base62string, nil
		}

		if resultErr != nil {
			break
		}
	}

	return "", errors.Join(domain.ErrUnavailableRepo, resultErr)
}

func New(cfg config.AppConfig, urlRepo domain.URLRepository) (Service, error) {
	m := otel.GetMeterProvider().Meter("application")
	c, err := m.Int64Counter(
		semconv.MetricURLHits,
		metric.WithDescription("Number of URL hits."),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		return nil, err
	}

	baseHost, err := url.Parse(cfg.BaseUrl)
	if err != nil {
		return nil, err
	}

	return &shortenerService{
		urlRepo:      urlRepo,
		hitCounter:   c,
		serviceMeter: m,
		svcURL:       *baseHost,
		cfg:          cfg,
	}, nil
}
