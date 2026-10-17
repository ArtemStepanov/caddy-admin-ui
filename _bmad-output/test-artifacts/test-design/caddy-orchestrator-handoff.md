---
workflowType: test-design-handoff
updated: 2026-10-17
status: corrected-planning
implementation_evidence: not-run
---

# QA → implementation handoff

Use the [current 33-story epics](../../planning-artifacts/epics.md), [architecture contracts](../../planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/configuration-contracts.md) and [122-scenario QA matrix](../test-design-qa.md). Current AC references are authoritative; earlier numeric test/story prefixes are aliases. No current story 3.6/3.7 exists. The complete old→new story allocation is in epics; each test intent has an explicit current AC.

## Delivery and dependencies

1. 1.1 secure seams and single writer mode; 1.2 encrypted source/secret versions; 1.3 functional local migration; 1.4 source/CST/import/preservation corpus; 1.5 target/build compatibility harness.
2. 1.6 adopt without remote mutation; 1.7 revisioned durable draft; 1.8 bound source/runtime review and rebase; 1.9 guarded durable apply/restore/reconciliation; 1.10 first real-service workspace with verification/recovery/restart.
3. Epic 2 completes route/global/server/editor forms and descriptor inventory; Epic 3 completes review/session/reconciliation UX; Epic 4 covers all TLS scenarios; Epic 5 completes typed observations/diagnosis; Epic 6 completes scoped history/backup/maintenance.

Later UI work reuses the core backend; no forward prerequisite or unsafe temporary immediate publication. Planning completion does not set ready-for-dev/done statuses. Story implementation begins only after consistent planning and its own prerequisite delivery.

## Mandatory integration assertions

Root POST /config/ carries exact reviewed If-Match, including restore; the external writer is inserted between final read and write. Snapshot+attempt commit precedes remote mutation; submission_possible persists before send. Every unfinished attempt locks on restart, including accepted response/local commit failure. Reconcile protected content, not ETag versus manager hash. Source mapping preserves nested/root raw data/order or blocks. Historical secret versions hydrate after rotation; missing material blocks. Expected draft revisions and applied probe identities prevent stale data overwrite/success. Browser opacity/policy is unavailable/unknown, not service failure. Functional migration retains local-only data; complete backup acquires remote source/startup/config/data through operator-side procedures.

## UI automation contract

Use native roles/labels first and stable hooks where semantics cannot identify dynamic complex controls:

| Hook | Owner |
| --- | --- |
| nav-routes / nav-changes / nav-recovery | 1.10; extended 2.4/3.1/6.1 |
| route-card-{id} / badge-route-status | 2.4 |
| btn-apply-batch / diff-viewer-container | 1.10; extended 3.2 |
| banner-conflict-409 | 1.10; extended 3.3 |
| banner-uncertain-lock | 1.10; extended 3.5 |
| caddyfile-editor-textarea | 2.2 |
| banner-stale-draft / draft-revision | 3.4 |
| badge-config-acceptance / badge-tls / badge-upstream | 5.4 |
| badge-vantage-manager / badge-vantage-client / evidence-applied-revision | 5.4 |

Hooks do not replace accessibility semantics or expose credentials. Fixture-generated evidence includes exact target/build/source/comparison versions, scenario/AC ID, result/time/vantage. Record manual public ACME or keyboard evidence explicitly. No harness, Playwright environment or new Signals dependency is marked ready merely because this handoff names it.

## Owners and closure

Developer implements fixtures and current ACs; QA reviews the risk/evidence mapping; Architect resolves any demonstrated protocol/source divergence within the approved contract. Product scope choices stay closed. Repeat readiness precedes sprint generation; product release acceptance follows actual scenario execution. The planning correction log records original-finding dispositions and the approved handoff.
