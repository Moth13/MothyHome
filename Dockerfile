FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bravia-app main.go

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/bravia-app /app/bravia-app

EXPOSE 8080

ENTRYPOINT ["/app/bravia-app"]
