# Workflows

- `ci.yml`: Go build/vet/race tests, real-Caddy scoped-write integration test, frontend lint/typecheck/coverage/build, Docker build, and an isolated local preview smoke test on `main` pushes and pull requests. The preview is not published.
- `security.yml`: blocking gosec, govulncheck, npm audit, Trivy filesystem/image scans, plus scheduled CodeQL.
- `release.yml`: release-please and multi-architecture GHCR images from `main`.

Toolchain versions match `mise.toml` and the Dockerfile: Go 1.25.13 and Node.js 24.

Run the main checks locally with:

```bash
make test
cd web && npm run lint && npx tsc --noEmit && npm run test:coverage && npm run build
docker build -t caddy-admin-ui:local .
```
