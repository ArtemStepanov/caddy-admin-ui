import { render, screen, fireEvent, waitFor } from "@testing-library/preact";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { RouteForm } from "./RouteForm";

function json(data: unknown, ok = true, status = 200) {
  return Promise.resolve({
    ok,
    status,
    json: () => Promise.resolve(data),
  } as Response);
}

function lastPostBody(calls: { url: string; init?: RequestInit }[]): {
  handler_type: string;
  strip_path_prefix?: string;
  headers?: unknown;
  config: Record<string, unknown>;
} {
  const post = [...calls]
    .reverse()
    .find((c) => c.url === "/api/routes" && c.init?.method === "POST");
  if (!post?.init?.body) throw new Error("no POST body captured");
  return JSON.parse(post.init.body as string) as {
    handler_type: string;
    strip_path_prefix?: string;
    headers?: unknown;
    config: Record<string, unknown>;
  };
}

describe("RouteForm kind isolation", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("submits reverse_proxy config with strip_path_prefix and route-level headers", async () => {
    const calls: { url: string; init?: RequestInit }[] = [];
    vi.stubGlobal(
      "fetch",
      vi.fn((url: string, init?: RequestInit) => {
        calls.push({ url, init });
        if (url === "/api/routes" && init?.method === "POST")
          return json({
            route: {
              id: "1",
              domain: "proxy.example.com",
              handler_type: "reverse_proxy",
              config: {},
            },
          });
        return json({});
      }),
    );

    render(<RouteForm />);
    await screen.findByText("New Route");

    fireEvent.input(screen.getByPlaceholderText("example.com"), {
      target: { value: "proxy.example.com" },
    });
    fireEvent.input(screen.getByPlaceholderText("/*"), {
      target: { value: "/api/*" },
    });
    fireEvent.input(screen.getByPlaceholderText("/api"), {
      target: { value: "/api" },
    });
    fireEvent.input(screen.getByPlaceholderText("localhost:8080"), {
      target: { value: "localhost:8080" },
    });
    fireEvent.click(screen.getAllByText("Add")[0]);

    fireEvent.click(screen.getByText("Create Route"));

    await waitFor(() => {
      const body = lastPostBody(calls);
      expect(body.handler_type).toBe("reverse_proxy");
      expect(body.strip_path_prefix).toBe("/api");
      expect(body.config.upstreams).toContain("localhost:8080");
      expect(body).toHaveProperty("headers");
    });
  });

  it("switching to file_server drops strip_path_prefix and proxy-only config", async () => {
    const calls: { url: string; init?: RequestInit }[] = [];
    vi.stubGlobal(
      "fetch",
      vi.fn((url: string, init?: RequestInit) => {
        calls.push({ url, init });
        if (url === "/api/routes" && init?.method === "POST")
          return json({
            route: {
              id: "2",
              domain: "files.example.com",
              handler_type: "file_server",
              config: {},
            },
          });
        return json({});
      }),
    );

    render(<RouteForm />);
    await screen.findByText("New Route");

    fireEvent.input(screen.getByPlaceholderText("example.com"), {
      target: { value: "files.example.com" },
    });
    fireEvent.input(screen.getByPlaceholderText("/*"), {
      target: { value: "/api/*" },
    });
    fireEvent.input(screen.getByPlaceholderText("/api"), {
      target: { value: "/api" },
    });
    fireEvent.input(screen.getByPlaceholderText("localhost:8080"), {
      target: { value: "localhost:8080" },
    });
    fireEvent.click(screen.getAllByText("Add")[0]);

    // Switch kind to file_server — must discard proxy-only settings
    fireEvent.click(screen.getByText("File Server"));
    fireEvent.input(screen.getByPlaceholderText("/var/www/html"), {
      target: { value: "/var/www" },
    });

    fireEvent.click(screen.getByText("Create Route"));

    await waitFor(() => {
      const body = lastPostBody(calls);
      expect(body.handler_type).toBe("file_server");
      expect(body).not.toHaveProperty("strip_path_prefix");
      expect(body.config).not.toHaveProperty("upstreams");
      expect(body.config.root).toBe("/var/www");
      expect(body).toHaveProperty("headers");
    });
  });

  it("switching to redir drops strip_path_prefix and proxy-only config", async () => {
    const calls: { url: string; init?: RequestInit }[] = [];
    vi.stubGlobal(
      "fetch",
      vi.fn((url: string, init?: RequestInit) => {
        calls.push({ url, init });
        if (url === "/api/routes" && init?.method === "POST")
          return json({
            route: {
              id: "3",
              domain: "redir.example.com",
              handler_type: "redir",
              config: {},
            },
          });
        return json({});
      }),
    );

    render(<RouteForm />);
    await screen.findByText("New Route");

    fireEvent.input(screen.getByPlaceholderText("example.com"), {
      target: { value: "redir.example.com" },
    });
    fireEvent.input(screen.getByPlaceholderText("/*"), {
      target: { value: "/api/*" },
    });
    fireEvent.input(screen.getByPlaceholderText("/api"), {
      target: { value: "/api" },
    });
    fireEvent.input(screen.getByPlaceholderText("localhost:8080"), {
      target: { value: "localhost:8080" },
    });
    fireEvent.click(screen.getAllByText("Add")[0]);

    fireEvent.click(screen.getByText("Redirect"));
    fireEvent.input(screen.getByPlaceholderText("https://example.com{uri}"), {
      target: { value: "https://target.example.com{uri}" },
    });

    fireEvent.click(screen.getByText("Create Route"));

    await waitFor(() => {
      const body = lastPostBody(calls);
      expect(body.handler_type).toBe("redir");
      expect(body).not.toHaveProperty("strip_path_prefix");
      expect(body.config).not.toHaveProperty("upstreams");
      expect(body.config.to).toBe("https://target.example.com{uri}");
      expect(body).toHaveProperty("headers");
    });
  });

  it("renders a preserved read-only route outside the editable form", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn((url: string) => {
        if (url === "/api/routes/r1")
          return json({
            route: {
              id: "r1",
              domain: "raw.example.com",
              handler_type: "unknown",
              config: {},
              enabled: true,
              readonly: true,
              support_status: "unsupported_readonly",
              readonly_reason: "unknown handler",
              created_at: "2026-01-01T00:00:00Z",
              updated_at: "2026-01-01T00:00:00Z",
            },
          });
        return json({});
      }),
    );

    render(<RouteForm id="r1" />);
    await screen.findByText("Read-only Route");
    expect(screen.getByText(/managed outside the UI/)).toBeInTheDocument();
    expect(screen.queryByText("Create Route")).not.toBeInTheDocument();
    expect(screen.queryByText("Update Route")).not.toBeInTheDocument();
  });
});
