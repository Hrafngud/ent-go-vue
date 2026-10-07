# Codebase quality standards

These standards apply to every agent task in this repository, including work outside
the `@start.md` loop. Read this file before making changes. Preserve the conventions
below when extending the project; change a standard only when the operator explicitly
requests that change, and update this document alongside the implementation.

This baseline was derived from the frontend and backend source and quality
configuration on 2026-10-07. Automated checks enforce only part of the standard;
architecture, behavior, accessibility, and operational compatibility also require review.

## Sources and baseline findings

| Area | Current convention and authoritative files |
| --- | --- |
| Backend architecture | Feature modules with domain interfaces, use cases, Ent repositories, and Huma HTTP delivery. See `backend/internal/user/`, `backend/internal/auth/`, and `backend/internal/httpapi/router.go`. |
| Frontend architecture | Vue Composition API with TypeScript; pages compose reusable components. HTTP calls, server state, client state, and validation have separate modules under `frontend/vue-ui/src/`. |
| UI conventions | FormKit forms, daisyUI controls, Tailwind utilities, semantic theme colors, shared modal/loading/feedback components, and individual Phosphor icon imports. See `src/forms/formkit.ts`, `src/components/`, and `src/style.css` in the frontend. |
| Automated checks | Root [Makefile](../Makefile), backend [Makefile](../backend/Makefile) and [.golangci.yml](../backend/.golangci.yml), frontend [package.json](../frontend/vue-ui/package.json), [ESLint config](../frontend/vue-ui/eslint.config.js), [TypeScript config](../frontend/vue-ui/tsconfig.app.json), and [jscpd config](../frontend/vue-ui/.jscpd.json). |
| Runtime and migrations | Same-origin `/api`, separate frontend and reverse proxy, PostgreSQL 17, and committed Atlas SQL applied before API startup. See [README.md](../README.md), Compose files, and `backend/ent/migrate/migrations/`. |

The full root gate currently includes backend unit and PostgreSQL integration tests and a small
frontend Node test suite for account validation and retry delays. It does not include
automated browser, accessibility, race, coverage, container-build, or dependency-security
checks. Do not describe those as passing unless they were separately run.

## Mandatory local quality gate by task scope

For backend, shared API-contract, or runtime/infrastructure changes, run from the
repository root:

```sh
make check
```

For frontend-only tasks, focus on code quality and run only the frontend build and
static checks:

```sh
cd frontend/vue-ui
npm run lint
npm run typecheck
npm run complexity
npm run duplication
npm run build
```

These commands are sufficient final validation for frontend-only work; the full root
gate and test suites are not required for that scope. Perform browser testing only
when explicitly demanded by the operator. Do not launch a browser, take screenshots,
or run browser automation as a default frontend validation step. Preserve accessibility
and responsive behavior through code review even when browser testing is not requested.
If a task also changes backend behavior, a shared API contract, or infrastructure, use
the full root gate and the applicable runtime checks below.

The workflow is: implement, run the gate for the task scope, fix failures, and repeat
until clean. Run it after the final change. For scopes requiring `make check`, individual
checks and a production build do not replace the complete gate. This frontend-only
exception also applies when workflow files direct agents to run `make check` generally.
Quality controls remain local; do not introduce CI as part of this workflow.

| Stage | Controls enforced today |
| --- | --- |
| `make lint` | Go `errcheck`, `govet`, `ineffassign`, `staticcheck`, and `unused`; frontend recommended JS/TypeScript rules and Vue essential rules, with zero warnings; `vue-tsc -b`. |
| `make complexity` | Go cyclop maximum cyclomatic complexity **15**; ESLint complexity maximum **15** in frontend JS/TS and Vue scripts. The frontend command runs the full ESLint configuration, as does frontend lint. |
| `make duplication` | Go dupl clone detection at **100 tokens**; frontend jscpd with a **5% ceiling**, **50-token** and **5-line** minimums, `mild` mode, and `failOnEmpty: true`. Frontend scanning includes Vue templates and CSS. |
| `make test` | `go test -count=1 ./...`, including PostgreSQL Testcontainers integration; `go vet ./cmd/... ./internal/...`; native API build; frontend `npm test`; `vue-tsc -b` and Vite production build. |

