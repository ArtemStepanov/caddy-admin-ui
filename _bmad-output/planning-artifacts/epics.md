---
status: final
updated: 2026-10-17
revision: course-correction-1
stepsCompleted: [step-01-validate-prerequisites, step-02-design-epics, step-03-create-stories, step-04-final-validation]
inputDocuments:
  - _bmad-output/planning-artifacts/prds/prd-caddy-orchestrator-2026-10-16/prd.md
  - _bmad-output/planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md
  - _bmad-output/planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/configuration-contracts.md
  - _bmad-output/planning-artifacts/ux-designs/ux-caddy-orchestrator-2026-10-16/DESIGN.md
  - _bmad-output/planning-artifacts/ux-designs/ux-caddy-orchestrator-2026-10-16/EXPERIENCE.md
  - _bmad-output/specs/spec-caddy-orchestrator/SPEC.md
---

# Caddy Orchestrator — Implementation Backlog

Approved [P01–P13](sprint-change-proposal-2026-10-17.md) reorganizes six epics into **33 stories**. These are backlog definitions, not implemented or ready-for-dev story files. Stable FR-1–FR-11 meanings come from the [corrected PRD](prds/prd-caddy-orchestrator-2026-10-16/prd.md). [Architecture AD-1–AD-11](architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md) and [configuration contracts](architecture/architecture-caddy-orchestrator-2026-10-17/configuration-contracts.md) define shared behavior. [DESIGN](ux-designs/ux-caddy-orchestrator-2026-10-16/DESIGN.md)/[EXPERIENCE](ux-designs/ux-caddy-orchestrator-2026-10-16/EXPERIENCE.md) retain detailed appearance/accessibility and corrected interactions. Full release scope remains; the first service slice is not a scope-reduced release.

## Requirements inventory

| Requirement | Outcome | Primary consumers |
| --- | --- | --- |
| FR-1 | Explicit existing-target whole-config adoption | 1.1–1.6, 1.10 |
| FR-2 | Publish, retain and recover existing route capabilities | 1.3, 1.7–1.10, 2.1, 6.4 |
| FR-3 | Complete supported forms and validated destinations | 1.4, 2.1–2.5, 4.4 |
| FR-4 | Source/nested/root preservation | 1.2, 1.4, 1.6, 1.8, 2.1, 2.2, 2.5 |
| FR-5 | Concurrent shared edits, guarded publication, uncertain recovery | 1.7–1.9, 3.1–3.5, 6.3 |
| FR-6 | Honest revision/time/vantage status | 1.10, 2.4, 3.1, 3.4, 3.5, 4.1, 4.3, 5.1–5.4 |
| FR-7 | Evidence-based failure diagnosis | 1.5, 4.1, 4.3, 5.1–5.5 |
| FR-8 | All four TLS scenarios | 1.2, 1.5, 1.10, 2.5, 4.1–4.4, 5.1 |
| FR-9 | Scoped protected recovery | 1.2, 1.3, 1.6, 1.9, 1.10, 3.5, 6.1–6.3 |
| FR-10 | Complete backup and functional upgrade | 1.1–1.3, 6.4 |
| FR-11 | Quiet Workshop accessible forms/editor/review | 1.1, 1.4, 1.10, Epics 2–6 |

### Cross-cutting acceptance

- **NFR-SEC-1/2/4:** listener-based Basic Auth, paired credentials, constant-time checks; CSRF header plus origin policy on all mutation methods; private Admin API. Delivered first in 1.1, retained on every added API.
- **NFR-SEC-3:** encrypted exact source/runtime/legacy data, immutable secret versions, safe keep/replace/remove and redaction, delivered in 1.2 before import. No permanent plaintext archives or masks saved as secrets.
- **NFR-DATA-1/2:** WAL/foreign keys/busy timeout; target/draft/review/attempt/applied identities; root remote guard; versioned transactional migration with functional data mappings. Stories 1.3, 1.7–1.9.
- **NFR-REL-1/2/3:** restart persistence, durable unresolved-operation locks, at most 100 recoveries with active pins/capacity rejection; complete external config/data backup. Stories 1.9/1.10/6.4. API connectivity is not evidence of external mount/resume flags.
- **NFR-COMPAT-1:** pinned Caddy 2.9.1 + reproducible provider fixture, real root guard/adapt wrapper/PKI behavior, honest unknown build state. Stories 1.4/1.5/2.5/4.2/4.3. The manager has no Caddy binary.
- **NFR-A11Y-1:** WCAG 2.2 AA target; keyboard/editor escape/focus restoration, textual diff, labels/errors, non-color state, 44px touch, 320px reflow, 200% zoom, reduced motion/transparency. All UI acceptance includes this floor; full audit is not claimed from document creation.
- **NFR-PERF-1:** 500ms adaptation debounce with stale-response rejection. Resource/latency budgets, scenario counts and CI duration targets require measurements and are not invented release thresholds.

The eleven UX components in DESIGN (header, route row, search, form, editor, change bar, diff, evidence, action button, recovery entry, feedback) are covered across 1.10, 2.1–2.5, 3.1–3.5, 5.3–5.5 and 6.1–6.3. Their token/spacing/light/dark contracts remain in DESIGN; no theme default or new stack is silently selected. UI stories include stable automation hooks and semantic roles/labels; test-only transport/fault hooks are delivered in 1.4/1.5/1.9 before dependent acceptance, never exposed as unauthenticated production endpoints.

## Delivery order and completion boundaries

Implement prerequisites listed on each story; all refer backward. Epic 1 infrastructure stories have an explicit eventual operator outcome and can be checked independently through ports/APIs; 1.10 integrates the first complete user workflow. Later stories extend its minimal UI/probes/settings and do not recreate the safety backend. Root-adopted targets have one server-enforced writer mode, including across restart. Preserve current legacy behavior while introducing the new seams; no unsafe interim immediate full-root writes.

