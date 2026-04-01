# NetGuard — Proxy Admin Dashboard

A React + TypeScript frontend for managing a local web proxy / network monitoring backend written in Go.

## Tech Stack

- **React 18** + **TypeScript** — UI framework
- **Vite** — build tool
- **Tailwind CSS** — utility-first styling with semantic design tokens
- **shadcn/ui** — component library (Radix primitives)
- **React Router** — client-side routing
- **TanStack Query** — server state, caching, and mutations

## Project Structure

```
src/
├── components/
│   ├── layout/          # AppLayout, AppSidebar, PageHeader, ThemeToggle
│   ├── shared/          # Reusable: StatCard, StatusBadge, ConfirmDialog, EmptyState
│   └── ui/              # shadcn/ui primitives (do not edit directly)
├── features/
│   └── blocked-domains/ # AddDomainDialog
├── hooks/
│   ├── use-api.ts       # TanStack Query hooks — all data fetching
│   ├── use-realtime.ts  # SSE hooks for live data (alerts, visited hosts)
│   └── use-theme.ts     # Dark/light mode
├── pages/               # Route-level components
├── services/
│   └── api/
│       ├── client.ts    # ★ CENTRALIZED API LAYER — edit this to connect backend
│       └── mock-data.ts # Mock data (remove after backend integration)
├── types/
│   └── index.ts         # All TypeScript interfaces
└── lib/
    └── utils.ts         # Tailwind merge helper
```

## Data Flow Patterns

The app uses three data flow patterns, each suited to different use cases:

| Pattern | Used for | Implementation |
|---------|----------|----------------|
| **Request/Response** | CRUD pages (blocked domains, activity) | Standard TanStack Query in `use-api.ts` |
| **Polling** | Proxy status (is the listener running?) | `useSystemStatus()` with `refetchInterval: 15000` |
| **Server Push (SSE)** | Live alerts, visited domain updates | `useSSE()` hook in `use-realtime.ts` |

## Connecting the Go Backend

### Step 1 — Set the API base URL

Open `src/services/api/client.ts` and change:

```ts
const API_BASE_URL = "/api";
// → e.g. "http://localhost:8080/api"
```

### Step 2 — Replace mock function bodies

Each exported function has a commented-out real implementation. Example:

```ts
// Before (mock):
export async function fetchBlockedDomains(): Promise<BlockedDomain[]> {
  await delay();
  return [...mockBlockedDomains];
}

// After (real):
export async function fetchBlockedDomains(): Promise<BlockedDomain[]> {
  return _get<BlockedDomain[]>("/blocked-domains");
}
```

The `_get`, `_post`, and `_del` helpers are already implemented and ready to use.

### Step 3 — Enable SSE (optional)

Open `src/hooks/use-realtime.ts` and set:

```ts
const ENABLE_SSE = true;
```

Ensure your Go backend exposes:
- `GET /api/stream/alerts` — streams `SecurityAlert[]` as SSE
- `GET /api/stream/visited-hosts` — streams `VisitedHost[]` as SSE

### Step 4 — Remove mock data

Once all endpoints are connected, delete `src/services/api/mock-data.ts` and remove its import from `client.ts`.

## Expected API Endpoints

| Method | Endpoint | Request Body | Response |
|--------|----------|-------------|----------|
| GET | `/api/dashboard` | — | `DashboardSummary` |
| GET | `/api/status` | — | `SystemStatus` |
| GET | `/api/blocked-domains` | — | `BlockedDomain[]` |
| GET | `/api/blocked-domains/:domain` | — | `BlockedDomain` |
| POST | `/api/blocked-domains` | `AddBlockedDomainPayload` | `ApiResponse<BlockedDomain>` |
| DELETE | `/api/blocked-domains/:domain` | — | `ApiResponse<null>` |
| GET | `/api/activity` | — | `ActivityItem[]` |
| POST | `/api/listener/start` | — | `ApiResponse<null>` |
| POST | `/api/listener/stop` | — | `ApiResponse<null>` |
| POST | `/api/cache/clear` | — | `ApiResponse<null>` |
| GET | `/api/visited-hosts` | — | `VisitedHost[]` |
| GET | `/api/proxy-services` | — | `ProxyService[]` |
| POST | `/api/proxy-services/:id/toggle` | — | `ApiResponse<ProxyService>` |

See `src/types/index.ts` for all TypeScript interfaces.

## Running Locally

web dev
```bash
npm install
npm run dev
```

local app
```bash
npm install
npm run electron-package
```
