package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neonmei/challenge_urlshortener/application"
	"github.com/neonmei/challenge_urlshortener/domain"
	"github.com/neonmei/challenge_urlshortener/domain/validators"
)

const (
	StatusNotFoundTemplate       = "404.html"
	InternalServiceErrorTemplate = "500.html"
)

func handleRedirect(e application.Service, c *gin.Context) {
	urlId := c.Param("url_id")
	if err := validators.ValidateId(urlId); err != nil {
		c.HTML(http.StatusNotFound, StatusNotFoundTemplate, nil)
		_ = c.Error(err)
	}

	newURL, err := e.Redirect(c.Request.Context(), urlId)
	if err == nil {
		c.Redirect(http.StatusFound, newURL)
		return
	}

	if errors.Is(err, domain.ErrURLNotFound) {
		c.HTML(http.StatusNotFound, StatusNotFoundTemplate, nil)
		_ = c.Error(err)
		return
	}

	c.HTML(http.StatusInternalServerError, InternalServiceErrorTemplate, nil)
	_ = c.Error(err)
}
