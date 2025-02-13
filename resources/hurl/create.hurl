POST http://127.0.0.1:8080/v1/urls/short
Authorization: example
{
  "full_url": "https://opentelemetry.io/"
}

HTTP 201

[Captures]
short_url: jsonpath "$['short_url']"