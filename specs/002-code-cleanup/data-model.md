# Data Model: Code Cleanup

No SQLite schema or migration changes are required.

## Imported Route Decision

A transient result of parsing and validating one route fetched from Caddy before preview or replacement.

| Field | Source | Rules |
|---|---|---|
| `support_status` | Parser / existing route metadata | `editable`, `partial_readonly`, or `unsupported_readonly` |
| `readonly` | `Route.IsReadOnly()` | Preserved routes are non-mutable |
| `readonly_reason` | Parser or legacy fallback | Required for preserved routes |
| `raw_caddy_route` | Original Caddy route | Required and valid for preserved routes; emitted unchanged on sync |
| editable `config` | Parsed handler configuration | Must be valid for `reverse_proxy`, `file_server`, or `redir` |
| decision | Shared config validation | Accept editable route, preserve read-only route, or reject the complete import before storage replacement |

### Decision transitions

```text
Caddy route
  -> parsed editable route -> accepted
  -> parsed partial/unsupported/legacy raw route -> preserved read-only
  -> missing/invalid required representation -> rejected; local routes unchanged
```

## Route Kind Configuration

The editor has three existing supported kinds. Handler-specific `config` is replaced when the selected kind changes, so fields from the previous kind are not submitted.

| Kind | `config` fields | Route-level fields |
|---|---|---|
| `reverse_proxy` | `upstreams`, `websocket`, request `headers`, `load_balancing` | domain, path, response headers, optional `strip_path_prefix` |
| `file_server` | `root`, `browse`, `index`, `hide`, `precompressed` | domain, path, response headers |
| `redir` | `to`, `code` | domain, path, response headers |

`response headers` are intentionally route-level: the existing builder supports them for every current kind. `strip_path_prefix` is not carried into file-server or redirect submissions. Read-only routes do not enter this editing model.
