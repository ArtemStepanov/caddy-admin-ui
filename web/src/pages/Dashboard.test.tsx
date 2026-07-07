import { render, screen, fireEvent, waitFor } from '@testing-library/preact';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Dashboard } from './Dashboard';

function json(data: unknown, ok = true, status = 200) {
  return Promise.resolve({ ok, status, json: () => Promise.resolve(data) } as Response);
}

describe('Dashboard read-only routes', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('shows drift warning, hides read-only mutations, and displays JSON details', async () => {
    vi.stubGlobal('fetch', vi.fn((url: string) => {
      if (url === '/api/routes') return json({ routes: [
        { id: '1', domain: 'raw.example.com', handler_type: 'unknown', config: {}, enabled: true, readonly: true, support_status: 'unsupported_readonly', readonly_reason: 'unknown handler', created_at: '2026-01-01T00:00:00Z', updated_at: '2026-01-01T00:00:00Z' },
        { id: '2', domain: 'app.example.com', handler_type: 'reverse_proxy', config: { upstreams: ['localhost:8080'] }, enabled: true, created_at: '2026-01-01T00:00:00Z', updated_at: '2026-01-01T00:00:00Z' },
      ] });
      if (url === '/api/routes/1/details') return json({ route: { id: '1', domain: 'raw.example.com', handler_type: 'unknown', support_status: 'unsupported_readonly', readonly_reason: 'unknown handler' }, raw_caddy_route: { handle: [{ handler: 'custom' }] } });
      return json({});
    }));

    render(<Dashboard />);
    await screen.findByText('raw.example.com');
    expect(screen.getByText(/Manual Caddy changes/)).toBeInTheDocument();
    expect(screen.getByText('Unsupported / managed outside UI')).toBeInTheDocument();
    expect(screen.getByText('View JSON')).toBeInTheDocument();
    expect(screen.getByText('Edit')).toBeInTheDocument();
    expect(screen.getByText('Delete')).toBeInTheDocument();

    fireEvent.click(screen.getByText('View JSON'));
    await waitFor(() => expect(screen.getByText(/custom/)).toBeInTheDocument());
  });
});
