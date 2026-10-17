# Course correction — technology reality check

Date: 2026-10-17. Reviewed `ARCHITECTURE-SPINE.md` and normative `configuration-contracts.md` against current repository manifests and primary Caddy sources. No input documents edited and no runtime tests or deployments executed.

**Verdict: pass with one low-severity seed correction.** No substantive API-contract error found in the reviewed root-write, adaptation, source/import, build-evidence or pinned-reference contracts. Planned compatibility fixtures remain delivery gates, not completed evidence.

## Finding

### T1 — Low: Tailwind stack seed reports the range floor instead of installed lock version

Location: `ARCHITECTURE-SPINE.md:131`.

The seed says 3.4.1, while `web/package-lock.json:5142` locks 3.4.19. `web/package.json` declares `^3.4.1`, which is a range rather than evidence that 3.4.1 is installed. Since the table otherwise reports resolved frontend versions and says it ratifies the manifest/lock, use 3.4.19, or explicitly label the range and lock separately. This is documentation-only and authorizes no dependency update.

## Verified technical decisions

- **Conditional whole-root replacement is valid for the pinned target.** Caddy 2.9.1 forwards the supplied If-Match value from the config handler into `changeConfig`; the latter locks shared raw configuration before checking the precondition and mutation. At the root map key, POST sets/replaces a non-array value, so `/config/` supports replacement and null-root bootstrap. The proposed root GET ETag and exact If-Match contract is supported by implementation, including rejection with 412 after a competing mutation. Integration tests are still necessary to establish the deployed fixture behavior. [Pinned handler](https://github.com/caddyserver/caddy/blob/v2.9.1/admin.go#L891-L959), [root assignment](https://github.com/caddyserver/caddy/blob/v2.9.1/admin.go#L1102-L1124), [locked precondition](https://github.com/caddyserver/caddy/blob/v2.9.1/caddy.go#L132-L244).
- **Adaptation boundaries are accurate.** The pinned handler accepts the content-type-selected adapter, invokes Adapt with nil options, and returns an object with `result` and optional `warnings`. No query-parameter semantic-validation branch exists. Neither successful adaptation nor root config retrieval demonstrates module inventory or successful provisioning. [Pinned adaptation implementation](https://github.com/caddyserver/caddy/blob/v2.9.1/caddyconfig/load.go#L120-L193).
- **Captured virtual imports require the declared implementation work.** Caddy 2.9.1 expands environment values before tokenization and resolves imports against token filenames using the target filesystem, including ordered glob handling. X1 correctly requires manager-side captured expansion, a transient import-free representation, original encrypted sources and target-parity fixtures. This is an implementable replacement of unavailable target filesystem access, not a capability of `/adapt` itself. Fixtures should include substitution values containing Caddyfile syntax or another environment token, import/snippet argument expansion and unmatched globs; until parity is proved the existing block-on-unsupported rule applies. [Pinned parser and environment processing](https://github.com/caddyserver/caddy/blob/v2.9.1/caddyconfig/caddyfile/parse.go#L27-L103), [pinned import resolution](https://github.com/caddyserver/caddy/blob/v2.9.1/caddyconfig/caddyfile/parse.go#L380-L445).
- **Build evidence is honestly bounded.** Operator version/module/image/build-lock capture is distinct from API evidence and unknown remains unknown. `list-modules --versions` exists; its output can omit version information for main-module builds, so the supplied build lock/digest and unknown-state handling remain necessary. [Official CLI documentation](https://caddyserver.com/docs/command-line#caddy-list-modules).
- **PKI behavior is accurate.** The pinned certificates endpoint emits intermediate followed by root, and can provision the default local CA. Root selection/validation and explicit disclosure correctly avoid treating the endpoint as a generic side-effect-free probe. [Pinned PKI handler](https://github.com/caddyserver/caddy/blob/v2.9.1/modules/caddypki/adminapi.go#L130-L213).
- **Durability is an external deployment prerequisite.** Persistent config/data and resume/equivalent are correctly required; connectivity cannot certify startup flags or mounts. The existing `docker-compose.yml` currently uses Caddy 2.11.4 and file-based startup without resume, so it is neither the promised 2.9.1 fixture nor evidence of restart durability. The documents explicitly call the fixture planned, so this is delivery work rather than a false certification. [Official run/resume semantics](https://caddyserver.com/docs/command-line#caddy-run).

## Existing-manifest reality check

Confirmed Go 1.25.13, Gin 1.12.0 and go-sqlite3 1.14.50 in `go.mod`; Node 24 in `mise.toml` and Dockerfile; locked Preact 10.29.8, Preact Router 4.1.2 and Vite 8.2.2; Alpine runtime 3.24 in Dockerfile. These are inherited version seeds, so no latest-version recommendation is inferred. Signals is correctly labelled uninstalled with compatible-version delivery assigned to Story 1.10.

The DNS provider/build lock is also explicitly assigned to Story 1.5. No Cloudflare dependency is currently locked by this repository and no particular provider release is asserted as proven compatible. Do not choose an unversioned latest provider during implementation: build the pinned 2.9.1 target with the selected module graph and record the resolved core/module versions, binary/image digest and compatibility evidence. This requirement is already part of X3 and does not reopen the approved provider selection.

## Disposition

Autofix T1 in the seed. Retain existing planned fixture gates for root preconditions, source expansion, provider provisioning and external restart; this review establishes source-level plausibility, not runtime completion.


## Follow-up disposition

The author updated the Tailwind seed to the existing lock resolution 3.4.19. The low-severity documentation correction is applied. No product dependency was changed.
