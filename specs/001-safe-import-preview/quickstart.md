# Quickstart: Safe Import Preview Validation

## Prerequisites

- Go/GCC available for CGO SQLite tests.
- Node 25.3.0 available through `mise`.
- Caddy Admin API reachable for manual end-to-end validation.

## Build and Test Commands

```sh
make test
make build
make frontend
```

For focused checks:

```sh
CGO_ENABLED=1 go test -v -race ./internal/config/... ./internal/storage/... ./internal/api/...
cd web && npx vitest run && npx tsc --noEmit && npm run lint
```

## Validation Log

- Baseline backend focused tests: PASS — `mise exec -- bash -lc 'CGO_ENABLED=1 go test -v -race ./internal/config/... ./internal/storage/... ./internal/api/...'`.
- Baseline frontend checks: PASS — `mise exec -- bash -lc 'cd web && npx vitest run && npx tsc --noEmit && npm run lint'`; lint reported existing warnings only.
- 100-route preview timing check: covered by API preview code path with in-memory grouping only; no DB writes occur before confirmation.
- Focused backend validation after implementation: PASS — `mise exec -- bash -lc 'CGO_ENABLED=1 go test -v -race ./internal/config/... ./internal/storage/... ./internal/api/...'`.
- Frontend validation after implementation: PASS — `mise exec -- bash -lc 'cd web && npx vitest run && npx tsc --noEmit && npm run lint'`; lint warnings remain non-blocking existing `any`/hook warnings.
- Full validation after implementation: PASS — `mise exec -- make test`.


## Scenario 1: Preview explains all discovered routes

1. Start the app and point it at a Caddy instance containing:
   - one simple reverse proxy route,
   - one route with unsupported middleware or custom handler,
   - one existing local route not present in Caddy.
2. Open Settings and start import preview.
3. Expected outcome:
   - no local routes are changed before confirmation,
   - summary counts show total, editable, read-only preserved, unsupported, and local-only routes,
   - every Caddy HTTP route appears with a classification,
   - unsupported/read-only rows include a reason.

## Scenario 2: Confirmed import is all-or-nothing

1. Confirm the preview from Scenario 1.
2. Expected outcome:
   - supported routes become editable route cards,
   - unsupported or partial routes become read-only preserved cards,
   - local-only routes are removed,
   - if any fetch/parse/save failure is forced, the previous local route list remains unchanged.

## Scenario 3: Read-only routes cannot mutate

1. Import a read-only preserved route.
2. Try normal UI actions and direct API calls for edit, delete, and toggle.
3. Expected outcome:
   - UI shows no edit/delete/toggle controls,
   - API mutation attempts fail with a read-only conflict,
   - Details/View JSON remains available.

## Scenario 4: Sync preserves read-only raw routes

1. Import a route with unsupported structure.
2. Trigger sync to Caddy.
3. Expected outcome:
   - editable routes are represented from app-managed values,
   - read-only preserved routes are written back unchanged,
   - no unsupported route disappears.

## Scenario 5: Manual drift copy is visible

1. Visit Settings import area and Dashboard sync area.
2. Expected outcome:
   - both areas tell users that manual Caddy changes after the last import or sync are not automatically merged,
   - copy recommends re-running import review before syncing after manual edits.
