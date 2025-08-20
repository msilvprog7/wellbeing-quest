# API

Placeholder for documentation for the API services.

## Setup

The API services use the following:

- Go
- Go dot env
- Gin web service framework
- Docker
- Postgres (localhost)

Run postgres from a docker container, specifying
a password:

```cmd
docker run --name pg-local \
  -e POSTGRES_PASSWORD= \
  -e POSTGRES_DB=yourdb \
  -p 5432:5432 \
  -d postgres
```

Create file `.env`, specifying the password
entered above:

```text
DB_MODE=localhost
DB_RESET=reset
DB_DRIVER=postgres
DB_USER=postgres
DB_PASSWORD=
DB_NAME=yourdb
```

## Run

Run the service on `http://localhost:8080`:

```cmd
cd api
go mod tidy
go run ./cmd/service
```

Send requests via curl or postman collection
in `examples/Wellbeing Quest.postman_collection.json`:

```cmd
curl http://localhost:8080/activities/v1 \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"name": "Meditate","feelings": ["Relaxed"]}'

curl http://localhost:8080/activities/v1/weeks/2025-08-10

curl http://localhost:8080/activities/v1/suggestions
```
