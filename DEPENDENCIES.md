# Dependency audit

Audited on 2026-10-07. The existing stack and application architecture are retained.

## Backend

Baseline commands: `go version`, `go list -m -u all`, `go mod tidy`, `go vet ./...`, and `go test ./...`. The initial installed Go version was `go1.27.0-X:nodwarf5`; baseline vet and tests passed, including the Testcontainers auth/user integration test.

| Dependency | Before | Decision |
| --- | --- | --- |
| Go module / toolchain | `go 1.25.0`, no toolchain | `go 1.26.4`, `toolchain go1.27.1`; Docker builder also uses 1.27.1 |
| Ent | 0.14.6 | Retained; current release and existing generated client remain compatible |
| Huma v2 | 2.36.0 | 2.39.1; stays on v2 with the existing standard-library adapter |
| `lib/pq` | 1.12.3 | Retained; current release and existing PostgreSQL driver works |
| Atlas library | 0.36.2 prerelease | 1.3.0, aligned with a pinned Atlas 1.3.0 CLI/image |
| Testcontainers + PostgreSQL module | 0.42.0 | Both 0.44.0; matched versions, tested with PostgreSQL 17 |
| Google UUID | 1.6.0 | Retained; no newer compatible release in the audit |
| godotenv | 1.5.1 | Retained; no newer compatible release in the audit |
| JWT v5 | 5.3.1 (previously marked indirect) | Retained; classified as direct by tidy |
| `x/crypto` | 0.52.0 (previously marked indirect) | 0.57.0, current maintained bcrypt dependency; requires Go 1.26.4 |

Go's [release policy and history](https://go.dev/doc/devel/release) guided moving off the older 1.25 release line. Huma [release notes](https://github.com/danielgtaylor/huma/releases/tag/v2.39.1) and Testcontainers [release notes](https://github.com/testcontainers/testcontainers-go/releases/tag/v0.44.0) were reviewed; the existing adapter and container API calls remain compatible. Ent and `lib/pq` release listings confirm their retained versions.

Atlas is the meaningful major-version change. Before updating the repository's requirement, a temporary `-modfile` selected Atlas 1.3.0 and the **existing** backend passed `go test ./...` and `go vet ./...`. The new Docker CLI then validated/applied the unchanged `atlas.sum` and both SQL files to PostgreSQL 17. The API still never applies schema synchronization at startup. See Atlas's [versioned apply workflow](https://atlasgo.io/versioned/apply).

An extra `go generate ./ent` check exposed incompatible transitive iterator APIs between the older text-width package and its selected Unicode dependency. Ent's CLI is now recorded with Go's `tool` directive and invoked with `go tool ent`. Updating its `tablewriter` dependency from 1.1.3 to 1.1.5 selects compatible `displaywidth` 0.10.0. Regeneration passes without changes to generated Ent source or migration history; tool dependencies survive `go mod tidy`.

Indirect updates are those required by the chosen direct libraries and generator. The entire transitive graph was not upgraded merely to match every entry in `go list -m -u all`.

## Frontend

Ran `npm outdated`, `npm update`, `npm audit`, `npm install`, and `npm run build`. Kept the declared compatible ranges and all existing major versions. `npm update` refreshed the lockfile within those ranges; `npm audit` reports **zero vulnerabilities**. The production builder uses Node 24; native development supports Node 22.12+ or 24+.

| Direct dependency | Final installed version |
| --- | --- |
| Vue | 3.5.43 |
| Vite | 8.3.3 |
| TypeScript | 6.0.3 |
| Tailwind / Tailwind Vite plugin | 4.3.3 |
| daisyUI | 5.7.47 |
| Pinia | 4.0.3 |
| TanStack Vue Query | 5.104.1 |
| Vite Vue plugin | 6.0.9 |
| Vue TS config | 0.9.1 |
| vue-tsc | 3.3.12 |
| Node type definitions | 24.19.1 |

`npm outdated` offered newer majors for TypeScript and Node types. Neither has a concrete benefit for this integration, so they remain on their current compatible major versions. No `npm audit fix --force` or framework migration was used.

## Container choices

- PostgreSQL 17 Alpine, matching the reference project's supported database generation without introducing a PostgreSQL data-directory migration.
- Go 1.27.1 Alpine builder, compiling a static API into Alpine 3.23 with certificates and a non-root user.
- Atlas 1.3.0 only in the separate migration image; no development tooling in the API runtime.
- Node 24 Alpine builder; unprivileged Nginx 1.30 Alpine runtime serving built assets in the frontend container, plus a separate unprivileged Nginx 1.30 Alpine reverse-proxy container.

Application containers drop capabilities and use `no-new-privileges`. PostgreSQL has a persistent named volume and a healthcheck; only Nginx has a published host port in the base Compose configuration.
