package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neonmei/challenge_urlshortener/application"
)

var readyProbe = true

func handleHealth(_ application.Service, c *gin.Context) {
	if readyProbe {
		c.String(http.StatusOK, "OK")
	} else {
		c.String(http.StatusServiceUnavailable, "Shutting down")
	}
}
