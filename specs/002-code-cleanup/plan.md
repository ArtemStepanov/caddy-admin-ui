# Implementation Plan: Code Cleanup

**Branch**: `002-code-cleanup` | **Date**: 2026-07-10 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-code-cleanup/spec.md`

## Summary

Narrow cleanup to two existing seams: make `config.ValidateRoutesForBuild` the only post-parse import-safety gate before preview or transactional replacement, and ensure the existing three-kind route form drops proxy-only path stripping after a kind switch while retaining only route-level settings that apply to every kind. Preserve current API shapes, read-only raw-route handling, warning copy, and sync behavior.

## Technical Context

**Language/Version**: Go 1.25.6 backend; TypeScript 5.3 with Preact on Node 25.3.0 frontend

**Primary Dependencies**: Existing Gin, SQLite (`mattn/go-sqlite3`), Preact, Vite, Vitest, and Testing Library; no new dependencies

**Storage**: Existing SQLite `routes` table and `storage.Route` preservation metadata; no schema changes

**Testing**: Go `go test` for config/API paths; Vitest, TypeScript, and ESLint for the route form; `make test` for the full suite

**Target Platform**: Linux/self-hosted web admin UI managing the Caddy Admin API

**Project Type**: Web application with Go API backend and Preact/TypeScript frontend

**Performance Goals**: No material performance change; validation remains in-process over the fetched route list before preview or replacement

**Constraints**: No silent configuration loss; no read-only mutation; no endpoint/database schema changes; no route kinds, warning-copy, dashboard, authentication, or dependency changes

**Scale/Scope**: Single-admin/self-hosted deployments; only imported-route validation and three existing route-form kinds

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Preserve Caddy Configuration**: PASS — preserved routes retain existing raw JSON, status, and reason; invalid imports fail before transactional replacement.
- **Secure the Admin Surface**: PASS — the plan changes neither bind/auth behavior nor client trust boundaries; it retains existing server-side validation and secret-safe errors.
- **Test Sync-Critical Paths**: PASS — focused Go tests cover import acceptance, preservation, rejection, and unchanged local storage; current sync validation remains the shared gate. Frontend tests cover emitted route data.
- **Keep the Stack Small**: PASS — delete the redundant API validator and use existing tooling; no new service, dependency, database, or abstraction.
- **Observable and Recoverable Operations**: PASS — existing import/sync error responses, warnings, recovery guidance, and route details remain unchanged.

## Project Structure

### Documentation (this feature)

```text
specs/002-code-cleanup/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── import-and-route-editor.md
└── tasks.md             # created later by /speckit.tasks
```

### Source Code (repository root)

```text
internal/api/handlers.go        # remove duplicate import-only validation call
internal/api/handlers_test.go   # import endpoint acceptance/preservation/rejection coverage
internal/config/builder.go      # shared validation gate used by import and sync
internal/config/parser.go       # existing support classification retained
internal/config/parser_test.go  # parser classification fixtures retained/extended as needed
internal/storage/models.go      # existing read-only predicate and route metadata
web/src/pages/RouteForm.tsx     # selected-kind submission and form state
web/src/pages/RouteForm.test.tsx # focused kind-switch request coverage (new)
web/src/lib/api.ts              # existing request shape retained
```

**Structure Decision**: Keep the change in the existing config/API seam and route form. Do not add packages, database structures, endpoint variants, or frontend state abstractions.

## Phase 0: Research Summary

See [research.md](./research.md). All technical choices are resolved: the existing shared config validator replaces the overlapping API helper; read-only storage remains unchanged; response headers remain route-level while proxy-only path stripping is dropped on incompatible form submissions.

## Phase 1: Design Summary

See [data-model.md](./data-model.md), [contracts/import-and-route-editor.md](./contracts/import-and-route-editor.md), and [quickstart.md](./quickstart.md).

**Agent Context Update**: No update-agent-context script is installed in this checkout. `AGENTS.md` and `CLAUDE.md` already describe the stack, layout, and commands used by this plan; no context change is needed.

## Post-Design Constitution Check

- **Preserve Caddy Configuration**: PASS — the contract requires raw preserved content and a reason, and failure occurs before `ReplaceAllRoutes`; legacy/raw routes remain read-only through existing storage and builder behavior.
- **Secure the Admin Surface**: PASS — API routes, validation placement, auth assumptions, and error-copy handling are unchanged; frontend isolation does not weaken the backend boundary.
- **Test Sync-Critical Paths**: PASS — the validation used before import remains the same one used before sync, with automated coverage for editable, partial, unsupported, legacy, and invalid inputs.
- **Keep the Stack Small**: PASS — the design removes duplication and adds only focused tests using installed tooling.
- **Observable and Recoverable Operations**: PASS — current sync warning and import recovery contract remains intact, including when invalid input prevents local replacement.

## Complexity Tracking

No constitution violations.
