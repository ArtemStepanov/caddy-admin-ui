---
project: caddy-orchestrator
date: 2026-10-17
intent: sprint-planning
status: complete
gate: PASS
tracking_generated: true
method: corrected-contract-and-implementability-review
runtime_certification: false
supersedes: implementation-readiness-before-course-correction.md
---

# Implementation readiness — Caddy Orchestrator

**Verdict: PASS for planning and sprint generation.** Approved P01–P13 have been propagated across requirements, architecture/technical contracts, UX, SPEC, 33 stories and QA. The previously missing shared decisions now have concrete rules and acceptance mappings. This permits implementation planning, not a claim that the expanded product already works.

The [prior FAIL report](implementation-readiness-before-course-correction.md) remains as historical evidence. Artem approved the [change proposal](sprint-change-proposal-2026-10-17.md) in this conversation. Scope decisions are retained; generic product discovery is closed.

## Inventory and implementability

| Artifact | Current role / evidence |
| --- | --- |
| [PRD](prds/prd-caddy-orchestrator-2026-10-16/prd.md), addendum and reconciliation | Confirmed full-config/forms/editor/draft/TLS/discovery scope propagated with stable FR-1–FR-11 and explicit authority. Current source is separately identified as legacy baseline. |
| [Architecture spine](architecture/architecture-caddy-orchestrator-2026-10-17/ARCHITECTURE-SPINE.md) and [shared contracts](architecture/architecture-caddy-orchestrator-2026-10-17/configuration-contracts.md) | Stable AD-1–AD-11; root precondition, captured source/paths, durable attempts/locks, protected history, session/probe identity, supported forms, migration/backup mechanisms recorded. |
| [DESIGN](ux-designs/ux-caddy-orchestrator-2026-10-16/DESIGN.md), [EXPERIENCE](ux-designs/ux-caddy-orchestrator-2026-10-16/EXPERIENCE.md) | Quiet Workshop retained; scope/technical unknowns reconciled; visual assumptions remain labelled rather than silently approved by plan status. |
| [SPEC](../specs/spec-caddy-orchestrator/SPEC.md), glossary/diagrams | Same publication/source/recovery/identity contract; WCAG 2.2 AA and manager port 3000 aligned; no unsupported transactional/fidelity claim. |
| [Epics](epics.md) | Six epics / 33 stories / 122 numbered ACs, explicit backward prerequisites and old→new aliases. First complete service slice includes safety/recovery. |
| [QA plan](../test-artifacts/test-design-qa.md), [architecture test plan](../test-artifacts/test-design-architecture.md), [handoff](../test-artifacts/test-design/caddy-orchestrator-handoff.md) | 122 uniquely AC-linked scenario intents plus original-finding mapping; test seams assigned before consumers; all tooling/results labelled planned. |
| [Independent review dispositions](architecture/architecture-caddy-orchestrator-2026-10-17/review-dispositions.md) | Fresh rubric, tech-fit and adversarial reviews; specific findings corrected and follow-up closure recorded. Earlier reviews preserved as historical. |
| [Structural validation](course-correction-validation.json) | Unique IDs, acyclic/backward explicit dependencies, full QA/AC equality, all 15 findings mapped, stable FR/AD IDs and local link existence. |
| Current CLAUDE/AGENTS/source/Compose | Existing runtime baseline; no product implementation or dependency change made in this correction. |
| Implementation tracking | Deterministically generated after PASS: six epics / 33 stories, all backlog; no standalone implementation story files or runtime progress are inferred. |

Could a developer implement the epics without inventing unrecorded shared decisions? **Yes, within the explicit declared compatibility boundary.** Source data acquisition, preservation collisions, import expansion, canonical identities, target write mode, concurrent writers, durable attempts, draft saves, restore, secrets, observations, migration and remote backups are specified. Unknown/unrepresentable data blocks safely with protected originals; arbitrary JSON inversion/plugin support is never assumed. Every registered setting of the declared supported build must obtain a form/codec/validation/secret fixture before that build can be declared fully supported. The reference DNS provider/build-lock delivery is assigned, not a missing provider scope decision.

## Original findings disposition

All items are **resolved in the plan**. The ACs below define implementation proof still to be collected, not executed results.

