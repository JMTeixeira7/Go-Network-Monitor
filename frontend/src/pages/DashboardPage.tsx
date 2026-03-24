/**
 * Dashboard Page
 *
 * Main overview page showing proxy status, blocked domain count,
 * security alerts, recent activity, and visited domains.
 *
 * Data sources:
 *  - useDashboardSummary() → GET /api/dashboard (polled via staleTime)
 *  - useStartListener()    → POST /api/listener/start
 *
 * @see src/services/api/client.ts for endpoint mapping
 */

import { PageHeader } from "@/components/layout/PageHeader";
import { StatCard } from "@/components/shared/StatCard";
import { StatusBadge } from "@/components/shared/StatusBadge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useDashboardSummary, useStartListener } from "@/hooks/use-api";
import {
  ShieldBan,
  Activity,
  AlertTriangle,
  Globe,
  Play,
} from "lucide-react";
import { useNavigate } from "react-router-dom";
import { format } from "date-fns";
import { toast } from "sonner";

export default function DashboardPage() {
  const { data, isLoading } = useDashboardSummary();
  const navigate = useNavigate();
  const startMutation = useStartListener();

  const handleStartProxy = () => {
    startMutation.mutate(undefined, {
      onSuccess: () => toast.success("Proxy started"),
      onError: () => toast.error("Failed to start proxy"),
    });
  };

  if (isLoading || !data) {
    return (
      <div className="space-y-6">
        <PageHeader title="Dashboard" description="Network monitoring overview" />
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-32" />
          ))}
        </div>
      </div>
    );
  }

  const listenerRunning = data.status.listenerRunning;

  return (
    <div className="space-y-6">
      <PageHeader title="Dashboard" description="Network monitoring overview" />

      {/* Stats row */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <StatCard
          title="Listener Status"
          value={
            listenerRunning ? (
              <StatusBadge status="running" />
            ) : (
              <Button
                size="sm"
                onClick={handleStartProxy}
                disabled={startMutation.isPending}
                className="bg-success text-success-foreground hover:bg-success/90"
              >
                <Play className="h-3.5 w-3.5 mr-1" />
                {startMutation.isPending ? "Starting…" : "Start Proxy"}
              </Button>
            )
          }
          icon={<Activity className="h-4 w-4" />}
          description={listenerRunning ? `Updated ${format(new Date(data.status.lastUpdated), "HH:mm:ss")}` : "Proxy is not running"}
        />
        <StatCard
          title="Blocked Domains"
          value={data.totalBlocked}
          icon={<ShieldBan className="h-4 w-4" />}
        />
        <StatCard
          title="Security Alerts"
          value={data.alerts.length}
          icon={<AlertTriangle className="h-4 w-4" />}
          description="Active alerts"
        />
      </div>

      <div className="grid gap-4 lg:grid-cols-2">
        {/* Recent Activity */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="text-base">Recent Activity</CardTitle>
            <Button variant="ghost" size="sm" onClick={() => navigate("/activity")}>
              View all
            </Button>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {data.recentActivity.map((item) => (
                <div key={item.id} className="flex items-center justify-between text-sm">
                  <div className="flex items-center gap-2 min-w-0">
                    <Globe className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
                    <span className="font-mono text-xs truncate">{item.domain}</span>
                  </div>
                  <div className="flex items-center gap-2 shrink-0">
                    <StatusBadge status={item.action} />
                    <span className="text-xs text-muted-foreground">
                      {format(new Date(item.timestamp), "HH:mm")}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Visited Domains Summary */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="text-base flex items-center gap-2">
              <Globe className="h-4 w-4" /> Visited Domains
            </CardTitle>
            <Button variant="ghost" size="sm" onClick={() => navigate("/visited-hosts")}>
              View all
            </Button>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {data.visitedHosts.map((host) => (
                <div key={host.domain} className="flex items-center justify-between text-sm">
                  <span className="font-mono text-xs truncate">{host.domain}</span>
                  <div className="flex items-center gap-3 shrink-0">
                    <span className="text-xs text-muted-foreground">{host.visitCount} visits</span>
                    <span className="text-xs text-muted-foreground">
                      {format(new Date(host.lastVisited), "HH:mm")}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Security Alerts */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Security Alerts</CardTitle>
          </CardHeader>
          <CardContent>
            {data.alerts.length === 0 ? (
              <p className="text-sm text-muted-foreground">No active alerts</p>
            ) : (
              <div className="space-y-3">
                {data.alerts.map((alert) => (
                  <div key={alert.id} className="flex items-start gap-3 text-sm">
                    <AlertTriangle className="h-4 w-4 text-warning shrink-0 mt-0.5" />
                    <div className="min-w-0">
                      <p className="font-medium truncate">{alert.domain}</p>
                      <p className="text-xs text-muted-foreground">{alert.message}</p>
                    </div>
                    <StatusBadge status="flagged" label={alert.type} className="shrink-0" />
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
