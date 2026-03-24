import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

type StatusVariant = "running" | "stopped" | "unknown" | "blocked" | "allowed" | "flagged" | "visited" | "scheduled" | "active" | "cleared";

const variantMap: Record<StatusVariant, string> = {
  running: "bg-success/15 text-success border-success/30",
  active: "bg-success/15 text-success border-success/30",
  stopped: "bg-destructive/15 text-destructive border-destructive/30",
  blocked: "bg-destructive/15 text-destructive border-destructive/30",
  unknown: "bg-warning/15 text-warning border-warning/30",
  flagged: "bg-warning/15 text-warning border-warning/30",
  allowed: "bg-primary/15 text-primary border-primary/30",
  visited: "bg-muted text-muted-foreground border-border",
  scheduled: "bg-primary/15 text-primary border-primary/30",
  cleared: "bg-muted text-muted-foreground border-border",
};

interface StatusBadgeProps {
  status: StatusVariant;
  label?: string;
  className?: string;
}

export function StatusBadge({ status, label, className }: StatusBadgeProps) {
  return (
    <Badge
      variant="outline"
      className={cn("text-xs font-medium capitalize", variantMap[status], className)}
    >
      {label ?? status}
    </Badge>
  );
}
