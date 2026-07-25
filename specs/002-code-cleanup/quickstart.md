# Quickstart: Code Cleanup Validation

## Prerequisites

- Go 1.25.6 and GCC for CGO SQLite tests.
- Node 25.3.0 through `mise`.
- Existing dependencies installed (`make deps` if needed).

## Automated validation

```sh
CGO_ENABLED=1 go test -v -race ./internal/config/... ./internal/api/...
cd web && npx vitest run && npx tsc --noEmit && npm run lint
make test
```

## Scenario 1: One import safety decision

1. Run focused Go tests for a complete editable route, a partially supported route, an unsupported route, and a legacy raw route.
2. Confirm the first is accepted as editable and every preserved case has its original raw JSON, read-only status, and reason.
3. Exercise an incomplete editable configuration through both import preview and import.
4. Expected: both reject it before `ReplaceAllRoutes`; the previous local route list is unchanged.

See [data-model.md](./data-model.md) and [the import contract](./contracts/import-and-route-editor.md).

## Scenario 2: Route-kind isolation

1. In a Vitest route-form test, enter a reverse-proxy-only value, select `file_server`, then submit; repeat for `redir`.
2. Confirm each submitted `config` has only the selected kind's defaults/values, and `strip_path_prefix` is absent outside `reverse_proxy`.
3. Confirm response headers survive the change because they are route-level settings.
4. Open a preserved route and confirm the existing read-only view still prevents editable configuration handling.

## Expected outcome

All focused and full checks pass. No endpoint schema, database schema, route kind, warning copy, dashboard behavior, authentication behavior, or dependency changes are introduced.
