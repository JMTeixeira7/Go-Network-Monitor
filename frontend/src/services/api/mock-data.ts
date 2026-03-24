import type {
  BlockedDomain,
  ScheduleRule,
  SystemStatus,
  ActivityItem,
  SecurityAlert,
  DashboardSummary,
  VisitedHost,
  ProxyService,
} from "@/types";

const schedules: ScheduleRule[] = [
  { id: "s1", startTime: "09:00", endTime: "17:00", weekday: "Monday" },
  { id: "s2", startTime: "09:00", endTime: "17:00", weekday: "Tuesday" },
  { id: "s3", startTime: "00:00", endTime: "23:59", weekday: "Wednesday" },
];

export const mockBlockedDomains: BlockedDomain[] = [
  { domain: "malware-site.com", schedulesCount: 2, createdAt: "2025-03-10T08:30:00Z", schedules: [schedules[0], schedules[1]] },
  { domain: "phishing-example.net", schedulesCount: 0, createdAt: "2025-03-12T14:00:00Z", schedules: [] },
  { domain: "ad-tracker.io", schedulesCount: 1, createdAt: "2025-03-15T10:15:00Z", schedules: [schedules[2]] },
  { domain: "typosquat-g00gle.com", schedulesCount: 0, createdAt: "2025-03-18T09:00:00Z", schedules: [] },
  { domain: "suspicious-redirect.xyz", schedulesCount: 1, createdAt: "2025-03-20T16:45:00Z", schedules: [schedules[0]] },
  { domain: "crypto-scam.org", schedulesCount: 0, createdAt: "2025-03-21T11:20:00Z", schedules: [] },
];

export const mockStatus: SystemStatus = {
  listenerRunning: true,
  cacheStatus: "active",
  lastUpdated: new Date().toISOString(),
};

export const mockActivity: ActivityItem[] = [
  { id: "a1", domain: "malware-site.com", timestamp: "2025-03-23T10:05:00Z", action: "blocked", reasons: ["Known malware distributor"] },
  { id: "a2", domain: "example.com", timestamp: "2025-03-23T10:04:00Z", action: "allowed", reasons: [] },
  { id: "a3", domain: "phishing-example.net", timestamp: "2025-03-23T10:02:00Z", action: "blocked", reasons: ["Phishing detected"] },
  { id: "a4", domain: "typosquat-g00gle.com", timestamp: "2025-03-23T09:58:00Z", action: "flagged", reasons: ["Typosquatting: google.com"] },
  { id: "a5", domain: "news.ycombinator.com", timestamp: "2025-03-23T09:50:00Z", action: "visited", reasons: [] },
  { id: "a6", domain: "suspicious-redirect.xyz", timestamp: "2025-03-23T09:45:00Z", action: "blocked", reasons: ["Suspicious redirect chain", "Low reputation score"] },
  { id: "a7", domain: "github.com", timestamp: "2025-03-23T09:30:00Z", action: "visited", reasons: [] },
  { id: "a8", domain: "crypto-scam.org", timestamp: "2025-03-23T09:15:00Z", action: "blocked", reasons: ["Scam site"] },
];

export const mockAlerts: SecurityAlert[] = [
  { id: "al1", domain: "typosquat-g00gle.com", type: "typosquatting", message: "Possible typosquatting of google.com", timestamp: "2025-03-23T09:58:00Z" },
  { id: "al2", domain: "phishing-example.net", type: "phishing", message: "SSL certificate mismatch detected", timestamp: "2025-03-23T10:02:00Z" },
  { id: "al3", domain: "suspicious-redirect.xyz", type: "suspicious", message: "Multiple redirect hops detected", timestamp: "2025-03-23T09:45:00Z" },
];

export const mockVisitedHosts: VisitedHost[] = [
  { domain: "github.com", lastVisited: "2025-03-23T10:05:00Z", visitCount: 42 },
  { domain: "stackoverflow.com", lastVisited: "2025-03-23T09:58:00Z", visitCount: 31 },
  { domain: "news.ycombinator.com", lastVisited: "2025-03-23T09:50:00Z", visitCount: 18 },
  { domain: "google.com", lastVisited: "2025-03-23T09:45:00Z", visitCount: 97 },
  { domain: "reddit.com", lastVisited: "2025-03-23T09:30:00Z", visitCount: 24 },
  { domain: "docs.microsoft.com", lastVisited: "2025-03-23T09:15:00Z", visitCount: 12 },
  { domain: "npmjs.com", lastVisited: "2025-03-23T09:00:00Z", visitCount: 8 },
  { domain: "developer.mozilla.org", lastVisited: "2025-03-23T08:45:00Z", visitCount: 15 },
  { domain: "cloudflare.com", lastVisited: "2025-03-23T08:30:00Z", visitCount: 5 },
  { domain: "vercel.com", lastVisited: "2025-03-23T08:15:00Z", visitCount: 3 },
  { domain: "docker.com", lastVisited: "2025-03-23T08:00:00Z", visitCount: 7 },
  { domain: "aws.amazon.com", lastVisited: "2025-03-23T07:45:00Z", visitCount: 11 },
];

export const mockProxyServices: ProxyService[] = [
  { id: "ps1", name: "Phishing Detection", description: "Detects and blocks known phishing domains", active: true },
  { id: "ps2", name: "Typosquatting Guard", description: "Identifies domains mimicking popular sites", active: true },
  { id: "ps3", name: "XSS Prevention", description: "Scans and blocks cross-site scripting attempts", active: false },
  { id: "ps4", name: "Malware Blocking", description: "Blocks domains associated with malware distribution", active: true },
  { id: "ps5", name: "Ad & Tracker Filtering", description: "Filters advertising and tracking domains", active: true },
  { id: "ps6", name: "DNS Rebinding Protection", description: "Prevents DNS rebinding attacks", active: false },
  { id: "ps7", name: "SSL Certificate Validation", description: "Validates SSL certificates for visited domains", active: true },
];

export const mockDashboardSummary: DashboardSummary = {
  status: mockStatus,
  totalBlocked: mockBlockedDomains.length,
  recentActivity: mockActivity.slice(0, 5),
  alerts: mockAlerts,
  visitedHosts: mockVisitedHosts.slice(0, 10),
};
