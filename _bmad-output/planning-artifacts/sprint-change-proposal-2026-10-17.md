---
project: caddy-orchestrator
date: 2026-10-17
created_at: 2026-10-17T08:26:20Z
workflow: bmad-correct-course
status: approved-applied
review_mode: batch
approval: approved-by-Artem
scope_classification: moderate
technical_risk: high
trigger: implementation-readiness.md — FAIL
documentation_changes_authorized: true
product_implementation_performed: false
approval_date: 2026-10-17
handoff_recipients: [Product Owner, Developer, Architect, QA]
---

# Sprint Change Proposal — Implementation Readiness Corrections

## 1. Issue summary

The [readiness assessment](implementation-readiness.md) found ten blockers (B01–B10) and five consistency concerns (C01–C05). Artem identified that assessment as the trigger for this course correction. These are planning and implementation-contract defects discovered before sprint tracking generation; no in-progress implementation story has been identified. There is no `sprint-status.yaml` to repair yet.

The most consequential examples are a read-before-`/load` race that can overwrite another writer, adoption without a defined preservation model, recovery without durable attempt records, and early stories depending on secrets/CST/snapshot work assigned to later epics. Other discrepancies concern forms coverage, stale session writes, verification evidence, Caddy APIs, migration and backups.

Confirmed direction remains: full existing-configuration management; forms for all settings plus a Caddyfile editor; shared batch drafts and explicit review/apply; all four TLS scenarios; manual destinations and optional Docker discovery; manager-only Compose against external Caddy; Quiet Workshop. This proposal does not reopen those choices. Arbitrary plugin compatibility and original-source reconstruction are not established facts.

Evidence examined: PRD and addendum; full epics; architecture AD-1–AD-11; DESIGN/EXPERIENCE and UX decision record; SPEC and companions; test plans/handoff; legacy storage models/schema; pinned Caddy source and official API documentation. This run performed document/source analysis, not product runtime certification.

## 2. Impact analysis

| Area | Impact |
| --- | --- |
| Epic 1: adoption | Needs source/preservation, encryption, migration, recovery and minimal safe publication prerequisites together. Connection testing cannot invent build/module inventory. |
| Epic 2: route management | CST identity, validation, all retained route features, global/server forms and source compatibility need explicit acceptance. Publication cannot depend on an undefined later pipeline. |
| Epic 3: draft/apply | Needs root conditional mutation, durable attempts, true three-way reconciliation, draft revisions and bounded failure claims. Core safety must move before the first publication. |
| Epic 4: TLS | All four scenarios remain required. Secret infrastructure moves earlier; capability evidence and provisioning errors need accurate contracts. |
| Epic 5: verification | Needs configuration identity, observation/causation separation, explicit expectations and browser limitations. Minimal bounded verification belongs in the first publication slice. |
| Epic 6: recovery/maintenance | Core snapshot/restore storage moves earlier; history UX remains later. Backup and migration must cover actual external assets and locally stored disabled routes/settings. |
| PRD/addendum | Narrow scope is superseded by confirmed decisions; retain FR-1–FR-11 identifiers while updating meanings and acceptance. |
| SPEC/glossary/diagrams | Remove unsupported atomicity/fidelity/API promises; use the same source, revision, apply and recovery model as architecture. |
| UX | Keep visual direction. Specify stale draft, incompatible source, renewed review, uncertain recovery and unavailable observations. |
| QA/reviews | Correct story mappings, API fixtures, readiness labels, invented budgets and review dispositions. |
| Code/operations later | Domain/application/adapter boundaries, DB migrations, API DTOs, SSE, CSP, Compose and runbooks will change during implementation. Current CLAUDE/README baseline descriptions stay distinguishable from proposed behavior. |

There is no evidence justifying a rollback of current code. Current scoped writes remain the working brownfield boundary until the new workflow is implemented and explicitly adopted. Do not introduce an interim unconditional full-config writer.

## 3. Recommended approach

**Direct adjustment with backlog reorganization; moderate change scope and high technical risk.** The architectural direction remains approved. The work is cross-document contract repair and sequencing, requiring Architect, PO/DEV and QA coordination rather than a new product plan.

| Option | Effort / risk | Decision |
| --- | --- | --- |
| Direct adjustment | Medium planning effort; high implementation risk until preservation/concurrency contracts are proved | Recommended. Fix contracts, propagate them, then reslice stories and validation. |
| Rollback | Additional effort and operational risk without resolving the planning defects | Not selected. No completed new-model implementation needs reverting. |
| MVP reduction | Would alter confirmed editing/TLS scope without solving safety | Not selected. Keep capabilities; order delivery into complete safe outcomes. |

**Planning estimate:** approximately 3–6 focused engineering days for technical contract completion, propagation, backlog/test mapping and repeat readiness review by one person. This is a provisional estimate, not measured velocity or a release commitment. Product implementation and acceptance execution are additional and need re-estimation after story slicing. Earlier 74–121-hour test estimates are not a validated estimate for this correction.

**Timeline:** sprint generation moves after correction and repeat readiness PASS. No release date was agreed, so a calendar slip cannot be computed. Stop downstream story implementation that would require inventing one of these contracts; independent document reconciliation can proceed after proposal approval.

## 4. Detailed change proposals

OLD excerpts below identify the current clauses to replace. NEW text is the proposed contract; supporting operational details are part of the replacement. Artem explicitly approved all proposals in this conversation; the before/after text below preserves the reviewed change record.

