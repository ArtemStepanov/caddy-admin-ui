> Historical readiness snapshot before approved P01–P13. Superseded by [current readiness](implementation-readiness.md).

---
project: caddy-orchestrator
date: 2026-10-17
intent: sprint-planning
status: complete
gate: FAIL
tracking_generated: false
method: artifact-consistency-and-implementability-review
---

# Implementation readiness — Caddy Orchestrator

**Verdict: FAIL.** The planning package exists, but developers cannot implement it consistently without inventing decisions about configuration preservation, concurrency, recovery, and supported editing. Resolve the blocking findings before generating sprint tracking.

This review follows `bmad-sprint-planning/references/readiness-gate.md`. It evaluates the recorded plan, with selective source inspection and official API references. It does not certify runtime behavior. No product tests, builds, deployments, or implementation changes were performed.

## Inventory and authority

| Artifact | State and role |
| --- | --- |
| `CLAUDE.md`, `AGENTS.md`, current source and Compose | Brownfield implementation baseline: scoped route-array ownership, guarded writes, existing routes/settings/snapshots. These describe current behavior, not the intended expanded product. |
| [Research](research/competitive-caddy-homelab-manager-2026-10-16/research.md) and brief | Complete; advisory evidence, not scope authorization. |
| [PRD](prds/prd-caddy-orchestrator-2026-10-16/prd.md), addendum and reconciliations | Draft; retains narrower ownership and candidate features that later user decisions supersede. |
| [DESIGN](ux-designs/ux-caddy-orchestrator-2026-10-16/DESIGN.md), [EXPERIENCE](ux-designs/ux-caddy-orchestrator-2026-10-16/EXPERIENCE.md), [UX memory](ux-designs/ux-caddy-orchestrator-2026-10-16/.memlog.md) and reconciliation | Final handoff with explicit proposals. Records confirmed expanded scope and unresolved technical contracts. Three representative static mocks exist. |
| [SPEC](../specs/spec-caddy-orchestrator/SPEC.md), glossary and diagrams | Declares itself the canonical contract, includes architecture and UX as companions. Some promises exceed the compatibility and safety decisions actually recorded. |
| [Architecture spine](architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md) and reviews | Final; AD-1–AD-11, schema and structure exist. Earlier review findings have mixed dispositions and some incorrect or superseded premises. |
| [Epics](epics.md) | Final; **6 epics, 28 stories**, unique IDs. FR coverage is listed, but it is not complete acceptance-level traceability. |
| [Architecture test design](../test-artifacts/test-design-architecture.md), [QA plan](../test-artifacts/test-design-qa.md), progress and [handoff](../test-artifacts/test-design/caddy-orchestrator-handoff.md) | System test strategy exists; story references and parts of the API assumptions diverge from the final epics. Planned tests are not execution evidence. |
| Implementation tracking | No existing `sprint-status.yaml` or implementation story files found. Their absence is normal before this planning step. |
| Configured `docs/` knowledge directory | Absent. Root documentation and planning artifacts provide context; this absence is not a blocker. |

Read-only structural checks found no duplicate story IDs and no broken relative Markdown links in the non-research planning/spec/test Markdown inspected. This is a file-existence check, not validation of headings, claims, or web links.

### Confirmed scope must not be reopened

EXPERIENCE Foundation and the UX memory explicitly record user choices for full existing configuration management, forms for all settings plus a Caddyfile editor, shared batch drafts, all four TLS scenarios, manual input plus Docker discovery, and the Quiet Workshop direction. These supersede the narrower PRD proposals. Missing propagation is a documentation problem; it is not evidence that the user still needs to choose those capabilities.

Artem additionally confirmed in this readiness conversation that the expanded scope was already approved with the architect. The architecture rubric itself calls root-level CST management a deliberate, approved evolution and concludes APPROVED with Minor Issues. Treat the expanded scope and architectural direction as settled. The findings below concern concrete implementation gaps and discrepancies; they do not request renewed scope approval or replacement of the chosen direction.

The user has also closed further generic discovery. Remaining gaps below concern implementation contracts and fidelity to the recorded choices. A new product decision is necessary only if a proposed technical solution changes those choices.

## Blocking findings, ordered by severity

### B01 — Critical: the planned apply/restore protocol does not prevent concurrent overwrites

