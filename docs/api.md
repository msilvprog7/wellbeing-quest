# API

Placeholder for documentation for the API services.

## Setup

The API services use the following:

- Go
- Gin web service framework

## Run

Run the service on `http://localhost:8080`:

```cmd
cd api
go mod tidy
go run ./cmd/service
```

Send requests via curl:

```cmd
curl http://localhost:8080/activities/v1 \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"name": "Meditate","feelings": ["Relaxed"]}'

curl http://localhost:8080/activities/v1/weeks/2025-08-10

curl http://localhost:8080/activities/v1/suggestions
```
