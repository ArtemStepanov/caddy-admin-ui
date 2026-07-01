## Project Overview

Caddy Admin UI is a web UI for managing Caddy server routes without manually editing the Caddyfile. It provides CRUD operations for reverse proxies, file servers, and redirects, with real-time sync to Caddy's Admin API.

## Build & Development Commands

### Prerequisites

Go and GCC are installed via Homebrew. Node.js 25.3.0 is managed by `mise.toml`.

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
| `GIN_MODE` | `debug` | Gin framework mode |
| `WEB_DIR` | `./web/dist` | Frontend static files directory |
| `ADMIN_USER` | _(empty, auth disabled)_ | Basic Auth username |
| `ADMIN_PASSWORD` | _(empty, auth disabled)_ | Basic Auth password |

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

- **`api/`** — HTTP handlers and route definitions. All mutation endpoints (create/update/delete/toggle) auto-sync to Caddy after persisting to SQLite. Sync failures return warnings but don't fail the request.
- **`config/builder.go`** — Converts internal Route models into Caddy JSON config. Separate builder functions per handler type (reverse_proxy, file_server, redir). Injects encode middleware globally when enabled.
- **`config/parser.go`** — Parses Caddy JSON config back into internal Route models for the import feature.
- **`caddy/client.go`** — HTTP client for Caddy's Admin API (health, get/load/set config).
- **`storage/`** — SQLite persistence layer with models. Routes store handler-specific config as `json.RawMessage` blobs. `RawCaddyRoute` field preserves unknown handlers during round-trip sync to prevent data loss.

### Frontend (`web/src/`)

- **`pages/`** — Dashboard (route list), RouteForm (create/edit), Settings (global config), NotFound (404)
- **`lib/`** — `api.ts` (TypeScript API client), `syncNotify.ts` (toast notification pub/sub)
- **`components/`** — Layout, StatusBadge, Toast, `forms/HeaderEditor`

### Key Design Decisions

1. **Raw route preservation**: Unknown Caddy handler types are stored in `RawCaddyRoute` and merged back during export, preventing data loss on round-trip.
2. **Dynamic Caddy URL**: Configurable at runtime via GlobalConfig (stored in DB), falling back to `CADDY_ADMIN_URL` env var.
3. **Preact over React**: Chosen for minimal bundle size with the same component model.
4. **Tailwind CSS**: Utility-first CSS framework used for all frontend styling.