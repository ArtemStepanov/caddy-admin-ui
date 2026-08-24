## Project Overview

Caddy Admin UI is a web UI for managing Caddy server routes without manually editing the Caddyfile. It provides CRUD operations for reverse proxies, file servers, and redirects, with real-time sync to Caddy's Admin API.

## Build & Development Commands

### Prerequisites

Go 1.25.13 and Node.js 24 are managed by `mise.toml`. GCC is required for CGO/SQLite builds.

### Common Commands

| Command | Purpose |
|---------|---------|
| `make build` | Build backend (requires GCC for CGO/SQLite) |
| `make frontend` | Install deps and build frontend |
| `make dev` | Run backend dev server (`go run ./cmd/server`) |
| `make dev-frontend` | Run frontend dev server with HMR (`cd web && npm run dev`) |
| `make test` | Run all tests (Go + frontend coverage) |
| `make docker-up` / `make docker-down` | Start/stop Docker Compose stack |
| `make deps` | Install all dependencies (Go + npm) |
| `make clean` | Remove build artifacts (`bin/`, `web/dist/`, `web/node_modules/`) |
| `make logs` | Tail Docker Compose logs |
| `make docker-up-build` | Build image and start Docker Compose stack |
| `make preview` / `make preview-down` | Build/start or remove the disposable loopback preview stack |
| `./scripts/preview-pr <number>` | Preview a PR in an isolated temporary git worktree |

### Running Individual Tests

**Go (backend):**
```sh
CGO_ENABLED=1 go test -v -race ./internal/api/...
CGO_ENABLED=1 go test -v -race ./internal/config/...
CGO_ENABLED=1 go test -v -race ./internal/storage/...
```

**Frontend:**
```sh
cd web && npx vitest run              # all tests
cd web && npx vitest run StatusBadge  # single test file by name
cd web && npx vitest run --coverage   # with coverage
cd web && npm run lint                # ESLint
cd web && npx tsc --noEmit            # type-check only
```

### Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `CADDY_ADMIN_URL` | `http://localhost:2019` | Caddy Admin API endpoint |
| `DB_PATH` | `/app/data/routes.db` | SQLite database path |
| `LISTEN_ADDR` | `127.0.0.1:3000` | Server listen address (loopback by default) |
| `GIN_MODE` | `release` | Gin framework mode |
| `WEB_DIR` | `./web/dist` | Frontend static files directory |
| `ADMIN_USER` | _(empty, auth disabled)_ | Basic Auth username |
| `ADMIN_PASSWORD` | _(empty, auth disabled)_ | Basic Auth password |
| `ALLOW_INSECURE_NO_AUTH` | `false` | Explicitly allow an unauthenticated non-loopback listener |

## Commit Conventions

Releases are automated with release-please: version bumps and the changelog are computed from commit messages, so every commit that lands on `main` MUST follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: ...` — new feature (minor version bump)
- `fix: ...` — bug fix (patch version bump)
- `feat!: ...` or a `BREAKING CHANGE:` footer — breaking change (minor bump while pre-1.0)
- `chore: ...`, `ci: ...`, `docs: ...`, `refactor: ...`, `test: ...` — no version bump, hidden from the changelog

PRs are squash-merged, so the **PR title becomes the commit message on main** — it must follow the convention too. Scopes are optional (`feat(api): ...`). Releases ship when the auto-generated release-please PR is merged.

## Architecture

```
Frontend (Preact/TS)  ──/api/*──►  Go Backend (Gin)  ──HTTP──►  Caddy Admin API
                                       │
                                   SQLite DB
```

### Backend (`internal/`)

- **`api/`** — HTTP handlers, managed-server setup, serialized mutations, ETag drift checks, and snapshot restore/export.
- **`config/builder.go`** — Builds only the managed Caddy route array. UI-owned routes receive stable `@id` markers; read-only routes are emitted unchanged.
- **`config/parser.go`** — Parses a selected server's raw route array without narrowing unknown JSON. External routes remain read-only.
- **`caddy/client.go`** — HTTP client for scoped Caddy config reads and ETag-guarded PATCH/PUT writes, including empty-Caddy bootstrap.
- **`storage/`** — SQLite in WAL mode with route order, managed-server state, and immutable pre-write snapshots.

### Frontend (`web/src/`)

- **`pages/`** — Dashboard (route list), RouteForm (create/edit), Settings (global config), NotFound (404)
- **`lib/`** — `api.ts` (TypeScript API client), `syncNotify.ts` (toast notification pub/sub)
- **`components/`** — Layout, StatusBadge, Toast, `forms/HeaderEditor`

### Key Design Decisions

1. **Scoped ownership**: The UI adopts one HTTP server after preview/confirmation and writes only its `routes` array. It never calls full-config `/load`.
2. **Optimistic concurrency**: Every write uses the current route-array ETag with `If-Match`. Drift aborts before local persistence.
3. **Raw route preservation**: Routes without a `caddy-admin-ui-route-*` ID are stored as exact raw JSON and remain read-only.
4. **Recovery**: A snapshot of the live route array is persisted before every write; restore is also ETag-guarded.
5. **Dynamic Caddy URL**: Configurable at runtime via GlobalConfig, falling back to `CADDY_ADMIN_URL`. Changing it clears setup ownership.
6. **Preact and Tailwind CSS**: Keep the frontend small while using a familiar component and utility-class model.
7. **Private manual previews**: Codespaces and the local worktree helper run the real two-container stack on demand. No PR deployment workflow or public Caddy Admin API is used.
