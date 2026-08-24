# Contributing

Thanks for improving Caddy Admin UI.

## Before opening a change

- Use an issue or discussion for broad product changes so the ownership and safety model stays coherent.
- Keep imported external Caddy JSON lossless and read-only.
- Never replace the full Caddy config for a routine route operation.
- Treat drift conflicts as a stop condition, not something to overwrite automatically.

## Local checks

Install Go 1.25.13 and Node.js 24 (`mise install` is supported), then run:

```bash
make deps
make test
cd web && npm run lint && npx tsc --noEmit && npm run test:coverage && npm run build
```

Docker changes should also pass:

```bash
docker build -t caddy-admin-ui:local .
docker compose config
make preview
make preview-down
```

## Pull requests

- Use a conventional-commit title such as `feat: add route reordering` or `fix: preserve matcher metadata`.
- Explain the user-visible behavior, safety impact, tests, and any migration.
- Add focused backend and frontend tests for behavior changes.
- Do not include generated frontend output or `node_modules`.

Pull requests are squash-merged, so the PR title becomes the release commit.
