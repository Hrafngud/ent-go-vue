# Vue UI

Vue 3 Composition API + TypeScript + Vue Router 4 + Vite + Tailwind CSS 4 + daisyUI 5 + Pinia + TanStack Vue Query.

```bash
npm install
npm run dev
npm run build
```

The browser uses the API client in `src/api/client.ts`, with relative `/api` URLs. Local Vite development proxies those requests to `http://localhost:8080`; optionally set `API_PROXY_TARGET` in a frontend `.env` to change the native backend target.

All routed screens live under `src/pages`; reusable UI lives under `src/components`. `App.vue` composes the route outlet and signed-in shell. UI controls use daisyUI 5 and semantic theme colors throughout, with the shared dark theme declared in `src/style.css`.

```text
src/
  pages/
    LoginPage.vue, RegisterPage.vue, WorkspacePage.vue, NotFoundPage.vue
    admin/
      UsersPage.vue, UserCreatePage.vue, UserDetailPage.vue, UserEditPage.vue
  components/
    layout/       BrandMark, AuthLayout, AppShell
    forms/        PasswordField
    ui/           PageHeader, LoadingState, FeedbackAlert
    users/        UserForm, UserTable, UserIdentity, UserActions,
                  UsersBreadcrumbs, DeleteUserPanel
    ConnectionStatus.vue
  router/         routes, session restoration, member/root guards
  stores/         tab-scoped session and dismissible feedback (Pinia)
  queries/        health/readiness and users (TanStack Query)
  api/            relative HTTP client, typed user endpoints, error messages
```

The hierarchy is `App → route page → reusable components` for public screens, and `App → AppShell → route page → reusable components` for signed-in screens. Login and registration share `AuthLayout`; registration, create, and edit share `UserForm` and `PasswordField`. Deletion is confirmed inline on the details page.

| Route | Access |
| --- | --- |
| `/login`, `/register` | Signed-out users |
| `/workspace` | Signed-in users |
| `/admin/users` | Configured root account |
| `/admin/users/new` | Configured root account |
| `/admin/users/:id` | Configured root account |
| `/admin/users/:id/edit` | Configured root account |

The root account is identified by the backend's `ROOT_EMAIL`; `is_admin` in `/api/users/me` controls navigation. The backend independently authorizes every admin request. Expired tokens clear the session and query cache and return to login. Registration handles an empty 201 response and returns to login with the email prefilled. Admin forms validate password byte limits, preserve existing passwords when left blank, and show duplicate-email errors. Root email editing and deletion are disabled and rejected by the API.

`src/queries/health.ts` owns the server health/readiness state. `ConnectionStatus.vue` displays API and database connectivity, loading/errors, and a refresh action. `src/queries/users.ts` owns user reads; successful mutations update or invalidate cached records. Search, sorting, and pagination operate on the backend's complete user list.

The production `frontend` service builds static assets and serves them with an unprivileged Nginx runtime using this directory's `nginx.conf`. Its Dockerfile's build context is `frontend/vue-ui`. The separate root Compose `nginx` service proxies `/` to this container and `/api/` to the backend, publishing one application port. Vite is used only for development.

```bash
# Run from repository root:
docker compose build frontend nginx
```

See the [root README](../../README.md) for the full stack, migrations, and native development setup.
