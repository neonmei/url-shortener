ARG GO_VERSION=1.23-alpine3.20
ARG CGO_ENABLED=0
ARG GOOS=linux

FROM docker.io/golang:$GO_VERSION AS compiler
WORKDIR /src
COPY . .
RUN apk add --no-cache ca-certificates tzdata git
RUN mkdir -p bin && \
    go mod download -x && \
    go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./bin/main ./cmd/*/*.go

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --chown=nonroot:nonroot --from=compiler /src/bin/main /app/main
COPY --chown=nonroot:nonroot assets /app/assets

ENTRYPOINT ["/app/main"]