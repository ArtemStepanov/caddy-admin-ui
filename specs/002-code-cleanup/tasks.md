# Tasks: Code Cleanup

**Input**: Design documents from `/specs/002-code-cleanup/`

**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contract](./contracts/import-and-route-editor.md), [quickstart.md](./quickstart.md)

**Tests**: Required. This feature changes import validation and must preserve raw Caddy configuration; it also explicitly requires automated coverage for every supported route kind and import classification.

**Organization**: Tasks are grouped by user story so each can be implemented and validated independently.

## Phase 1: Setup

**Purpose**: Confirm the existing test baseline before the bounded cleanup.

- [X] T001 Run the focused backend and frontend baseline commands documented in `specs/002-code-cleanup/quickstart.md`

---

## Phase 2: Foundational

**Purpose**: No foundational implementation is needed. The existing `storage.Route` metadata, parser classification, API routes, and shared config validator are the required foundation.

**Checkpoint**: User-story work can begin after T001.

---

## Phase 3: User Story 1 - Safe import validation (Priority: P1) 🎯 MVP

**Goal**: Make the existing shared config validator the only post-parse safety gate while retaining editable, preserved, legacy, and rejected route behavior.

**Independent Test**: Import complete editable, partial, unsupported, legacy-preserved, and invalid routes; accepted/preserved results retain their status and raw content, while invalid input leaves local routes unchanged.

### Tests for User Story 1

- [X] T002 [P] [US1] Extend import endpoint coverage for editable, partial, unsupported, legacy-preserved, and invalid route decisions (including unchanged local routes after rejection) in `internal/api/handlers_test.go`

### Implementation for User Story 1

- [X] T003 [US1] Remove the overlapping `validateImportedRoutes` helper and retain `config.ValidateRoutesForBuild` as the single post-parse gate in `internal/api/handlers.go`

**Checkpoint**: Preview and import apply one server-side decision before local replacement, and preserved routes remain non-mutating and sync-safe.

---

## Phase 4: User Story 2 - Mode-safe route editing (Priority: P2)

**Goal**: Keep submissions isolated to the selected existing route kind without clearing valid route-level response headers.

**Independent Test**: Switch from a populated reverse proxy to file server and redirect, submit each, and verify the request has only the selected `config` shape and no proxy-only `strip_path_prefix`; verify a read-only route stays outside editable controls.

### Tests for User Story 2

- [X] T004 [P] [US2] Create route-form submission and read-only rendering coverage for all three handler kinds in `web/src/pages/RouteForm.test.tsx`

### Implementation for User Story 2

- [X] T005 [US2] Omit proxy-only `strip_path_prefix` from file-server and redirect submissions while preserving selected-kind config resets and route-level response headers in `web/src/pages/RouteForm.tsx`

**Checkpoint**: Supported kind switches cannot persist incompatible proxy-only configuration, and read-only routes remain outside the editor.

---

## Phase 5: Polish & Cross-Cutting Validation

**Purpose**: Run the documented end-to-end checks without expanding scope.

- [X] T006 Run all focused and full validation commands in `specs/002-code-cleanup/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies.
- **Foundational (Phase 2)**: No implementation task; existing infrastructure is reused.
- **US1 (Phase 3)** and **US2 (Phase 4)**: Start after T001; implement in P1 → P2 order for the MVP, or in parallel when separately staffed.
- **Polish (Phase 5)**: Depends on the desired user stories being complete.

### User Story Dependencies

- **US1 (P1)**: Independent; it changes only the Go import-validation seam.
- **US2 (P2)**: Independent; it changes only the Preact route form.

### Parallel Opportunities

- T002 and T004 touch different test files and can be prepared in parallel after T001.
- US1 and US2 have no source-file dependency and can be implemented by separate developers after their respective tests are in place.

## Parallel Example: User Story 1

```text
Task: "Extend import endpoint coverage in internal/api/handlers_test.go"
```

## Parallel Example: User Story 2

```text
Task: "Create route-form coverage in web/src/pages/RouteForm.test.tsx"
```

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete T001.
2. Add T002, then remove the duplicate gate in T003.
3. Run the focused Go tests and validate import rejection leaves stored routes untouched.

### Incremental Delivery

1. Deliver US1 as the import-safety MVP.
2. Add T004 and T005 for route-kind isolation.
3. Complete T006 before merge.

## Format Validation

All six tasks use the required checkbox, sequential task ID, optional `[P]` marker only for parallel work, required story label for story tasks, and an exact file path.
