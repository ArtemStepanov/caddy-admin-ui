> Historical review of the pre-correction spine. Its original verdict/claims are retained as provenance. Current dispositions and fresh revision-1 review links are in [review-dispositions](../review-dispositions.md); this old verdict is not the current readiness gate.

# Adversarial Architecture Review: Caddy Orchestrator Spine

**Document Under Review:** [`ARCHITECTURE-SPINE.md`](file:///home/artem/git/caddy-orchestrator/_bmad-output/planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md)  
**Reviewer Role:** System Adversary & Edge-Case Vulnerability Analyst  
**Date:** 2026-10-17  
**Verdict:** **Critical Issues** (Architecture contains fatal composability gaps, entity identity contradictions, and security/state synchronization catch-22s that will cause downstream implementation failure).

---

## Executive Summary

The `ARCHITECTURE-SPINE.md` establishes a compelling high-level vision: Hexagonal architecture, Caddyfile-first configuration authority, staged draft pipelines, and multi-vantage verification.

However, when attacked from an adversarial perspective, the spine exhibits **critical architectural incompatibilities one level down**. Two independent engineering teams implementing adjacent units according to the Architectural Decisions (ADs) will produce components that strictly satisfy every rule yet fail catastrophically upon integration.

The three primary structural contradictions are:
1. **The "Phantom Route" Identity Paradox (AD-1, AD-7 vs ERD & API Conventions):** The spine declares Caddyfile text as the sole authoritative truth (stored as a raw text blob in the `DRAFT` table), yet requires a RESTful API `/api/routes` with persistent lowercase UUIDv4 identifiers and links probe results to `route_id`. Because standard Caddyfiles have no concept of persistent UUIDs, any Caddyfile re-parse destroys route identity, invalidating frontend state signals and orphaning historical verification records.
2. **The Remote Adapt / Secret Injection Catch-22 (AD-3, AD-5, AD-9):** Remote compilation (`POST /adapt?validate=true`) executed on an external Caddy container cannot resolve orchestrator-local environment variables. Injecting plaintext secrets before transmission exposes them over HTTP and embeds them into Caddy's `/config/` JSON, which is then captured unencrypted in SQLite `RECOVERY_POINT.full_caddy_json`, directly violating AD-5's zero-trust AES-GCM encryption mandate.
3. **The Non-Deterministic Hash Lockout Loop (AD-3, AD-8, AD-10):** Computing SHA-256 over raw `GET /config/` responses without a canonicalization standard (such as RFC 8785) causes constant false-positive 409 Conflict drift rejections due to JSON key ordering, Caddy runtime certificate state mutations, and internal defaults.

---

## Downstream Unit Clashes (Concrete Composition Breakdowns)

### Unit Clash 1: Backend CST Parser / Draft Manager vs Frontend State Signal & Verification Subsystem

#### The Setup
- **Unit A (`internal/domain/cst` & `internal/app/draft`):** Implements AD-1 ("canonical master source of truth... is the formatted Caddyfile text") and the ERD (`DRAFT` schema contains only `caddyfile_content text`). Unit A treats Caddyfile AST/CST nodes as the underlying domain model. To satisfy the Consistency Conventions (`/api/routes`), Unit A exposes route projections on the fly. Because standard Caddyfile grammar does not include UUIDs, Unit A generates transient synthetic IDs (e.g. hash of domain + path or ephemeral UUIDv4) upon parsing.
- **Unit B (`web/src/lib/state.ts` & `internal/app/verification`):** Implements AD-4 (Verification), AD-7 (Centralized Signals), and Consistency Conventions ("Entity Identifiers: UUIDv4 formatted as lowercase string"). Unit B maintains UI route state signals keyed by `route.id`, binds route edit forms (`PUT /api/routes/:id`), and records probe executions into `VERIFICATION_RESULT` using `route_id`.

#### The Failure Scenario
1. Operator loads the UI. Unit A parses `caddyfile_content` and serves 5 routes with transient IDs `[R1, R2, R3, R4, R5]`.
2. Unit B launches async verification probes for route `R1` (`app.homelab.local`). The probe succeeds and writes `VERIFICATION_RESULT { route_id: "R1", passed: true }`.
3. Operator switches to `CaddyfileEditor.tsx` and adds a shared snippet or reorders blocks:
   ```caddyfile
   (security_headers) {
       header Strict-Transport-Security "max-age=31536000"
   }
   app.homelab.local {
       import security_headers
       reverse_proxy 192.168.1.50:8080
   }
   ```
4. Draft is saved to SQLite `DRAFT.caddyfile_content`. Unit A re-parses the Caddyfile. The snippet block shifts token offsets; Unit A generates a fresh set of routes with new transient IDs `[R_new_1, R_new_2, ...]`.
5. **The Crash:**
   - The frontend's Preact Signals store retains references to `R1`, but the backend now identifies `app.homelab.local` as `R_new_1`. Form submissions targeting `/api/routes/R1` fail with HTTP 404.
   - The `VERIFICATION_RESULT` table in SQLite retains `route_id: "R1"`. `R1` no longer exists. Historical evidence cards and status badges on the dashboard permanently revert to `unverified` because the link between the route and its verification history has been obliterated.
   - Both units strictly complied with the spine: Unit A kept raw Caddyfile text as the sole truth (AD-1), and Unit B enforced UUIDv4 REST entities and centralized signals (AD-7). Yet they cannot function together.

---

### Unit Clash 2: Secret Vault & Remote Caddy Client vs Snapshot Storage

#### The Setup
- **Unit A (`internal/adapters/secrets` & `internal/domain/cst`):** Implements AD-5 ("All DNS provider credentials... encrypted at rest in SQLite using AES-GCM-256... Secret values are masked with `***` in all HTTP API responses, diff reviews, snapshot exports, and logs"). To avoid leaking secrets in Caddyfile text, Unit A writes environment placeholder tokens (e.g. `api_token {env.CLOUDFLARE_API_TOKEN}`) or injects them dynamically into HTTP payloads.
- **Unit B (`internal/adapters/caddy` & `internal/app/apply`):** Implements AD-3 (Pre-Write Snapshots, Remote Pre-check) and AD-9 (Remote Compilation via Caddy Admin API `POST /adapt`). Unit B runs in a separate orchestrator container and communicates with external Caddy over HTTP (`http://caddy:2019`).

#### The Failure Scenario
1. Operator configures a wildcard domain `*.homelab.net` using Cloudflare DNS-01. The token is saved in `SECRET_CREDENTIAL` encrypted with AES-GCM.
2. Operator clicks "Apply Changes".
3. **Branch 1 (Compilation Failure):** Unit B sends the draft Caddyfile (`POST /adapt?validate=true`) containing `{env.CLOUDFLARE_API_TOKEN}` to external Caddy. External Caddy runs in an independent container/host environment where `CLOUDFLARE_API_TOKEN` is NOT set. Caddy's validation fails with `error adapting config: environment variable not found: CLOUDFLARE_API_TOKEN`. Apply is aborted.
4. **Branch 2 (Plaintext Security Breach):** To fix Branch 1, Unit B replaces `{env.CLOUDFLARE_API_TOKEN}` with the decrypted raw secret token before calling `POST /adapt`. Caddy compiles it successfully to runtime JSON. Caddy's `/load` succeeds.
5. In accordance with AD-3 Step 3 ("A full-configuration snapshot containing the active Caddy JSON... is committed to SQLite before any network call"), Unit B queries Caddy `GET /config/` or uses the adapted JSON and inserts it into SQLite `RECOVERY_POINT.full_caddy_json`.
6. **The Breach:** SQLite now contains the plaintext DNS-01 API token in the `RECOVERY_POINT` table (`full_caddy_json` column). Anyone reading database backups or using the `/api/snapshots` export receives the raw Cloudflare secret, directly violating AD-5.
7. Furthermore, if Caddy returns a validation error on `/load`, Caddy's RFC 7807 problem detail response frequently prints the failing config block, echoing the plaintext secret into orchestrator application logs.

---

### Unit Clash 3: Apply Coordinator Drift Engine vs Caddy Admin Live State

#### The Setup
- **Unit A (`internal/app/apply` & `internal/app/diff`):** Implements AD-3 Step 2 and AD-8 ("queries Caddy Admin API (`GET /config/`), computes its SHA-256 hash (or reads the ETag), and compares it to `base_caddy_hash`. If the hashes differ: 1. The apply operation is aborted with HTTP 409 Conflict").
- **Unit B (`internal/adapters/caddy`):** Implements the HTTP client for Caddy Admin API (`GET /config/`). External Caddy is running official `caddy:2.9` with automatic HTTPS enabled.

#### The Failure Scenario
1. Draft is created at 10:00:00 UTC. Unit B calls `GET /config/`, receives JSON, hashes it to `H0`, and records `DRAFT.base_caddy_hash = H0`.
2. Operator prepares a route and clicks "Apply" at 10:02:00 UTC.
3. Unit A calls `GET /config/` to verify drift.
4. In the intervening 2 minutes:
   - External Caddy rotated an internal certificate or renewed a Let's Encrypt challenge in its internal PKI memory.
   - Or, Caddy's Go runtime serialized internal struct map keys in a non-deterministic order (standard Go `json.Marshal` on `map[string]any` sorts keys, but nested anonymous structs or custom extension modules inside Caddy do not guarantee byte-for-byte canonical JSON).
   - Or, Caddy updated an internal dynamic state field under `apps.tls` or `apps.pki`.
5. The SHA-256 hash computed at 10:02:00 UTC is `H1` (`H1 != H0`).
6. Unit A aborts with HTTP 409 Conflict and prompts the operator to "Rebase Draft".
7. Operator clicks "Rebase". `base_caddy_hash` is updated to `H1`. Operator immediately clicks "Apply" again.
8. Caddy's JSON byte serialization shifts again or another internal state field mutates. Apply aborts with HTTP 409 Conflict again.
9. The orchestrator is trapped in a permanent drift lockout loop where valid user configurations can never be committed.

---

## Comprehensive Hole, Ambiguity, and Loophole Inventory across ADs

### AD-1: Caddyfile-First Authoritative Configuration Model
- **Lossless CST Trivia Deletion Gap:** When a route is removed via visual forms, how does the CST engine decide which surrounding comments belong to that route vs the general file? The spine provides no comment association rules (leading trivia vs trailing trivia). Deleting a route will either orphan trailing comments or delete comments belonging to subsequent routes.
- **Snippet & Global Block Blindspot:** AD-1 states visual forms operate directly on CST nodes. But visual forms only edit routes, proxies, redirects, and file servers. If the user defines reusable snippets (e.g. `(cors) { ... }`) or global options (`{ email admin@example.com }`), where do these live in the CST? If a visual form updates a route that imports `(cors)`, does the visual form mutate the snippet (breaking other routes) or detach the snippet?
- **Syntax Error Trapping in Editor:** AD-1 says if Caddyfile editor input contains syntax errors, "switching back to visual forms is blocked with an error banner, and unparsed text is safely preserved in the draft buffer." What if the user closes the browser tab while in this state? The draft table only has `caddyfile_content text`. Does SQLite store invalid syntax? If SQLite stores invalid syntax, any other device polling via SSE (AD-7) will attempt to parse it and crash their visual forms.

### AD-2: Progressive Root-Level Configuration Ownership
- **Merge Collision Mechanics in `apps.http`:** Caddy's runtime JSON wraps HTTP routes under `apps.http.servers[server_name].routes`. AD-2 states: "merges the preserved `passthrough_apps` JSON at the root level".
  - If `passthrough_apps` contains non-HTTP apps (`logging`, `storage`, `pki`), root-level key merging works.
  - BUT what if adopted Caddy already contained a server listening on `:80`/`:443` with custom sub-routes or layer 4 TCP apps?
  - `POST /adapt` generates a new `apps.http` block. If `passthrough_apps` also contains `apps.http`, a shallow merge overwrites either the managed routes or the passthrough HTTP configuration. The spine lacks a deep merge and conflict resolution specification for shared Caddy applications.
- **Remote Caddyfile Access Paradox:** AD-2 states: "If an existing Caddyfile is mounted/accessible at `/etc/caddy/Caddyfile`, it is adopted directly as the baseline text." In standard homelab deployments (`docker-compose.prod.yml`), the orchestrator runs in container A, while Caddy runs in container B. The orchestrator cannot access `/etc/caddy/Caddyfile` unless the user manually bind-mounts it into both containers. If it is NOT mounted, adoption falls back to synthesizing a Caddyfile from JSON, permanently losing all comments, snippets, and formatting from the user's original Caddyfile on day one.

### AD-3: Staged Draft-and-Apply Pipeline
- **Draft Versioning & Concurrent Edit Void:** The `DRAFT` table schema has no `version` or `sequence_id` column. If two browser sessions (desktop and mobile) submit edits via REST, the database executes blind Last-Write-Wins (LWW) based on timestamps. There is no optimistic locking or conflict detection on draft updates, only on Apply against live Caddy.
- **Rollback Semantic Gap:** "ApplyFailed: Caddy /load returns HTTP 4xx/5xx -> Rollback to Draft with Errors". What does "rollback to draft with errors" mean? Was the pre-write snapshot restored to Caddy? If `/load` failed with 4xx/5xx, Caddy aborted the load internally and kept its previous running config. But what state is the orchestrator's draft in? Is it marked `error`, `draft_saved`, or `editing`?

### AD-4: Multi-Vantage Asynchronous Verification Subsystem
- **Browser Fetch Sandbox Impossibility:** AD-4 mandates a `client_perspective` observation vantage point executed via the operator browser. Modern web browser sandboxes strictly enforce W3C Web APIs:
  - Web browsers **cannot** inspect TLS certificates, SANs, or expiry dates via JavaScript.
  - Web browsers **cannot** execute arbitrary HTTP probes to backend services lacking CORS headers (`Access-Control-Allow-Origin`).
  - Web browsers **cannot** bypass self-signed certificate warnings programmatically to report TLS status.
  - Stating that the client perspective verifies TLS certificates or checks non-CORS private internal routes from the browser is physically impossible in standard web browsers.
- **Uncorrelated Verification Results:** The `VERIFICATION_RESULT` schema lacks any foreign key to `RECOVERY_POINT` or a `caddy_hash`. If an operator applies config A (hash `HA`), then immediately applies config B (hash `HB`), async probes from apply A will write results alongside apply B without any indication of which configuration version was actually verified.

### AD-5: Zero-Trust Secret and Credential Handling
- **Remote Caddy Plugin Verification Hole:** AD-5 mandates: "Before saving or applying a DNS-01 route configuration, the orchestrator verifies via `POST /adapt` dry-run or Caddy module probe that the external Caddy instance actually contains the required DNS challenge module."
  - How is a "Caddy module probe" executed? Caddy Admin API has `GET /config/` and `POST /adapt`, but does NOT have a public `/modules` endpoint in default builds unless querying `/config/` or parsing `/adapt` error strings.
  - If the check relies on `POST /adapt`, Caddy only fails with "unknown module" if the module name is syntactically evaluated. If credentials are missing, it fails on credential syntax first, masking whether the module is present.

### AD-6: Decoupled Read-Only Docker Container Discovery
- **Network Resolution Ambiguity:** `GET /containers/json` returns container private IP addresses (e.g. `172.18.0.4`) and host port bindings (e.g. `0.0.0.0:8080`).
  - In a homelab bridge network, Caddy container might not be able to reach `172.18.0.4` if Caddy is on a different Docker bridge network or on host networking.
  - The discovery adapter does not specify whether suggested upstream addresses use container names, internal bridge IPs, Docker DNS aliases, or host IP + published ports. Selecting the wrong target creates unroutable reverse proxy upstreams.

### AD-7: Centralized Frontend State via Preact Signals
- **Unidirectional SSE Synchronization Hole:** AD-7 says: "Edits initiated on desktop are synchronized via Server-Sent Events (SSE) to mobile browser sessions".
  - SSE is strictly server-to-client streaming.
  - When the desktop user types into an input, how does the state reach the server? Via HTTP `POST/PUT`?
  - If every keystroke emits an HTTP request, the backend is flooded with SQLite write transactions. If debounced, mobile sessions lag.
  - If mobile and desktop users both edit simultaneously, SSE will continuously overwrite their active form inputs, causing cursor jumping and lost input.

### AD-8 & AD-10: Optimistic Concurrency and Ambiguous Apply States
- **ETag / Hash Calculation Specification Gap:** AD-8 allows either SHA-256 computation of `GET /config/` or reading the `ETag` header. Caddy standard Admin API does not send an `ETag` HTTP response header on `GET /config/`. If the implementer relies on SHA-256, they must implement RFC 8785 (JSON Canonicalization Scheme) to strip non-deterministic formatting; otherwise drift detection fails 100% of the time.
- **Deadlock in Reconciliation State:** When `POST /load` times out, draft transitions to `uncertain / reconciliation_required`. AD-10 states: "All further mutations and automated applies are locked until the operator triggers a reconciliation probe".
  - What if the external Caddy instance crashed or the admin port is unresponsive?
  - The operator cannot edit, cannot discard draft, cannot apply, and cannot rollback because all mutations are locked. There is no emergency unlock or force-reset mechanism specified.

### AD-11: Recovery Point Retention and Backup Package Specification
- **FIFO Eviction of Critical Baselines:** AD-11 mandates a maximum of 100 recovery points with "automatic FIFO cleanup on each successful apply".
  - If an operator executes 100 small draft applies (e.g. tweaking headers), the original **Initial Adoption Baseline Recovery Point** (created in AD-2) is automatically deleted from SQLite!
  - Once the baseline is pruned, full disaster recovery to the pre-orchestrator state is permanently impossible. There is no concept of "pinned" or "named" snapshots exempt from FIFO pruning.

---

## Critical Data Model & API Contradictions

| Spine Section / Requirement | Schema / API Definition in Spine | Contradiction / Fatal Flaw |
| :--- | :--- | :--- |
| **Route Identity** (`AD-1`, Consistency Conventions) | Consistency Conventions: `/api/routes` with UUIDv4; `VERIFICATION_RESULT` has `route_id string`. | The database schema has **no `ROUTE` table**. The `DRAFT` table stores only raw `caddyfile_content text`. Caddyfile format has no UUID syntax. Route IDs are purely synthetic and ephemeral, breaking all foreign keys upon text edits. |
| **Snapshot Secrets** (`AD-3`, `AD-5`) | `RECOVERY_POINT` table has `full_caddy_json text` in plaintext SQLite. | AD-5 forbids secrets in snapshots/logs and mandates AES-GCM encryption. Caddy runtime JSON requires plain API tokens for DNS-01, meaning `full_caddy_json` stores raw credentials in unencrypted text. |
| **Draft Concurrency** (`AD-3`, `AD-7`) | `DRAFT` table schema: `id PK`, `status`, `base_caddy_hash`, `caddyfile_content`, `updated_at`. | No `version` or revision column. Concurrent edits from multiple tabs or devices will silently clobber each other via unversioned HTTP PUTs. |
| **Verification Traceability** (`AD-4`) | `VERIFICATION_RESULT` schema: `id`, `route_id`, `domain`, `check_type`, `passed`, `evidence_json`. | No `caddy_hash`, `draft_id`, or `recovery_point_id`. Probe results cannot be attributed to any specific configuration deployment. |

---

## Actionable Recommendations for Spine Hardening

1. **Establish a Deterministic Route Identity Protocol in the CST:**
   - Either introduce a lightweight, comment-based metadata pragma inside the Caddyfile that the CST engine preserves (e.g. `# @id: 550e8400-e29b-41d4-a716-446655440000` placed directly above site blocks), OR define route identity strictly by canonical host matchers (`domain:port`), completely removing ephemeral UUIDv4 keys from the domain route model.
2. **Standardize on RFC 8785 Canonical JSON for Live Drift Detection:**
   - Explicitly mandate RFC 8785 (JSON Canonicalization Scheme) for hashing `GET /config/`. Define a filter that strips dynamic runtime state keys (`apps.pki`, active certificate statuses) prior to computing `base_caddy_hash`.
3. **Resolve the Remote Adapt / Secret Injection Architecture:**
   - Clarify how secrets reach external Caddy. If external Caddy is running in Docker, specify that secrets are written to an encrypted environment file or passed via Caddy's dynamic JSON API, and mandate that `RECOVERY_POINT.full_caddy_json` must undergo secret masking/redaction before being stored in SQLite.
4. **Scope Browser-Perspective Verification Realistically:**
   - Remove raw TLS certificate inspection and non-CORS header checks from the `client_perspective` probe definition. Restrict `client_perspective` strictly to public HTTPS GET requests to verify end-to-end reachability from the operator's browser, leaving TLS certificate chain analysis solely to `manager_network` probers.
5. **Protect the Baseline Recovery Point from FIFO Pruning:**
   - Modify the `RECOVERY_POINT` schema to include `is_pinned boolean DEFAULT 0`. Pin the initial adoption baseline and user-tagged snapshots so FIFO eviction only purges unpinned intermediate apply states.
6. **Add Draft Optimistic Locking:**
   - Add `version integer NOT NULL DEFAULT 1` to the `DRAFT` table and require `If-Match: "version"` on all draft mutation API calls to prevent multi-device race conditions over SSE.
