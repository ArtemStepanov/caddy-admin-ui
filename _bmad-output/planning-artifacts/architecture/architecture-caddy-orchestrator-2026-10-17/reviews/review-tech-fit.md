> Historical review of the pre-correction spine. Its original verdict/claims are retained as provenance. Current dispositions and fresh revision-1 review links are in [review-dispositions](../review-dispositions.md); this old verdict is not the current readiness gate.

# Technology Fit & Reality-Check Review — Architecture Spine

- **Target Artifact:** `_bmad-output/planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md`
- **Reviewer:** Technology Fit & Reality-Check Reviewer (Independent Subagent)
- **Evaluation Date:** 2026-10-17
- **Verification Method:** Live container execution (Caddy 2.9.1, Docker 29.9, Go 1.25.13 toolchain, Alpine Linux manifests) and AST/dependency verification
- **Verdict:** **Minor issues** (Technology selections are exceptionally well-grounded and compatible; 3 specific runtime API & dependency gaps require targeted correction)

---

## 1. Executive Summary

Every technology, library version, protocol endpoint, and operational capability committed in `ARCHITECTURE-SPINE.md` was subjected to concrete empirical reality-checking against live runtimes and official specifications:

1. **Caddy 2.9 Admin API:** Tested directly against official `caddy:2.9` (v2.9.1).
   - `POST /adapt` with `Content-Type: text/caddyfile` is fully supported, parses Caddyfile to runtime JSON, and accurately catches missing plugin modules (e.g. `dns.providers.cloudflare`).
   - `GET /config/` reliably emits `Etag: "/config/ <hash>"`.
   - Conditional writes with `If-Match` work on `/config/...` endpoints (returning `412 Precondition Failed` on mismatch).
   - However, **`POST /load` ignores `If-Match` headers entirely**. Optimistic concurrency drift protection relies strictly on the orchestrator's pre-flight check or using `POST /config/`.
   - The PKI endpoint referenced in the spine (`/pki/authorities/local/root.crt`) does not exist in Caddy (returns 404). Caddy's actual endpoints are `GET /pki/ca/<id>` and `GET /pki/ca/<id>/certificates`.
2. **Go & Gin Backend:** Go 1.25.13 and Gin 1.12.0 are completely compatible and functional. Full test suite passes cleanly with race detector on `linux/amd64`. Official Docker image `golang:1.25.13-alpine` is available on Docker Hub.
3. **SQLite WAL Mode:** Verified with `mattn/go-sqlite3 v1.14.50`. Query parameter configuration `_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on` works seamlessly with single-writer connection pooling.
4. **Frontend Dependencies:** Preact 10.29, Preact Router 4.1.2, Tailwind CSS 3.4.1, and Vite 8.2.2 are verified and compatible. However, **`@preact/signals` is missing from `web/package.json` and omitted from the spine's Stack table**, despite being a core requirement of AD-7.

---

## 2. Reality-Check Evidence Matrix

| Technology / Capability | Spine Specification | Live Runtime / Empirical Result | Status | Notes / Adjustments Needed |
| :--- | :--- | :--- | :--- | :--- |
| **Caddy Target Engine** | 2.9.x | `v2.9.1` (tested in container) | **Pass** | Clean baseline compatibility. |
| **Caddy `POST /adapt`** | `Content-Type: text/caddyfile`, `?validate=true` | Returns `200 OK` with adapted JSON and warnings; returns `400 Bad Request` with syntax/module errors | **Pass** | Verified that missing DNS modules trigger descriptive errors (`module not registered: dns.providers.cloudflare`). Caddy requires newline after `{` in Caddyfile blocks. |
| **Caddy `GET /config/` ETag** | Returns ETag for concurrency check | Emits `Etag: "/config/ <hash>"` | **Pass** | Verified. Returns 64-bit hexadecimal hash. |
| **Caddy `If-Match` Enforcement** | Optimistic concurrency locking on commit | **Supported on `/config/...`**, but **Ignored on `POST /load`** | **Caution / Minor Issue** | `POST /load` does not inspect `If-Match`. The orchestrator must enforce optimistic locking in Go application code prior to `/load`, or use `POST /config/` for atomic replacement. |
| **Caddy PKI Endpoints** | `/pki/authorities/local/root.crt` (AD-4) | **404 Not Found** for `/pki/authorities/...` | **Minor Issue** | Real endpoints: `GET /pki/ca/<id>` (JSON with PEM strings) and `GET /pki/ca/<id>/certificates` (`application/pem-certificate-chain`). |
| **Go Toolchain** | 1.25.13 | `go version go1.25.13 linux/amd64` | **Pass** | Verified via `mise.toml` and clean execution of `go test ./...`. |
| **Gin Framework** | 1.12.0 | `github.com/gin-gonic/gin v1.12.0` | **Pass** | Verified in `go.mod` and active handlers. |
| **SQLite WAL Mode** | `PRAGMA journal_mode=WAL;`, busy timeout 5000ms | Supported via DSN `_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on` | **Pass** | Tested with `mattn/go-sqlite3 v1.14.50`, single connection pool (`SetMaxOpenConns(1)`). |
| **Alpine Linux Base** | 3.24 | `alpine:3.24` / `golang:1.25.13-alpine` | **Pass** | Verified multi-arch OCI manifests on Docker Hub. |
| **Preact** | 10.29.8 | `preact: ^10.29.8` (10.29.3 installed) | **Pass** | Verified in `web/package.json`. |
| **Preact Signals** | AD-7 core state management | **Missing in `web/package.json`** | **Minor Issue** | `@preact/signals` must be added to `web/package.json` dependencies and listed in the Spine Stack table. |
| **Tailwind CSS** | 3.4.1 | `tailwindcss: ^3.4.1` with PostCSS 8 | **Pass** | Standard PostCSS setup in `web/package.json`. |
| **Vite** | 8.2.2 | `vite: ^8.2.2` | **Pass** | Verified in `web/package.json`. |

