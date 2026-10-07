# Implementation validation

Validated on 2026-10-07 against the repository worktree and local Docker Engine. The reference architecture was inspected at `/home/joaod/gungnr/templates/webgalidor` before edits.

## Results

| Commands / checks | Result |
| --- | --- |
| Initial `go version`, `go list -m -u all`, `go mod tidy`, `go vet ./...`, `go test ./...` | Passed; dependency audit recorded in DEPENDENCIES.md |
| Atlas 1.3.0 isolated `-modfile` compatibility `go test ./...` / `go vet ./...` | Passed against the original code before changing the major-version requirement |
| Final `go mod tidy`, `go test ./...`, `go vet ./...`, `go build ./cmd/api` | Passed on Go 1.27.1; includes PostgreSQL 17 Testcontainers integration |
| `make test` at root | Passed: Go tests/vet/build and frontend type-check/production build |
| `go generate ./ent` | Passed; generated Ent files unchanged |
| `make atlas-install` | Passed on Linux amd64; pinned CLI reports 1.3.0 |
| `atlas migrate validate --dir file://ent/migrate/migrations` | Passed; existing checksum and SQL unchanged |
| `atlas migrate diff compatibility_check --dir file:///tmp/ent-go-vue-migration-compat --to ent://ent/schema --dev-url docker://postgres/17-alpine/dev?search_path=public` | Passed against a copy of the migration directory: no schema changes detected |
| `npm outdated`, `npm update`, `npm install`, `npm audit` | Compatible updates applied; zero vulnerabilities; newer TypeScript/Node-types majors deliberately retained at current ranges |
| `npm run build` | Passed: vue-tsc and Vite production build |
| `docker compose config --quiet` | Passed for base and explicit development override |
| `docker compose build` | Passed for backend, migration, and Nginx production images |
| `docker compose up -d`, `docker compose ps --all` | Nginx/backend/PostgreSQL healthy; migration job exited 0 |
| `docker compose exec -T nginx nginx -t` | Nginx configuration valid |
| `docker compose run --rm migrate` | First run applied two migrations; repeats report no files to execute |
| `HTTP_PORT=18082 docker compose -p ent-go-vue-fresh-check up -d --build --wait --wait-timeout 90` | Passed from a fresh named volume: database initialization → both Atlas migrations → healthy backend → healthy Nginx |
| `curl --fail http://localhost:18080/api/health` | HTTP 200, `status: ok` |
| `curl --fail http://localhost:18080/api/ready` | HTTP 200, `status: ok`, `database: connected` |
| `curl --fail http://localhost:18080/api/users` | HTTP 200; actual Ent repository queried PostgreSQL |
| HTTP auth/user smoke requests through Nginx | Register 201; login, protected profile, public detail/list 200; password absent from response |
| `/api/openapi.json`, nested SPA route, unknown API route | Correct `/api` OpenAPI server; SPA fallback 200; unknown API path 404 rather than HTML fallback |
| Native `PORT=18081 go run ./cmd/api`, `API_PROXY_TARGET=http://localhost:18081 npm run dev -- --host 127.0.0.1 --port 15173 --strictPort` | Both processes ran with loopback development PostgreSQL; browser health and readiness succeeded through Vite |
| Headless Chromium via Playwright, native and production | Backend/database visibly connected; same-origin `/api` requests; refresh action; simulated API failure and recovery; no page exceptions; no overflow at 375px |
| `BACKEND_PORT=8082 docker compose up -d` | Backend ready and Nginx rendered/proxied port 8082 successfully; final stack restored to configured port 8080 |
| Real `docker compose stop postgres` / `start postgres` | Health stayed 200; readiness became 503; browser showed database unavailable while backend remained connected; recovered to connected without restarting the API |
| `docker compose down` / `docker compose up -d --wait --wait-timeout 90` | A user created through the API remained readable after container recreation, proving named-volume persistence |
| `docker compose logs nginx backend postgres`, final container inspection | Healthy services, zero automatic restarts, no crash loops; only Nginx has host port bindings |
| `git diff --check` | Passed |

## Local runtime state

### Resource-limit follow-up

