package domain

import (
	"errors"
)

var (
	ErrEmptyId           = errors.New("empty URL identifier")
	ErrInvalidId         = errors.New("invalid URL identifier")
	ErrInvalidURL        = errors.New("invalid, insecure or empty URL")
	ErrInvalidAuthor     = errors.New("invalid author")
	ErrEmptyAuthor       = errors.New("empty author")
	ErrEmptyTime         = errors.New("empty time")
	ErrCreatedInFuture   = errors.New("creation dates in the future are not accepted")
	ErrURLNotFound       = errors.New("url not found")
	ErrURLTooLong        = errors.New("URL is too long")
	ErrCannotUseDisabled = errors.New("URL exist but is currently disabled")
	ErrUnavailableRepo   = errors.New("unavailable repository")
	ErrRepoSchema        = errors.New("repository anticorruption layer is erroring")
)
