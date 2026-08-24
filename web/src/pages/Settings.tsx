import { useEffect, useState } from "preact/hooks";
import { api, GlobalConfig, SetupPreview, Snapshot } from "../lib/api";
import { DRIFT_WARNING, errorMessage } from "../lib/messages";

function RoutePreview({ preview }: { preview: SetupPreview }) {
  return (
    <div class="mt-5 space-y-4">
      <div class="grid grid-cols-3 gap-3 text-center text-sm">
        <div class="bg-slate-900/60 rounded p-3">
          <div class="text-xl font-bold">{preview.route_count}</div>
          <div class="text-slate-400">Routes</div>
        </div>
        <div class="bg-slate-900/60 rounded p-3">
          <div class="text-xl font-bold text-green-300">{preview.editable}</div>
          <div class="text-slate-400">UI-owned</div>
        </div>
        <div class="bg-slate-900/60 rounded p-3">
          <div class="text-xl font-bold text-amber-300">{preview.readonly}</div>
          <div class="text-slate-400">Preserved</div>
        </div>
      </div>
      <div class="p-3 bg-blue-900/20 border border-blue-800/60 rounded text-blue-200 text-sm">
        {preview.ownership_notice}
      </div>
      {preview.caddy_empty && (
        <div class="p-3 bg-amber-900/20 border border-amber-800/60 rounded text-amber-200 text-sm">
          Caddy is empty. Confirmation will bootstrap one dedicated server on
          ports 80 and 443.
        </div>
      )}
      {preview.routes.length > 0 && (
        <div class="max-h-64 overflow-auto space-y-2">
          {preview.routes.map((route, index) => (
            <div
              key={`${route.id}-${index}`}
              class="bg-slate-900/60 rounded p-3 text-sm flex items-center justify-between gap-3"
            >
              <div>
                <span class="font-medium">{route.domain}</span>
                {route.path && (
                  <span class="text-slate-400"> {route.path}</span>
                )}
              </div>
              <span
                class={`text-xs px-2 py-0.5 rounded ${route.readonly ? "bg-amber-800/50 text-amber-200" : "bg-green-800/50 text-green-200"}`}
              >
                {route.readonly ? "preserved read-only" : "UI-owned"}
              </span>
            </div>
          ))}
        </div>
      )}
      {(preview.local_drafts || []).length > 0 && (
        <div class="p-3 bg-slate-900/60 border border-slate-700 rounded text-slate-300 text-sm">
          {preview.local_drafts.length} local UI-owned draft
          {preview.local_drafts.length === 1 ? "" : "s"} will remain local and
          apply only on the next guarded sync.
        </div>
      )}
    </div>
  );
}

