# Cloneslist

Craigslist like backend server with location based post searching.

## Prerequisites:

1. Have the Go programming language installed
2. Install PostgreSQL and PostGIS and have a database created.
3. Create a free Geocodio account and get an API key. You can create an account [here](https://dash.geocod.io/register).
4. Install the Goose database migration tool. You can do so with the following command.

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Setup

All these commands assume a Unix-like system

```bash
git clone https://github.com/Eval-99/cloneslist/
```

```bash
cd cloneslist
```

```bash
go mod tidy
```

```bash
touch .env
```

Fill out the .env file. Here is an example.

```bash
DB_URL="postgres://<user_name>:<password>@localhost:<port>/<database_name>?sslmode=disable"
SECRET="X7mQ2vL9aK4pT8cR1jW6yH3sD5uF0zN8bE2nG7kC4qM1rY9tP6hJ3xV5wS0dL8uA2iO7eR4f"
GEOCODIO_API_KEY="K7mQ2xN9pL4vT8cR1jW6yH3sD5uF0aZkE2nB7gC"
```

For the URL, add the PostgreSQL username, password, port, and database name.

You can create a secret with this command.

```bash
openssl rand -base64 64
```

Run database migrations.

```bash
goose postgres <postgres_connection_url> up
```

You can use the same URL as earlier.

## Build and run server

```bash
go build -o cloneslist && ./cloneslist
```

## How to use

### Create user account.

```bash
curl -X POST "http://localhost:8080/user/signup" -d '{"email": "example@email.com", "password": "somepassword", "address": "1234 Some St", "city": "Towntown", "state": "CA", "zip": "12345"}'
```

The response will be something like this.

```json
{
  "id": "7dc6d38d-0d22-463d-9929-a13572bcb00c",
  "created_at": "2026-06-01T17:05:40.784057Z",
  "updated_at": "2026-06-01T17:05:40.784057Z",
  "email": "example@email.com",
  "token": "",
  "refresh_token": ""
}
```

### Login

```bash
curl -X POST "http://localhost:8080/user/login" -d '{"email": "example@email.com", "password": "somepassword"}'
```

The response will be something like this.

```json
{
  "id": "7dc6d38d-0d22-463d-9929-a13572bcb00c",
  "created_at": "2026-06-01T17:05:40.784057Z",
  "updated_at": "2026-06-01T17:05:40.784057Z",
  "email": "example@email.com",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXN",
  "refresh_token": "c5802f6eb17743cc81e4fe102973d904e37429fe401d0629535beb6148074324"
}
```