| Finding | Recorded resolution | Current acceptance references |
| --- | --- | --- |
| B01 | Atomic root precondition and reviewed rebase/restore; AD-3/8; companion P02/X1/X2 | 1.8.AC2, 1.8.AC3, 1.9.AC2, 6.3.AC2 |
| B02 | Source association and nested preservation; AD-1/2; companion P03/X1 | 1.4.AC2, 1.4.AC3, 1.6.AC1, 1.6.AC2, 1.6.AC3 |
| B03 | Durable attempts and crash/finalization reconciliation; AD-3/10; companion P04/X2 | 1.9.AC1, 1.9.AC3, 1.9.AC4, 1.9.AC6, 3.5.AC3 |
| B04 | Protected source and historical secret hydration; AD-5/11; companion P05/X2 | 1.2.AC1, 1.2.AC2, 1.2.AC3, 1.2.AC4, 4.2.AC3 |
| B05 | Forward implementable complete service slice; epics delivery order / aliases | 1.10.AC1, 1.10.AC2, 1.10.AC3 |
| B06 | All-settings forms and retained baseline features; AD-1/9; companion P07/X3 | 2.1.AC1, 2.1.AC2, 2.1.AC4, 2.5.AC1, 2.5.AC2, 2.5.AC3 |
| B07 | Draft/session/probe revisions; AD-4/7; companion P06/X2 | 1.7.AC2, 1.7.AC4, 3.4.AC1, 5.1.AC3, 5.4.AC2 |
| B08 | Bounded typed and browser evidence; AD-4; companion P08 | 5.1.AC2, 5.2.AC3, 5.3.AC1, 5.3.AC2, 5.3.AC3, 5.5.AC1 |
| B09 | Pinned Caddy request/response and build evidence; AD-9; companion P08/X3 | 1.5.AC2, 1.5.AC3, 1.5.AC4, 4.2.AC2, 4.3.AC2 |
| B10 | Functional transactional migration; AD-11; companion P09/X4 | 1.3.AC1, 1.3.AC2, 1.3.AC3, 1.3.AC4, 6.4.AC3 |
| C01 | Scope and artifact authority reconciled; PRD 1/4/8; SPEC authority | 2.5.AC1, 4.2.AC1, 4.3.AC1, 4.4.AC1 |
| C02 | Current test/story mapping and honest gates; QA matrix / architecture test design | 1.1.AC3, 1.5.AC4, 1.9.AC3, 1.10.AC4 |
| C03 | Complete backup and remote operator acquisition; AD-11; companion P10/X4 | 6.4.AC1, 6.4.AC2 |
| C04 | Identity and input edge contracts; AD-1; companion P11/X3 | 1.4.AC1, 1.4.AC4, 2.1.AC3, 4.4.AC3 |
| C05 | Previous review dispositions and consistent targets; review-dispositions.md; SPEC diagrams | 1.5.AC4, 1.10.AC4, 2.4.AC3 |

## Independent review corrections

Fresh adversarial review found immutable connection and enabled-metadata seams; the rubric also found an early UI dependency wording issue. These are corrected in the spine, companion and ACs: attempts/history pin endpoint/transport descriptors, source bundles include historical enabled/node maps, preparation acquires a durable target/draft lock with CAS and finalization cannot erase a later revision, and vault tests use a port/API harness before UI delivery. Follow-up reviewers report no remaining planning blocker. Technology check confirmed target API/source semantics and corrected Tailwind seed to the existing lock resolution. Architecture mechanical lint has zero findings.

## Remaining implementation work and evidence limits

Source corpus, exact module/image/build locks, form inventory, Signals installation, target/fault/E2E harnesses and all 122 scenario executions are required story deliverables. Their absence before development is not a missing decision or fabricated Ready status. No real product tests, builds, migrations, HTTPS issuance/renewal, restore, deployment, rendered accessibility or backup execution ran during this correction. Primary-source verification and document validation are the evidence gathered here.

No performance/hardware/route-count/CI-duration budget has been supplied; measure before binding one. Full-config editing remains constrained by Caddyfile representability and verified path mappings; expose incompatibility rather than dropping data or substituting a JSON editor. Reference Caddy 2.9.1 is a pinned compatibility fixture, not a recommendation to upgrade a target. Operator-side remote asset acquisition and external startup prerequisites are documented responsibilities rather than invented manager capabilities.

## Next step

Sprint tracking is generated and validated mechanically from current epics; all stories remain backlog. Developer begins **1.1 Secure manager-only startup and application seams**, then follows prerequisites. Story readiness/done states require real implementation artifacts/evidence, not this planning PASS. Future findings discovered in execution must update contracts/aliases/tests consistently before conflicting implementations proceed.
