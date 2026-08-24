# Feature Specification: Safe Import Preview

> **Archived:** This specification records the original import-preview design. It has been superseded by the managed setup and ownership flow documented in the root README.

**Feature Branch**: `docs/adopt-constitution-v1`

**Created**: 2026-07-03

**Status**: Draft

**Input**: User description: "Trustworthy Adoption milestone notes, especially Spec Candidate 1: Safe Import Preview"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Review Caddy routes before importing (Priority: P1)

As a self-hoster with an existing Caddy setup, I want to preview what the app found in my current Caddy configuration before importing, so I can trust that the app understands what it will manage and what it will preserve.

**Why this priority**: This is the core adoption moment. Users will not trust the app if import is a blind replace operation.

**Independent Test**: Can be fully tested by connecting to a Caddy configuration with supported, partially supported, unsupported, and local-only routes, opening import preview, and verifying the review explains every route and planned change before any local route list changes.

**Acceptance Scenarios**:

1. **Given** Caddy has supported and unsupported HTTP routes, **When** the user opens import preview, **Then** the user sees totals for routes found, editable routes, read-only preserved routes, unsupported routes, and routes that will replace local entries.
2. **Given** an imported route cannot be safely edited by the app, **When** the preview lists that route, **Then** it is clearly marked read-only with a specific reason such as unsupported matcher, unknown handler, multi-handler pipeline, custom module, or another detected safety reason.
3. **Given** the app has local routes that are not present in Caddy, **When** the preview is shown, **Then** those routes appear in a local-only group that explains they will be removed if the import is confirmed.

---

### User Story 2 - Commit import without losing unsupported config (Priority: P2)

As a user, I want confirming import to bring all discovered HTTP routes into the app while preserving routes the app cannot edit, so unsupported Caddy configuration is not silently dropped or changed.

**Why this priority**: Safe commit behavior is the safety promise behind the preview. Without it, the preview is only cosmetic.

**Independent Test**: Can be fully tested by confirming an import containing supported and unsupported HTTP routes, then verifying the local route list is replaced as one complete change and unsupported routes remain present as read-only preserved entries.

**Acceptance Scenarios**:

1. **Given** the user confirms an import preview, **When** all routes can be classified and saved, **Then** editable routes become normal managed routes and unsupported or partially supported routes become read-only preserved routes.
2. **Given** import confirmation fails because Caddy cannot be fetched, the response cannot be understood, or a route cannot be classified or saved, **When** the failure is reported, **Then** the previous local route list remains unchanged.
3. **Given** unsupported HTTP routes are found, **When** import is confirmed, **Then** import still succeeds and those routes are preserved as read-only unless a separate failure occurs.

---

### User Story 3 - Keep preserved routes truly read-only (Priority: P3)

As a user, I want unsupported imported routes to remain visible but not editable, so I can understand what exists without accidentally breaking configuration the app does not understand.

**Why this priority**: Visibility without mutation prevents the app from hiding important config while avoiding unsafe partial edits.

**Independent Test**: Can be fully tested after importing a read-only preserved route by checking that its summary and technical details are available, while edit, delete, enable/disable, and field-changing controls are unavailable.

**Acceptance Scenarios**:

1. **Given** a read-only preserved route exists, **When** the user views the dashboard, **Then** the route appears with its detected domain, path when available, handler or destination summary when available, read-only label, and preservation reason.
2. **Given** a read-only preserved route exists, **When** the user opens its actions, **Then** the only mutation-free action available is viewing details or raw JSON.
3. **Given** Caddy sync is run after importing read-only preserved routes, **When** the sync completes, **Then** the preserved routes are included unchanged alongside editable managed routes.

---

### User Story 4 - Warn users about manual Caddy drift (Priority: P4)

As a user who may still edit Caddy manually, I want clear copy explaining that live Caddy may have changed since the last import or sync, so I know when to re-run import review before syncing.

**Why this priority**: Drift detection can wait, but users need honest guidance in v1 to avoid false confidence.

**Independent Test**: Can be tested by viewing import and sync areas and confirming they include a visible warning that manual Caddy edits require re-running import review before syncing.

**Acceptance Scenarios**:

1. **Given** the user is about to sync or review import status, **When** the page explains current safety limits, **Then** it states that manual Caddy edits after the last import or sync are not automatically merged and recommends re-running import review first.

### Edge Cases

