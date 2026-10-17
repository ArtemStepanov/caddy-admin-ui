---
workflowStatus: corrected-planning
updated: 2026-10-17
implementation_evidence: not-run
reference_target: Caddy 2.9.1
---

# System test design — QA execution plan

This plan supersedes the earlier inconsistent test/story numbering under [approved P01–P13](../planning-artifacts/sprint-change-proposal-2026-10-17.md). It is **planned acceptance**, not Ready-for-Execution tooling or passing product evidence. [Epics](../planning-artifacts/epics.md) defines 33 current stories and 122 numbered acceptance scenarios. Each scenario below maps to exactly one `CC` test intent; implementations may split it into several unit/integration/E2E cases while retaining the AC mapping. Test count does not measure runtime correctness.

## Inputs and common contract

[PRD](../planning-artifacts/prds/prd-caddy-orchestrator-2026-10-16/prd.md), [architecture](../planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md), [configuration contract](../planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/configuration-contracts.md), [SPEC](../specs/spec-caddy-orchestrator/SPEC.md) and corrected UX supply assertions. Root writes are POST /config/ + exact root If-Match, remote 412→application 409; target adapt returns result/warnings and is not semantic provisioning validation. Config export is not module inventory. Review/attempt/target/generation/revision identity is shared by all tests.

## Findings → acceptance matrix

Every original readiness item below has a concrete contract plus current acceptance references. This is planning closure intent; runtime closure requires executed evidence.

| Finding | Contract/coverage | Recorded location | Current ACs |
| --- | --- | --- | --- |
| B01 | Atomic root precondition and reviewed rebase/restore | AD-3/8; companion P02/X1/X2 | 1.8.AC2, 1.8.AC3, 1.9.AC2, 6.3.AC2 |
| B02 | Source association and nested preservation | AD-1/2; companion P03/X1 | 1.4.AC2, 1.4.AC3, 1.6.AC1, 1.6.AC2, 1.6.AC3 |
| B03 | Durable attempts and crash/finalization reconciliation | AD-3/10; companion P04/X2 | 1.9.AC1, 1.9.AC3, 1.9.AC4, 1.9.AC6, 3.5.AC3 |
| B04 | Protected source and historical secret hydration | AD-5/11; companion P05/X2 | 1.2.AC1, 1.2.AC2, 1.2.AC3, 1.2.AC4, 4.2.AC3 |
| B05 | Forward implementable complete service slice | epics delivery order / aliases | 1.10.AC1, 1.10.AC2, 1.10.AC3 |
| B06 | All-settings forms and retained baseline features | AD-1/9; companion P07/X3 | 2.1.AC1, 2.1.AC2, 2.1.AC4, 2.5.AC1, 2.5.AC2, 2.5.AC3 |
| B07 | Draft/session/probe revisions | AD-4/7; companion P06/X2 | 1.7.AC2, 1.7.AC4, 3.4.AC1, 5.1.AC3, 5.4.AC2 |
| B08 | Bounded typed and browser evidence | AD-4; companion P08 | 5.1.AC2, 5.2.AC3, 5.3.AC1, 5.3.AC2, 5.3.AC3, 5.5.AC1 |
| B09 | Pinned Caddy request/response and build evidence | AD-9; companion P08/X3 | 1.5.AC2, 1.5.AC3, 1.5.AC4, 4.2.AC2, 4.3.AC2 |
| B10 | Functional transactional migration | AD-11; companion P09/X4 | 1.3.AC1, 1.3.AC2, 1.3.AC3, 1.3.AC4, 6.4.AC3 |
| C01 | Scope and artifact authority reconciled | PRD 1/4/8; SPEC authority | 2.5.AC1, 4.2.AC1, 4.3.AC1, 4.4.AC1 |
| C02 | Current test/story mapping and honest gates | QA matrix / architecture test design | 1.1.AC3, 1.5.AC4, 1.9.AC3, 1.10.AC4 |
| C03 | Complete backup and remote operator acquisition | AD-11; companion P10/X4 | 6.4.AC1, 6.4.AC2 |
| C04 | Identity and input edge contracts | AD-1; companion P11/X3 | 1.4.AC1, 1.4.AC4, 2.1.AC3, 4.4.AC3 |
| C05 | Previous review dispositions and consistent targets | review-dispositions.md; SPEC diagrams | 1.5.AC4, 1.10.AC4, 2.4.AC3 |

## Scenario matrix