---

## 3. Detailed Findings & Reality-Check Analysis

### Finding 1: Caddy PKI Endpoint URI Mismatch in AD-4

- **Location in Spine:** Section *AD-4 — Multi-Vantage Asynchronous Verification Subsystem*, line 180:
  > *"Internal CA Guidance: For local CA (`tls internal`), the orchestrator provides an endpoint to download Caddy's root CA certificate (`/pki/authorities/local/root.crt`) and displays client trust installation guidance."*
- **Empirical Reality Check:**
  Invoking `curl http://127.0.0.1:2019/pki/authorities/local/root.crt` against Caddy 2.9.1 returns:
  ```http
  HTTP/1.1 404 Not Found
  Content-Type: application/json
  
  {"error":"resource not found: /pki/authorities/local/root.crt"}
  ```
  In Caddy 2.x, the PKI app registers the following endpoints:
  1. `GET /pki/ca/<id>` (e.g. `/pki/ca/local`): returns a JSON object with fields:
     - `id`: `"local"`
     - `name`: `"Caddy Local Authority"`
     - `root_certificate`: `"-----BEGIN CERTIFICATE-----\n..."`
     - `intermediate_certificate`: `"-----BEGIN CERTIFICATE-----\n..."`
  2. `GET /pki/ca/<id>/certificates` (e.g. `/pki/ca/local/certificates`): returns `Content-Type: application/pem-certificate-chain` containing the raw PEM certificate chain directly.
- **Impact:** Any client or adapter expecting `/pki/authorities/local/root.crt` directly from Caddy will fail with a 404 error.
- **Remediation:** Update line 180 of `ARCHITECTURE-SPINE.md` to specify:
  - The orchestrator fetches the root certificate from Caddy via `GET /pki/ca/local` (extracting `root_certificate`) or `GET /pki/ca/local/certificates`.
  - The orchestrator exposes its own clean UI download endpoint (e.g. `GET /api/pki/local/root.crt` with `Content-Disposition: attachment; filename="caddy-local-root.crt"` and `Content-Type: application/x-x509-ca-cert`).

---

### Finding 2: Missing `@preact/signals` in Frontend Stack and `package.json`

- **Location in Spine:** 
  - Primary Adapters Diagram (line 41): `UI["Web SPA (Preact + Signals)"]`
  - Layer Boundaries (line 92): `Preact, TypeScript, Tailwind CSS, and Preact Signals`
  - AD-7 (lines 201–208): `Centralized Frontend State via Preact Signals [ADOPTED]`
  - Stack Table (lines 261–275): Lists `Preact 10.29.8`, but omits `@preact/signals`.
  - Directory Structure (line 459): `web/src/lib/state.ts # Centralized Preact Signals store`
- **Empirical Reality Check:**
  Inspecting `web/package.json` shows only:
  ```json
  "dependencies": {
    "preact": "^10.29.8",
    "preact-router": "^4.1.2"
  }
  ```
  Neither `@preact/signals` nor `@preact/signals-core` is listed in `dependencies` or installed in `node_modules`.
- **Impact:** Frontend implementation following AD-7 will immediately encounter compile errors (`Cannot find module '@preact/signals'`) unless installed.
- **Remediation:**
  1. Add `@preact/signals` (e.g. `^2.0.0` or `^1.3.1` compatible with Preact 10.29) to `web/package.json` dependencies.
  2. Add `@preact/signals` to the Spine's **Stack** table under Frontend technologies.

---

### Finding 3: Caddy Admin API `POST /load` Concurrency Enforcement Contract

- **Location in Spine:**
  - Section *AD-3 — Staged Draft-and-Apply Pipeline*, Step 2 & 4 (lines 165–168).
  - Section *AD-8 — Optimistic Concurrency and External Drift Locking*, line 213:
    > *"Before committing changes, the orchestrator queries Caddy Admin API (`GET /config/`), computes its SHA-256 hash (or reads the ETag), and compares it to `base_caddy_hash`. If the hashes differ: 1. The apply operation is aborted with HTTP 409 Conflict."*
