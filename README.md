# Caddy Admin UI

A lightweight web UI for managing [Caddy](https://caddyserver.com/) routes without editing the Caddyfile. Create, update, and delete reverse proxies, file servers, and redirects — synced to Caddy's Admin API in real time.

> **Alpha Software** — This project is under active development. Expect breaking changes, incomplete features, and rough edges. Use with caution in production environments.

## Features

- **Reverse Proxy** — forward requests to backend services
- **File Server** — serve static files from directories
- **Redirects** — set up URL redirects
- **Headers** — add security and CORS headers
- **Basic Auth** — password-protect routes
- **Compression** — enable gzip/zstd encoding
- **Import Review** — preview Caddy routes before replacing local state
- **Read-only Preservation** — keep unsupported Caddy routes visible and sync-preserved without editing them

## Quick Start

### Docker Compose

```bash
git clone https://github.com/ArtemStepanov/caddy-admin-ui.git
cd caddy-admin-ui
docker compose up -d --build
# UI at http://localhost:3000
```

This starts the orchestrator alongside a test Caddy instance with the admin API enabled.

### Connect to an Existing Caddy Instance

Use the production compose file to connect to a Caddy server running elsewhere:

```bash
CADDY_ADMIN_URL=http://your-caddy:2019 docker compose -f docker-compose.prod.yml up -d

```

Or run the container directly:

```bash
docker run -d \
  -p 3000:3000 \
  -e CADDY_ADMIN_URL=http://your-caddy:2019 \
  -v caddy-admin-ui-data:/app/data \
  ghcr.io/artemstepanov/caddy-admin-ui
```

### Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `CADDY_ADMIN_URL` | `http://localhost:2019` | Caddy Admin API URL |
| `DB_PATH` | `/app/data/routes.db` | SQLite database path |
| `LISTEN_ADDR` | `127.0.0.1:3000` | Server listen address (loopback by default) |
| `GIN_MODE` | `debug` | Gin mode (`debug` / `release`) |
| `ADMIN_USER` | _(empty)_ | Basic Auth username (optional) |
| `ADMIN_PASSWORD` | _(empty)_ | Basic Auth password (optional) |

The server binds to localhost by default. In Docker, it binds to `0.0.0.0:3000` so port mapping works.

When `ADMIN_USER` and `ADMIN_PASSWORD` are both set, the entire app (UI + API) is protected with HTTP Basic Auth. The browser shows a native login dialog.

The Caddy URL can also be changed at runtime from the Settings page.

## Development

Go and GCC are installed via Homebrew. Node.js is managed by `mise.toml`.

```bash
make deps            # install Go + JS dependencies
make dev             # run backend dev server
make dev-frontend    # run frontend dev server with HMR
make test            # run all tests
make build           # build backend binary (requires GCC for CGO/SQLite)
make frontend        # build frontend
```

## Architecture

```
Preact/TypeScript UI  ──/api/*──►  Go backend (Gin)  ──HTTP──►  Caddy Admin API
                                        │
                                    SQLite DB
```

- **Backend** (`internal/`) — Go with Gin. Routes are stored in SQLite and synced to Caddy on every mutation.
- **Frontend** (`web/`) — Preact + TypeScript. Minimal bundle, same React component model.
- **SQLite** via `mattn/go-sqlite3` (requires CGO, built with GCC).

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/routes` | List all routes |
| `POST` | `/api/routes` | Create a route |
| `GET` | `/api/routes/:id` | Get a route |
| `PUT` | `/api/routes/:id` | Update a route |
| `DELETE` | `/api/routes/:id` | Delete a route |
| `POST` | `/api/routes/:id/toggle` | Enable/disable a route |
| `GET/PUT` | `/api/config` | Global configuration |
| `GET` | `/api/status` | Caddy connection status |
| `POST` | `/api/sync` | Sync all routes to Caddy |
| `POST` | `/api/import-preview` | Preview import summary/groups from Caddy without changing local routes |
| `POST` | `/api/import` | Transactionally replace local routes from Caddy |
| `GET` | `/api/routes/:id/details` | View preserved raw JSON for read-only routes |

## Import Safety

Import preview classifies discovered HTTP routes as editable or read-only preserved. Unsupported routes remain visible, are emitted back to Caddy unchanged on sync, and cannot be edited, deleted, or toggled from the UI/API. Manual Caddy changes are not auto-merged; re-run import review before syncing after manual edits.

## License

MIT
