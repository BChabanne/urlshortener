FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /urlshortener

FROM ubuntu

EXPOSE 8000
USER 1000
CMD ["/urlshortener", "-addr", "0.0.0.0:8000"]

COPY --from=builder /urlshortener /urlshortener
