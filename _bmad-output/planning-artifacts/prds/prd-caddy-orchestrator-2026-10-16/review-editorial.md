# Editorial review — Caddy Orchestrator PRD

These documents help Artem and downstream UX/implementation contributors decide the scope of a dependable homelab Caddy manager and preserve its operating constraints.

Reviewed `prd.md` and `addendum.md` as human-facing documents using the resolved `bmad-review` configuration: **structure first, then prose using structure findings**, baseline **Microsoft Writing Style Guide**, no custom standing directives. Structure model: **Strategic/Context (Pyramid)** for the PRD; **Reference/Database** for the supporting notes. No length target was supplied. Content and product decisions were not challenged or edited. The PRD remains a proposal, not an approved implementation contract.

## Exact word metrics

Computed with `uv run .agents/skills/bmad-review/scripts/word_metrics.py <path>` on 2026-10-16. Script counts include frontmatter as preamble and exclude heading text; proposed savings below are estimates grounded in these exact counts. Parent may subsequently update research status or accept edits, so these counts describe the reviewed snapshot.

| Document / section | Words |
|---|---:|
| **prd.md total** | **2,219** |
| Preamble | 18 |
| §1 Proposal and status | 156 |
| §2 Target user and jobs | 211 |
| §3 Current baseline and product gap | 271 |
| §4 Candidate first-release requirements — introductory text | 37 |
| Setup and publication | 216 |
| Safe changes and honest feedback | 251 |
| HTTPS, recovery and operation | 219 |
| Modern interface | 84 |
| §5 Cross-cutting release constraints | 145 |
| §6 Proposed scope boundary | 110 |
| §7 Success and release acceptance | 132 |
| §8 Decisions and UX handoff | 299 |
| **addendum.md total** | **561** |
| Preamble | 14 |
| Source boundaries | 114 |
| Architecture constraints to carry forward | 137 |
| UX collaboration contract | 113 |
| Verification considerations for downstream planning | 161 |
| **Combined total** | **2,780** |

## Findings

The structure pass found that the PRD already opens with the proposal and status, then groups baseline, requirements, constraints, acceptance, and decisions coherently. The main opportunity is to show proposed scope before the detailed requirements. In the prose pass, preserve direct requirements, deliberate confirmed/proposed/unverified distinctions, technical terms needed for precision, and concrete examples; replace only phrases that slow comprehension. Prose suggestions below avoid overlapping or contradicting the structural condensation.

| Pass | Original Text | Revised Text | Changes |
|---|---|---|---|
| structure | `prd.md` §6 Proposed scope boundary (110 words) follows the complete FR list and release constraints. | **MOVE** §6 before §4 Candidate first-release requirements; renumber section references consistently if accepted. | A decision-maker sees the core, conditional scope, and exclusions before reading individual FRs. Word impact: 0; no content removed. Keeping §6 in place is also defensible if stable section numbers matter more than scan order. |
| structure | `addendum.md` §UX collaboration contract (113 words) repeats the main PRD's UX handoff and visual-scope boundary before giving its distinct evidence/assumption instructions. | **CONDENSE** to: “Use a separate UX workspace and keep PRD requirement IDs stable. Re-read PRD and research before finalizing DESIGN/EXPERIENCE. The PRD defines behavioral goals; UX defines visual and interaction details. Treat illustrative sessions as hypotheses until Artem confirms them. Where scope is unresolved, show alternatives and dependencies rather than choosing by visual convenience.” | Retains independent operational guidance, stable IDs, source refresh, hypothesis status, and scope dependencies. Main §8 already carries the elicitation checklist and output links. Saves approximately 60 of the section's 113 words. Trade-off: the addendum reader must follow the main PRD for the complete UX task checklist. |
| structure | `prd.md` §3 baseline table and FRs repeat ownership, snapshots, and application/service-state distinctions (within the 271-word §3). | **PRESERVE** the baseline table and the corresponding FR protections. | The table answers “what exists?”; the FRs answer “what must hold in the release?” This is useful reinforcement rather than duplication. Word impact: 0. |
| prose | `prd.md` §2: “The working operator wants to expose a home service through a memorable domain” | “The operator wants to expose a home service through a memorable domain” | Removes an unclear qualifier without changing the audience or job. |
| prose | `prd.md` FR-5: “An external edit causes a conflict before candidate local persistence, with a review/reconciliation path.” | “An external edit causes a conflict before the candidate Routes are saved locally, with a review/reconciliation path.” | Replaces a compressed noun phrase with an explicit action; retains the exact local-persistence boundary. |
| prose | `prd.md` §5: “successful failure/recovery acceptance must demonstrate that unrelated configuration survives” | “failure and recovery checks must demonstrate that unrelated configuration survives” | Removes a phrase that momentarily reads as a contradiction and states the acceptance activity directly. |
| prose | `addendum.md` §Verification considerations: “Record whether each observation is a source assertion, automated test, manual runtime check, or unresolved.” | “For each claim, record whether it is supported by source inspection, an automated test, a manual runtime check, or remains unverified.” | Uses parallel evidence categories and distinguishes the claim from its supporting test. Preserves all four evidence states. |

## Reduction summary

Seven recommendations: two structural changes, one explicit preserve decision, and four prose fixes. If all changes are accepted, the combined documents shrink by approximately **50 words (1.8% of 2,780)**: about 60 words from condensation, partly offset by clearer prose. The move and preservation do not change length. No length target applies; the only comprehension trade-off is that the shorter addendum relies on the linked main PRD for the full UX elicitation checklist.

The final snapshot includes integrated research status, explicit local-draft/failed/unknown definitions, and the bootstrap Snapshot exception. All Q1–Q5 decisions remain with Artem; editorial suggestions do not resolve them. No source documents were modified by this review.
