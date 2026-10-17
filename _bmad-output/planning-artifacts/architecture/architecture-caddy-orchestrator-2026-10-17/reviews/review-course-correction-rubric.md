# Course-correction architecture rubric review

Date: 2026-10-17
Reviewer: independent rubric walker
Verdict: **Pass with two bounded clarifications; no critical or high findings.** Resolve R1 in the shared contract before its storage/source/recovery implementations are split. R2 is a story-boundary clarification and does not require changing product scope.

## Scope and method

Reviewed corrected `ARCHITECTURE-SPINE.md`, `configuration-contracts.md`, PRD, SPEC and all six epics/33 stories against the good-spine checklist in `bmad-architecture/references/reviewer-gate.md`. Read `CLAUDE.md`, `AGENTS.md`, `go.mod`, `mise.toml`, `web/package.json` and `Dockerfile` for the brownfield boundary and stack reality check. No product/planning artifacts were edited. QA documents were being edited separately and were outside this review.

The deterministic spine lint passed with zero findings. This review assesses planning coherence, not runtime certification. It accepts P01–P13 and the user-confirmed scope as binding. The root mutation protocol is `POST /config/` with the reviewed Caddy root `If-Match`, never an unconditional `/load` fallback.

## Rubric results

| Criterion | Result | Assessment |
| --- | --- | --- |
| Fix real divergence points for the next level | Pass with R1 | Source ownership, structural composition, immutable imports, guarded root writes, interruption recovery, shared sessions and typed observations have concrete cross-component contracts. Enabled metadata needs an explicit immutable envelope association. |
| Rules enforce their stated prevention | Pass | Fail-closed cases are named: unmapped overlap, mixed arrays, missing secret material, stale review, remote precondition failure and unprovable reconciliation. No distributed transaction or historical-causation guarantee is asserted. |
| Deferred work leaves no contradictory implementations | Pass | Non-goals are separated from required capability. Arbitrary plugins/JSON inversion are compatibility limits; complete supported forms and all four TLS scenarios remain required. Measurements do not silently become invented release budgets. |
| Technology reality check | Pass at planning altitude | Existing stack versions are ratified against project manifests. Signals installation/version verification belongs explicitly to 1.10; provider build locks belong to 1.5. These are assigned deliverables, not fabricated completed evidence. No latest-version assertion was used. |
| Brownfield continuity | Pass | Legacy and root writers are mutually excluded. Migration precedes adoption, preserves disabled/local-only payloads and scoped snapshots, and has a prior-version backup/downgrade route. Incremental seams preserve current implementation while replacement is built. |
| PRD/SPEC capabilities preserved | Pass | Whole adopted configuration, forms plus Caddyfile, shared staged changes, four TLS scenarios, optional read-only Docker, manager-only Compose, Quiet Workshop, desktop/phone, recovery and full backup are represented. The first service slice is explicitly not the full release. |
| No forward prerequisites | Pass with R2 | Every explicit dependency points backward. Foundation outcomes are testable at their ports/APIs; however, 1.2.AC2 literally invokes later form/editor integration despite its stated backend-only boundary. |
| Operational/environmental envelope | Pass | Private Admin API, manager/external-Caddy separation, startup resume/persistence, encrypted backup assets, operator-side remote acquisition, restart, clean restore and downgrade limits are addressed. API access is not misrepresented as filesystem or Docker control. |
| Inherited invariant conflicts | Not applicable | This is the corrected feature-level spine; no separate parent spine was supplied. Stable AD-1–AD-11 IDs remain in place. |

## R1 — Medium: bind enabled metadata to the immutable compilation/recovery envelope

**Evidence:** Companion X1 defines the `SourceBundleV1` contents, including exact files, imports, substitutions, source maps and build identity. X3 says enabled metadata is persisted "alongside source nodes" and disabled source remains in the bundle but is omitted from compilation. X2's shared records do not explicitly require a version/digest or immutable snapshot of enabled metadata in `SOURCE_BUNDLE`, `APPLY_ATTEMPT` or `RECOVERY_POINT`. P11 calls the metadata versioned, but does not define what binds that version to source/review/recovery identity.

**Concrete divergence:** A source implementer can put enabled flags in a mutable node table, while a recovery implementer copies the exact source bundle and runtime JSON listed in X2. Both retain disabled source and encrypted snapshots. After a route is enabled later, restoring an older point may restore its runtime bytes but regenerate a different candidate on the next edit because it reads today's enabled flags. A review that binds only source-file/bundle identity can also miss an enabled-state change if flags are external to that identity. The high-level requirements clearly intend preservation; the missing detail is the shared envelope relationship that enforces it.

**Clear fix:** State that immutable bundle identity covers node identities and the enabled-state map/version as well as files/dependencies; every review, durable attempt, applied revision and full-config recovery point retains the exact bundle identity needed to reproduce compilation. Restoring that point creates a new draft/applied revision with the historical enabled-state snapshot. Alternatively define a separate immutable node-state record and require its ID/digest in the same chain. Neither approach requires a new product decision.

