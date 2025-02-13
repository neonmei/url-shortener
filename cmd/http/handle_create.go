package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neonmei/challenge_urlshortener/application"
	"github.com/neonmei/challenge_urlshortener/platform/dtos"
)

func handleCreate(e application.Service, c *gin.Context) {
	createRequest := dtos.URLCreateRequest{}
	if err := json.NewDecoder(c.Request.Body).Decode(&createRequest); err != nil {
		_ = c.Error(errors.Join(ErrHttpRequestDecode, err))
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Error: ErrHttpRequestDecode.Error()})
		return
	}

	shortURL, err := e.Shorten(c.Request.Context(), createRequest.Upstream, c.GetString(UserContextKey))
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dtos.URLCreateResponse{
		ShortURL: shortURL.String(),
	})
}
