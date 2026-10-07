# Ent Go Vue

A full-stack boilerplate for starting Go and Vue projects. Includes registration, JWT authentication, optional administrator access, user management, database migrations, and Docker Compose setup.

- **Backend:** Go, Huma, Ent, Atlas, and PostgreSQL.
- **Frontend:** Vue 3, TypeScript, Vite, Tailwind CSS, and daisyUI.

## Quick start

Requires Docker with the Compose plugin. Run commands from the repository root.

```bash
cp .env.example .env
```

Edit `.env`:

- Set `POSTGRES_PASSWORD` and update `DATABASE_URL` with matching credentials. Keep the hostname as `postgres`; percent-encode special characters in URL credentials.
- Set `JWT_SECRET` to at least 32 random characters. Generate one with `openssl rand -hex 32`.
- Optionally set both `ROOT_EMAIL` and `ROOT_PASSWORD` for an administrator account. Use an 8–72 byte password; startup reapplies it.

```bash
docker compose up --build -d
```

Open <http://127.0.0.1:18080/> and register an account or sign in as the configured administrator. Change `HTTP_PORT` in `.env` to use another port. API documentation is available at <http://127.0.0.1:18080/api/docs>.

Database migrations run automatically when the stack starts. Database data persists when you stop it.

## Common commands

| Command | Purpose |
| --- | --- |
| `docker compose up --build -d` | Start the stack or rebuild after source changes |
| `make down` | Stop the stack and keep database data |
| `make logs` | Follow service logs |
| `make migrate` | Apply pending database migrations |
| `make check` | Run lint, complexity, duplication, tests, and builds |

`docker compose down -v` also deletes database data.

## Local development

Requires Docker, GNU Make, Go 1.26.4+ (the module selects toolchain 1.27.1), Node 22.13+ or 24+, and npm.

Configure the root `.env` as above, then start the database and apply migrations:

```bash
make dev-db
make migrate
cp backend/.env.example backend/.env
```

In `backend/.env`, set matching database credentials and a random JWT secret. Use `localhost:55432` in `DATABASE_URL`, or the port configured by `POSTGRES_PORT`.

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

Open the URL printed by Vite. It proxies `/api` to `http://localhost:8080`. To use a different API port, copy `frontend/vue-ui/.env.example` to `frontend/vue-ui/.env` and update `API_PROXY_TARGET`.

Before switching back to the full Docker stack, stop the development containers with `docker compose -f docker-compose.yml -f docker-compose.dev.yml down`.

## Test and operate

Run `make check` before submitting changes. It requires the local development tools, a running Docker daemon for integration tests, and golangci-lint 2.14.0:

```bash
curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/v2.14.0/install.sh -o /tmp/ent-go-vue-install-golangci-lint.sh
sh /tmp/ent-go-vue-install-golangci-lint.sh -b backend/bin v2.14.0
npm --prefix frontend/vue-ui ci
make check
```

For quick backend tests without Docker, run `make -C backend test`.

See [backend commands and schema migrations](backend/README.md), [frontend development](frontend/vue-ui/README.md), [.env.example](.env.example) for configuration, and [dependency notes](DEPENDENCIES.md) for version details.
