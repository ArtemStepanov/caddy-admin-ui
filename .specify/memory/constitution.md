<!--
Sync Impact Report
Version change: template -> 1.0.0
Modified principles:
- PRINCIPLE_1_NAME -> I. Preserve Caddy Configuration
- PRINCIPLE_2_NAME -> II. Secure the Admin Surface
- PRINCIPLE_3_NAME -> III. Test Sync-Critical Paths
- PRINCIPLE_4_NAME -> IV. Keep the Stack Small
- PRINCIPLE_5_NAME -> V. Make Operations Observable and Recoverable
Added sections:
- Operational Constraints
- Development Workflow & Quality Gates
Removed sections: None
Templates requiring updates:
- ✅ .specify/templates/plan-template.md
- ✅ .specify/templates/spec-template.md
- ✅ .specify/templates/tasks-template.md
- ✅ .specify/templates/commands/*.md (no command template directory present)
Runtime guidance reviewed:
- ✅ README.md (no change required)
- ✅ CLAUDE.md (no change required)
- ✅ AGENTS.md (no change required)
Follow-up TODOs: None
-->
# Caddy Admin UI Constitution

## Core Principles

### I. Preserve Caddy Configuration
Caddy configuration managed by this project MUST round-trip without losing unsupported
or unknown Caddy route data. Import, edit, export, and sync paths MUST preserve raw
Caddy route content that the UI does not understand, unless a user explicitly deletes
or replaces it. Sync failures MUST NOT erase persisted routes.

Rationale: The UI is a safer front end for Caddy, not an authority that may discard
valid Caddy configuration.

### II. Secure the Admin Surface
The server MUST default to loopback-only access outside container networking, and any
internet- or LAN-exposed deployment MUST enable authentication. API handlers MUST
validate user input before writing SQLite data or sending Caddy Admin API requests.
Secrets MUST NOT be logged or returned in error responses.

Rationale: The application controls a reverse proxy; weak defaults can expose every
route behind it.

### III. Test Sync-Critical Paths
Changes to API mutations, Caddy config building/parsing, import/export, storage, or
authentication MUST include automated tests that fail without the change. All pull
requests MUST pass `make test` before merge. Integration or contract tests MUST cover
changes to Caddy Admin API behavior or serialized route schemas.

Rationale: Route changes affect live traffic, and regressions can break or delete
proxy behavior.

### IV. Keep the Stack Small
New services, queues, databases, frontend frameworks, or runtime dependencies MUST be
justified by a current requirement in the plan. The default stack is Go/Gin, SQLite,
Preact/TypeScript, Tailwind CSS, and direct Caddy Admin API calls. Simpler standard
library or already-installed dependency solutions MUST be preferred.

Rationale: This project is an admin tool; operational weight must stay lower than the
Caddy service it manages.

### V. Make Operations Observable and Recoverable
User-visible mutations MUST report Caddy sync status, including warnings when local
persistence succeeds but Caddy sync fails. Errors MUST include actionable context
without leaking secrets. Health/status endpoints and UI indicators MUST reflect Caddy
connectivity and sync state when those states affect user decisions.

Rationale: Operators need to know whether the UI state and running Caddy config agree,
and they need a safe recovery path when they do not.

## Operational Constraints

- Backend code MUST remain in Go under `cmd/` and `internal/` unless a plan justifies
  a different boundary.
- Frontend code MUST remain in Preact/TypeScript under `web/` unless a plan justifies
  a different boundary.
- SQLite remains the persistence layer for route and global configuration data.
- Route mutations MUST persist locally before sync attempts, and sync failures MUST be
  surfaced as warnings rather than silent success.
- Configuration values that vary by deployment MUST use documented environment
  variables or persisted global config.
- Release-impacting commits MUST follow Conventional Commits because release-please
  derives versions and changelog entries from commit messages.

## Development Workflow & Quality Gates

- Specs MUST define independently testable user stories and measurable success
  criteria before implementation planning.
- Plans MUST run the Constitution Check before Phase 0 research and repeat it after
  Phase 1 design.
- Tasks MUST group implementation by independently deliverable user story and include
  constitution-mandated tests for sync-critical, security, and data-preservation work.
- Pull requests MUST document any constitution violation in the plan's Complexity
  Tracking table, including the simpler alternative rejected.
- Development commands and architecture notes in `CLAUDE.md` and `README.md` are the
  runtime guidance for local work and MUST stay accurate when workflows change.

## Governance

This constitution supersedes conflicting project guidance. Amendments require a pull
request that explains the change, updates affected templates and runtime guidance, and
includes a Sync Impact Report in this file.

Versioning follows semantic versioning:
- MAJOR for incompatible governance changes or removed/redefined principles.
- MINOR for new principles, new required sections, or materially expanded obligations.
- PATCH for clarifications, typo fixes, and non-semantic wording changes.

Compliance is reviewed during spec planning, task generation, and pull request review.
Unresolved violations block merge unless they are explicitly recorded with rationale
and an owner-approved mitigation plan.

**Version**: 1.0.0 | **Ratified**: 2026-07-02 | **Last Amended**: 2026-07-02
