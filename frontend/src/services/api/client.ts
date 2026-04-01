/**
 * Centralized API Client
 * ======================
 *
 * All backend communication flows through this single module.
 * Components never call fetch() directly — they use hooks from
 * src/hooks/use-api.ts which delegate to the functions below.
 *
 * ── Current state ────────────────────────────────────────────────────
 * Returns mock data with simulated latency. The real Go backend
 * is not yet connected.
 *
 * ── To connect the Go backend ────────────────────────────────────────
 *
 *   1. Set API_BASE_URL to your Go server address:
 *        const API_BASE_URL = "http://localhost:8080/api";
 *
 *   2. Replace each mock function body with a real fetch call.
 *      Use the provided helpers (_get, _post, _del). Example:
 *
 *        // Before (mock):
 *        export async function fetchBlockedDomains(): Promise<BlockedDomain[]> {
 *          await delay();
 *          return [...mockBlockedDomains];
 *        }
 *
 *        // After (real):
 *        export async function fetchBlockedDomains(): Promise<BlockedDomain[]> {
 *          return _get<BlockedDomain[]>("/blocked-domains");
 *        }
 *
 *   3. Keep the return types unchanged — components rely on them.
 *
 *   4. Remove the mock-data import and the delay() helper once done.
 *
 * ── Endpoint map ─────────────────────────────────────────────────────
 *
 *   Function                 Method   Endpoint
 *   ───────────────────────  ──────   ──────────────────────────────
 *   fetchDashboardSummary    GET      /api/dashboard
 *   fetchBlockedDomains      GET      /api/blocked-domains
 *   fetchBlockedDomain       GET      /api/blocked-domains/:domain
 *   addBlockedDomain         POST     /api/blocked-domains
 *   deleteBlockedDomain      DELETE   /api/blocked-domains/:domain
 *   fetchActivity            GET      /api/activity
 *   fetchSystemStatus        GET      /api/status          (polled every 15s)
 *   startListener            POST     /api/listener/start
 *   stopListener             POST     /api/listener/stop
 *   clearCache               POST     /api/cache/clear
 *   fetchVisitedHosts        GET      /api/visited-hosts   (or via SSE)
 *   fetchProxyServices       GET      /api/proxy-services
 *   toggleProxyService       POST     /api/proxy-services/:id/toggle
 *
 * ── SSE / real-time endpoints (see src/hooks/use-realtime.ts) ───────
 *
 *   /api/stream/alerts          → pushes SecurityAlert[]
 *   /api/stream/visited-hosts   → pushes VisitedHost[]
 */

import type {
  ApiResponse,
  BlockedDomain,
  DashboardSummary,
  ActivityItem,
  SystemStatus,
  AddBlockedDomainPayload,
  VisitedHost,
  ProxyService,
} from "@/types";

import {
  mockBlockedDomains,
  mockDashboardSummary,
  mockActivity,
  mockStatus,
  mockVisitedHosts,
  mockProxyServices,
} from "./mock-data";

// Configuration
// Point this to the Go backend when ready (e.g. "http://localhost:8080/api")
const API_BASE_URL = "http://127.0.0.1:8081/api";

// Simulated network latency (remove when using real backend)
const delay = (ms = 400) => new Promise((r) => setTimeout(r, ms));

// HTTP Helpers
// Ready-to-use fetch wrappers. Uncomment the bodies when connecting.

function extractErrorMessage(value: unknown, fallback = "Request failed"): string {
  if (typeof value === "string" && value.trim()) return value;

  if (typeof value === "object" && value !== null) {
    const obj = value as Record<string, unknown>;

    if (typeof obj.message === "string" && obj.message.trim()) return obj.message;
    if (typeof obj.error === "string" && obj.error.trim()) return obj.error;
  }

  return fallback;
}

function isApiResponse<T>(value: unknown): value is ApiResponse<T> {
  return (
    typeof value === "object" &&
    value !== null &&
    "success" in value &&
    typeof (value as { success?: unknown }).success === "boolean" &&
    "data" in value
  );
}

async function _get<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    method: "GET",
    headers: {
      Accept: "application/json",
    },
  });

  const raw = await res.text();

  let parsed: unknown;
  try {
    parsed = raw ? JSON.parse(raw) : null;
  } catch {
    throw new Error(
      raw
        ? `Expected JSON but received: ${raw.slice(0, 200)}`
        : "Empty response from server"
    );
  }

  if (!res.ok) {
    throw new Error(extractErrorMessage(parsed, `HTTP ${res.status}`));
  }

  if (isApiResponse<T>(parsed)) {
    if (!parsed.success) {
      throw new Error(parsed.message || "Request failed");
    }
    return parsed.data;
  }

  return parsed as T;
}

