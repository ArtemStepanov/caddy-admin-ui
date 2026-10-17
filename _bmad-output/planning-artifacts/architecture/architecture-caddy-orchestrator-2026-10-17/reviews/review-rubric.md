> Historical review of the pre-correction spine. Its original verdict/claims are retained as provenance. Current dispositions and fresh revision-1 review links are in [review-dispositions](../review-dispositions.md); this old verdict is not the current readiness gate.

# Rubric Walker Review — Architecture Spine

- **Target Artifact:** `_bmad-output/planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md`
- **Reviewer:** Rubric Walker Reviewer (Independent Subagent)
- **Evaluation Date:** 2026-10-17
- **Checklist Reference:** `.agents/skills/bmad-architecture/references/reviewer-gate.md`
- **Verdict:** **Minor issues** (Strong, well-structured spine ready for implementation with minor clarifications)

---

## 1. Executive Summary

The Architecture Spine for Caddy Orchestrator is an exemplary, rigorous architectural specification. It successfully unifies the brownfield technical foundation (Go, Gin, SQLite, Preact, Tailwind CSS) with the expanded requirements from the PRD and UX reconciliation (root-level configuration ownership, Caddyfile CST source of truth, staged draft/apply batch pipeline, multi-vantage verification, and zero-trust secrets).

Mechanical validation (`lint_spine.py`) completed with **0 findings**. The semantic evaluation against the Good-Spine Checklist confirms that all 11 Architectural Decisions (ADs) are actionable and enforceable, with clear bindings and divergence prevention. The driving inputs (PRD FR-1 to FR-11 and UX U1–U23) are mapped to components and architectural rules.

A small number of minor clarifications are highlighted below regarding brownfield database migration, external Caddyfile `import` directives, Caddy Admin API network isolation, and editor validation debouncing. None of these represent critical architectural blockers.

---

## 2. Good-Spine Rubric Checklist Evaluation

### Question 1: Does it fix the real divergence points for the level below and miss none?
**Rating:** **Pass (with minor edge-case recommendations)**

- **Source of Truth & Syntax (AD-1, AD-9):** Definitively eliminates ambiguity between visual forms and code editor by establishing the formatted Caddyfile text in SQLite as authoritative, backed by a lossless Concrete Syntax Tree (CST) and validated remotely via Caddy's `POST /adapt`.
- **Configuration Ownership Scope (AD-2):** Resolves the tension between single-server ownership and full instance management by introducing progressive root-level ownership with `passthrough_apps` JSON isolation for non-HTTP applications.
- **Publication & Mutation Model (AD-3, AD-8, AD-10):** Eliminates premature runtime writes by enforcing an explicit SQLite `DraftState` singleton, pre-write immutable recovery snapshots, live hash drift locking (409 Conflict), and an uncertain state handler for network timeouts.
- **Verification Semantics (AD-4):** Clarifies that Caddy configuration acceptance (`POST /load` 200 OK) is not equivalent to service availability; establishes asynchronous typed probes with explicit vantage point tagging (`manager_network` vs `client_perspective`).
- **Secrets & Discovery (AD-5, AD-6):** Decouples Docker discovery as a read-only auxiliary module and protects DNS credentials using AES-GCM-256 encryption at rest with UI masking.
- *Omission / Edge-Case Note:* The spine does not explicitly address how Caddyfile multi-file `import` directives (e.g., `import /etc/caddy/snippets/*`) are handled if adopted from an existing multi-file Caddy setup (see Finding 2).

### Question 2: Is every AD's Rule enforceable and does it actually prevent its stated divergence?
**Rating:** **Pass**

