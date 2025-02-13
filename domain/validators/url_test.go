package validators

import (
	"net/url"
	"testing"
	"time"

	"github.com/neonmei/challenge_urlshortener/domain"
	"github.com/stretchr/testify/assert"
)

var (
	URLBadSchema, _ = url.Parse("http://opentelemetry.io")
	URLBadHost, _   = url.Parse("https://")
	validURL, _     = url.Parse("https://opentelemetry.io")
	validAuthor     = "root@neonmei.cloud"
	validId         = "asd"
)

func TestValidateShouldPass(t *testing.T) {
	u := domain.ShortURL{
		ID:        validId,
		Upstream:  *validURL,
		CreatedBy: validAuthor,
		CreatedAt: time.Now(),
		Enabled:   true,
	}

	assert.NoError(t, ValidateShortURL(u))
}

func TestValidateShouldFail(t *testing.T) {
	cases := []struct {
		err  error
		item domain.ShortURL
	}{
		{
			err: domain.ErrEmptyId,
			item: domain.ShortURL{
				Upstream:  *validURL,
				CreatedBy: validAuthor,
				CreatedAt: time.Now(),
				Enabled:   true,
			},
		},
		{
			err: domain.ErrInvalidId,
			item: domain.ShortURL{
				ID:        "asd-asd",
				Upstream:  *validURL,
				CreatedBy: validAuthor,
				CreatedAt: time.Now(),
				Enabled:   true,
			},
		},
		{
			err: domain.ErrInvalidURL,
			item: domain.ShortURL{
				ID:        validId,
				Upstream:  *URLBadSchema,
				CreatedBy: validAuthor,
				CreatedAt: time.Now(),
				Enabled:   true,
			},
		},
		{
			err: domain.ErrInvalidURL,
			item: domain.ShortURL{
				ID:        validId,
				Upstream:  *URLBadHost,
				CreatedBy: validAuthor,
				CreatedAt: time.Now(),
				Enabled:   true,
			},
		},
		{
			err: domain.ErrInvalidAuthor,
			item: domain.ShortURL{
				ID:        validId,
				Upstream:  *validURL,
				CreatedBy: "free text",
				CreatedAt: time.Now(),
				Enabled:   true,
			},
		},
		{
			err: domain.ErrCreatedInFuture,
			item: domain.ShortURL{
				ID:        validId,
				Upstream:  *validURL,
				CreatedBy: validAuthor,
				CreatedAt: time.Now().Add(time.Hour),
				Enabled:   true,
			},
		},
		{
			err: domain.ErrEmptyTime,
			item: domain.ShortURL{
				ID:        validId,
				Upstream:  *validURL,
				CreatedBy: validAuthor,
				Enabled:   true,
			},
		},
	}

	for _, testCase := range cases {
		assert.ErrorContains(t, ValidateShortURL(testCase.item), testCase.err.Error())
	}
}
