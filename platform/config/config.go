package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	// Port is the HTTP server port
	Port int `split_words:"true" default:"8080" `

	// TraceIdSampleRatio enables extra troubleshooting information
	TraceIdSampleRatio float64 `split_words:"true" default:"0" `

	// BaseUrl is the host of the service
	BaseUrl string `split_words:"true" default:"https://me.li" `

	// ApiKey implements an auth token for admin endpoint for demo purposes
	ApiKey string `split_words:"true" default:"example"`

	// ApiUser for demo purposes
	ApiUser string `split_words:"true" default:"root@neonmei.cloud"`

	// MaxLength is the HTTP server port
	MaxLength int `split_words:"true" default:"1024" `

	// ShutdownTimeout how much to wait for pending operations
	ShutdownTimeout time.Duration `split_words:"true" default:"5s" `

	// ShutdownWait how much to wait before initiating shutdown
	ShutdownWait time.Duration `split_words:"true" default:"60s" `

	Dynamo struct {
		// DynamoTableName sets where the storage backend will search for url data
		TableName string `split_words:"true" default:"url_shortener" `

		// ReadTimeout how much to wait for DynamoDB read operations
		ReadTimeout time.Duration `split_words:"true" default:"50ms" `

		// WriteTimeout how much to wait for DynamoDB write operations
		WriteTimeout time.Duration `split_words:"true" default:"900ms" `
	}

	Hasher struct {
		// RandomMaxValue sets the maximum value cap for the random hasher
		RandomMaxValue uint64 `split_words:"true" default:"3521614606207" `

		// MaxRounds indicates how many time to try hashing before giving up
		MaxRounds uint64 `split_words:"true" default:"4" `
	}

	Cache struct {
		// Counter is the number of keys to track frequency of
		NumCounters int64 `split_words:"true" default:"100000" `

		// MaxCost is the capacity limit of the cache in a arbitrary unit of cost
		MaxCost int64 `split_words:"true" default:"50000" `

		// BufferItems is number of keys per Get buffer
		BufferItems int64 `split_words:"true" default:"64" `

		// MetricsEnabled optionally enables metrics
		MetricsEnabled bool `split_words:"true" default:"false" `
	}
}

func Load() AppConfig {
	var result AppConfig
	err := envconfig.Process("shortener", &result)
	if err != nil {
		panic(err)
	}

	return result
}
