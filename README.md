# Cloneslist: Craiglist like backend server with location based post searching

## Prerequisites:

1. Have the Go programming language installed
2. Install PostgreSQL and PostGIS and have a database created.
3. Create a free Geocodio account and get an API key. You can make an account [here](https://dash.geocod.io/register).

## Setup

These commands assume a Unix like system

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

You can create a secret with this command.

```bash
openssl rand -base64 64
```

## Build and run server

```bash
go build -o cloneslist && ./cloneslist
```
