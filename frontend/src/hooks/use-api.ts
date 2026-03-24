/**
 * TanStack Query hooks — the bridge between UI and API.
 *
 * Every component that needs server data uses one of these hooks.
 * They wrap the functions from src/services/api/client.ts and manage
 * caching, refetching, and cache invalidation automatically.
 *
 * ── Data flow patterns ──────────────────────────────────────────────
 *
 *   Pattern        Hook example          When to use
 *   ─────────────  ────────────────────  ──────────────────────────
 *   Request/Reply  useBlockedDomains()   Standard CRUD pages
 *   Polling        useSystemStatus()     Proxy status (every 15s)
 *   SSE (push)     useRealtimeVisitedHosts()  Live data (see use-realtime.ts)
 *   Mutation       useAddBlockedDomain() Any write operation
 *
 * ── Adding a new endpoint ───────────────────────────────────────────
 *   1. Add the fetch function to src/services/api/client.ts
 *   2. Add a query key to queryKeys below
 *   3. Create a useXxx hook that calls useQuery/useMutation
 *   4. Use the hook in your component
 */

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import * as api from "@/services/api/client";
import type { AddBlockedDomainPayload } from "@/types";

// ── Query keys ──────────────────────────────────────────────────────
// Centralised keys prevent typos and enable targeted invalidation.
export const queryKeys = {
  dashboard: ["dashboard"] as const,
  blockedDomains: ["blocked-domains"] as const,
  blockedDomain: (domain: string) => ["blocked-domains", domain] as const,
  activity: ["activity"] as const,
  systemStatus: ["system-status"] as const,
  visitedHosts: ["visited-hosts"] as const,
  proxyServices: ["proxy-services"] as const,
};

// ── Queries ─────────────────────────────────────────────────────────

/** Dashboard summary — fetched on page load, cached for 30s. */
export const useDashboardSummary = () =>
  useQuery({ queryKey: queryKeys.dashboard, queryFn: api.fetchDashboardSummary });

/** All blocked domains — used by the Blocked Domains table. */
export const useBlockedDomains = () =>
  useQuery({ queryKey: queryKeys.blockedDomains, queryFn: api.fetchBlockedDomains });

/** Single blocked domain with schedules — used by Domain Details page. */
export const useBlockedDomain = (domain: string) =>
  useQuery({ queryKey: queryKeys.blockedDomain(domain), queryFn: () => api.fetchBlockedDomain(domain), enabled: !!domain });

/** Activity log — request/response pattern. */
export const useActivity = () =>
  useQuery({ queryKey: queryKeys.activity, queryFn: api.fetchActivity });

/** Proxy status — POLLED every 15 seconds for near-real-time status. */
export const useSystemStatus = () =>
  useQuery({ queryKey: queryKeys.systemStatus, queryFn: api.fetchSystemStatus, refetchInterval: 15000 });

/** Visited hosts — standard fetch (see use-realtime.ts for SSE variant). */
export const useVisitedHosts = () =>
  useQuery({ queryKey: queryKeys.visitedHosts, queryFn: api.fetchVisitedHosts });

/** Proxy services list with active/inactive status. */
export const useProxyServices = () =>
  useQuery({ queryKey: queryKeys.proxyServices, queryFn: api.fetchProxyServices });

// ── Mutations ───────────────────────────────────────────────────────

/** Add a new blocked domain (with optional schedules). */
export const useAddBlockedDomain = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: AddBlockedDomainPayload) => api.addBlockedDomain(payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.blockedDomains });
      qc.invalidateQueries({ queryKey: queryKeys.dashboard });
    },
  });
};

/** Delete / unblock a domain. */
export const useDeleteBlockedDomain = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (domain: string) => api.deleteBlockedDomain(domain),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.blockedDomains });
      qc.invalidateQueries({ queryKey: queryKeys.dashboard });
    },
  });
};

/** Start the proxy listener. */
export const useStartListener = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: api.startListener,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.systemStatus });
      qc.invalidateQueries({ queryKey: queryKeys.dashboard });
    },
  });
};

/** Stop the proxy listener. */
export const useStopListener = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: api.stopListener,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.systemStatus });
      qc.invalidateQueries({ queryKey: queryKeys.dashboard });
    },
  });
};

/** Clear the proxy session cache. */
export const useClearCache = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: api.clearCache,
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.systemStatus }),
  });
};

/** Toggle a proxy service on/off. */
export const useToggleProxyService = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (serviceId: string) => api.toggleProxyService(serviceId),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.proxyServices }),
  });
};
