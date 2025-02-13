# Enable .env file support for local configuration
set dotenv-load

# Use bash with strict error checking
set shell := ["bash", "-uc"]

alias t := test
alias b := build
alias r := run

out_dir  := "dist"
docs_dir := "docs"
mocks_dir := "mocks"
service_name := env_var_or_default("OTEL_SERVICE_NAME", "unknown")
dataset_items := "1000"

export GOPATH := env_var_or_default("GOPATH", `go env GOPATH`)
export GOOS := env_var_or_default("GOOS", `go env GOOS`)
export GOARCH := env_var_or_default("GOARCH", `go env GOARCH`)
export GOROOT := env_var_or_default("GOROOT", `go env GOROOT`)
#export CGO_ENABLED := env_var_or_default("CGO_ENABLED", "0")

clean:
  rm -rf "{{out_dir}}"
  mkdir -p "{{out_dir}}"

build: clean
  go build -o "{{out_dir}}" ./...

run: build
  {{out_dir}}/http

build-container:
  docker build -t "{{service_name}}:latest" .

run-container: build-container
  docker run \
    --rm -it --net host \
    -e AWS_ACCESS_KEY_ID \
    -e AWS_SECRET_ACCESS_KEY \
    -e AWS_ENDPOINT_URL_DYNAMODB \
    -e AWS_REGION \
    -e OTEL_SERVICE_NAME \
    -e OTEL_SERVICE_VERSION \
    -e OTEL_EXPORTER_OTLP_ENDPOINT \
    -e OTEL_EXPORTER_OTLP_INSECURE \
    -e OTEL_EXPORTER_OTLP_TIMEOUT \
    -e OTEL_EXPORTER_OTLP_PROTOCOL \
    -e OTEL_RESOURCE_ATTRIBUTES \
    -e SHORTENER_API_TOKEN \
    "{{service_name}}:latest"

fmt:
	go fmt ./...
	fd '\.go' -x gofumpt -w '{}'

lint:
	golangci-lint run --out-format checkstyle --issues-exit-code 0

test:
  go test -json > report.json -cover -coverprofile=coverage.out -race ./...

test_local:
  go test -cover ./...

coverage:
  # enable -race?
  go test -coverprofile=coverage.out $(go list ./... | rg -v 'mocks|cmd/http')
  go tool cover -html=coverage.out

coverage-lcov: coverage
  gcov2lcov -infile=coverage.out -outfile=lcov.info

coverwatch:
  watchexec -vvw 'coverage.out' "gcov2lcov -infile=coverage.out -outfile=lcov.info"

vet:
  go vet ./...

mocks:
  rm -rf "{{mocks_dir}}"
  mockery

docs:
    mkdir -p {{docs_dir}}
    go doc -all > {{docs_dir}}/API.md

scylla:
  docker run \
    --rm -it \
    --name scylla \
    --net host \
    -e SCYLLA_DEVELOPER_MODE=1 \
    docker.io/scylladb/scylla:6.1.5 \
    --listen-address=127.0.0.1 \
    --rpc-address=127.0.0.1 \
    --seed-provider-parameters=seeds=127.0.0.1 \
    --alternator-write-isolation=always \
    --alternator-address=127.0.0.1 \
    --alternator-port=8000 \
    --memory 5Gi

ddb-run:
  docker run \
    --rm -it \
    -p 8000:8000 \
    -p 4040:4040 \
    docker.io/amazon/dynamodb-local:latest
  # -sharedDb -inMemory

ddb-list:
  aws dynamodb list-tables \
    --profile local \
    --endpoint-url http://localhost:8000

ddb-create:
  aws dynamodb create-table \
    --profile local \
    --endpoint-url http://127.0.0.1:8000 \
    --cli-input-json file://resources/table_schema.json

ddb-makeprofile:
  aws configure set region dummyRegion --profile local;
  aws configure set output json --profile local;
  aws configure set aws_access_key_id dummyAccessKey --profile local;
  aws configure set aws_secret_access_key dummySecretKey --profile local;

perftest-dataset:
  mkdir -p dist
  seq 1 {{dataset_items}} | \
    xargs -I '{}' hurl resources/hurl/create.hurl | jq -r .short_url |\
    awk '{ print "\""$0"\""}'  | jq -s | tee "{{out_dir}}/dataset.json"