import { useState, useEffect } from 'preact/hooks';
import { api, GlobalConfig, ImportPreview, ImportRouteRow, ImportResult } from '../lib/api';
import { DRIFT_WARNING, errorMessage } from '../lib/messages';

function ImportRows({ title, rows }: { title: string; rows: ImportRouteRow[] }) {
  if (rows.length === 0) return null;
  return (
    <div>
      <h3 class="font-medium text-slate-200 mb-2">{title} ({rows.length})</h3>
      <div class="space-y-2">
        {rows.map((row, i) => (
          <div key={`${row.change_type}-${row.domain}-${row.path}-${row.handler_type}-${i}`} class="bg-slate-900/60 rounded p-3 text-sm">
            <div class="flex items-center justify-between gap-3">
              <span class="font-medium">{row.domain}{row.path ? ` ${row.path}` : ''}</span>
              <span class="text-xs px-2 py-0.5 rounded bg-slate-700">{row.support_status}</span>
            </div>
            <div class="text-slate-400 mt-1">{row.handler_type}{row.destination ? ` → ${row.destination}` : ''}</div>
            {row.readonly_reason && <div class="text-amber-300 mt-1">{row.readonly_reason}</div>}
            {row.raw_caddy_route && (
              <details class="mt-2">
                <summary class="cursor-pointer text-slate-300">View JSON</summary>
                <pre class="mt-2 text-xs overflow-auto bg-black/30 rounded p-2">{JSON.stringify(row.raw_caddy_route, null, 2)}</pre>
              </details>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}

export function Settings() {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [testing, setTesting] = useState(false);
  const [testResult, setTestResult] = useState<{ success: boolean; latency?: number; error?: string } | null>(null);
  const [importing, setImporting] = useState(false);
  const [importPreview, setImportPreview] = useState<ImportPreview | null>(null);
  const [importResult, setImportResult] = useState<ImportResult | null>(null);

  const [config, setConfig] = useState<GlobalConfig>({
    caddy_admin_url: 'http://localhost:2019',
    enable_encode: true,
  });

  async function loadConfig() {
    try {
      const { config: cfg } = await api.getConfig();
      setConfig(cfg);
    } catch (err: unknown) {
      setError(errorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadConfig();
  }, []);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    setSaving(true);
    setError(null);
    setSuccess(false);

    try {
      await api.updateConfig(config);
      await api.sync();
      setSuccess(true);
      setTimeout(() => setSuccess(false), 3000);
    } catch (err: unknown) {
      setError(errorMessage(err));
    } finally {
      setSaving(false);
    }
  }

  async function testConnection() {
    setTesting(true);
    setTestResult(null);
    try {
      const result = await api.testConnection(config.caddy_admin_url);
      setTestResult(result);
    } catch (err: unknown) {
      setTestResult({ success: false, error: errorMessage(err) });
    } finally {
      setTesting(false);
    }
  }

  async function handleImport() {
    setImporting(true);
    setError(null);
    setImportResult(null);
    setImportPreview(null);
    try {
      setImportPreview(await api.previewImport());
    } catch (err: unknown) {
      setError("Import preview failed: " + errorMessage(err));
    } finally {
      setImporting(false);
    }
  }

  async function confirmImport() {
    setImporting(true);
    setError(null);
    try {
      const result = await api.importFromCaddy();
      setImportResult(result);
      setSuccess(true);
      setTimeout(() => setSuccess(false), 3000);
    } catch (err: unknown) {
      setError("Import failed: " + errorMessage(err));
    } finally {
      setImporting(false);
    }
  }

  if (loading) {
    return (
      <div class="text-center py-8">
        <div class="text-slate-400">Loading settings...</div>
      </div>
    );
  }

  return (
    <div class="max-w-2xl mx-auto">
      <h1 class="text-2xl font-bold mb-6">Settings</h1>

      {error && (
        <div class="bg-red-900/50 border border-red-700 rounded-lg p-4 mb-6 text-red-300">
          {error}
        </div>
      )}

      {success && (
        <div class="bg-green-900/50 border border-green-700 rounded-lg p-4 mb-6 text-green-300">
          Settings saved and applied successfully!
        </div>
      )}

      <form onSubmit={handleSubmit} class="space-y-6">
        {/* Caddy Connection */}
        <div class="card">
          <h2 class="text-lg font-semibold mb-4">Caddy Connection</h2>

          <div>
            <label class="label">Caddy Admin API URL</label>
            <input
              type="url"
              value={config.caddy_admin_url}
              onInput={(e) => setConfig({ ...config, caddy_admin_url: (e.target as HTMLInputElement).value })}
              placeholder="http://localhost:2019"
              class="input"
              required
            />
            <p class="text-sm text-slate-500 mt-1">
              The URL of your Caddy server's admin API (default port is 2019)
            </p>
          </div>

          <div class="flex gap-2 mt-4">
            <button type="button" onClick={testConnection} disabled={testing} class="btn btn-secondary">
              {testing ? 'Testing...' : 'Test Connection'}
            </button>

            {testResult && (
              <div class={`flex items-center gap-2 text-sm ${testResult.success ? 'text-green-400' : 'text-red-400'}`}>
                {testResult.success ? (
                  <>✓ Connected ({testResult.latency}ms)</>
                ) : (
                  <>✗ {testResult.error}</>
                )}
              </div>
            )}
          </div>
        </div>

        {/* Compression */}
        <div class="card">
          <h2 class="text-lg font-semibold mb-4">Response Compression</h2>

          <label class="flex items-center gap-3 cursor-pointer">
            <input
              type="checkbox"
              checked={config.enable_encode}
              onChange={(e) => setConfig({ ...config, enable_encode: (e.target as HTMLInputElement).checked })}
              class="w-5 h-5 rounded bg-slate-900 border-slate-700"
            />
            <div>
              <div class="font-medium">Enable Compression</div>
              <div class="text-sm text-slate-500">
                Compress responses with gzip and zstd for faster loading
              </div>
            </div>
          </label>

          <div class="mt-4 bg-slate-800/50 rounded-lg p-4 text-sm">
            <div class="font-medium text-slate-300 mb-2">About Compression</div>
            <ul class="space-y-1 text-slate-400">
              <li>• Reduces response size by 60-90% for text content</li>
              <li>• Supports gzip (universal) and zstd (faster, newer)</li>
              <li>• Automatically negotiates best format with browser</li>
            </ul>
          </div>
        </div>

        {/* Configuration Sync */}
        <div class="card">
          <h2 class="text-lg font-semibold mb-4">Configuration Import</h2>

          <div class="mb-4 text-sm text-slate-400">
            Import configuration from the running Caddy instance.
            <div class="mt-2 p-3 bg-amber-900/20 border border-amber-900/50 rounded text-amber-200">
              <strong>Warning:</strong> Confirming import replaces all local routes. Unsupported Caddy routes are preserved read-only.
              <div class="mt-2">{DRIFT_WARNING}</div>
            </div>
          </div>

          <div class="flex gap-2">
            <button
              type="button"
              onClick={handleImport}
              disabled={importing}
              class="btn bg-slate-800 hover:bg-slate-700 border-slate-600 text-slate-200"
            >
              {importing ? 'Loading preview...' : 'Review Import from Caddy'}
            </button>
            {importPreview && (
              <button type="button" onClick={confirmImport} disabled={importing} class="btn btn-primary">
                {importing ? 'Importing...' : 'Confirm Import'}
              </button>
            )}
          </div>

          {importResult && (
            <div class="mt-4 p-3 bg-green-900/30 border border-green-800 rounded text-green-200 text-sm">
              Imported {importResult.imported} routes ({importResult.editable} editable, {importResult.readonly_preserved} read-only preserved, {importResult.unsupported} unsupported).
            </div>
          )}

          {importPreview && (
            <div class="mt-6 space-y-5">
              <div class="grid grid-cols-2 md:grid-cols-6 gap-3 text-center text-sm">
                <div class="bg-slate-900/60 rounded p-3"><div class="text-xl font-bold">{importPreview.summary.total_found}</div><div class="text-slate-400">Found</div></div>
                <div class="bg-slate-900/60 rounded p-3"><div class="text-xl font-bold text-green-300">{importPreview.summary.editable}</div><div class="text-slate-400">Editable</div></div>
                <div class="bg-slate-900/60 rounded p-3"><div class="text-xl font-bold text-blue-300">{importPreview.summary.will_update}</div><div class="text-slate-400">Will update</div></div>
                <div class="bg-slate-900/60 rounded p-3"><div class="text-xl font-bold text-amber-300">{importPreview.summary.readonly_preserved}</div><div class="text-slate-400">Read-only</div></div>
                <div class="bg-slate-900/60 rounded p-3"><div class="text-xl font-bold text-red-300">{importPreview.summary.unsupported}</div><div class="text-slate-400">Unsupported</div></div>
                <div class="bg-slate-900/60 rounded p-3"><div class="text-xl font-bold">{importPreview.summary.local_only}</div><div class="text-slate-400">Local-only removed</div></div>
              </div>
              {(importPreview.warnings || []).map((warning) => (
                <div key={warning} class="p-3 bg-amber-900/20 border border-amber-900/50 rounded text-amber-200 text-sm">{warning}</div>
              ))}
              <ImportRows title="New from Caddy" rows={importPreview.groups.new_from_caddy} />
              <ImportRows title="Will update" rows={importPreview.groups.will_update} />
              <ImportRows title="Read-only preserved" rows={importPreview.groups.readonly_preserved} />
              <ImportRows title="Local-only routes removed on confirm" rows={importPreview.groups.local_only} />
            </div>
          )}
        </div>

        {/* About */}
        <div class="card">
          <h2 class="text-lg font-semibold mb-4">About</h2>

          <div class="space-y-2 text-sm text-slate-400">
            <div>
              <span class="text-slate-300">Caddy Admin UI</span>
            </div>
            <div>
              A simple web UI for managing Caddy server routes.
            </div>
          </div>
        </div>

        {/* Submit */}
        <div class="flex justify-end">
          <button type="submit" class="btn btn-primary" disabled={saving}>
            {saving ? 'Saving...' : 'Save Settings'}
          </button>
        </div>
      </form>
    </div>
  );
}
