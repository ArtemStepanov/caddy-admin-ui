import { notifySyncResult } from './syncNotify';

const API_BASE = '/api';

export interface HeaderConfig {
  set?: Record<string, string>;
  add?: Record<string, string>;
  delete?: string[];
}

export type SupportStatus = 'editable' | 'partial_readonly' | 'unsupported_readonly';

export interface Route {
  id: string;
  domain: string;
  path?: string;
  strip_path_prefix?: string;
  handler_type: string;
  config: any;
  headers?: HeaderConfig;
  enabled: boolean;
  readonly?: boolean;
  support_status?: SupportStatus;
  readonly_reason?: string;
  created_at: string;
  updated_at: string;
}

export interface ImportRouteRow {
  route_id?: string;
  domain: string;
  path?: string;
  handler_type: string;
  destination?: string;
  support_status: SupportStatus;
  readonly_reason?: string;
  raw_caddy_route?: unknown;
  change_type: 'new' | 'update' | 'local_only_remove' | 'readonly_preserve';
}

export interface ImportPreview {
  summary: {
    total_found: number;
    editable: number;
    readonly_preserved: number;
    unsupported: number;
    local_only: number;
    will_update: number;
    will_replace_local: boolean;
  };
  groups: {
    new_from_caddy: ImportRouteRow[];
    will_update: ImportRouteRow[];
    local_only: ImportRouteRow[];
    readonly_preserved: ImportRouteRow[];
  };
  warnings: string[];
}

export interface ImportResult {
  imported: number;
  editable: number;
  readonly_preserved: number;
  unsupported: number;
  message: string;
  warnings: string[];
}

export interface RouteDetails {
  route: Pick<Route, 'id' | 'domain' | 'path' | 'handler_type' | 'support_status' | 'readonly_reason'>;
  raw_caddy_route: unknown;
}

export interface GlobalConfig {
  caddy_admin_url: string;
  enable_encode: boolean;
}

export interface StatusResponse {
  status: 'online' | 'offline';
  latency?: number;
  error?: string;
  admin_url?: string;
  route_count?: number;
  last_synced_at?: string;
  last_sync_error?: string;
  version?: string;
}

class ApiClient {
  private async request<T>(path: string, options?: RequestInit): Promise<T> {
    const res = await fetch(`${API_BASE}${path}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    const data = await res.json();

    if (!res.ok) {
      throw new Error(data.error || `HTTP ${res.status}`);
    }

    return data;
  }

  // Routes
  async listRoutes(): Promise<{ routes: Route[] }> {
    return this.request('/routes');
  }

  async getRoute(id: string): Promise<{ route: Route }> {
    return this.request(`/routes/${id}`);
  }

  async createRoute(route: Partial<Route>): Promise<{ route: Route; warning?: string }> {
    const res = await this.request<{ route: Route; warning?: string }>('/routes', {
      method: 'POST',
      body: JSON.stringify(route),
    });
    if (res.warning) notifySyncResult('error', res.warning);
    else notifySyncResult('success', 'Route created and synced');
    return res;
  }

  async updateRoute(id: string, route: Partial<Route>): Promise<{ route: Route; warning?: string }> {
    const res = await this.request<{ route: Route; warning?: string }>(`/routes/${id}`, {
      method: 'PUT',
      body: JSON.stringify(route),
    });
    if (res.warning) notifySyncResult('error', res.warning);
    else notifySyncResult('success', 'Route updated and synced');
    return res;
  }

  async deleteRoute(id: string): Promise<{ message: string }> {
    const res = await this.request<{ message: string }>(`/routes/${id}`, { method: 'DELETE' });
    if (res.message && res.message.includes('failed')) notifySyncResult('error', res.message);
    else notifySyncResult('success', 'Route deleted and synced');
    return res;
  }

  async toggleRoute(id: string): Promise<{ route: Route; warning?: string }> {
    const res = await this.request<{ route: Route; warning?: string }>(`/routes/${id}/toggle`, { method: 'POST' });
    if (res.warning) notifySyncResult('error', res.warning);
    else notifySyncResult('success', `Route ${res.route.enabled ? 'enabled' : 'disabled'} and synced`);
    return res;
  }

  // Config
  async getConfig(): Promise<{ config: GlobalConfig }> {
    return this.request('/config');
  }

  async updateConfig(config: GlobalConfig): Promise<{ config: GlobalConfig }> {
    return this.request('/config', {
      method: 'PUT',
      body: JSON.stringify(config),
    });
  }

  // Status
  async getStatus(): Promise<StatusResponse> {
    return this.request('/status');
  }

  async sync(): Promise<{ message: string }> {
    return this.request('/sync', { method: 'POST' });
  }

  async testConnection(url: string): Promise<{ success: boolean; latency?: number; error?: string }> {
    return this.request('/test-connection', {
      method: 'POST',
      body: JSON.stringify({ url }),
    });
  }

  async previewImport(): Promise<ImportPreview> {
    return this.request('/import-preview', { method: 'POST' });
  }

  async importFromCaddy(): Promise<ImportResult> {
    return this.request('/import', { method: 'POST' });
  }

  async getRouteDetails(id: string): Promise<RouteDetails> {
    return this.request(`/routes/${id}/details`);
  }
}

export const api = new ApiClient();
