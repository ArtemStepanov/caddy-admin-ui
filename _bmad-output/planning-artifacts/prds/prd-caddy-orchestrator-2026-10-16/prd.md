---
title: Caddy Orchestrator for Homelab Operators
status: approved-scope
created: 2026-10-16
updated: 2026-10-17
research_status: integrated
decision_status: confirmed-scope-corrected-contracts
---

# PRD: Caddy Orchestrator for Homelab Operators

## 1. Product direction and authority

Build a dependable Caddy manager for homelabbers, comparable in role to Nginx Proxy Manager, completing **connect → publish → understand the result → change → recover**. Confirmed choices are full existing-configuration management, forms for all settings plus a Caddyfile editor, shared persistent batch drafts/diff/apply, all four TLS scenarios, manual destinations and optional Docker discovery, existing-Caddy manager-only Compose, and Quiet Workshop.

Artem approved [the course correction](../../sprint-change-proposal-2026-10-17.md) on 2026-10-17. It propagates choices already settled in UX/architecture and repairs technical contracts; it does not establish running product acceptance. Scope is confirmed; detailed visual assumptions remain labelled in UX. Numeric performance targets, release date/version and complete NPM parity are not supplied.

Latest explicit user decisions govern scope. This PRD records outcomes; [architecture and its contracts](../../architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md) fix implementation mechanisms; [SPEC](../../../specs/spec-caddy-orchestrator/SPEC.md), epics and tests must agree with both. A conflict is surfaced rather than decided by a blanket canonical label. [CLAUDE.md](../../../../CLAUDE.md) and current source describe today's scoped-route implementation; they are not the final product boundary. [Research](../../research/competitive-caddy-homelab-manager-2026-10-16/research.md) is complete and advisory. No runtime checks were performed during this correction.

## 2. Target user and jobs

The confirmed audience is homelabbers broadly, including Artem. The operator wants to expose a home service through a memorable domain, change its destination without losing existing configuration, understand a failure, and restore service after a bad change. [ASSUMPTION A2: documented Compose, DNS and networking prerequisites are an acceptable setup path; prior Caddy JSON expertise is unnecessary for routine proxy management.] Acceptance uses a documented reproducible reference service. Artem's own services can supply additional validation, but are not a scope prerequisite. A broader audience does not imply fleet management, RBAC or multiple application users.

No real user session or observed journey has been supplied. This operator-focused capability PRD leaves named journeys to UX elicitation; it does not invent personas or present hypothetical sequences as customer evidence.

**Terms:** a Route maps hosts/paths to proxy, redirect or file-server behavior; an Upstream is its backend destination. An adopted Target is one Caddy instance under explicit root-configuration ownership. A Source Bundle holds exact protected Caddyfile/dependencies. A Draft is a revisioned local candidate; saving it does not publish. A Recovery Point has explicit full-config or legacy route-array scope and excludes external certificate keys unless a complete backup includes them. Managed Server/External Route/Snapshot identify legacy scoped-route concepts where explicitly used. Applied means configuration acceptance plus known observed revision, not working service. Failed states name the known remote effect; unknown means it cannot be established. Service verified names the typed expectations, time, revision and vantage actually observed.

## 3. Current baseline and product gap

| Area | Current source evidence | First-release work to assess |
|---|---|---|
| Route management | Proxy/redirect/file-server CRUD, enable/disable, multiple Upstreams, host/path rules, headers; dashboard search/filter | Prove real-service workflows; improve validation and forms |
| Ownership and changes | Explicit server setup preview/confirmation, scoped Routes writes, ETag conflicts, raw External Route preservation | Preserve as regression gates; explain conflicts and ambiguous outcomes |
| Recovery | Pre-change Snapshots, guarded restore/export; fixed retention of 100, list of 50 | Demonstrate recovery; distinguish Snapshots from complete backups |
| Runtime status | Admin API availability, setup state, last sync result held in memory | Distinguish stored/applied/unknown, certificate, and Upstream observations |
| HTTPS and access | Caddy can provide HTTPS; application does not manage certificate/TLS modes or host access policies | Deliver primary public HTTPS and supported choices; expose usable diagnosis without unsupported promises |
| Deployment | Local bundled-Caddy fixture and existing-Caddy Compose option, private bindings, optional single Basic Auth account with enforced listener rules | Production: manager-only Compose connects to existing Caddy; verify install/restart/update/backup |