Every numbered acceptance item is a stable `<story>.AC<n>` reference for QA. Planned fixture/corpus/build-lock delivery is part of these stories; missing implementation evidence does not imply planning decisions remain unspecified. If implementation sizing splits a story, preserve aliases and update the dependency/test map before work starts.

## Epic 1: Safely adopt, publish and recover a first service

The operator can deploy securely, preserve prior data/source, adopt without rewrite, stage/review/apply a public-HTTPS proxy, inspect bounded evidence and perform immediate guarded recovery. Core safety is complete before this first publication.

### Story 1.1: Secure manager-only startup and application seams

As a homelab operator,
I want to launch the manager securely with persistent data,
So that my existing Caddy remains separate and protected.

**Depends on:** none. **Requirements:** FR-1, FR-10, FR-11. **Legacy source:** 1.1 platform (pre-correction IDs).

**Acceptance Criteria:**

**1.1.AC1** — **Given** a non-loopback listener without credentials, **when** startup runs, **then** it fails unless ALLOW_INSECURE_NO_AUTH is explicit; a partial credential pair always fails.

**1.1.AC2** — **Given** valid paired credentials and manager-only Compose, **when** authenticated and unauthenticated requests arrive, **then** 401/authorized responses follow listener policy, using constant-time comparison; /app/data is persistent and no production Caddy is bundled.

**1.1.AC3** — **Given** any browser mutation using POST/PUT/PATCH/DELETE, **when** origin or required CSRF header is invalid, **then** the request is rejected before side effects; correct same-origin authenticated requests work.

**1.1.AC4** — **Given** current legacy source and an existing DB, **when** new application/port seams are introduced, **then** the current baseline remains functional, schema detection is nondestructive, and root-adopted targets cannot use legacy mutation paths.

**Completion boundary:** Core backend/storage/API outcome; no later UI story is a prerequisite for these acceptance checks.

### Story 1.2: Protected source and immutable secret-version storage

As a homelab operator,
I want to retain configuration and credentials safely,
So that drafts and old configurations remain recoverable without disclosing credentials.

**Depends on:** 1.1. **Requirements:** FR-4, FR-8, FR-9, FR-10. **Legacy source:** 4.1 core (pre-correction IDs).

**Acceptance Criteria:**

**1.2.AC1** — **Given** source/runtime/legacy payloads containing known and unclassified secrets, **when** they are persisted, **then** authenticated encrypted envelopes cover every exact payload; safe projections contain references and unclassified opaque content is withheld.

**1.2.AC2** — **Given** keep/replace/remove secret operations, **when** the draft is round-tripped through the secret port/API harness (literal forms/editor integration follows in 2.1/2.2), **then** unchanged references stay unchanged; a displayed mask never replaces a credential.

**1.2.AC3** — **Given** a retained history point referencing credential version A, **when** the current credential rotates to B or is deleted, **then** A remains immutable and hydratable for its retained users; missing key/version blocks operation.

**1.2.AC4** — **Given** adaptation errors, logs, diffs or ordinary exports, **when** secret-bearing values are processed, **then** safe output reveals no plaintext, including remote error echoes; key loss and sanitized-export limitations are explicit.

**Completion boundary:** Core backend/storage/API outcome; no later UI story is a prerequisite for these acceptance checks.

### Story 1.3: Transactional legacy migration and scoped history access

As a homelab operator,
I want to upgrade without losing locally stored routes and settings,
So that disabled routes and recovery history remain usable.

**Depends on:** 1.2. **Requirements:** FR-2, FR-9, FR-10. **Legacy source:** 1.1 migration; 6.4 upgrade (pre-correction IDs).

**Acceptance Criteria:**

**1.3.AC1** — **Given** a prior-release DB with active/disabled/local-only routes, raw route order, custom globals and snapshots, **when** migration runs, **then** a versioned transaction maps every record and ID with exact payload checks and preserves snapshot server/scope; no Caddy write occurs.

**1.3.AC2** — **Given** migration failure or process interruption, **when** the service restarts, **then** the old committed state or completed new version remains usable, with protected prior-version recovery and no partial completion marker.

**1.3.AC3** — **Given** migrated route-scoped ownership and ETag, **when** root adoption has not occurred, **then** legacy state remains accessible and is never presented as root adoption or a root ETag.

**1.3.AC4** — **Given** a consistent protected migration backup, **when** downgrade is needed, **then** the documented compatible prior image/backup route works and encrypted legacy copies do not create a permanent plaintext archive.

**Completion boundary:** Core backend/storage/API outcome; no later UI story is a prerequisite for these acceptance checks.

### Story 1.4: Lossless source, identity and preservation mapping engine

As a homelab operator,
I want to keep my Caddyfile and unrelated runtime data intact,
So that forms cannot silently narrow my configuration.

**Depends on:** 1.2. **Requirements:** FR-3, FR-4, FR-11. **Legacy source:** 2.1 (pre-correction IDs).

**Acceptance Criteria:**

**1.4.AC1** — **Given** unchanged Caddyfiles with multiple route nodes/site, comments, nesting and snippets, **when** parse/serialize executes, **then** bytes match exactly; explicit save injects only missing IDs and duplicate IDs block with repair locations.

**1.4.AC2** — **Given** a captured source bundle with file imports/environment inputs, **when** transient compilation input is built, **then** virtual paths, expansion order, explicit substitutions and source map are reproducible, originals remain intact, and unsupported/cyclic/escaping input blocks.

**1.4.AC3** — **Given** live nested raw routes, listeners, other servers, root logging/storage and apps.pki, **when** generated changes are composed, **then** only proven owned paths change; preserved arrays retain order and mixed-array modifications block without a complete mapping.

**1.4.AC4** — **Given** a golden corpus for supported constructs, **when** forms/text/mapping and versioned enabled metadata round-trip, **then** stable IDs survive, disabled source is retained, source/adaptation parity fixtures exist and historical restore retains its original enabled map, and no blanket JSON inversion is assumed.