- Caddy cannot be reached: preview and commit fail with no local route changes.
- Caddy returns data that cannot be understood: preview and commit fail with no local route changes.
- A route has no easily detected domain, path, destination, or root: the route still appears with the best available summary and a reason if it is read-only.
- A route uses multiple unsupported features: the app shows at least one clear primary reason and may show additional reasons.
- A route is only partially supported: it is preserved as read-only rather than edited as if fully supported.
- No HTTP routes are discovered: preview explains that confirming import will remove local managed HTTP routes.
- Imported route details may contain sensitive-looking values: technical details are exposed only through the existing admin UI/API surface, protected by Basic Auth when configured, and errors or summaries do not expose secrets unnecessarily.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST provide an import preview before any import confirmation changes the local route list.
- **FR-002**: The preview MUST show summary counts for total discovered HTTP routes, editable routes, read-only preserved routes, unsupported routes, and local routes that will be removed or replaced.
- **FR-003**: The preview MUST place each route row in one planned-change group, using this priority: local only/will be removed, read-only preserved, already exists/will update, or new from Caddy.
- **FR-004**: Each preview route row MUST show the best available domain, path, route type, destination or root, support status, and read-only reason when applicable, falling back from host matcher to path matcher to handler destination/root to an unknown placeholder.
- **FR-005**: The system MUST classify every discovered HTTP route as editable/supported, partially supported/read-only preserved, or unsupported/read-only preserved.
- **FR-006**: The system MUST NOT silently skip discovered HTTP routes during preview or confirmed import.
- **FR-007**: Unsupported or partially supported HTTP routes MUST NOT block import solely because they are unsupported; they MUST be preserved as read-only unless a separate failure occurs.
- **FR-008**: Confirming import MUST replace the local route list as one complete change: either all classified routes are saved and the previous local route list is replaced, or no local route changes are made.
- **FR-009**: If Caddy cannot be fetched, the response cannot be understood, or routes cannot be classified or saved, the system MUST report the failure and leave the previous local route list unchanged.
- **FR-010**: Read-only preserved routes MUST be visible in the main route list with a clear "Unsupported / managed outside UI" or equivalent label.
- **FR-011**: Read-only preserved routes MUST NOT allow edit, delete, enable/disable, or field-changing actions.
- **FR-012**: Read-only preserved routes MUST allow viewing technical details or raw JSON for troubleshooting through the existing admin UI/API surface, protected by Basic Auth when configured.
- **FR-013**: When syncing, editable managed routes MUST be represented from the app-managed values, while read-only preserved routes MUST be included unchanged from the imported route details.
- **FR-014**: The feature MUST include user-visible copy explaining that manual Caddy changes after the last import or sync are not automatically merged in v1 and that users should re-run import review before syncing after manual edits.
- **FR-015**: The feature MUST NOT include per-route import include/exclude checkboxes, a raw JSON editor, full line-by-line JSON diffs, automatic drift detection, editable Layer4 support, or editable advanced matcher/custom module support in the first slice.
- **FR-016**: The system MUST validate imported route data before saving or syncing: preserved read-only routes require valid raw route JSON, support status, and read-only reason; editable routes require valid domain, handler type, and handler config; import errors and summaries MUST avoid exposing credentials, tokens, or full sensitive raw config.

### Constitution Requirements *(mandatory for affected features)*

- **CR-001**: The feature imports and syncs Caddy routes, so it MUST preserve unsupported or unknown Caddy configuration unless the user explicitly deletes or replaces it.
- **CR-002**: The feature exposes admin UI/API behavior, so it MUST define input validation, authentication impact, and secret-handling expectations: import and technical details use the existing admin UI/API surface with Basic Auth when configured, imported route data is validated before local changes or sync, and summaries/errors avoid unnecessary secret exposure.
- **CR-003**: If local persistence can succeed while Caddy sync fails, the sync result MUST show a user-visible warning and a recovery path that recommends reviewing the warning, checking Caddy connectivity, and retrying sync or import review as appropriate.

### Key Entities *(include if feature involves data)*

- **Import Preview**: A pending review of discovered Caddy HTTP routes, current local routes, classifications, counts, and planned local changes before confirmation.
- **Imported Route**: A route discovered from Caddy with the best available domain, path, route type, destination/root, support status, and preservation details.
- **Support Classification**: The safety status for a discovered route: editable/supported, partially supported/read-only preserved, or unsupported/read-only preserved.
- **Read-only Reason**: A user-facing explanation of why a route cannot be edited safely, such as unsupported matcher, unknown handler, multi-handler pipeline, custom module, or other unsupported structure.
- **Route Change Summary**: The preview grouping that explains whether a route is new from Caddy, will update an existing local route, is local-only and will be removed, or will be preserved read-only.
- **Technical Details**: The original route details shown only for troubleshooting read-only or unsupported routes, without allowing raw editing.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In test imports, 100% of discovered HTTP routes are either shown in the preview with a classification or the preview fails before any local route changes occur.
- **SC-002**: A user can identify how many routes are editable, read-only preserved, unsupported, and local-only within 60 seconds of opening the import preview.
- **SC-003**: Confirmed import failure cases leave the prior local route list unchanged in 100% of covered failure tests.
- **SC-004**: Read-only preserved routes expose zero edit, delete, enable/disable, or field-changing actions in user-facing route controls.
- **SC-005**: Sync after importing preserved routes keeps those routes byte-for-byte or semantically unchanged in 100% of covered preservation tests.
- **SC-006**: At least 90% of users in a review or demo can correctly state what will happen before confirming import using only the preview screen.

## Assumptions

- The first slice focuses on HTTP routes discovered from Caddy; Layer4 and other non-HTTP apps are roadmap unless they can be safely shown as unmanaged without expanding scope.
- All discovered HTTP routes are imported/classified together when the user confirms; per-route include/exclude selection is deferred.
- Existing admin UI/API routing, optional Basic Auth settings, and Caddy connection settings are reused for import, preview, details, and sync access.
- Local-only routes are removed when the user confirms import because the import adopts Caddy as the current source for HTTP routes at that moment.
- Read-only preserved route details are for troubleshooting only, not editing.
- Automatic live drift detection is deferred; v1 relies on clear user guidance to re-run import review after manual Caddy edits.