All entries are **planned**. `CC` IDs are current stable test-intent IDs; their independent Current AC column is authority. IDs in the historical QA document are aliases only and must not be parsed as current story numbers. Assigned owner is Developer for implementation/fixtures and QA for scenario evidence/review. Evidence files are planned `evidence/<test-id>.json` with actual version/environment, checks, result and timestamp; manual evidence may use linked Markdown with the same identity. No empty output file counts as execution.

| Test ID | Current AC | Priority | Planned level | Original finding links |
| --- | --- | --- | --- | --- |
| CC-1.1-01 | 1.1.AC1 | P0 | Integration | retained requirement |
| CC-1.1-02 | 1.1.AC2 | P0 | Integration | retained requirement |
| CC-1.1-03 | 1.1.AC3 | P0 | Integration | C02 |
| CC-1.1-04 | 1.1.AC4 | P0 | Integration | retained requirement |
| CC-1.2-01 | 1.2.AC1 | P0 | Integration | B04 |
| CC-1.2-02 | 1.2.AC2 | P0 | Integration | B04 |
| CC-1.2-03 | 1.2.AC3 | P0 | Integration | B04 |
| CC-1.2-04 | 1.2.AC4 | P0 | Integration | B04 |
| CC-1.3-01 | 1.3.AC1 | P0 | Integration | B10 |
| CC-1.3-02 | 1.3.AC2 | P0 | Integration | B10 |
| CC-1.3-03 | 1.3.AC3 | P0 | Integration | B10 |
| CC-1.3-04 | 1.3.AC4 | P0 | Integration | B10 |
| CC-1.4-01 | 1.4.AC1 | P0 | Unit + target parity | C04 |
| CC-1.4-02 | 1.4.AC2 | P0 | Unit + target parity | B02 |
| CC-1.4-03 | 1.4.AC3 | P0 | Unit + target parity | B02 |
| CC-1.4-04 | 1.4.AC4 | P0 | Unit + target parity | C04 |
| CC-1.5-01 | 1.5.AC1 | P1 | Integration | retained requirement |
| CC-1.5-02 | 1.5.AC2 | P1 | Integration | B09 |
| CC-1.5-03 | 1.5.AC3 | P1 | Integration | B09 |
| CC-1.5-04 | 1.5.AC4 | P1 | Integration | B09, C02, C05 |
| CC-1.6-01 | 1.6.AC1 | P0 | Integration | B02 |
| CC-1.6-02 | 1.6.AC2 | P0 | Integration | B02 |
| CC-1.6-03 | 1.6.AC3 | P0 | Integration | B02 |
| CC-1.6-04 | 1.6.AC4 | P0 | Integration | retained requirement |
| CC-1.7-01 | 1.7.AC1 | P0 | Integration | retained requirement |
| CC-1.7-02 | 1.7.AC2 | P0 | Integration | B07 |
| CC-1.7-03 | 1.7.AC3 | P0 | Integration | retained requirement |
| CC-1.7-04 | 1.7.AC4 | P0 | Integration | B07 |
| CC-1.8-01 | 1.8.AC1 | P0 | Integration | retained requirement |
| CC-1.8-02 | 1.8.AC2 | P0 | Integration | B01 |
| CC-1.8-03 | 1.8.AC3 | P0 | Integration | B01 |
| CC-1.8-04 | 1.8.AC4 | P0 | Integration | retained requirement |
| CC-1.9-01 | 1.9.AC1 | P0 | Integration | B03 |
| CC-1.9-02 | 1.9.AC2 | P0 | Integration | B01 |
| CC-1.9-03 | 1.9.AC3 | P0 | Integration | B03, C02 |
| CC-1.9-04 | 1.9.AC4 | P0 | Integration | B03 |
| CC-1.9-05 | 1.9.AC5 | P0 | Integration | retained requirement |
| CC-1.9-06 | 1.9.AC6 | P0 | Integration | B03 |
| CC-1.10-01 | 1.10.AC1 | P0 | E2E / integration | B05 |
| CC-1.10-02 | 1.10.AC2 | P0 | E2E / integration | B05 |
| CC-1.10-03 | 1.10.AC3 | P0 | E2E / integration | B05 |
| CC-1.10-04 | 1.10.AC4 | P0 | E2E / integration | C02, C05 |
| CC-2.1-01 | 2.1.AC1 | P1 | Unit + target parity | B06 |
| CC-2.1-02 | 2.1.AC2 | P1 | Unit + target parity | B06 |
| CC-2.1-03 | 2.1.AC3 | P1 | Unit + target parity | C04 |
| CC-2.1-04 | 2.1.AC4 | P1 | Unit + target parity | B06 |
| CC-2.2-01 | 2.2.AC1 | P1 | Component + target parity | retained requirement |
| CC-2.2-02 | 2.2.AC2 | P1 | Component + target parity | retained requirement |
| CC-2.2-03 | 2.2.AC3 | P1 | Component + target parity | retained requirement |
| CC-2.2-04 | 2.2.AC4 | P1 | Component + target parity | retained requirement |
| CC-2.3-01 | 2.3.AC1 | P1 | Integration | retained requirement |
| CC-2.3-02 | 2.3.AC2 | P1 | Integration | retained requirement |
| CC-2.3-03 | 2.3.AC3 | P1 | Integration | retained requirement |
| CC-2.3-04 | 2.3.AC4 | P1 | Integration | retained requirement |
| CC-2.4-01 | 2.4.AC1 | P1 | E2E / integration | retained requirement |
| CC-2.4-02 | 2.4.AC2 | P1 | E2E / integration | retained requirement |
| CC-2.4-03 | 2.4.AC3 | P1 | E2E / integration | C05 |
| CC-2.4-04 | 2.4.AC4 | P1 | E2E / integration | retained requirement |
| CC-2.5-01 | 2.5.AC1 | P1 | Integration | B06, C01 |
| CC-2.5-02 | 2.5.AC2 | P1 | Integration | B06 |
| CC-2.5-03 | 2.5.AC3 | P1 | Integration | B06 |
| CC-2.5-04 | 2.5.AC4 | P1 | Integration | retained requirement |
| CC-3.1-01 | 3.1.AC1 | P1 | E2E / integration | retained requirement |
| CC-3.1-02 | 3.1.AC2 | P1 | E2E / integration | retained requirement |
| CC-3.1-03 | 3.1.AC3 | P1 | E2E / integration | retained requirement |
| CC-3.1-04 | 3.1.AC4 | P1 | E2E / integration | retained requirement |
| CC-3.2-01 | 3.2.AC1 | P1 | E2E / integration | retained requirement |
| CC-3.2-02 | 3.2.AC2 | P1 | E2E / integration | retained requirement |
| CC-3.2-03 | 3.2.AC3 | P1 | E2E / integration | retained requirement |
| CC-3.2-04 | 3.2.AC4 | P1 | E2E / integration | retained requirement |
| CC-3.3-01 | 3.3.AC1 | P1 | E2E / integration | retained requirement |
| CC-3.3-02 | 3.3.AC2 | P1 | E2E / integration | retained requirement |
| CC-3.3-03 | 3.3.AC3 | P1 | E2E / integration | retained requirement |
| CC-3.4-01 | 3.4.AC1 | P1 | E2E / integration | B07 |
| CC-3.4-02 | 3.4.AC2 | P1 | E2E / integration | retained requirement |
| CC-3.4-03 | 3.4.AC3 | P1 | E2E / integration | retained requirement |
| CC-3.5-01 | 3.5.AC1 | P1 | E2E / integration | retained requirement |
| CC-3.5-02 | 3.5.AC2 | P1 | E2E / integration | retained requirement |
| CC-3.5-03 | 3.5.AC3 | P1 | E2E / integration | B03 |
| CC-4.1-01 | 4.1.AC1 | P1 | E2E / integration | retained requirement |
| CC-4.1-02 | 4.1.AC2 | P1 | E2E / integration | retained requirement |
| CC-4.1-03 | 4.1.AC3 | P1 | E2E / integration | retained requirement |
| CC-4.1-04 | 4.1.AC4 | P1 | E2E / integration | retained requirement |
| CC-4.2-01 | 4.2.AC1 | P0 | Integration | C01 |
| CC-4.2-02 | 4.2.AC2 | P0 | Integration | B09 |
| CC-4.2-03 | 4.2.AC3 | P0 | Integration | B04 |
| CC-4.2-04 | 4.2.AC4 | P0 | Integration | retained requirement |
| CC-4.3-01 | 4.3.AC1 | P1 | E2E / integration | C01 |
| CC-4.3-02 | 4.3.AC2 | P1 | E2E / integration | B09 |
| CC-4.3-03 | 4.3.AC3 | P1 | E2E / integration | retained requirement |
| CC-4.4-01 | 4.4.AC1 | P1 | Integration | C01 |
| CC-4.4-02 | 4.4.AC2 | P1 | Integration | retained requirement |
| CC-4.4-03 | 4.4.AC3 | P1 | Integration | C04 |
| CC-5.1-01 | 5.1.AC1 | P1 | Integration | retained requirement |
| CC-5.1-02 | 5.1.AC2 | P1 | Integration | B08 |
| CC-5.1-03 | 5.1.AC3 | P1 | Integration | B07 |
| CC-5.1-04 | 5.1.AC4 | P1 | Integration | retained requirement |
| CC-5.2-01 | 5.2.AC1 | P1 | Integration | retained requirement |
| CC-5.2-02 | 5.2.AC2 | P1 | Integration | retained requirement |
| CC-5.2-03 | 5.2.AC3 | P1 | Integration | B08 |
| CC-5.3-01 | 5.3.AC1 | P1 | E2E / integration | B08 |
| CC-5.3-02 | 5.3.AC2 | P1 | E2E / integration | B08 |
| CC-5.3-03 | 5.3.AC3 | P1 | E2E / integration | B08 |
| CC-5.4-01 | 5.4.AC1 | P1 | E2E / integration | retained requirement |
| CC-5.4-02 | 5.4.AC2 | P1 | E2E / integration | B07 |
| CC-5.4-03 | 5.4.AC3 | P1 | E2E / integration | retained requirement |
| CC-5.5-01 | 5.5.AC1 | P1 | E2E / integration | B08 |
| CC-5.5-02 | 5.5.AC2 | P1 | E2E / integration | retained requirement |
| CC-5.5-03 | 5.5.AC3 | P1 | E2E / integration | retained requirement |
| CC-6.1-01 | 6.1.AC1 | P1 | E2E / integration | retained requirement |
| CC-6.1-02 | 6.1.AC2 | P1 | E2E / integration | retained requirement |
| CC-6.1-03 | 6.1.AC3 | P1 | E2E / integration | retained requirement |
| CC-6.2-01 | 6.2.AC1 | P1 | E2E / integration | retained requirement |
| CC-6.2-02 | 6.2.AC2 | P1 | E2E / integration | retained requirement |
| CC-6.2-03 | 6.2.AC3 | P1 | E2E / integration | retained requirement |
| CC-6.3-01 | 6.3.AC1 | P0 | E2E / integration | retained requirement |
| CC-6.3-02 | 6.3.AC2 | P0 | E2E / integration | B01 |
| CC-6.3-03 | 6.3.AC3 | P0 | E2E / integration | retained requirement |
| CC-6.4-01 | 6.4.AC1 | P0 | E2E / integration | C03 |
| CC-6.4-02 | 6.4.AC2 | P0 | E2E / integration | C03 |
| CC-6.4-03 | 6.4.AC3 | P0 | E2E / integration | B10 |
| CC-6.4-04 | 6.4.AC4 | P0 | E2E / integration | retained requirement |

