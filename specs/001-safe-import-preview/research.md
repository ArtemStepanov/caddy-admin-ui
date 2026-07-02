# Research: Safe Import Preview

## Decision: Keep the current parser/import architecture and harden it

**Rationale**: The repo already has Caddy fetch, parse, preview, import, route storage, raw route preservation, and dashboard read-only signals. The smallest safe change is to make those paths truthful and strict instead of replacing them.

**Alternatives considered**:
- Rewrite import as a separate subsystem: rejected as unnecessary churn.
- Add a raw config editor: rejected as a safety footgun and out of scope.

## Decision: Classify every discovered HTTP route during parse

**Rationale**: Classification belongs next to route parsing because that is where handler/matcher support is known. Each route should leave parsing with support status, read-only flag, and reason metadata.

**Alternatives considered**:
- Classify only in the UI: rejected because API/storage/sync also need the same safety decision.
- Treat unsupported routes as import errors: rejected because unsupported is expected and must be preserved.

## Decision: Persist read-only reason metadata with routes

**Rationale**: The dashboard needs reasons after import, not only during preview. Persisting a short reason avoids reparsing raw JSON on every list request and keeps the UI simple.

**Alternatives considered**:
- Compute reasons live from raw JSON: rejected as more fragile and repeated work.
- Show only a generic imported badge: rejected because the spec requires explicit reasons.

## Decision: Commit import with one storage transaction

**Rationale**: Current delete-then-insert behavior can leave partial imports. A single replace operation guarantees either all classified routes replace local routes or nothing changes.

**Alternatives considered**:
- Keep delete-then-create and count successes: rejected because it violates the no-partial-import requirement.
- Import route-by-route with best effort: rejected because it hides failures and weakens trust.

## Decision: Preserve read-only routes unchanged on sync

**Rationale**: Read-only means the app does not understand enough to mutate safely. Sync should emit the preserved raw route without matcher/handler normalization.

**Alternatives considered**:
- Merge edited domain/path into raw routes: rejected because the spec now forbids any read-only mutation.
- Rebuild unknown handlers around managed handlers: rejected because ordering and nested subroute semantics are easy to break.

## Decision: Enforce read-only in backend and frontend

**Rationale**: Hiding buttons is not enough. API mutation endpoints must reject edit, delete, and toggle for read-only preserved routes so direct requests cannot damage preserved config.

**Alternatives considered**:
- Frontend-only enforcement: rejected because API clients could still mutate.
- Allow delete as an escape hatch: rejected for v1 because the spec says read-only means no delete.

## Decision: Keep existing import endpoints and expand their response shapes

**Rationale**: Existing `/api/import-preview` and `/api/import` already model the user action. Expanding responses is smaller than adding parallel endpoints.

**Alternatives considered**:
- Add a separate import session resource: rejected as unnecessary without per-route selection or long-lived previews.
- Add WebSocket/progress flow: rejected because import is small and synchronous for v1 scope.

## Decision: Use a best-effort route identity for preview grouping

**Rationale**: The preview must explain new/update/local-only groups. Domain + path + handler type is enough for common homelab routes and avoids adding a new identity system.

**Alternatives considered**:
- Stable Caddy route IDs: rejected because existing Caddy config may not contain app-owned IDs.
- Deep JSON diffing: rejected by scope and unnecessary for the first slice.

## Decision: Manual drift gets clear warning copy, not detection

**Rationale**: The spec explicitly defers automatic drift detection. A visible warning covers v1 without extra state or hash comparisons.

**Alternatives considered**:
- Store full Caddy config hashes and block sync on drift: rejected as later-roadmap scope.
- Auto-merge live Caddy before every sync: rejected as complex and risky.
