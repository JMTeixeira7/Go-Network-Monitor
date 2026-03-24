import { PageHeader } from "@/components/layout/PageHeader";
import { StatusBadge } from "@/components/shared/StatusBadge";
import { LoadingTable } from "@/components/shared/LoadingTable";
import { EmptyState } from "@/components/shared/EmptyState";
import { useActivity } from "@/hooks/use-api";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { format } from "date-fns";

export default function ActivityPage() {
  const { data: activity, isLoading } = useActivity();

  return (
    <div className="space-y-6">
      <PageHeader title="Activity" description="Request history and inspection results" />

      {isLoading ? (
        <LoadingTable columns={5} />
      ) : !activity || activity.length === 0 ? (
        <EmptyState title="No activity yet" description="Activity will appear when the backend starts processing requests." />
      ) : (
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Domain</TableHead>
                <TableHead>Action</TableHead>
                <TableHead>Reasons</TableHead>
                <TableHead>Time</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {activity.map((item) => (
                <TableRow key={item.id}>
                  <TableCell className="font-mono text-sm">{item.domain}</TableCell>
                  <TableCell>
                    <StatusBadge status={item.action} />
                  </TableCell>
                  <TableCell className="text-sm text-muted-foreground">
                    {item.reasons.length > 0 ? item.reasons.join(", ") : "—"}
                  </TableCell>
                  <TableCell className="text-sm text-muted-foreground">
                    {format(new Date(item.timestamp), "MMM d, HH:mm:ss")}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </div>
  );
}
