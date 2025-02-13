package dtos

import "github.com/neonmei/challenge_urlshortener/domain"

type URLFetchResponse struct {
	URL       string `json:"full_url"`
	Enabled   bool   `json:"enabled"`
	CreatedAt int64  `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

func FromDomain(item domain.ShortURL) URLFetchResponse {
	return URLFetchResponse{
		URL:       item.Upstream.String(),
		Enabled:   item.Enabled,
		CreatedAt: item.CreatedAt.Unix(),
		CreatedBy: item.CreatedBy,
	}
}