**Matrix totals:** 122 scenario intents: 49 P0, 73 P1; 33 stories. These replace historical 33/34/35-test totals. No execution duration or implementation estimate is inferred from these counts.

## Risk coverage and mandatory details

- **R-01 source fidelity/preservation:** unchanged byte corpus, snippets/import arguments/env expansion, source-map errors, duplicate/missing IDs, disabled nodes, multiple routes/site, unknown root and nested config, mixed arrays and server regrouping. Test B/G0 association and composed candidate, not merely a top-level app merge.
- **R-02 concurrent overwrite:** inject an external write **after final read and before conditional root submission** on real pinned Caddy, not just before a mock comparison. Assert remote rejection and unchanged competing config. Repeat for restore. Rebase covers nonoverlap, ancestor overlap, deletion/edit and arrays with renewed source association/review.
- **R-03 ambiguous remote state:** crash at every boundary (before submission_possible, before/after send, after response, before/after readback/local finalization); DB failures and malformed/proxy errors; candidate/baseline/other/unavailable equality, no-op candidates and unsupported normalization. Startup must lock from durable intent even if later error recording failed.
- **R-04 secret lifecycle:** plaintext detection in DB active payloads, source/invalid buffers, retained migration archives, API/editor/diffs/logs/errors/ordinary exports; historical rotation/deletion, missing key/version, reference token forgery, masked-value round trip. Authenticated encrypted backups are distinct from sanitized exports.
- **R-05 persistence/maintenance:** independent manager/Caddy restart/recreation, resume/config/data prerequisites, functional prior-release migration including disabled/local intent, rollback/interruption and clean remote-asset backup restoration. API reads cannot prove remote mounts/flags.
- **R-06 false evidence:** protected-service expected statuses, TLS trust/SAN versus near-expiry warnings, direct upstream versus published URL, redirect first hop, known file resource, browser CORS/CSP/opaque/mixed content, late probes after restore and explicit stale time.
- **R-07 auth/CSRF:** listener-based auth, partial credentials even with override, all POST/PUT/PATCH/DELETE browser-origin/header failures, no unauthenticated fault endpoints.
- **R-08 adaptation flood/stale response:** 500ms debounce, coalesced explicit triggers, result/warnings wrapper, stale revision response ignored. No universal semantic dry run.
- **R-09 bounded recoverability:** history cap, oldest safely unpinned pruning, all-pinned rejection, secret/attempt references retained; never prune active recovery material for a count assertion.
- **R-10 stale shared state:** same-revision concurrent saves, lost SSE IDs/reconnect, target switch, unsaved input retention and no offline auto-replay.
- **R-11 unsupported editing/build:** enumerate reference build settings against descriptors/forms/CST/runtime/secret/validation fixtures, including globals/server settings and retained handler options; unknown module/native-JSON mapping blocks rather than falsely passes.

