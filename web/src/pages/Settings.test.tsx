import { fireEvent, render, screen, waitFor } from "@testing-library/preact";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { Settings } from "./Settings";

function json(data: unknown, ok = true, status = 200) {
  return Promise.resolve({
    ok,
    status,
    json: () => Promise.resolve(data),
  } as Response);
}

const baseConfig = {
  caddy_admin_url: "http://caddy:2019",
  enable_encode: true,
  setup_complete: false,
};

describe("Settings managed setup", () => {
  beforeEach(() => vi.restoreAllMocks());

  it("reviews ownership and confirms the content-bound preview", async () => {
    const fetchMock = vi.fn((url: string, init?: RequestInit) => {
      if (url === "/api/config") return json({ config: baseConfig });
      if (url === "/api/snapshots") return json({ snapshots: [] });
      if (url === "/api/setup/preview")
        return json({
          url: "http://caddy:2019",
          servers: ["srv0"],
          selected_server: "srv0",
          route_count: 1,
          editable: 0,
          readonly: 1,
          caddy_empty: false,
          can_bootstrap: false,
          preview_token: "token-1",
          ownership_notice:
            "Caddy Admin UI will manage only this server's routes array.",
          routes: [
            {
              id: "raw-1",
              domain: "external.example.com",
              handler_type: "unknown",
              config: {},
              enabled: true,
              readonly: true,
              position: 0,
              created_at: "",
              updated_at: "",
            },
          ],
        });
      if (url === "/api/setup/confirm" && init?.method === "POST")
        return json({
          config: {
            ...baseConfig,
            setup_complete: true,
            managed_server: "srv0",
            last_etag: "etag-1",
          },
          imported: 1,
          message: "ok",
        });
      return json({});
    });
    vi.stubGlobal("fetch", fetchMock);

    render(<Settings />);
    await screen.findByText("Review Caddy Servers");
    fireEvent.click(screen.getByText("Review Caddy Servers"));
    await screen.findByText("Ownership Review");
    expect(screen.getByText("external.example.com")).toBeInTheDocument();
    expect(screen.getByText("preserved read-only")).toBeInTheDocument();
    expect(screen.getByText(/ETags/)).toBeInTheDocument();

    fireEvent.click(screen.getByText("Confirm Ownership"));
    await waitFor(() =>
      expect(screen.getByText(/Connected to srv0/)).toBeInTheDocument(),
    );
    const confirmCall = fetchMock.mock.calls.find(
      ([url]) => url === "/api/setup/confirm",
    );
    expect(JSON.parse(confirmCall?.[1]?.body as string)).toEqual({
      url: "http://caddy:2019",
      server: "srv0",
      preview_token: "token-1",
    });
  });

  it("shows snapshots with export and guarded restore actions", async () => {
    vi.stubGlobal(
      "confirm",
      vi.fn(() => true),
    );
    vi.stubGlobal(
      "fetch",
      vi.fn((url: string) => {
        if (url === "/api/config")
          return json({
            config: {
              ...baseConfig,
              setup_complete: true,
              managed_server: "srv0",
            },
          });
        if (url === "/api/snapshots")
          return json({
            snapshots: [
              {
                id: "snap-1",
                server: "srv0",
                etag: "e1",
                reason: "before route update",
                created_at: "2026-01-01T00:00:00Z",
              },
            ],
          });
        if (url === "/api/snapshots/snap-1/restore")
          return json({ message: "snapshot restored", restored: 2 });
        return json({});
      }),
    );

    render(<Settings />);
    await screen.findByText("before route update");
    expect(screen.getByText("Export")).toHaveAttribute(
      "href",
      "/api/snapshots/snap-1/export",
    );
    fireEvent.click(screen.getByText("Restore"));
    await screen.findByText(/Snapshot restored/);
  });

  it("clears stale ownership preview when refresh fails", async () => {
    let previews = 0;
    vi.stubGlobal(
      "fetch",
      vi.fn((url: string) => {
        if (url === "/api/config") return json({ config: baseConfig });
        if (url === "/api/snapshots") return json({ snapshots: [] });
        if (url === "/api/setup/preview") {
          previews++;
          if (previews === 1)
            return json({
              url: "http://caddy:2019",
              servers: ["srv0"],
              selected_server: "srv0",
              route_count: 0,
              editable: 0,
              readonly: 0,
              caddy_empty: false,
              can_bootstrap: false,
              preview_token: "token",
              ownership_notice: "Scoped ownership.",
              routes: [],
            });
          return json({ error: "offline" }, false, 502);
        }
        return json({});
      }),
    );

    render(<Settings />);
    await screen.findByText("Review Caddy Servers");
    fireEvent.click(screen.getByText("Review Caddy Servers"));
    await screen.findByText("Ownership Review");
    fireEvent.click(screen.getByText("Review Caddy Servers"));
    await screen.findByText(/Setup preview failed/);
    expect(screen.queryByText("Ownership Review")).not.toBeInTheDocument();
  });
});
