# Input reconciliation — current disposition, 2026-10-17

Artem approved [P01–P13](../../sprint-change-proposal-2026-10-17.md). The PRD now propagates full existing configuration, all-settings forms and Caddyfile, shared drafts, four TLS scenarios, Docker/manual destinations and Quiet Workshop. FR identifiers remain stable. Research is complete; former candidate exclusions and generic discovery requests are superseded.

Architecture defines source/preservation, root conditional writes, durable attempts, secret history, session/probe revisions, migration and backup. [Epics](../../epics.md) and QA trace that corrected contract. Current-code facts remain in [baseline extraction](reconcile-code-baseline.md); the new design is future work. Any remaining build/corpus/harness absence is labelled planned; no runtime certification is claimed.

The prior reconciliation below is a dated historical record, superseded where it describes unsettled scope or research. It is not a current requirements decision.

---

# Input reconciliation — Caddy Orchestrator PRD

Reviewed: `prd.md`, `addendum.md`, `reconcile-code-baseline.md`, the supplied user messages, and current code previously extracted. Date: 2026-10-16. No runtime validation performed. Competitive research is still in progress and requires a final integration pass.

## Traceability

| Input / evidence | PRD treatment | Assessment |
|---|---|---|
| User: “довести этот проект … [до] хорошего состояния” in functionality/usefulness for homelabbers | §1 workflow connect → publish → understand → change → recover; FR-1–FR-10, SM-1/SM-2 | Directly represented. Reliability is tied to evidence rather than source presence. |
| User: “дизайн более современный” | FR-11, Q4, UX handoff; addendum UX collaboration contract | Represented without inventing a palette or treating visual decisions as approved. |
| User: “условная копия nginx proxy manager но для caddy” | §1 role comparison, §6 parity explicitly proposed as a non-goal, Q2 essential capability decision | Fair interpretation pending confirmation; does not silently convert “copy” into a promise of complete feature parity. The user must still decide which NPM-like capabilities are essential. |
| User authorized Deep Recon and PRD via subagents, with UX potentially separate | §1 draft/research status; §8 separate-chat handoff | Correctly represented. Authorization to draft is not treated as implementation approval. |
| `CLAUDE.md`, `internal/api/managed.go`, `internal/config/builder.go`: scoped Routes writes, ETags, snapshots, preservation | §3 existing baseline; FR-1/4/5/9 retain protections | Correct: does not repeat the prior assistant's outdated full-config-overwrite claim or schedule existing safety mechanisms as entirely missing features. |
| `internal/storage/models.go`, `internal/api/routes.go`: route CRUD/toggle, single instance/server, no TLS/access management | §3, FR-2/3/8, A1/Q2/Q3 | Correctly separates supported operations from candidate new work and scope decisions. |
| `internal/storage/sqlite.go`: 100 retained snapshots, 50 listed | §3 and FR-9 | Accurate; configurable policy and full backup remain distinct from existing fixed retention. |
| `cmd/server/main.go`, Compose files: loopback defaults, paired Basic Auth credentials, insecure override, private Admin API | §3 deployment baseline, §5 security, Q3 | Accurate; first-admin UI/users/RBAC are not incorrectly claimed. |
| `docker-compose.yml`, `compose.preview.yml`: static Caddyfile startup without `--resume`; UI has no startup sync | §3 unverified; FR-2, SM-2; addendum restart verification | Correctly escalated as a verification target, not reported as a demonstrated runtime failure. |
| `internal/api/managed.go`: PATCH may succeed before later ETag/local persistence error | FR-5/6; addendum observation distinctions | Correctly avoids claiming atomic cross-system transactions and calls for persistent honest outcome states. |
| `internal/storage/models.go`, builder: no certificate state, TLS modes, TLS upstream transport, health/access policies | FR-6/7/8 and Q2 | Correctly treats these as gaps or choices, not existing behavior. |
| Archived `specs/001`/`002` | §3 historical, source contract in addendum | Correct; old destructive local-only import semantics are not reintroduced. |
| Latest §3 research synthesis and §8 UX draft links | Competitor examples explicitly selected, anecdotal and non-prevalence evidence; DESIGN/EXPERIENCE identified as discovery drafts | Appropriate limits: research does not establish Caddy defects, and document existence does not settle product/design decisions. Final research report is still in progress. |

## Material gaps and unresolved interpretation

1. **No material confirmed user requirement is omitted.** The PRD deliberately remains a review draft: A1–A3 and Q1–Q3 surface the main inferred scope. The one-instance restriction and parity exclusions must not become implementation assumptions before those decisions.
2. **Research integration is pending.** The draft must update `research_status` and the supporting evidence after the competitive report is complete. Competitor capabilities should inform Q2 prioritization without automatically becoming FRs.
3. **Validation boundaries need a later decision artifact.** FR-3 proposes supported inputs and examples but does not yet define accepted host/path/upstream forms, proxy transport defaults, or redirect codes. This is suitable for a product discovery draft; it must resolve before engineering can judge complete input validation.
4. **Recovery acceptance should name both branches.** FR-9 says to verify restoration against “a changed live configuration”; it should explicitly distinguish successful restore with unchanged ETag from drift rejection after an external change. Current source already rejects the latter. This clarification prevents an accidental expectation that recovery bypasses existing protection.
5. **Broader homelabber usefulness remains a hypothesis.** SM-1 validates Artem's real workflow, appropriate for the proposed solo-lab first release. It cannot establish newcomer/public adoption without a later audience decision and representative users. Q1 honestly carries that limitation.

Verdict: faithful to the available user request and code baseline; ready for UX discovery and user scope decisions, not yet an implementation contract. No product-code changes are requested by this reconciliation.

## Final disposition — 2026-10-16

Research is now complete and the final report is integrated into the main PRD. The original review above records its earlier snapshot. The author addressed the input-contract finding with concrete destination examples and an explicit pre-story boundary decision; split restore success from drift rejection; and defined local draft, failed, and unknown without implying that failure leaves live configuration unchanged. The addendum also distinguishes ordinary pre-write Snapshots from empty bootstrap. No unresolved critical/high finding remains. This remains a reviewable proposal: Q1–Q5 and release-scope approval are still Artem's decisions.

## Latest direct-user decisions — focused update

The user clarified that the product serves homelabbers broadly including Artem; primarily uses HTTPS for public domains with a mode/settings choice; and deploys production Compose without Caddy, connecting to an existing instance. Local bundled Caddy is for development/testing. The PRD/addendum/status now reflect those decisions; Q1/Q3 and primary Q2 are resolved. A reproducible reference service replaces any mandatory personal-service acceptance dependency. The exact alternative mode/provider menu remains a nonblocking UX/architecture detail; it and the entire candidate FR set were not approved by this clarification. Historical review findings above describe the earlier draft. Focused structure/prose reconciliation preserves stable FR IDs and does not reopen settled choices.
