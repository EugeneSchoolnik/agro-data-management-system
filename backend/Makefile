include .env
export

# Змінні для міграцій
MIGRATE_PATH=./migrations
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable

.PHONY: up down migrate-up migrate-down build

# Запуск усієї інфраструктури (БД, Python-сервіс)
up:
	docker-compose up -d

# Зупинка контейнерів
down:
	docker-compose down

# Запуск міграцій "вгору"
migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

# Відкат міграцій "вниз"
migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down

# Збірка та запуск Go-бекенду
run:
	go run cmd/api/main.go