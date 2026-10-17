---
name: Caddy Orchestrator
description: Existing-Caddy adoption, shared drafts, reviewed batch apply and type-aware verification.
status: final
phase: handoff
decision_status: confirmed-direction-with-explicit-proposals
validation: basic-checks-only
updated: 2026-10-17
language: en
visual-reference: DESIGN.md
sources:
  - .memlog.md
  - ../../../../CLAUDE.md
  - ../../research/competitive-caddy-homelab-manager-2026-10-16/research.md
  - ../../prds/prd-caddy-orchestrator-2026-10-16/prd.md
  - ../../prds/prd-caddy-orchestrator-2026-10-16/reconcile-code-baseline.md
---

# Caddy Orchestrator — Experience

## Foundation

Responsive web for homelabbers broadly, including Artem. Desktop is primary for configuration; phone supports status and quick corrections. Production connects to already running Caddy; bundled Caddy is a test fixture. DESIGN.md defines Quiet Workshop. No new library, frontend migration or native app is selected.

Confirmed choices in [.memlog.md](.memlog.md): manage the entire existing configuration, including previous Routes/global settings; forms for all settings plus a **Caddyfile editor**; multiple edits in one shared draft → diff → apply → type-aware verification. DNS service records stay operator-managed. Manual Upstream input and Docker discovery are both desired. Static website means domain plus prepared Caddy-side directory. All four scenarios are wanted immediately: public automatic HTTPS, private/public-domain DNS-01 and wildcards, local names/internal CA, and HTTPS Upstreams. Broad Caddy capability is the intent, not an established promise of arbitrary plugin/build compatibility.

Discovery is closed for this feature. Remaining technical unknowns are downstream handoff items, not a new questionnaire. All concrete grouping/defaults/behaviors below are **[ASSUMPTION — proposed interaction contract]** unless explicitly identified as a user choice. Walkthroughs are designed from supplied tasks/anonymized infrastructure context, not observed sessions. Product implementation/deployment remains separate from this approved planning correction.

Final status records completion of this session's UX capture. Marked assumptions remain proposals, and the corrected PRD records approved scope; visual assumptions retain their labelled status.

### Source reconciliation and compatibility

The [PRD](../../prds/prd-caddy-orchestrator-2026-10-16/prd.md) now records the confirmed expanded scope; Q1–Q3 settle audience, primary public HTTPS and external Caddy. [Research](../../research/competitive-caddy-homelab-manager-2026-10-16/research.md) is complete and advisory. Latest explicit choices supersede narrower draft proposals:

| Topic | Existing baseline / PRD proposal | User choice to carry forward |
|---|---|---|
| Ownership | FR-1/FR-4: one Managed Server's Routes; External Routes read-only. | Full existing configuration, including previous Routes/global settings; explicit whole-config adoption/recovery required. |
| Apply | Legacy post-setup CRUD applies immediately. | Shared batch draft, diff, explicit apply, then verification; draft save must not publish. |
| TLS | FR-8 leaves alternatives conditional. | All four listed scenarios desired now; supported builds/modules/credentials need design. |
| Coverage | Proposed exclusions include discovery/forms for every module. | Docker discovery plus manual input; forms for all settings plus Caddyfile editor. Preserve these goals and expose compatibility limitations. |

Route, Upstream, applied, failed, unknown and service verified retain PRD meanings. Managed Server/External Route/Snapshot identify current implementation concepts when used as such. Proposed configuration draft holds pending edits against a known baseline; proposed configuration recovery point is distinct from a route-array Snapshot and excludes certificate keys/CA data unless explicitly covered.

