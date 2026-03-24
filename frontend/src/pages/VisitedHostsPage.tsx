import { PageHeader } from "@/components/layout/PageHeader";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useRealtimeVisitedHosts } from "@/hooks/use-realtime";
import { Globe } from "lucide-react";
import { format } from "date-fns";

export default function VisitedHostsPage() {
  const { data: hosts, isLoading } = useRealtimeVisitedHosts();

  if (isLoading || !hosts) {
    return (
      <div className="space-y-6">
        <PageHeader title="Visited Hosts" description="All visited domains" />
        <Skeleton className="h-64" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader title="Visited Hosts" description="All domains visited through the proxy" />
      <Card>
        <CardContent className="pt-6">
          <div className="space-y-2">
            {hosts.map((host) => (
              <div
                key={host.domain}
                className="flex items-center justify-between rounded-md border border-border px-4 py-3 text-sm"
              >
                <div className="flex items-center gap-2 min-w-0">
                  <Globe className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
                  <span className="font-mono text-xs truncate">{host.domain}</span>
                </div>
                <div className="flex items-center gap-4 shrink-0">
                  <span className="text-xs text-muted-foreground">{host.visitCount} visits</span>
                  <span className="text-xs text-muted-foreground">
                    {format(new Date(host.lastVisited), "MMM d, HH:mm")}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
