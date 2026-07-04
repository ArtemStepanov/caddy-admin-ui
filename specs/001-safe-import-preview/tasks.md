# Tasks: Safe Import Preview

**Input**: Design documents from `/specs/001-safe-import-preview/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/import-flow.md, quickstart.md

**Tests**: Required because this feature touches import/export/sync, storage, API mutations, data preservation, and admin UI behavior.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm the existing project state and avoid adding unnecessary stack changes.

- [X] T001 Run current backend import/parser/storage tests and record baseline command result in specs/001-safe-import-preview/quickstart.md
- [X] T002 Run current frontend type/test checks and record baseline command result in specs/001-safe-import-preview/quickstart.md
- [X] T003 [P] Review existing import endpoints and note current response gaps in specs/001-safe-import-preview/research.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared route metadata, classification, validation, and transaction primitives needed by all stories.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

### Tests for Foundation

- [X] T004 Add storage tests for support_status, readonly_reason, and raw route persistence in internal/storage/sqlite_test.go
- [X] T005 Add storage transaction test proving route replacement rolls back on insert failure in internal/storage/sqlite_test.go
- [X] T006 [P] Add parser tests for editable, partial_readonly, unsupported_readonly, and readonly_reason classification in internal/config/parser_test.go
- [X] T007 [P] Add builder test proving a read-only raw Caddy route is validated and emitted unchanged on sync in internal/config/builder_test.go
- [X] T008 Add API tests for import validation errors and secret-safe error responses in internal/api/handlers_test.go
- [X] T009 Add API preview failure tests for Caddy fetch failure, parse failure, zero HTTP routes, and unchanged local state in internal/api/handlers_test.go

### Implementation for Foundation

- [X] T010 Add SupportStatus and ReadOnlyReason fields to Route in internal/storage/models.go
- [X] T011 Add SQLite migrations and scan/write support for support_status and readonly_reason in internal/storage/sqlite.go
- [X] T012 Implement transactional ReplaceAllRoutes(routes []*Route) in internal/storage/sqlite.go
- [X] T013 Add route classification constants/helpers in internal/config/parser.go
- [X] T014 Update ParseCaddyConfig classification and reason assignment in internal/config/parser.go
- [X] T015 Change buildRoute to validate and return preserved raw routes unchanged for read-only routes in internal/config/builder.go
- [X] T016 Update frontend Route and import preview TypeScript types in web/src/lib/api.ts
- [X] T017 Add import route validation helper and secret-safe error formatting in internal/api/handlers.go
- [X] T018 Apply import validation to PreviewImport and ImportFromCaddy before response or storage writes in internal/api/handlers.go

**Checkpoint**: Foundation ready - user story implementation can now begin.

---

## Phase 3: User Story 1 - Review Caddy routes before importing (Priority: P1) 🎯 MVP

**Goal**: User can open import preview and understand every discovered route plus planned local changes before anything is modified.

**Independent Test**: Connect to Caddy with supported, unsupported, and local-only routes; open import preview; verify counts, groups, statuses, reasons, and no local route changes.

### Tests for User Story 1

- [X] T019 [P] [US1] Add API preview contract test for summary counts, groups, support status, and no storage mutation in internal/api/handlers_test.go
- [X] T020 [P] [US1] Add frontend import preview rendering test for counts, groups, and read-only reasons in web/src/pages/Settings.test.tsx

### Implementation for User Story 1

- [X] T021 [US1] Add import preview summary/group response structs and route row summary helpers in internal/api/handlers.go
- [X] T022 [US1] Update PreviewImport to return summary, grouped rows, reasons, local-only removals, and warnings in internal/api/handlers.go
- [X] T023 [US1] Update previewImport client return types in web/src/lib/api.ts
- [X] T024 [US1] Replace confirm-only import prompt with an inline import review section in web/src/pages/Settings.tsx
- [X] T025 [US1] Render preview summary cards and grouped route rows in web/src/pages/Settings.tsx

**Checkpoint**: User Story 1 is fully functional and testable independently.

---

## Phase 4: User Story 2 - Commit import without losing unsupported config (Priority: P2)

**Goal**: User can confirm import and get an all-or-nothing local route replacement that preserves unsupported routes read-only.

**Independent Test**: Confirm an import with supported and unsupported routes; verify all routes are saved with correct status, and forced failures leave previous local routes unchanged.

### Tests for User Story 2

- [X] T026 [US2] Add API import success test for editable/read-only counts and saved classifications in internal/api/handlers_test.go
- [X] T027 [US2] Add API import failure test proving previous local routes remain unchanged in internal/api/handlers_test.go
- [X] T028 [P] [US2] Add frontend import confirmation result test for imported counts and warnings in web/src/pages/Settings.test.tsx

### Implementation for User Story 2

- [X] T029 [US2] Update ImportFromCaddy to use ReplaceAllRoutes and return editable/read-only/unsupported counts in internal/api/handlers.go
- [X] T030 [US2] Remove delete-then-create partial import behavior from ImportFromCaddy in internal/api/handlers.go
- [X] T031 [US2] Update importFromCaddy client return type in web/src/lib/api.ts
- [X] T032 [US2] Add confirm import action and success/error state for preview results in web/src/pages/Settings.tsx

**Checkpoint**: User Story 2 works independently and preserves previous state on failure.

---

## Phase 5: User Story 3 - Keep preserved routes truly read-only (Priority: P3)

**Goal**: Read-only preserved routes are visible with details, but cannot be edited, deleted, toggled, or field-mutated.

**Independent Test**: Import a read-only route; verify UI shows details only, API mutation attempts fail with conflict, and sync keeps raw route unchanged.

### Tests for User Story 3

- [X] T033 [US3] Add API mutation rejection tests for update, delete, and toggle of read-only routes in internal/api/handlers_test.go
- [X] T033a [US3] Add API sync failure warning test proving persisted routes remain and recovery guidance is returned in internal/api/handlers_test.go
- [X] T034 [US3] Add API details endpoint test for preserved raw route JSON in internal/api/handlers_test.go
- [X] T035 [P] [US3] Add dashboard test proving read-only route cards hide edit/delete/toggle and show View JSON/Details in web/src/pages/Dashboard.test.tsx

### Implementation for User Story 3

- [X] T036 [US3] Reject update/delete/toggle for read-only routes with 409 Conflict in internal/api/handlers.go
- [X] T036a [US3] Ensure sync failure responses include warning status and recovery guidance without deleting persisted routes in internal/api/handlers.go
- [X] T037 [US3] Add GET /api/routes/:id/details route wiring under the existing /api route group in internal/api/routes.go
- [X] T038 [US3] Implement read-only route details handler returning raw_caddy_route in internal/api/handlers.go
- [X] T039 [US3] Add getRouteDetails client method and response types in web/src/lib/api.ts
- [X] T040 [US3] Update dashboard route cards to show Unsupported/managed outside UI label and hide mutation controls in web/src/pages/Dashboard.tsx
- [X] T041 [US3] Add read-only Details/View JSON display on the dashboard in web/src/pages/Dashboard.tsx
- [X] T042 [US3] Block direct edit form access for read-only routes with details-only copy in web/src/pages/RouteForm.tsx

**Checkpoint**: User Story 3 works independently and read-only routes are truly non-mutating.

---

## Phase 6: User Story 4 - Warn users about manual Caddy drift (Priority: P4)

**Goal**: Users see honest v1 copy that manual Caddy changes are not automatically merged and should re-run import review before syncing.

**Independent Test**: View Settings import area and Dashboard sync area; verify manual drift warning copy is visible.

### Tests for User Story 4

- [X] T043 [P] [US4] Add frontend settings drift warning test in web/src/pages/Settings.test.tsx
- [X] T044 [P] [US4] Add frontend dashboard drift warning test in web/src/pages/Dashboard.test.tsx

### Implementation for User Story 4

- [X] T045 [US4] Add manual drift warning copy to Configuration Import section in web/src/pages/Settings.tsx
- [X] T046 [US4] Add manual drift warning copy near Sync to Caddy controls in web/src/pages/Dashboard.tsx

**Checkpoint**: User Story 4 works independently.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final validation, cleanup, and docs aligned with the safety promise.

- [X] T047 [P] Update README supported/unsupported import behavior and read-only route documentation in README.md
- [X] T048 [P] Update quickstart validation results and manual test notes in specs/001-safe-import-preview/quickstart.md
- [X] T049 Add or document a 100-route import preview timing check in internal/api/handlers_test.go or specs/001-safe-import-preview/quickstart.md
- [X] T050 Run focused backend test suite and fix regressions in internal/config/parser_test.go, internal/config/builder_test.go, internal/storage/sqlite_test.go, and internal/api/handlers_test.go
- [X] T051 Run frontend tests, type-check, and lint; fix regressions in web/src/pages/Settings.tsx, web/src/pages/Dashboard.tsx, web/src/pages/RouteForm.tsx, and web/src/lib/api.ts
- [X] T052 Run full make test and record final validation command result in specs/001-safe-import-preview/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies.
- **Foundational (Phase 2)**: Depends on Setup; blocks all user stories.
- **User Stories (Phases 3-6)**: Depend on Foundation. Implement in priority order for simplest delivery: US1 → US2 → US3 → US4.
- **Polish (Phase 7)**: Depends on completed desired stories.

### User Story Dependencies

- **US1 (P1)**: Requires Foundation only; MVP preview without changing local routes.
- **US2 (P2)**: Requires Foundation; benefits from US1 UI but API import can be tested independently.
- **US3 (P3)**: Requires Foundation; can be tested with seeded read-only routes even before full UI import polish.
- **US4 (P4)**: Requires existing Settings/Dashboard pages and follows Foundation for the planned delivery sequence.

### Within Each User Story

- Write tests first and confirm they fail.
- Backend contract/model changes before frontend client usage.
- API behavior before UI rendering.
- Story checkpoint before moving to the next priority.

## Parallel Opportunities

- T003 can run alongside T001/T002 because it writes research.md, not quickstart.md.
- T006 and T007 can run in parallel with storage/API foundation tests because they touch different files.
- T019 and T020 can run in parallel for US1.
- T028 can run in parallel with backend US2 tests.
- T035 can run in parallel with backend US3 tests.
- T043 and T044 can run in parallel for US4.
- T047 and T048 can run in parallel during polish.

## Parallel Example: User Story 1

```bash
Task: "T019 [US1] Add API preview contract test in internal/api/handlers_test.go"
Task: "T020 [US1] Add frontend import preview rendering test in web/src/pages/Settings.test.tsx"
```

## Parallel Example: User Story 3

```bash
Task: "T033 [US3] Add API mutation rejection tests in internal/api/handlers_test.go"
Task: "T035 [US3] Add dashboard read-only actions test in web/src/pages/Dashboard.test.tsx"
```

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup.
2. Complete Phase 2: Foundation.
3. Complete Phase 3: User Story 1.
4. Stop and validate import preview without committing any import.

### Incremental Delivery

1. Foundation: route metadata, classification, validation, transactional storage, raw preservation.
2. US1: truthful preview.
3. US2: safe confirmed import.
4. US3: strict read-only behavior and details.
5. US4: manual drift warning copy.
6. Polish: docs and full validation.

### Lazy Scope Guard

No new dependencies, no raw JSON editor, no per-route import checkboxes, no drift detection, no Layer4 editing, and no broad import-session abstraction unless a later spec requires it.

---

## Phase 8: Convergence

- [X] T053 CRITICAL: Make sync/build validation fail with recovery guidance instead of silently omitting invalid read-only raw routes per FR-013 and Constitution I (partial)
- [X] T054 CRITICAL: Add handler-specific import and sync validation so invalid editable route configs are rejected or preserved read-only instead of later being dropped per FR-016 and Constitution II (partial)
- [X] T055 Implement preview row fallback summaries for routes without host/domain using path, destination/root, or an unknown placeholder per FR-004 (partial)
- [X] T056 Review, document, or remove raw_caddy_route exposure in import-preview rows so it matches the import-flow contract and FR-016 secret-handling expectations per Contract: import preview row (unrequested)