Sources: current [CLAUDE.md](../../../../CLAUDE.md), [README.md](../../../../README.md), and [code extraction](reconcile-code-baseline.md). Historical `specs/001` and `specs/002` are not the current API contract. Source presence does not establish runtime correctness. In particular, Caddy startup persistence and public/LAN HTTPS remain unverified. See FR-2 and FR-8.

Research supports the workflow proposal, not a uniqueness claim: NPM and several Caddy managers already document proxy/TLS convenience, diagnosis and recovery. Three intentionally selected recent NPM accounts illustrate vague save errors, hard-to-find certificate diagnostics and a reported mismatch between renewal and served certificates. They support testing requested/applied/observed states; they do not establish prevalence or Caddy defects. See the [competition](../../research/competitive-caddy-homelab-manager-2026-10-16/digests/competition-r1-1.md), [user-voice](../../research/competitive-caddy-homelab-manager-2026-10-16/digests/user-voice-r1-1.md) and [native constraints](../../research/competitive-caddy-homelab-manager-2026-10-16/digests/caddy-native-r1-1.md) extracts.

## 4. Confirmed scope boundary

FR-1–FR-11 are the release scope. Preserve existing route functionality while expanding to full adopted configuration. Required TLS scenarios: public automatic HTTPS; DNS-01/wildcards; local names/internal CA with trust guidance; and independently configured HTTPS upstreams. Public HTTPS is the primary workflow. Optional Docker discovery supplements manual input; discovery availability cannot become a setup prerequisite.

Forms cover all settings of each declared supported build, including global/server options. Dedicated common-task controls and structured directive/matcher forms are peer projections of the Caddyfile CST. Build/module/source compatibility is explicit; missing form descriptors or unsafe source mapping are incomplete capabilities with blocking guidance, not a silent permanent editor-only exclusion. Arbitrary native JSON cannot be assumed reconstructable as Caddyfile. Preserve incompatible content and exact originals; do not drop it to make a form work.

Non-goals remain fleet management, RBAC/SSO, DNS service-record provisioning, container lifecycle management, traffic analytics, automatic NPM migration and a stack migration. Full NPM feature parity is not promised. First internal service delivery is a safe vertical slice; it does not silently reduce release scope.

## 5. First-release requirements

Confirmed scope is recorded above. Acceptance uses manager-only Compose against a separately deployed Caddy fixture and a reproducible reference service. Planned tests and documentation are not runtime evidence.

### Setup and publication

**FR-1 — Connect and adopt existing configuration.** Test the target without adopting or publishing. Review live root configuration, source/build/dependency provenance, preservation gaps and exact ownership; confirm only an unchanged preview. Supply source or a verified converter when runtime JSON cannot provide it. Empty-target bootstrap is an explicit reviewed draft/apply, never an adoption side effect. Changing targets invalidates ownership/reviews, detaches and retains prior drafts/history, and requires renewed adoption.

**FR-2 — Publish and retain a service.** Create/edit/enable/disable/delete proxy, redirect and file-server routes through forms or Caddyfile, accumulate drafts, review the shared source/runtime delta and apply explicitly. Prove a reference published service and safe recovery. Retain multi-upstream/load balancing, host/path/prefix, headers, encoding, ordering and all existing handler settings. Accepted configuration survives independent manager/Caddy restarts and manager recreation under documented persistence. File roots/permissions resolve on Caddy, not the manager.

