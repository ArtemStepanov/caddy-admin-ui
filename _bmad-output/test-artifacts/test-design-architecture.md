---
workflowStatus: corrected-planning
updated: 2026-10-17
implementation_evidence: not-run
---

# Architecture test design — shared safety contracts

Scope: [33 stories](../planning-artifacts/epics.md) implementing [AD-1–AD-11](../planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md) and [shared configuration contracts](../planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/configuration-contracts.md). The [QA matrix](test-design-qa.md) is the single current scenario register; this file specifies testability seams and evidence boundaries rather than duplicating inconsistent counts.

## Contract seams

| Boundary | Shared model / proof | Story delivering seam |
| --- | --- | --- |
| Source ↔ forms/editor/compiler | Exact encrypted SourceBundleV1, descriptor forms, stable node IDs, deterministic captured import expansion, target parity and preserved path composition | 1.2, 1.4; complete inventory 2.5 |
| Target ↔ review/apply/restore | Exact root ETag, POST /config/ If-Match, remote 412→app 409, review-bound identities and no unconditional fallback | 1.5, 1.8, 1.9 |
| Storage ↔ apply | Transactional snapshot+attempt, durable submission_possible, protected baseline/candidate, finalization/readback and startup lock | 1.2, 1.9 |
| Shared sessions ↔ draft | Expected revision CAS, target generation, event identity and unsaved-buffer resync | 1.7; complete UI 3.4 |
| Publication ↔ verification | Applied revision/route/expectation/job identity, explicit method/vantage/time/stale state | 1.10; full engine/UI 5.1–5.4 |
| Secrets ↔ all source/runtime history | Encrypted exact envelopes + immutable typed versions, authorized keep/replace/remove, no mask literals, safe errors/exports | 1.2; provider use 4.2 |
| Legacy ↔ adopted root mode | Transactional functional record mapping, retained local intent/snapshot scope, no legacy writer after root adoption | 1.1, 1.3, 1.6 |
| Manager ↔ remote host recovery | Operator-side acquisition, consistent DB/key/current source/startup/build/env/config/data assets, incomplete-package state | 6.4 |

## Planned fault controls

Inject through ports, httptest transports, subprocess lifecycle and disposable DB handles: an external writer just before submission; reset/drop after remote acceptance; process kill at each durable transition; readback error; SQLite pre-write/finalization failure; unavailable secret version; stale review/draft/event/result. No production API switch can bypass authentication, create fake success or mutate Caddy in a test-only way.

Real Caddy 2.9.1 integration must test root null/empty behavior and exact ETag handling. The adapter decodes {result,warnings}; no validate query establishes provisioning. Provisioning errors and unknown outcomes are distinct from syntax errors. GET config cannot return fictitious module inventory. A successful root HTTP response followed by failed local finalization remains reconciliation-required until content readback and local commit establish current presence. A changed ETag is never compared to a manager candidate hash.

## Risk and scope assessment

Use QA R-01–R-11 and B01–B10/C01–C05 mapping. Concurrency, source loss, secret lifecycle, crash recovery and migration are critical data-safety risks. Shared sessions remain relevant despite single-operator scope. Read-only Docker suggestions are optional; a raw socket ro mount does not enforce API permissions. No remote-host write privilege is added for backup.

Supported-build coverage is a descriptor inventory validation problem, not a claim that arbitrary modules already have forms. JSON-only/source-incompatible content must be preserved and visibly blocked. Golden corpus and module locks are planned deliverables with owners, not completed evidence; their absence before implementation is normal, but no dependent story can pass without them.

## Operational evidence limits

Current source establishes existing dependencies/security/storage behavior, not the new domain/adapters already implemented. The manager API cannot discover remote volume mounts/resume flags by reading config; operator declarations plus fixture restart evidence satisfy the deployment contract honestly. Exact hardware, route counts, RAM/CPU/latency and CI duration budgets remain unspecified; measure before binding. Test counts and prior 74–121-hour estimates are historical, not this correction's delivery commitment.

Complete release proof covers independent restart, real publication/typed observations, safe concurrent writes, interruption/local-failure reconciliation, credential rotation, stale clients/jobs, prior-release functional upgrade, protected complete backup and clean restore. Automated scans complement manual keyboard/zoom/focus and editor escape checks. A planning gate certifies recorded decisions, never these product outcomes.
