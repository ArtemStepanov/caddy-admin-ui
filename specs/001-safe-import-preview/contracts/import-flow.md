# Contract: Safe Import Flow

## POST /api/import-preview

Builds a preview from the live Caddy configuration. Does not change local routes.

### Request

No body.

### Success Response: 200

```json
{
  "summary": {
    "total_found": 3,
    "editable": 1,
    "readonly_preserved": 2,
    "unsupported": 1,
    "local_only": 1,
    "will_replace_local": true
  },
  "groups": {
    "new_from_caddy": [],
    "will_update": [],
    "local_only": [],
    "readonly_preserved": []
  },
  "warnings": [
    "Manual Caddy changes after the last import or sync are not automatically merged. Re-run import review before syncing after manual edits."
  ]
}
```

Each group row uses this shape:

```json
{
  "route_id": "existing-local-id-if-any",
  "domain": "app.example.com",
  "path": "/api/*",
  "handler_type": "reverse_proxy",
  "destination": "http://app:8080",
  "support_status": "editable",
  "readonly_reason": "",
  "raw_caddy_route": {},
  "change_type": "new"
}
```

`raw_caddy_route` is present only for read-only preview rows so the user can validate the exact unsupported handler before confirming import. It is exposed through the existing admin UI/API surface only; editable rows omit it.

### Error Responses

- `502`: Caddy cannot be reached. Local routes unchanged.
- `500`: Caddy response cannot be parsed or classified. Local routes unchanged.

## POST /api/import

Fetches live Caddy config again, classifies routes, and replaces local routes transactionally.

### Request

No body.

### Success Response: 200

```json
{
  "imported": 3,
  "editable": 1,
  "readonly_preserved": 2,
  "unsupported": 1,
  "message": "Configuration imported successfully",
  "warnings": [
    "Manual Caddy changes after the last import or sync are not automatically merged. Re-run import review before syncing after manual edits."
  ]
}
```

### Error Responses

- `502`: Caddy cannot be reached. Previous local routes unchanged.
- `500`: Caddy response cannot be parsed, classified, or saved. Previous local routes unchanged.

## GET /api/routes/:id/details

Returns technical details for a preserved read-only route.

### Success Response: 200

```json
{
  "route": {
    "id": "route-id",
    "domain": "app.example.com",
    "path": "",
    "handler_type": "unknown",
    "support_status": "unsupported_readonly",
    "readonly_reason": "unknown handler"
  },
  "raw_caddy_route": {}
}
```

### Error Responses

- `404`: Route not found.
- `404`: Route has no preserved raw details.

## Existing Route Mutation Rules

The existing mutation endpoints must reject preserved read-only routes:

- `PUT /api/routes/:id`
- `DELETE /api/routes/:id`
- `POST /api/routes/:id/toggle`

### Read-only Rejection Response

```json
{
  "error": "route is read-only and managed outside the UI"
}
```

Recommended status: `409 Conflict`.

## UI Contract

- Import action opens a review screen/section before confirmation.
- Preview shows summary counts and grouped rows.
- Read-only route cards show an `Unsupported / managed outside UI` label or equivalent.
- Read-only route cards do not show edit, delete, or enable/disable controls.
- Read-only route cards show `View JSON` or `Details` only.
- Sync/import areas show the manual-drift warning copy.
