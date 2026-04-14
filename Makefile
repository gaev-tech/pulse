include .env
export

.PHONY: dev prod prod-down prod-down-v lint test-api test-e2e test-e2e-ui swagger seed

## Запустить dev-окружение (останавливает prod, если запущен)
dev:
	@if docker ps --format '{{.Ports}}' 2>/dev/null | grep -q '0\.0\.0\.0:5173->' ; then \
		: ; \
	elif docker ps --filter 'name=pulse-nginx-1' --format '{{.Names}}' 2>/dev/null | grep -q '.' ; then \
		echo "Останавливаю prod..."; \
		docker compose down; \
	fi
	docker compose -f docker-compose.dev.yml up --build

## Запустить prod-окружение в фоне (выдаёт ошибку, если запущен dev)
prod:
	@if docker ps --format '{{.Ports}}' 2>/dev/null | grep -q '0\.0\.0\.0:5173->' ; then \
		echo "Ошибка: запущен dev. Остановите его (Ctrl+C) и повторите."; \
		exit 1; \
	fi
	docker compose up --build -d

## Остановить prod-окружение
prod-down:
	docker compose down

## Остановить prod-окружение и удалить тома
prod-down-v:
	docker compose down -v

## Запустить линтеры (backend + frontend)
lint:
	cd backend && golangci-lint run ./...
	cd frontend && ./node_modules/.bin/eslint .

## Запустить API-тесты (make test-api t=TestName для конкретного теста)
test-api:
ifdef t
	cd backend && go test ./tests/... -run $(t) -v
else
	cd backend && go test ./tests/... -v
endif

## Запустить e2e-тесты локально
test-e2e:
	cd e2e && npx --registry https://registry.npmjs.org/ playwright test

## Запустить e2e-тесты с UI Playwright
test-e2e-ui:
	cd e2e && npx --registry https://registry.npmjs.org/ playwright test --ui

## Сгенерировать Swagger-документацию
swagger:
	cd backend && swag init -g cmd/api/main.go -o api

## Заполнить базу тестовыми данными
seed:
	cd backend && go run cmd/seed/main.go