## Fixture and harness delivery order

| Owner story | Required delivery | Dependent acceptance |
| --- | --- | --- |
| 1.1 | Auth/origin middleware fixtures and writer-mode seam | Every later inbound mutation |
| 1.2 | Protected payload/secret-version factories, safe output capture | Adoption/draft/history/restore |
| 1.3 | Actual prior-release DB fixture with local-only data and rollback harness | Upgrade/release maintenance |
| 1.4 | Golden source/import/preservation corpus and target parity runner | Form/editor/mapping acceptance |
| 1.5 | Pinned Caddy 2.9.1 + reference DNS-enabled build lock and correct API wrapper fixtures | Target compatibility/TLS |
| 1.9 | Deterministic competing writer, process-kill, network reset and local-commit failure seams | Apply/restore/reconciliation |
| 1.10 | Real manager/external-Caddy/reference-service E2E harness and locked Signals version | Complete operator flow |
| 2.5 | Exhaustive declared-build setting descriptor/inventory validator | All-settings coverage |
| 5.3 | Browser-policy endpoint fixtures and scoped CSP checks | Client-vantage evidence |
| 6.4 | Operator-side external asset capture/clean restore environment | Complete backup/upgrade release gate |

Corpus size is driven by grammar/features/edge coverage; the prior arbitrary 50+ sample claim is replaced by a coverage manifest and exact-case golden results. Fixtures/mocks enable determinism but do not replace real pinned-target protocol/adapter cases. Public ACME/domain evidence records DNS/network prerequisites; controlled local ACME simulation is labelled simulation and never substituted for a public-domain issuance claim. Default no outbound live-service tests in ordinary unit suites.

