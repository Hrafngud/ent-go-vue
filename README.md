# Ent Go Vue

A full-stack boilerplate for building applications with Go and Vue. Includes registration, JWT authentication, an optional administrator account, user management, database migrations, and a Docker Compose stack.

- **Backend:** Go, `net/http`, Huma, Ent, Atlas, and PostgreSQL.
- **Frontend:** Vue 3, TypeScript, Vite, Tailwind CSS, daisyUI, and Phosphor Icons. TanStack Vue Query manages server state; Pinia manages client state.

## Architecture

```text
Browser → Cloudflare Tunnel or VPS Nginx (HTTPS) → application Nginx (HTTP)
          ├── /      → Vue frontend
          └── /api/* → Go API → PostgreSQL
```

The frontend and API share one origin through a reverse proxy. The backend groups features under `internal/auth` and `internal/user`, with domain, use case, repository, and HTTP layers. Ent defines the schema; Atlas applies committed SQL migrations before the API starts.

Inputs are validated in both layers, API bodies and request rates are limited, and user management requires root access.

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
- `HTTP_PORT` defaults to `18080`, published only on `127.0.0.1` for the outer proxy.

```bash
docker compose config --quiet
docker compose up --build -d
docker compose ps --all
```

Open <http://127.0.0.1:18080/> (or your configured port). Register an account or sign in with the configured administrator. API documentation is at `/api/docs`; OpenAPI is at `/api/openapi.json`. Use `/api/health` for liveness and `/api/ready` for database readiness. Existing `.env` files retain their chosen `HTTP_PORT`.

Compose starts PostgreSQL, applies migrations, then starts the API, frontend, and proxy. The migration container exits successfully when done. Only Nginx publishes an application port; database data persists in `postgres_data`.

```bash
docker compose logs -f
docker compose down          # stop and preserve database data
```

Use `docker compose up --build -d` after source changes. `docker compose down -v` deletes database data. Database initialization credentials only take effect when the volume is first created.

## External HTTPS and proxy trust

TLS terminates at Cloudflare or the VPS-wide Nginx. Application Nginx stays HTTP,
uses relative redirects, owns the common response headers, and preserves route CSP
and cache policies. The API trusts only application Nginx, using its fixed internal
address. Application Nginx trusts no outer forwarding headers by default.

Set `TRUSTED_EDGE_CIDRS` to the exact outer proxy socket peers **as seen by application
Nginx**. For host-local VPS Nginx or `cloudflared` through the published Docker port,
this is commonly the bridge gateway (for example `172.30.90.1/32`); confirm it in
application Nginx access logs before enabling trust. List multiple numeric CIDRs
with commas. For a containerized connector, give it a fixed private address and
trust that `/32` (or IPv6 `/128`). Do not trust all addresses or the whole application
subnet. Empty trust is safe but visitors behind one edge share its IP quota until
the deployment configures the correct peer.

Trusted edges must provide a trustworthy `X-Forwarded-For` chain and a single
`X-Forwarded-Proto: http` or `https`. Nginx restores the last untrusted visitor IP,
forwards only that IP to the API, and accepts scheme hints only from a trusted socket
peer. Untrusted hints are ignored; alternative forwarding fields are stripped.
Both proxy locations inherit this normalization; adding location-level
`proxy_set_header` directives requires preserving the complete set.

For **VPS Nginx + Certbot**, use
[vps-edge.conf.example](docker/nginx/vps-edge.conf.example) as a starting point.
Replace the hostname/paths/port, obtain the certificate before loading its HTTPS
block, validate with `nginx -t`, and check automatic renewal with `certbot renew
--dry-run`. The HTTP listener permits ACME validation and redirects application
traffic to a fixed HTTPS hostname. The HTTPS listener owns one HSTS field, including
errors. Assign a deployment owner to observe the initial one-day HSTS rollout and
increase its duration after certificate renewal and route checks pass.

For **Cloudflare Tunnel**, point a host-local connector at
`http://127.0.0.1:18080`, or a containerized connector at `http://nginx:8080` on a
private network. Configure Cloudflare HTTPS redirects and HTTPS-only HSTS for the
intended hostname. Check the actual forwarded IP chain before enabling peer trust;
Cloudflare appends visitor/proxy information to `X-Forwarded-For` and overwrites
`X-Forwarded-Proto`. Direct public access to application Nginx must remain blocked.
Tunnel development should target application Nginx to receive its document policy.

Keep HSTS ownership at the external edge. Decide `includeSubDomains` only after
inventorying affected hosts, delegated services, and HTTPS readiness. Decide preload
separately; emitting the directive does not enroll a domain. Avoid zone-wide changes
that would affect unrelated development hosts. Nginx suppresses its version but
still emits `Server: nginx`; product-header removal depends on outer-edge support.

After deployment, check HTTP redirects and HTTPS headers for `/`, `/api`, `/api/docs`,
a real asset, a missing asset, and an unauthenticated API request. Confirm redirects
preserve the external scheme/port, HSTS appears exactly once on HTTPS, API responses
retain `no-store`, and per-visitor limits remain distinct. Test outer-edge error
responses too. Local proxy regression checks run with
`python3 docker/nginx/security_test.py` after building its image.

References: [Nginx real-IP handling](https://nginx.org/en/docs/http/ngx_http_realip_module.html),
[Cloudflare forwarding headers](https://developers.cloudflare.com/fundamentals/reference/http-headers/),
[Cloudflare HTTPS redirects](https://developers.cloudflare.com/ssl/edge-certificates/additional-options/always-use-https/),
[Cloudflare HSTS](https://developers.cloudflare.com/ssl/edge-certificates/additional-options/http-strict-transport-security/),
and [OWASP HSTS guidance](https://cheatsheetseries.owasp.org/cheatsheets/HTTP_Strict_Transport_Security_Cheat_Sheet.html).

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
