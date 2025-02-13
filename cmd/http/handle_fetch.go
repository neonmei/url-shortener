package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neonmei/challenge_urlshortener/application"
	"github.com/neonmei/challenge_urlshortener/domain"
	"github.com/neonmei/challenge_urlshortener/domain/validators"
	"github.com/neonmei/challenge_urlshortener/platform/dtos"
)

func handleFetch(e application.Service, c *gin.Context) {
	urlId := c.Param("url_id")
	if err := validators.ValidateId(urlId); err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	item, err := e.Fetch(c.Request.Context(), urlId)
	if err == nil {
		c.JSON(http.StatusOK, dtos.FromDomain(*item))
		return
	}

	if errors.Is(err, domain.ErrURLNotFound) {
		c.Status(http.StatusNotFound)
		return
	}

	_ = c.Error(err)
	c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: err.Error()})
}