**Evidence:** Architecture AD-3 step 4 explicitly says `/load` does not evaluate `If-Match`; AD-8, SPEC CAP-5, and stories 3.4/3.5 promise protection using a preceding `GET /config/`. Story 6.3 repeats the same guard for restore. AD-8 also describes rebase without defining conflict semantics.

An external writer can change Caddy after the comparison and before `/load`; the manager then overwrites that change. A manager-local mutex cannot synchronize an external writer. Merely refreshing the draft's base ETag also does not merge external changes into the candidate. Caddy documents conditional mutation at `/config/` and the absence of transactions spanning requests. Its 2.9.1 `/load` handler calls `caddy.Load` without a precondition. See [concurrency documentation](https://caddyserver.com/docs/api#concurrent-config-changes), [config handler](https://github.com/caddyserver/caddy/blob/v2.9.1/admin.go#L891-L959), and [load handler](https://github.com/caddyserver/caddy/blob/v2.9.1/caddyconfig/load.go#L63-L119).

**Required resolution:** Record a write protocol with an atomic remote precondition, or a technically enforced exclusive-writer contract with honest limits. Define rebase against baseline/live/candidate, overlapping edit conflicts, renewed review, restore target checks, and the ETag captured by the restore review. Update the shared contract and acceptance criteria, including an external write inserted between final read and write.

**Fix workflow:** `bmad-architecture`, then `bmad-correct-course` to propagate the protocol across SPEC, epics, UX and test design.

### B02 — Critical: existing-config adoption and passthrough merging have no lossless contract

**Evidence:** Architecture AD-2; SPEC CAP-1/CAP-2; stories 1.3 and 2.1; EXPERIENCE Source reconciliation and Interaction Primitives. AD-2 synthesizes HTTP Caddyfile text from JSON and retains non-convertible routes in `passthrough_apps`, but describes merging at root level. Story 1.3 only explicitly preserves non-HTTP top-level applications.

No contract specifies how to combine unsupported routes inside `apps.http.servers`, server/listener settings, other HTTP servers, and generated HTTP configuration without dropping or overwriting one side. `logging` and `storage` are root configuration fields, while `pki` is an application under `apps`; a single unspecified app merge does not define preservation. The plan also leaves source acquisition, freshness against live JSON, import files and environment-dependent expansion unresolved for remote Caddy.

