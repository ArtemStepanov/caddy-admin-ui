# Architecture Diagrams — Caddy Orchestrator

The [spine](../../planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md) owns stable AD rules; its [companion](../../planning-artifacts/architecture/architecture-caddy-orchestrator-2026-10-17/configuration-contracts.md) owns the shared record/API model. Diagrams summarize those contracts and do not supersede them.

```mermaid
flowchart TD
  Forms[Forms / structured directives] <--> Draft[Revisioned protected Caddyfile bundle]
  Editor[Caddyfile editor] <--> Draft
  Draft --> Compile[Transient bundle expansion and target adapt]
  Baseline[Live root + preservation manifest] --> Compose[Proven mapped runtime delta]
  Compile --> Compose
  Compose --> Review[Bound source/runtime review]
  Review --> Intent[SQLite pre-write recovery + durable attempt]
  Intent --> Guard[POST /config/ with exact If-Match]
  Guard --> Readback[Observe content and finalize applied revision]
  Guard --> Uncertain[Unfinished attempt locks mutations]
  Uncertain --> Reconcile[Candidate / baseline / conflict / unavailable]
  Readback --> Verify[Revision-bound typed observations]
```

```mermaid
flowchart LR
  Browser[Operator browser] -->|HTTP REST / SSE :3000| Manager[Manager-only Compose]
  Browser -->|Published HTTPS URL| Caddy[External Caddy]
  Manager -->|GET config / POST adapt / conditional root config| Caddy
  Manager -->|GET-only metadata, optional| Docker[Restricted Docker proxy]
  Manager --> DB[Persistent SQLite + encryption key]
  Caddy --> Assets[Persistent startup config / source / data]
  Caddy --> Upstream[Upstream services]
  Operator[Operator-side backup commands] --> DB
  Operator --> Assets
```

No transaction spans SQLite and Caddy. Restore uses the same reviewed guarded attempt protocol, generating a new applied revision. Manager-network, direct-upstream and client observations are distinct targets/vantages. No automatic source inversion or remote filesystem write is implied.
