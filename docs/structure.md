# Structure

## Monorepo

```
pulse/
├── backend/          # Go API
├── frontend/         # SvelteKit
├── e2e/              # Playwright e2e тесты
├── nginx/            # Nginx конфиг
├── docs/             # Документация проекта
├── docker-compose.yml
├── docker-compose.dev.yml
├── Makefile
├── .env.example
├── README.md
└── CLAUDE.md
```

## Backend

```
backend/
├── cmd/
│   └── api/
│       └── main.go           # Точка входа, старт HTTP сервера
├── internal/
│   ├── domain/               # Доменные сущности и интерфейсы репозиториев
│   │   ├── user/
│   │   ├── team/
│   │   ├── task/
│   │   ├── filter/
│   │   ├── permission/
│   │   └── event/
│   ├── usecase/              # Бизнес-логика
│   │   ├── user/
│   │   ├── team/
│   │   ├── task/
│   │   ├── filter/
│   │   └── permission/
│   ├── handler/
│   │   ├── v1/               # HTTP handlers + router.go
│   │   └── middleware/       # Auth JWT, permission check
│   ├── repository/
│   │   ├── postgres/         # Реализации интерфейсов репозиториев
│   │   └── migrations/       # SQL-миграции (golang-migrate)
│   └── infrastructure/
│       ├── config/           # Конфигурация из переменных окружения
│       ├── postgres/         # Подключение к БД
│       └── logger/           # Zap logger
├── pkg/
│   ├── apierror/             # Стандартизированные ошибки API
│   ├── validator/            # Валидация запросов
│   └── pagination/           # Cursor-based пагинация
├── tests/                    # API-тесты (go test)
├── api/
│   └── openapi.yaml          # OpenAPI спецификация
├── go.mod
└── Dockerfile
```

Зависимости направлены внутрь: `handler → usecase → domain ← repository`.
`domain/` не импортирует БД и HTTP. `pkg/` содержит утилиты без бизнес-логики,
доступные из любого слоя.

## Frontend

```
frontend/
├── src/
│   ├── api/              # API клиент
│   │   ├── client.ts     # Базовый fetch + JWT заголовок
│   │   ├── tasks.ts
│   │   ├── teams.ts
│   │   ├── filters.ts
│   │   └── ...
│   ├── components/
│   │   ├── feed/
│   │   ├── tasks/
│   │   ├── teams/
│   │   └── ui/           # Переиспользуемые UI-компоненты
│   ├── stores/
│   │   ├── session.ts    # Сессии, персистируется в localStorage
│   │   ├── feed.ts       # Cursor для пагинации ленты
│   │   └── ...
│   └── routes/           # Файловый роутинг SvelteKit
│       ├── +layout.svelte
│       ├── feed/
│       ├── tasks/[key]/
│       ├── teams/[slug]/
│       └── filters/[id]/
├── svelte.config.js
├── vite.config.ts
└── Dockerfile
```

## Nginx

```
nginx/
└── conf.d/
    └── default.conf          # Роутинг: /api/ → backend, / → frontend
```

## Docker Compose

- `docker-compose.yml` — production стек
- `docker-compose.dev.yml` — dev стек с hot reload (air для Go, vite dev server для SvelteKit)
