run:
	go run cmd/main.go

up:
	docker compose up -d

down:
	docker compose down

remove volume:
	docker volume rm slack-clone-backend_db

init:
	go run cmd/migrations/main.go db init

migrate:
	go run cmd/migrations/main.go db migrate

create_sql:
	go run cmd/migrations/main.go db create_sql