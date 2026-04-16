# API (Go Backend)

Backend-сервис на Go с авторизацией, CRUD API и кэшированием через Redis.  
Проект сделан как pet-project с упором на архитектуру (handlers → service → repository) и практики, используемые в реальной разработке.

## Stack

- Go (Gin)
- PostgreSQL
- Redis
- JWT (auth)
- Docker / Docker Compose

## Features

- регистрация и логин пользователей
- JWT авторизация (middleware)
- CRUD API
- фильтрация + пагинация
- Redis кэш для GET запросов
- инвалидция кэша при изменениях
- миграции базы данных

## Project structure

cmd/                # entry point (main)
internal/
  user/             # auth + user logic
  APIs/             # основной API (handlers/service/repo)
  jwt/              # middleware авторизации
  errs/             # кастомные ошибки
migrations/         # SQL миграции
Dockerfile
docker-compose.yml

## Run with Docker

1. clone repo

git clone https://github.com/ArturM-r/API.git
cd API

2. create .env

POSTGRES_USER=kobazaz
POSTGRES_PASSWORD=123456
POSTGRES_DB=apis

DATABASE_URL=postgres://example:123456@postgres:5432/apis?sslmode=disable
REDIS_URL=redis://redis:6379
HMAC_KEY=supersecretkey

3. run

docker compose up --build

API будет доступен на:
http://localhost:8080

## API

Auth:

POST /user/register  
POST /user/login  

Protected routes (JWT required):

GET    /api  
GET    /api/:id  
POST   /api  
PATCH  /api/:id  
DELETE /api/:id  

Authorization header:

Authorization: Bearer <token>

## Example requests

Create:

POST /api

{
  "title": "task"
}

Update:

PATCH /api/:id

{
  "title": "new title",
  "completed": true
}

## Architecture

Проект построен по классической схеме:

- handler — работа с HTTP (Gin)
- service — бизнес-логика
- repository — работа с БД

Это позволяет:
- легко тестировать код
- разделять ответственность
- менять инфраструктуру без переписывания логики

## Notes

- PostgreSQL используется как основная БД
- Redis используется как кэш (GET /api)
- кэш инвалидируется при create/update/delete
- используется offset/limit пагинация
