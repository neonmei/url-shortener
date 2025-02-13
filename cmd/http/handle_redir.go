package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neonmei/challenge_urlshortener/application"
	"github.com/neonmei/challenge_urlshortener/domain/validators"
)

const (
	StatusNotFoundTemplate       = "assets/404.html"
	InternalServiceErrorTemplate = "assets/500.html"
)

func handleRedirect(e application.Service, c *gin.Context) {
	urlId := c.Param("url_id")
	if err := validators.ValidateId(urlId); err != nil {
		_ = c.Error(err)
		c.File(StatusNotFoundTemplate)
		c.Status(http.StatusNotFound)
	}

	newURL, err := e.Redirect(c.Request.Context(), urlId)
	if err == nil {
		c.Redirect(http.StatusFound, newURL)
		return
	}

	_ = c.Error(err)
	c.File(InternalServiceErrorTemplate)
	c.Status(http.StatusNotFound)
}