All 11 ADs are articulated with concrete, testable, imperative rules:
- **AD-1 (CST Model):** Enforces trivia preservation, error banners on invalid editor input, and fallback to "Editor-managed" blocks for unsupported directives.
- **AD-2 (Root-Level Ownership):** Enforces a 5-step adoption algorithm and baseline snapshot archiving.
- **AD-3 (Staged Pipeline):** Enforces pre-check, drift verification, pre-write snapshot creation, and max 100 FIFO pruning.
- **AD-4 (Multi-Vantage Verification):** Enforces typed probe assertions (reverse proxy, redirect, file server, TLS inspection) and distinguishes container reachability from client browser reachability.
- **AD-5 (Zero-Trust Secrets):** Enforces AES-GCM-256 encryption at rest, masking in APIs/logs/diffs, and upfront Caddy module capability checks before saving DNS-01 configurations.
- **AD-6 (Docker Discovery):** Restricts Docker integration to read-only endpoints, mandates graceful degradation to manual IP:port input, and bans lifecycle operations.
- **AD-7 (Centralized State):** Enforces single-active-draft singleton in SQLite and forbids UI components from holding private draft copies.
- **AD-8 (Drift Locking):** Enforces SHA-256/ETag comparison before commit and mandates explicit operator rebase/discard choices.
- **AD-9 (Remote Adapt):** Forbids bundling a local Caddy binary in the manager container, delegating compilation to the target Caddy instance.
- **AD-10 (Ambiguous Apply & Persistence):** Mandates state freeze on dropped connection and specifies external Caddy `--resume` and volume mounts.
- **AD-11 (Recovery Retention & Backup):** Enforces FIFO cleanup and defines the three-part backup package.

### Question 3: Does anything under Deferred let two units diverge?
**Rating:** **Pass**

The Deferred section explicitly lists four capabilities:
1. *Multi-Instance Fleet Management*: Deferred with clear single-instance boundary (`CADDY_ADMIN_URL`).
2. *Multi-User RBAC & OAuth/OIDC*: Deferred with single-user Basic Auth security boundary.
3. *Automated External DNS Record Creation*: Deferred with clear boundary that manager configures Caddy ACME challenges, not external DNS zones.
4. *Live Streaming Access Log Ingestion*: Deferred with targeted active diagnostic probes as the designated observability mechanism.

None of these deferred capabilities create dangling interfaces, conflicting contracts, or seams where one module assumes another module will deliver the capability.

### Question 4: Does it ratify rather than contradict the brownfield codebase (Gin, SQLite, Preact, Tailwind)?
**Rating:** **Pass (with migration note)**

- **Stack Alignment:** The spine faithfully retains Go 1.25.13, Gin 1.12.0, SQLite3 (mattn/go-sqlite3 1.14.50), Node 24-alpine, Preact 10.29.8, Preact Router 4.1.2, Tailwind CSS 3.4.1, and Vite 8.2.2.
- **Architectural Refactoring:** The shift from the brownfield scoped server route array (`internal/config/builder.go`, `internal/config/parser.go`) to Hexagonal CST-based root management (`internal/domain/cst`, `internal/app/draft`) is a deliberate, approved evolution based on the PRD-UX reconciliation.
- *Migration Note:* The brownfield database schema in `internal/storage/sqlite.go` currently uses `routes` and `snapshots` tables. The spine introduces `CONFIG_METADATA`, `DRAFT`, `RECOVERY_POINT`, `SECRET_CREDENTIAL`, and `VERIFICATION_RESULT`. The spine should clarify whether brownfield database files will be migrated or whether existing Caddy live state is re-adopted into the new schema on first boot (see Finding 1).

### Question 5: Does it cover the driving inputs (PRD FR-1 to FR-11, UX U1-U23)?
**Rating:** **Pass**

- The Capability → Architecture Map explicitly maps FR-1 through FR-11 and key UX features to specific application packages and governing ADs.
- All UX decisions (U1 to U23) from `reconcile-inputs.md` and `.memlog.md` are addressed:
  - U2, U14 (Quiet Workshop visual style, responsive ergonomics): Consistency Conventions and layout structure.
  - U7, U11 (Core proxy, redirect, and static site flows): AD-1, AD-4.
  - U8 (Docker discovery + manual fallback): AD-6.
  - U9 (Desktop primary + mobile quick status/corrections): AD-7, Consistency Conventions (touch targets, responsive cards).
  - U15 (Type-aware verification distinct from config acceptance): AD-4.
  - U16 (All 4 TLS scenarios immediately): AD-1, AD-4, AD-5.
  - U17, U22 (Forms + Caddyfile editor peer views): AD-1, AD-9.
  - U18 (Batch draft, shared diff, apply, verify): AD-3, AD-7, AD-8.
  - U19, U20 (Connect existing Caddy, adopt entire configuration): AD-2.

### Question 6: Is the operational/environmental envelope complete (deployment, persistence, external Caddy, backup, operations)?
**Rating:** **Pass (with minor security note)**

