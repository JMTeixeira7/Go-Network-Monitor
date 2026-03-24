/**
 * Web Proxy Page
 *
 * Controls for starting/stopping the proxy listener, clearing cache,
 * and managing individual proxy services (activate/deactivate).
 *
 * Data sources:
 *  - useSystemStatus()        → GET /api/status (polled every 15s)
 *  - useProxyServices()       → GET /api/proxy-services
 *  - useToggleProxyService()  → POST /api/proxy-services/:id/toggle
 *  - useStartListener()       → POST /api/listener/start
 *  - useStopListener()        → POST /api/listener/stop
 *  - useClearCache()          → POST /api/cache/clear
 *
 * @see src/services/api/client.ts for endpoint mapping
 */

import { useState } from "react";
import { PageHeader } from "@/components/layout/PageHeader";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import {
  useSystemStatus,
  useStartListener,
  useStopListener,
  useClearCache,
  useProxyServices,
  useToggleProxyService,
} from "@/hooks/use-api";
import { Play, Square, Trash2, Network, ChevronDown, Power, PowerOff } from "lucide-react";
import { toast } from "sonner";

export default function SystemPage() {
  const { data: status, isLoading: statusLoading } = useSystemStatus();
  const { data: services, isLoading: servicesLoading } = useProxyServices();
  const startMutation = useStartListener();
  const stopMutation = useStopListener();
  const cacheMutation = useClearCache();
  const toggleMutation = useToggleProxyService();

  const [openServices, setOpenServices] = useState<Record<string, boolean>>({});

  const toggleServicePanel = (id: string) =>
    setOpenServices((prev) => ({ ...prev, [id]: !prev[id] }));

  const handleStart = () => {
    startMutation.mutate(undefined, {
      onSuccess: () => toast.success("Listener started"),
      onError: () => toast.error("Failed to start listener"),
    });
  };

  const handleStop = () => {
    stopMutation.mutate(undefined, {
      onSuccess: () => toast.success("Listener stopped"),
      onError: () => toast.error("Failed to stop listener"),
    });
  };

  const handleClearCache = () => {
    cacheMutation.mutate(undefined, {
      onSuccess: () => toast.success("Cache cleared"),
      onError: () => toast.error("Failed to clear cache"),
    });
  };

  const handleToggleService = (serviceId: string) => {
    toggleMutation.mutate(serviceId, {
      onSuccess: (res) => toast.success(res.message),
      onError: () => toast.error("Failed to toggle service"),
    });
  };

  const isLoading = statusLoading || servicesLoading;

  if (isLoading || !status || !services) {
    return (
      <div className="space-y-6">
        <PageHeader title="Web Proxy" description="Proxy services and controls" />
        <Skeleton className="h-48" />
        <Skeleton className="h-64" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader title="Web Proxy" description="Proxy services and controls" />

      {/* Controls */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Controls</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-wrap gap-3">
          <Button
            onClick={handleStart}
            disabled={startMutation.isPending || status.listenerRunning}
            className="bg-success text-success-foreground hover:bg-success/90"
          >
            <Play className="h-4 w-4 mr-1" />
            {startMutation.isPending ? "Starting…" : "Start Listener"}
          </Button>
          <Button
            variant="destructive"
            onClick={handleStop}
            disabled={stopMutation.isPending || !status.listenerRunning}
          >
            <Square className="h-4 w-4 mr-1" />
            {stopMutation.isPending ? "Stopping…" : "Stop Listener"}
          </Button>
          <Button
            variant="outline"
            onClick={handleClearCache}
            disabled={cacheMutation.isPending}
          >
            <Trash2 className="h-4 w-4 mr-1" />
            {cacheMutation.isPending ? "Clearing…" : "Clear Session Cache"}
          </Button>
        </CardContent>
      </Card>

      {/* Proxy Services */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base flex items-center gap-2">
            <Network className="h-4 w-4" /> Proxy Services
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {services.map((service) => (
              <Collapsible
                key={service.id}
                open={openServices[service.id] ?? false}
                onOpenChange={() => toggleServicePanel(service.id)}
              >
                <div className="rounded-md border border-border">
                  <div className="flex items-center justify-between px-4 py-3">
                    <CollapsibleTrigger className="flex items-center gap-2 min-w-0 cursor-pointer hover:text-primary transition-colors">
                      <ChevronDown
                        className={`h-4 w-4 shrink-0 text-muted-foreground transition-transform duration-200 ${
                          openServices[service.id] ? "rotate-180" : ""
                        }`}
                      />
                      <span className="text-sm font-medium">{service.name}</span>
                    </CollapsibleTrigger>
                    <Button
                      size="sm"
                      variant={service.active ? "outline" : "default"}
                      disabled={toggleMutation.isPending}
                      onClick={() => handleToggleService(service.id)}
                      className={service.active
                        ? "text-destructive border-destructive/30 hover:bg-destructive/10"
                        : "bg-success text-success-foreground hover:bg-success/90"
                      }
                    >
                      {service.active ? (
                        <>
                          <PowerOff className="h-3.5 w-3.5 mr-1" /> Deactivate
                        </>
                      ) : (
                        <>
                          <Power className="h-3.5 w-3.5 mr-1" /> Activate
                        </>
                      )}
                    </Button>
                  </div>
                  <CollapsibleContent>
                    <div className="px-4 pb-3 pt-0">
                      <p className="text-sm text-muted-foreground pl-6">
                        {service.description}
                      </p>
                    </div>
                  </CollapsibleContent>
                </div>
              </Collapsible>
            ))}
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="pt-6">
          <p className="text-sm text-muted-foreground">
            These actions communicate with the Go backend through the centralized API layer.
            Proxy status is polled every 15 seconds. See <code className="text-xs">src/services/api/client.ts</code> for endpoint mapping.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