- **Empirical Reality Check:**
  Testing conditional HTTP requests against Caddy 2.9.1 reveals a critical nuance in Caddy's HTTP API:
  - `POST /config/` (and subpaths) enforces `If-Match`: sending an invalid or mismatched ETag header (`If-Match: "/config/ 0000000000000000"`) returns `412 Precondition Failed`.
  - However, **`POST /load` does NOT inspect the `If-Match` header**. Sending `If-Match: "/config/ 0000000000000000"` to `POST /load` returns `200 OK` and unconditionally replaces the configuration.
- **Impact:**
  - There is a narrow race window if an external process modifies Caddy *between* the orchestrator's `GET /config/` drift check and its subsequent `POST /load`. Caddy will not reject the `/load` call.
- **Remediation:**
  Clarify in AD-3 and AD-8 that:
  1. Because Caddy's `/load` endpoint is inherently unconditional, optimistic concurrency locking is **orchestrator-enforced** via the pre-flight `GET /config/` check.
  2. Alternatively, if strict protocol-level atomicity is required, the orchestrator can execute root-level replacement via `POST /config/` with `If-Match: "/config/ <hash>"`. However, `POST /load` is preferred for complete app lifecycle reloads; therefore, the spine should explicitly note that the concurrency guard is enforced in `internal/app/apply` before the network call.

---

### Finding 4: Docker Container Network Access to Caddy Admin API (`:2019`)

- **Location in Spine:**
  - System Context Diagram (lines 282–300): Orchestrator communicates with `CaddyServer (Admin API :2019)` via HTTP.
  - Environment defaults: `CADDY_ADMIN_URL=http://localhost:2019`.
- **Empirical Reality Check:**
  By default, Caddy's admin endpoint binds to `localhost:2019` with origin and host enforcement. When Caddy runs in a standard container:
  - Inside the Caddy container, curl to `http://localhost:2019/config/` succeeds.
  - If queried from an external container using a different port or hostname without configuring Caddy's `admin` block, Caddy returns:
    `HTTP/1.1 403 Forbidden: Port 12019 not allowed for HTTP` or blocks the connection due to binding to loopback only.
- **Impact:** In a multi-container Docker Compose deployment where the Orchestrator and Caddy are separate containers, Caddy's default `localhost:2019` binding will refuse connections from the orchestrator container.
- **Remediation:** 
  The deployment architecture (`docker-compose.yml` / `docker-compose.prod.yml`) and AD-2/AD-10 must document the required Caddy admin configuration:
  - Either Caddy must be started with admin listening on all interfaces in its container network:
    ```caddyfile
    {
        admin 0.0.0.0:2019 {
            origins localhost:2019 caddy:2019
        }
    }
    ```
  - Or the orchestrator and Caddy share the network namespace (`network_mode: "service:caddy"` or loopback host networking). This must be made explicit in the infrastructure setup instructions.

---

## 4. Verification of Core Decisions

### CST Caddyfile Model vs. Runtime JSON (AD-1 & AD-9)
- **Verdict:** **Fully Validated.**
- **Details:** Remote compilation via `POST /adapt` with `Content-Type: text/caddyfile` works cleanly and eliminates the need for bundling a Caddy binary in the orchestrator image. Caddy's error reporting on `/adapt` provides exact line numbers and messages, which directly supports the visual editor error banners required by AD-1.

### Root-Level Progressive Ownership & `passthrough_apps` (AD-2)
- **Verdict:** **Fully Validated.**
- **Details:** Caddy allows selective extraction and injection of top-level JSON keys under `/config/apps`. Non-HTTP apps (e.g. `apps.logging`, `apps.pki`, `apps.storage`) can be safely retained as raw JSON while `apps.http` is generated from the adapted Caddyfile.

### SQLite Storage and WAL Concurrency (AD-3, AD-11)
- **Verdict:** **Fully Validated.**
- **Details:** The SQLite DSN `path+"?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on"` with `SetMaxOpenConns(1)` guarantees zero database lock contention during concurrent background probe writes and UI reads.

### Zero-Trust Secrets via AES-GCM-256 (AD-5)
- **Verdict:** **Fully Validated.**
- **Details:** Standard Go standard library `crypto/cipher` and `crypto/aes` provide high-performance, non-external cryptographic primitives. Injection into Caddy during `/adapt` dry-runs avoids token persistence in plain Caddyfile text.

---

## 5. Summary of Recommended Adjustments to Spine

1. **AD-4 (Line 180):** Replace `/pki/authorities/local/root.crt` with Caddy's actual endpoint `GET /pki/ca/local` (and `GET /pki/ca/local/certificates`), noting the orchestrator provides a user-facing download route `GET /api/pki/local/root.crt`.
2. **Stack Table & Frontend Dependencies:** Add `@preact/signals` (`^2.0.0` or `^1.3.1`) to the Stack table in `ARCHITECTURE-SPINE.md` and `web/package.json`.
3. **AD-3 & AD-8:** Clarify that Caddy's `POST /load` is an unconditional full configuration replacement, so optimistic locking is guarded by the orchestrator comparing live `GET /config/` ETag against draft `base_caddy_hash` prior to calling `/load`.
4. **Deployment Contracts:** Document the Caddy admin listener requirement (`admin 0.0.0.0:2019`) or shared container network namespace in the production deployment guidelines.