**Completion boundary:** Core backend/storage/API outcome; no later UI story is a prerequisite for these acceptance checks.

### Story 1.5: Target connection and honest build/source context

As a homelab operator,
I want to test a Caddy connection and inspect compatibility evidence,
So that I can distinguish reachable from safely editable.

**Depends on:** 1.1. **Requirements:** FR-1, FR-7, FR-8. **Legacy source:** 1.2 (pre-correction IDs).

**Acceptance Criteria:**

**1.5.AC1** — **Given** a reachable Caddy 2.9.1 reference target, **when** connection is tested, **then** GET /config/ yields JSON and exact root ETag; status is connected-unadopted and live configuration is unchanged.

**1.5.AC2** — **Given** operator-supplied version/modules/image evidence or no build manifest, **when** compatibility is displayed, **then** the evidence provenance is recorded or unknown is shown; config export is never an invented module/version inventory.

**1.5.AC3** — **Given** text/caddyfile adaptation returning a result wrapper or errors, **when** the adapter decodes it, **then** result/warnings are parsed and errors sanitized; adaptation is never claimed to provision the config.

**1.5.AC4** — **Given** stock and reference DNS-enabled builds, **when** compatibility harness runs, **then** root conditional writes, null/empty bootstrap, missing module and provisioning failure contracts have pinned fixtures; exact provider/build dependency lock is recorded.

**Completion boundary:** Core backend/storage/API outcome; no later UI story is a prerequisite for these acceptance checks.

### Story 1.6: Reviewed adoption without remote mutation

As a homelab operator,
I want to review and adopt existing whole configuration,
So that I understand ownership without an implicit rewrite.

**Depends on:** 1.3, 1.4, 1.5. **Requirements:** FR-1, FR-4, FR-9. **Legacy source:** 1.3 (pre-correction IDs).

**Acceptance Criteria:**

**1.6.AC1** — **Given** live B, root ETag E, exact source and preservation manifest, **when** preview is confirmed, **then** the unchanged preview is locally adopted with protected baseline/history, without any remote config write.

**1.6.AC2** — **Given** live/source/build/dependency identity changed since preview, **when** confirmation is attempted, **then** adoption is rejected with preserved input and renewed preview guidance.

**1.6.AC3** — **Given** missing/unrepresentable source or unproven mapping, **when** the operator connects, **then** raw configuration stays protected and the exact compatibility gap is shown; no synthesized original source is claimed.

**1.6.AC4** — **Given** empty Caddy or a target change with a prior draft, **when** adoption/target switch occurs, **then** bootstrap awaits reviewed apply, previous draft/history are detached and retained, and new generation invalidates old review/jobs.

**Completion boundary:** Core backend/storage/API outcome; no later UI story is a prerequisite for these acceptance checks.

### Story 1.7: Revisioned durable shared draft

As a homelab operator,
I want to save multiple edits across sessions without publication,
So that my live services stay unchanged until review and apply.

**Depends on:** 1.6. **Requirements:** FR-2, FR-5, FR-6. **Legacy source:** 3.1 backend (pre-correction IDs).

**Acceptance Criteria:**

**1.7.AC1** — **Given** an adopted target and matching expected draft revision, **when** route or global edits are saved, **then** only a protected singleton draft updates atomically, with incremented revision and no Caddy mutation.

**1.7.AC2** — **Given** desktop and phone save from the same revision, **when** the later save arrives, **then** it receives 409 before persistence and its unsaved buffer remains available.

**1.7.AC3** — **Given** a manager/browser restart and a saved invalid editor buffer, **when** state reloads, **then** candidate, baseline ETag, source and invalid input remain protected and recoverable.

**1.7.AC4** — **Given** SSE event gaps, reconnect or new target generation, **when** the client synchronizes, **then** it fetches authoritative projections, keeps unsaved input separate and never silently replays offline changes.

**Completion boundary:** Core backend/storage/API outcome; no later UI story is a prerequisite for these acceptance checks.

### Story 1.8: Bound source/runtime review and conflict-safe rebase

As a homelab operator,
I want to see exactly what will change and resolve external drift,
So that my edits cannot conceal unrelated changes.

**Depends on:** 1.7. **Requirements:** FR-4, FR-5, FR-11. **Legacy source:** 3.3 core; 3.4 core (pre-correction IDs).

**Acceptance Criteria:**

**1.8.AC1** — **Given** a draft and current compatible source/runtime baseline, **when** review is requested, **then** a token binds target/generation/revision/root ETag/candidate/source/build/dependency/preservation identity; source and root JSON deltas are safely displayed.

**1.8.AC2** — **Given** any bound identity changes, **when** an old review is submitted, **then** it is invalidated; candidate edits are preserved and a new review is required.

**1.8.AC3** — **Given** baseline B, external L and candidate C, **when** rebase runs, **then** nonoverlap prepares a candidate only through verified source mapping; ancestor/deletion/array overlaps require explicit resolution.

**1.8.AC4** — **Given** live cannot be represented safely after drift, **when** discard or rebase is selected, **then** no ETag-only reset or fabricated Caddyfile occurs; existing inputs remain protected with source-required guidance.

**Completion boundary:** Core backend/storage/API outcome; no later UI story is a prerequisite for these acceptance checks.

### Story 1.9: Durable guarded apply and immediate restore

As a homelab operator,
I want to apply a reviewed candidate and recover safely after interruption,
So that a competing writer or uncertain result cannot cause a blind overwrite.

**Depends on:** 1.8. **Requirements:** FR-2, FR-5, FR-9. **Legacy source:** 3.5 core; 6.1 core; 6.3 core (pre-correction IDs).

**Acceptance Criteria:**

