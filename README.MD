# BiteDash

Backend и frontend для food delivery приложения. Пользователь может зарегистрироваться, войти, просматривать рестораны, добавлять продукты в корзину и оформлять заказ.

## Стек

| Часть | Технологии |
|-------|------------|
| Backend | Go, Gin, PostgreSQL, gRPC, JWT |
| Frontend | React, TypeScript, Vite, React Router |
| Observability | Prometheus, Grafana, Loki, Alloy, Tempo |

## Быстрый старт (Docker)

### 1. Подготовка

```bash
git clone https://github.com/CodeForgeHQ/BiteDash.git
cd BiteDash
```

Скопируйте `.env.example` в `.env` и задайте значения:

```bash
cp .env.example .env
```

Обязательные переменные:

- `POSTGRES_PASSWORD`
- `JWT_SECRET`
- `GRAFANA_ADMIN_PASSWORD`

### 2. Запуск всего стека

```powershell
docker compose up -d --build
```

Поднимаются сервисы:

| Сервис | URL |
|--------|-----|
| Frontend | http://localhost:3001 |
| Backend HTTP API | http://localhost:8080 |
| Backend gRPC | localhost:9090 |
| Swagger | http://localhost:8080/swagger/index.html |
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9091 |
| Loki | http://localhost:3100 |
| Tempo | http://localhost:3200 |
| Alloy | http://localhost:12345 |
| PostgreSQL | localhost:5433 |

Проверка контейнеров:

```powershell
docker ps
```

Проверка backend:

```powershell
curl.exe http://localhost:8080/health
```

### 3. Сценарий проверки

1. Откройте http://localhost:3001
2. Зарегистрируйтесь на `/register`
3. Войдите на `/login`
4. Загрузите рестораны на `/restaurants` (кнопка **Search**)
5. Откройте ресторан и добавьте продукт в корзину
6. Оформите заказ на `/cart`

Пересборка только frontend (после изменения `frontend/.env.production`):

```powershell
docker compose up -d --build frontend
```

## Локальная разработка

### Backend

```bash
go run ./cmd/api
```

API по умолчанию: http://localhost:8080

### Frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Приложение: http://localhost:5173

В `.env` укажите:

```env
VITE_API_BASE_URL=http://localhost:8080
```

Frontend в Docker собирается в production-режиме (Vite build + Nginx). Для локальной разработки используйте `npm run dev`.

## Структура проекта

```text
BiteDash/
├── cmd/                  # точки входа (api, grpc-client)
├── internal/             # бизнес-логика, handlers, middleware
├── frontend/             # React SPA
├── deploy/               # конфиги Prometheus, Grafana, Loki, Alloy, Tempo
├── docker-compose.yml
└── Dockerfile
```

## Логи в Grafana

В Grafana → Explore → Loki:

```logql
{service="api"}
{service="frontend"}
```

## CORS

Backend разрешает запросы с origins:

- `http://localhost:3001` (frontend в Docker)
- `http://localhost:5173` (frontend в dev-режиме)

---

## Работа над задачами

### Создание ветки

```bash
git checkout main
git pull origin main
git checkout -b feature/название-вашей-фичи
```

### Коммит и push

```bash
git add .
git commit -m "feat: добавил форму входа"
git push origin feature/название-вашей-фичи
```

### Pull Request

1. Перейдите на [GitHub BiteDash](https://github.com/CodeForgeHQ/BiteDash)
2. Нажмите **Compare & pull request**
3. Отправьте ссылку на PR в рабочую группу

### Именование коммитов

- `feat:` — новая функциональность
- `fix:` — исправление бага
- `docs:` — документация
- `refactor:` — рефакторинг без изменения поведения
- `test:` — тесты
- `chore:` — обслуживание (зависимости, конфиги)
- `build:` — сборка и CI
- `ci:` — CI/CD

Примеры:

- `fix: исправил баг с отображением списка заказов`
- `feat: добавил возможность добавления новых ресторанов`
- `docs: обновил README`

## Ссылки

- **Архитектура и БД:** [Excalidraw](https://excalidraw.com/#room=4750d506c502e2215d47,1820LWrTrDuJJ0KAhPSlMA)
