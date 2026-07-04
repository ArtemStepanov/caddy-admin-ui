import { render, screen, fireEvent, waitFor } from '@testing-library/preact';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Settings } from './Settings';

function json(data: unknown, ok = true, status = 200) {
  return Promise.resolve({ ok, status, json: () => Promise.resolve(data) } as Response);
}

describe('Settings import review', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('renders drift warning, preview groups, reasons, and import result', async () => {
    const fetchMock = vi.fn((url: string, init?: RequestInit) => {
      if (url === '/api/config') return json({ config: { caddy_admin_url: 'http://localhost:2019', enable_encode: true } });
      if (url === '/api/import-preview') return json({
        summary: { total_found: 2, editable: 1, readonly_preserved: 1, unsupported: 1, local_only: 1, will_replace_local: true },
        groups: {
          new_from_caddy: [{ domain: 'app.example.com', handler_type: 'reverse_proxy', destination: 'localhost:8080', support_status: 'editable', change_type: 'new' }],
          will_update: [],
          local_only: [{ domain: 'old.example.com', handler_type: 'file_server', support_status: 'editable', change_type: 'local_only_remove' }],
          readonly_preserved: [{ domain: 'raw.example.com', handler_type: 'unknown', support_status: 'unsupported_readonly', readonly_reason: 'unknown handler', raw_caddy_route: { handle: [{ handler: 'custom_handler' }] }, change_type: 'readonly_preserve' }],
        },
        warnings: ['Manual Caddy changes after the last import or sync are not automatically merged. Re-run import review before syncing after manual edits.'],
      });
      if (url === '/api/import' && init?.method === 'POST') return json({ imported: 2, editable: 1, readonly_preserved: 1, unsupported: 1, message: 'ok', warnings: [] });
      return json({});
    });
    vi.stubGlobal('fetch', fetchMock);

    render(<Settings />);
    await screen.findByText('Configuration Import');
    expect(screen.getByText(/Manual Caddy changes/)).toBeInTheDocument();

    fireEvent.click(screen.getByText('Review Import from Caddy'));
    await screen.findByText('New from Caddy (1)');
    expect(screen.getByText('Read-only preserved (1)')).toBeInTheDocument();
    expect(screen.getByText('Local-only routes removed on confirm (1)')).toBeInTheDocument();
    expect(screen.getByText('unknown handler')).toBeInTheDocument();
    expect(screen.getByText('View JSON')).toBeInTheDocument();
    expect(screen.getByText(/custom_handler/)).toBeInTheDocument();

    fireEvent.click(screen.getByText('Confirm Import'));
    await waitFor(() => expect(screen.getByText(/Imported 2 routes/)).toBeInTheDocument());
  });
});
