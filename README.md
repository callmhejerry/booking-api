# Booking API

A full-stack appointment booking application. The repository contains a Go
REST API backed by PostgreSQL and a Flutter web client. Users can find a
provider's available slots, create a booking, look up bookings, and cancel a
booking.

## Features

- Search appointment slots by provider.
- Create a booking for an available future slot.
- List bookings for an email address with cursor-based pagination.
- Retrieve a booking by its UUID.
- Cancel a booking and make its slot available again.
- Create slots and replace the complete slot schedule through API endpoints.
- Protect API routes with an API key, apply CORS restrictions, and rate-limit
  booking creation.
- Deploy the API, PostgreSQL database, and static web frontend on Render using
  the repository's Blueprint configuration.

## Architecture

```text
booking-api/
├── cmd/api/                 # API process entry point
├── internal/
│   ├── bookings/            # Booking controller, service, repository, DTOs
│   ├── slots/               # Slot controller, service, repository, DTOs
│   ├── middleware/          # API-key auth, CORS, and rate limiting
│   ├── config/              # Environment-based configuration
│   └── db/                  # PostgreSQL connection and generated queries
├── sql/
│   ├── migrations/          # Database schema migrations
│   └── queries/             # SQL query sources
├── tests/
│   ├── e2e/                 # Go API end-to-end tests
│   └── helpers/             # Test server and database setup
├── booking_web/             # Flutter web client
├── docker-compose.yml       # Local PostgreSQL service
├── Makefile                 # Common local development commands
└── render.yaml              # Render Blueprint for API, web, and database
```

The API is organized around controllers, services, and repositories. Controllers
bind HTTP requests and write responses; services apply booking and slot rules;
repositories use PostgreSQL transactions and queries. The Flutter client calls
the API over HTTP. Database migrations run when the API starts.

## Requirements

- Go **1.27.2** (the version declared in `go.mod`).
- Docker with Docker Compose, or an available PostgreSQL database.
- Flutter/Dart for the web client. The Render build uses Flutter **3.47.7**;
  the Dart SDK constraint is specified in `booking_web/pubspec.yaml`.
- `make` is optional but provides the shortcuts listed below.

## Run locally

### 1. Configure the API and database

Create a local environment file from the example:

```sh
cp .env.example .env
```

Set the following values in `.env` for the supplied Docker Compose database:

```dotenv
DATABASE_URL=
POSTGRES_USER=booking
POSTGRES_PASSWORD=change-this-local-password
POSTGRES_DB=booking
POSTGRES_PORT=7001
POSTGRES_HOST=localhost
POSTGRES_SSL_MODE=disable

PORT=8080
MODE=debug
API_KEY=replace-with-a-local-development-key
CORS_ALLOWED_ORIGINS=http://localhost:5000
```

Use your own local API key and password; do not commit real credentials. The
Compose service maps PostgreSQL to host port `7001`. If you use another
database, configure either `DATABASE_URL` or all of `POSTGRES_USER`,
`POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_HOST`, and `POSTGRES_PORT`.
`DATABASE_URL`, when set, takes precedence over the individual connection
fields.

Start PostgreSQL:

```sh
docker compose up -d
```

Or, if the environment file is configured as above:

```sh
make up
```

### 2. Start the API

```sh
go run ./cmd/api
```

The API listens on the configured `PORT` (for example, `http://localhost:8080`).
On startup it applies the SQL migrations and connects to PostgreSQL. Stop the
database when finished with `docker compose down` or `make down`. The Compose
database uses a named volume, so its data remains between restarts.

### 3. Start the web client

In another terminal:

```sh
cd booking_web
flutter pub get
flutter run -d web-server --web-hostname 0.0.0.0 --web-port 5000
```

Open [http://localhost:5000](http://localhost:5000). Enter the API base URL
(`http://localhost:8080`), the same API key configured for the API, and a
provider name, then select **Find slots**.

The API's `CORS_ALLOWED_ORIGINS` must contain the exact origin used by the
browser. The default in the API is `http://localhost:5000`. Origins can be
comma-separated. If you run the web client on another port, update the API
configuration and restart it.

## API reference

All routes require the `X-API-Key` request header, including reads:

```http
X-API-Key: your-api-key
```

Responses are JSON. API errors use an object shaped like
`{"error":"message"}`. Timestamps are RFC 3339 strings, and booking IDs are
UUIDs.

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/slots?provider={name}` | List slots for a provider, including availability. |
| `POST` | `/slots` | Add one slot. |
| `POST` | `/reset` | Delete and replace the complete slot schedule. |
| `POST` | `/bookings` | Book an available slot. |
| `GET` | `/bookings?email={email}&limit={n}&cursor={cursor}` | List bookings for an email address. |
| `GET` | `/bookings/{id}` | Retrieve a booking by UUID. |
| `DELETE` | `/bookings/{id}` | Cancel a booking and reopen its slot. |

### Search for slots

```sh
curl --get "http://localhost:8080/slots" \
  --data-urlencode "provider=Dr Example" \
  -H "X-API-Key: $API_KEY"
```

Example response:

```json
{
  "slots": [
    {
      "id": "slot-1",
      "provider": "Dr Example",
      "start": "2027-01-15T09:00:00Z",
      "end": "2027-01-15T09:30:00Z",
      "available": true
    }
  ]
}
```

### Create a slot

Both `start` and `end` must be RFC 3339 timestamps. Choose a future start time;
bookings are rejected after a slot has started. Slot IDs must be unique.

```sh
curl -X POST "http://localhost:8080/slots" \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "slot-1",
    "provider": "Dr Example",
    "start": "2027-01-15T09:00:00Z",
    "end": "2027-01-15T09:30:00Z"
  }'