Caddyfile cannot represent every native JSON configuration; automatic inverse reconstruction cannot recover original comments or imports. This limitation is already acknowledged by UX and supported by the [Caddyfile documentation](https://caddyserver.com/docs/caddyfile#note). Logging a synthesis notice is insufficient to establish equivalence or safe whole-config editing.

**Required resolution:** Define source acquisition and its association with a live baseline; the supported conversion boundary; preservation at exact JSON paths; merge ownership/collisions/order; import ownership and dependency checks; and explicit blocking states for incompatible source. Adoption must not publish a rewrite implicitly, must reject a changed preview, and must retain the original recoverable configuration. Apply, diff, discard and rebase must use this same model.

**Fix workflow:** `bmad-architecture`; reconcile SPEC/UX and add adoption/preservation cases with `bmad-create-epics-and-stories`.

### B03 — High: uncertain apply recovery covers transport failure but not crash/finalization failure

**Evidence:** Architecture AD-3/AD-10, state diagram and DRAFT schema; SPEC CAP-8; story 3.5. PRD FR-5 explicitly includes acknowledgement, ETag-read and local-persistence failures.

The plan marks uncertainty after a network error, but does not persist a write intent before submission or define startup behavior after a process dies during the write. A successful remote response followed by a failed ETag read or SQLite commit has no specified transition. Reading a live ETag alone cannot determine whether the candidate, old baseline, or a third-party configuration is present without persisted comparison material and a defined comparison algorithm. The test-design architecture proposes comparing an ETag to a candidate hash without specifying equivalent encoding.

**Required resolution:** Record durable apply/restore attempt identity, baseline and candidate identity, transitions before/after remote submission, startup reconciliation, and safe locking when local persistence itself fails. Define outcomes for live matching candidate, baseline, another config, or unavailable Caddy. Do not infer success from a changed ETag alone. Include post-acceptance local failures and process death in acceptance criteria.

**Fix workflow:** `bmad-architecture`, then `bmad-create-epics-and-stories`; propagate to SPEC diagrams and test design.

### B04 — High: secret preservation and historical restore are underspecified

**Evidence:** Architecture AD-5 and schema; stories 1.3, 4.1, 4.3, 6.1 and 6.3; SPEC secret constraints.

Encrypted credentials and sanitized recovery JSON are recorded, but there is no concrete representation for secrets in `DRAFT.caddyfile_content` or `RECOVERY_POINT.caddyfile_content`. Initial adoption may already contain credentials before the vault story is delivered. Replacing runtime secrets with `secret_ref:<id>` is not enough to restore an immutable historical configuration if a credential is later rotated, deleted, or unavailable; restore does not specify reference hydration. Masking all API output also needs an update contract that distinguishes unchanged secret references from literal replacement text.

**Required resolution:** Define extraction and protected storage on adoption, references in every persisted text/JSON representation, historical credential version retention, hydration for adapt/apply/restore, failure when key/reference material is missing, and API/editor semantics. Cover redaction of remote error text and round-trip editing without resaving `***` as a credential. The vault prerequisite must exist before any workflow stores a secret-bearing baseline.

**Fix workflow:** `bmad-architecture`, followed by story/dependency and test-design updates.

### B05 — High: epic order contains forward dependencies and incomplete user outcomes

**Evidence:** Epics 1–3 and 6; stories 1.3, 2.1/2.2/2.5, 3.1/3.3/3.5, 4.1 and 6.1.

| Earlier outcome | Later prerequisite currently assigned elsewhere |
| --- | --- |
| 1.3 adopts canonical text and captures a safe immutable baseline | CST/conversion in 2.1, secret treatment in 4.1, recovery storage in 6.1 |
| Epic 2 promises service publication and restart persistence | Persistent draft in 3.1 and the only explicit batch apply pipeline in 3.5 |
| 3.3/3.5 require redacted diffs and pre-write recovery points | Secret lifecycle in 4.1 and recovery storage/retention in 6.1 |

As recorded, implementing in epic/story order either leaves earlier outcomes unfinishable or duplicates later stories. The plan also needs an explicit path for introducing the new application/domain/adapters while preserving the running brownfield behavior.

**Required resolution:** Re-slice the epics into implementable user outcomes, move the minimum prerequisites into the first outcome that needs them, and state story dependencies and completion boundaries. Make the first publishable slice include safe draft/review/apply/recovery. Do not solve ordering by allowing temporary immediate writes or unsafe snapshots.

**Fix workflow:** `bmad-create-epics-and-stories`, using the architecture decisions after their correction; `bmad-correct-course` if the changes span the full handoff.

### B06 — High: confirmed editing breadth does not trace into concrete stories

**Evidence:** UX memory and EXPERIENCE Foundation confirm forms for all settings and management of previous global settings. Architecture AD-1 and story 2.5 leave unrepresented directives editor-managed. SPEC CAP-2 names global options, but the story set contains no concrete global/server/settings editing outcome or form-coverage contract.

Listing a `Configuration.tsx` page does not specify what can be edited or how build-specific options gain forms. The coverage map claims all FRs while losing this confirmed capability. Existing multi-Upstream, host/path, headers and compression functionality also lacks an explicit preservation/migration acceptance mapping; general route CRUD does not settle those regressions.

**Required resolution:** Define supported settings/build coverage and how forms expose settings across that coverage, including global/server options and temporarily unsupported constructs. Map the confirmed intent to stories and measurable acceptance criteria without silently making some required settings editor-only. Explicitly map retained brownfield route features to new behavior and regression checks. Preserve the Caddyfile editor choice.

**Fix workflow:** `bmad-architecture` for compatibility/form strategy, `bmad-prd` for scope reconciliation, and `bmad-create-epics-and-stories` for coverage.

### B07 — High: shared edits and asynchronous evidence lack configuration revisions

**Evidence:** Architecture AD-7, DRAFT and VERIFICATION_RESULT schema; stories 3.1 and 5.1/5.4. The shared desktop/phone draft is required, but only Caddy-side concurrency is specified.

Two sessions can save based on the same draft and overwrite one another; SSE notifications do not prevent this. No expected draft revision or stale-write behavior is recorded. Verification rows have route/time/vantage but no deployment or configuration revision, allowing a late result from configuration A to be displayed for configuration B after a quick edit or restore.

**Required resolution:** Define draft revision checks, stale-save rejection or safe merge, local unsaved-input handling, and reload/reconnect semantics. Bind probe jobs/results to target and applied configuration revision and specify rejection/labeling of stale results. Keep the proposed API, schema, state store and acceptance criteria consistent.

**Fix workflow:** `bmad-architecture`, followed by stories 3.1 and 5.x and their test mappings.

### B08 — High: verification can report causes and client outcomes that its observations do not establish

**Evidence:** Architecture AD-4; stories 5.1–5.5; EXPERIENCE Verification contract. Story 5.1 treats an HTTP 502 as an upstream reachability failure; 5.2 identifies filesystem permission errors from a failed HTTP assertion. Story 5.3 treats browser Fetch/CORS as client reachability without specifying unavailable/opaque outcomes. Current `cmd/server/main.go` CSP defaults to same-origin connections.

An HTTP status alone does not prove its cause or directly establish upstream reachability. A protected service may intentionally return 401/403. Browser requests can fail under CORS, CSP, mixed-content or client trust restrictions even when a service is reachable. An opaque Fetch response exposes neither status nor response headers; see the [Fetch response model](https://fetch.spec.whatwg.org/#concept-filtered-response-opaque). The existing UX explicitly requires bounded evidence and unknown states.

**Required resolution:** Define observations separately from inferred causes, direct upstream probes when supported, per-type expected response/path/location policies, and browser limitations with unknown/unavailable fallbacks. Specify probe timeouts, retries, staleness, internal-CA trust handling and what evidence permits final success. Include the CSP change needed for supported client probes as implementation work. Tests must not equate a Docker upstream name with the browser's published service URL.

**Fix workflow:** `bmad-architecture`, align UX technical behavior, and refine epic 5 through `bmad-create-epics-and-stories`.

### B09 — High: Caddy capability and validation assumptions do not match the pinned API

**Evidence:** Architecture AD-3/AD-9 and stack target Caddy 2.9.x; story 1.2 promises version and modules from `GET /config/`; stories 2.5/3.5/4.3 rely on remote adapt checks.

The 2.9.1 adapt handler does not inspect a `validate` query parameter and calls the adapter with nil options; `POST /adapt?validate=true` therefore does not establish full module provisioning or runtime validity. Some missing-module problems can still be detected during adaptation, but adaptation is not a complete semantic dry run. `GET /config/` exports configuration and is not a documented installed-module/version inventory. See the pinned [adapt handler](https://github.com/caddyserver/caddy/blob/v2.9.1/caddyconfig/load.go#L120-L193) and [config endpoint documentation](https://caddyserver.com/docs/api#get-configpath).

**Required resolution:** Specify exact supported API request/response shapes, what adapt proves, what can only fail on application, and a supported capability/version strategy with unknown states. Define provider/build support and plugin error handling without inventing endpoints or treating successful compilation as runtime validation. Make real pinned-Caddy compatibility scenarios part of the plan rather than encoding the mistaken assumption only in mocks.

**Fix workflow:** `bmad-architecture`, then SPEC, epics and `bmad-testarch-test-design` updates.

### B10 — High: archiving legacy tables does not define a working upgrade

**Evidence:** Architecture Consistency Conventions, story 1.1 and story 6.4; PRD FR-10. Existing storage includes `routes`, `snapshots`, `global_config`, raw external route data, ordering and disabled routes.

The only explicit migration action archives `routes` and `snapshots` before new schema creation. It does not specify how existing active routes, disabled/local-only routes, settings, target ownership and route-scoped snapshots remain usable in the new model. Re-adopting live Caddy cannot recover data that exists only locally. Keeping an inaccessible backup table is not equivalent to preserving functioning routes/settings across upgrade.

**Required resolution:** Define a versioned, transactional, restart-safe migration with explicit old-to-new data mapping, target/ownership treatment, retention of old recovery scope, migration failure handling and downgrade limits. Include a prior-release database fixture with disabled routes and custom settings, and document how the new canonical configuration is established without silently publishing.

**Fix workflow:** `bmad-architecture`, followed by migration and upgrade stories in `bmad-create-epics-and-stories` and the test plan.

## Additional consistency concerns

### C01 — Superseded PRD and conflicting contract status

PRD terms, FR-1/4/8/9, scope exclusions, addendum and remaining decisions still describe the earlier scope. UX reconciliation correctly flags this; SPEC calls itself canonical while including those UX files as companions, so its stronger fidelity/atomicity promises do not settle their cautions. Reconcile the existing user decisions into stable FR IDs and an explicit source-precedence/change record; do not describe this as needing approval of already answered choices. Keep current-code documentation distinguishable from future architecture.

**Fix workflow:** `bmad-prd` and `bmad-correct-course`.

### C02 — QA references and gates are not aligned to the final story set

The handoff refers to **stories 3.6 and 3.7**, which do not exist; uncertainty is in 3.5. QA rows for public HTTPS, plugin detection, restore and backup use numbers that differ from the current story subjects. If test IDs are intentionally immutable, provide a separate explicit mapping rather than inferring story identity from the number. CSRF/origin acceptance and proposed testability hooks are not fully incorporated into the relevant stories. API P95 <200ms appears in the QA plan without a chosen environment or corresponding agreed performance budget; other artifacts leave the budget unknown. Tooling marked `Ready` is not supported by implementation evidence in this review.

The test harness, corpus and hooks may legitimately be future work; their absence alone does not block planning. Assign their delivery before dependent acceptance checks, correct mappings, and label proposed thresholds and unverified tooling honestly.

**Fix workflow:** `bmad-testarch-test-design`, then update epic acceptance mappings.

### C03 — Backup scope and remote execution model need completion

AD-10 requires external `/config` and `/data`, whereas AD-11/story 6.4 call database + key + `/data` the complete package. The procedure must preserve or explicitly reconstruct external startup configuration and required Caddyfile/import/env/plugin assets; a pre-write history point alone need not contain the current applied config. A manager-only deployment also needs an operator-side acquisition/restore contract for remote Caddy volumes. Do not silently grant the manager filesystem or Docker write access to obtain a backup.

**Fix workflow:** `bmad-architecture` and story 6.4/runbook refinement.

### C04 — Input and identity edge contracts remain incomplete

PRD FR-3 explicitly requires supported host/path/file-root examples before final stories. Story 2.2 defines some upstream and redirect cases but not the full host/path, IPv6/IDN, URI/root and multi-Upstream boundary; story 4.5 adds HTTPS destinations beyond the earlier host:port-only wording. CST story 2.1 both injects IDs on parse and promises byte-identical unchanged round trips, whereas AD-1 injects missing IDs on first save. Define lifecycle for missing/duplicate IDs and disabled routes, preserving source until the intended mutation. These can be refined without revisiting the user's feature choices.

**Fix workflow:** `bmad-architecture`/`bmad-create-epics-and-stories`; propagate stable validation rules to the SPEC and tests.

### C05 — Earlier reviews require explicit disposition

The adversarial review still contains a Critical Issues verdict alongside positive reviews, but some premises are wrong or superseded: Caddy does supply config ETags, and the current spine already includes comment-based route IDs and sanitized recovery JSON. Do not carry those claims forward as new defects. Other concerns, such as source/merge fidelity, draft revisions and evidence attribution, remain unresolved in the current text. Record each disposition against the current spine and refresh affected reviews. Minor drift also includes SPEC WCAG 2.1 versus UX/epics WCAG 2.2 and a diagram's admin UI port 8080 versus the current deployment's 3000.

**Fix workflow:** `bmad-architecture` validation after the corrections, then companion/test alignment.

## Resolution order and repeat gate

1. Synchronize stale PRD clauses and cross-references with the already approved UX and architecture scope. This is documentation maintenance, not a new scope decision or approval gate.
2. Resolve B01–B04 and B07–B10 in architecture: write safety, source/preservation, durable recovery, secret lifecycle, revisions, evidence, API compatibility and upgrade mapping.
3. Re-slice epics for B05, complete coverage for B06, and map every retained requirement/decision to concrete acceptance criteria. Add the required backup, validation and testability work.
4. Refresh test-design mappings, thresholds and previous review dispositions against the corrected plan.
5. Re-run `bmad-sprint-planning` readiness. On PASS, use its deterministic generator for `implementation-artifacts/sprint-status.yaml` and validate that report. Do not manufacture progress or ready-for-dev status from document existence.

Apply focused corrections to architecture, stories and test mappings within the approved direction. **`bmad-correct-course`** can coordinate cross-document propagation if useful; restarting product discovery or repeating scope approval is not a prerequisite. No new broad product questionnaire is needed.

## Workflow outcome

The readiness review is complete; sprint tracking generation is intentionally pending. The invoked skill says to load generation instructions **“then on PASS”** and, for FAIL, **“Present findings ordered by severity … and stop”**. The gate cannot honestly be marked PASS while these implementation decisions remain open.
