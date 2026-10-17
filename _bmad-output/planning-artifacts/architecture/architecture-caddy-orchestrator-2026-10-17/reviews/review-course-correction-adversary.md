# Course-correction architecture adversarial seam review

Date: 2026-10-17  
Inputs: `ARCHITECTURE-SPINE.md`, `configuration-contracts.md`, corrected `epics.md`, corrected `SPEC.md`.  
Scope: independent implementation units composing under the corrected AD-1–AD-11 and companion contracts. Inputs were not edited.

**Verdict: tighten two shared-record seams before handoff.** The corrected contracts close the previous root-preservation, opaque-secret, ETag/hash, route-identity and browser-evidence holes. The remaining findings below concern information needed across lifecycle boundaries, rather than missing implementation fixtures or reopening approved P01–P13 scope.

## F1 — High: detached target generations have no immutable connection record

**Contract locations:** AD-2 target-change rule; AD-8 target/generation checks; AD-10 unfinished-attempt reconciliation; companion P03 paragraph 6, P06 paragraph 2, and X2 `TARGET` / `APPLY_ATTEMPT` records.

**Two independently compliant units:**

- The settings/adoption unit stores the declared `TARGET {id, endpoint, adoption_generation, ...}`. Changing the configured URL updates `endpoint`, advances the generation, detaches the previous draft/history and invalidates its reviews/jobs. It retains protected prior records, as required.
- The attempt/reconciliation unit persists the declared attempt fields, including target ID/generation and exact baseline/candidate. It resolves a target's connection through `TARGET.id`, because no generation-specific endpoint record or attempt endpoint snapshot is specified.

**Composition failure:** An attempt against endpoint A becomes `submission_possible`, then loses its response. The operator switches to endpoint B. If target configuration uses the same ID and increments its generation, the old attempt still exists but its connection identity no longer does. The reconciler must either read B while handling the A attempt, or refuse to reconcile forever because A's endpoint cannot be recovered. Coincidentally equal configuration at B is particularly dangerous: content equality cannot identify which Caddy target was observed. None of the explicit required attempt fields fixes this lookup. Conversely, a settings implementation that allocates a new target ID for every URL change also satisfies the text, so independently written consumers disagree on the lifecycle.

**Smallest correction:** Define a durable immutable target-generation connection record, including the exact configured endpoint/connection profile identity. Alternatively require a new immutable target ID whenever endpoint identity changes and retain prior targets while referenced. Bind each attempt/readback/reconciliation to that record; never resolve a historical attempt through mutable current settings. Define whether switching away from an unresolved attempt is blocked or allowed with its original reconciliation lock retained. Credentials may remain immutable secret references; this does not require plaintext connection secrets.

**Acceptance:** Lose an apply response at A, switch to B, restart, and reconcile the old attempt. Requests must reach A only; matching content at B cannot finalize it. Test both return-to-A and an unavailable A. Keep the new generation's UI/events separate from the old result.

## F2 — Medium: enabled-state metadata has no shared snapshot identity

**Contract locations:** AD-1 disabled-node retention; AD-11 recovery; companion X1 `SourceBundleV1`, X2 `SOURCE_BUNDLE` / `RECOVERY_POINT` / `DRAFT`, X3 final paragraph; stories 1.4, 1.9, 6.2 and 6.3.

**Two independently compliant units:**

- The source/editor unit retains every node and its comment ID in exact source files, and persists versioned enabled flags in a separate node-metadata table "alongside source nodes," as X3 permits. Transient compilation omits disabled nodes. The declared bundle fields cover files, imports, substitutions, secret references, source-map spans and build/comparison identity, but do not explicitly include enabled-state metadata.
- The recovery unit captures the required encrypted source bundle, runtime JSON and secret references, then restores through the shared guarded pipeline. It expects the bundle to provide the complete editable historical source state. No declared recovery field tells it to acquire the source unit's separate enabled-state revision.

