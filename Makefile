COMPOSE ?= docker compose

.PHONY: up down logs build test migrate dev-db
up:
	$(COMPOSE) up -d
down:
	$(COMPOSE) down
logs:
	$(COMPOSE) logs -f
build:
	$(COMPOSE) build
test:
	cd backend && go test ./... && go vet ./... && go build -o /tmp/ent-go-vue-api ./cmd/api
	cd frontend/vue-ui && npm run build
migrate:
	$(COMPOSE) run --rm migrate
dev-db:
	$(COMPOSE) -f docker-compose.yml -f docker-compose.dev.yml up -d postgres
