.PHONY: swagger build run docker-up docker-down

swagger:
	swag init -g cmd/server/main.go -o docs

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

install-deps:
	go install github.com/swaggo/swag/cmd/swag@latest