Frontend application typing inherits Vue's strict TypeScript configuration and enables
`noUnusedLocals`, `noUnusedParameters`, `erasableSyntaxOnly`, and
`noFallthroughCasesInSwitch`. The project build checks application source and Vite
configuration; `npm test` executes its explicitly named test file with Node type
stripping. New test files must be wired into that script to run in the gate.

Prerequisites and pinned installation commands are in
[README.md](../README.md#test-and-operate). The current baseline uses golangci-lint
**2.14.0**, the Go toolchain selected by `backend/go.mod` (**1.27.1**, module minimum
**1.26.4**), and Node satisfying frontend `engines.node` (`^22.13.0 || >=24.0.0`).
Use `npm ci` with the committed lockfile. Docker must be available for the full Go suite;
frontend-only validation does not require Docker or Go. Application services and
`.env` files are unnecessary for these gates.

- Fix causes rather than suppressing diagnostics. Do not raise limits, disable rules,
  broaden exclusions, weaken assertions, or skip integration tests to obtain a pass.
- If a prerequisite or environment failure prevents validation, report the command,
  failure, and remaining check explicitly. Never claim completion with a passing gate
  when it did not run successfully.
- `make -C backend test` uses `-short` and skips container integration. It is a quick
  feedback command, not final validation.
- Report the task scope, actual final gate result, and any additional checks performed.

### Exclusions and existing debt

Go lint targets only `./cmd/...` and `./internal/...`; Ent is compiled as a dependency
but excluded from lint, complexity, and clone detection. Frontend checks exclude
dependencies, build output, coverage, generated directories, and declaration files;
jscpd also honors `.gitignore`. Keep handwritten application code inside the checked
paths. Excluded generated code still has to compile and behave correctly.

The current `.golangci.yml` permits only these existing debt cases:

- `TestUserAndAuthAPI_Integration` in `internal/user/delivery/http/user_test.go` has
  complexity **69**. Its exception matches the exact function, path, and diagnostic;
  another over-limit value fails.
- Unchecked `defer client.Close()` in `cmd/api/main.go` and the integration test;
  unchecked `defer db.Close()` and response-body `Close()` cleanup in that test.
  These exceptions match specific paths, messages, and source expressions.

The cleanup exceptions can also match newly added identical expressions in those files;
review must prevent adding more unchecked cleanup. Remove exceptions when their debt
is fixed. There is no blanket exemption for tests and no current complexity exception
for `config.Load`.

Vue files disable `no-useless-assignment` because ESLint cannot see template reads;
TypeScript and Vue-aware unused-binding checks remain active. This parser limitation
does not justify disabling other rules. A passing 5% frontend duplication ceiling is
not permission to introduce avoidable copy/paste.

At this baseline, jscpd reports one seven-line template clone between
`src/components/users/UserRecord.vue` and `src/pages/admin/UsersPage.vue` (0.25% of
scanned lines). It is below the existing ceiling, not a separate exclusion.

`VALIDATION.md` is an ignored, historical local report and may be absent or stale in
other checkouts. Use the committed configurations and this file for current standards.

## Shared implementation standards

- Read the affected feature, its callers, tests, and configuration before editing.
  Follow established patterns and keep changes within the requested scope.
- Keep responsibilities explicit. Extract shared behavior when it has the same meaning
  across callers; do not add abstraction merely to reduce a metric or split a large
  function into meaningless fragments.
- Use descriptive names, small cohesive functions, and comments explaining intent,
  constraints, or non-obvious behavior. Remove unused code and debugging output.
- Match surrounding formatting. Format changed Go files with `gofmt`; frontend source
  convention uses two-space indentation, single-quoted JS/TS strings, and no statement
  semicolons. No Prettier or Go formatting check is currently part of `make check`.
- Preserve existing API and user behavior unless the task changes it. Update consumers,
  tests, and documentation together when a contract changes.
- Keep credentials, tokens, personal data, and real `.env` contents out of source,
  logs, and documentation. Use placeholders in committed examples.
- Add or update meaningful backend regression tests for changed logic and bug fixes
  using Go `testing`. Frontend-only tasks require the build and static checks above;
  add/run frontend tests when the operator requests them, using the existing `node:test`
  approach for unit tests. Documentation-only edits do not need new tests. Do not add
  tests that merely repeat implementation details.
- Update setup instructions and environment examples when commands or configuration
  change. Keep dependency manifests and lockfiles consistent; introduce dependencies
  for a concrete need and preserve the current stack unless directed otherwise.
- Keep the root README short and developer-facing: boilerplate purpose, stack,
  prerequisites, setup/run commands, and essential checks. Put architecture,
  implementation, and detailed deployment guidance in dedicated docs and link to them.
  Update existing instructions rather than appending change history or feature-by-feature
  implementation summaries.

## File and folder structure

Respect the existing organization in every feature and change request. Extend the
current structure before introducing another convention. Place files according to
their responsibility and existing neighbors so a reader can predict where code lives.
Reuse established modules and components rather than creating a parallel organization.

### Backend placement and naming

| Location | Responsibility and convention |
| --- | --- |
| `backend/cmd/api/main.go` | API process entry point and lifecycle; keep feature implementations under `internal/`. |
| `backend/internal/config/` | Environment loading and validation. |
| `backend/internal/httpapi/` | Shared HTTP routing, feature wiring, and request security boundary. |
| `backend/internal/<feature>/` | Feature domain types, errors, validation, and interfaces, following `auth/` and `user/`. |
| `backend/internal/<feature>/usecase/` | Business rules and workflows. |
| `backend/internal/<feature>/repository/` | Thin persistence adapters over generated Ent APIs. |
| `backend/internal/<feature>/delivery/http/` | Huma operations and HTTP request/response mapping. |
| `backend/internal/auth/middleware/` | Existing authentication middleware; keep feature-owned middleware with its owner. |
| `backend/ent/schema/` | Handwritten Ent data models. |
| `backend/ent/` | Generated ORM code and generation entry point; preserve the generator's layout. |
| `backend/ent/migrate/migrations/` | Versioned SQL migrations and `atlas.sum`. |

Use lowercase Go package/directory names and descriptive snake_case filenames, as in
`admin_user.go` and `user_ent.go`. Place tests beside the package they exercise using
`*_test.go`. Extend an existing feature when it owns the behavior; add a feature folder
only for a distinct domain concern. Create only layers the feature actually needs,
following the current boundaries; do not scaffold empty packages or create a package
for every entity, operation, interface, or helper.

### Frontend placement and naming

All application source stays under `frontend/vue-ui/src/`:

| Location | Responsibility and convention |
| --- | --- |
| `App.vue`, `main.ts` | Application composition and plugin/bootstrap wiring. |
| `pages/` | Routed screens, named `...Page.vue`; preserve meaningful groups such as `pages/admin/`. |
| `components/` | Reusable Vue components; preserve `layout/`, `ui/`, `forms/`, and feature groups such as `users/`. |
| `api/` | Shared transport/errors and typed endpoint modules such as `client.ts` and `users.ts`. |
| `queries/` | Server queries, mutations, and cache behavior grouped by resource/concern. |
| `stores/` | Pinia client state grouped by responsibility, such as `session.ts` and `feedback.ts`. |
| `composables/` | Reusable composition/lifecycle behavior, with `use...` filenames such as `useCooldown.ts`. |
| `forms/`, `validation/` | Shared FormKit configuration and reusable validation rules. |
| `router/` | Route definitions, navigation state, and guards. |
| `assets/`, `style.css` | Bundled assets and shared theme/component styles. |

Keep Vue component filenames in PascalCase and ordinary TypeScript modules consistent
with their neighboring descriptive names. Frontend Node tests belong in
`frontend/vue-ui/tests/`; directly served static files belong in `frontend/vue-ui/public/`.
Group related reusable feature components as `components/users/` already does. Do not
introduce a competing top-level feature tree that duplicates `pages/`, `api/`,
`queries/`, or `stores/`, or mix their responsibilities into convenience folders.

### Limits on structural changes

- Add a folder or subdivision only when a real feature boundary, reusable concern, or
  navigation benefit justifies it. A new feature group within the existing convention
  is an extension; replacing that convention is an architectural change.
- Keep related code together and maintain a readable hierarchy. Avoid excessive
  nesting, one-file wrapper folders, fragmented micro-modules, generic dumping grounds
  such as unowned `utils/` or `common/`, and inconsistent parallel names for the same role.
- Do not split files merely to make them smaller, satisfy a metric, or impose another
  project's preferred structure. Split when cohesive responsibilities or reusable
  behavior make the result easier to understand. Preserve useful granularity already
  present in this repository.
- Keep reorganizations scoped to the requested work. A change to the established
  organization requires a concrete reason and operator direction; document the new
  convention here and update affected imports, tests, and path-based tooling together.
- Review every new file and folder in the final diff: its owner, location, name, and
  level of abstraction should be immediately understandable from the existing structure.

## Backend standards

### Boundaries and Go conventions

- Keep features under `backend/internal/<feature>/`. Domain types, sentinel errors,
  and repository/use-case interfaces belong at the feature root. Business rules belong
  in `usecase/`, Ent access in `repository/`, and request/response handling in
  `delivery/http/`. Follow the current constructor-based dependency injection.
- Keep use cases independent of Ent and HTTP delivery. Repositories translate Ent
  records and persistence errors into domain types and errors. HTTP handlers map domain
  failures to Huma errors rather than exposing database internals.
- Wire dependencies and cross-cutting middleware in `internal/httpapi`; keep process
  startup and shutdown in `cmd/api` and environment validation in `internal/config`.
- Pass `context.Context` through handlers, use cases, and database operations. Preserve
  cancellation and deadlines. Use typed private context keys, as in JWT middleware.
- Handle returned errors, including cleanup errors in new code. Use `errors.Is`/`As`
  when identifying errors and `%w` when wrapping them with useful internal context.

### HTTP, authorization, and input contracts

- Register typed Huma operations under `/api` using the `net/http` adapter. Keep operation
  IDs, methods, paths, response statuses, schema tags, and security metadata coherent.
  Reuse the current `data` envelopes for user reads/writes and empty collection arrays.
- Preserve empty **201** registration responses, token login responses, and **204**
  deletion responses; the frontend client supports successful empty bodies.
- Require JWT authentication and backend authorization for protected operations.
  Route guards and hidden frontend controls are convenience, not authorization.
- Preserve HS256-only verification, required expiration, and UUID subject validation.
  Administrator access derives from the configured root email and current user record;
  root email changes and deletion are rejected by the backend.
- Preserve bcrypt hashing, password exclusion from JSON, generic credential failures,
  and the missing-account dummy-hash comparison. Never return or log password hashes.
- Preserve the shared security boundary: 16 KiB bodies, UTF-8 JSON objects, unique keys,
  bounded nesting, valid Unicode escapes, content-type/encoding checks, request quotas,
  `Retry-After`, and unknown-query rejection. New routes must use the same boundary.
- Trust forwarding headers only from configured proxy CIDRs. Preserve bounded limiter
  storage, concurrency protection, and quota failure behavior; do not bypass limits
  for public readiness requests.

### Persistence and lifecycle

- Future backend features must use Ent's schema-first ORM generation workflow. Focus
  agent effort on data modeling, relationships, constraints, and business logic;
  let Ent generate entities, type-safe CRUD/query builders, predicates, and relation
  traversal rather than writing equivalent CRUD persistence plumbing by hand. See the
  official [Ent documentation](https://entgo.io/docs/),
  [code generation guide](https://entgo.io/docs/code-gen/), and
  [generated CRUD API](https://entgo.io/docs/crud/).
- Model fields, edges, indexes, and constraints in `backend/ent/schema/`, then run
  `make -C backend generate` from the repository root (equivalent to
  `cd backend && go generate ./ent`). The existing directive invokes
  `go tool ent generate --feature sql/versioned-migration ./schema` using the module's
  pinned Ent tool. Do not hand-edit generated Ent files. Review and commit generated
  changes with the schema.
- Consume generated Ent builders through thin repository adapters that satisfy domain
  interfaces; keep authorization, business workflows, and HTTP contracts in their
  existing layers. Use generated ORM operations for routine persistence. Custom SQL
  requires a concrete need that the generated API cannot reasonably cover. Ent ORM
  generation does not replace the project's Huma delivery or generate its HTTP handlers.
- Generate migrations with `make -C backend migrate name=describe_schema_change`.
  Review and commit SQL and `atlas.sum`. Root `make migrate` applies committed migrations;
  the backend target generates them. Do not rewrite migrations already applied.
- Apply migrations separately before API startup. Preserve the startup schema check;
  do not add runtime schema synchronization. Integration tests initialize PostgreSQL
  from committed migration SQL rather than generating a schema on the fly.
- Preserve connection-pool bounds, server timeouts, header limits, SIGTERM shutdown,
  and bounded readiness database calls. Liveness must remain independent of database
  availability. Root provisioning must preserve account ID and creation time.

## Frontend standards

### Mandatory daisyUI primitives and Tailwind Sane skill

- Every frontend feature and change request must use daisyUI primitives as the UI
  foundation. Reuse existing shared components built on those primitives first;
  otherwise compose daisyUI controls such as `btn`, `input`, `select`, `modal`, `alert`,
  `table`, and `card`. Do not recreate an available primitive with custom CSS or a
  competing component library. Keep FormKit as the form/validation layer, using the
  shared daisyUI section classes and existing field components.
- Before implementing frontend features or changes, read and apply the installed
  `tailwind-sane` skill's `SKILL.md`. This is a mandatory implementation discipline,
  not an npm command or a replacement for static checks. Locate it through the current
  agent's available skills rather than assuming a machine-specific path. If the skill
  is unavailable, report the missing requirement before proceeding with dependent edits.
- Inspect nearby components and `src/style.css` first. Prefer existing semantic
  components/variants, then project theme tokens and named utilities, then standard
  Tailwind utilities and configured breakpoints. Use arbitrary values only for an
  important requirement that these options cannot express cleanly.
- Keep spacing, typography, dimensions, and responsive variants on the project's scale:
  prefer `p-4`, `gap-4`, `text-xs`, and `max-w-xl` over incidental arbitrary measurements.
  Use mobile-first configured breakpoints and simple flex/grid layouts. Preserve precise
  values required by an explicit specification or to prevent a layout defect.
- Keep class lists easy to scan. Remove redundant/conflicting utilities and simplify
  layout before extracting abstractions. Roughly ten utilities on an ordinary element
  is a review signal, not a hard limit. Extract a component or variant for a coherent
  repeated pattern; do not hide a confusing utility pile behind `@apply` or a variable.
- Keep justified arbitrary values local and simple; promote repeated concepts to
  semantic tokens or reusable component APIs within the authorized task scope. Avoid
  unrelated class reordering or restyling. Do not claim visual equivalence without
  verification. The skill's browser-inspection guidance remains subject to the operator's
  policy: browser testing occurs only when explicitly demanded.

### Components and state ownership

- Use Vue 3 `<script setup lang="ts">`, typed props/events, and Composition API.
  Component and page filenames use PascalCase; composables use `use...` names.
  Prefer typed interfaces and `unknown` with narrowing over `any` and unsafe casts.
- Routed screens belong in `src/pages/`; reusable feature/UI/layout controls belong in
  `src/components/`. Reuse `AuthLayout`, `AppShell`, `UserForm`, `PasswordField`,
  `AppModal`, `LoadingState`, `PageHeader`, and `FeedbackAlert` where appropriate.
- Keep relative `/api` requests and shared transport behavior in `src/api/client.ts`,
  typed endpoint functions in `src/api/`, server reads/writes in `src/queries/`, client
  session/feedback in Pinia `src/stores/`, and reusable lifecycle logic in composables.
- Use TanStack Vue Query for server state. Include changing resource identities in
  query keys, pass cancellation signals, gate queries by access, and update/invalidate
  caches after mutations. Preserve cached data during background refreshes.
- Do not automatically retry writes. Retain the existing limited read retry policy and
  avoid retrying client errors. Surface pending, empty, error, and recovery states.
- Clear session and query data on logout or expiration. Keep tab-scoped session storage
  with its in-memory fallback; do not introduce another token persistence scheme casually.
- Keep lazy routes and access metadata in `src/router/index.ts`. Preserve nested user
  modals, direct links, Back navigation, and directory state while opening/closing them.

### Forms and validation parity

- Use named FormKit inputs inside `<FormKit type="form">` and shared configuration
  in `src/forms/formkit.ts`. Reuse validators in `src/validation/account.ts`; keep them
  aligned with `backend/internal/user/validation.go`. Backend validation remains required.
- Names are trimmed, NFC-normalized, nonempty plain text up to **255 UTF-8 bytes**,
  without angle brackets or control characters. Emails are trimmed, retain case, and
  follow the shared practical ASCII syntax and total/local/domain-label length limits.
- Passwords retain their exact bytes. New passwords require **8–72 UTF-8 bytes**;
  login permits existing passwords starting at **1 byte**. Reject invalid Unicode,
  all-whitespace values, and NUL. Blank edit passwords are omitted to preserve the hash.
- Seed drafts once per mounted record so background reads cannot overwrite user input.
  Keep field errors and form-level API errors visible. Focus the first invalid field.
- Disable repeat submissions and modal dismissal during writes. On HTTP 429, preserve
  the draft and honor the bounded retry delay through the shared cooldown behavior.
  Registration includes password confirmation; password visibility controls stay accessible.

### Visual consistency and accessibility

- Use Tailwind scale utilities and daisyUI components with semantic colors such as
  `base-*`, `primary`, and `error`. Keep shared theme tokens, typography, and component
  styles in `src/style.css`; prefer existing patterns over isolated hard-coded styles.
- Keep bundled DM Sans/Newsreader fonts and their licenses. Import individual Phosphor
  Vue icons, inheriting text color and regular weight; use 20px for controls and 16px
  for compact controls. Hide decorative icons with `aria-hidden="true"`.
- Use semantic elements, visible action labels, field labels, accessible icon-only
  names, focus indicators, and status/error announcements. Preserve keyboard behavior
  in `CustomSelect` and native dialogs, including focus containment and opener restoration.
- Preserve responsive desktop/mobile navigation, usable small-screen forms/modals,
  long-content handling, and reduced-motion preferences through code review. Only when
  explicitly demanded by the operator, test changed UI in a browser at the requested
  sizes and with keyboard navigation; report what was actually checked.
- Render user content through Vue's escaped text interpolation. Do not introduce raw
  HTML rendering for account data or external assets that conflict with the existing CSP.

## Runtime changes and completion

Preserve the production boundary: browser → Cloudflare Tunnel or VPS-wide Nginx
(external HTTPS) → application Nginx (internal HTTP) → static frontend or `/api` backend
→ PostgreSQL. Base Compose publishes only the reverse proxy; native database access
belongs in the development override. Keep separate migration execution, health-based
startup dependencies, non-root application runtimes, capability restrictions, and
configured resource limits. Vite remains a development server.

Keep application publication on loopback and outer-edge peer trust explicit and
narrow. Preserve relative redirects, normalized client IP/scheme forwarding, one
effective common header set, route-owned CSP/cache policies, and backend trust only
for application Nginx. Public HTTPS redirects and HSTS belong to the external edge;
subdomain coverage and preload require a deployment/domain-owner decision. Run
`python3 docker/nginx/security_test.py` after building the proxy image when changing
this boundary, alongside the runtime checks below.

For changes to Docker, Compose, Nginx, or migrations, supplement `make check` with the
relevant configuration validation, image builds, migration checks, and runtime smoke
checks. The local gate does not perform these. Preserve CSP/security headers and API
no-store behavior; do not weaken them to accommodate a new UI shortcut.

Before reporting work complete, review the final diff for scope, whitespace, accidental
generated output, secrets, duplicated logic, and outdated documentation. Record the
passing gate for the task scope and any additional checks. Browser checks require an
explicit operator request. If validation is blocked, state that limit clearly rather
than declaring the codebase verified.
