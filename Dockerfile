FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o home-client main.go

FROM scratch AS runner

WORKDIR /app
COPY --from=builder /app/home-client /app/home-client

EXPOSE 8080

ENTRYPOINT ["/app/home-client"]
