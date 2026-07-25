# Feature Specification: Code Cleanup

**Feature Branch**: `main`

**Created**: 2026-07-09

**Status**: Implemented

**Input**: User description: "we are alreadt on spec 002, but i need you to use ponytail and other lsp/best practices to narrow the cleanup scope in the current codebase. we have backend on go and FE on react. go ahead, and lets improve the spec."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Safe import validation (Priority: P1)

As an operator importing server configuration, I want every imported route to receive one consistent safety decision, so valid routes import and preserved routes remain protected without conflicting checks.

**Why this priority**: Import decisions can affect live traffic and configuration preservation. This is the smallest high-risk cleanup with direct operator value.

**Independent Test**: Can be fully tested by importing editable, partially supported, unsupported, legacy-preserved, and invalid routes and verifying each receives the expected accept, preserve, or reject outcome.

**Acceptance Scenarios**:

1. **Given** an imported route is editable and complete, **When** it is reviewed for import, **Then** it is accepted exactly once and remains editable.
2. **Given** an imported route must be preserved, **When** it is reviewed for import, **Then** its original content and its non-mutable status are retained.
3. **Given** an imported route is incomplete or unsafe to process, **When** it is reviewed for import, **Then** import stops before replacing local routes and reports an actionable error.

---

### User Story 2 - Mode-safe route editing (Priority: P2)

As an operator creating or editing a supported route, I want each route kind to expose and save only its applicable settings, so switching route kinds cannot carry incompatible settings into the saved route.

**Why this priority**: This narrows the client-side cleanup to the form's existing supported route kinds, improving maintainability without adding route capabilities.

**Independent Test**: Can be fully tested by creating, editing, and switching among every supported route kind, then verifying saved settings match the selected kind and the existing route behavior is unchanged.

**Acceptance Scenarios**:

1. **Given** an operator selects a supported route kind, **When** the editor is displayed, **Then** it shows only the settings applicable to that kind.
2. **Given** an operator switches route kinds before saving, **When** the route is saved, **Then** settings from the previous kind are not included.
3. **Given** an operator opens a preserved read-only route, **When** it is viewed, **Then** it remains non-mutable and its preserved content is not interpreted as an editable route.

### Edge Cases

- A legacy preserved route has incomplete metadata but retains original route content; it must remain protected rather than become editable during cleanup.
- A route changes kind after settings have been entered; incompatible settings must not survive the switch.
- An imported route passes one former check but fails another; the unified decision must reject it before local routes are replaced.
- A proposed cleanup does not directly improve either import safety or route-kind isolation; it must be left out of this feature.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST make one authoritative safety decision for each imported route before local routes can be replaced.
- **FR-002**: The authoritative import decision MUST accept complete editable routes, preserve non-editable routes with their original content and reason, and reject incomplete or unsafe routes with an actionable error.
- **FR-003**: Existing legacy and unsupported routes MUST retain their non-mutable status and original content throughout import, viewing, and synchronization.
- **FR-004**: For every supported route kind, the route editor MUST present and save only settings applicable to the selected kind.
- **FR-005**: Changing a route kind before save MUST discard settings that are not applicable to the newly selected kind.
- **FR-006**: Preserved read-only routes MUST remain outside editable route configuration handling.
- **FR-007**: Automated checks MUST cover all supported route kinds and imported editable, partially supported, unsupported, legacy-preserved, and invalid route cases.
- **FR-008**: This cleanup MUST be limited to import-route validation and route-editor configuration isolation; it MUST NOT add route kinds, user-facing features, warning-copy changes, dashboard changes, or unrelated restructuring.

### Constitution Requirements *(mandatory for affected features)*

- **CR-001**: The cleanup MUST preserve unsupported or unknown route configuration unless an operator explicitly deletes or replaces it.
- **CR-002**: Existing input validation, authentication behavior, and secret-handling expectations MUST remain unchanged.
- **CR-003**: If local persistence succeeds while synchronization fails, the existing user-visible warning and recovery guidance MUST remain available and unchanged.

### Key Entities *(include if feature involves data)*

- **Imported Route Decision**: The single result that classifies an imported route as editable, preserved, or rejected before import is committed.
- **Preserved Route**: A route whose original content is retained because it is unsupported, only partially supported, or otherwise unsafe to edit.
- **Route Kind**: One of the currently supported operator-selectable route behaviors, with its own applicable configuration settings.
- **Route Configuration**: The settings saved for the selected route kind.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of covered imported-route cases receive the same accept, preserve, or reject outcome wherever import safety is evaluated.
- **SC-002**: 100% of covered preserved and legacy route cases retain their original content and non-mutable status after import and synchronization checks.
- **SC-003**: For each currently supported route kind, 100% of create, edit, and kind-switch checks save only settings applicable to the selected kind.
- **SC-004**: All existing automated checks pass, and the new focused checks fail if import safety decisions or route-kind isolation regress.
- **SC-005**: A maintainer can confirm the two in-scope cleanup areas and the explicit exclusions from this specification in under 3 minutes.

## Assumptions

- This revision replaces the original broad cleanup backlog; planning for this pass uses only the bounded scope in this specification.
- The reviewed codebase identifies two justified cleanup targets: overlapping imported-route safety checks and unbounded editable-route configuration state.
- The current supported route kinds and operator-facing behavior remain unchanged.
- Simpler consolidation is preferred over new layers, shared infrastructure, or a repository-wide cleanup.
