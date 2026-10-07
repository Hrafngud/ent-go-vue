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

### Separate frontend container follow-up

Validated on 2026-10-07 after inspecting the approved reference at `/home/joaod/Documentos/Github/Galidor`. The integrated static-serving/proxy image was split into a `frontend` image and a separate `nginx` reverse-proxy image. Nginx routes `/` to `frontend:8080` and `/api/` to the backend, publishing the existing application port.

- Base and development-override `docker compose config --quiet` checks passed.
- `docker compose build frontend nginx` passed, including the frontend's `vue-tsc` type check and Vite production build.
- `docker compose up -d --wait --wait-timeout 120` passed. Frontend, reverse proxy, backend, and PostgreSQL were healthy; the migration job exited 0. The existing database volume was retained.
- `docker compose exec -T frontend nginx -t` and `docker compose exec -T nginx nginx -t` passed.
- HTTP checks at `http://localhost:18080/` passed for the SPA, `/index.html`, a nested SPA route with a query string, both built JS/CSS assets, `/api/health`, `/api/ready`, `/api/users`, `/api/docs`, and `/api/openapi.json`. Readiness reported the database connected. Missing assets and unknown API paths returned 404; unknown API paths did not receive the SPA fallback.
- HTML retained `Cache-Control: no-cache`; built assets retained their one-year cache policy through the proxy. Security headers were present on successful and error responses.
- Docker inspection confirmed only `nginx` publishes a host port. The frontend runs as a non-root user with 0.5 CPU, 128 MiB memory/memory+swap, 64 PIDs, and soft/hard `nofile` limits of 4096. Running containers had no OOM kills or automatic restarts.
- `git diff --check` passed.

This follow-up changes `.env.example`, `docker-compose.yml`, `docker/nginx/nginx.conf`, `frontend/vue-ui/Dockerfile`, `README.md`, `frontend/vue-ui/README.md`, `DEPENDENCIES.md`, and this validation record. It adds `docker/nginx/Dockerfile` and `frontend/vue-ui/nginx.conf`. Application source, dependencies, and migrations are unchanged.

### Resource-limit follow-up

Validated on 2026-10-07 after adding limits to `docker-compose.yml`, configurable defaults to `.env.example`, and operating guidance to `README.md`. No application code or images changed in this follow-up.

- Base and development-override Compose configuration validation passed. Parsed configuration assertions verified CPU, memory, memory+swap, PID, and `nofile` defaults for all four services. Independent environment override checks also passed for each service's CPU, memory, and PID settings.
- `docker compose up -d --wait --wait-timeout 120` passed, recreating the application containers with their resource limits while preserving the PostgreSQL volume. Backend, Nginx, and PostgreSQL were healthy; Atlas exited 0 with no pending migrations.
- `docker inspect` confirmed all four containers had the configured CPU, memory, memory+swap, PID, and open-file ceilings. Live `/sys/fs/cgroup/{memory.max,memory.swap.max,cpu.max,pids.max}` checks confirmed enforcement for every running service, with swap capped at zero. PID 1's `/proc/1/limits` confirmed soft/hard `nofile` limits of 4096.
- HTTP checks through Nginx passed for `/`, a nested SPA route, `/api/health`, `/api/ready`, and `/api/users`. Readiness reported the database connected, and the users endpoint exercised Ent against PostgreSQL.
- A bounded concurrent smoke check passed 300/300 health, readiness, and user-list requests with 12 workers. This verifies basic operation under the limits; it is not a production capacity benchmark.
- Container checks found no OOM kills or automatic restarts. Startup/migration logs were clean, and `git diff --check` passed.

The main application is left running at **http://localhost:18080/**. Port 80 already belonged to another application, so the ignored local root `.env` sets `HTTP_PORT=18080`; the committed example and Compose default remain 80. The frontend, backend, and PostgreSQL have no host port bindings. The application volume is retained.

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

## Registration and root user management (2026-10-07)

Added routed login/registration/workspace pages under `src/pages`, admin list/create/detail/edit pages under `src/pages/admin`, and reusable daisyUI layouts, forms, alerts, and user components under `src/components`. The hierarchy and route access policy are documented in `frontend/vue-ui/README.md`.

Validation passed:

- Backend `go test ./...`, including PostgreSQL Testcontainers integration coverage of all root-only CRUD endpoints, 401 without authentication, 403 for members, disabled root access, duplicate email conflicts, input validation (including bcrypt byte limits), hashed passwords, password preservation/rotation, root email/deletion protection, and missing-user responses.
- Backend `go vet ./...` and API binary build.
- Frontend `vue-tsc` and Vite production builds, including the final responsive layout build inside Docker.
- Headless Chromium against the rebuilt local app at `http://127.0.0.1:18080`: registration/password confirmation, empty 201 response handling, login, reload/session restoration, member route guards, real root CRUD, duplicate-email feedback, optional password preservation, password rotation, delete confirmation/cancellation, root protections, and missing-user errors.
- Browser response interception verified pagination/sorting/search, empty lists, service error/retry recovery, and expired-session cleanup. Desktop (1280px) and mobile (375px) screenshots were inspected. Mobile users render as a stacked daisyUI list with visible View/Edit actions and full-width search; no page overflow or browser runtime errors occurred.
- Temporary browser-test accounts were deleted afterward. Rebuilt backend/frontend/proxy and PostgreSQL services report healthy. No database schema migration was needed.

