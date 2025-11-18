include .env
export

# Змінні для міграцій
MIGRATE_PATH=./migrations
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable
TEST_DB_URL=postgres://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@localhost:$(TEST_DB_PORT)/$(TEST_DB_NAME)?sslmode=disable

.PHONY: up down migrate-create migrate-up migrate-down migrate-test-up migrate-test-down build

# Запуск усієї інфраструктури (БД, Python-сервіс)
up:
	docker-compose up -d

# Зупинка контейнерів
down:
	docker-compose down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Запуск міграцій "вгору"
migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

# Відкат міграцій "вниз"
migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down

# Міграції для тестової БД
migrate-test-up:
	migrate -path $(MIGRATE_PATH) -database "$(TEST_DB_URL)" up

migrate-test-down:
	migrate -path $(MIGRATE_PATH) -database "$(TEST_DB_URL)" down

# Збірка та запуск Go-бекенду
run:
	go run cmd/api/main.go