**1.9.AC1** — **Given** an approved review and no unresolved attempt, **when** apply/restore begins, **then** a transaction commits pre-write recovery, immutable endpoint/transport descriptor, historical enabled metadata and protected attempt material while rechecking draft CAS/acquiring a durable draft/target lock, then submission_possible before any root mutation; local preparation failure sends nothing; concurrent draft saves or endpoint changes are rejected until finalization/reconciliation.

**1.9.AC2** — **Given** an external write inserted after final GET, **when** POST /config/ uses reviewed If-Match, **then** remote 412 becomes app 409 without overwrite, no /load fallback or automatic retry; restore uses identical target/generation/scope checks.

**1.9.AC3** — **Given** process death before/after send or response, readback failure or local finalization failure, **when** the manager resumes, **then** unfinished attempts lock mutations, use the original immutable endpoint descriptor and retain durable comparison material even if recording the later error failed.

**1.9.AC4** — **Given** candidate/baseline/other/unavailable live content, **when** reconciliation runs, **then** versioned structural equality determines current presence/conflict/unresolved; never compare Caddy ETag to a manager hash or infer historical causation.

**1.9.AC5** — **Given** 100 retained entries or a restored history point, **when** write/restore is prepared, **then** oldest safely unpinned history is pruned transactionally or capacity blocks; pre-restore state and referenced secret versions remain recoverable.

**1.9.AC6** — **Given** authoritative rejection or an ambiguous server/proxy failure, **when** result is finalized, **then** known rejection is distinct from uncertain remote effect; candidate remains available and readback is required where effect is unproven.

**Completion boundary:** Core backend/storage/API outcome; no later UI story is a prerequisite for these acceptance checks.

### Story 1.10: First service workspace with verification and recovery

As a homelab operator,
I want to publish a public-HTTPS proxy through the complete safe flow,
So that I can see the bounded result and immediately recover.

**Depends on:** 1.9. **Requirements:** FR-1, FR-2, FR-6, FR-8, FR-9, FR-11. **Legacy source:** 1.4 shell; 2.2 proxy subset; 3.2 minimum; 4.2 baseline; 5.1 minimum; 5.4 minimum (pre-correction IDs).

**Acceptance Criteria:**

**1.10.AC1** — **Given** the manager-only/external-Caddy/reference-service fixture, **when** the operator adopts, adds a public HTTPS proxy, saves, reviews and applies, **then** the minimal Quiet Workshop workspace completes the entire flow without immediate publication, and acceptance differs from bounded manager verification.

**1.10.AC2** — **Given** a bad service change or interrupted apply, **when** the operator uses immediate reviewed recovery, **then** guarded restore and reconciliation are available through a minimal UI before this epic is complete.

**1.10.AC3** — **Given** independent manager/Caddy restarts and manager recreation, **when** the fixture restarts with documented persistence, **then** applied route, enabled state, drafts/recovery and current observed revision survive; required browser/public checks not performed are unknown.

**1.10.AC4** — **Given** the new frontend state dependency and core controls, **when** the workspace is built and used by keyboard/phone, **then** a compatible Signals version is verified and locked; labelled focusable controls, text status and non-overlapping change bar work; no unverified dependency pin is assumed.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

## Epic 2: Complete route and global configuration editing

The operator can maintain all supported route/global/server settings through forms and Caddyfile, with validated manual/Docker destinations and a responsive searchable dashboard.

### Story 2.1: Complete route forms and retained feature validation

As a homelab operator,
I want to edit proxies, redirects and file servers with all prior settings,
So that routine management preserves existing capabilities.

**Depends on:** 1.10. **Requirements:** FR-2, FR-3, FR-4, FR-11. **Legacy source:** 2.2 remaining (pre-correction IDs).

**Acceptance Criteria:**

**2.1.AC1** — **Given** existing proxies including multiple upstreams/load balancing, host/path/prefix, request/response headers and encoding, **when** forms are edited/saved/reviewed, **then** all retained values survive source/runtime round-trip with explicit draft/apply effect.

**2.1.AC2** — **Given** redirect and file-server settings including browse/index/hide/precompression, **when** forms are saved, **then** supported codes/URI/root and Caddy-filesystem boundaries are validated without losing valid input.

**2.1.AC3** — **Given** IPv4, IDNA hosts, wildcard, bracketed IPv6 and malformed destinations, **when** fields are validated, **then** contract X3 valid/invalid boundaries and transport constraints hold; no unsupported path/userinfo is silently dropped.

**2.1.AC4** — **Given** disabled/local-only or ordered routes, **when** they are edited/reordered/enabled/deleted, **then** source/identity, immutable enabled metadata and ordering remain stable and every operation is staged, with regression checks for all legacy handler fields.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 2.2: Peer Caddyfile editor and guarded form switching

As a homelab operator,
I want to choose forms or Caddyfile over one source bundle,
So that switching views never drops my edits or source.

**Depends on:** 2.1. **Requirements:** FR-3, FR-4, FR-11. **Legacy source:** 2.5 (pre-correction IDs).

**Acceptance Criteria:**

**2.2.AC1** — **Given** one candidate source bundle, **when** the operator alternates forms and Caddyfile, **then** both views edit the same scope/revision without dropping comments/imports/unknown data.

**2.2.AC2** — **Given** invalid editor text or unsafe source mapping, **when** switch/save is attempted, **then** text remains protected, errors link to source spans, unsafe form conversion blocks and saved draft remains distinct from publication.

**2.2.AC3** — **Given** rapid text entry and explicit blur/switch/save, **when** validation runs, **then** remote adapt uses a 500ms debounce and coalesced explicit triggers, parses result/warnings and discards responses for stale revisions.

**2.2.AC4** — **Given** supported settings and secret references, **when** editor text is round-tripped, **then** reference tokens are validated server-side and no display mask becomes a credential; keyboard/textual alternatives and focus escape work.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 2.3: Optional Docker suggestions with manual destinations

