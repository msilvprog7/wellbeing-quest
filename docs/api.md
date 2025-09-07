# API

Placeholder for documentation for the API services.

## Setup

The API services use the following:

- Go
- Go dot env
- Gin web service framework
- Docker
- Postgres (localhost)

Create file `.env`, specifying the password
entered above:

```text
DB_MODE=localhost
DB_RESET=
DB_DRIVER=postgres
DB_USER=postgres
DB_PASSWORD=
DB_HOST=db
DB_NAME=yourdb
DB_SQLDIRECTORY=cmd/service
```

## Run

Run the service on `http://localhost:8080`:

```cmd
docker-compose up --build
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

## Test

Run unit tests:

```cmd
go test ./tests/unittests
```

Create `tests/e2etests/.env`:

```text
DB_MODE=localhost
DB_RESET=reset
DB_DRIVER=postgres
DB_USER=postgres
DB_PASSWORD=
DB_HOST=db
DB_NAME=yourdb
DB_SQLDIRECTORY=../../cmd/service
```

Run e2e tests:

```cmd
go test ./tests/e2etests
```
