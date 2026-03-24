import { useParams, useNavigate } from "react-router-dom";
import { PageHeader } from "@/components/layout/PageHeader";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { StatusBadge } from "@/components/shared/StatusBadge";
import { Skeleton } from "@/components/ui/skeleton";
import { useBlockedDomain } from "@/hooks/use-api";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { ArrowLeft, Clock, Plus, ShieldAlert, Eye } from "lucide-react";
import { format } from "date-fns";

export default function DomainDetailsPage() {
  const { domain } = useParams<{ domain: string }>();
  const navigate = useNavigate();
  const { data, isLoading, isError } = useBlockedDomain(domain ?? "");

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-64" />
        <Skeleton className="h-48" />
      </div>
    );
  }

  if (isError || !data) {
    return (
      <div className="space-y-4">
        <Button variant="ghost" onClick={() => navigate("/blocked-domains")}>
          <ArrowLeft className="h-4 w-4 mr-1" /> Back
        </Button>
        <p className="text-muted-foreground">Domain not found.</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2">
        <Button variant="ghost" size="icon" onClick={() => navigate("/blocked-domains")}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <PageHeader title={data.domain} description={`Added ${format(new Date(data.createdAt), "MMM d, yyyy")}`} />
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="text-base flex items-center gap-2">
              <Clock className="h-4 w-4" /> Schedules
            </CardTitle>
            <Button variant="outline" size="sm" disabled>
              <Plus className="h-3.5 w-3.5 mr-1" /> Add
            </Button>
          </CardHeader>
          <CardContent>
            {!data.schedules || data.schedules.length === 0 ? (
              <p className="text-sm text-muted-foreground">
                No schedules — domain is always blocked.
              </p>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Day</TableHead>
                    <TableHead>Start</TableHead>
                    <TableHead>End</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {data.schedules.map((s) => (
                    <TableRow key={s.id}>
                      <TableCell>
                        <StatusBadge status="scheduled" label={s.weekday} />
                      </TableCell>
                      <TableCell className="font-mono text-sm">{s.startTime}</TableCell>
                      <TableCell className="font-mono text-sm">{s.endTime}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </CardContent>
        </Card>

        <div className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="text-base flex items-center gap-2">
                <Eye className="h-4 w-4" /> Visit History
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                Visit history will be available when connected to the backend.
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-base flex items-center gap-2">
                <ShieldAlert className="h-4 w-4" /> Security Analysis
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                Security analysis and blocking reasons will appear here.
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
