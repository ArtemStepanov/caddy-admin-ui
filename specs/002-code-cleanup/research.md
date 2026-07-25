# Research: Code Cleanup

## Decision: Reuse `config.ValidateRoutesForBuild` as the one import-safety gate

**Rationale**: `internal/api/handlers.go` currently runs both its local `validateImportedRoutes` helper and `config.ValidateRoutesForBuild` after `ParseCaddyConfig`. The config validator already validates editable handler configuration and preserved raw routes, and it is also the validation used before sync. Remove the weaker API-local pass; parsing remains responsible for classifying editable versus preserved routes, and the shared validator makes the accept/preserve/reject decision before either preview is returned or storage is replaced.

**Alternatives considered**:
- Add an import-decision type/service: rejected; it adds a layer around the existing parser and validator without a new workflow.
- Keep both validators: rejected; their overlapping checks can diverge and the API helper does not validate handler-specific configuration.
- Validate only in the UI: rejected; Caddy-derived data must be checked on the server before persistence.

## Decision: Preserve legacy and unsupported routes through the existing read-only model

**Rationale**: `storage.Route.IsReadOnly`, `RawCaddyRoute`, `SupportStatus`, and `ReadOnlyReason` already preserve legacy raw routes, reject mutations, and write preserved Caddy JSON back unchanged. The cleanup will exercise those cases through the shared import gate rather than change their representation or migration.

**Alternatives considered**:
- Add a new database decision table: rejected; import decisions are derived from the fetched route and need no persistence beyond the existing route metadata.
- Reclassify legacy raw routes: rejected; it risks turning preserved routes into editable ones.

## Decision: Keep the three existing route kinds and route-level response headers

**Rationale**: The form and builder support `reverse_proxy`, `file_server`, and `redir`. Each has its own `config` shape; response headers are a separate route-level setting that the builder applies to every supported kind, so they must survive a kind change. `strip_path_prefix` is only meaningful in this UI for reverse-proxy forwarding and must be omitted after switching to another kind. The existing kind-change reset already replaces handler-specific config; focused tests will lock that behavior down.

**Alternatives considered**:
- Add a generic route-config abstraction or state machine: rejected; three fixed kinds only need the existing switch and one conditional submission field.
- Clear response headers on every kind change: rejected; headers are supported by all current route builders and clearing them loses valid configuration.
- Add new route kinds or a raw editor: rejected as explicitly out of scope.

## Decision: Use existing Go and Vitest test tooling

**Rationale**: Parser/API tests already use table-like Caddy fixtures, and the frontend already uses Vitest with Testing Library. Add focused cases alongside those tests; no dependency, integration service, or database migration is needed.

**Alternatives considered**:
- Add browser automation: rejected; mocked API requests and existing Go HTTP tests cover this bounded behavior.
- Add a new test framework: rejected; current tooling is sufficient.
