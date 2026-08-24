# Roadmap

The project's logical finish line is a focused, dependable single-instance Caddy route manager, not a replacement for every Caddy module or a general hosting control plane.

## Stabilization (current)

- [x] Explicit server discovery and ownership confirmation
- [x] Scoped route-array writes with ETag drift protection
- [x] Exact preservation of external routes
- [x] Pre-write snapshots, export, and guarded restore
- [x] Safe-by-default container and authentication posture
- [x] Real-Caddy integration coverage and blocking security checks

## 0.4: Operability

- [ ] Route reordering with a clear precedence preview
- [ ] Certificate inventory and expiry status
- [ ] Upstream health indicators
- [ ] Filtered access-log viewer with redaction controls
- [ ] Snapshot retention policy and database backup command

## 0.5: Route capabilities

- [ ] Route-level Basic Auth with passwords hashed before storage
- [ ] Health checks and transport settings for reverse proxies
- [ ] Safer matcher editor beyond host and path
- [ ] Caddy JSON export without applying changes

## 1.0 criteria

- Documented database and API compatibility policy
- Upgrade and rollback guide tested across two consecutive releases
- No known lossy path for preserved external route JSON
- End-to-end coverage for setup, drift, CRUD, restore, and authentication
- Maintainer-reviewed threat model and stable deployment examples

## Explicit non-goals for 1.0

- Managing every Caddy app or plugin
- Editing arbitrary external route JSON in the UI
- Multi-tenant RBAC or active-active UI replicas
- Exposing the Caddy Admin API to public networks