**Verification allocation:** Extend 1.4.AC4 and 1.9 restore/reconciliation acceptance with a route disabled in recovery A, enabled in B, restored to A, followed by another unrelated edit/recompile/restart. The route stays disabled with the same retained source/identity. This is additional precision for planned implementation, not a demand for the fixture to exist now.

**Triage:** Autofix the shared contract before separately implementing source/storage/recovery.

## R2 — Low: make story 1.2's secret round-trip acceptance match its foundation boundary

**Evidence:** 1.2 depends only on 1.1. Its AC2 requires a draft round-trip "through forms and editor", while its completion boundary says no later UI story is a prerequisite. The source/codec engine arrives in 1.4, minimum workspace in 1.10, full forms in 2.1 and peer editor in 2.2. The explicit graph is acyclic; this issue is the literal acceptance wording.

**Concrete divergence:** One builder can mark 1.2 complete after testing only the secret store's keep/replace/remove operations; another can insist the forms/CST/editor exist to execute AC2 and pull later work forward. The completion boundary expresses the intended former scope but does not specify the foundation round-trip mechanism.

**Clear fix:** Define 1.2.AC2 as a port/API contract-harness round-trip of typed secret references and operations, independent of actual forms/editor. Retain the real UI round-trip in 2.2.AC4 and the source/codec cases in 1.4. No new dependency or removal of secret safety is needed.

**Triage:** Autofix story wording/test allocation; not an architectural blocker.

## Explicitly not findings

- The new CST engine, full descriptor inventory, pinned provider fixture, protocol harness, fault hooks, Signals dependency and runtime acceptance evidence do not exist yet. Their delivery is assigned to concrete stories; that is expected for a build substrate.
- Identity-only initial runtime mapping and blocking modifications to mixed generated/preserved arrays are explicit safe compatibility behavior, not permission to omit the full-settings release goal.
- Remote source capture and backup are operator-supplied/read-only procedures. The documents do not need a Docker or filesystem mutation integration to satisfy them.
- The corrected stack is inherited from the codebase. A newer upstream version is not by itself a reason to reopen the approved no-stack-migration scope.
- Older pre-correction review files contain historical claims that this corrected review does not ratify. The reviewed revision and current companion are the basis of this verdict.

## Handoff assessment

The corrected architecture is implementable as a staged build plan with fail-closed compatibility boundaries. No missing planning decision requires renewed product discovery. R1 closes a cross-component persistence identity seam; R2 removes an avoidable interpretation of the early acceptance boundary. Neither changes approved scope, architecture paradigm, root protocol or six-epic delivery order.

## Follow-up disposition — 2026-10-17

**Final verdict: PASS for planning handoff. R1 and R2 are closed; no remaining implementation-decision blocker was identified in the reviewed corrected inputs.** This supersedes the initial conditional verdict above while preserving the review history.

- **R1 closed:** Companion X1 now makes immutable `enabled_nodes` and `node_identity_map` part of `SourceBundleV1` identity/digest and explicitly binds the enabled map to review, attempt and recovery payloads. Historical restore uses the historical map. X2 requires those maps in the encrypted bundle; spine AD-1 carries the invariant. Story 1.4.AC4 requires versioned metadata round-trip and historical-map preservation; the QA matrix links that AC to `CC-1.4-04`. Its source-engine test can exercise immutable bundle restoration without requiring the later remote restore coordinator.
- **R2 closed:** 1.2.AC2 now explicitly uses the secret port/API harness and assigns literal forms/editor integration to 2.1/2.2. Its dependency and backend completion boundary agree. QA preserves the current AC reference in `CC-1.2-02`.
- **Related corrections checked:** X2 now requires immutable connection descriptors retained by history/attempts; readback/reconciliation dereference the captured descriptor. Target/draft locks, revision recheck during preparation and finalization CAS prevent concurrent saves or connection changes from changing the operation's meaning. Spine AD-2/AD-3 carry these rules. Tailwind 3.4.19 matches the actual package lock. The explicit per-story dependency graph remains backward-directed.
- **Input consistency rechecked:** Corrected PRD/SPEC outcomes remain compatible with the tightened source/session/recovery rules and preserve all confirmed scope. The rewritten QA document labels its 122 AC-linked scenarios as planned acceptance and maps B01–B10/C01–C05 without asserting runtime execution. Missing fixture/corpus/provider-lock implementation remains assigned build work, not an unresolved architectural choice.
- **Mechanical recheck:** Spine lint again returned `ok: true`, zero findings.

No product/planning input was edited during this follow-up. This verdict establishes coherent planning contracts; product implementation, target compatibility and release acceptance still require the assigned evidence.
