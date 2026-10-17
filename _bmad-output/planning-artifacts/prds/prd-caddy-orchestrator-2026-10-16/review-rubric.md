# PRD Quality Review — Caddy Orchestrator for Homelab Operators

## Overall verdict

This is a useful and unusually honest brownfield review draft: it connects the homelab goal to observable service operation, preserves implemented safety mechanisms, and exposes the important audience/TLS/deployment trade-offs. It is ready to ground a UX conversation and scope decisions; selected scenario boundaries and input examples still need resolution before implementation stories. Competitive research is in progress, so this verdict does not certify its final integration.

## Decision-readiness — adequate

§1 explicitly distinguishes confirmed direction from proposed scope; §6 and Q1–Q3 identify the main decisions rather than smoothing them into neutral prose. The precedence of persistence/recovery over ownership expansion is actionable. A1–A3, `status: draft`, and §8 prevent the apparent completeness of FR-1–FR-11 from implying user approval.

The remaining uncertainty is real product uncertainty, not a defect to solve by inventing answers. Research should help Artem select Q2 scenarios; it cannot make that choice on his behalf.

## Substance over theater — strong

The document avoids fabricated personas, novelty claims, adoption metrics, and generic performance promises. Admin API availability versus service verification, local/internal trust, immutable snapshot scope, and restart persistence all follow directly from this application's failure modes. §5 explicitly leaves resource budgets unclaimed until a representative environment is selected.

## Strategic coherence — strong

The thesis in §1 — connect → publish → understand → change → recover — supplies a consistent priority for setup, publication, state honesty, diagnostics, and restore. SM-1 evaluates a complete useful operator task, SM-2 evaluates trustworthiness, and SM-C1 prevents cosmetic simplification from hiding risk. Preservation of existing route kinds avoids a redesign regression while the conditional TLS/access scope remains visible.

## Done-ness clarity — adequate

FR-1/2/4/5/6/9/10 have concrete observable consequences, including independent restarts, guarded conflicts, unknown outcomes, clean restore, and prior-release upgrade. FR-7 sensibly permits documented next checks where universal diagnosis is impossible. FR-11 gives meaningful keyboard, label, focus, and error-retention requirements without pretending to have validated a new visual style.

### Findings

- **medium** Supported-input acceptance remains underspecified (§4 FR-3; §8 Q2) — The draft asks for validation and examples but does not identify the accepted host/path/upstream syntax or redirect-code set. HTTPS transport is called out, but domains with schemes/ports, wildcard handling, and malformed dial strings could receive inconsistent behavior across UI and API. *Fix:* after Q2/Q3, add a compact supported-input table or link a contract with accepted and rejected examples, transport defaults, and error expectations; keep this at a boundary-setting level rather than prescribing a validator implementation.
- **low** Restore acceptance can be read two ways (§4 FR-9) — “Verify restoration against … a changed live configuration” follows a conflict-guarded requirement but does not explicitly say whether an external edit must block the restore. *Fix:* name two cases: an unchanged-current-state restore succeeds and creates a pre-restore Snapshot; external drift rejects restore until reconciliation, preserving live configuration.

## Scope honesty — strong

§6 states proposed non-goals and marks conditional NPM-like capabilities that may change the scope. Q2 explicitly acknowledges that access policies, HTTPS upstreams, LAN CA, and wildcards may be essential rather than quietly omitting them. The assumptions index matches the three inline assumptions. The draft does not transform the prior assistant's recommendations or competitors' features into accepted user requirements.

The density of open items is proportionate for this early solo-homelab draft. It would be too high for a green-light-to-build contract, which §1/§8 expressly say this is not.

## Downstream usability — adequate

FR-1–FR-11, SM-1/SM-2/SM-C1, and Q1–Q5 are stable and cross-referenced. Terms in §2 distinguish Route, Managed Server, External Route, Snapshot, applied, and service verified; the addendum explains the same boundaries. Separate UX output ownership and a final reread of PRD/research reduce concurrent-chat conflicts. Named UJs are unnecessary for the chosen single-operator capability shape.

### Findings

- **low** FR-6 state vocabulary needs one authoritative definition (§2 terms; §4 FR-6) — “local draft,” “applied,” “failed,” and “unknown/reconciliation-required” appear as product outcomes, but only applied has an explicit definition. *Fix:* define whether failed means rejected before application, known retained prior state, or any attempted operation; ensure a failed operation cannot falsely imply unchanged live configuration after an ambiguous remote result. A small state glossary is enough; the architecture can define persistence mechanics later.

## Shape fit — strong

A compact capability PRD fits a solo operator tool with a meaningful UX refresh and substantial brownfield safety constraints. Existing versus missing versus unverified behavior is kept explicit in §3 and the code extraction. The document adds release evidence without imposing enterprise personas, telemetry, a large architecture template, or a fabricated timeline.

## Mechanical notes

- FR, SM, Q, and assumption identifiers are unique and resolve. A1–A3 roundtrip with the assumptions index.
- Relative source, research, and addendum links resolve structurally from the PRD workspace; research content still needs the final integration pass.
- Capitalization of Route/Upstream is mostly consistent; ordinary generic plural “route array” is understandable and does not change the model.
- No named user journeys are present by design; §2 explicitly says why and delegates observed workflows to UX elicitation.
- No critical or high issue was found. Resolve product decisions before treating the draft as approved first-release scope.

## Final disposition — 2026-10-16

Research is now complete and the final report is integrated into the main PRD. The original review above records its earlier snapshot. The author addressed the input-contract finding with concrete destination examples and an explicit pre-story boundary decision; split restore success from drift rejection; and defined local draft, failed, and unknown without implying that failure leaves live configuration unchanged. The addendum also distinguishes ordinary pre-write Snapshots from empty bootstrap. No unresolved critical/high finding remains. This remains a reviewable proposal: Q1–Q5 and release-scope approval are still Artem's decisions.