As a homelab operator,
I want to choose Docker destinations or enter them manually,
So that discovery saves effort without becoming a dependency.

**Depends on:** 2.1. **Requirements:** FR-3, FR-11. **Legacy source:** 2.3 (pre-correction IDs).

**Acceptance Criteria:**

**2.3.AC1** — **Given** a configured restricted Docker metadata proxy, **when** suggestions are requested, **then** only GET container metadata is used; suggestions label network context and do not assert Caddy reachability.

**2.3.AC2** — **Given** absent, denied or unreachable Docker, **when** the editor opens, **then** an availability state plus empty suggestions appears and manual host:port remains usable.

**2.3.AC3** — **Given** an accessible raw socket mounted read-only, **when** deployment instructions are read, **then** they explain that mount flags do not enforce API permissions and recommend the GET-only proxy.

**2.3.AC4** — **Given** a suggested container alias/port, **when** it is chosen, **then** it becomes a validated editable draft destination, never an automatic publication or container lifecycle action.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 2.4: Searchable responsive route dashboard

As a homelab operator,
I want to find and inspect routes on desktop and phone,
So that technical values and pending states remain clear.

**Depends on:** 2.1. **Requirements:** FR-6, FR-11. **Legacy source:** 2.4 (pre-correction IDs).

**Acceptance Criteria:**

**2.4.AC1** — **Given** routes with long domains/destinations, **when** the dashboard renders, **then** name/host/action/enabled/live versus pending/status text remain readable and selectable without hover-only actions.

**2.4.AC2** — **Given** name/domain/path/destination filters, **when** the operator navigates away/back, **then** query state remains and empty configuration differs from zero matches.

**2.4.AC3** — **Given** 320px viewport, 200% zoom and keyboard/touch navigation, **when** routes are inspected, **then** reflow avoids viewport overflow, focus remains visible, targets are at least 44px and selected tint is not health.

**2.4.AC4** — **Given** offline or stale/removed selected route, **when** the view refreshes, **then** last known data carries timestamps/revisions and pending edits never replace live facts.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 2.5: Complete supported global/server forms and coverage inventory

As a homelab operator,
I want to edit previous global/server settings through forms and Caddyfile,
So that all supported settings participate in the same review workflow.

**Depends on:** 2.2. **Requirements:** FR-3, FR-4, FR-8, FR-11. **Legacy source:** new global/server forms (pre-correction IDs).

**Acceptance Criteria:**

**2.5.AC1** — **Given** every setting registered for the pinned supported build, including global/server/snippets/matchers, **when** coverage inventory is generated, **then** each has a dedicated or structured form, CST codec/runtime mapper, validation, secret policy and valid/invalid fixture; missing entries fail supported-build completion.

**2.5.AC2** — **Given** existing global/server options or shared snippets, **when** forms or Caddyfile change them, **then** the same candidate is produced and review names affected scopes/routes while preserving unrelated raw paths.

**2.5.AC3** — **Given** rare supported settings, **when** the structured directive form is used, **then** ordered labelled argument fields, nested rows and matchers are editable without requiring editor-only work.

**2.5.AC4** — **Given** new/unknown modules or unsafe native-JSON representation, **when** compatibility is checked, **then** gaps stay explicit/protected and unadvertised until descriptors/codecs/fixtures exist; no silent narrowing of all-settings goal.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

## Epic 3: Complete review, shared-draft and reconciliation experience

The operator can manage unsaved/saved work, inspect full redacted differences, resolve drift and stale sessions, and recover interrupted operations without blind retry.

### Story 3.1: Complete change bar and unsaved-input navigation

As a homelab operator,
I want to see saved versus unsaved changes while navigating,
So that I cannot confuse local input with publication.

**Depends on:** 1.10. **Requirements:** FR-5, FR-6, FR-11. **Legacy source:** 3.2 remaining (pre-correction IDs).

**Acceptance Criteria:**

**3.1.AC1** — **Given** unsaved input, saved draft or a clean state, **when** the Change bar renders, **then** the state/count and save/review actions name their actual effect without publication on blur/navigation.

**3.1.AC2** — **Given** unsaved edits and navigation or Escape, **when** the operator leaves/closes transient UI, **then** save/discard/cancel preserves intent; Escape never silently discards.

**3.1.AC3** — **Given** mobile soft keyboard, 320px/zoom or unavailable backdrop filter, **when** the Change bar renders, **then** it wraps without covering focused fields/actions and uses opaque fallback.

**3.1.AC4** — **Given** an in-flight action, **when** the operator clicks again, **then** duplicate submission is prevented and progress/disabled reason stays visible.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 3.2: Accessible redacted diff across full configuration

As a homelab operator,
I want to review source and runtime effects in a readable diff,
So that I know the exact scope without exposing credentials.

**Depends on:** 3.1, 2.5. **Requirements:** FR-4, FR-5, FR-11. **Legacy source:** 3.3 UI (pre-correction IDs).

**Acceptance Criteria:**

**3.2.AC1** — **Given** full root settings and source changes, **when** diff renders, **then** affected-area summary and readable line-numbered signs/text show both source delta and actual runtime paths, including preserved scope.

**3.2.AC2** — **Given** secret-bearing or unclassified opaque data, **when** diff/copy/export is requested, **then** only safe sanitized summaries/reference changes appear, and hidden scope is explicitly stated.

**3.2.AC3** — **Given** keyboard/screen-reader/reduced-motion use, **when** the operator reviews, **then** text alternative, focus restoration and non-color additions/removals preserve meaning.

**3.2.AC4** — **Given** stale review or dependency identity, **when** Apply is rendered, **then** it is disabled with explanation and renewal preserves the draft.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 3.3: Explicit drift and rebase conflict-resolution UI

As a homelab operator,
I want to resolve external changes against my candidate,
So that I can review a safe merged result without force-overwriting.

**Depends on:** 3.2. **Requirements:** FR-5, FR-11. **Legacy source:** 3.4 UI (pre-correction IDs).

