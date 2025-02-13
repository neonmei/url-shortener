package semconv

const (
	UrlId      = "short_url.id"
	UrlAuthor  = "short_url.author"
	UrlEnabled = "short_url.enabled"
	UrlFull    = "short_url.full"
	UrlCreated = "short_url.created"

	HasherRounds = "hasher.rounds"
	HasherLength = "hasher.length"

	CacheHit      = "cache.access.hit"
	CacheMiss     = "cache.access.miss"
	CacheAdded    = "cache.keys.added"
	CacheEvicted  = "cache.keys.evicted"
	CacheRejected = "cache.keys.rejected"
)

const (
	MetricURLHits = "meli.shortener.url.hits"
)
