COMPOSE ?= docker compose

.PHONY: up down logs build test lint complexity duplication check migrate dev-db
up:
	$(COMPOSE) up -d
down:
	$(COMPOSE) down
logs:
	$(COMPOSE) logs -f
build:
	$(COMPOSE) build
test:
	cd backend && go test -count=1 ./... && go vet ./cmd/... ./internal/... && go build -o /tmp/ent-go-vue-api ./cmd/api
	cd frontend/vue-ui && npm run build
lint:
	$(MAKE) -C backend lint
	cd frontend/vue-ui && npm run lint && npm run typecheck
complexity:
	$(MAKE) -C backend complexity
	cd frontend/vue-ui && npm run complexity
duplication:
	$(MAKE) -C backend duplication
	cd frontend/vue-ui && npm run duplication
check:
	$(MAKE) lint
	$(MAKE) complexity
	$(MAKE) duplication
	$(MAKE) test
migrate:
	$(COMPOSE) run --rm migrate
dev-db:
	$(COMPOSE) -f docker-compose.yml -f docker-compose.dev.yml up -d postgres
