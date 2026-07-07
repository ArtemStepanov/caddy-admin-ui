# Implementation Plan: Safe Import Preview

**Branch**: `docs/adopt-constitution-v1` | **Date**: 2026-07-03 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-safe-import-preview/spec.md`

## Summary

Replace the blind Caddy import flow with a trustworthy review-and-commit flow: preview every discovered HTTP route, classify each route as editable or read-only preserved, explain unsupported reasons, commit imports transactionally, and make preserved routes truly non-mutating while still visible and sync-preserved.

## Technical Context

**Language/Version**: Go 1.25.6 backend; TypeScript 5.3 frontend on Node 25.3.0

**Primary Dependencies**: Existing Gin HTTP server, SQLite driver, Preact, Vite, Vitest; no new runtime dependencies planned

**Storage**: Existing SQLite route storage, extended only as needed for import classification/reason metadata and transactional replacement

**Testing**: Go `go test ./internal/...`; frontend Vitest, TypeScript check, ESLint via existing npm scripts and `make test`

**Target Platform**: Linux/self-hosted web application managing a local or containerized Caddy Admin API

**Project Type**: Web admin UI with Go API backend and Preact frontend

**Performance Goals**: Import preview for typical homelab configs, up to about 100 HTTP routes, should be understandable and returned within 2 seconds when Caddy is reachable

**Constraints**: No silent route loss; no mutation of read-only preserved routes; no raw JSON editor; no per-route import selection in the first slice; no new service/database/framework

**Scale/Scope**: Single-admin/self-hosted deployments; first slice covers Caddy HTTP routes and visible preservation of unsupported HTTP routes

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Preserve Caddy Configuration**: PASS — preview and import classify every discovered HTTP route; read-only preserved routes retain original route JSON and sync unchanged.
- **Secure the Admin Surface**: PASS — existing whole-app Basic Auth still protects API/UI when configured; import data is validated before persistence; raw details are read-only and errors avoid secret dumps.
- **Test Sync-Critical Paths**: PASS — plan requires parser, storage transaction, API mutation guard, builder preservation, and frontend route-action tests.
- **Keep the Stack Small**: PASS — uses existing Go/Gin, SQLite, Preact, and current test tooling; no new dependencies.
- **Observable and Recoverable Operations**: PASS — preview explains planned changes, import failures leave prior routes unchanged, sync/import copy tells users how to recover and when to re-run preview.

## Project Structure

### Documentation (this feature)

```text
specs/001-safe-import-preview/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── import-flow.md
└── tasks.md              # created later by /speckit.tasks
```

### Source Code (repository root)

```text
cmd/server/              # Existing server/auth wiring remains unchanged
internal/api/            # Import preview/commit contracts, read-only mutation guards
internal/caddy/          # Existing Caddy Admin API client reused
internal/config/         # Route classification, read-only reasons, raw preservation on export
internal/storage/        # Route metadata and transactional route replacement
web/src/components/      # Reusable route/import UI pieces if needed
web/src/pages/           # Settings import review flow and Dashboard read-only actions
web/src/lib/             # API client types for import preview/details
```

**Structure Decision**: Keep the feature inside existing backend, storage, config, and frontend page boundaries. Add no new top-level packages or frontend framework structure.

## Phase 0: Research Summary

See [research.md](./research.md). All technical unknowns are resolved with existing-stack decisions.

## Phase 1: Design Summary

See [data-model.md](./data-model.md), [contracts/import-flow.md](./contracts/import-flow.md), and [quickstart.md](./quickstart.md).

**Agent Context Update**: No update-agent-context script is installed in this repo. Existing `AGENTS.md` and `CLAUDE.md` already document the stack and commands used by this plan.

## Post-Design Constitution Check

- **Preserve Caddy Configuration**: PASS — data model stores read-only reason/classification and raw route details; builder contract preserves read-only routes unchanged.
- **Secure the Admin Surface**: PASS — contracts keep import/details under the existing admin API surface with Basic Auth when configured and forbid editing raw JSON.
- **Test Sync-Critical Paths**: PASS — quickstart and planned tasks cover parser classification, transactional import, mutation rejection, raw preservation, and UI action hiding.
- **Keep the Stack Small**: PASS — design uses only existing dependencies and tables with a small storage extension.
- **Observable and Recoverable Operations**: PASS — contract includes preview summaries, failure responses with unchanged local state, and manual-drift warning copy.

## Complexity Tracking

No constitution violations.
