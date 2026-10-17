---
id: SPEC-caddy-orchestrator
updated: 2026-10-17
status: corrected-contract
companions:
  - glossary.md
  - architecture-diagrams.md
  - ../../planning-artifacts/ux-designs/ux-caddy-orchestrator-2026-10-16/DESIGN.md
  - ../../planning-artifacts/ux-designs/ux-caddy-orchestrator-2026-10-16/EXPERIENCE.md
  - ../../planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md
  - ../../planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/configuration-contracts.md
sources:
  - ../../planning-artifacts/prds/prd-caddy-orchestrator-2026-10-16/prd.md
  - ../../planning-artifacts/sprint-change-proposal-2026-10-17.md
---

# Caddy Orchestrator for Homelab Operators

## Contract authority and purpose

Homelab operators need to publish, change, understand and recover services on an existing Caddy without losing configuration or mistaking acceptance for working service. This specification summarizes the approved direction and [course-correction contracts](../../planning-artifacts/sprint-change-proposal-2026-10-17.md). Latest explicit user decisions govern scope; PRD records outcomes; architecture/companion define implementation rules; UX supplies interaction and labelled visual assumptions; stories/tests trace acceptance. Surface conflicts rather than invent authority from a canonical label. Current root documentation/source describes legacy scoped-route behavior, not completed expanded capability.

## Capabilities

- **CAP-1 — Connect and adopt existing configuration.** Intent: explicit whole-config ownership without disruption. Success: connected-unadopted state, live/source/build/preservation preview, protected baseline and unchanged-preview confirmation without remote mutation; missing source/incompatible mapping blocks with retained originals. Empty bootstrap is reviewed publication, not adoption.
- **CAP-2 — Complete supported forms and Caddyfile authoring.** Intent: edit previous routes/global/server settings in either peer view. Success: exact protected source/CST, stable node IDs/versioned enabled maps and declared-build settings inventory; dedicated or structured directive forms cover every supported setting with matching codecs/validation/secret classification. Imports use captured bundle semantics; unknown native JSON is preserved and unsafe conversion blocks. No automatic original-source reconstruction is promised.
- **CAP-3 — Destinations and Docker suggestions.** Intent: validated manual destinations plus optional discovery. Success: host:port/IPv4/bracketed-IPv6 and HTTPS transport boundaries follow companion X3; GET-only metadata suggestions state network context; absent Docker leaves manual input functional. No Docker lifecycle mutation or automatic publication.
- **CAP-4 — Shared durable drafts and review.** Intent: collect multiple changes across sessions. Success: expected-revision saves, protected invalid buffers, explicit unsaved-versus-saved state, SSE resync without input overwrite and redacted source/root-runtime review bound to target/generation/revision/source/dependencies/preservation/ETag. Save never publishes.
- **CAP-5 — Guarded publication and conflict resolution.** Intent: reviewed batch changes without blind overwrite. Success: durable pre-write recovery/attempt with immutable connection identity and target/draft mutation lock, one root POST /config/ with reviewed If-Match, remote 412→app 409, and three-way conflict handling with new source association/review. Never use unconditional /load as a fallback. This is one guarded remote mutation, not a Caddy/SQLite transaction or automatic rollback promise.
- **CAP-6 — Revision-bound typed verification.** Intent: distinguish acceptance from observed operation. Success: separate config/TLS/upstream/published-service observations with target/applied revision/expectation/time/method/vantage; route-specific expected statuses/resources/Location/trust; late results historical; browser restrictions unknown/unavailable; retry never applies.
- **CAP-7 — All four TLS scenarios.** Intent: public automatic HTTPS, DNS-01/wildcards, internal CA and HTTPS upstreams. Success: supported-build evidence, protected immutable credentials and explicit trust/SNI controls with real scenario evidence. Public HTTPS remains primary. Cloudflare is the first reference DNS adapter; arbitrary providers require manifests/descriptors/fixtures. PKI root download is explicit and explains possible default-CA provisioning.
- **CAP-8 — Honest diagnosis and interrupted-operation recovery.** Intent: useful next checks and safe ambiguous-state handling. Success: facts separated from possible causes; durable submission-possible attempt survives process death/response loss/readback/local commit failure; startup locks until content reconciliation establishes current candidate/baseline/other/unavailable state. ETag change alone proves neither candidate equality nor historical causation.
- **CAP-9 — Scoped history and complete maintenance.** Intent: compare/export/restore known configuration and recover deployment after upgrade/failure. Success: protected full-config versus legacy-route scopes; reviewed guarded restore/new applied revision; at most 100 retained history points with active pins; historical secret hydration; functional migration of disabled/local-only data/settings/history; consistent DB/key plus external startup/config/source/import/env/build/data backup with operator-side remote acquisition and clean-restore evidence.
- **CAP-10 — Clear accessible desktop and phone operation.** Intent: Quiet Workshop for routine settings, status and quick corrections. Success: legible technical values, labelled fields/errors, keyboard/focus/textual diff, WCAG 2.2 AA target, 44px touch, 320px reflow, 200% zoom and reduced motion/transparency. Exact visual proposals remain labelled; a generated document is not a completed usability/accessibility audit.

## Constraints

Manager-only Compose connects to external Caddy; bundled target/reference service is a test fixture. Preserve Go/Gin/SQLite/Preact/Tailwind and single-operator/single-instance boundary. Root adoption enforces one writer mode so legacy CRUD cannot bypass the new lifecycle. External persistence requires documented resume/equivalent and config/data storage; API access alone cannot prove host mounts/flags. Admin API stays private. Keys, exact source/runtime/legacy payloads and historical credentials use the protected storage contract before import. Current reference target is Caddy 2.9.1; target adaptation is POST /adapt with text/caddyfile and decoded result/warnings, not a semantic provisioning dry run or module inventory endpoint. No Caddy binary is bundled in the manager.

## Non-goals

Fleet/cluster control; RBAC/SSO; DNS service-record provisioning; container lifecycle actions; continuous traffic analytics; automatic NPM migration; new frontend/backend stack. No blanket arbitrary-plugin compatibility or universal Caddyfile inversion claim. These limits do not silently remove the confirmed all-settings forms or four TLS scenarios within declared compatibility.

## Success signal and handoff

Following documented manager-only deployment, an operator adopts compatible existing config without rewrite, stages a reference proxy, reviews exact source/runtime scope, publishes through the root precondition, sees bounded HTTPS/service evidence, changes it and restores safely. Independent restart, competing writer, interrupted operation, stale session/probe, credential rotation, prior-release migration and clean backup restore have recorded evidence. Required release capability cannot be declared complete while its critical acceptance remains unresolved. The first internal service slice is not the full release.

[Epics](../../planning-artifacts/epics.md) define 33 stories; [QA](../../test-artifacts/test-design-qa.md) maps acceptance. Build locks, corpus, test hooks and all runtime evidence are assigned implementation work, not already delivered. No numeric latency/memory budget or calendar release commitment is inferred.