### P01 — Propagate confirmed scope and record authority (C01, B06)

**Artifacts:** PRD §§1–5, FR-1/3/4/5/8/9/10/11, §8 decisions; addendum; both reconciliation documents; SPEC opening authority statement.

**OLD:** PRD FR-4: “Ordinary managed changes modify only the Managed Server's Routes.” FR-8: “HTTP-only, internal CA and wildcard/DNS-01 are candidates, not all approved modes.” §4 excludes automated Docker discovery and a visual editor for every Caddy module.

**NEW:** “The intended product manages the entire adopted Caddy configuration, including existing routes and global/server settings. Forms and the Caddyfile editor operate on one shared persistent candidate. Saved edits do not publish. Public automatic HTTPS, DNS-01/wildcards, internal CA and HTTPS upstreams are required scenarios. Manual destinations and optional read-only Docker discovery are included. All settings within each declared supported build must have form coverage; compatibility limitations are explicit and cannot silently become permanent editor-only scope.”

Replace FR-1 adoption with source-aware whole-config preview and confirmation without a remote write; FR-4 with the preservation contract in P03; FR-5 with P02/P04/P06; FR-8 with all four scenarios and honest capability evidence; FR-9 with scope-labelled whole-config recovery plus separately retained legacy route snapshots; FR-10 with P09/P10. FR-3/11 explicitly include global/server forms and retained feature/validation coverage in P07/P11. Keep all FR identifiers.

Record authority: latest explicit user choices govern scope; current-code documents govern current behavior; the corrected PRD records those choices; architecture defines technical mechanisms; SPEC summarizes that consistent contract; epics/QA trace it. Unresolved conflicts remain visible and block the affected implementation rather than being settled by a blanket “canonical” label. Update D1/D2 and stale addendum candidate language accordingly. Detailed layout/default assumptions retain their proposal status.

**Rationale / release impact:** documentation reconciliation, not additional MVP scope. Preserve fleet/RBAC/DNS-record/container-lifecycle non-goals.

### P02 — Use remotely guarded root writes and reviewed reconciliation (B01)

**Artifacts:** architecture AD-3/8/11, ports and diagrams; old stories 3.4/3.5/6.3; SPEC CAP-5/9; QA concurrency fixtures.

**OLD:** AD-3: “The merged configuration is posted to Caddy's `/load` endpoint.” Story 3.4 checks an ETag before applying and describes rebase as updating the base ETag while preserving edits.

**NEW:** “Apply and restore submit one complete JSON candidate using `POST /config/`, `Content-Type: application/json`, and the exact root ETag from the reviewed baseline as `If-Match`. The root conditional mutation is the remote concurrency guard. A preceding GET is for review/recovery, not protection by itself. Remote 412 becomes application HTTP 409 conflict, preserving the candidate; there is no fallback to `/load` or unconditional retry.”

Bind a review token to target/adoption generation, draft revision, baseline ETag, candidate identity and preservation/dependency manifest. Any change invalidates that review. Restore binds the same fields to the selected recovery point and live baseline captured during comparison; target mismatch blocks before submission. Require real Caddy 2.9.1 tests for root replacement, empty/null root bootstrap and a competing write inserted between final read and submission. Block unsupported targets rather than weaken the guard.

Rebase compares original baseline B, latest live L and draft candidate C. Nonoverlapping, provably mapped edits may form a proposed merged draft; overlapping edits, deletion-versus-edit and ambiguous array/source identities require explicit resolution. Preserve original inputs, show conflicts and require a new diff review. Updating only the ETag is forbidden. Discard needs a fresh safe source association under P03; if live cannot be represented, retain the draft and show compatibility blocking rather than fabricate Caddyfile.

External writers can still replace configuration *after* a successful guarded write; detect subsequent drift and mark observations stale. This is not an exclusive-writer guarantee or a transaction spanning Caddy and SQLite.