- **Deployment:** Multi-stage `Dockerfile`, production manager-only Compose (`docker-compose.prod.yml`), and local development stack with fixture (`docker-compose.yml`).
- **Persistence Contract:** Explicitly mandates external Caddy `--resume` / `autosave.json` and persistent volumes (`/config`, `/data`).
- **Backup & Disaster Recovery:** Explicitly specifies the complete 3-part backup package (SQLite DB, encryption key, and external Caddy `/data` volume) and prevents unbounded DB growth with 100-snapshot FIFO retention.
- **Operations & Security:** Enforces HTTP Basic Auth on non-loopback listeners, paired credentials, CSRF headers (`X-Requested-With`), constant-time credential comparison, and drift locking.
- *Security Note:* Exposing Caddy's Admin API (`:2019`) across containers/hosts requires explicit network segmentation guidance (see Finding 3).

---

## 3. Findings & Detailed Observations

### Finding 1 (Medium): Brownfield SQLite Database Schema Migration Path
- **Observation:** Existing installations maintain a SQLite database with `routes`, `global_config`, and `snapshots` tables (defined in `internal/storage/sqlite.go`). The spine's ERD defines a clean new schema (`CONFIG_METADATA`, `DRAFT`, `RECOVERY_POINT`, `SECRET_CREDENTIAL`, `VERIFICATION_RESULT`).
- **Risk:** When an existing user upgrades their container, the existing `routes.db` file may cause schema conflicts or leave orphaned tables if the storage adapter expects the new schema.
- **Recommendation:** Add an explicit rule or note in `internal/adapters/storage/sqlite`: during startup migration, if legacy `routes` table exists and `DRAFT` does not, either execute an automatic one-time migration or archive the legacy tables and bootstrap `DRAFT` directly from the live Caddy instance via initial adoption (AD-2).

### Finding 2 (Medium): Handling of External Multi-File Caddyfile `import` Directives
- **Observation:** In real-world Caddy configurations (including the homelab reference reviewed in UX discovery), users frequently use `import` directives (e.g. `import /etc/caddy/snippets/*` or `import conf.d/*.conf`). The spine models the configuration as a single authoritative `caddyfile_content` text field in SQLite.
- **Risk:** If an existing Caddyfile uses external file imports, calling `POST /adapt` remotely on Caddy works only if those files exist on the Caddy host filesystem, but the manager has no visibility into them, and modifying imported blocks from visual forms is undefined.
- **Recommendation:** Clarify in AD-1/AD-2 that imported snippets defined within the Caddyfile (`(snippet_name) { ... }`) are fully parsed and editable, while external file path imports (`import /path/*`) are treated as opaque, read-only directive blocks that are preserved verbatim in the CST.

### Finding 3 (Low): Network Isolation & Security for Caddy Admin API
- **Observation:** AD-9 and AD-2 require the manager to communicate with Caddy's Admin API over HTTP (`CADDY_ADMIN_URL`). Caddy's Admin API does not have native authentication.
- **Risk:** In multi-container or cross-host deployments, if the operator binds port 2019 to `0.0.0.0`, the Admin API could be exposed to the local network without authentication.
- **Recommendation:** In deployment documentation and `docker-compose.prod.yml`, specify that external Caddy and Caddy Orchestrator should share an internal, non-exposed Docker bridge network, or that Caddy's admin endpoint must listen strictly on `127.0.0.1:2019` when running on the same host network.

### Finding 4 (Low): Editor Validation Trigger and Debounce Specification
- **Observation:** AD-1 states: "Mutations in the Caddyfile editor re-parse the CST. If the editor input contains syntax errors (as verified by Caddy remote compilation), switching back to visual forms is blocked with an error banner...".
- **Risk:** Calling remote `POST /adapt` on every single keystroke in the Monaco/Code editor could flood the Caddy Admin API and introduce typing latency.
- **Recommendation:** Clarify that local CST lexical/syntactic checks run in the browser/backend on debounce (e.g., 500ms), while full remote semantic compilation via Caddy's `POST /adapt` is triggered on explicit blur, form switch, or save draft action.

---

## 4. Conclusion & Acceptance Decision

The Architecture Spine `ARCHITECTURE-SPINE.md` is **APPROVED with Minor Issues**. The architecture is robust, well-delimited, and adheres to the project's non-negotiable invariants and conventions. The 4 findings above are non-blocking refinements that can be addressed directly in story definitions or during sprint execution.
