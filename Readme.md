# URL Shortener

This is a short url shortener writen with golang

## API

```
POST /
request {
    "url": "<url>"
}
response {
   "slug": "<slug>"
}

GET /:slug
response : 302 <url>
```

You can use go client with `NewClient` function

There is also a basic html ui available on listening addresse.

## Run

### With docker

Build and run using docker by running

```
docker-compose up --build
```

### With go

```
go run .
```

## Test

```
go test .
go test -bench=.
```
