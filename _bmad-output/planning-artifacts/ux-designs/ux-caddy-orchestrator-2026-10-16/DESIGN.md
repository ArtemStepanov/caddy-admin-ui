---
name: Caddy Orchestrator — Quiet Workshop
description: Calm engineering surfaces for understanding and changing a homelab Caddy configuration.
status: final
phase: handoff
decision_status: confirmed-direction-with-explicit-proposals
validation: basic-checks-only
updated: 2026-10-17
language: en
sources:
  - .memlog.md
  - ../../../../CLAUDE.md
  - ../../research/competitive-caddy-homelab-manager-2026-10-16/research.md
  - ../../prds/prd-caddy-orchestrator-2026-10-16/prd.md
  - .working/caddy-design-directions.html
colors:
  canvas: '#F5F3EF'
  panel: '#FFFEFA'
  ink: '#252C28'
  muted: '#626A62'
  line: '#DEDFD6'
  soft: '#EDECE4'
  accent: '#435E47'
  selected: '#E8EEE2'
  glass: '#F7F8F3DF'
  shadow: '#292E2612'
  canvas-dark: '#202220'
  panel-dark: '#282B28'
  ink-dark: '#E8EBE5'
  muted-dark: '#B1B9AE'
  line-dark: '#41473F'
  soft-dark: '#343B32'
  accent-dark: '#B4C9A7'
  selected-dark: '#364331'
  glass-dark: '#303630E8'
  shadow-dark: '#00000024'
  control-line: '#777F73'
  control-line-dark: '#919B89'
  success: '#355E40'
  success-dark: '#AED0AF'
  warning: '#805318'
  warning-dark: '#E7C184'
  danger: '#A13630'
  danger-dark: '#F0AAA3'
typography:
  title: {fontFamily: '-apple-system, BlinkMacSystemFont, Segoe UI, sans-serif', fontSize: 26px, fontWeight: '500', lineHeight: '1.2', letterSpacing: -0.6px}
  section: {fontFamily: '-apple-system, BlinkMacSystemFont, Segoe UI, sans-serif', fontSize: 16px, fontWeight: '500', lineHeight: '1.5'}
  body: {fontFamily: '-apple-system, BlinkMacSystemFont, Segoe UI, sans-serif', fontSize: 14px, fontWeight: '400', lineHeight: '1.5'}
  meta: {fontFamily: '-apple-system, BlinkMacSystemFont, Segoe UI, sans-serif', fontSize: 12px, fontWeight: '400', lineHeight: '1.5'}
  code: {fontFamily: 'ui-monospace, SFMono-Regular, Consolas, monospace', fontSize: 12px, fontWeight: '400', lineHeight: '1.5'}
rounded: {sm: 8px, md: 10px, lg: 16px}
spacing:
  '1': 4px
  '2': 8px
  '3': 12px
  '4': 16px
  '5': 24px
  '6': 30px
  row-block: 21px
  row-inline: 20px
  mobile-margin: 16px
components:
  Workspace header: {background: '{colors.panel}', padding: '{spacing.6}'}
  Route row: {background: '{colors.panel}', selected: '{colors.selected}', padding-block: '{spacing.row-block}'}
  Search and filters: {background: '{colors.panel}', radius: '{rounded.md}'}
  Settings form: {background: '{colors.panel}', radius: '{rounded.lg}'}
  Caddyfile editor: {background: '{colors.panel}', font: '{typography.code.fontFamily}'}
  Change bar: {background: '{colors.glass}', radius: '{rounded.md}'}
  Diff review: {background: '{colors.panel}', radius: '{rounded.lg}'}
  Status evidence: {foreground: '{colors.muted}'}
  Action button: {background: '{colors.accent}', foreground: '{colors.panel}', radius: '{rounded.md}'}
  Recovery entry: {background: '{colors.panel}', radius: '{rounded.md}'}
  Feedback panel: {background: '{colors.soft}', radius: '{rounded.md}'}
---

## Brand & Style

Artem selected **Quiet Workshop / Тихая мастерская**: minimal, cozy engineering software, with restrained Apple-like glass. Addresses, configuration scope and consequences remain easy to read during ordinary infrastructure work.

The overall direction is confirmed. Exact tokens and component extensions below are **[ASSUMPTION — review proposals]**, extracted from the chosen mock or supplied to complete this draft. No component library or frontend migration is selected; Preact/Tailwind remains the implementation baseline.

Final status records completion of this session's UX capture. Explicit assumptions remain proposals; it does not convert the remaining visual assumptions into measured usability/accessibility evidence. Approved scope and safety contracts are recorded in the PRD and course correction.

[Quiet Workshop exploration](.working/caddy-design-directions.html) illustrates the spacious list, warm surfaces, system typography and detail block. Only its first variant is selected. Sample services and external-read-only labels are illustrative; desired ownership is defined in EXPERIENCE.md. DESIGN.md and EXPERIENCE.md take precedence over mocks on conflict.

Key-screen references apply this direction to the [Routes list](mockups/key-routes.html), [form and Caddyfile views](mockups/key-route-editor.html), and [batch diff with result](mockups/key-changes-result.html). They are static examples with anonymized data, not a working product or proof of source fidelity.

## Colors

Warm `{colors.canvas}` surrounds near-white `{colors.panel}`. `{colors.ink}` carries names and decisions; `{colors.muted}` supports context. `{colors.line}` divides content; `{colors.control-line}` defines interactive boundaries where a soft divider is insufficient. Green `{colors.accent}` identifies primary actions; `{colors.selected}` indicates selection, never health. The `-dark` tokens are corresponding appearances extracted from the mock. **[ASSUMPTION]** Follow system appearance with light/dark/system preferences; no default theme is confirmed.

