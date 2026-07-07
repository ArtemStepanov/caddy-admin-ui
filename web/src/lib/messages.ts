export const DRIFT_WARNING = 'Manual Caddy changes after the last import or sync are not automatically merged. Re-run import review before syncing after manual edits.';

export function errorMessage(err: unknown): string {
  return err instanceof Error ? err.message : String(err);
}
