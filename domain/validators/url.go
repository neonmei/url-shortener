package validators

import (
	"errors"
	"net/mail"
	"net/url"
	"time"

	"github.com/neonmei/challenge_urlshortener/domain"
)

func ValidateURL(u *url.URL) error {
	if u == nil {
		return domain.ErrInvalidURL
	}

	if u.Scheme != "https" {
		return domain.ErrInvalidURL
	}

	if len(u.Host) < 1 {
		return domain.ErrInvalidURL
	}

	return nil
}

func ValidateAuthor(author string) error {
	if len(author) < 1 {
		return domain.ErrEmptyAuthor
	}

	// Validamos RFC 5322, por ahora no validamos dominios
	if _, err := mail.ParseAddress(author); err != nil {
		return errors.Join(domain.ErrInvalidAuthor, err)
	}

	return nil
}

func ValidateCreated(t time.Time) error {
	if (t == time.Time{}) {
		return domain.ErrEmptyTime
	}

	if t.After(time.Now().UTC()) {
		return domain.ErrCreatedInFuture
	}

	return nil
}

func ValidateId(u string) error {
	if len(u) < 1 {
		return domain.ErrEmptyId
	}

	for _, r := range u {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return domain.ErrInvalidId
		}
	}

	return nil
}

func ValidateShortURL(u domain.ShortURL) error {
	return errors.Join(
		ValidateAuthor(u.CreatedBy),
		ValidateCreated(u.CreatedAt),
		ValidateURL(&u.Upstream),
		ValidateId(u.ID),
	)
}