**Acceptance Criteria:**

**3.3.AC1** — **Given** external drift, **when** the conflict UI opens, **then** baseline/live/candidate differences and affected mapped paths are shown without exposing secrets.

**3.3.AC2** — **Given** overlapping edits or indivisible arrays, **when** resolution is attempted, **then** operator chooses a valid compatible resolution and performs renewed review; no force/ETag-only rebase exists.

**3.3.AC3** — **Given** missing new source association, **when** rebase/discard is selected, **then** the precise source/compatibility requirement is shown while old input remains recoverable.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 3.4: Cross-session stale-save and reconnect experience

As a homelab operator,
I want to continue a shared draft from desktop or phone,
So that other session events cannot overwrite unsaved input.

**Depends on:** 3.1. **Requirements:** FR-5, FR-6, FR-11. **Legacy source:** 3.1 UI (pre-correction IDs).

**Acceptance Criteria:**

**3.4.AC1** — **Given** desktop/phone concurrent edits, **when** the server rejects a stale save, **then** the client retains local input and offers inspect latest/reconcile/resave rather than overwrite.

**3.4.AC2** — **Given** SSE reconnect with gaps or an event for an old target, **when** state sync runs, **then** latest projections are fetched, old-generation updates ignored and unsaved text kept.

**3.4.AC3** — **Given** offline edits, **when** connection returns, **then** no automatic replay/publication occurs and input can be explicitly reconciled.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 3.5: Interrupted-operation reconciliation experience

As a homelab operator,
I want to understand and resolve an interrupted apply or restore,
So that I know the currently observed state before another write.

**Depends on:** 3.2. **Requirements:** FR-5, FR-6, FR-9, FR-11. **Legacy source:** 3.5 UI (pre-correction IDs).

**Acceptance Criteria:**

**3.5.AC1** — **Given** an unfinished apply or restore after restart, **when** Changes/result loads, **then** a persistent reconciliation-required panel names known effect and locked operations.

**3.5.AC2** — **Given** live candidate, baseline, other or unavailable, **when** the operator reconciles, **then** UI presents current presence/renew-review/conflict/unresolved outcomes accurately and never offers blind retry.

**3.5.AC3** — **Given** readback/local commit failure after acceptance, **when** feedback renders, **then** it distinguishes accepted response from unfinished local state and retains attempt/recovery references.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

## Epic 4: Complete all four TLS scenarios

The operator can configure public automatic HTTPS, DNS-01/wildcards, internal CA and HTTPS upstreams with protected credentials, honest build evidence and explicit trust controls.

### Story 4.1: Complete public HTTPS settings and lifecycle evidence

As a homelab operator,
I want to configure public automatic HTTPS and observe its lifecycle,
So that I can diagnose issuance and renewal separately from config acceptance.

**Depends on:** 1.10, 2.5, 3.5. **Requirements:** FR-8, FR-6, FR-7. **Legacy source:** 4.2 remaining (pre-correction IDs).

**Acceptance Criteria:**

**4.1.AC1** — **Given** a public-domain route, **when** TLS forms open, **then** automatic public HTTPS is primary with explicit optional email/global setting impact and 80/443/DNS prerequisites.

**4.1.AC2** — **Given** configuration accepted but certificate not observed, **when** status renders, **then** it remains applied/unverified or issuance-pending/unknown, not certificate-ready.

**4.1.AC3** — **Given** the reference public HTTPS environment, **when** issuance, renewal, failure diagnosis and restart scenarios run, **then** certificate/response observations record exact environment/time/vantage; simulated ACME renewal differs from external public evidence.

**4.1.AC4** — **Given** existing global TLS settings, **when** a route is added, **then** they remain intact unless explicitly reviewed changes require otherwise.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 4.2: DNS-01 and wildcards with build and secret evidence

As a homelab operator,
I want to use DNS-01 for private public-domain services and wildcards,
So that certificates can be issued with supported modules and protected tokens.

**Depends on:** 4.1. **Requirements:** FR-8, FR-4. **Legacy source:** 4.3 (pre-correction IDs).

**Acceptance Criteria:**

**4.2.AC1** — **Given** wildcard/private public-domain DNS-01 using the Cloudflare reference adapter, **when** a provider is selected, **then** its pinned module/build manifest, structured fields and versioned secret references are required; DNS service-record setup remains external.

**4.2.AC2** — **Given** missing module, unknown build or provisioning failure, **when** validation/application runs, **then** narrow evidence or actionable failure is shown without claiming adaptation guarantees runtime success.

**4.2.AC3** — **Given** credentials rotate after a recovery point, **when** adapt/apply/restore occurs, **then** the selected immutable version is hydrated transiently and absent material blocks safely.

**4.2.AC4** — **Given** stock and DNS-enabled reference builds, **when** compatibility scenarios run, **then** actual target tests cover module errors, challenge/certificate evidence and restart, without requiring every arbitrary provider.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 4.3: Internal CA with explicit root and client-trust workflow

As a homelab operator,
I want to encrypt local-name services and install client trust,
So that I understand why manager trust differs from my device.

**Depends on:** 4.1. **Requirements:** FR-8, FR-6, FR-7. **Legacy source:** 4.4 (pre-correction IDs).

**Acceptance Criteria:**

**4.3.AC1** — **Given** local-name internal TLS, **when** it is configured, **then** tls internal is staged/reviewed and client trust is separate from manager certificate observations.

**4.3.AC2** — **Given** explicit local-root download request, **when** PKI chain is fetched, **then** the real endpoint/chain format is parsed, the root is validated/selected, and possible default-CA provisioning is disclosed.

**4.3.AC3** — **Given** a client without installed trust, **when** instructions/results display, **then** OS-specific trust guidance is available and manager trust does not mark client trusted.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 4.4: Independent HTTPS upstream transport and trust

