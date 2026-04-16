# API (Go Backend)

Backend-сервис на Go с авторизацией, CRUD API, Redis-кэшем и PostgreSQL.

## Стек

- Go (Gin)
- PostgreSQL
- Redis
- JWT авторизация
- Docker / Docker Compose

## Функционал

- Регистрация и логин пользователей
- JWT авторизация
- CRUD для API (todo-подобные сущности)
- Кэширование через Redis
- Миграции базы данных

## Структура проекта

cmd/                # точка входа (main)
internal/
  user/             # логика пользователей (auth, service, repo)
  APIs/             # основное API (handlers, service, repo)
  jwt/              # middleware авторизации
  errs/             # кастомные ошибки
migrations/         # SQL миграции
Dockerfile
docker-compose.yml

## Быстрый старт

1. Клонировать репозиторий

git clone https://github.com/ArturM-r/API.git
cd API

2. Создать .env

POSTGRES_USER=kobazaz
POSTGRES_PASSWORD=123456
POSTGRES_DB=hardtodo

DATABASE_URL=postgres://kobazaz:123456@postgres:5432/hardtodo?sslmode=disable
REDIS_URL=redis://redis:6379
HMAC_KEY=supersecretkey

3. Запуск через Docker

docker compose up --build

Сервис будет доступен на:
http://localhost:8080

## API

Регистрация  
POST /user/register

{
  "email": "test@example.com",
  "password": "123456"
}

Логин  
POST /user/login

{
  "user": {...},
  "token": "JWT_TOKEN"
}

Авторизация  
Authorization: Bearer <token>

Получить все записи  
GET /api

Query:
- limit
- offset
- completed
- title

Получить по ID  
GET /api/:id

Создать  
POST /api

{
  "title": "task"
}

Обновить  
PATCH /api/:id

{
  "title": "new title",
  "completed": true
}

Удалить  
DELETE /api/:id
