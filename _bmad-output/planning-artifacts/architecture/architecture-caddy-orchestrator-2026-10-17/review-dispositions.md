# Architecture review dispositions — course-correction revision 1

The [approved proposal](../../sprint-change-proposal-2026-10-17.md) retains AD-1–AD-11 and expanded scope. Earlier reviews are dated provenance; their old positive/critical verdicts do not certify the corrected revision. Fresh reviews are [rubric](reviews/review-course-correction-rubric.md), [technology](reviews/review-course-correction-tech-fit.md) and [adversarial seams](reviews/review-course-correction-adversary.md). Follow-up closure checks found no remaining planning blocker. Product compatibility/acceptance remains planned execution.

| Historical review topic | Disposition for current contract |
| --- | --- |
| Adversary: no persistent route identity | Original absence premise superseded by comment IDs; AD-1 and X1/X3 now also define distinct site/node IDs, unchanged parsing, explicit save, duplicate repair and immutable enabled/node maps in history. |
| Adversary: compilation/secret snapshot catch-22 | Earlier plaintext-JSON premise superseded, but lifecycle concern was valid. AD-5/P05 encrypts exact source/runtime/legacy envelopes, protects invalid buffers, versions credentials and gates safe projection/hydration. |
| Adversary: no usable ETag / unstable manager hash | No-ETag premise rejected by pinned Caddy API/source. AD-8 uses exact server ETag for remote root precondition; P04 separately defines structural content reconciliation without comparing manager hashes to ETags. |
| Adversary: arbitrary source conversion/imports and HTTP app merge | Valid concern resolved through AD-2/X1 captured source/dependency provenance, exact mapped deltas, nested preservation and explicit mixed-array/rename/incompatible blocking. No universal inversion claim. |
| Adversary: ambiguous state, competing sessions and provenance | AD-3/7/10 plus X2 fixes durable intent, target/draft locks/CAS, immutable connection descriptors, recovery maps and revision-bound observations. |
| Adversary: Docker permissions and retention | AD-6 distinguishes GET-only proxy from read-only mount; AD-11 pins active recovery and rejects all-pinned capacity rather than deleting required evidence. |
| Rubric: migration can archive and readopt live | Archive-only path replaced by functional transactional mapping, protected prior backup, local-only/disabled intent and legacy scoped history; no remote write during migration. |
| Rubric: imports are editor-managed read-only | Captured immutable bundle expansion/source maps/parity in X1 and structured form descriptors in X3 replace permanent editor-only exclusion. Unsafe constructs remain visible compatibility gaps until implemented. |
| Rubric: network isolation/debounce | Retained private reachable Admin API and exact 500ms debounce, with real startup/resume requirements and no fictitious mount inspection API. |
| Tech-fit: /load preflight sufficient or preferred | Rejected. Root POST /config/ with exact If-Match prevents the race; pinned source uses config locking/application for the conditional mutation. |
| Tech-fit: PKI endpoint and Signals pin | PKI chain/root selection and possible CA-provisioning side effect are explicit. Signals remains chosen but not installed; story 1.10 verifies/locks a compatible dependency. |
| Tech-fit: logging/storage as apps; universal semantic adapt | Corrected root/app namespaces and wrapper/result semantics. Adapt is not provisioning or module inventory; real pinned compatibility tests are assigned. |
| Tech-fit: technology seed | Existing manifests/lockfiles ratified without stack migration. Tailwind seed corrected from package range minimum 3.4.1 to lock resolution 3.4.19. |
| Current rubric R1 and adversarial enabled metadata | Closed: SourceBundleV1 identity includes immutable enabled_nodes/node_identity_map; review/attempt/recovery copy that version; historical restore cannot reuse current flags. |
| Current adversarial connection/finalization seams | Closed: immutable endpoint/transport descriptors retained by attempts/history; original-descriptor reconciliation; durable target/draft lock and revision CAS from preparation through finalization/unresolved state. |
| Current rubric R2 early UI dependency | Closed: 1.2.AC2 uses secret port/API harness; later literal form/editor integration does not block vault completion. |

Fresh reviews retain their original findings with appended disposition, preserving the evidence trail. Refer to the new readiness assessment for the planning gate; no review here is a claim of executed product acceptance.
