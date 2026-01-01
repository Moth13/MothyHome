air:
	air

server: 
	go run cmd/main.go

make build-docker:
	docker build -t home_client:latest .

PHONY: air server