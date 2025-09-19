DB_URL=$(DATABASE_URL)

.PHONY: migrate-up migrate-down migrate-new

run:
	go run ./cmd/

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1
