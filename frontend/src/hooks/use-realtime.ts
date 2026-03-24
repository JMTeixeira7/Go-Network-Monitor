/**
 * Real-time data hooks using Server-Sent Events (SSE).
 *
 * These hooks are designed for data that the backend pushes in real-time,
 * such as security alerts and visited domain updates.
 *
 * ── How it works ─────────────────────────────────────────────────────
 * Each hook pairs an SSE connection with a TanStack Query cache entry.
 * When the backend sends a message, the query cache is updated instantly,
 * which re-renders any component consuming that data.
 *
 * ── Current state ────────────────────────────────────────────────────
 * SSE is disabled by default (ENABLE_SSE = false). Data is fetched via
 * normal HTTP queries as a fallback.
 *
 * ── To connect the real backend ──────────────────────────────────────
 *   1. Set ENABLE_SSE = true below.
 *   2. Ensure SSE_ENDPOINTS match your Go server's streaming routes.
 *   3. Each SSE message should be a JSON-encoded payload matching the
 *      corresponding TypeScript type (VisitedHost[], SecurityAlert[]).
 *   4. The hooks will start receiving push updates automatically.
 *
 * ── Expected SSE message format ──────────────────────────────────────
 *   data: [{"domain":"example.com","lastVisited":"...","visitCount":5}]
 *
 * ── Reconnection ─────────────────────────────────────────────────────
 * On connection error, the hook waits 5 seconds then reconnects.
 */

import { useEffect, useRef, useCallback } from "react";
import { useQueryClient, useQuery } from "@tanstack/react-query";
import { queryKeys } from "./use-api";
import * as api from "@/services/api/client";

// ── Feature flag — flip to true when the Go backend exposes SSE ─────
const ENABLE_SSE = false;

// ── SSE endpoint configuration ──────────────────────────────────────
const SSE_ENDPOINTS = {
  /** Streams VisitedHost[] updates */
  alerts: "/api/stream/alerts",
  /** Streams SecurityAlert[] updates */
  visitedHosts: "/api/stream/visited-hosts",
} as const;

type SSEChannel = keyof typeof SSE_ENDPOINTS;

// ── Generic SSE Hook ────────────────────────────────────────────────
/**
 * Connects to an SSE endpoint and updates the corresponding TanStack
 * Query cache entry whenever a message is received.
 *
 * No-op when ENABLE_SSE is false.
 */
export function useSSE<T>(channel: SSEChannel, queryKey: readonly string[]) {
  const qc = useQueryClient();
  const sourceRef = useRef<EventSource | null>(null);

  const connect = useCallback(() => {
    if (!ENABLE_SSE) return;

    const url = SSE_ENDPOINTS[channel];
    const es = new EventSource(url);

    es.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as T;
        qc.setQueryData(queryKey, data);
      } catch {
        console.error(`[SSE:${channel}] Failed to parse message`);
      }
    };

    es.onerror = () => {
      es.close();
      setTimeout(connect, 5000);
    };

    sourceRef.current = es;
  }, [channel, qc, queryKey]);

  useEffect(() => {
    connect();
    return () => sourceRef.current?.close();
  }, [connect]);
}

// ── Domain-Specific Real-Time Hooks ─────────────────────────────────

/**
 * Visited hosts — receives push updates via SSE when enabled,
 * otherwise falls back to a standard HTTP query.
 */
export function useRealtimeVisitedHosts() {
  useSSE("visitedHosts", queryKeys.visitedHosts);
  return useQuery({
    queryKey: queryKeys.visitedHosts,
    queryFn: api.fetchVisitedHosts,
    staleTime: ENABLE_SSE ? 5 * 60 * 1000 : 30_000,
  });
}

/**
 * Security alerts — receives push updates via SSE when enabled,
 * otherwise falls back to polling via the dashboard summary query.
 */
export function useRealtimeAlerts() {
  useSSE("alerts", [...queryKeys.dashboard, "alerts"] as unknown as readonly string[]);
  return useQuery({
    queryKey: queryKeys.dashboard,
    queryFn: api.fetchDashboardSummary,
    staleTime: ENABLE_SSE ? 5 * 60 * 1000 : 30_000,
    select: (data) => data.alerts,
  });
}
