# Caddy Admin UI

A small, self-hosted web UI for managing one Caddy HTTP server's routes without replacing the rest of your Caddy configuration.

> Alpha software. Back up production configuration and review the [current limitations](#current-limitations) before adopting an existing server.

## What it does

- Creates reverse-proxy, file-server, and redirect routes.
- Configures response and upstream request headers, path-prefix stripping, load balancing, and gzip/zstd compression.
- Discovers Caddy HTTP servers before making changes and requires explicit ownership confirmation.
- Preserves external or unsupported routes as exact, read-only JSON.
- Updates only `/config/apps/http/servers/<server>/routes`.
- Guards every write with Caddy `ETag` and `If-Match`; external drift stops the write.
- Saves a SQLite restore point before every Caddy write, with restore and JSON export in Settings.

The UI never uses Caddy's full-config `/load` endpoint. After explicit first-run confirmation on a completely empty Caddy instance, it creates only the initial `apps` object through an ETag-guarded `/config/apps` request.

## Quick start

The included Compose stack publishes the UI only on loopback and keeps Caddy's Admin API inside the Docker network:

```bash
git clone https://github.com/ArtemStepanov/caddy-admin-ui.git
cd caddy-admin-ui
docker compose up -d
```

Open <http://127.0.0.1:3000>, select **Settings**, review the discovered Caddy server, and confirm ownership.

To build the image from the checked-out source instead:

```bash
make docker-up-build
```

## Review a pull request privately

GitHub Codespaces can run any selected pull request as a disposable Admin UI +
Caddy stack without creating a deployment, check, or comment on the pull
request:

1. Open the pull request on GitHub.
2. Select **Code → Codespaces → Create codespace**.
3. Wait for the dev container to build the preview stack.
4. Open **Caddy Admin UI** from the forwarded ports panel. The **Caddy data
   plane** port is available separately for end-to-end route checks.

Both forwarded ports are declared private and require authentication as the
Codespace owner. Caddy's Admin API is reachable only inside the preview Docker
network and is never forwarded. The preview stores data only in disposable
Docker volumes; deleting the Codespace removes the whole environment.

To review a pull request locally without switching the current checkout, use
the worktree helper:

```bash
./scripts/preview-pr 127
# Admin UI: http://127.0.0.1:3000
# Caddy:    http://127.0.0.1:8080/preview-seed

./scripts/preview-pr down 127
```

The helper fetches `refs/pull/<number>/head`, builds it in an isolated temporary
git worktree, and binds both ports to loopback. Only one preview can use the
default ports at a time; set `PREVIEW_UI_PORT` and `PREVIEW_CADDY_PORT` to run
more than one. Review untrusted pull-request code before running it locally.

For the current checkout, the equivalent commands are `make preview` and
`make preview-down`.

## Connect an existing Caddy

Caddy Admin UI must reach the Admin API, but that API should not be exposed to an untrusted network. Set credentials for the UI and an Admin API URL reachable from its container:

```bash
export CADDY_ADMIN_URL=http://caddy.internal:2019
export ADMIN_USER=admin
export ADMIN_PASSWORD='replace-with-a-long-random-password'
docker compose -f docker-compose.prod.yml up -d
```

The production Compose file binds the UI to `127.0.0.1:3000` by default. Set `UI_BIND_ADDRESS` only when you also have an appropriate network boundary and TLS termination.

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `CADDY_ADMIN_URL` | `http://localhost:2019` | Initial Caddy Admin API URL |
| `DB_PATH` | `/app/data/routes.db` in the image | SQLite database path |
| `LISTEN_ADDR` | `127.0.0.1:3000` | HTTP listen address |
| `WEB_DIR` | `./web/dist` | Built frontend directory |
| `GIN_MODE` | `release` in the image | Gin runtime mode |
| `ADMIN_USER` | empty | UI/API HTTP Basic Auth user |
| `ADMIN_PASSWORD` | empty | UI/API HTTP Basic Auth password |
| `ALLOW_INSECURE_NO_AUTH` | `false` | Explicit escape hatch for an unauthenticated non-loopback listener |

Both auth variables must be set together. The server refuses to start without authentication on a non-loopback address unless `ALLOW_INSECURE_NO_AUTH=true` is explicitly set.

## Safety model

1. **Review:** setup fetches the full Caddy inventory only to list HTTP servers and preview their routes.
2. **Confirm:** confirmation re-fetches Caddy and verifies a content-derived preview token. A changed config is rejected.
3. **Own one scope:** after setup, the UI writes only the selected server's `routes` array. Listeners, TLS automation, logging, storage, and other apps/servers are untouched.
4. **Preserve:** routes without a `caddy-admin-ui-route-*` ID are stored as exact raw JSON and remain read-only.
5. **Detect drift:** a mismatched ETag returns a conflict before local or remote mutation.
6. **Recover:** a snapshot of the live route array is stored before every write and can be exported or restored.

See the [Caddy Admin API documentation](https://caddyserver.com/docs/api) for scoped configuration and ETag semantics.

## Development

Go 1.25.13 and Node.js 24 are declared in `mise.toml`.

```bash
make deps
make test
make build
make frontend
```

For local development:

```bash
make dev
make dev-frontend
```

### Architecture

```text
Preact/TypeScript UI ── /api/* ──> Go/Gin API ── scoped HTTP + ETag ──> Caddy
                                      |
                                    SQLite
```

CI runs race-enabled Go tests, frontend lint/typecheck/coverage/build, a real-Caddy integration test, container builds, dependency audits, gosec, Trivy, and scheduled CodeQL.

## API highlights

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/setup/preview` | Inventory and preview a selected Caddy HTTP server |
| `POST` | `/api/setup/confirm` | Confirm the unchanged preview and adopt its routes |
| `GET/POST/PUT/DELETE` | `/api/routes...` | Manage UI-owned routes |
| `POST` | `/api/sync` | Guarded manual sync of the selected route array |
| `GET` | `/readyz` | Minimal database + Caddy readiness probe |
| `GET` | `/api/snapshots` | List recent restore points |
| `GET` | `/api/snapshots/:id/export` | Export route-array JSON |
| `POST` | `/api/snapshots/:id/restore` | Guarded snapshot restore |

Mutating requests require the `X-Caddy-Admin-UI: 1` header. The bundled web client sends it automatically.

## Current limitations

- One Caddy instance and one HTTP server are managed per UI database.
- External routes are deliberately read-only even if their shape looks editable; this prevents lossy round trips.
- Route-level authentication, certificate inspection, access-log viewing, upstream health, route reordering, and multi-instance management are not implemented yet.
- The SQLite database is local to one UI process; active-active replicas are unsupported.

The intended path to a stable 1.0 is tracked in [ROADMAP.md](ROADMAP.md).

## Contributing and security

See [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request. Please report vulnerabilities using the process in [SECURITY.md](SECURITY.md), not a public issue.

## License

[MIT](LICENSE)
