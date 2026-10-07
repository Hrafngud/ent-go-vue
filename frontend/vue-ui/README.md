# Vue UI

Vue 3 Composition API + TypeScript + Vite + Tailwind CSS 4 + daisyUI 5 + Pinia + TanStack Vue Query.

```bash
npm install
npm run dev
npm run build
```

The browser uses the API client in `src/api/client.ts`, with relative `/api` URLs. Local Vite development proxies those requests to `http://localhost:8080`; optionally set `API_PROXY_TARGET` in a frontend `.env` to change the native backend target.

`src/queries/health.ts` owns the server health/readiness state in TanStack Query. `ConnectionStatus.vue` displays API and database connectivity, loading/errors, and a refresh action. Pinia is reserved for client/application state.

Production builds static assets into the root Compose Nginx service; Vite is used only for development. The Dockerfile's build context is the repository root so it can include `docker/nginx/nginx.conf`:

```bash
# Run from repository root:
docker compose build nginx
```

See the [root README](../../README.md) for the full stack, migrations, and native development setup.
