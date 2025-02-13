package dtos

type URLCreateRequest struct {
	Upstream string `json:"full_url"`
}

type URLCreateResponse struct {
	ShortURL string `json:"short_url"`
}
