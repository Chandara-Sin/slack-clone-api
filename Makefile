run:
	go run main.go

format:
	go fmt ./...

up:
	docker compose up -d

down:
	docker compose down