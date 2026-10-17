---
title: Caddy Orchestrator PRD Supporting Notes
status: reconciled
updated: 2026-10-17
---

# Supporting notes

[PRD](prd.md) records confirmed scope after [approved P01–P13](../../sprint-change-proposal-2026-10-17.md). Full existing configuration, all-settings forms/Caddyfile, shared drafts, four TLS scenarios and Docker suggestions/manual input supersede earlier candidate exclusions. Historical discovery is complete. No competitor feature adds unapproved scope.

[CLAUDE.md](../../../../CLAUDE.md) and the current source remain the brownfield implementation baseline: scoped route writes, ETags and route-array snapshots. New root ownership requires explicit adoption and one enforced writer mode; do not relabel legacy snapshots as full-config history. Historical `specs/001` and `specs/002` are not the active future contract.

Use [architecture and its companion](../../architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md) for source association, nested preservation, durable guarded apply/restore, secret versions, migration and backup. Keep Go/Gin/SQLite/Preact/Tailwind. Manager-only Compose targets existing Caddy; bundled Caddy is a test fixture. The manager has no implied remote filesystem or Docker write access.

A separate Caddy fixture and reference service establish reproducible acceptance; Artem's private services are optional additional evidence. Compare root config before/after managed changes, conflict and restore. Include other servers/listeners, PKI/logging/storage, nested unsupported routes, disabled/local-only legacy data, secret rotation and interrupted operations. Configuration acceptance, certificate readiness, manager reachability and client trust are different observations.

Record exact version/environment and evidence type. Caddy startup persistence, clean backup restoration, upgrade, issuance/renewal, keyboard/zoom/focus and typed verification require actual execution before release claims. Planned corpus/build locks and test hooks are not implemented evidence. Numeric performance budgets remain unchosen. Scope approval does not claim deployment or product verification.
