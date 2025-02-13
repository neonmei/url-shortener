package application

import (
	"context"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/neonmei/challenge_urlshortener/domain"
	"github.com/neonmei/challenge_urlshortener/domain/validators"
	mockDomain "github.com/neonmei/challenge_urlshortener/mocks/domain"
	"github.com/neonmei/challenge_urlshortener/platform/config"
	"github.com/neonmei/challenge_urlshortener/platform/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	baseURL, _    = url.Parse("https://base.url")
	validURL, _   = url.Parse("https://opentelemetry.io")
	validAuthor   = "root@neonmei.cloud"
	validId       = "someId"
	invalidURL, _ = url.Parse("http://opentelemetry.io")
)

func TestBasicShortenShouldGenerate(t *testing.T) {
	ctx := context.Background()
	cfg := config.Load()
	cfg.BaseUrl = baseURL.String()
	svc, err := New(cfg, repositories.NewMemory())
	assert.NoError(t, err)

	u, err := svc.Shorten(ctx, validURL.String(), validAuthor)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, u.Scheme, baseURL.Scheme)
	assert.Equal(t, u.Host, baseURL.Host)
	assert.NoError(t, validators.ValidateId(u.Path))
}

func TestBadURLShouldNotValidate(t *testing.T) {
	ctx := context.Background()
	cfg := config.Load()
	svc, err := New(cfg, repositories.NewMemory())
	assert.NoError(t, err)

	u, err := svc.Shorten(ctx, invalidURL.String(), validAuthor)
	assert.ErrorIs(t, err, domain.ErrInvalidURL)
	assert.Nil(t, u)
}

func TestBadURLShouldNotParse(t *testing.T) {
	ctx := context.Background()
	cfg := config.Load()
	svc, err := New(cfg, repositories.NewMemory())
	assert.NoError(t, err)

	u, err := svc.Shorten(ctx, "hello!", validAuthor)
	assert.ErrorIs(t, err, domain.ErrInvalidURL)
	assert.Nil(t, u)
}

func TestBrokenBackendErrShouldPropagate(t *testing.T) {
	ctx := context.Background()
	cfg := config.Load()
	repo := mockDomain.NewMockURLRepository(t)
	repoErr := errors.New("unknown storage error")
	repo.On("Get", mock.Anything, mock.Anything).Return(nil, repoErr)

	svc, err := New(cfg, repo)
	assert.NoError(t, err)

	u, err := svc.Shorten(ctx, validURL.String(), validAuthor)
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, u)
}

func TestBrokenBackendSaveErrShouldPropagate(t *testing.T) {
	ctx := context.Background()
	cfg := config.Load()
	repo := mockDomain.NewMockURLRepository(t)
	repoErr := errors.New("unknown storage error")

	repo.On("Get", mock.Anything, mock.Anything).Return(nil, domain.ErrURLNotFound)
	repo.On("Save", mock.Anything, mock.Anything).Return(repoErr)

	svc, err := New(cfg, repo)
	assert.NoError(t, err)

	u, err := svc.Shorten(ctx, validURL.String(), validAuthor)
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, u)
}

func TestRedirectOk(t *testing.T) {
	ctx := context.Background()
	cfg := config.Load()
	cfg.BaseUrl = baseURL.String()
	svc, err := New(cfg, repositories.NewMemory())
	assert.NoError(t, err)

	u, err := svc.Shorten(ctx, validURL.String(), validAuthor)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	upstream, err := svc.Redirect(ctx, u.Path)
	assert.NoError(t, err)
	assert.Equal(t, validURL.String(), upstream)
}

func TestRedirectDisabled(t *testing.T) {
	ctx := context.Background()
	cfg := config.Load()
	cfg.BaseUrl = baseURL.String()

	repo := mockDomain.NewMockURLRepository(t)

	repo.On("Get", mock.Anything, mock.Anything).Return(&domain.ShortURL{
		ID:        validId,
		Upstream:  *validURL,
		CreatedBy: validAuthor,
		CreatedAt: time.Now(),
		Enabled:   false,
	}, nil)

	svc, err := New(cfg, repo)
	assert.NoError(t, err)

	upstream, err := svc.Redirect(ctx, validURL.Path)
	assert.ErrorIs(t, err, domain.ErrCannotUseDisabled)
	assert.Equal(t, "", upstream)
}

func TestDeleteOk(t *testing.T) {
	ctx := context.Background()
	cfg := config.Load()
	cfg.BaseUrl = baseURL.String()
	svc, err := New(cfg, repositories.NewMemory())
	assert.NoError(t, err)

	u, err := svc.Shorten(ctx, validURL.String(), validAuthor)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	err = svc.Delete(ctx, u.Path)
	assert.NoError(t, err)

	upstream, err := svc.Redirect(ctx, u.Path)
	assert.ErrorIs(t, domain.ErrURLNotFound, err)
	assert.Equal(t, "", upstream)
}

func TestFetchOk(t *testing.T) {
	ctx := context.Background()
	cfg := config.Load()
	cfg.BaseUrl = baseURL.String()
	svc, err := New(cfg, repositories.NewMemory())
	assert.NoError(t, err)

	u, err := svc.Shorten(ctx, validURL.String(), validAuthor)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	upstream, err := svc.Fetch(ctx, u.Path)
	assert.NoError(t, err)
	assert.Equal(t, validURL.String(), upstream.Upstream.String())
}

func TestFetchNonExistant(t *testing.T) {
	ctx := context.Background()
	cfg := config.Load()
	cfg.BaseUrl = baseURL.String()
	svc, err := New(cfg, repositories.NewMemory())
	assert.NoError(t, err)

	upstream, err := svc.Fetch(ctx, validId)
	assert.ErrorIs(t, domain.ErrURLNotFound, err)
	assert.Nil(t, upstream)
}