**FR-3 — Validate supported settings and destinations.** Field-linked validation preserves valid input. The architecture companion X3 defines host/path/IDNA/IPv6, destination, URI/root and identity boundaries with examples. Accept IP:port, service:port and bracketed IPv6:port; HTTPS upstreams have separate transport/SNI/trust settings. Reject malformed destinations before apply without silently replacing accepted config. Structured forms include previous global/server settings and descriptor validation; unknown build/grammar states are explicit compatibility gaps. Optional read-only Docker suggestions retain manual input.

### Safe changes and honest feedback

**FR-4 — Preserve the entire adopted configuration.** Manage full existing root configuration after explicit adoption. A path ownership/preservation manifest protects unrelated root and nested data, array order, listeners, other servers, PKI/logging/storage, unknown directives and exact available source. Generated changes replay only through verified mappings. Unsafe conversion, overlapping ownership and unrepresentable content block with retained protected originals and actionable scope, rather than silently becoming lost or globally read-only product scope.

**FR-5 — Handle shared, concurrent and uncertain changes.** Expected draft revisions reject stale desktop/phone saves while retaining local input. Apply and restore use the exact reviewed root ETag as the remote precondition; changed live state rejects overwrite, including a competing write after the final read. Rebase compares baseline/live/candidate and requires conflict resolution and renewed review. Persist pre-write snapshots and durable attempts before submission. Acknowledgement/readback/local-finalization failure or crash requires content reconciliation; do not promise a Caddy/SQLite transaction or infer outcome from an ETag change alone.

**FR-6 — Show what is known.** Separate saved draft, accepted/applied revision, enabled state, rejected/conflicted operation and reconciliation-required state, durably across restart. Display observation time, method, expectation, target/revision and vantage. Certificate, upstream and published-service evidence are separate; unsupported or stale observations stay unknown/unavailable/stale. Late checks cannot verify newer config. Compression/settings changes identify their draft/apply effect.

**FR-7 — Diagnose a failed publication.** Provide evidence-based next checks for unreachable Caddy, invalid/adaptation/provisioning config, DNS/client limitations, trust failures and upstream problems. An HTTP 502/403 alone does not establish its cause. Direct upstream and published URL observations are separate. Browser CORS/CSP/opaque/mixed-content/trust limits produce bounded unknown/unavailable states. Retry verifies without applying. All results state manager versus client vantage and the expectation checked.

### HTTPS, recovery and operation

**FR-8 — Cover all four TLS scenarios.** Public automatic HTTPS is primary; also deliver DNS-01/wildcards with encrypted versioned credentials and build evidence, local names/internal CA with explicit root download/client trust, and HTTPS upstream SNI/trust independent of frontend TLS. Cloudflare is the first reference DNS adapter; additional providers need descriptors/build locks/fixtures, not a new scope discovery cycle. Unknown plugin evidence remains unknown; adaptation success is not semantic provisioning validation. Prove public issuance, renewal and restart behavior, diagnose failure honestly, and preserve existing adopted TLS settings.

**FR-9 — Recover with explicit scope.** Inspect, compare, safely export and restore full-config recovery points; keep migrated legacy route-array snapshots with original server/scope. Display time, target/generation, reason, retention and exclusions. Restore is a reviewed candidate against captured live root ETag, with an immediate pre-restore point, durable attempt and separate post-restore checks. Target/scope mismatch or drift blocks overwrite. Keep up to 100 history entries while protecting active recovery evidence; full capacity of pinned entries blocks a new write. History excludes external TLS private keys and is not complete disaster backup.

**FR-10 — Back up and upgrade predictably.** A complete recovery package covers consistent SQLite, secret key/provenance, current applied configuration/source/import/environment/build/startup assets and external configuration/data or equivalent storage. Document operator-side acquisition/restoration for remote Caddy without granting manager Docker/filesystem mutation. Verify clean-environment restore and prior-release transactional migration, including disabled/local-only routes, raw route order, custom settings, legacy ownership and scoped snapshots. A table archive alone is not a working upgrade. State downgrade via prior backup/image and all missing-asset/key limits.

### Modern interface

