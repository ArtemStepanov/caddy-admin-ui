# Contract: Import Safety and Route Editor

This feature changes no endpoint paths and adds no response fields. It consolidates the existing behavior below.

## `POST /api/import-preview` and `POST /api/import`

Both endpoints fetch Caddy configuration, parse every route, and apply the same server-side validation before returning a preview or replacing local routes.

| Parsed route | Required representation | Outcome |
|---|---|---|
| Editable | domain, supported handler type, and valid handler config | Included as `support_status: "editable"` |
| Preserved | `raw_caddy_route` and `readonly_reason` | Included as a read-only preserved route; raw content remains available only through existing protected details/preview responses |
| Invalid or incomplete | Any required representation missing or invalid | `500` error; local routes are not replaced |

Existing `502` Caddy-connectivity failures, `500` error shape, recovery guidance on `POST /api/import`, success response fields, and drift warnings remain unchanged. See [data-model.md](../data-model.md) for the decision fields.

## Route editor submissions

The existing `POST /api/routes` and `PUT /api/routes/:id` request shape is unchanged. The editor submits only the selected handler's `config` shape:

| `handler_type` | `config` |
|---|---|
| `reverse_proxy` | reverse-proxy fields; may include `strip_path_prefix` |
| `file_server` | file-server fields; omits `strip_path_prefix` |
| `redir` | redirect fields; omits `strip_path_prefix` |

Response headers remain route-level and may be sent for any of the three existing kinds. Preserved read-only routes remain outside the editor and existing mutation endpoints continue to reject them.