**[ASSUMPTION — provisional semantic extensions]** `{colors.success}`, `{colors.warning}`, `{colors.danger}` and their dark partners support verified results, pending/uncertain conditions and errors. Apply as text/icons on the mode's panel, with written state and evidence. Unknown is neutral. Diff additions/removals also need signs and labels. Applied configuration can remain unverified.

Contrast targets: normal text ≥4.5:1, large text ≥3:1, essential control/focus boundaries ≥3:1. Soft dividers cannot be the sole input boundary. Calculated panel pairs pass in both modes: muted text 5.54:1 / 7.11:1, primary action 7.10:1 / 8.08:1, control boundaries 4.11:1 / 4.95:1. Muted text on glass composited over canvas is 5.21:1 / 6.30:1. Verify actual rendered focus, overlays and zoom; these calculations are not a completed accessibility audit.

## Typography

System sans follows the mock: `{typography.title.fontSize}` page titles, `{typography.section.fontSize}` section headings, `{typography.body.fontSize}` work text. `{typography.meta.fontSize}` serves supporting context, never the sole error or apply consequence. Medium weight identifies service names; no oversized dashboard display type.

Use `{typography.code.fontFamily}` for domain, path, Upstream, Caddyfile and diff values. Preserve selectable text and Cyrillic/IDN rendering; wrap addresses in summaries. **[ASSUMPTION]** Editor text grows to 14px for sustained work; phone inputs use at least 16px. Browser zoom must not truncate controls.

## Layout & Spacing

Quiet Workshop uses a top header, title/search area, spacious Route list and selected detail below. Proposed navigation destinations are defined in EXPERIENCE.md, not inferred from unselected variants. Desktop margin is `{spacing.6}`, row padding `{spacing.row-block}` by `{spacing.row-inline}`, major gaps `{spacing.5}`.

**[ASSUMPTION]** Forms use comfortable reading columns; diff/editor can expand. Below 620px, use `{spacing.mobile-margin}` margins, full-width search and stacked name/domain/destination/state. The Change bar wraps without hiding content or focused fields. Breakpoints remain provisional until key-screen rendering at realistic content/zoom.

## Elevation & Depth

Use spacing, tone and dividers for hierarchy. **[ASSUMPTION]** Reserve `{colors.shadow}` / `{colors.shadow-dark}` for overlaid feedback or a floating Change bar. Optional glass is limited to that bar or navigation chrome via `{colors.glass}`. Configuration, fields and error text stay opaque; reduced-transparency/unsupported environments use the panel fallback.

## Shapes

Use `{rounded.lg}` for list/detail/form containers, `{rounded.md}` for fields/buttons and `{rounded.sm}` for small service markers. Favor soft rectangles over decorative pills. Icons are restrained and labeled for consequential actions; no icon library is selected.

## Components

Names match EXPERIENCE.md. **[ASSUMPTION]** New components extend the selected list/detail vocabulary. The three key-screen references show the core composition; remaining surfaces use these visual contracts and EXPERIENCE.md tables.

| Component | Visual contract |
|---|---|
| Workspace header | Name, connection context and pending count in one quiet band; thin divider, body typography. No fleet switcher implied. |
| Route row | Marker/name, domains, destination/action, enabled and application/verification text; panel, line divider, selected tint. Wrap long values; actions visible without hover. |
| Search and filters | Explicit label, panel field, control-line boundary, clear focus outline; compact group without collapsing the page on no matches. |
| Settings form | Scope, section headings, labels, fields, inline errors and preservation notices; panel sections with generous gaps. Advanced settings remain part of the form. |
| Caddyfile editor | Monospace opaque surface, scope/version heading and peer form/editor switch; text indicators accompany syntax/validation color. |
| Change bar | Draft state/count and review action; wraps on phone. Optional restrained glass with opaque fallback; never covers final actions. |
| Diff review | Affected-area summary, labeled before/after or unified text, semantic ink plus signs; destructive scope visible before Apply. |
| Status evidence | Separate configuration, TLS, Upstream and published-service results with time/vantage; marker supplements written state, never one universal green dot. |
| Action button | Accent fill/panel text; dark uses accent-dark/panel-dark. Secondary outline, labeled danger action; visible focus distinct from selection. Proposed 40px desktop/44px touch height. |
| Recovery entry | Readable timestamp, scope and reason; detail reveals contents/restore. Route-array versus full-configuration scope remains explicit. |
| Feedback panel | Persistent operation/known effect/next action; soft fill and semantic heading. Loading is restrained; error text is not a saturated full-panel wash. |

## Do's and Don'ts

| Do | Don't |
|---|---|
| Preserve the warm spacious list as settings grow. | Substitute the unselected compact table or glass split view. |
| Distinguish pending, applied and observed states. | Let enabled or online imply a working service. |
| Keep technical values legible, selectable and copyable. | Truncate the only destination or display secrets in summaries. |
| Keep glass restrained and legible. | Blur configuration, evidence or errors. |
| Pair semantic color with text/signs. | Use selection as proof of health or color alone for diff. |
| Treat forms and Caddyfile as peer views. | Claim complete module/form coverage before compatibility is established. |

## Course-correction linkage — 2026-10-17

Quiet Workshop and existing component/token proposals are retained. EXPERIENCE now binds source compatibility, shared revision checks, guarded review/apply/restore and bounded observations to the approved architecture. [PRD](../../prds/prd-caddy-orchestrator-2026-10-16/prd.md) carries confirmed full-settings scope; the [33-story plan](../../epics.md) supplies concrete form/global/server acceptance. This update does not claim a new rendered accessibility audit or implementation.
