# Ent Go Vue

A full-stack boilerplate for building applications with Go and Vue. Includes registration, JWT authentication, an optional administrator account, user management, database migrations, and a Docker Compose stack.

- **Backend:** Go, `net/http`, Huma, Ent, Atlas, and PostgreSQL.
- **Frontend:** Vue 3, TypeScript, Vite, Tailwind CSS, daisyUI, and Phosphor Icons. TanStack Vue Query manages server state; Pinia manages client state.

## Architecture

```text
Browser → Nginx
          ├── /      → Vue frontend
          └── /api/* → Go API → PostgreSQL
```

The frontend and API share one origin through a reverse proxy. The backend groups features under `internal/auth` and `internal/user`, with domain, use case, repository, and HTTP layers. Ent defines the schema; Atlas applies committed SQL migrations before the API starts.

| Path | Contents |
| --- | --- |
| `backend/cmd/api/` | API entry point |
| `backend/internal/` | Features, configuration, and HTTP routing |
| `backend/ent/` | Ent schemas, generated client, and SQL migrations |
| `frontend/vue-ui/src/` | Pages, components, routing, API client, stores, and queries |
| `docker/nginx/` | Reverse proxy |
| `docker-compose*.yml` | Application stack and native development database override |

## Run with Docker Compose

Requires Docker with the Compose plugin. Go, Node, and Atlas run inside containers.

```bash
cp .env.example .env
```

Edit `.env` before starting:

- Replace `POSTGRES_PASSWORD` and use the same credentials in `DATABASE_URL`. Keep the URL hostname as `postgres`; percent-encode special characters in credentials.
- Replace `JWT_SECRET` with at least 32 random characters (`openssl rand -hex 32`). The API rejects the example placeholder.
- Optionally set both `ROOT_EMAIL` and `ROOT_PASSWORD` to provision an administrator. Use your own 8–72 byte password. Startup reapplies this password.
- Set `HTTP_PORT` if port 80 is unavailable; for example, `HTTP_PORT=18080`.

```bash
docker compose config --quiet
docker compose up --build -d
docker compose ps --all
```

Open <http://localhost/> (or your configured port). Register an account or sign in with the configured administrator. API documentation is at `/api/docs`; OpenAPI is at `/api/openapi.json`. Use `/api/health` for liveness and `/api/ready` for database readiness.

Compose starts PostgreSQL, applies migrations, then starts the API, frontend, and proxy. The migration container exits successfully when done. Only Nginx publishes an application port; database data persists in `postgres_data`.

```bash
docker compose logs -f
docker compose down          # stop and preserve database data
```

Use `docker compose up --build -d` after source changes. `docker compose down -v` deletes database data. Database initialization credentials only take effect when the volume is first created.

## Make commands

Run these from the repository root:

| Command | Purpose |
| --- | --- |
| `make up` | Start the Compose stack |
| `make down` | Stop the stack and preserve database data |
| `make logs` | Follow service logs |
| `make build` | Build production Docker images |
| `make dev-db` | Start PostgreSQL on `localhost:55432` for native development |
| `make migrate` | Apply pending SQL migrations with the Compose migration job |
| `make test` | Run Go tests, vet, API build, and frontend production build |
| `make lint` | Run Go/frontend lint and frontend type checks |
| `make complexity` | Check function complexity |
| `make duplication` | Check duplicated code |
| `make check` | Run all quality checks, tests, and native builds |

## Native development

Requires Go 1.26.4+ (the module selects toolchain 1.27.1), Node 22.13+ or 24+, npm, and Docker.

Configure the root `.env` as above, then:

```bash
make dev-db
make migrate
cp backend/.env.example backend/.env
```

Edit `backend/.env` with matching database credentials and a random JWT secret. Use `localhost:55432` in its `DATABASE_URL` (or the root `.env`'s `POSTGRES_PORT`).

Run the API:

```bash
make -C backend run
```

In another terminal, run the frontend:

```bash
cd frontend/vue-ui
npm ci
npm run dev
```

Open the URL printed by Vite. Vite proxies `/api` to `http://localhost:8080`. If you change the API's `PORT`, copy `frontend/vue-ui/.env.example` to `.env` and set `API_PROXY_TARGET` to match.

Before switching back to the full stack, stop the development containers with `docker compose -f docker-compose.yml -f docker-compose.dev.yml down`.

## Schema changes

After editing `backend/ent/schema/`, regenerate the Ent client and create a migration:

```bash
make -C backend atlas-install
make -C backend generate
make -C backend migrate name=describe_schema_change
```

Review and commit the generated SQL and `atlas.sum`. The backend target **generates** migrations; root `make migrate` **applies** them. Rebuild the migration image before applying new SQL:

```bash
docker compose build migrate
make migrate
```

The API does not modify the schema at startup. Avoid editing migrations already applied to a database.

## Test and operate

Run `make check` before completing changes. It requires the native development tools, GNU Make, frontend dependencies, golangci-lint 2.14.0, and a running Docker daemon for PostgreSQL integration tests. The application stack and `.env` files are unnecessary for these checks.

```bash
curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/v2.14.0/install.sh -o /tmp/ent-go-vue-install-golangci-lint.sh
sh /tmp/ent-go-vue-install-golangci-lint.sh -b backend/bin v2.14.0
npm --prefix frontend/vue-ui ci
make check
```

For quick backend tests without Docker, use `make -C backend test`.

See [backend commands](backend/README.md), [frontend commands](frontend/vue-ui/README.md), and [dependency notes](DEPENDENCIES.md) for more detail. Configuration defaults, including container resource limits, are listed in [.env.example](.env.example).