export function Settings() {
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [testResult, setTestResult] = useState<{
    success: boolean;
    latency?: number;
    error?: string;
  } | null>(null);
  const [preview, setPreview] = useState<SetupPreview | null>(null);
  const [snapshots, setSnapshots] = useState<Snapshot[]>([]);
  const [config, setConfig] = useState<GlobalConfig>({
    caddy_admin_url: "http://localhost:2019",
    enable_encode: true,
    setup_complete: false,
  });

  async function refresh() {
    try {
      const [{ config: current }, { snapshots: recent }] = await Promise.all([
        api.getConfig(),
        api.listSnapshots(),
      ]);
      setConfig(current);
      setSnapshots(recent || []);
    } catch (err: unknown) {
      setError(errorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  async function saveSettings(e: Event) {
    e.preventDefault();
    setBusy(true);
    setError(null);
    setSuccess(null);
    try {
      const { config: saved } = await api.updateConfig(config);
      setConfig(saved);
      setPreview(null);
      setSuccess(
        saved.setup_complete
          ? "Settings saved."
          : "Settings saved. Review and confirm Caddy ownership before syncing.",
      );
    } catch (err: unknown) {
      setError(errorMessage(err));
    } finally {
      setBusy(false);
    }
  }

  async function testConnection() {
    setBusy(true);
    setTestResult(null);
    try {
      setTestResult(await api.testConnection(config.caddy_admin_url));
    } catch (err: unknown) {
      setTestResult({ success: false, error: errorMessage(err) });
    } finally {
      setBusy(false);
    }
  }

  async function reviewSetup(server?: string) {
    setBusy(true);
    setError(null);
    setSuccess(null);
    setPreview(null);
    try {
      setPreview(await api.previewSetup(config.caddy_admin_url, server));
    } catch (err: unknown) {
      setError("Setup preview failed: " + errorMessage(err));
    } finally {
      setBusy(false);
    }
  }

  async function confirmSetup() {
    if (!preview) return;
    setBusy(true);
    setError(null);
    try {
      const result = await api.confirmSetup(preview);
      setConfig(result.config);
      setPreview(null);
      setSuccess(
        `Connected to ${result.config.managed_server}. Imported ${result.imported} routes.`,
      );
      await refresh();
    } catch (err: unknown) {
      setError("Setup confirmation failed: " + errorMessage(err));
    } finally {
      setBusy(false);
    }
  }

  async function restoreSnapshot(snapshot: Snapshot) {
    if (
      !confirm(
        `Restore routes from ${new Date(snapshot.created_at).toLocaleString()}?`,
      )
    )
      return;
    setBusy(true);
    setError(null);
    try {
      const result = await api.restoreSnapshot(snapshot.id);
      setSuccess(`Snapshot restored (${result.restored} routes).`);
      await refresh();
    } catch (err: unknown) {
      setError("Restore failed: " + errorMessage(err));
    } finally {
      setBusy(false);
    }
  }

  if (loading)
    return (
      <div class="text-center py-8 text-slate-400">Loading settings...</div>
    );

  return (
    <div class="max-w-3xl mx-auto space-y-6">
      <h1 class="text-2xl font-bold">Settings</h1>
      {error && (
        <div class="bg-red-900/50 border border-red-700 rounded-lg p-4 text-red-300">
          {error}
        </div>
      )}
      {success && (
        <div class="bg-green-900/50 border border-green-700 rounded-lg p-4 text-green-300">
          {success}
        </div>
      )}

      <form onSubmit={saveSettings} class="card space-y-5">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h2 class="text-lg font-semibold">Caddy Connection</h2>
            <p class="text-sm text-slate-400 mt-1">
              Connect explicitly, review the server, then confirm ownership.
            </p>
          </div>
          <span
            class={`text-xs px-2 py-1 rounded ${config.setup_complete ? "bg-green-800/50 text-green-200" : "bg-amber-800/50 text-amber-200"}`}
          >
            {config.setup_complete
              ? `Managing ${config.managed_server}`
              : "Setup required"}
          </span>
        </div>
        <div>
          <label class="label">Caddy Admin API URL</label>
          <input
            type="url"
            value={config.caddy_admin_url}
            onInput={(e) =>
              setConfig({
                ...config,
                caddy_admin_url: (e.target as HTMLInputElement).value,
              })
            }
            placeholder="http://caddy:2019"
            class="input"
            required
          />
          <p class="text-sm text-slate-500 mt-1">
            Changing this URL clears ownership; saving never writes to the new
            Caddy instance.
          </p>
        </div>
        <div class="flex flex-wrap gap-2 items-center">
          <button
            type="button"
            onClick={testConnection}
            disabled={busy}
            class="btn btn-secondary"
          >
            Test Connection
          </button>
          <button
            type="button"
            onClick={() => reviewSetup()}
            disabled={busy}
            class="btn btn-secondary"
          >
            Review Caddy Servers
          </button>
          <button type="submit" disabled={busy} class="btn btn-primary">
            Save Settings
          </button>
          {testResult && (
            <span
              class={
                testResult.success
                  ? "text-green-400 text-sm"
                  : "text-red-400 text-sm"
              }
            >
              {testResult.success
                ? `Connected (${testResult.latency}ms)`
                : testResult.error}
            </span>
          )}
        </div>

        <label class="flex items-center gap-3 cursor-pointer pt-2 border-t border-slate-700">
          <input
            type="checkbox"
            checked={config.enable_encode}
            onChange={(e) =>
              setConfig({
                ...config,
                enable_encode: (e.target as HTMLInputElement).checked,
              })
            }
            class="w-5 h-5 rounded bg-slate-900 border-slate-700"
          />
          <div>
            <div class="font-medium">Response compression</div>
            <div class="text-sm text-slate-500">
              Add gzip and zstd handlers to UI-owned routes.
            </div>
          </div>
        </label>
      </form>

      {preview && (
        <div class="card">
          <h2 class="text-lg font-semibold">Ownership Review</h2>
          {preview.servers.length > 1 && (
            <div class="mt-4">
              <label class="label">HTTP Server</label>
              <select
                class="input"
                value={preview.selected_server}
                onChange={(e) =>
                  reviewSetup((e.target as HTMLSelectElement).value)
                }
              >
                {preview.servers.map((server) => (
                  <option key={server} value={server}>
                    {server}
                  </option>
                ))}
              </select>
            </div>
          )}
          <RoutePreview preview={preview} />
          <div class="p-3 bg-amber-900/20 border border-amber-800/60 rounded text-amber-200 text-sm mt-4">
            {DRIFT_WARNING}
          </div>
          <div class="flex justify-end mt-4">
            <button
              type="button"
              onClick={confirmSetup}
              disabled={busy}
              class="btn btn-primary"
            >
              Confirm Ownership
            </button>
          </div>
        </div>
      )}

      <div class="card">
        <h2 class="text-lg font-semibold mb-1">Restore Points</h2>
        <p class="text-sm text-slate-400 mb-4">
          A snapshot is saved before every Caddy write. Restore is also
          ETag-guarded.
        </p>
        {snapshots.length === 0 ? (
          <p class="text-sm text-slate-500">No snapshots yet.</p>
        ) : (
          <div class="space-y-2">
            {snapshots.map((snapshot) => (
              <div
                key={snapshot.id}
                class="bg-slate-900/60 rounded p-3 flex items-center justify-between gap-3"
              >
                <div>
                  <div class="text-sm font-medium">{snapshot.reason}</div>
                  <div class="text-xs text-slate-500">
                    {snapshot.server} ·{" "}
                    {new Date(snapshot.created_at).toLocaleString()}
                  </div>
                </div>
                <div class="flex gap-2">
                  <a
                    class="btn btn-secondary text-sm"
                    href={api.snapshotExportURL(snapshot.id)}
                  >
                    Export
                  </a>
                  <button
                    type="button"
                    class="btn btn-secondary text-sm"
                    disabled={busy}
                    onClick={() => restoreSnapshot(snapshot)}
                  >
                    Restore
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div class="card text-sm text-slate-400">
        <h2 class="text-lg font-semibold text-slate-200 mb-2">
          About Caddy Admin UI
        </h2>
        Safely manage a selected Caddy HTTP server's routes while preserving
        unsupported configuration.
      </div>
    </div>
  );
}
