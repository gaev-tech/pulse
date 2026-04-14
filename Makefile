include .env
export

.PHONY: dev prod lint test-api test-e2e swagger seed

dev:
	docker compose -f docker-compose.dev.yml up --build

prod:
	docker compose up --build

lint:
	cd backend && golangci-lint run ./...
	cd frontend && ./node_modules/.bin/eslint .

test-api:
ifdef t
	cd backend && go test ./tests/... -run $(t) -v
else
	cd backend && go test ./tests/... -v
endif

test-e2e:
	cd e2e && npx playwright test --ui

swagger:
	cd backend && swag init -g cmd/api/main.go -o api

seed:
	cd backend && go run cmd/seed/main.go
