package domain

import (
	"net/url"
	"time"
)

type ShortURL struct {
	// ID is an alphanumeric identifier
	ID string

	// Upstream is the original address where content reside
	Upstream url.URL

	// CreatedBy is an RFC 5322 compliant email address identifying creator
	CreatedBy string

	// CreatedAt indicates creation time
	CreatedAt time.Time

	// Enabled flags if current URL is active
	Enabled bool
}