**FR-11 — Make routine work clear and accessible.** Use Quiet Workshop with legible routes/settings, search, connected target, pending changes, shared review and recovery. Forms and Caddyfile are peers for supported settings; fields retain input and expose errors. Desktop configuration and phone status/quick corrections use the same safe draft lifecycle. Target WCAG 2.2 AA with keyboard/focus, labels, non-color status/diff, 44px touch targets, 320px reflow and 200% zoom. Exact visual tokens/layout assumptions remain in UX and require rendered checks.

## 6. Cross-cutting release constraints

- **Security:** retain authenticated non-loopback operation or an explicit insecure override, paired credentials, mutation origin/CSRF protections, private Caddy Admin API and redaction of credentials. Configuration and Snapshot exports must warn about their actual sensitive contents; never expose passwords through ordinary status/errors. These are specific regression checks, not a claim of a security audit.
- **Data safety:** failure and recovery checks demonstrate that unrelated configuration survives, conflicts prevent blind replacement, and uncertain remote results remain visible. Empty bootstrap is assessed separately from pre-change Snapshot behavior.
- **Operational evidence:** release notes label each claim by its supporting evidence: documentation, automated tests, or manual runtime checks. Test results from a previous conversation are not current release certification. Record version/environment and remaining limitations for each supported scenario.
- **Resource fit:** retain a small self-hosted deployment as a design constraint; no hardware or latency targets have been supplied. Measure a representative chosen environment before setting performance budgets (Q5).

## 7. Success and release acceptance

**SM-1 — Useful routine operation:** a homelab operator following the documented manager-only Compose/existing-Caddy workflow can publish the reproducible reference service through public HTTPS, identify a failed change, correct it and recover a previous configuration. Validate with a separately deployed Caddy fixture; Artem's lab is optional additional evidence. Validates FR-1–FR-3, FR-5–FR-11. Task completion and points requiring manual intervention are recorded; no invented completion-time target.

**SM-2 — Trustworthy operation:** each selected scenario passes restart, upgrade, conflict, External Route preservation and recovery checks with recorded evidence. Validates FR-2 and FR-4–FR-10. The release cannot claim the selected scenario works while a critical acceptance check remains unresolved.

**Counter-metric SM-C1:** fewer clicks or a shorter setup is not success if it hides ownership, uncertain application state, TLS prerequisites or restore scope. This counterbalances SM-1. Route count and feature count do not establish usefulness. No telemetry service is required; a small release evidence record is sufficient.

## 8. Decisions and implementation handoff

Confirmed decisions: broad homelabber audience; external existing Caddy with manager-only Compose; primary public HTTPS and all four requested TLS scenarios; whole existing configuration; all-settings forms plus Caddyfile; shared reviewed drafts; Docker discovery plus manual input; Quiet Workshop and desktop/phone needs. Former Q1–Q4/D1/D2 candidate scope questions are resolved by the user decisions and approved course correction. Do not reopen them.

The [33-story backlog](../../epics.md) assigns implementation of the source/preservation engine, descriptor/build inventory, guarded attempts, sessions/probes, upgrade and backup. [Architecture contracts](../../architecture/architecture-caddy-orchestrator-2026-10-17/configuration-contracts.md) define their common models and boundaries. [Test mappings](../../../test-artifacts/test-design-qa.md) provide acceptance intent. Reference provider build locks, golden corpus, hooks and runtime evidence are implementation deliverables; their current absence is not a reason to invent a passing result.

Open measurement work: select representative hardware/browser and record resource/latency observations before binding budgets (former Q5). No release date is assigned. Full native-JSON-to-Caddyfile inversion and arbitrary plugin support are not established capabilities; compatibility states preserve the requested editing goal while preventing destructive conversions. [DESIGN](../../ux-designs/ux-caddy-orchestrator-2026-10-16/DESIGN.md) and [EXPERIENCE](../../ux-designs/ux-caddy-orchestrator-2026-10-16/EXPERIENCE.md) retain their detailed interaction/visual assumptions. No customer session, accessibility audit or product runtime certification is implied by approval of this plan.

