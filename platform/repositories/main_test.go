package repositories

import "net/url"

var (
	URLBadSchema, _ = url.Parse("http://opentelemetry.io")
	URLBadHost, _   = url.Parse("https://")
	validURL, _     = url.Parse("https://opentelemetry.io")
	validAuthor     = "root@neonmei.cloud"
	validId         = "asd"
)
