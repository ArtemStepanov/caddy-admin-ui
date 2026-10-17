# Caddy native constraints — round 1

Accessed 2026-10-16; undated living official documentation unless stated. Medium confidence / unverified independently except DNS-01 properties corroborated by Let's Encrypt (high). No live deployment test.

- Claim: Caddy automates public/local HTTPS and renewals; local clients must trust its CA. Public issuance depends on challenge reachability and durable writable data. Source: Caddy maintainers, https://caddyserver.com/docs/automatic-https ; class=operations.
- Claim: DNS-01 works without inbound public ports and supports wildcard certificates; provider credentials/configuration are necessary. Sources: Caddy above; Let's Encrypt, https://letsencrypt.org/docs/challenge-types/ ; class=operations, high independent agreement. Provider module presence must be checked rather than assumed (X4Applegate README explicitly distinguishes stock from custom image).
- Claim: Caddy supports scoped config writes and ETag/If-Match protection; multi-request changes are not one transaction; /adapt only converts config without running it. Source: Caddy maintainers, https://caddyserver.com/docs/api ; class=compatibility, medium.
- Claim: API-oriented service uses --resume with persisted autosave.json, unlike file-oriented service. Source: Caddy maintainers, https://caddyserver.com/docs/running ; class=operations, medium.
- Claim: WebSockets, health checks and upstream TLS controls are native reverse_proxy capabilities; disabling TLS verification removes protection. Source: Caddy maintainers, https://caddyserver.com/docs/caddyfile/directives/reverse_proxy ; class=features, medium.

Implications: define durable source of truth before adding editing surfaces; verify restart and disconnected-manager traffic continuity. Distinguish configuration acceptance, certificate lifecycle and upstream reachability. DNS-01 is valuable to private labs, but custom build/provider secrets extend operational scope. Preserve unknown configuration and refuse drift rather than claiming arbitrary-config editing is safe. An adaptation success is not proof that a complete deployment works.

Stop: coverage; exact Caddy versions/provider modules supported by this product remain architecture decisions, not externally inferred facts.