**Composition failure:** Recovery point R contains a disabled route's definition in source and its absence in runtime. The operator subsequently enables or edits that route. Restoring R can recover the original source bytes and runtime JSON while leaving current enabled metadata attached to the same stable node ID. The next edit/compile can silently reenable the historical disabled node, or the implementation blocks because the restored source no longer compiles to its restored baseline. Runtime absence alone cannot reconstruct the flag reliably: preserved/unmapped content and conditional/source expansion already exist in this model. Source-byte identity therefore is not the full source-state identity.

**Smallest correction:** Make enabled-state metadata an explicit immutable part of `SourceBundleV1` (or an equivalently versioned metadata manifest referenced by every bundle/draft/attempt/recovery). Key it by stable node ID and include its digest in bundle/candidate/review identity. Capture and restore that version atomically with source and runtime; never reuse present-day metadata for a historical bundle. Specify missing/unknown metadata as blocking rather than default-enabled. This implements the existing disable/restore promises without changing product scope.

**Acceptance:** Disable node N, create R, enable/edit N, restore R, then edit an unrelated node and review. N must remain disabled with R's original source. Repeat after restart, with disabled imported nodes, and with migrated local-only intent. Changing enabled metadata must invalidate an outstanding review even when original source-file bytes are unchanged.

## Other requested seams

- **Source/runtime mapping:** X1's initial identity mapping, generated/preserved disjointness, indivisible mixed arrays, and explicit server-renaming fixture requirement prevent the unsafe generic merge from the earlier review. Unsupported constructs fail closed.
- **Secrets/editor:** protected exact envelopes, safe reference tokens, server-side reference validation, withheld opaque blocks and immutable historical versions give editor/storage/apply units a consistent boundary. No plaintext snapshot or mask-as-value recommendation is made.
- **Durable attempts and draft/review concurrency:** prepared payloads, submission-possible persistence, unresolved mutation lock, exact review identities and revision-checked saves provide the required core constraints. Implementations must apply that mutation lock consistently; the wording already requires it, so ignoring it is an implementation defect rather than a new architecture finding.
- **Migration/write modes:** X4 explicitly orders encryption before import, migration before network acquisition, and durable server-side legacy/root writer exclusion. Missing prior-release fixtures are assigned story delivery.
- **Restore/retention:** scoped legacy restore, pre-restore recovery, pinned unresolved evidence and capacity rejection close the old FIFO loss issue. F2 concerns a different missing metadata reference.
- **Probe identity:** target/generation/applied-revision/expectation identity and historical late results prevent old probes from verifying current state. Unknown browser evidence remains unknown. Actual harness delivery remains implementation work.

These findings need narrow contract and acceptance additions, not further product elicitation or a change of approved feature scope.

## Closure verification — 2026-10-17

**Disposition: both findings closed; no remaining concrete blocker in this focused review.** Re-read the updated AD-1/2/3, companion X1/X2, and stories 1.4.AC4 / 1.9.AC1 / 1.9.AC3.

- **F1 closed:** `CONNECTION_DESCRIPTOR` is immutable; attempts and recovery reference it; readback/reconciliation uses the original descriptor. Connection changes create a new target/descriptor identity and cannot repoint historical records. Durable locking starts at preparation and rejects endpoint changes until resolution. Missing original transport evidence keeps the old attempt unresolved. The explicit same-URL physical-instance limitation avoids claiming an identity proof the API cannot supply.
- **F2 closed:** `SourceBundleV1` explicitly contains immutable `enabled_nodes` and `node_identity_map`; enabled metadata participates in bundle/review/attempt/recovery identity. Restore uses the historical map, and missing legacy metadata cannot default to reenabled routes. Updated story acceptance makes preservation of that map an implementation check.
- **Concurrent-finalization seam closed:** preparation atomically rechecks the reviewed revision and acquires the target/draft lock. Draft/source/secret-reference mutations are rejected through prepared/submission-possible/unresolved states. Finalization uses captured-revision CAS; unexpected mismatch remains unresolved. These constraints prevent a successful older attempt from silently clearing a newer saved candidate.

This is contract closure, not executed product evidence. Fault, restoration and compatibility fixtures remain assigned implementation deliverables.
