# Agent Instructions

Full project documentation (build commands, architecture, environment variables) lives in [CLAUDE.md](CLAUDE.md) — read it first.

## Commit Conventions

Releases are automated with release-please: version bumps and the changelog are computed from commit messages, so every commit that lands on `main` MUST follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: ...` — new feature (minor version bump)
- `fix: ...` — bug fix (patch version bump)
- `feat!: ...` or a `BREAKING CHANGE:` footer — breaking change (minor bump while pre-1.0)
- `chore: ...`, `ci: ...`, `docs: ...`, `refactor: ...`, `test: ...` — no version bump, hidden from the changelog

PRs are squash-merged, so the **PR title becomes the commit message on main** — it must follow the convention too. Scopes are optional (`feat(api): ...`). Releases ship when the auto-generated release-please PR is merged.
