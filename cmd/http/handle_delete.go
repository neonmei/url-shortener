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

func handleDelete(e application.Service, c *gin.Context) {
	urlId := c.Param("url_id")
	if err := validators.ValidateId(urlId); err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusNotFound, dtos.ErrorResponse{Error: err.Error()})
		return
	}

	err := e.Delete(c.Request.Context(), urlId)
	if err == nil {
		c.Status(http.StatusNoContent)
		return
	}

	if errors.Is(err, domain.ErrURLNotFound) {
		c.Status(http.StatusNotFound)
		return
	}

	_ = c.Error(err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
