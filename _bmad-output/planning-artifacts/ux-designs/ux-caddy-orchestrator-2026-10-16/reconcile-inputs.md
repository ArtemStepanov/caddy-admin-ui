# Current reconciliation disposition — 2026-10-17

The [approved course correction](../../sprint-change-proposal-2026-10-17.md) propagates U1–U23 decisions into PRD/architecture/SPEC/epics/QA. Ownership, forms+Caddyfile, shared batch drafts, four TLS scenarios, Docker/manual destinations and Quiet Workshop are settled; no scope question below is reopened. Architecture contracts fix the technical source/preservation/revision/probe/migration/backup handoff. Build locks/corpus/harness and product evidence remain planned implementation deliverables. Detailed visual assumptions remain labelled.

The original reconciliation below is preserved as historical input provenance; its divergence/pending wording is superseded by the current PRD and architecture contract.

---

---
title: Caddy Orchestrator UX input reconciliation
status: draft
created: 2026-10-17
source_reviewed: 2026-10-17
decision_status: product-and-architecture-reconciliation-required
---

# UX input reconciliation

This records the latest UX conversation alongside the current upstream proposal. It does not change the PRD, authorize implementation, or resolve architecture by implication. Confirmed user choices are retained even where the PRD proposes narrower scope; research recommendations do not override them.

Sources read: the twenty-three entries in [UX conversation memory](.memlog.md), the complete [research](../../research/competitive-caddy-homelab-manager-2026-10-16/research.md), the current [PRD](../../prds/prd-caddy-orchestrator-2026-10-16/prd.md), and the selected “Тихая мастерская” direction in the [exploration](.working/caddy-design-directions.html). U1–U23 below identify memory entries in order, including events and assumptions rather than treating every entry as a decision. PRD status is `draft`, `research_status: integrated`, and `decision_status: human-review-required`; research status is `complete`.

## Source → UX mapping