**Rationale:** closes the read/write race while retaining root ownership. The proposed endpoint follows the [official concurrency protocol](https://caddyserver.com/docs/api#concurrent-config-changes); pinned [mutation implementation](https://github.com/caddyserver/caddy/blob/v2.9.1/caddy.go#L132-L244) checks the precondition while holding its config lock. Runtime integration evidence is still required.

### P03 — Make source acquisition and preservation explicit (B02)

**Artifacts:** architecture AD-1/2, schema; old stories 1.3/2.1/2.5; SPEC CAP-1/2; EXPERIENCE F1/F4; adoption tests.

**OLD:** AD-2 synthesizes an initial Caddyfile from runtime JSON and merges `passthrough_apps` at root level; a readable `/etc/caddy/Caddyfile` is assumed available.

**NEW:** “Adoption stores a protected exact runtime baseline and an explicit source bundle, provenance and dependency manifest. Source is operator-supplied or obtained through a configured read-only source path; manager-only network access does not imply access to Caddy's filesystem. Missing source is a source-required state. JSON-only conversion is allowed only for a tested converter whose output passes the same equivalence/preservation checks; it never recreates original comments/imports.”

Use this concrete composition model:

1. B is the live root JSON, E its root ETag, S the exact supplied source bundle, and G0 the target-adapted result for S. Store all identities and source/build/dependency provenance. Adoption preview binds them; refetch live and reject a changed preview before local confirmation. Adoption performs no configuration write, including on empty Caddy.
2. Produce a path ownership manifest: generated/mapped subtrees versus preserved raw subtrees. It includes nested routes, server/listener settings, other servers, root `logging`/`storage`, and `apps.pki`; these are not all applications. Unknown/overlapping ownership or unexplained B/G0 differences block editing/publication of the affected configuration.
3. For candidate source S1, adapt to G1 and compute intended generated changes G0→G1. Apply only mapped changes onto B. Copy preserved subtrees without reconstruction; preserve array order. A change to a parent containing preserved descendants, route regrouping or a server-name remap needs a proven mapping and preservation test; otherwise block with the exact collision. Never use a generic recursive merge or replace the entire HTTP app by default.
4. Diff exposes both source edits and the resulting root-JSON delta, including preserved scope. The review identifies every changed runtime path. Source and semantic equivalence use versioned comparison rules under P04, not text equality of JSON.
5. External imports require a bundle manifest with ordered resolved file paths, exact contents/digests, expansion context and ownership. Adapter-visible dependencies must be reproduced on target and frozen for the compilation/review window, or supplied as a deterministic compiled bundle with source mapping. If dependency freshness/immutability cannot be established, block publication. No remote filesystem-write mechanism is implied. Keep originals separately; do not claim expansion preserves source formatting.
6. Source syntax errors, missing dependencies, unrepresentable runtime JSON and unknown module handlers retain protected source/runtime copies and actionable compatibility states. They do not silently narrow the requested all-settings goal. Changing targets detaches and retains prior drafts/history; it never binds them to the new target or silently clears unsaved work.

**Rationale:** full-config ownership needs exact preservation below the root, not just top-level passthrough. [Caddyfile limitations](https://caddyserver.com/docs/caddyfile#note) preclude assuming arbitrary inverse reconstruction. The manifest/mapping design is a proposed manager mechanism requiring corpus tests and Architect validation.

### P04 — Persist attempts and reconcile content after interruption (B03)

**Artifacts:** architecture AD-3/10, state machine/schema; old stories 3.5/6.3; SPEC CAP-8; fault-injection tests.

**OLD:** mark uncertain on network failure; compare a live ETag to a candidate hash to decide whether apply succeeded.

**NEW:** “Before submission, one local transaction commits the pre-write recovery point and a durable apply/restore attempt: attempt ID, operation, target/adoption generation, draft revision, reviewed ETag, protected baseline and exact candidate, comparison version and status `prepared`. No remote mutation occurs if this transaction fails. Commit `submission_possible` before sending; startup treats every unfinished attempt as reconciliation-required.”

Successful response alone does not finalize local state. Read back root state, compare content, and atomically record the observed applied revision plus finalized attempt. Response loss, readback failure, crash or failed local finalization keeps the attempt unresolved and blocks further mutations. If even recording uncertainty fails, fail closed in memory; the precommitted unfinished attempt reinstates the lock on restart. Never retry the mutation automatically.

Content comparison: parse JSON with exact number handling; compare objects independently of key order, arrays in order, numbers by exact numeric value, and missing separately from null. Do not remove unknown fields/defaults or reorder arrays. Version the algorithm and golden fixtures. If Caddy normalization prevents a proven comparison, outcome remains unresolved. Do not compare opaque Caddy ETags to manager hashes. Protected comparison material includes hydrated historical secret versions; equality digests exposed locally use a keyed digest rather than a secret oracle.

| Reconciliation observation | Result |
| --- | --- |
| Live equals candidate | Record candidate currently present, finalize local baseline, enqueue verification; lost response does not prove historical causation. |
| Live equals baseline and differs from candidate | Record baseline currently present, preserve unapplied candidate for renewed review; no assertion that it was never temporarily active. |
| Live differs from both | Conflict; retain attempt/draft and require explicit resolution and renewed review. |
| Live unavailable / comparison or secret material unavailable | Remain reconciliation-required; explain missing evidence. |

If candidate equals baseline, treat as a no-op and report current equality without inferring a unique remote outcome. Authoritative rejection/precondition responses are distinct from arbitrary proxy/server errors; ambiguous failures still require readback.

**Rationale:** covers process death and local failure after acceptance, not just transport errors. Restore uses precisely the same attempt protocol.

### P05 — Protect every source representation and historical secret version (B04)

**Artifacts:** architecture AD-5/11 and schema; old stories 1.3/4.1/4.3/6.1/6.3; SPEC security constraint; secret tests.

**OLD:** only recovery JSON has `secret_ref:<id>`; all API secret fields are masked with `***`.

**NEW:** “Initialize the encryption key and vault before importing any baseline. Protect exact source bundles, runtime originals, legacy archives, draft buffers and recovery payloads using authenticated encrypted envelopes; extracted known secret fields use typed references to immutable secret versions. Unknown opaque content is encrypted and withheld from ordinary API/editor output until its display/redaction mapping is proven. No secret-bearing source text is stored as plaintext merely because JSON was sanitized.”

Keep immutable secret versions while referenced by a retained draft, attempt or recovery point. Rotation creates a new version; deleting a current credential cannot destroy retained recovery material. Missing key/version blocks adapt/apply/restore with an explicit error. Garbage collection follows reference retention, including unresolved attempts.

Use API actions `keep(reference, version)`, `replace(value)` and `remove`; absence and display masks never mean replacement. The editor shows stable typed reference tokens, not masked literal credentials; validate tokens against server-side ownership/version metadata. Materialize secrets in memory for remote adaptation/submission only. Sanitize remote error text, warnings, logs, diffs and exports. Unknown secret-bearing blocks require safe opaque handling rather than a false universal redaction claim. Ordinary exports are sanitized and labelled incomplete for standalone secret recovery; full protected backup includes key material by the documented procedure.

**Rationale:** makes adoption safe before DNS feature UI and enables historical restore after rotation. Encryption protects storage; redaction separately governs display.

### P06 — Version shared drafts and asynchronous observations (B07)

**Artifacts:** architecture AD-7, DRAFT/VERIFICATION_RESULT schema, REST/SSE contracts; old stories 3.1/5.1/5.4; EXPERIENCE F7.

**OLD:** SSE synchronizes the singleton draft; probe rows contain route/time/vantage but no applied revision.

**NEW:** “Every persisted draft mutation includes the expected draft revision. SQLite atomically increments it only on a matching revision; stale writes return HTTP 409 with latest revision and preserve the client's unsaved buffer. SSE carries target generation/revision and monotonic event identity. Reconnect fetches authoritative state and offers explicit reconciliation for unsaved input; incoming events never overwrite it.”

Add target ID/adoption generation, draft revision, applied configuration revision, attempt ID, probe job ID, expectation version and source/dependency manifest identity where relevant. Jobs/results bind to the exact applied revision and target. Late results remain historical, visibly stale, and cannot mark the current revision verified. Restore creates a new applied revision even when content matches older history. Target changes cancel/detach jobs and invalidate review tokens.

**Rationale:** one operator with desktop and phone can still issue conflicting saves; SSE notification is not write concurrency control.

### P07 — Define complete form and brownfield regression coverage (B06)

**Artifacts:** architecture AD-1/9 and compatibility companion; PRD FR-3/11; old story 2.5 plus a new global/server forms story; SPEC CAP-2/3.

**OLD:** “Directives lacking dedicated visual form fields … [are] read-only ‘Editor-managed’ blocks.” No global/server editing story exists.

**NEW:** “Maintain a versioned supported-build/settings inventory mapping every setting to form representation, Caddyfile/CST representation, runtime path, validation, secret classification and round-trip test. A required setting without a form is an implementation gap, not accepted editor-only completion. Expose compatibility gaps and block destructive conversion. Preserve Caddyfile as the chosen peer editor.”

New story: **Edit existing global and server settings through forms and Caddyfile.** Acceptance: existing global/server values appear in forms; form edits and text edits produce the same candidate; shared snippets and affected routes are shown in review; unknown constructs survive unchanged; incompatibility blocks with retained input; no supported inventory setting is marked complete without a form/round-trip mapping. Build-specific forms may use declarative schemas, with custom controls where schemas cannot express directive semantics. Pin the reference build/module inventory before completing this story; arbitrary plugins are never advertised without coverage.

Explicit retained-feature matrix: multiple upstreams and load balancing; host/path matchers and prefix stripping; request/proxy/response headers; encoding/compression and when it applies; redirects; file roots/browse/index/hide/precompression; ordering, disabled/local-only routes, raw external data and global connection settings. Each maps to migration, editing, apply and regression checks or an explicit blocking incompatibility. No feature silently disappears behind generic “CRUD”.

**Rationale:** traces the already requested editing breadth to measurable completion without pretending universal module support exists.

### P08 — Report bounded verification and correct Caddy API claims (B08, B09)

**Artifacts:** architecture AD-4/9, old stories 1.2/2.5/3.5/4.3/5.1–5.5, SPEC CAP-2/6/7/8, EXPERIENCE Verification contract, QA fixtures.

**OLD:** `GET /config/` returns target version/modules; `/adapt?validate=true` provides semantic validation; HTTP 502 proves upstream unreachability, 403 proves filesystem permissions, and browser Fetch establishes client reachability.

**NEW API contract:** “`GET /config/` supplies live JSON and its ETag. Build/version/module inventory needs separately supplied, provenance-labelled build evidence or a supported capability mechanism; absence is unknown. `POST /adapt` with `text/caddyfile` returns the `result` JSON and optional warnings, or an error; decode the wrapper. It proves adaptation only. Remove the `validate=true` semantic guarantee. Provisioning may still fail on guarded application.”

Known missing-module adaptation errors block. Successful adaptation is narrower evidence, not an installed-module inventory. DNS-01 supported builds/providers require pinned build/module manifests and real compatibility cases including missing-plugin and provisioning failure. Unknown capability is shown explicitly; unsupported advertised scenarios cannot be declared complete. Use `/pki/ca/local/certificates` for chain retrieval, validate/select the root for a root-only download, and explain that the PKI request can provision the default CA. It is an explicit requested operation, not a side-effect-free generic connection probe. These contracts follow [pinned adapter code](https://github.com/caddyserver/caddy/blob/v2.9.1/caddyconfig/load.go#L120-L193) and [PKI API documentation](https://caddyserver.com/docs/api#get-pkicaidcertificates).

**NEW verification contract:** “Store observed DNS/TCP/TLS/HTTP facts separately from diagnostic hypotheses. A 502 is an unexpected HTTP response unless additional evidence identifies the cause; a file response 403 does not prove a permission problem. Direct upstream probes are separate observations with manager-network limitations. Published-service checks use the published URL, not Docker upstream names.”

Per-route expectations name scheme, hostname/SNI, port, path, allowed status codes, redirect Location policy and trust roots. Protected services may expect 401/403. Redirect probes do not follow the first hop. File probes use a known resource. Proposed initial timing defaults: 10-second per-check deadline, no automatic configuration resubmission, manual retry plus a bounded two-retry schedule at 5/15 seconds for transient observations; evidence becomes stale after five minutes or immediately on revision change. These are reviewable defaults, not supplied product performance targets. Certificate near-expiry warnings are distinct from certificate validation failure.

Browser Fetch failure may reflect CORS/CSP/mixed content/client trust; classify unavailable/unknown unless evidence establishes more. Opaque responses do not establish expected status/headers. Preserve a manual open-and-check path. Design CSP narrowly for explicitly selected probe origins; do not add wildcard outbound access. Client and manager results stay separate; final success names only the type, expectations, revision, time and vantage actually verified. Required unavailable checks leave the outcome applied/unverified or partially verified. [Fetch response model](https://fetch.spec.whatwg.org/#concept-filtered-response-opaque) supports these browser limits.

**Rationale:** prevents inaccurate success/failure causes and mock-only API assumptions. Include full API wrapper/error, CSP and browser-policy cases against the pinned reference target.

### P09 — Specify a functional upgrade rather than archive-only migration (B10)

**Artifacts:** architecture migration convention and new migration companion; old stories 1.1/6.4; PRD FR-10; upgrade fixtures.

**OLD:** archive `routes`/`snapshots`, create new schema and establish state from live Caddy.

**NEW:** “Use a versioned, transactional, restart-safe migration preserving usable local data without contacting or writing Caddy during database migration. Preserve a protected prior-schema backup and record migration version/checksum. If conversion fails, roll back and retain a usable recovery path; never mark migration complete from table existence.”

| Legacy data | Target treatment |
| --- | --- |
| Route ID/order/domain/path/handler/config/headers/prefix stripping | Mapped records with old-to-new identity and exact protected payload until source adoption validates canonical projection. Preserve all handler config fields. |
| Enabled and disabled/local-only routes | Preserve enabled state and locally stored intent separately from observed live baseline; do not silently publish missing/disabled routes on upgrade. |
| Raw external route, support status and reason | Protected raw preservation with exact order and scoped provenance; promotion to editable requires verified mapping. |
| `global_config` URL/encoding/server/setup/ETag | Preserve settings and legacy ownership record; root ownership remains pending explicit adoption. A route ETag cannot become a root ETag. |
| Route-array snapshots | Retain original ID/time/server/reason/ETag and `legacy_route_array` scope; allow inspection/export and only scoped reviewed restore where compatible. Never relabel full-config history. |

Canonical source becomes active only through P03 association; migrated local intent remains accessible and is reconciled against that baseline. Build an actual prior-release DB fixture including disabled routes, custom settings, raw routes and snapshot history. Cover repeated startup, interrupted migration, failure rollback and preserved editing/restore access. Document downgrade via protected prior-version backup and compatible image; no unsupported in-place downgrade promise.

**Rationale:** live configuration cannot recover local-only data. Source basis: `internal/storage/models.go` and `internal/storage/sqlite.go`.

### P10 — Complete backup scope and remote operator procedure (C03)

**Artifacts:** architecture AD-10/11; old story 6.4; PRD FR-10; runbook and QA backup scenarios.

**OLD:** complete backup is database + key + external `/data`.

**NEW:** “A recoverable deployment package includes a consistent SQLite backup, encryption key/key provenance, current applied configuration plus source/dependency bundle, external startup configuration/autosave (`/config` where applicable), certificate/CA assets (`/data` or configured storage equivalent), and deployment/build/import/environment assets required to reproduce the target. Record exact scope and acquisition time in a manifest.”

Use SQLite's consistent backup mechanism or an explicit stopped-manager procedure, not a naive copy of an active WAL database. Quiesce participating writers for a coordinated point-in-time package or label independently captured pieces and require reconciliation. External assets are acquired/restored by documented operator-side commands on the Caddy host. No manager Docker mutation or remote filesystem write privileges are added. Missing required assets make the package incomplete, not successful. Protect sensitive backup material and state key-loss consequences.

Clean-environment acceptance reconstructs both manager and external target, establishes source association, verifies current routes/settings/secrets/certificates and checks restart. UI recovery points exclude certificate private keys; sanitized configuration export alone is not this package.

**Rationale:** reconciles the existing external-Caddy deployment contract with complete recovery.

### P11 — Finish identity and input boundaries (C04)

**Artifacts:** architecture AD-1; old stories 2.1/2.2/4.5; PRD FR-3; glossary/tests.

**OLD:** parser injects IDs while unchanged parse/serialize is byte-identical; destination grammar is host:port only but HTTPS upstreams appear later.

**NEW:** “Parsing never mutates source bytes. On explicit save, assign missing IDs and show the source delta; duplicate IDs block save with locations and a repair action. Keep existing IDs stable across formatting edits. Disable retains source and identity through versioned enabled metadata, omits the route from compilation, and restores the same definition on enable. Editor deletion and disabled-source edits have explicit candidate diffs.”

Proposed supported examples: `service.example.com`, `*.lab.example.com`, Unicode hostname with validated IDNA projection; path `/api/*` plus separate prefix-strip setting; `192.0.2.10:8080`, `whoami:80`, `[2001:db8::10]:8080`; HTTPS destination `https://backend.example.com:8443` with separate upstream TLS/SNI/trust settings. Ports are 1–65535; reject `://`, unbracketed IPv6 with ambiguous port, embedded credentials, and unsupported destination URL paths rather than silently dropping them. Multi-upstream order and transport constraints must survive form/editor round trips.

Redirect accepts the supported 301/302/307/308 codes and validated absolute or relative Location with documented placeholder handling; reject CR/LF. File roots resolve on Caddy's filesystem; accept an absolute root such as `/srv/www`, reject NUL, and flag unsupported expansion syntax without claiming the manager can inspect the directory. Pin full matcher/URI/IDNA/placeholder grammar in a validation companion with valid/invalid examples and actual target tests before story readiness. These examples define a proposed boundary, not a claim of complete Caddy syntax coverage.

**Rationale:** resolves contradictory identity lifecycle and makes validation testable without silently losing supported settings.

### P12 — Re-slice six epics into complete user outcomes (B05)

**Artifact:** replace the Epic List, story allocation and FR coverage map in `epics.md`; add explicit dependencies and old-to-new alias mapping to QA/handoff. Old references below identify the existing 28-story document; they are not the proposed new IDs.

**OLD:** adoption needs 2.1/4.1/6.1; publication in Epic 2 needs drafts/apply in Epic 3; restore is available only after later recovery infrastructure.

**NEW target allocation:**

| Target story | Outcome and old-story provenance | Prerequisite |
| --- | --- | --- |
| 1.1 | Secure manager-only startup, origin/CSRF and adapter seams (old 1.1 platform portion) | None |
| 1.2 | Protected source/secret version storage (old 4.1 core, P05) | 1.1 |
| 1.3 | Transactional legacy migration and scoped history access (old 1.1 migration; P09) | 1.2 |
| 1.4 | Lossless source/identity and preservation mapping engine (old 2.1; P03/P11) | 1.2 |
| 1.5 | Target connection and honest capability/source context (old 1.2) | 1.1 |
| 1.6 | Reviewed adoption without remote mutation (old 1.3) | 1.3–1.5 |
| 1.7 | Revisioned durable shared draft (old 3.1 backend) | 1.6 |
| 1.8 | Source/runtime diff, review binding and conflict-safe rebase (old 3.3/3.4 core) | 1.7 |
| 1.9 | Durable conditional apply and restore with snapshot retention/reconciliation (old 3.5, 6.1 and 6.3 core) | 1.8 |
| 1.10 | Minimal workspace: publish public-HTTPS proxy, review/apply, bounded verification and recovery; prove restart (old 1.4 shell, 2.2 subset, 3.2 minimal bar, 4.2 baseline, 5.1/5.4 minimum) | 1.9 |
| 2.1 | Full proxy/redirect/file forms, retained settings and validation (remaining old 2.2) | Epic 1 |
| 2.2 | Full peer Caddyfile editor with safe switching and adaptation (old 2.5) | 2.1 |
| 2.3 | Optional read-only Docker suggestions/manual destination (old 2.3) | 2.1 |
| 2.4 | Search/filter/responsive route dashboard (old 2.4) | 2.1 |
| 2.5 | All declared supported global/server settings forms and coverage manifest (new; P07) | 2.2 |
| 3.1 | Complete Change bar/unsaved input/navigation (remaining old 3.2) | Epic 1 |
| 3.2 | Accessible redacted diff UI across all scope (remaining old 3.3) | 3.1 |
| 3.3 | Rich drift/rebase conflict-resolution UI (remaining old 3.4) | 3.2 |
| 3.4 | Cross-session stale-save/reconnect UX (remaining old 3.1) | 3.1 |
| 3.5 | Interrupted-operation reconciliation UX and fault cases (remaining old 3.5 UI) | 3.2 |
| 4.1 | Complete public HTTPS settings and issuance/renewal evidence (remaining old 4.2) | Epics 1–3 |
| 4.2 | DNS-01/wildcards, build evidence, secret version use (old 4.3) | 4.1 |
| 4.3 | Internal CA and explicit root/trust workflow (old 4.4) | 4.1 |
| 4.4 | HTTPS upstream TLS independently of client-facing TLS (old 4.5) | 4.1 |
| 5.1 | Full asynchronous manager probe expectations/revisions (remaining old 5.1) | Epics 1–4 |
| 5.2 | Redirect/file known-resource checks (old 5.2) | 5.1 |
| 5.3 | Bounded browser observations and scoped CSP (old 5.3) | 5.1 |
| 5.4 | Complete status evidence, stale results and retries (remaining old 5.4) | 5.2/5.3 |
| 5.5 | Evidence-based diagnostic guidance (old 5.5) | 5.4 |
| 6.1 | Complete history list/export and retention UI (remaining old 6.1) | Epic 1 |
| 6.2 | History comparison (old 6.2) | 6.1 |
| 6.3 | Full reviewed history restore UX (remaining old 6.3) | 6.2 |
| 6.4 | Complete backup/restore/upgrade runbooks and release evidence (old 6.4) | Epics 1–5; 6.3 |

Epic 1 becomes **Safely adopt, publish and recover a first service**. Epic 2 becomes **Complete route and global configuration editing**. Epics 3–6 extend review UX, TLS, verification and operational history. Minimal first-slice versions and later extensions have explicitly disjoint acceptance; developers do not implement the same service twice. New total is 33 stories, subject to final sizing within these outcomes. Splitting an oversized story must update the alias/dependency map, not introduce hidden prerequisite work.

Epic 1 completion means a real external-Caddy fixture can be adopted without rewrite, a proxy can be staged/reviewed/applied through the guard, evidence is bounded, and an immediate reviewed recovery works. It is an internal delivery slice, not a scope-reduced release: release still requires all confirmed capabilities.

Add required golden corpus, root-write compatibility harness and fault-injection seams to the stories that first need them (1.4, 1.5, 1.9); frontend automation hooks to relevant UI stories. Introduce domain/application/adapter seams alongside these use cases while keeping legacy mode functional; prohibit mixed legacy/root writers for the same adopted target. No generic rewrite-only epic or unsafe temporary immediate publication.

**Rationale:** every prerequisite precedes its consumer and the first publishable slice includes recovery. Preserve source aliases so QA numbering cannot masquerade as current story identity.

### P13 — Repair QA mapping, evidence labels and review dispositions (C02, C05)

**Artifacts:** both system test plans, progress checkpoint and handoff; architecture review files; SPEC glossary/diagrams.

**OLD:** handoff cites nonexistent old stories 3.6/3.7; public HTTPS/plugin/restore/backup rows point to different story subjects; tooling is marked Ready; API P95 <200ms and memory/CPU assumptions are presented alongside requirements; the SPEC says WCAG 2.1 and a diagram says port 8080.

**NEW:** “Use explicit test ID → current story ID → acceptance criterion → requirement/finding mappings. Existing test IDs may remain stable aliases, with a separate current-story column; their numeric prefix is not authority. Map old 3.6 uncertainty tests to target 1.9/3.5; public HTTPS to 1.10/4.1; DNS plugin cases to 4.2; guarded restore to 1.9/6.3; backup/upgrade to 1.3/6.4. Update all remaining mappings from P12, then recompute scenario totals from the actual matrix.”

Replace stale `/load` and semantic-adapt fixtures with the real request/response contracts. Add P0 scenarios for write race, source mismatch/nested preservation/import drift, post-submit crash/local commit failure, historical secret rotation/missing key, concurrent sessions/late probes, and prior-release migration. Add origin checking alongside header checking for all mutation methods, including PATCH. Test hooks remain planned work with owners/dependencies until implemented; they are not production mutation endpoints.

Performance/environment budgets and CI durations remain provisional until measured/approved; remove P95 <200ms, <100MB/<1% CPU and 10–200-route assumptions as binding gates. Retain the specified 500ms debounce as an interaction setting. Separate planning readiness, harness availability and passing runtime evidence. Automated accessibility scans complement keyboard/zoom/focus checks; normalize the target to WCAG 2.2 AA and default manager port 3000.

Record review dispositions instead of overwriting history: no-ETag claim rejected by pinned API/source; absent-identity and plaintext-recovery-JSON premises superseded by current AD-1/5, but P05/P11 still close their edge cases; read-before-`/load` recommendation replaced by P02; whole-HTTP passthrough validation replaced by P03. Secret lifecycle, source fidelity, revisions and evidence attribution remain open until their corrected contracts/tests exist. Runtime claims in earlier reviews remain dated evidence, not a new certification. Refresh review verdicts only against the corrected document revision.

**Rationale:** makes acceptance traceable and removes contradictory evidence claims without deleting provenance.

## 5. Implementation handoff and success criteria

This is a **moderate** scope correction: backlog reorganization plus concrete architectural contract completion inside the already approved direction. High technical risk warrants Architect ownership of protocol/preservation validation; it does not imply a new product discovery cycle.

| Order | Recipient / responsibility | Deliverable and completion condition |
| --- | --- | --- |
| 1 | Product owner/maintainer | Apply P01 scope reconciliation with stable FR IDs and source authority record. No repeated approval of existing product choices. |
| 2 | Solution Architect | Apply P02–P11 technical decisions to ADs, schema/API and compatibility/migration companions; resolve exact path mappings, supported build/settings manifest and import expansion model. Confirm all mechanisms are mutually implementable. |
| 3 | PO / Developer planning | Apply P12 to epics, produce old→new story aliases, size stories, assign dependencies and measurable acceptance. Full settings coverage must be explicit. |
| 4 | QA / test architect | Apply P13 to plans/handoff and map every B/C finding to acceptance/test intent. Assign harness and fault-hook delivery before consumers. |
| 5 | Architect / maintainer | Refresh earlier review dispositions and synchronize SPEC/glossary/diagrams/UX. Check cross-document contracts and links. |
| 6 | Sprint planner | Repeat readiness. Generate and validate sprint tracking only after PASS. All stories start at actual backlog status; document existence is not implementation progress. |
| 7 | Developer | Implement approved stories using corrected contracts; collect real runtime/upgrade/recovery evidence. Product implementation is subsequent work. |

Required measurable closure:

1. Every B01–B10 and C01–C05 has an artifact/section, owner, acceptance/test mapping and disposition; no unresolved blocker is hidden by a status label.
2. Apply and restore share root `If-Match`, reviewed target/revision binding, durable attempts, protected snapshots and reconciliation. The external-write race and crash/finalization cases are explicit acceptance scenarios.
3. Adoption has no implicit publish; nested raw data, source/import provenance, unknown constructs and retained local-only data have concrete preservation/blocking rules.
4. Every declared supported setting has form/editor/validation/secret coverage; all four TLS scenarios and retained baseline features trace to stories.
5. Draft/session/probe revisions, browser unknown states, truthful API contracts, complete backup and functional prior-release migration are consistent across artifacts.
6. Story dependencies are acyclic and forward implementable; the first service slice includes draft/review/apply/verification/recovery. All QA story references resolve.
7. Repeat readiness passes before sprint status generation; later runtime evidence is required before claiming product acceptance.

At initial proposal creation only this file had been written. After Artem’s explicit approval, documentation/backlog corrections were applied and the repeat planning gate passed. See the completion record below. Product implementation and runtime acceptance remain subsequent work.

## Checklist and execution record

| Checklist item | Status | Finding / disposition |
| --- | --- | --- |
| 1.1 Triggering story | [N/A] | Cross-artifact readiness FAIL before sprint creation; affected old story IDs recorded in each proposal. |
| 1.2 Core problem | [x] | Technical contract gaps and incomplete propagation of confirmed choices. |
| 1.3 Evidence | [x] | Readiness findings verified against current artifacts, storage source and primary API/source references. |
| 2.1 Current epic viability | [x] | Current ordered plan cannot be implemented safely as written. |
| 2.2 Epic modifications | [x] | P12 complete first outcome and prerequisite moves. |
| 2.3 Remaining epics | [x] | All six assessed in §2/P12. |
| 2.4 Obsolete/new epics | [x] | Keep six outcomes; add explicit global/server forms and migration stories. |
| 2.5 Priority/order | [x] | Safety/storage/source/draft/apply/restore before publication. |
| 3.1 PRD impact | [x] | P01/P07; confirmed scope retained. |
| 3.2 Architecture impact | [x] | P02–P11 and schema/state/diagram changes. |
| 3.3 UX impact | [x] | P03/P06/P08; visual direction preserved. |
| 3.4 Other artifacts | [x] | QA, reviews, operations and SPEC in P09/P10/P13. |
| 4.1 Direct adjustment | [x] | Viable and recommended. |
| 4.2 Rollback | [x] | Evaluated, not justified. |
| 4.3 MVP review | [x] | Evaluated; scope reduction not selected. |
| 4.4 Recommendation | [x] | Direct adjustment plus backlog reorganization, §3. |
| 5.1 Issue summary | [x] | §1. |
| 5.2 Impact | [x] | §2 and P01–P13. |
| 5.3 Rationale | [x] | §3 and per-proposal rationale. |
| 5.4 Scope/action plan | [x] | Confirmed release breadth remains; sequencing in P12/§5. |
| 5.5 Handoff plan | [x] | Ordered responsibilities in §5. |
| 6.1 Checklist review | [x] | Applicable analysis complete; pending review items explicit below. |
| 6.2 Proposal accuracy | [x] | Source/claim, link, coverage and dependency checks performed; technical proposals still require application and repeat readiness. |
| 6.3 User approval | [x] | Artem continued to approval and explicitly answered yes; corrections applied. |
| 6.4 Sprint tracking | [x] | Generated mechanically after repeat readiness PASS; six epics / 33 stories at backlog. |
| 6.5 Confirm handoff | [x] | Corrected artifacts, QA mapping and backlog delivered; Developer next story 1.1. |

Proposal drafting record: initialized configured workflow; loaded required/optional project artifacts; accepted readiness report as trigger; selected batch presentation to avoid repeating closed product discovery; assessed trigger/epics/artifacts/options; prepared explicit proposals and handoff. No subagents, external messages, code changes, deployment or sprint generation performed.

## Approved execution and handoff — 2026-10-17

Artem selected Continue, then explicitly answered yes to applying this proposal to planning documents, backlog and test mappings. P01–P13 were applied; no further product-scope approval is required. Classification remains Moderate, with high implementation risk managed by explicit safety acceptance. The historical FAIL review is preserved in `implementation-readiness-before-course-correction.md`; [current readiness](implementation-readiness.md) is PASS for planning.

Delivered corrections: PRD/addendum/reconciliations; architecture spine plus shared configuration-contract companion; UX linkage; SPEC/glossary/diagrams; six epics and 33 current stories with old→new aliases and 122 ACs; QA/architecture test plans/handoff/checkpoint; independent review dispositions and structural validation; deterministic sprint tracking.

Architecture skill required three independent review lenses. Their findings were addressed: immutable endpoint/transport identity for historical attempts; immutable enabled/node maps for source/recovery; durable draft/target lock and revision CAS during submission/finalization; early vault harness wording; actual Tailwind lock seed. These tighten the approved shared-contract design without changing feature scope. Fresh reviews and follow-up closure found no remaining planning blocker; mechanical spine lint is clean.

Handoff is recorded to Product Owner/Developer for ordered story delivery, with Architect/QA contracts and AC mapping supplied. Developer starts **1.1 Secure manager-only startup and application seams** from [epics](epics.md), using [sprint tracking](../implementation-artifacts/sprint-status.yaml). All stories remain backlog, with no standalone story completion or product runtime test results claimed. Success means implementing the specified safe workflows and collecting acceptance evidence before full-release claims.

Correct Course workflow deliverables are complete: approved before/after proposal, applied artifact corrections, implementable backlog and handoff plan. No product code/dependency/Compose change, deploy, commit or external message was performed. Pre-existing workspace changes were retained. Logs and validation distinguish document correctness from runtime verification.