async function _post<T>(path: string, body?: unknown): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: body ? JSON.stringify(body) : undefined,
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

async function _del<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, { method: "DELETE" });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

// ── Public API (mock implementations) ───────────────────────────────
// Replace each function body with the commented real call when ready.

/** GET /api/dashboard */
export async function fetchDashboardSummary(): Promise<DashboardSummary> {
  await delay();
  // return _get<DashboardSummary>("/dashboard");
  return { ...mockDashboardSummary, status: { ...mockStatus, lastUpdated: new Date().toISOString() } };
}

/** GET /api/blocked-domains */
export async function fetchBlockedDomains(): Promise<BlockedDomain[]> {
  return _get<BlockedDomain[]>("/blocked-domains");
}

/** GET /api/blocked-domains/:domain */
export async function fetchBlockedDomain(domain: string): Promise<BlockedDomain> {
  return _get<BlockedDomain>(`/blocked-domains/${encodeURIComponent(domain)}`);
  //const found = mockBlockedDomains.find((d) => d.domain === domain);
  //if (!found) throw new Error(`Domain ${domain} not found`);
  //return { ...found };
}

/** POST /api/blocked-domains */
export async function addBlockedDomain(payload: AddBlockedDomainPayload): Promise<ApiResponse<BlockedDomain>> {
  await delay(600);
  return _post<ApiResponse<BlockedDomain>>("/blocked-domains", payload);
  //const newDomain: BlockedDomain = {
  //  domain: payload.domain,
  //  schedulesCount: payload.schedules.length,
  //  createdAt: new Date().toISOString(),
  //  schedules: payload.schedules.map((s, i) => ({ ...s, id: `new-${i}` })),
  //};
  //mockBlockedDomains.push(newDomain);
  //return { data: newDomain, success: true, message: "Domain added" };
}

/** DELETE /api/blocked-domains/:domain */
export async function deleteBlockedDomain(domain: string): Promise<ApiResponse<null>> {
  await delay(600);
  // return _del<ApiResponse<null>>(`/blocked-domains/${encodeURIComponent(domain)}`);
  const idx = mockBlockedDomains.findIndex((d) => d.domain === domain);
  if (idx !== -1) mockBlockedDomains.splice(idx, 1);
  return { data: null, success: true, message: "Domain removed" };
}

/** GET /api/activity */
export async function fetchActivity(): Promise<ActivityItem[]> {
  await delay();
  // return _get<ActivityItem[]>("/activity");
  return [...mockActivity];
}

/** GET /api/status — polled every 15 seconds */
export async function fetchSystemStatus(): Promise<SystemStatus> {
  return _get<SystemStatus>("/status");
}

/** POST /api/listener/start */
export async function startListener(): Promise<ApiResponse<null>> {
  await delay(500);
  return _post<ApiResponse<null>>("/listener/start");
}

/** POST /api/listener/stop */
export async function stopListener(): Promise<ApiResponse<null>> {
  await delay(800);
  return _post<ApiResponse<null>>("/listener/stop");
}

/** POST /api/cache/clear */
export async function clearCache(): Promise<ApiResponse<null>> {
  await delay(600);
  return _post<ApiResponse<null>>("/cache/clear");
}

/** GET /api/visited-hosts (also available via SSE at /api/stream/visited-hosts) */
export async function fetchVisitedHosts(): Promise<VisitedHost[]> {
  await delay();
  // return _get<VisitedHost[]>("/visited-hosts");
  return [...mockVisitedHosts];
}

/** GET /api/proxy-services */
export async function fetchProxyServices(): Promise<ProxyService[]> {
  await delay();
  // return _get<ProxyService[]>("/proxy-services");
  return [...mockProxyServices];
}

/** POST /api/proxy-services/:id/toggle */
export async function toggleProxyService(serviceId: string): Promise<ApiResponse<ProxyService>> {
  await delay(400);
  // return _post<ApiResponse<ProxyService>>(`/proxy-services/${serviceId}/toggle`);
  const service = mockProxyServices.find((s) => s.id === serviceId);
  if (!service) throw new Error(`Service ${serviceId} not found`);
  service.active = !service.active;
  return { data: { ...service }, success: true, message: `Service ${service.active ? "activated" : "deactivated"}` };
}

// Re-export helpers for direct use if needed
export { _get, _post, _del };