## Automation and execution

Retain existing Go testing/httptest and frontend Vitest/test-library tooling. Playwright, accessibility tooling and any new harness are planned dependencies with locked versions selected when delivered, not claimed installed. Use semantic locators; add stable data-testid hooks for complex editor/diff/state controls where needed. Core hook names are specified in the handoff.

Fast deterministic unit/component/contract checks run per PR; pinned-container integration and E2E checks run in a supported container environment. Long renewal/restart/backup scenarios may be scheduled separately with required release execution. Exact CI budgets/cadence are decided from measured suite duration; no undocumented workflow is already configured. Fault injection is test-only port/transport substitution, never a production mutation API.

## Entry and exit evidence

Entry for a scenario: its prerequisites/harness are implemented and buildable, contract and AC refs resolve, target/source/build locks are recorded, and required secrets/remote assets are acquired safely. Planning PASS is not runtime entry for all tests.

Release exit: all required functional acceptance and P0 safety scenarios pass with no unresolved high-severity data/security defect. Any incomplete/unavailable required scenario prevents a full-scope release claim; observed evidence can still describe narrower facts. Coverage percentages and performance thresholds are measurements unless explicitly agreed; neither 80% line coverage nor P95<200ms nor <100MB/<1% CPU is an invented acceptance gate. Accessibility requires keyboard/focus/editor escape/zoom/manual inspection in addition to automated scans. Evidence records current version/environment and known limits.

## Historical alias corrections

Old test IDs such as `3.6-INT-001/002` remain traceable to uncertainty now in 1.9/3.5, not nonexistent current stories. Old public HTTPS `4.3-E2E-001` maps to 1.10/4.1; plugin `4.2-INT-001` to 4.2; restore `6.1-INT-001` to 1.9/6.3; backup `6.2-INT-001` to 6.4; restart `6.3-E2E-001` to 1.10/6.4; accessibility `6.4-E2E-001` to each UI story and final release checks. The complete old-story allocation is in epics. Test intent names and Current AC fields govern meaning, not historical numeric prefixes.
