# Competition digest — round 1

Accessed: 2026-10-16. Publication: undated living documentation. Method: documentation teardown, no running product verification. Each vendor-specific capability is medium confidence / unverified independently. Two-source checks establish category patterns, not individual implementation quality.

| Claim | Source / publisher | Class | Confidence |
|---|---|---|---|
| NPM centers low-friction home-service exposure: prebuilt Docker, proxy hosts, redirects, streams/404, managed or custom TLS, host access lists/basic auth, advanced configuration, users/permissions/audit. | https://nginxproxymanager.com/guide/ — NPM maintainers | features | medium; single primary source |
| NPM documents shared Docker networks so upstream services need not publish host ports; service names can be upstream hostnames. | https://nginxproxymanager.com/advanced-config/ — NPM maintainers | operations | medium; single primary source |
| X4Applegate/caddyui documents proxy/redirect forms, live certificate status, first-run journey, existing-config import, pre-save Caddy validation, post-apply checks/rollback, roles, snapshots, logs, upstream health and much broader fleet/modules scope. | https://github.com/X4Applegate/caddyui — X4Applegate | features | medium; vendor claim, not exercised |
| zackwag/caddy-ui explicitly keeps the Caddyfile as source of truth and documents editor validation, history/rollback, route health/search, TLS expiry, CA download and logs. | https://github.com/zackwag/caddy-ui — zackwag | features | medium; vendor claim, not exercised |
| rhad00/CaddyManager documents auth/RBAC/TOTP, proxy/multi-domain/TLS, health, logs/alerts, backups, Docker label discovery, GitOps and service templates. | https://github.com/rhad00/CaddyManager — rhad00 | features | medium; vendor claim, not exercised |
| Proxy/TLS convenience and operational visibility/recovery recur across independent vendors. | NPM + X4Applegate + zackwag + rhad00 above | category-pattern | high for documented convergence; no adoption inference |

Implications (researcher inference; not approved scope): a generic CRUD GUI is already offered by several competitors. A credible first release should optimize the complete publish → diagnose → change → recover journey. Avoid claiming unique GUI, automatic HTTPS or modern styling. Treat TLS diagnosis, meaningful live state and rollback as core release candidates. Existing-config users and greenfield appliance users need distinct ownership/onboarding paths.

Leads: verify Caddy restart persistence, Admin API boundary, local CA trust, DNS-01 module prerequisites. Need bounded first-hand issues to test the diagnosis/recovery hypothesis. No product quality ranking, adoption/market-share estimate, pricing-tier assertion or user-preference conclusion is supported.

Stop: coverage for competitor positioning within scope; one targeted second round reserved for current release/change evidence if needed.
