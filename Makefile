.PHONY: generate
generate:
	@echo "Generating proto files"
	protoc -I=./internal/protos --go_out=./internal/protos ./internal/protos/*.proto
	@echo "Generated proto files"
up:
	@echo "Spinning up containers"
	docker compose up --build -d
	@echo "Containers are up... maybe :thinking:"

down:
	@echo "Spinning down containers"
	docker compose down
	@echo "Containers are down... maybe :thinking:"

start-both:
	@echo "Starting both services"
	go run ./cmd/consumer/main.go &
	go run ./cmd/producer/main.go
	@echo "Both services are running"