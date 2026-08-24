import { notifySyncResult } from "./syncNotify";

const API_BASE = "/api";

export interface HeaderConfig {
  set?: Record<string, string>;
  add?: Record<string, string>;
  delete?: string[];
}

export type SupportStatus =
  | "editable"
  | "partial_readonly"
  | "unsupported_readonly";

export interface ReverseProxyConfig {
  upstreams: string[];
  headers?: Record<string, string>;
  load_balancing?: string;
}

export interface FileServerConfig {
  root: string;
  browse?: boolean;
  index?: string[];
  hide?: string[];
  precompressed?: boolean;
}

export interface RedirectConfig {
  to: string;
  code?: number;
}

export type RouteConfig =
  | ReverseProxyConfig
  | FileServerConfig
  | RedirectConfig
  | Record<string, unknown>;

export interface Route {
  id: string;
  domain: string;
  path?: string;
  strip_path_prefix?: string;
  handler_type: string;
  config: RouteConfig;
  headers?: HeaderConfig;
  enabled: boolean;
  readonly?: boolean;
  support_status?: SupportStatus;
  readonly_reason?: string;
  created_at: string;
  updated_at: string;
  position: number;
}

export interface RouteDetails {
  route: Pick<
    Route,
    | "id"
    | "domain"
    | "path"
    | "handler_type"
    | "support_status"
    | "readonly_reason"
  >;
  raw_caddy_route: unknown;
}

export interface GlobalConfig {
  caddy_admin_url: string;
  enable_encode: boolean;
  managed_server?: string;
  setup_complete: boolean;
  last_etag?: string;
}

export interface SetupPreview {
  url: string;
  servers: string[];
  selected_server: string;
  routes: Route[];
  local_drafts: Route[];
  route_count: number;
  editable: number;
  readonly: number;
  caddy_empty: boolean;
  can_bootstrap: boolean;
  preview_token: string;
  ownership_notice: string;
}

export interface Snapshot {
  id: string;
  server: string;
  etag: string;
  reason: string;
  created_at: string;
}

export interface StatusResponse {
  status: "online" | "offline";
  latency?: number;
  error?: string;
  admin_url?: string;
  route_count?: number;
  last_synced_at?: string;
  last_sync_error?: string;
  version?: string;
  setup_complete?: boolean;
  managed_server?: string;
}

class ApiClient {
  private async request<T>(path: string, options?: RequestInit): Promise<T> {
    const res = await fetch(`${API_BASE}${path}`, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        "X-Caddy-Admin-UI": "1",
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
    return this.request("/routes");
  }

  async getRoute(id: string): Promise<{ route: Route }> {
    return this.request(`/routes/${id}`);
  }

  async createRoute(
    route: Partial<Route>,
  ): Promise<{ route: Route; warning?: string }> {
    const res = await this.request<{ route: Route; warning?: string }>(
      "/routes",
      {
        method: "POST",
        body: JSON.stringify(route),
      },
    );
    if (res.warning) notifySyncResult("error", res.warning);
    else notifySyncResult("success", "Route created and synced");
    return res;
  }

  async updateRoute(
    id: string,
    route: Partial<Route>,
  ): Promise<{ route: Route; warning?: string }> {
    const res = await this.request<{ route: Route; warning?: string }>(
      `/routes/${id}`,
      {
        method: "PUT",
        body: JSON.stringify(route),
      },
    );
    if (res.warning) notifySyncResult("error", res.warning);
    else notifySyncResult("success", "Route updated and synced");
    return res;
  }

  async deleteRoute(id: string): Promise<{ message: string }> {
    const res = await this.request<{ message: string }>(`/routes/${id}`, {
      method: "DELETE",
    });
    if (res.message && res.message.includes("failed"))
      notifySyncResult("error", res.message);
    else notifySyncResult("success", "Route deleted and synced");
    return res;
  }

  async toggleRoute(id: string): Promise<{ route: Route; warning?: string }> {
    const res = await this.request<{ route: Route; warning?: string }>(
      `/routes/${id}/toggle`,
      { method: "POST" },
    );
    if (res.warning) notifySyncResult("error", res.warning);
    else
      notifySyncResult(
        "success",
        `Route ${res.route.enabled ? "enabled" : "disabled"} and synced`,
      );
    return res;
  }

  // Config
  async getConfig(): Promise<{ config: GlobalConfig }> {
    return this.request("/config");
  }

  async updateConfig(config: GlobalConfig): Promise<{ config: GlobalConfig }> {
    return this.request("/config", {
      method: "PUT",
      body: JSON.stringify(config),
    });
  }

  // Status
  async getStatus(): Promise<StatusResponse> {
    return this.request("/status");
  }

  async sync(): Promise<{ message: string }> {
    return this.request("/sync", { method: "POST" });
  }

  async testConnection(
    url: string,
  ): Promise<{ success: boolean; latency?: number; error?: string }> {
    return this.request("/test-connection", {
      method: "POST",
      body: JSON.stringify({ url }),
    });
  }

  async previewSetup(url: string, server?: string): Promise<SetupPreview> {
    return this.request("/setup/preview", {
      method: "POST",
      body: JSON.stringify({ url, server }),
    });
  }

  async confirmSetup(
    preview: SetupPreview,
  ): Promise<{ config: GlobalConfig; imported: number; message: string }> {
    return this.request("/setup/confirm", {
      method: "POST",
      body: JSON.stringify({
        url: preview.url,
        server: preview.selected_server,
        preview_token: preview.preview_token,
      }),
    });
  }

  async listSnapshots(): Promise<{ snapshots: Snapshot[] }> {
    return this.request("/snapshots");
  }

  async restoreSnapshot(
    id: string,
  ): Promise<{ message: string; restored: number }> {
    return this.request(`/snapshots/${id}/restore`, { method: "POST" });
  }

  snapshotExportURL(id: string): string {
    return `${API_BASE}/snapshots/${id}/export`;
  }

  async getRouteDetails(id: string): Promise<RouteDetails> {
    return this.request(`/routes/${id}/details`);
  }
}

export const api = new ApiClient();