## Local quality gate

Implemented and validated on 2026-10-07. The gate is local-only: root `make lint`,
`make complexity`, `make duplication`, and canonical `make check`. `make check`
serializes those stages and then reuses `make test` (uncached full Go tests,
auth/user PostgreSQL integration, vet, API build, Vue type-check/production build).
No CI workflows or unrelated application refactoring were added.

Prerequisites and installation are in [README.md](README.md#test-and-operate).
Validated with Go 1.27.1, Node 22.17.0, npm 11.6.4, Docker Engine 29.7.2,
golangci-lint 2.14.0 (official release checksum verified), ESLint 10.12.0,
typescript-eslint 8.71.1, eslint-plugin-vue 10.11.1, vue-eslint-parser 10.3.0,
and jscpd 5.4.0. New frontend tooling is pinned in the lockfile; existing package
versions were unchanged. ESLint 10 requires raising the native Node minimum from
22.12 to 22.13. A fresh `npm ci` succeeded and reported zero vulnerabilities.

### Pre-existing debt

The first unfiltered Go scan, before any application source changes, reported the
following. `backend/.golangci.yml` contains narrowly scoped exceptions so this
integration preserves the existing code rather than refactoring it to satisfy a
new tool. The general threshold stays 15, including tests.

| Existing finding | Baseline treatment |
| --- | --- |
| `internal/config/config.go:21`, `Load`: complexity 24 | Exception matches this path, function, and exact complexity 24; any different value above 15 fails |
| `internal/user/delivery/http/user_test.go:27`, `TestUserAndAuthAPI_Integration`: complexity 70 | Exception matches this path, function, and exact complexity 70 |
| `cmd/api/main.go:49`: unchecked deferred `client.Close()` | Exception matches the file, specific errcheck message, and exact source expression |
| `internal/user/delivery/http/user_test.go:66,72`: unchecked deferred client/database close | Same narrow path/message/source treatment |
| `internal/user/delivery/http/user_test.go:83,137,156,169,198,226,238,263,289`: unchecked response-body close | Exception restricted to response-body cleanup calls in this existing integration-test file |

Cleanup exceptions also match identical expressions added in these same files;
reviewers must not extend that debt. Other unchecked calls and paths still fail.
Remove exceptions when the underlying debt is fixed; do not broaden them to pass
future tasks. There were no Go clones at 100 tokens. Frontend ESLint and complexity
passed without debt exceptions. The Vue-only `no-useless-assignment` override
addresses the [known template-read false positive](https://github.com/vuejs/eslint-plugin-vue/issues/2660);
unused bindings still receive TypeScript and typescript-eslint checks.

jscpd at 50 tokens/5 lines reports one existing seven-line template clone between
`UserDetailPage.vue:22-28` and `UserEditPage.vue:45-51`, approximately 0.4% duplication.
The initial 5% ceiling tolerates this small amount of declarative repetition while
still detecting substantial copy/paste. Existing code was not changed to erase it.

### Executed checks and probes

| Check | Result |
| --- | --- |
| `golangci-lint config verify` | Passed |
| Root `make lint` | Passed: Go static analysis, ESLint, `vue-tsc` |
| Root `make complexity` | Passed: cyclop and ESLint complexity at 15, with the documented Go baseline |
| Root `make duplication` | Passed: dupl at 100 tokens; jscpd at 50 tokens/5 lines and 5% ceiling |
| Root `make test` | Passed: uncached PostgreSQL-backed integration tests, unit tests, vet, API build, frontend type-check and Vite production build |
| Final root `make check` after fresh `npm ci` and removal of all probes | Passed, exit 0: every quality stage, uncached full Go tests including PostgreSQL integration (13.773s), vet/API build, and Vue type-check/Vite production build; jscpd reported 0.43%, below 5% |
| Temporary Go function at complexity 16 | `make complexity` exited 2 with cyclop finding |
| Temporary TS and Vue script functions at complexity 16 | Each `make complexity` exited 2 with ESLint complexity finding |
| Existing `Load` temporarily raised from 24 to 25 | `make -C backend complexity` exited 2; exact-value baseline does not allow increases |
| Temporary duplicated Go and frontend source blocks | Each `make duplication` exited 2 with its detector's finding |
| Temporary TypeScript type mismatch | `make lint` exited 2 with TS2322 |
| Temporary explicit-any lint violation | Root `make check` exited 2 at lint, without reaching later stages |
| Ent exclusion fixture with complexity, duplication, and unused code, without generated marker | All three individual root checks passed; Ent is never selected as a lint target |
| Parsing-error/duplicate fixtures in frontend dependency/build/coverage/generated directories | All three individual root checks passed; ESLint API confirmed every fixture had no applicable config |
| `git diff --check` | Passed; no application source or generated Ent changes |

All temporary source fixtures were removed, and the original configuration source
was restored byte for byte after its complexity probe. Generated Ent source and
migrations were unchanged. The ignored native golangci-lint binary remains in
`backend/bin/` so the documented root commands work locally.