| Topic and source | Latest confirmed input or source role | UX consequence and reconciliation |
|---|---|---|
| Audience and deployment — PRD §§1–2, Q1/Q3; U6, U19 | Upstream decisions: homelabbers broadly, including Artem; manager-only Compose connects to existing Caddy. UX independently confirms connecting an installed Caddy. Existing Ansible is context, not a read-only product requirement. | Design existing-configuration adoption for a broad audience. Do not require Artem's private services or a bundled production Caddy. Ansible migration/coexistence remains open. |
| Ownership — U20; PRD FR-1, FR-4, FR-9, §4 | User explicitly wants the **entire existing Caddy configuration**, including previously created routes and global settings, after being told this extends today's boundary. | Adoption must review the whole proposed ownership boundary. **Divergence:** current PRD retains one Managed Server's route array, read-only External Routes and route-scoped snapshots. Product/architecture must reconcile ownership, apply and recovery scope. PRD A1 is already retired; the remaining narrow boundary is in the terms and FRs, not an active A1 decision. Whole-config intent does not establish fleet management. |
| Capability breadth — U16–U17; PRD A3/§4 | User's broad intent is to cover everything Caddy makes possible, using forms for **all settings** plus an editor. | Preserve this ambition explicitly. Do not silently reduce it to common-route forms or make rare capabilities editor-only. **Divergence:** PRD proposes no visual editor for every module and prioritizes narrower release capabilities. Installed-build/plugin compatibility and the exact meaning of exhaustive coverage still need definition; no blanket compatibility claim follows. |
| Form/editor relationship — U17, U22 | Home Assistant is the user's behavioral reference; forms and editor represent the same configuration. **Caddyfile is the chosen editor format**, not JSON or dual syntax. | Both contribute to one batch draft. Granularity, switching, validation and temporarily unrepresentable constructs remain unresolved. Existing Caddyfile source availability, imports, comments/format/directive preservation and form↔text fidelity require an explicit contract; runtime JSON does not establish reconstructable source. These questions must not silently replace the chosen syntax or narrow universal forms. The Home Assistant reference does not select its visual style or establish tested behavior of that product. |
| Application model — U18; PRD FR-2, FR-5–6 | Collect multiple edits → review a shared diff → apply the batch → verify results. Saving a draft must not publish it. | **Divergence/refinement:** route CRUD is an implementation baseline, not the desired immediate-publication interaction. Shared pending state must cover forms and editor. Retain conflict/unknown-outcome safeguards; batch UX must not imply a transaction across Caddy and the manager database. |
| TLS — U16; PRD Q2, FR-3/FR-8, D1 | Upstream primary workflow remains public-domain HTTPS. Latest UX decision requires public automatic certificates, private public-domain DNS-01/wildcards, local names/internal CA, and HTTPS upstreams **immediately**. | A primary public workflow can coexist with all these required scenarios. **Divergence:** PRD treats several as additional candidates or unapproved alternatives. Do not defer them silently. Product/architecture must reconcile release scope, supported builds and provider/module capabilities without reopening the settled primary direction. |
| Docker destinations — U7–U8; PRD FR-3/§4 | Both manual address/container-and-port entry and choosing discovered Docker containers/ports. | Preserve both destination paths. **Divergence:** automated discovery is a proposed PRD non-goal and a research deferral recommendation. Discovery desire does not authorize socket access, container lifecycle management or automatic publication. Discovery source/network model needs design. |
| Core jobs — U7, U11 | Add services, change routes and diagnose unavailability; expose domains, map container destinations, and serve a static website from a prepared directory on Caddy's side. | Cover proxy, redirect and file-server flows. Static publishing takes a domain plus existing directory; no upload, site builder or deployment platform was requested. PRD FR-2's Caddy-filesystem boundary aligns. |
| DNS and success — U15; PRD FR-6–8 | Operator manages DNS separately. Manager verifies that the domain was added and works for the selected route type before reporting final success. | No DNS-record management requirement. Config acceptance is an intermediate result. Show type-aware observations and incomplete/failed/unknown verification; distinguish probe vantage from the operator's client. Exact expected response and “ordinary setup” remain to clarify. |
| Visual direction — U2, U14 | Minimal, restrained, cozy engineering UI with some Apple-like glass; “Тихая мастерская” selected. | Retain calm hierarchy, soft warmth, generous spacing and readable technical identifiers. Selection approves personality, not exact navigation, palette values, density, theme default, library or every sample interaction. Glass remains restrained and subordinate to readability. |
| Devices — U9 | Desktop is primary; phone is for status and quick corrections. | Support those phone tasks without assuming full parity. Exact quick-edit boundaries and cross-device draft behavior are still open. |
| Exploratory mock — U12/U14; selected variant | Four anonymized external/read-only examples, “Пример · не live”, no known API or service-health observations. | These examples are historical exploration, **not a restriction on the now-confirmed whole-configuration management**. The chosen list/detail arrangement is not a signed-off IA. No live health, certificate readiness or API success may be inferred. |
| Research and safety — research recommendations/limitations; PRD FR-5–10 | Research supports explicit ownership, honest observations, diagnosis, persistence and recovery; its breadth cautions are recommendations. Competitor claims and selected anecdotes are not usability or runtime validation. | Use these as rationale and design constraints, not new user requirements or proof the existing product works. Preserve safeguards while reconciling broader scope; existing route snapshots must not be relabeled as complete-config backups. |
| Discovery closure — U23 | User has no further scenarios for now; future ideas can be separate features. Focus on the substantial scope already discussed. | Proceed to a concrete UX draft and lightweight sign-off. Do not broaden features or continue generic elicitation; remaining unknowns are explicit handoff items. |

## Missing inputs and handoff decisions

Already answered: audience, existing-Caddy deployment, primary public HTTPS, all listed TLS scenarios, discovery plus manual input, full-configuration intent, batch application, forms plus Caddyfile editor, visual personality and primary devices. Do not ask these again as if no choice exists. U23 closes additional discovery; the details below are recorded gaps for downstream work, not a new questionnaire or a reason to delay the concrete draft.

Remaining user-facing detail:

- Ansible transition/coexistence and who may write configuration after adoption; instance/server-count expectations beyond the confirmed whole-configuration boundary.
- Editor switching behavior, detailed form coverage for build-specific modules, Caddyfile source/import ownership and preservation, and acceptable behavior when a construct cannot yet be represented in forms. Preserve the selected Caddyfile syntax and all-settings request while resolving these details; source fidelity is not yet established.
- Which corrections must be possible on a phone; whether drafts must continue across devices/sessions.
- What proves success for proxy, redirect and file-server routes; intended client/probe vantage, protected-service expectations, timing and meaning of “ordinary setup”.

Product/architecture follow-up must reconcile the explicit PRD divergences before implementation scope is treated as agreed: whole-config adoption/ownership/recovery, exhaustive forms and editor, all requested TLS scenarios, Docker discovery and shared batch application. Draft durability, remote drift, ambiguous/partial verification and installed-module boundaries require concrete designs. None of these open implementation questions cancels the confirmed capability requests, and this reconciliation does not automatically edit upstream documents.
