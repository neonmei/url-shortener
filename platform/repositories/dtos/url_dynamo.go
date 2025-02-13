package dtos

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/neonmei/challenge_urlshortener/domain"
	"github.com/neonmei/challenge_urlshortener/domain/validators"
)

const DynamoTimeFormat = time.RFC3339

type URLItem struct {
	Id      string `dynamodbav:"url_id"`
	Created string `dynamodbav:"created_at"`
	Author  string `dynamodbav:"created_by"`
	Enabled bool   `dynamodbav:"enabled"`
	FullURL string `dynamodbav:"full_url"`
}

func FromDomain(u domain.ShortURL) URLItem {
	return URLItem{
		Id:      u.ID,
		Created: u.CreatedAt.Format(DynamoTimeFormat),
		Author:  u.CreatedBy,
		Enabled: u.Enabled,
		FullURL: u.Upstream.String(),
	}
}

func (i URLItem) Domain() (*domain.ShortURL, error) {
	t, err := time.Parse(DynamoTimeFormat, i.Created)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("cannot parse time dynamodb"), err)
	}

	u, err := url.Parse(i.FullURL)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("cannot parse URL"), err)
	}

	shortUrl := domain.ShortURL{
		ID:        i.Id,
		CreatedBy: i.Author,
		Enabled:   i.Enabled,
		Upstream:  *u,
		CreatedAt: t,
	}

	if err := validators.ValidateShortURL(shortUrl); err != nil {
		return nil, err
	}

	return &shortUrl, err
}
