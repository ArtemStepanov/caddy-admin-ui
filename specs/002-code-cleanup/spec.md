# Feature Specification: Code Cleanup

**Feature Branch**: `main`

**Created**: 2026-07-05

**Status**: Draft

**Input**: User description: "separate spec to address #7 #8 and other code cleanups. for now start with obvious places, we will refine it later. i dont want to forget about these findings."

## Clarifications

### Session 2026-07-05

- Q: What should this cleanup spec be named? → A: code-cleanup

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Consistent read-only route behavior (Priority: P1)

As a maintainer reviewing or changing route safety behavior, I want read-only route decisions to be defined once and applied consistently, so routes cannot be treated as editable in one place and protected in another.

**Why this priority**: Inconsistent read-only decisions can create unsafe mutations, confusing UI behavior, or sync failures. This is the highest-risk cleanup from the review findings.

**Independent Test**: Can be fully tested by creating representative editable, partially read-only, unsupported read-only, and legacy preserved routes, then verifying every user-facing and mutation-protection path classifies each route the same way.

**Acceptance Scenarios**:

1. **Given** a route is classified as read-only, **When** the route appears in route lists, mutation controls, sync validation, and technical details, **Then** every area treats it as read-only.
2. **Given** a route is classified as editable, **When** the same areas inspect it, **Then** none of them block it solely because of stale or conflicting read-only metadata.
3. **Given** a future cleanup changes the read-only rule, **When** the shared rule is updated, **Then** all affected behavior uses the updated rule without separate copy changes.

---

### User Story 2 - Single drift warning message (Priority: P2)

As an operator using import and sync screens, I want manual Caddy drift guidance to appear once in each relevant context with consistent wording, so I receive clear safety guidance without duplicate or conflicting messages.

**Why this priority**: Duplicate warnings reduce trust and make future copy changes error-prone, but the cleanup is lower risk than route safety classification.

**Independent Test**: Can be fully tested by viewing import preview, import results, and sync guidance areas and confirming the warning appears once per context with the same approved wording.

**Acceptance Scenarios**:

1. **Given** the user opens an import or sync area, **When** manual drift guidance is shown, **Then** the wording is consistent across surfaces.
2. **Given** the backend or frontend provides a drift warning, **When** the page renders it, **Then** the user does not see the same sentence twice in the same section.
3. **Given** the warning wording changes later, **When** maintainers update the approved message, **Then** all relevant surfaces can be verified from one expected copy source or one documented contract.

---

### User Story 3 - Remove redundant cleanup debt (Priority: P3)

As a maintainer, I want obvious duplicated or redundant cleanup findings recorded and addressed in a bounded pass, so cleanup work does not get forgotten or grow into speculative refactoring.

**Why this priority**: These cleanups improve maintainability but should not expand beyond review-confirmed duplication and redundant behavior.

**Independent Test**: Can be fully tested by reviewing the cleanup candidate list before and after the change and confirming each candidate is either removed, consolidated, or explicitly deferred with rationale.

**Acceptance Scenarios**:

1. **Given** a cleanup candidate duplicates an existing validation or summary calculation, **When** the cleanup is completed, **Then** only one necessary source of the behavior remains.
2. **Given** repeated test setup or expected data makes tests harder to maintain, **When** the cleanup is completed, **Then** repeated setup is reduced without weakening test coverage.
3. **Given** a candidate is not worth changing now, **When** the cleanup pass completes, **Then** the candidate is documented as intentionally deferred instead of silently forgotten.

### Edge Cases

- A cleanup candidate looks duplicated but protects a distinct edge case; it should be documented and left intact.
- A copy consolidation could accidentally remove user-visible guidance from one surface; the affected screen must still show the warning where required.
- A read-only classification cleanup could change behavior for legacy or partially imported routes; representative legacy and imported route cases must be checked.
- A cleanup pass could grow into broad refactoring; work must stay limited to review-confirmed candidates and obvious nearby duplication.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST use one consistent read-only route classification rule across route display, mutation prevention, route details, import handling, and sync validation.
- **FR-002**: Read-only classification behavior MUST remain unchanged for supported, partially supported, unsupported, and legacy preserved routes unless a separate safety fix explicitly requires the change.
- **FR-003**: Manual Caddy drift warning copy MUST be consistent across import and sync guidance surfaces.
- **FR-004**: The same manual Caddy drift warning MUST NOT be rendered twice in the same visible context.
- **FR-005**: Redundant import validation, summary calculation, or route-safety checks MUST be consolidated when they do not protect distinct behavior.
- **FR-006**: Cleanup of duplicated test setup or expected data MUST preserve coverage for import preview, read-only route behavior, and drift warning visibility.
- **FR-007**: The cleanup pass MUST record any reviewed candidate that is intentionally not changed, including the reason it remains separate.
- **FR-008**: The cleanup pass MUST NOT add new user-facing features, new dependencies, new route editing capabilities, or broader import behavior changes.

### Constitution Requirements *(mandatory for affected features)*

- **CR-001**: Because the cleanup touches route import, sync, and read-only behavior, it MUST preserve unsupported or unknown Caddy configuration and MUST NOT weaken read-only protection.
- **CR-002**: Because the cleanup touches admin UI/API behavior, it MUST preserve existing input validation, authentication assumptions, and secret-handling expectations.
- **CR-003**: Because the cleanup touches sync-related behavior and warnings, it MUST preserve user-visible recovery guidance when sync or import state is uncertain.

### Key Entities *(include if feature involves data)*

- **Cleanup Candidate**: A known duplicated, redundant, or inconsistent behavior identified during review, with a decision to consolidate, leave intact, or defer.
- **Read-only Route Classification**: The shared decision of whether a route is safe to mutate through the app or must remain preserved and non-mutating.
- **Manual Drift Warning**: The approved user-facing guidance that manual Caddy edits are not automatically merged and import review should be rerun before syncing.
- **Deferred Cleanup**: A reviewed candidate intentionally left for later with a short reason, so it is not rediscovered as forgotten debt.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of covered route safety scenarios produce the same read-only/editable decision across all checked user-facing and mutation-protection areas.
- **SC-002**: The manual drift warning appears no more than once in each relevant visible context during review.
- **SC-003**: At least the two originally reported cleanup findings are either completed or explicitly deferred with rationale.
- **SC-004**: No existing import, sync, or read-only route safety checks lose coverage during the cleanup pass.
- **SC-005**: Reviewers can identify the remaining cleanup candidates and their status within 5 minutes using the feature artifacts or resulting notes.

## Assumptions

- This is an overall maintenance and debt-tracking feature, not a user-facing feature expansion.
- The first cleanup pass starts with the confirmed review findings: read-only predicate duplication, duplicate drift warning copy/rendering, redundant import validation, redundant summary calculation, and obvious duplicated test setup.
- Behavior-preserving consolidation is preferred; any behavior change belongs in the active safe-import PR tasks or a separate bug-fix spec.
- No new runtime dependency or architectural layer is justified for this cleanup.
