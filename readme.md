
# EM Subscription Manager

REST-сервис для управления и агрегации данных об онлайн-подписках пользователей.

## Основные возможности
* **CRUDL подписок**: Создание, чтение, обновление, удаление и список всех записей.
* **Агрегация**: Расчет суммарной стоимости подписок пользователя за выбранный период с фильтрацией по названию сервиса.
* **Автоматические миграции**: Использование `goose` для инициализации схемы БД.
* **Swagger UI**: Интерактивная документация API.
* **Логирование**: Структурированные логи с помощью `zerolog`.

## Технологический стек
* **Language**: Go 1.26
* **Framework**: Fiber v3
* **Database**: PostgreSQL (Driver: pgx)
* **Migrations**: Goose
* **Containerization**: Docker / Docker Compose
* **Documentation**: Swaggo (Swagger 2.0)

## Быстрый запуск

Для запуска всего окружения (приложение + база данных) выполните команду:

```bash
docker-compose up --build
```

Сервис будет доступен по адресу: `http://localhost:3000`

## API Документация
После запуска проекта Swagger UI доступен по адресу:
`http://localhost:3000/swagger`

## Формат данных
* **Цена**: Целое число (рубли).
* **ID пользователя**: Формат UUID.
* **Даты**: Строка в формате `MM-YYYY` (например, `07-2025`).

### Пример запроса на создание подписки (POST /subscriptions)
```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "07-2026"
}
```

### Пример запроса суммы (GET /subscriptions/summary)
```text
GET /subscriptions/summary?user_id=...&start_date=01-2024&end_date=12-2024&service_name=Netflix
```

## Конфигурация
Настройки приложения находятся в `configs/config.yaml`.  
Основные параметры:
* `app.port`: Порт сервера (по умолчанию 3000).
* `postgres.dsn`: Строка подключения к БД.
* `logger.level`: Уровень логирования (debug, info, error).

## Структура проекта
* `cmd/main.go` — точка входа.
* `internal/models` — доменные модели.
* `internal/services` — бизнес-логика.
* `internal/repo` — уровень работы с БД.
* `internal/handlers` — HTTP-контроллеры.
* `migrations/` — SQL файлы миграций.

---

### Разработка (локальный запуск)
Если вы хотите запустить сервис без Docker:
1. Установите зависимости: `go mod download`
2. Настройте БД и укажите DSN в `configs/config.yaml`.
3. Запустите: `CONFIG_PATH=configs/config.yaml go run cmd/main.go`