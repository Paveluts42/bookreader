# BookReader

Веб-сервис для чтения и аннотирования PDF-книг.

## Возможности

- Загрузка и чтение PDF-книг
- Добавление и просмотр заметок к книгам
- Закладки и аудиокниги
- Авторизация пользователей (JWT)
- Админ-панель для управления пользователями и книгами

## Технологии

- **Frontend:** React, TypeScript, Vite, PDF.js, Nginx
- **Backend:** Go, gRPC-Web, PostgreSQL, Docker

## Быстрый старт (локально)

1. Клонируйте репозиторий:
    ```bash
    git clone <repo-url>
    cd bookreader
    ```

2. Запустите сервисы:
    ```bash
    docker compose up -d --build
    ```

3. Откройте фронтенд:  
    [http://localhost/](http://localhost/)

## Доступ на сервере

- Фронтенд: [http://212.113.119.120/](http://212.113.119.120/)
- Бэкенд: [http://212.113.119.120/api/](http://212.113.119.120/api/)

## Структура проекта

- `/frontend` — клиентская часть (React + Vite)
- `/backend` — серверная часть (Go + gRPC)
- `/docker-compose.yaml` — запуск всех сервисов

## Переменные окружения

- `.env` — UID/GID для контейнеров

## Основные команды

- `docker compose up -d --build` — сборка и запуск всех сервисов
- `docker compose logs -f frontend` — логи фронтенда
- `docker compose logs -f backend` — логи бэкенда



---

> _Для тестирования используйте [http://212.113.119.120/](http://212.113.119.120/)_