```

### Create a booking

The `id` value is the slot's ID, not a booking UUID. The API accepts an
optional `Idempotency-Key` header for retry-safe booking requests.

```sh
curl -X POST "http://localhost:8080/bookings" \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: demo-request-001" \
  -d '{
    "id": "slot-1",
    "name": "Alex Example",
    "email": "alex@example.com"
  }'
```

Successful creation returns HTTP `201` and a booking object:

```json
{
  "id": "booking-uuid",
  "slot_id": "slot-1",
  "name": "Alex Example",
  "email": "alex@example.com",
  "status": "confirmed",
  "created_at": "2027-01-01T12:00:00Z"
}
```

Booking creation is rate-limited per client IP. A slot can only be booked once
while it is available. The database also prevents duplicate bookings for the
same email and slot.

### List bookings

The default page size is 10; `limit` can be between 1 and 50. When more results
are available, the response includes `next_cursor`. Pass that value back as the
`cursor` query parameter to retrieve the next page.

```sh
curl --get "http://localhost:8080/bookings" \
  --data-urlencode "email=alex@example.com" \
  --data-urlencode "limit=10" \
  -H "X-API-Key: $API_KEY"
```

### Retrieve or cancel a booking

Use the booking UUID returned when the booking was created:

```sh
curl "http://localhost:8080/bookings/booking-uuid" \
  -H "X-API-Key: $API_KEY"

curl -X DELETE "http://localhost:8080/bookings/booking-uuid" \
  -H "X-API-Key: $API_KEY"
```

Cancellation updates the booking status to `cancelled` and makes the associated
slot available again.

### Replace the complete slot schedule

`POST /reset` accepts an object containing a `slots` array. It deletes all
existing slots before creating the supplied schedule. Because bookings reference
slots with `ON DELETE CASCADE`, this also deletes related bookings. Treat this
as a destructive administrative operation; do not use it against production
data unless that is explicitly intended.

```sh
curl -X POST "http://localhost:8080/reset" \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "slots": [
      {
        "id": "slot-2",
        "provider": "Dr Example",
        "start": "2027-01-16T10:00:00Z",
        "end": "2027-01-16T10:30:00Z"
      }
    ]
  }'
```

The web client's reset panel presents a confirmation dialog, but the API route
itself does not require a separate confirmation token.

## Configuration

| Variable | Required | Description |
| --- | --- | --- |
| `API_KEY` | Yes | Shared API key expected in the `X-API-Key` header. |
| `PORT` | Yes | Port for the Go HTTP server. Render provides this in deployment. |
| `MODE` | No | Gin mode; use `release` for production. |
| `CORS_ALLOWED_ORIGINS` | No | Comma-separated allowed browser origins. Defaults to `http://localhost:5000`. |
| `DATABASE_URL` | One DB option | PostgreSQL connection URL. If set, takes precedence over the individual `POSTGRES_*` settings. |
| `POSTGRES_USER` | One DB option | PostgreSQL username. |
| `POSTGRES_PASSWORD` | One DB option | PostgreSQL password. |
| `POSTGRES_DB` | One DB option | Database name. |
| `POSTGRES_HOST` | One DB option | PostgreSQL host. |
| `POSTGRES_PORT` | One DB option | PostgreSQL port. |
| `POSTGRES_SSL_MODE` | No | PostgreSQL TLS mode. Release mode requires TLS; if omitted in release mode, the API defaults to `require`. |

The API loads `.env` when present. Environment variables supplied by the
runtime can also configure it. Keep credentials out of source control. The web
client takes the API URL and key in its UI; it holds the key in page memory and
does not persist it. A key distributed to browser users must be considered
public and should not be treated as a secret for privileged operations.

## Tests and checks

### Go tests

The end-to-end tests use `.env.test` and a PostgreSQL database. Configure that
file with a **separate test database** before running the suite; the test setup
applies migrations and truncates the `slots` and `bookings` tables during
cleanup.

```sh
go test ./...
```

Or run the repository's end-to-end test target:

```sh
make test
```

### Flutter tests and analysis

```sh
cd booking_web
flutter test
flutter analyze
```

To build the static web app locally:

```sh
cd booking_web
flutter build web --release
```

## Deployment on Render

The root [`render.yaml`](./render.yaml) defines a Render Blueprint containing:

- A Go web service named `booking-api`.
- A Flutter static site named `booking-web`.
- A PostgreSQL database named `booking-api-db`.

Create or update a Render Blueprint from this repository and deploy the
services. The API service builds with `go build` and starts the resulting
binary. The web service checks out the pinned Flutter SDK, builds
`booking_web` for web, and publishes `booking_web/build/web`.

The Blueprint generates an API key and connects the API service to the managed
database. After deployment:

1. Confirm the API and static site have finished deploying.
2. Set the API service's `CORS_ALLOWED_ORIGINS` to the exact public origin of
   the static site if it differs from the value in `render.yaml`.
3. Open the web app and enter the API service's public HTTPS URL and the API
   key from the API service's environment settings.
4. Search for a provider with available slots.

Do not put the API key into the Flutter build or commit it to the repository.
For more frontend-specific setup and deployment notes, see
[`booking_web/README.md`](./booking_web/README.md).
