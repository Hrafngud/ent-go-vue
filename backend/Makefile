ATLAS ?= ./bin/atlas
ATLAS_VERSION ?= v1.3.0
GOLANGCI_LINT ?= ./bin/golangci-lint
# Ent is compiled as a dependency, but never passed as a linter target.
LINT_PACKAGES := ./cmd/... ./internal/...

.PHONY: run generate tidy migrate atlas-install apply test test-integration lint complexity duplication

run:
	go run ./cmd/api

generate:
	go generate ./ent

tidy:
	go mod tidy

test:
	go test -short ./...

test-integration:
	go test ./... -v

lint:
	$(GOLANGCI_LINT) run --config .golangci.yml --disable cyclop,dupl $(LINT_PACKAGES)

complexity:
	$(GOLANGCI_LINT) run --config .golangci.yml --enable-only cyclop $(LINT_PACKAGES)

duplication:
	$(GOLANGCI_LINT) run --config .golangci.yml --enable-only dupl $(LINT_PACKAGES)

migrate:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate name=<migration_name>"; \
		exit 1; \
	fi
	mkdir -p ent/migrate/migrations
	$(ATLAS) migrate diff $(name) \
		--dir "file://ent/migrate/migrations" \
		--to "ent://ent/schema" \
		--dev-url "docker://postgres/17/dev?search_path=public"

atlas-install:
	mkdir -p bin
	@atlas_os=$$(go env GOOS); atlas_arch=$$(go env GOARCH); \
		curl -fsSL "https://release.ariga.io/atlas/atlas-$$atlas_os-$$atlas_arch-$(ATLAS_VERSION)" -o bin/atlas
	chmod +x bin/atlas

apply:
	$(ATLAS) migrate apply --env app
