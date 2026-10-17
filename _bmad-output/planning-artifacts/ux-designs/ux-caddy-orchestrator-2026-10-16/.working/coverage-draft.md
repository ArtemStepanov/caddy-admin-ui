# Review-draft coverage — 2026-10-17

Proactive Pass 1 from bmad-ux/references/validate.md, not opt-in independent validation or final sign-off. Read CLAUDE.md first, current completed research, current draft PRD, canonical memlog, selected Quiet Workshop variant, all configured design/experience examples and the spec/validation references. Discovery remains closed. Documents total 334 lines before further parent edits.

## 1. Flow coverage

The PRD has no named UJs. Its candidate requirements map to named Artem walkthroughs with numbered steps, explicit climax and failure path:

| Requirement / input | Flow |
|---|---|
| FR-1 — Connect with explicit ownership | F1 |
| FR-2 — Publish and retain a service | F2, F3; F6 carries restart-evidence guidance |
| FR-3 — Validate supported destinations | F2–F4 |
| FR-4 — Preserve ownership boundaries | F1, F4; explicit user supersession/reconciliation, not silent retention |
| FR-5 — Handle concurrent and uncertain changes | F4–F6 |
| FR-6 — Show what is known | F2, F5, F7 |
| FR-7 — Diagnose a failed publication | F3, F5 |
| FR-8 — Public HTTPS with supported mode/settings choice | F2, F4 plus all four confirmed scenario rows in Verification contract |
| FR-9 — Recover Routes without overstating coverage | F6; current route Snapshot distinct from proposed full-config recovery |
| FR-10 — Back up and upgrade predictably | F6 guidance, no invented backup UI |
| FR-11 — Make routine work clear and accessible | F2, F5, F7 plus accessibility/component rules |
| Docker manual + discovery, static prepared folder, Caddyfile/forms, full config, batch, phone | F2, F3, F4, F1/F4, F2–F6, F7 respectively |
| SM-1 / SM-2 / SM-C1 | F2/F5/F6 task chain; stated evidence boundaries; no speed improvement hides scope/result |

**Real gaps:** broad goals are covered, but exhaustive supported module/form field coverage and concrete automated verification expectations cannot be asserted yet. All TLS scenarios have contract rows; they do not yet have rendered, executable prerequisite/credential/error paths. These are downstream compatibility/design work, not reasons to restart generic elicitation.

## 2. Token completeness

Python/YAML mechanical check passed: every token reference in both documents resolves, every color is six/eight-digit hex, frontmatter parses, all source paths exist. Contrast targets are stated for load-bearing pairs. Exact extensions are labeled assumptions.

**Real gap:** rendered contrast/focus/glass/zoom verification has not been performed. Semantic colors/control outlines are provisional additions to the chosen mock, not separately user-approved.

## 3. Component coverage

All 11 component names in DESIGN frontmatter have substantive rows in both DESIGN.Components and EXPERIENCE.Component Patterns: Workspace header, Route row, Search and filters, Settings form, Caddyfile editor, Change bar, Diff review, Status evidence, Action button, Recovery entry, Feedback panel. No mismatch found.

## 4. State coverage

All eight proposed IA surfaces have state rows, including cold/empty where applicable, validation/error, offline/unavailable, denied, stale/conflict and pending/unknown outcomes; global focus/announcement rules apply everywhere. Applied versus verified, timestamps and vantage are explicit.

**Real gap:** durable draft/result storage, source-import incompatibility recovery, unknown remote outcome reconciliation and whole-config restore cannot be specified as working guarantees before architecture determines the contracts.

## 5. Visual references

The only artifact is the exploratory .working/caddy-design-directions.html; both spines link it inline and identify its limited meaning. Selected variant only informs extraction; other alternatives are not adopted. No promoted mockups/wireframes/imports exist. Spines-win statement appears once in DESIGN.

**Real gap:** no key-screen mock for form/Caddyfile switching, shared diff/apply result, full-config adoption, TLS/discovery/verification, history or phone quick correction. Surface table records exploration-only versus spine-only. Final mock coverage choice and any additional rendering remain parent-led workflow steps.

## Handoff status

Both spines intentionally remain draft/review-draft. The current PRD still conflicts with confirmed whole-config management, full TLS set, discovery and all-forms intent; conflicts are recorded without changing upstream sources. Caddyfile source/runtime compatibility is cited and treated as inference plus architecture handoff. No new features or technical permissions are inferred. This file does not substitute for the optional Reviewer Gate.
