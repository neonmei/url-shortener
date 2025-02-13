package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/gin-gonic/gin"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"github.com/neonmei/challenge_urlshortener/application"
	"github.com/neonmei/challenge_urlshortener/platform/clients"
	"github.com/neonmei/challenge_urlshortener/platform/config"
	"github.com/neonmei/challenge_urlshortener/platform/o11y"
	"github.com/neonmei/challenge_urlshortener/platform/repositories"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	cfg := config.Load()
	otelShutdown, err := otelconfig.ConfigureOpenTelemetry(buildOtelOpts(cfg)...)
	if err != nil {
		panic(err)
	}

	defer otelShutdown()

	// TODO: Sumar instrumentaciones async de OTel si se habilita metrics
	cache, err := ristretto.NewCache(&repositories.URLCacheConfig{
		NumCounters: cfg.Cache.NumCounters,
		MaxCost:     cfg.Cache.MaxCost,
		BufferItems: cfg.Cache.BufferItems,
		Metrics:     cfg.Cache.MetricsEnabled,
	})
	if err != nil {
		panic(err)
	}
	defer cache.Close()

	if cfg.Cache.MetricsEnabled {
		o11y.InstrumentCacheAsync(cache)
	}

	dynamoClient, err := clients.NewDynamoClient(cfg)
	if err != nil {
		panic(err)
	}

	urlRepository := repositories.NewCached(repositories.NewDynamoURLRepository(cfg, dynamoClient), cache)
	app, err := application.New(cfg, urlRepository)
	if err != nil {
		panic(app)
	}

	if err := gracefulServe(cfg, app); err != nil {
		slog.Error(err.Error())
	}
}

func gracefulServe(cfg config.AppConfig, e application.Service) error {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())

	apiAddress := fmt.Sprintf(":%d", cfg.Port)
	ginRouter.Use(otelgin.Middleware(os.Getenv("OTEL_SERVICE_NAME"), otelgin.WithFilter(func(r *http.Request) bool {
		return r.URL.Path != "/healthz"
	})))

	routes(ginRouter, cfg, e)

	httpServer := &http.Server{
		Addr:    apiAddress,
		Handler: ginRouter.Handler(),
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Handle SIGTERM/SIGINT signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	readyProbe = false
	time.Sleep(cfg.ShutdownWait)

	ctx, cancelFunc := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancelFunc()

	if err := httpServer.Shutdown(ctx); err != nil {
		return errors.Join(errors.New("Shutdown error"), err)
	}

	select {
	case <-ctx.Done():
		return errors.New("Shutdown timeout")
	}
}

func buildOtelOpts(cfg config.AppConfig) []otelconfig.Option {
	otelOpts := []otelconfig.Option{}
	if cfg.TraceIdSampleRatio > 0 {
		otelOpts = append(otelOpts, otelconfig.WithSampler(trace.TraceIDRatioBased(0.01)))
	}

	return otelOpts
}
