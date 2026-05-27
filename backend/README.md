# Agro Data Management System - Backend

Система управління даними для сільськогосподарських полів з підтримкою датчиків, прогнозування та аналізу шкідників.

## 🚀 Вибір стартер

### 1. Запустити API сервер

```bash
go run cmd/api/main.go
```

### 2. Створити нового користувача

```bash
go run cmd/createuser/main.go --email user@example.com --password pass123 --role user
```

### 3. Заповнити БД тестовими даними

```bash
go run cmd/seeddata/main.go -config config/local.yaml
```

### 4. Перевірити створені дані

```bash
go run cmd/verifydata/main.go
```

## 📋 Вимоги

- Go 1.19+
- PostgreSQL 15+
- Docker & Docker Compose (для запуску БД)

## 🛠️ Встановлення

### 1. Клонувати репозиторій та встановити залежності

```bash
cd backend
go mod download
```

### 2. Запустити Docker контейнер з БД

```bash
docker-compose up -d
```

### 3. Застосувати міграції

```bash
# Використовувати migrate CLI
migrate -path migrations -database "postgres://postgres:9210@localhost:5432/agro-db?sslmode=disable" up
```

## ⚙️ Конфігурація

### .env файл

```env
# Database
DB_USER=postgres
DB_PASSWORD=9210
DB_NAME=agro-db

# Server
PORT=8080

# JWT
JWT_SECRET=qwerty

# Weather API
WEATHER_API_LOGIN=your_email@example.com
WEATHER_API_PASSWORD=your_api_key
```

## 📊 Структура проекту

```
backend/
├── cmd/
│   ├── api/              # REST API сервер
│   ├── createuser/       # Утиліта для створення користувачів
│   ├── seeddata/         # Утиліта для заповнення тестовими даними
│   └── verifydata/       # Утиліта для перевірки даних
├── internal/
│   ├── config/           # Конфігурація
│   ├── handler/          # HTTP обробники
│   ├── models/           # Структури даних
│   ├── repository/       # Доступ до БД
│   ├── service/          # Бізнес логіка
│   └── weather/          # Інтеграція з погодою
├── migrations/           # SQL міграції
├── docs/                 # Документація
└── docker-compose.yml    # Docker конфігурація
```

## 🗄️ Структура БД

### Основні таблиці

- **crops** - типи культур
- **fields** - сільськогосподарські поля
- **sensors** - датчики на полях
- **metrics** - дані з датчиків
- **pests** - типи шкідників
- **forecasts** - прогнози ризиків
- **weather_parameters** - параметри погоди
- **weather_stations** - метеостанції
- **users** - користувачі системи

Детальну документацію див. в [docs/TESTDATA.md](docs/TESTDATA.md)

## 🔧 Розробка

### Запустити тести

```bash
go test ./...
```

### Запустити тести з покриттям

```bash
go test -cover ./...
```

### Запустити API сервер в режимі розробки

```bash
go run cmd/api/main.go
```

## 📈 Робочий процес - Практичний приклад

1. **Запустити БД**

   ```bash
   docker-compose up -d
   ```

2. **Створити тестового користувача**

   ```bash
   go run cmd/createuser/main.go --email test@agro.com --password test123 --role admin
   ```

3. **Заповнити БД тестовими даними**

   ```bash
   go run cmd/seeddata/main.go -config config/local.yaml
   ```

4. **Перевірити дані**

   ```bash
   go run cmd/verifydata/main.go
   ```

5. **Запустити API сервер**

   ```bash
   go run cmd/api/main.go
   ```

6. **Тестувати API**
   ```bash
   curl http://localhost:8080/api/health
   ```

## 📚 API Документація

Основні endpoints:

- `GET /api/fields` - отримати всі поля
- `GET /api/fields/:id` - отримати конкретне поле
- `GET /api/sensors/field/:fieldId` - отримати датчики поля
- `GET /api/metrics/sensor/:sensorId` - отримати метрики датчика
- `GET /api/pests` - отримати список шкідників
- `GET /api/forecasts/field/:fieldId` - отримати прогнози для поля

## 🐛 Вирішення проблем

### Помилка підключення до БД

```
failed to connect to db: error
```

Перевірте:

- Docker контейнер запущено: `docker ps`
- Правильні мережеві параметри в `.env`

### Міграції не застосовані

```bash
migrate -path migrations -database "postgres://postgres:9210@localhost:5432/agro-db?sslmode=disable" up
```

### Тестові дані не створились

```bash
go run cmd/seeddata/main.go -config config/local.yaml
```

## 📄 Ліцензія

MIT
