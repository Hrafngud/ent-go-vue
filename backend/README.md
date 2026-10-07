# Backend

Go + `net/http` + Huma + Ent + Atlas + PostgreSQL, retaining the feature-based Clean Architecture under `internal/auth` and `internal/user`.

The integrated Docker environment now lives at the repository root. See the [root README](../README.md) for configuration, migration, deployment, native development, and validation commands, and the [dependency audit](../DEPENDENCIES.md) for version decisions.

```bash
# From the repository root, configure .env first:
make dev-db
make migrate
cp backend/.env.example backend/.env
# Edit backend/.env: use matching credentials, localhost:55432, and a random JWT secret.
cd backend
go run ./cmd/api
```

The native API listens on port 8080 by default. Routes and Huma documentation are under `/api`; `/api/health` is liveness, `/api/ready` verifies PostgreSQL connectivity.

Backend Make targets:

- `make run`: run the API.
- `make test`: fast tests (`-short`, no containers).
- `make test-integration`: full tests with Testcontainers PostgreSQL.
- `make lint`: normal static analysis with golangci-lint 2.14.0.
- `make complexity`: cyclop, maximum 15.
- `make duplication`: dupl, minimum 100 tokens.
- `make tidy`: tidy module dependencies.
- `make generate`: regenerate the Ent client after schema changes.
- `make atlas-install`: install pinned Atlas locally under `bin/` for the host OS/architecture.
- `make migrate name=...`: generate versioned SQL against PostgreSQL 17 (Docker required).
- `make apply`: apply committed SQL with native Atlas using exported `DATABASE_URL`.

The API checks the migrated schema on startup, supports SIGTERM shutdown, and never changes the database schema. Atlas application runs as a separate Compose job or CLI command. The root and backend migration Make targets have different roles: root **applies**, backend **generates**.

Install the pinned linter into `bin/golangci-lint` as documented in the
[root README](../README.md#test-and-operate), or override `GOLANGCI_LINT` with its
installed path. Lint targets select `cmd` and `internal`, excluding Ent. Existing
debt exceptions are recorded in [VALIDATION.md](../VALIDATION.md#local-quality-gate).
Every implementation must finish with a passing `make check` at the repository root,
which also runs the frontend checks and full integration tests/builds.
