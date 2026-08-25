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

describe("Settings connection flow", () => {
  beforeEach(() => vi.restoreAllMocks());

  it("reviews routes and confirms the content-bound connection", async () => {
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
    await screen.findByText("Preview Routes");
    fireEvent.click(screen.getByText("Preview Routes"));
    await screen.findByText("Review Connection");
    expect(screen.getByText(/No changes have been made/)).toBeInTheDocument();
    expect(
      screen.getByText(
        /Save the current route array as the initial restore point/,
      ),
    ).toBeInTheDocument();
    expect(screen.getByText("external.example.com")).toBeInTheDocument();
    expect(screen.getAllByText("Preserved read-only")).toHaveLength(2);
    expect(screen.getByText(/ETags/)).toBeInTheDocument();

    fireEvent.click(screen.getByText("Connect and Import"));
    await waitFor(() =>
      expect(
        screen.getByText(
          /Connected to srv0.*Imported 1 route.*initial restore point/,
        ),
      ).toBeInTheDocument(),
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
    const fetchMock = vi.fn((url: string) => {
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
              reason: "initial setup",
              created_at: "2026-01-01T00:00:00Z",
            },
          ],
        });
      if (url === "/api/snapshots/snap-1/restore")
        return json({ message: "snapshot restored", restored: 2 });
      return json({});
    });
    vi.stubGlobal("fetch", fetchMock);

    render(<Settings />);
    await screen.findByText("Initial Caddy baseline");
    expect(screen.getByText("Export")).toHaveAttribute(
      "href",
      "/api/snapshots/snap-1/export",
    );
    fireEvent.click(screen.getByText("Restore..."));
    await screen.findByRole("dialog", { name: "Restore Caddy routes?" });
    expect(
      screen.getByText(/current live routes are backed up first/),
    ).toBeInTheDocument();
    expect(
      fetchMock.mock.calls.some(
        ([url]) => url === "/api/snapshots/snap-1/restore",
      ),
    ).toBe(false);

    fireEvent.click(screen.getByText("Restore Live Routes"));
    await screen.findByText(/Restore complete.*2 routes now active/);
  });

  it("clears a stale connection preview when refresh fails", async () => {
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
    await screen.findByText("Preview Routes");
    fireEvent.click(screen.getByText("Preview Routes"));
    await screen.findByText("Review Connection");
    fireEvent.click(screen.getByText("Preview Routes"));
    await screen.findByText(/Setup preview failed/);
    expect(screen.queryByText("Review Connection")).not.toBeInTheDocument();
  });
});
