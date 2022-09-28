run:
	go run main.go

format:
	go fmt ./...

up:
	docker compose up -d

down:
	docker compose down

remove volume:
	docker volume rm slack-clone-backend_db