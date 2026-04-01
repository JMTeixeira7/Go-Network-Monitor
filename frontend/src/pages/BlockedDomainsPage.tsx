import { useState } from "react";
import { PageHeader } from "@/components/layout/PageHeader";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { StatusBadge } from "@/components/shared/StatusBadge";
import { LoadingTable } from "@/components/shared/LoadingTable";
import { EmptyState } from "@/components/shared/EmptyState";
import { ConfirmDialog } from "@/components/shared/ConfirmDialog";
import { useBlockedDomains, useDeleteBlockedDomain } from "@/hooks/use-api";
import { Plus, Search, Trash2, ExternalLink, Calendar } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { format } from "date-fns";
import { toast } from "sonner";
import { AddDomainDialog } from "@/features/blocked-domains/AddDomainDialog";

export default function BlockedDomainsPage() {
  const { data: domains, isLoading } = useBlockedDomains();
  const deleteMutation = useDeleteBlockedDomain();
  const navigate = useNavigate();
  const [search, setSearch] = useState("");
  const [deleteTarget, setDeleteTarget] = useState<string | null>(null);
  const [addOpen, setAddOpen] = useState(false);

  const filtered = domains?.filter((d) =>
    d.domain.toLowerCase().includes(search.toLowerCase())
  );

  const handleDelete = () => {
    if (!deleteTarget) return;
    deleteMutation.mutate(deleteTarget, {
      onSuccess: () => {
        toast.success(`Removed ${deleteTarget}`);
        setDeleteTarget(null);
      },
      onError: () => toast.error("Failed to remove domain"),
    });
  };

  return (
    <div className="space-y-6">
      <PageHeader title="Blocked Domains" description="Manage domains blocked by the proxy">
        <Button size="sm" onClick={() => setAddOpen(true)}>
          <Plus className="h-4 w-4 mr-1" /> Add Domain
        </Button>
      </PageHeader>

      <div className="flex items-center gap-2">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Filter domains…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-8"
          />
        </div>
      </div>

      {isLoading ? (
        <LoadingTable columns={4} />
      ) : !filtered || filtered.length === 0 ? (
        <EmptyState
          title="No blocked domains"
          description={search ? "No domains match your filter." : "Add a domain to get started."}
        >
          {!search && (
            <Button size="sm" onClick={() => setAddOpen(true)}>
              <Plus className="h-4 w-4 mr-1" /> Add Domain
            </Button>
          )}
        </EmptyState>
      ) : (
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Domain</TableHead>
                <TableHead>Schedules</TableHead>
                <TableHead>Added</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map((domain) => (
                <TableRow key={domain.domain}>
                  <TableCell className="font-mono text-sm">{domain.domain}</TableCell>
                  <TableCell>
                    {domain.schedulesCount > 0 ? (
                      <StatusBadge status="scheduled" label={`${domain.schedulesCount} rule(s)`} />
                    ) : (
                      <span className="text-xs text-muted-foreground">Always blocked</span>
                    )}
                  </TableCell>
                  <TableCell className="text-sm text-muted-foreground">
                    {domain.createdAt && !Number.isNaN(new Date(domain.createdAt).getTime())
                      ? format(new Date(domain.createdAt), "MMM d, yyyy")
                      : "—"}
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex items-center justify-end gap-1">
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => navigate(`/blocked-domains/${encodeURIComponent(domain.domain)}`)}
                        title="View details"
                      >
                        <ExternalLink className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => navigate(`/blocked-domains/${encodeURIComponent(domain.domain)}`)}
                        title="Schedules"
                      >
                        <Calendar className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => setDeleteTarget(domain.domain)}
                        title="Remove"
                      >
                        <Trash2 className="h-4 w-4 text-destructive" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}

      <AddDomainDialog open={addOpen} onOpenChange={setAddOpen} />

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(o) => !o && setDeleteTarget(null)}
        title="Remove blocked domain"
        description={`Are you sure you want to unblock "${deleteTarget}"? This action can be undone by re-adding the domain.`}
        onConfirm={handleDelete}
        loading={deleteMutation.isPending}
        destructive
      />
    </div>
  );
}
