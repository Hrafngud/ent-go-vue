# Ent Go Vue

Vue 3 + TypeScript + Vite + Tailwind CSS 4 + daisyUI 5, backed by Go's `net/http`, Huma, Ent, Atlas, and PostgreSQL. TanStack Vue Query owns server state; Pinia is available for client state.

## Architecture

```text
Browser → Nginx :8080 (published on HTTP_PORT, default 80)
             ├── /           Frontend :8080 → Vue production assets and SPA fallback
             └── /api/*      Go backend :BACKEND_PORT → PostgreSQL :5432
```

The application follows the separate reverse-proxy, frontend, and backend services in the local `Galidor` reference. The frontend image builds Vue and serves its static assets with an unprivileged Nginx runtime. A separate Nginx reverse proxy sends `/` to the frontend and `/api/*` to the backend, so the SPA and API share one origin and port. Only the reverse proxy publishes a port in the base configuration. PostgreSQL persists in `postgres_data`. The backend keeps its existing feature-based Clean Architecture, Ent schemas, and Atlas SQL history.

```text
backend/
  cmd/api/                    configuration and application lifecycle
  internal/config/            environment configuration
  internal/httpapi/           Huma wiring, liveness and readiness
  internal/{auth,user}/       existing domain, usecase, repository, HTTP layers
  ent/                        schemas and generated client
  ent/migrate/migrations/     versioned SQL and atlas.sum
  atlas.hcl                   Atlas environment using DATABASE_URL
  Dockerfile                  compiled API runtime and separate migration target
frontend/vue-ui/
  src/api/                    relative /api HTTP client
  src/queries/                TanStack Query health/readiness requests
  src/components/             visible connectivity display
  Dockerfile                  Node build → static assets in Nginx
  nginx.conf                  frontend static serving and SPA fallback
  vite.config.ts              local /api proxy
docker/nginx/Dockerfile       separate reverse-proxy image
docker/nginx/nginx.conf       frontend and API proxy configuration template
docker-compose.yml           integrated application
docker-compose.dev.yml        optional loopback database port for native development
.env.example                 Compose configuration template
Makefile                     thin command wrappers
```

## Start the integrated stack

Requires Docker with the Compose plugin and BuildKit. Native Node, Go, and Atlas installations are unnecessary for this workflow.

```bash
cp .env.example .env
```

Edit `.env`: replace `POSTGRES_PASSWORD`, put the same password in `DATABASE_URL`, and replace `JWT_SECRET` with at least 32 random characters. Generate values with `openssl rand -hex 32`. Percent-encode reserved characters in URL credentials. The URL's hostname must be `postgres` for containers. Example placeholders are not real credentials; the API rejects the example JWT secret.

```bash
docker compose config --quiet
docker compose build
docker compose up -d
docker compose ps --all
```

`docker compose up -d` also builds missing images on a fresh checkout. After source or dependency changes, rebuild with `docker compose up --build -d`.

Open <http://localhost/>. The page should display **Backend: connected** and **Database: connected**. If port 80 is occupied, set `HTTP_PORT=18080` in `.env` and open <http://localhost:18080/> instead.

```bash
curl --fail http://localhost/api/health
curl --fail http://localhost/api/ready
curl --fail http://localhost/api/users
```

`/api/health` reports `status: ok` without database work. `/api/ready` pings PostgreSQL with a two-second timeout, reports `database: connected`, and returns 503 when the database is unavailable. `/api/users` also exercises the Ent repository against the migrated schema. Responses may include Huma's `$schema` metadata.

Huma documentation and OpenAPI are available at `/api/docs` and `/api/openapi.json`. Feature routes live under `/api/auth/*` and `/api/users*` in both native and container execution.

```bash
docker compose logs nginx frontend backend postgres migrate
docker compose down
```

Stopping preserves PostgreSQL data. `docker compose down -v` intentionally deletes the named database volume. PostgreSQL initialization variables apply only when creating a new data directory; changing `.env` does not rotate credentials in an existing database.

## Migrations

The `migrate` service is an explicit one-shot Atlas deployment job. On `compose up`, PostgreSQL must be healthy, then Atlas applies pending **committed SQL migrations**, validates `atlas.sum`, and exits successfully before the backend starts. A migration error blocks backend startup. The backend never generates migrations or calls `ent.Schema.Create()` at runtime, and `docker compose restart backend` does not migrate.

For a controlled deployment, review migrations before running:

```bash
docker compose up -d postgres
docker compose run --rm migrate
docker compose up -d backend frontend nginx
```

`make migrate` is the root wrapper for the same application step. Run it after rebuilding the migration image when SQL files change. It is safe to repeat when nothing is pending; Atlas records applied versions in its revision table. The exited `migrate` container is expected, not a failed application service.

Generate a new migration on the host, preserving the existing Ent workflow:

```bash
cd backend
make atlas-install
make generate
make migrate name=describe_schema_change
```

The backend Makefile pins Atlas 1.3.0 and selects the host OS/architecture. `make migrate` there generates SQL using `ent://ent/schema` and a disposable PostgreSQL 17 database; it does not apply SQL to the application database. Review generated SQL and its checksum, then commit them. Do not rewrite applied migrations. Rebuild and apply the root migration image before deploying new backend code.

For a database that predates Atlas tracking, inspect its schema and establish the appropriate baseline using Atlas's [existing database workflow](https://atlasgo.io/versioned/import); do not bypass dirty-database checks indiscriminately.

## Native development

Requires Go 1.26.4 or newer (the module selects toolchain 1.27.1), Node 22.12+ or 24+, npm, and Docker for the optional database and integration tests.

Start only PostgreSQL with an explicit development port:

```bash
make dev-db
make migrate
cp backend/.env.example backend/.env
```

