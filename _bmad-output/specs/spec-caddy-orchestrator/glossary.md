# Glossary — Caddy Orchestrator

- **Target:** one existing Caddy endpoint identity plus adoption generation and exclusive manager write mode (legacy routes or adopted root).
- **Route:** a stable source node mapping host/path to a handler; multiple routes may share a site. Site and route identities are distinct UUIDs.
- **Upstream:** a validated destination with optional independent TLS/SNI/trust settings, distinct from the published client URL.
- **Source Bundle:** exact encrypted Caddyfile/dependencies, immutable enabled/node-identity metadata, explicit environment inputs, ordered import resolution, source map and build provenance. Original source is distinct from transient expanded compilation text.
- **CST:** source tree preserving syntax, comments, whitespace and unknown constructs; parsing is byte-preserving. Explicit save assigns missing IDs; duplicates require repair.
- **Preservation manifest:** nonoverlapping generated versus raw runtime path ownership; protects root/nested content and order. It replaces the insufficient top-level passthrough-apps model.
- **Draft revision:** integer used for stale-save protection, separate from opaque target ETags and applied revision IDs. A saved draft has not published.
- **Review:** single-use binding of target/generation/draft revision/baseline ETag/candidate/source/build/dependencies/preservation scope to a sanitized source/runtime diff.
- **Apply attempt:** durable pre-submission record with protected baseline/candidate, recovery identity and lifecycle, shared by ordinary apply and restore. It pins the immutable connection descriptor and locks draft/target mutation until resolved.
- **Applied revision:** identity of observed accepted config, distinct from draft revision and verification outcome. Restore creates a new applied revision.
- **Recovery Point:** protected configuration/source/secret-version record with explicit target/generation/scope/time/reason. Full config excludes external certificate/private-key/host assets. Legacy route-array Snapshots retain their narrower server scope.
- **Secret version:** immutable encrypted credential value retained while referenced by drafts, attempts or history. Display masks are never credentials.
- **Service verified:** required typed expectations observed for the named applied revision, time and vantage. Unsupported/opaque checks stay unknown; no global client-reachability claim follows.
- **Complete backup:** consistent DB/key/current config/source/import/env/build/startup/external-data package with explicit operator-side acquisition and clean-restore proof, distinct from sanitized exports/history.
