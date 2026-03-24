export interface ScheduleRule {
  id: string;
  startTime: string;
  endTime: string;
  weekday: string;
}

export interface BlockedDomain {
  domain: string;
  schedulesCount: number;
  createdAt: string;
  schedules?: ScheduleRule[];
}

export interface SystemStatus {
  listenerRunning: boolean;
  cacheStatus: "active" | "cleared" | "unknown";
  lastUpdated: string;
}

export interface ActivityItem {
  id: string;
  domain: string;
  timestamp: string;
  action: "blocked" | "allowed" | "flagged" | "visited";
  reasons: string[];
}

export interface SecurityAlert {
  id: string;
  domain: string;
  type: "phishing" | "typosquatting" | "xss" | "suspicious";
  message: string;
  timestamp: string;
}

export interface VisitedHost {
  domain: string;
  lastVisited: string;
  visitCount: number;
}

export interface ProxyService {
  id: string;
  name: string;
  description: string;
  active: boolean;
}

export interface DashboardSummary {
  status: SystemStatus;
  totalBlocked: number;
  recentActivity: ActivityItem[];
  alerts: SecurityAlert[];
  visitedHosts: VisitedHost[];
}

export interface ApiResponse<T> {
  data: T;
  success: boolean;
  message?: string;
}

export interface AddBlockedDomainPayload {
  domain: string;
  schedules: Omit<ScheduleRule, "id">[];
}
