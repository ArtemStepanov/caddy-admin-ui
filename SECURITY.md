# Security policy

## Supported versions

The project is pre-1.0. Security fixes are made on the latest release and `main`; older releases are not maintained.

## Reporting a vulnerability

Please use GitHub's **Report a vulnerability** flow in the repository Security tab. Do not open a public issue with exploit details.

Include the affected version or commit, deployment topology, reproduction steps, impact, and any suggested mitigation. You should receive an acknowledgement within seven days. Disclosure timing will be coordinated after a fix is available.

## Deployment assumptions

- Caddy's Admin API is highly privileged and should remain on loopback or a trusted private network.
- Caddy Admin UI should be loopback-bound or protected by authentication and TLS at a trusted reverse proxy.
- `ALLOW_INSECURE_NO_AUTH=true` is an explicit risk acceptance, not a recommended production setting.
- Back up the SQLite database and Caddy configuration before upgrades while the project remains pre-1.0.