As a homelab operator,
I want to connect to HTTPS backends with explicit SNI and trust,
So that upstream security is independent of client-facing TLS.

**Depends on:** 4.1. **Requirements:** FR-8, FR-3. **Legacy source:** 4.5 (pre-correction IDs).

**Acceptance Criteria:**

**4.4.AC1** — **Given** an HTTPS backend, **when** transport fields are edited, **then** scheme/SNI/trust roots and multi-upstream compatibility validate independently of frontend TLS.

**4.4.AC2** — **Given** verification bypass is explicitly enabled, **when** draft is saved/reviewed, **then** the insecure trust consequence is clearly stated with no silent bypass.

**4.4.AC3** — **Given** HTTP and HTTPS destinations plus malformed URL paths/userinfo, **when** validation runs, **then** supported X3 grammar holds and preserved source/settings survive round-trip.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

## Epic 5: Typed verification and evidence-based diagnosis

The operator can interpret revision-bound manager/client observations, retry safely and act on bounded diagnostics for the chosen route type.

### Story 5.1: Complete asynchronous revision-bound manager probes

As a homelab operator,
I want to observe DNS, TCP, TLS, upstream and published-service facts,
So that each result identifies exactly what was checked.

**Depends on:** 4.2, 4.3, 4.4. **Requirements:** FR-6, FR-7, FR-8. **Legacy source:** 5.1 remaining (pre-correction IDs).

**Acceptance Criteria:**

**5.1.AC1** — **Given** an applied revision, **when** asynchronous probes execute, **then** jobs/results bind target/generation/revision/expectation and separate DNS/TCP/TLS/direct-upstream/published-URL evidence.

**5.1.AC2** — **Given** protected service expectations or near-expiry certificates, **when** checks run, **then** allowed statuses may be 401/403; near-expiry warning differs from invalid certificate, SAN/trust and time are recorded.

**5.1.AC3** — **Given** late result after edit/restore, timeout or retry, **when** result arrives, **then** old revision remains historical; 10s deadlines, 5/15s transient retries and five-minute/revision staleness follow the contract; retry never reapplies.

**5.1.AC4** — **Given** unknown normalization or inaccessible upstream from manager network, **when** results finalize, **then** unknown/unavailable stays bounded and cannot become global service success.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 5.2: Redirect and known-resource file checks

As a homelab operator,
I want to verify behavior appropriate to the route type,
So that a generic HTTP badge cannot hide a wrong redirect or resource.

**Depends on:** 5.1, 2.1. **Requirements:** FR-6, FR-7. **Legacy source:** 5.2 (pre-correction IDs).

**Acceptance Criteria:**

**5.2.AC1** — **Given** a redirect expectation, **when** the first hop is requested, **then** configured code and Location policy are compared without following redirects; destination verification is separate.

**5.2.AC2** — **Given** file server with known test resource, **when** it is requested, **then** expected response/path is checked without claiming directory mount/permission knowledge from status alone.

**5.2.AC3** — **Given** unexpected 502/403 or Location, **when** diagnosis is produced, **then** the observed mismatch is stated and possible causes remain hypotheses unless independently established.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 5.3: Bounded browser observations and scoped CSP

As a homelab operator,
I want to compare supported client evidence with manager evidence,
So that browser policy failures do not become false service failures.

**Depends on:** 5.1. **Requirements:** FR-6, FR-7, FR-11. **Legacy source:** 5.3 (pre-correction IDs).

**Acceptance Criteria:**

**5.3.AC1** — **Given** a selected published URL and supported origin policy, **when** browser check runs, **then** only its explicit origin is allowed in CSP and method/vantage/revision is recorded, without manager Admin API exposure.

**5.3.AC2** — **Given** CORS denial, opaque response, CSP/mixed-content or trust failure, **when** fetch returns/fails, **then** unavailable/unknown is shown when expected status/headers cannot be observed and manual open-and-check guidance remains.

**5.3.AC3** — **Given** manager reaches upstream Docker alias while client uses published host, **when** evidence is displayed, **then** the two targets/vantages are never equated; raw certificate inspection remains manager-only.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 5.4: Complete evidence, stale-result and retry interface

As a homelab operator,
I want to read independent timestamped results and retry checks,
So that old observations cannot verify a new configuration.

**Depends on:** 5.2, 5.3, 2.4. **Requirements:** FR-6, FR-7, FR-11. **Legacy source:** 5.4 remaining (pre-correction IDs).

**Acceptance Criteria:**

**5.4.AC1** — **Given** route/detail/result UI, **when** evidence displays, **then** acceptance, TLS, upstream and published/client observations have independent labels/method/time/vantage/revision and unavailable states.

**5.4.AC2** — **Given** a late job for a previous applied revision, **when** it completes, **then** it cannot mark current service verified and is labelled stale/historical.

**5.4.AC3** — **Given** Retry verification, **when** it is selected, **then** only checks restart against current expected revision, no config apply occurs, and async announcements do not steal focus.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 5.5: Evidence-based diagnosis and remediation

As a homelab operator,
I want to see observations, likely causes and relevant next checks,
So that I can correct problems without trusting unsupported diagnosis.

**Depends on:** 5.4. **Requirements:** FR-7, FR-11. **Legacy source:** 5.5 (pre-correction IDs).

**Acceptance Criteria:**

**5.5.AC1** — **Given** DNS/TCP/TLS/HTTP evidence, **when** diagnosis is rendered, **then** facts, candidate causes and two relevant next checks are distinguishable; 502 or 403 alone never establishes root cause.

**5.5.AC2** — **Given** unavailable client vantage or timeout, **when** help is opened, **then** known observations are preserved and the missing evidence/next observation is named.

**5.5.AC3** — **Given** a correction from phone/desktop, **when** settings are changed, **then** it follows the same save/review/apply/verify lifecycle rather than diagnosis retry publishing changes.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

## Epic 6: Scoped history, complete backup and maintenance

