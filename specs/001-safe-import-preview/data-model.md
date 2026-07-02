# Data Model: Safe Import Preview

## Route

Represents one route known to the app.

### Fields

- `id`: Stable local route identifier.
- `domain`: Best available host matcher summary, or `*`/unknown summary for global routes.
- `path`: Optional path matcher summary.
- `handler_type`: Route type shown to users: `reverse_proxy`, `file_server`, `redir`, or `unknown`.
- `config`: Editable handler settings for supported routes, or an empty/minimal summary for read-only routes.
- `headers`: Optional supported header settings for editable routes.
- `strip_path_prefix`: Optional supported path-strip setting for editable routes.
- `enabled`: Whether an editable route participates in sync. Imported read-only routes remain enabled because disabling them would mutate preserved config.
- `readonly`: True when the route is preserved but not safely editable.
- `support_status`: `editable`, `partial_readonly`, or `unsupported_readonly`.
- `readonly_reason`: Human-readable primary reason when `readonly` is true.
- `raw_caddy_route`: Original Caddy route details for read-only preservation; not included in normal route list responses.
- `created_at`, `updated_at`: Local timestamps.

### Validation Rules

- Editable routes require domain, handler type, and valid handler config.
- Read-only routes require `raw_caddy_route`, `support_status`, and `readonly_reason`.
- Read-only routes cannot be edited, deleted, toggled, or field-mutated through API or UI.
- Raw route details are available only through a read-only details action.

### State Transitions

```text
Caddy route discovered
  -> editable
  -> partial_readonly
  -> unsupported_readonly

editable -> edited/toggled/deleted by UI
partial_readonly -> details only
unsupported_readonly -> details only
```

## Import Preview

Represents the review shown before import confirmation.

### Fields

- `summary.total_found`: Count of discovered HTTP routes.
- `summary.editable`: Count classified editable.
- `summary.readonly_preserved`: Count classified read-only preserved.
- `summary.unsupported`: Count classified unsupported read-only.
- `summary.local_only`: Count current local routes missing from discovered Caddy routes.
- `summary.will_replace_local`: Whether confirmation replaces the local route list.
- `groups.new_from_caddy`: Routes discovered in Caddy with no matching local route.
- `groups.will_update`: Routes matching existing local routes.
- `groups.local_only`: Local routes that will be removed on confirmation.
- `groups.readonly_preserved`: Discovered routes that will be preserved read-only.
- `warnings`: User-visible safety/drift warnings.

### Validation Rules

- Every discovered HTTP route must appear in exactly one discovered-route group and have a support status.
- Preview failure must not modify local routes.
- Empty Caddy HTTP route sets are valid and must clearly show that confirmation removes local routes.

## Import Route Row

Preview-only summary for one route.

### Fields

- `route_id`: Local route ID when the row refers to an existing local route.
- `domain`: Best available domain summary.
- `path`: Best available path summary.
- `handler_type`: Best available route type.
- `destination`: Best available upstream, redirect destination, file root, or blank.
- `support_status`: `editable`, `partial_readonly`, or `unsupported_readonly`.
- `readonly_reason`: Required for read-only rows.
- `change_type`: `new`, `update`, `local_only_remove`, or `readonly_preserve`.

### Validation Rules

- Read-only rows must include a reason.
- Local-only rows must say they will be removed if import is confirmed.

## Technical Route Details

Read-only troubleshooting view for preserved routes.

### Fields

- `route_id`: Local route ID.
- `summary`: Same user-facing summary as the route row.
- `raw_caddy_route`: Original route details as formatted JSON.

### Validation Rules

- Available only for routes with preserved raw details.
- Does not accept edits.
- Must not be returned in normal list responses.
