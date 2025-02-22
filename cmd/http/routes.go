package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neonmei/challenge_urlshortener/application"
	"github.com/neonmei/challenge_urlshortener/platform/config"
)

const (
	UserContextKey = "auth.user"
)

func routes(apiRouter *gin.Engine, cfg config.AppConfig, e application.Service) {
	// Public endpoints /v1/urls/redirect/:url_id
	apiRouter.GET("/:url_id", func(ctx *gin.Context) { handleRedirect(e, ctx) })

	// Administrative endpoints
	groupUrls := apiRouter.Group("/v1/urls").Use(TokenAuthMiddleware(cfg))
	groupUrls.POST("/short", func(ctx *gin.Context) { handleCreate(e, ctx) })
	groupUrls.DELETE("/short/:url_id", func(ctx *gin.Context) { handleDelete(e, ctx) })
	groupUrls.GET("/short/:url_id", func(ctx *gin.Context) { handleFetch(e, ctx) })

	// Platform endpoints
	apiRouter.GET("/platform/healthz", func(ctx *gin.Context) { handleHealth(e, ctx) })

	apiRouter.LoadHTMLFiles(
		fmt.Sprintf("assets/%s", StatusNotFoundTemplate),
		fmt.Sprintf("assets/%s", InternalServiceErrorTemplate),
	)
}

func TokenAuthMiddleware(cfg config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerToken := c.Request.Header.Get("Authorization")

		if headerToken != cfg.ApiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		c.Set(UserContextKey, cfg.ApiUser)
		c.Next()
	}
}
