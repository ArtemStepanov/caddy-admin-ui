# Historical specifications

The numbered directories in this folder are retained as implementation history. They describe the earlier import-preview design and are not the current API contract.

The current behavior is documented in the root [README](../README.md): setup uses `POST /api/setup/preview` and `POST /api/setup/confirm`, ownership is limited to one selected HTTP server's route array, external routes are preserved read-only, and writes are guarded by Caddy ETags.