Validated on 2026-10-07 after adding limits to `docker-compose.yml`, configurable defaults to `.env.example`, and operating guidance to `README.md`. No application code or images changed in this follow-up.

- Base and development-override Compose configuration validation passed. Parsed configuration assertions verified CPU, memory, memory+swap, PID, and `nofile` defaults for all four services. Independent environment override checks also passed for each service's CPU, memory, and PID settings.
- `docker compose up -d --wait --wait-timeout 120` passed, recreating the application containers with their resource limits while preserving the PostgreSQL volume. Backend, Nginx, and PostgreSQL were healthy; Atlas exited 0 with no pending migrations.
- `docker inspect` confirmed all four containers had the configured CPU, memory, memory+swap, PID, and open-file ceilings. Live `/sys/fs/cgroup/{memory.max,memory.swap.max,cpu.max,pids.max}` checks confirmed enforcement for every running service, with swap capped at zero. PID 1's `/proc/1/limits` confirmed soft/hard `nofile` limits of 4096.
- HTTP checks through Nginx passed for `/`, a nested SPA route, `/api/health`, `/api/ready`, and `/api/users`. Readiness reported the database connected, and the users endpoint exercised Ent against PostgreSQL.
- A bounded concurrent smoke check passed 300/300 health, readiness, and user-list requests with 12 workers. This verifies basic operation under the limits; it is not a production capacity benchmark.
- Container checks found no OOM kills or automatic restarts. Startup/migration logs were clean, and `git diff --check` passed.

The main application is left running at **http://localhost:18080/**. Port 80 already belonged to another application, so the ignored local root `.env` sets `HTTP_PORT=18080`; the committed example and Compose default remain 80. The backend and PostgreSQL have no host port bindings. The application volume is retained.

The temporary fresh-check Compose project and its test-only volume were removed. The temporary user used for auth/persistence acceptance was deleted. Native Go/Vite validation processes were stopped. Random local credentials are in ignored `.env` files, not committed examples. Atlas's installed `backend/bin/atlas` and build outputs are also ignored.

## Exact changed files

`A` = added, `M` = modified, `D` = deleted. Generated Ent source and migration SQL/checksums were preserved; only the generator directive changed.

- `A` `.dockerignore`
- `A` `.env.example`
- `A` `.gitignore`
- `A` `DEPENDENCIES.md`
- `A` `Makefile`
- `A` `README.md`
- `A` `VALIDATION.md`
- `A` `backend/.dockerignore`
- `M` `backend/.env.example`
- `A` `backend/Dockerfile`
- `M` `backend/Makefile`
- `M` `backend/README.md`
- `A` `backend/atlas.hcl`
- `M` `backend/cmd/api/main.go`
- `D` `backend/docker-compose.yml`
- `M` `backend/ent/generate.go`
- `M` `backend/go.mod`
- `M` `backend/go.sum`
- `M` `backend/internal/auth/middleware/jwt.go`
- `A` `backend/internal/auth/middleware/jwt_test.go`
- `M` `backend/internal/auth/usecase/auth.go`
- `A` `backend/internal/config/config.go`
- `A` `backend/internal/config/config_test.go`
- `A` `backend/internal/httpapi/router.go`
- `A` `backend/internal/httpapi/router_test.go`
- `M` `backend/internal/user/delivery/http/my_user.go`
- `M` `backend/internal/user/delivery/http/user_test.go`
- `A` `docker-compose.dev.yml`
- `A` `docker-compose.yml`
- `A` `docker/nginx/nginx.conf`
- `A` `frontend/vue-ui/.dockerignore`
- `A` `frontend/vue-ui/.env.example`
- `A` `frontend/vue-ui/Dockerfile`
- `M` `frontend/vue-ui/README.md`
- `M` `frontend/vue-ui/package-lock.json`
- `M` `frontend/vue-ui/package.json`
- `M` `frontend/vue-ui/src/App.vue`
- `A` `frontend/vue-ui/src/api/client.ts`
- `A` `frontend/vue-ui/src/components/ConnectionStatus.vue`
- `A` `frontend/vue-ui/src/queries/health.ts`
- `M` `frontend/vue-ui/vite.config.ts`
