export const DRIFT_WARNING =
  "Sync is guarded by Caddy ETags. External route changes stop writes until you review setup again.";

export function errorMessage(err: unknown): string {
  return err instanceof Error ? err.message : String(err);
}