Caddy documents [Caddyfile](https://caddyserver.com/docs/caddyfile) as less expressive than native JSON; its [API](https://caddyserver.com/docs/api) reads runtime JSON, and [adapters](https://caddyserver.com/docs/config-adapters) output JSON. **Inference:** arbitrary runtime configuration cannot be assumed to yield an original or losslessly equivalent Caddyfile. User intent remains Caddyfile plus forms; the architecture companion defines source acquisition, compatibility and preservation. Do not promise automatic reconstruction, comments/import recovery or universal round-trip fidelity.

## Approved technical interaction contract — 2026-10-17

The course correction applies the architecture's durable source/revision/recovery rules to F1–F7. Source reconciliation below records historical baseline contrasts, not unresolved scope choices. For current API and compatibility behavior use the [configuration contract](../../architecture/architecture-caddy-orchestrator-2026-10-17/configuration-contracts.md). Visual assumptions remain proposals; scope and approved safety contracts are settled. Private unsaved buffers are permitted alongside shared persisted state.

## Information Architecture

**[ASSUMPTION]** Routes, Configuration and History are top-level destinations; a global Changes entry exposes the shared draft. Connection setup appears when needed. Details are pages/sections, not nested modal workflows. No login/users/fleet/metrics/general log viewer or separate L4 surface is added.

| Proposed surface | Need / entry | Flow | Visual coverage |
|---|---|---|---|
| Connection and adoption | Existing Caddy endpoint, source/build capability, current config and ownership review. | F1 | Spine-only |
| Routes and Route detail | Search name/domain/path/destination; inspect settings/results; add/edit/disable/delete proxy, redirect, file server. | F2–F5 | [Representative mock](mockups/key-routes.html) |
| Route settings | Hosts/path/action, manual/discovered Upstream, prepared root, TLS and advanced fields; form/editor views. | F2–F4 | [Form/Caddyfile mock](mockups/key-route-editor.html); remaining types/states spine-only |
| Configuration | Existing global/server/settings constructs, source/capability context and Caddyfile peer view. | F1, F4 | Spine-only |
| Changes and result | Shared diff, validation, apply progress/result and links to affected checks/conflicts. | F2–F6 | [Diff/result mock](mockups/key-changes-result.html) |
| Verification and diagnosis | Detail/result drill-in with typed expectations, evidence, retry and next checks. | F2–F5 | Mixed checks shown in [result mock](mockups/key-changes-result.html); full drill-in spine-only |
| History and recovery | Compare/export/review restore by actual scope; backup/restart/upgrade guidance. | F6 | Spine-only |
| Preferences / connection details | Appearance and connection context within Configuration; changing target requires renewed adoption. | F1, F7 | Spine-only |

[Quiet Workshop exploration](.working/caddy-design-directions.html) illustrates the selected visual direction; its read-only labels are not adopted behavior. The three newer key-screen references illustrate the proposed core layout. Connection/adoption, global configuration, full diagnosis, history and preferences remain spine-only; no separate mobile mock is supplied. Mock-coverage choices were offered; absent a selection, this is an assistant default, not user approval of every layout. Detailed interaction proposals remain explicitly marked assumptions.

## Voice and Tone

Calm, specific: name operation, known effect and next action. Language/default locale is undecided; Russian examples follow the conversation, not a localization commitment. Use «Черновик сохранён. Caddy ещё не изменён» only after confirmed draft persistence. After acceptance: «Caddy принял конфигурацию. Проверяем сервис». For bounded success: «Проверка прокси прошла из сети менеджера в 14:32». Unknown: «Ответ Caddy потерян. Сначала сверим текущую конфигурацию». DNS: «Запись DNS настраивается отдельно».

Avoid generic “saved successfully” for publication and universal “works” after one probe. Ordinary summaries/errors redact secrets; sensitive previews/exports need an explicit handling policy.

## Component Patterns

Names match DESIGN.md. Behaviors are proposals, not current implementation claims.

| Component | Behavioral contract |
|---|---|
| Workspace header | Current connection/adoption and pending-change entry. Navigation never applies. Target stays visible during diff/apply. |
| Route row | Open detail; separate enabled, live/pending values and last known result. Toggle/delete create draft edits, with no hover-only action. |
| Search and filters | Filter without applying; retain query on return; distinguish no Routes from no matches. |
| Settings form | Shared draft, field-linked errors, input retention; advanced sections remain form-accessible. Manual destination remains available when discovery fails. |
| Caddyfile editor | Same candidate/scope as forms. Validate before switching; preserve invalid text and unsupported directives. Block destructive conversion instead of dropping content. |
| Change bar | Separate unsaved input, saved draft and validation. Review opens all changes. No publication on blur; saved drafts are durable; unsaved input is distinct and protected as specified in the architecture contract. |
| Diff review | Baseline, target, affected scope and redacted changes. Live drift requires reconciliation and renewed review; disabled Apply states its reason. |
| Status evidence | Result, expectation, time, vantage/method; distinct configuration/TLS/Upstream/published-service observations. Retry checks without reapplying. |
| Action button | Explicit save-draft/review/apply/retry/restore verbs; prevent duplicates in flight. Unknown apply cannot become blind retry. |
| Recovery entry | Actual scope/exclusions/time/target; restore prepares a reviewed candidate against live state. Existing Snapshot remains route-only. |
| Feedback panel | Persistent operation/progress/error/known effect/next step; distinguishes offline, denied and unsupported. Applied revisions, attempt state and observations persist with explicit staleness across restart. |

## State Patterns

**[ASSUMPTION]** Baseline, candidate and observations remain separate. For each failure, state the remote effect: not sent, rejected, accepted with a later error, or unknown. Failure alone never means unchanged Caddy.

| Surface | Applicable states and treatment |
|---|---|
| Connection/adoption | Loading/probing, unavailable/denied, connected-unadopted, stale preview, source missing/incompatible, unknown build, reviewed. Preserve endpoint; test neither adopts nor applies. |
| Routes/detail | Loading, empty, no matches, selected/focused, stale/removed detail, manager offline/Caddy unavailable. Last known data has timestamps; pending overlay does not replace live facts. |
| Route settings/Configuration | Clean, edited, saved draft, invalid, parse error, unsupported construct, denied capability, live drift. Preserve directives/imports/comments where source is available; never silently narrow to form fields. |
| Changes/result | Empty, reviewable, validation blocked, baseline drift, applying, accepted/checking, rejected, uncertain, partly verified. Per-item verification outcomes; no remote/local transaction or automatic rollback promise. |
| Verification/diagnosis | Not checked, checking, passed, failed expectation, timeout, unavailable vantage, unknown/stale. Sanitized evidence and next action; unknown checks preserve known acceptance. |
| History/recovery | Loading/empty/error, incompatible/missing scope, target mismatch, stale baseline, prepared/applied/unknown restore. Scope-sensitive export; up to 100 history entries, with active recovery pins and explicit capacity blocking. |
| Preferences/connection | Loading/editable/invalid/unsaved, target-change review, offline/denied. Appearance does not change Caddy; an old draft is not silently rebound. |

All surfaces retain visible keyboard focus using `{colors.accent}` / mode equivalent and `{colors.control-line}` boundaries. Async feedback does not steal focus. Offline mode cannot promise replay, draft persistence or cross-device sync without support. Access denial needs no invented account system.

## Interaction Primitives

Forms and Caddyfile are peers, as requested with Home Assistant as interaction reference. **[ASSUMPTION]** One candidate and baseline feed both views. Unavailable existing source, missing imports, unsupported directives/modules or nonrepresentable JSON trigger an explicit compatibility/preservation state. Preserve recoverable source and raw runtime data separately; unsafe conversion/apply blocks with an explanation. Temporary form gaps must remain visible while the goal stays forms for all supported settings.

Confirmed sequence: **multiple edits → shared diff → apply → type-aware verification**. Proposed detail: Save draft records only candidate; Review compares baseline; Apply rechecks live state, validates, prepares an appropriately scoped recovery point and submits. Architecture specifies durable content reconciliation and guarded root safety; Caddy's [API](https://caddyserver.com/docs/api) does not make multiple requests one transaction. External drift preserves candidate and requires manual review, never force/auto-merge. External writes are guarded at submission and detected as later drift; no exclusive control over Ansible or other writers is implied.

Mouse/touch/keyboard can perform primary actions. Escape closes transient UI without discarding; leaving unsaved input offers save/discard/cancel. No required drag/shortcut. Reduced motion removes transitions. Secrets stay out of ordinary copy/error summaries.

### Verification contract

DNS service records remain external. DNS-01 is certificate challenge automation, not service-record provisioning. **[ASSUMPTION]** Checks start after acceptance, report per-item results and have explicit retry/time limits; unavailable checks yield applied/unverified, not fabricated success.

| Type/scenario | Proposed observation / boundary |
|---|---|
| Reverse proxy | Request published host against a defined expected response; Upstream evidence separate. Any HTTP response may still be a login/error page. |
| Redirect | Inspect first-hop code and Location; following destination is a separate observation. |
| Static file server | Request known resource from prepared root; path alone proves neither mount nor permissions. No upload/deployment. |
| Public automatic HTTPS | Observed hostname/certificate/response, pending issuance or unknown; config acceptance proves neither issuance nor renewal/restart survival. |
| Private/public-domain DNS-01 and wildcard | Show build/module/credential prerequisites and available challenge/certificate evidence; issuance does not create DNS records. |
| Local names/internal CA | Separate service response and certificate observation; manager trust cannot prove client trust. |
| HTTPS Upstream | Separate upstream TLS/hostname/trust settings from client-facing TLS; no silent verification bypass. |

Every result names time and vantage: manager network, Caddy observation or operator client. Expected status/path/Location/trust policy is explicit per type; the approved contract specifies deadlines, retries, staleness and unavailable observations. Global changes without endpoints get relevant configuration checks with narrower scope stated. User receives final success only for the selected type's verified expectations; partial/unknown outcomes remain explicit.

## Accessibility Floor

**[ASSUMPTION]** Target WCAG 2.2 AA; contrast lives in DESIGN.md. Native semantics, labels, error summary/field links, visible focus and return-focus restoration. Preserve input on errors. Diff/editor require keyboard and screen-reader operation, textual diff alternative and escape from editor focus. Announce operation transitions politely, not every poll. Touch targets ≥44px; sticky content must not obscure fields at zoom/soft keyboard. Reflow at 320 CSS px, 200% text zoom, reduced motion/transparency. State/diff/selection never rely on color alone.

## Key Flows

Designed Artem walkthroughs below are not observed sessions; examples remain anonymous. Each follows the proposed surface map.

### F1 — Connect with explicit ownership (FR-1; expanded ownership)

1. On desktop Artem enters existing Caddy endpoint, tests connection, reviews current source/runtime configuration and build/capabilities.
2. He reviews whole-config ownership, previous Routes/global settings, preservation gaps and persistence/external-writer prerequisites.
3. **Climax:** explicit adoption exposes existing settings without implicit rewrite/publication.

Failure: denied/unreachable/stale preview preserves input; missing source/incompatible config shows preservation limits rather than claiming complete adoption. Ansible migration is not assumed.

### F2 — Publish and retain a service (FR-2, FR-3, FR-6, FR-8, FR-11)

1. Artem adds proxy domain and manual address/port or discovered container/port; DNS is prepared externally.
2. He chooses applicable HTTPS settings, adds to shared draft, reviews diff and applies.
3. **Climax:** acceptance and subsequent typed proxy/HTTPS result appear separately with expectation/time/vantage.

Failure: invalid input remains; discovery failure leaves manual input. DNS/certificate/response failures offer specific diagnosis; restart survival still needs runtime evidence.

### F3 — Static site or redirect (FR-2, FR-3, FR-7)

1. Artem supplies domain plus prepared Caddy-side directory, or redirect destination/code.
2. He reviews/applies and checks a known static resource or first-hop redirect expectation.
3. **Climax:** the result validates the chosen type rather than presenting generic proxy health.

Failure: missing root/permissions or unexpected Location leads to the relevant next check. No file upload, site builder or DNS management is implied.

### F4 — Change through forms and Caddyfile (FR-3, FR-4 reconciliation, FR-8)

1. Artem opens previous Route/global settings, edits forms and switches to the Caddyfile peer view.
2. Supported values and recoverable source remain intact; TLS/Upstream/global edits accumulate in one candidate.
3. **Climax:** shared diff exposes intended changes/preserved content before batch apply and relevant verification.

Failure: invalid text remains recoverable; missing modules/source or unsafe round-trip blocks the affected conversion/apply. Whole-config fidelity is a handoff contract, not established current behavior.

### F5 — Diagnose unavailability (FR-5–FR-8, FR-11)

1. Artem opens a Route on desktop/phone and reads enabled, acceptance and timestamped observations separately.
2. He inspects failing/unknown layer, retries its check or corrects settings through draft/review/apply.
3. **Climax:** evidence identifies the next action without equating online Caddy with a working domain.

Failure: unavailable vantage stays unknown with client instructions; lost apply acknowledgement requires live reconciliation before another write. Diagnosis retry never reapplies.

### F6 — Recover and maintain predictably (FR-5, FR-9, FR-10; SM-2)

1. Artem opens History, reads recovery scope/exclusions and compares the chosen point with live config.
2. He prepares/reviews restore candidate, applies and examines separate service checks.
3. **Climax:** he knows what was restored/verified and which data lies outside it.

Failure: stale baseline/target mismatch blocks overwrite; unknown acknowledgement stays uncertain. Guidance covers restart/upgrade/full backup prerequisites; existing route Snapshot is not full-config recovery or certificate/CA backup. No automatic downgrade or full-system restore is promised.

### F7 — Phone status and quick correction (FR-6, FR-11)

1. Artem finds a service on phone, inspects known status and corrects a small field in stacked forms.
2. He follows the same draft/review/apply sequence; appearance preferences remain separate.
3. **Climax:** the correction's actual result is readable without desktop layout or bypassed review.

Failure: offline/validation/soft-keyboard states retain context without hidden actions. Cross-device continuity uses expected draft revisions and SSE resync; unsaved phone input remains separate. Exact quick-edit bounds do not imply full mobile parity.

## Inspiration & Anti-patterns

Home Assistant is the user's forms-plus-editor behavioral reference, not verified implementation evidence or visual imitation. Quiet Workshop supplies the chosen composition; Apple glass remains restrained. NPM provides comparable purpose, not parity. Avoid immediate publication on edits, universal green status, dropped unknown directives, blanket read-only previous config and route-only history labeled full backup.

## Responsive & Platform

**[ASSUMPTION]** Desktop retains spacious list/detail and comfortable form columns; diff/editor can expand. Below 620px stack values/fields and wrap Change bar; use full-width search and readable unified diff. Meaning/review remain consistent across devices. Exact quick edits, theme/language and editor ergonomics need review; native/mobile parity is not assumed.

## Downstream Handoff

Discovery stays closed. [Approved course correction](../../sprint-change-proposal-2026-10-17.md) and [architecture contracts](../../architecture/architecture-caddy-orchestrator-2026-10-17/configuration-contracts.md) settle source bundles/import provenance, nested preservation, build/form inventory, durable review/apply/restore, secret versions, revisioned sessions/probes, external persistence, migration and complete backup. The earlier handoff unknowns in this historical UX capture are resolved by those contracts; they are not new generic elicitation gates.

Saved drafts/invalid editor buffers are protected and persistent; unsaved input has a separate buffer. Stale save returns 409 with preserved input. Reconnect fetches current state without overwriting that buffer. Review binds source/dependencies/build/preservation/target/generation/revision/root ETag; drift or any bound change requires renewed review. Interrupted attempts lock until content reconciliation names candidate present, baseline present, conflict or unavailable; never offer blind retry. Restore checks captured live review ETag, target and actual scope. Changing target detaches and retains work rather than silently clearing it.

Missing source, unsafe mixed arrays or unknown codecs show a precise compatibility block with protected originals. Every declared supported setting must have a dedicated or structured directive form; temporary gaps cannot be labelled completed editor-only scope. Browser policy/opaque observations remain unavailable/unknown and late checks remain stale for newer applied revisions. PKI root download is explicit and describes default-CA provisioning. Complete backup requires operator-side external asset acquisition; UI history remains narrower.

[Epics](../../epics.md) and [QA](../../../test-artifacts/test-design-qa.md) map all user outcomes to acceptance. DESIGN keeps Quiet Workshop and labelled token/layout assumptions. Three mocks remain illustrative static references, not live verification or exact UI completion. Earlier basic rendering checks are dated evidence; this correction performed no new browser/product runtime validation. No new fleet, Docker lifecycle, arbitrary plugin universe or automatic remote/local transaction guarantee is added.
