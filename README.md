# Authentication Service

Микросервис для аутентификации пользователей с поддержкой JWT, refresh-токенов, хранения сессий в Postgres и ревокации access-токенов через Redis. Есть Swagger-документация, миграции, тесты, docker-compose для локального запуска.

## Стек
- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker, docker-compose
- Swagger (swaggo)
- sqlx, squirrel, bcrypt
- Тесты: testify, sqlmock

## Переменные окружения (пример .env)
Для запуска приложения нужно создать файл .env в котором будет указано, например, следующее

```
POSTGRES_DB=auth_db
POSTGRES_USER=postgres
POSTGRES_PASSWORD=13
DB_HOST=db
DB_PORT=5432
DB_USER=authuser
DB_PASSWORD=authpass
DB_NAME=authdb
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0
TTL_ACCESS_TOKEN=3600 # в секундах
JWT_SECRET_KEY=supersecretkey
WEBHOOK_URL=http://example.com/webhook
```

## Быстрый старт

```bash
git clone <REPO_URL>
cd authentication-service
docker-compose up --build
```


## Архитектура
- **Postgres**: хранит сессии (user_id, refresh_token_hash, user_agent, ip_addr)
- **Redis**: хранит revoked access-токены (blacklist)
- **Swagger**: автогенерируется из Go-комментариев
- **Миграции**: в internal/migrations, применяются через migrate/migrate

## Основные эндпоинты

- `GET /api/v1/auth/tokens?user_id=...` — получить пару access/refresh токенов
- `POST /api/v1/auth/tokens/refresh` — обновить пару токенов (тело: {access_token, refresh_token})
- `GET /api/v1/auth/guid` — получить user_id из access_token (требует Authorization)
- `POST /api/v1/auth/logout` — разлогинить пользователя (требует Authorization)

**Полное описание и схемы ошибок — в Swagger!**

## Токены
- **Access**: JWT (HS512), не хранится в БД, revocation через Redis
- **Refresh**: случайная строка, хранится в БД только bcrypt-хеш


## Миграции
Миграции лежат в internal/migrations. Применяются автоматически через сервис `migrate` в docker-compose.

## Тесты

```bash
go test ./...
```
Покрытие: сервисная логика, repo, handlers (моки через testify/sqlmock).

## Swagger

- Автогенерируется через swaggo (см. internal/handlers/*)
- После запуска: http://localhost:8080/swagger/index.html


## Интеграционные тесты

Интеграционные тесты находятся в директории `integration/` и покрывают полный сценарий работы сервиса:
- Получение access/refresh токенов
- Обновление токенов
- Получение GUID пользователя по access token
- Logout (ревокация токенов)
- Проверка невозможности refresh/logout/access после logout (black-list)
- Повторный логин и работа с новыми токенами
- Проверка edge cases (повторный logout, refresh с отозванными токенами и т.д.)

### Как запускать

1. Убедись, что весь стек поднят через docker-compose:
   ```bash
   docker-compose up --build -d
   ```
2. Запусти интеграционные тесты:
   ```bash
   go test -v ./integration/
   ```

Тесты автоматически ждут поднятия API и используют реальные endpoints, взаимодействуя с Postgres и Redis (revocation store).

### Что проверяют
- Корректность выдачи и обновления токенов
- Сохранение и удаление сессий в базе
- Работу black-list (revocation store) для access/refresh токенов
- Защиту эндпоинтов через access token
- Обработку edge cases и ошибок (дубликаты, невалидные токены, повторные действия)

Тесты используют уникальные UUID для user_id, чтобы не было конфликтов между запусками.

