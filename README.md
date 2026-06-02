# Cloneslist

Craigslist like backend server with location based post searching.

This Project is WIP.

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

The address must be valid

```bash
curl -X POST "http://localhost:8080/user/signup" -d '{"email": "example@email.com", "password": "somepassword", "address": "1234 Some St", "city": "Towntown", "state": "CA", "zip": "12345"}'
```

Response:

```json
{
  "id": "7dc6d38d-0d22-463d-9929-a13572bcb00c",
  "created_at": "2026-06-01T17:05:40.784057Z",
  "updated_at": "2026-06-01T17:05:40.784057Z",
  "email": "example@email.com"
}
```

### Login

```bash
curl -X POST "http://localhost:8080/user/login" -d '{"email": "example@email.com", "password": "somepassword"}'
```

Response:

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

### Update user account

```bash
curl -X PUT "http://localhost:8080/user/update" -d '{"email": "example2@email.com", "password": "newpassword", "address": "1234 Someother St", "city": "Towntwo", "state": "CA", "zip": "54321"}' -H "Authorization: Bearer <token>"
```

Updating the address, city, state, and zip is optional.
The address must be valid.

Response:

```json
{
  "id": "7dc6d38d-0d22-463d-9929-a13572bcb00c",
  "created_at": "2026-06-01T17:05:40.784057Z",
  "updated_at": "2026-06-01T17:49:31.237689Z",
  "email": "example2@email.com"
}
```

### Get user info

```bash
curl -X GET "http://localhost:8080/user/<user_id>"
```

Response:

```json
{
  "id": "7dc6d38d-0d22-463d-9929-a13572bcb00c",
  "email": "example2@email.com"
}
```

### Delete user account

```bash
curl -X DELETE "http://localhost:8080/user/delete" -d '{"user_id": <user_id>}' -H "Authorization: Bearer <token>"
```

### Get new token using refresh token

```bash
curl -X POST "http://localhost:8080/api/refresh" -H "Authorization: Bearer <refresh_token>"
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXN"
}
```

### Revoke refresh token

```bash
curl -X POST "http://localhost:8080/api/revoke" -H "Authorization: Bearer <refresh_token>"
```

### Create post

```bash
curl -X POST "http://localhost:8080/user/post" -d '{"title": "This is an item I am trying to sell", "description": "This is my description", "price": 20.99, "category": "forsale"}' -H "Authorization: Bearer <token>"
```

Response:

```json
{
  "id": "4419e35e-59b7-43c5-8b3a-f5c3000ab897",
  "user_id": "7dc6d38d-0d22-463d-9929-a13572bcb00c",
  "title": "This is an item I am trying to sell",
  "description": "This is my description",
  "price": 20.99,
  "created_at": "2026-06-01T18:03:29.516589Z",
  "updated_at": "2026-06-01T18:03:29.516589Z"
}
```

Available categories are:

1. forsale
2. housing
3. jobs
4. services
5. community

### Update post

```bash
curl -X PUT "http://localhost:8080/user/post/<post_id>" -d '{"title": "New post title", "description": "New post description", "price": 900.99, "status": "sold", "category": "housing"}' -H "Authorization: Bearer <token>"
```

Response:

```json
{
  "id": "4419e35e-59b7-43c5-8b3a-f5c3000ab897",
  "user_id": "7dc6d38d-0d22-463d-9929-a13572bcb00c",
  "title": "New post title",
  "description": "New post description",
  "price": 900.99,
  "created_at": "2026-06-01T18:03:29.516589Z",
  "updated_at": "2026-06-01T18:17:34.237167Z",
  "status": "sold"
}
```

Category is optional. The available ones are:

1. forsale
2. housing
3. jobs
4. services
5. community

The available statuses are:

1. active
2. sold

### Delete post

```bash
curl -X DELETE "http://localhost:8080/user/post/<post_id>" -H "Authorization: Bearer <token>"
```

### Search posts

The only thing absolutely needed to search posts is a location. You can set the location in 2 ways.

1. Utilize the location used at account creation by setting the auth header.

```bash
curl -X GET "http://localhost:8080/posts/search" -H "Authorization: Bearer <token>"
```

2. Or set the location in the URL.

```bash
curl -X GET "http://localhost:8080/posts/search?city=Towntown&state=CA"
```

In the example above, the location is city=Towntown and state=CA.

For the sake of simplicity, the following example sets the location via the URL, but you can use the auth header instead if you like.

```bash
curl -X GET "http://localhost:8080/posts/search?city=Towntown&state=CA&distance=100&category=forsale&s=item&sort=pricedesc"
```

Here we see the following fields set.

1. distance (Distance from search location in miles. In this case "100". Default is 50)
2. category (Category of the posts. In this case "forsale")
3. s (Search term to look for in post title and description. In this case "item")
4. sort (Sorting of results. In this case "pricedesc". Default is "timedesc")

You have to separate the fields with a “&”. All of these are optional so you can mix and match whatever you need.

Available categories are:

1. forsale
2. housing
3. jobs
4. services
5. community

Available sorting styles are:

1. timedesc (Newest first)
2. timeasc (Oldest first)
3. pricedesc (Lowest first)
4. priceasc (Highest first)

Response:

```json
[
  {
    "id": "083fe225-3f08-4c14-90a0-1c4807bf480a",
    "user_id": "7904c467-aac3-4857-b1da-0e7e53025b9c",
    "title": "This is an item I am trying to sell",
    "description": "This is my description",
    "price": 20.99,
    "created_at": "2026-06-01T18:36:58.258378Z",
    "updated_at": "2026-06-01T18:36:58.258378Z",
    "status": "active"
  },
  {
    "id": "07289346-d0c8-4405-8da7-4a48094b1365",
    "user_id": "7dc6d38d-0d22-463d-9929-a13572bcb00c",
    "title": "This is another post of an item I am trying to sell",
    "description": "This is another description",
    "price": 400.99,
    "created_at": "2026-06-01T18:52:21.223938Z",
    "updated_at": "2026-06-01T18:52:21.223938Z",
    "status": "active"
  }
]
```

### Get post info

```bash
curl -X GET "http://localhost:8080/posts/<post_id>"
```

Response:

```json
{
  "id": "07289346-d0c8-4405-8da7-4a48094b1365",
  "user_id": "7904c467-aac3-4857-b1da-0e7e53025b9c",
  "title": "This is another post of an item I am trying to sell",
  "description": "This is another description",
  "price": 400.99,
  "created_at": "2026-06-01T18:52:21.223938Z",
  "updated_at": "2026-06-01T18:52:21.223938Z",
  "status": "active"
}
```

### Get all posts from a certain user

```bash
curl -X GET "http://localhost:8080/posts/user/<user_id>"
```

Response:

```json
[
  {
    "id": "083fe225-3f08-4c14-90a0-1c4807bf480a",
    "user_id": "7904c467-aac3-4857-b1da-0e7e53025b9c",
    "title": "This is an item I am trying to sell",
    "description": "This is my description",
    "price": 20.99,
    "created_at": "2026-06-01T18:36:58.258378Z",
    "updated_at": "2026-06-01T18:36:58.258378Z",
    "status": "active"
  },
  {
    "id": "8bdbd56a-d14c-4491-884b-ebb17f1fc374",
    "user_id": "7904c467-aac3-4857-b1da-0e7e53025b9c",
    "title": "This is another post",
    "description": "This is another description",
    "price": 400.99,
    "created_at": "2026-06-01T18:39:23.417677Z",
    "updated_at": "2026-06-01T18:39:23.417677Z",
    "status": "active"
  }
]
```