Edit `backend/.env` to use the same database credentials and JWT secret as the root file, with hostname **localhost** and port **55432**. If you change `POSTGRES_PORT` in the root file, use that port in the native URL too. The backend loads only its working directory's `.env`; it does not accidentally inherit the root container URL.

```bash
cd backend
go run ./cmd/api
```

Then, in another terminal:

```bash
cd frontend/vue-ui
npm install
npm run dev
```

Vite proxies `/api` to `http://localhost:8080`; browser requests stay relative and need no CORS configuration. To change the native API port, set `PORT` in `backend/.env`, copy the frontend `.env.example` to `.env`, and adjust `API_PROXY_TARGET` to match. That setting is used by Vite only and is never embedded as a production hostname.

To apply migrations using the native Atlas CLI instead of the Compose job:

```bash
cd backend
make atlas-install
# Export trusted local .env values for the CLI (the Go API loads them itself).
set -a
. ./.env
set +a
make apply
```

When switching from native development back to the integrated stack, remove the development containers first and recreate from the base file, preserving the volume:

```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml down
docker compose up -d
```

## Configuration

| Variable | Purpose |
| --- | --- |
| `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` | Required database initialization values |
| `DATABASE_URL` | Required Compose URL; same credentials, hostname `postgres`, port `5432` |
| `BACKEND_PORT` | Internal backend listening/proxy port, default `8080` |
| `HTTP_PORT` | Only publicly published application port, default `80` |
| `JWT_SECRET` | Required random HS256 signing secret, at least 32 characters |
| `POSTGRES_PORT` | Loopback host port used only by the development override, default `55432` |
| `PORT` | Native backend listening port, default `8080`; set from `BACKEND_PORT` in Compose |
| `API_PROXY_TARGET` | Optional local Vite proxy target, default `http://localhost:8080` |

Native Go execution still supports `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, and optional `DB_SSLMODE` when `DATABASE_URL` is unset. `DATABASE_URL` takes precedence. Credentials are URL-encoded when building the legacy connection string. Container configuration uses one `DATABASE_URL` for both Atlas and the API.

The reverse-proxy Nginx configuration is rendered at startup using only `BACKEND_PORT`; all Nginx variables are preserved. API requests retain `/api`, proxy headers, HTTP/1.1 upgrade support, and unbuffered streaming. Docker DNS is re-resolved for both upstreams so replacing a backend or frontend container does not leave Nginx pointing to its old IP. The frontend serves nested SPA routes with `index.html`, returns 404 for missing assets, and preserves asset and HTML cache policies through the proxy. Nginx starts after both upstreams are healthy; its healthcheck requests both the SPA and API readiness endpoint.

## Container resource limits

The base Compose file enforces CPU, memory, process/thread, and open-file ceilings for every service, including the one-shot migration job. These also apply when using the development database override and `make migrate`.

| Service | CPU quota | Memory ceiling | PID ceiling | Open files per process |
| --- | --- | --- | --- | --- |
| PostgreSQL | 1 CPU | 1 GiB | 256 | 4096 |
| Backend | 1 CPU | 512 MiB | 128 | 4096 |
| Frontend | 0.5 CPU | 128 MiB | 64 | 4096 |
| Nginx | 0.5 CPU | 128 MiB | 64 | 4096 |
| Migration job | 0.5 CPU | 256 MiB | 64 | 4096 |

These are starting defaults for a small deployment, not capacity guarantees. CPU quotas allow fractional CPUs and throttle sustained use rather than reserving a core. Memory limits are hard ceilings; exceeding them can trigger the container's OOM killer. Each service's combined memory-and-swap limit equals its memory limit, preventing container swap use on supported hosts. PID limits include threads; file limits apply separately to each process.

See Docker's [Compose service reference](https://docs.docker.com/reference/compose-file/services/) for the runtime settings used here.

Override `<SERVICE>_CPUS`, `<SERVICE>_MEMORY_LIMIT`, and `<SERVICE>_PIDS_LIMIT` in the root `.env` using the names in `.env.example`. Memory values accept units such as `512m` and `1g`. Use positive, finite limits, and size them against representative traffic and migration workloads. The fixed soft and hard `nofile` limits are defined per service in `docker-compose.yml`. Restart policies retain automatic recovery for long-running services; migration jobs never restart automatically.

After changing limits, validate and recreate the containers to apply them:

```bash
docker compose config --quiet
docker compose up -d --wait --wait-timeout 120
docker compose ps --all
docker stats --no-stream $(docker compose ps -q)
```

`docker compose restart` alone does not apply changed resource configuration. These are runtime limits; Docker image builds and native Go/Vite processes use their own host resources.

## Test and operate

```bash
make up                 # docker compose up -d
make down               # stop, preserve database volume
make logs               # follow service logs
make build              # build production images
make test               # Go tests/vet/build and frontend production build
make migrate            # apply pending migrations using the migration image
```

Standalone validation:

```bash
cd backend
go test ./...            # includes Testcontainers PostgreSQL and auth/user requests
go vet ./...
go build ./cmd/api
# Fast checks without Docker:
go test -short ./...
cd ../frontend/vue-ui
npm run build
cd ../..
docker compose config --quiet
docker compose build
docker compose up -d
docker compose ps --all
docker compose logs nginx frontend backend postgres
```

Backend integration tests initialize PostgreSQL 17 from the repository's actual SQL migrations, then use the same `/api` router as the binary. Unit tests cover configuration precedence and credentials, database failure readiness, and JWT signature/expiration checks. The frontend production build type-checks with `vue-tsc`.

See [DEPENDENCIES.md](DEPENDENCIES.md) for the dependency audit, version decisions, and major-version compatibility checks. [VALIDATION.md](VALIDATION.md) records the executed checks, results, and exact changed-file manifest.