The operator can inspect/export/compare/restore scoped history and recover the whole deployment with tested upgrade and remote-asset runbooks.

### Story 6.1: Complete scoped history list and export

As a homelab operator,
I want to inspect and safely export retained configurations,
So that I know which target and recovery scope each entry covers.

**Depends on:** 1.10. **Requirements:** FR-9, FR-11. **Legacy source:** 6.1 remaining (pre-correction IDs).

**Acceptance Criteria:**

**6.1.AC1** — **Given** full-config and legacy route-array entries, **when** History loads, **then** target/time/reason/scope/exclusions/retention appear distinctly; empty/loading/denied states are usable.

**6.1.AC2** — **Given** sanitized export, **when** the operator downloads it, **then** secret references remain safe and the export warns that standalone full-system/secret recovery needs the protected package.

**6.1.AC3** — **Given** retention or referenced historical credentials, **when** history changes, **then** active recovery points/versions remain pinned or write capacity blocks, never silently invalidating restore.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 6.2: Historical comparison against current live state

As a homelab operator,
I want to compare a history point to the live target,
So that I can inspect changes before preparing restore.

**Depends on:** 6.1, 3.2. **Requirements:** FR-9, FR-11. **Legacy source:** 6.2 (pre-correction IDs).

**Acceptance Criteria:**

**6.2.AC1** — **Given** a compatible point and current live target, **when** Compare is requested, **then** a sanitized source/runtime comparison captures exact live ETag and target/generation in a new review.

**6.2.AC2** — **Given** scope/target mismatch or live drift, **when** comparison is refreshed, **then** the mismatch blocks restore or renews review without reusing the snapshot-era ETag.

**6.2.AC3** — **Given** opaque/unknown source areas, **when** diff is displayed, **then** preserved scope/limits remain explicit rather than fabricated equivalent source.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 6.3: Complete reviewed history restore experience

As a homelab operator,
I want to restore a selected compatible history point,
So that recovery has the same review, guard and evidence as ordinary apply.

**Depends on:** 6.2, 3.5. **Requirements:** FR-5, FR-9, FR-11. **Legacy source:** 6.3 remaining (pre-correction IDs).

**Acceptance Criteria:**

**6.3.AC1** — **Given** a selected full or compatible legacy-scope point, **when** restore confirmation is shown, **then** target/scope/exclusions and effect are explicit, then the shared guarded durable attempt creates immediate pre-restore recovery.

**6.3.AC2** — **Given** live root changes after comparison, **when** restore submits, **then** If-Match rejects without overwrite and a fresh review is required.

**6.3.AC3** — **Given** lost response, missing secret key/version or target mismatch, **when** restore runs, **then** shared reconciliation/locking or pre-submission blocking applies and a successful restore gets a new applied revision plus separate verification.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

### Story 6.4: Complete backup, upgrade and maintenance runbooks

As a homelab operator,
I want to back up and restore the whole deployment predictably,
So that config, local intent, secrets and certificates survive maintenance.

**Depends on:** 6.3, 5.5, 1.3. **Requirements:** FR-2, FR-9, FR-10. **Legacy source:** 6.4 remaining (pre-correction IDs).

**Acceptance Criteria:**

**6.4.AC1** — **Given** consistent SQLite plus key/current source/config/startup/import/env/build/external data assets, **when** backup is acquired with operator-side remote-host commands, **then** a manifest validates scope/time/provenance and missing pieces mark it incomplete; no manager remote write privileges are added.

**6.4.AC2** — **Given** a clean disposable environment and complete protected package, **when** documented restore runs, **then** manager/target configuration, local-only intent, settings, credentials and certificates recover with explicit source association and independent restart checks.

**6.4.AC3** — **Given** prior-release DB/image and interrupted upgrade, **when** maintenance procedure runs, **then** functional mapping/rollback/downgrade limitations from 1.3 are exercised; current applied config is captured, not inferred from old history.

**6.4.AC4** — **Given** a release scenario matrix, **when** release evidence is compiled, **then** every required flow names exact versions/environment and automated/manual evidence, with no budget or universal reachability claim unsupported by measurements.

**Completion boundary:** Extends the already guarded shared pipeline; this story cannot bypass safety, conceal unsupported scope or claim unexecuted acceptance evidence.

## Legacy-to-current aliases

Old IDs refer to the pre-correction 28-story file. New IDs are authoritative; identical numbers may now name a different outcome. QA test IDs may remain immutable but must have a separate current-story/AC field.

| Old story | Current allocation |
| --- | --- |
| 1.1 | 1.1, 1.3 |
| 1.2 | 1.5 |
| 1.3 | 1.6 |
| 1.4 | 1.10 |
| 2.1 | 1.4 |
| 2.2 | 1.10, 2.1 |
| 2.3 | 2.3 |
| 2.4 | 2.4 |
| 2.5 | 2.2, 2.5 |
| 3.1 | 1.7, 3.4 |
| 3.2 | 1.10, 3.1 |
| 3.3 | 1.8, 3.2 |
| 3.4 | 1.8, 3.3 |
| 3.5 | 1.9, 3.5 |
| 4.1 | 1.2 |
| 4.2 | 1.10, 4.1 |
| 4.3 | 4.2 |
| 4.4 | 4.3 |
| 4.5 | 4.4 |
| 5.1 | 1.10, 5.1 |
| 5.2 | 5.2 |
| 5.3 | 5.3 |
| 5.4 | 1.10, 5.4 |
| 5.5 | 5.5 |
| 6.1 | 1.9, 6.1 |
| 6.2 | 6.2 |
| 6.3 | 1.9, 6.3 |
| 6.4 | 1.3, 6.4 |

New explicit global/server coverage is 2.5. New source/protocol/migration acceptance is integrated into the allocations above, with no hidden prerequisite epic. See QA finding-to-AC matrix for B01–B10/C01–C05 closure intent.
