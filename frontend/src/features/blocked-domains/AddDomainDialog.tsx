import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAddBlockedDomain } from "@/hooks/use-api";
import { toast } from "sonner";
import { Plus, Trash2 } from "lucide-react";
import type { ScheduleRule } from "@/types";

const WEEKDAYS = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];

interface ScheduleInput {
  startTime: string;
  endTime: string;
  weekday: string;
}

interface AddDomainDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function AddDomainDialog({ open, onOpenChange }: AddDomainDialogProps) {
  const [domain, setDomain] = useState("");
  const [schedules, setSchedules] = useState<ScheduleInput[]>([]);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const mutation = useAddBlockedDomain();

  const addScheduleRow = () => {
    setSchedules([...schedules, { startTime: "", endTime: "", weekday: "" }]);
  };

  const removeScheduleRow = (index: number) => {
    setSchedules(schedules.filter((_, i) => i !== index));
  };

  const updateSchedule = (index: number, field: keyof ScheduleInput, value: string) => {
    const updated = [...schedules];
    updated[index] = { ...updated[index], [field]: value };
    setSchedules(updated);
  };

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};
    if (!domain.trim()) {
      newErrors.domain = "Domain is required";
    } else if (!/^[a-zA-Z0-9][a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(domain.trim())) {
      newErrors.domain = "Enter a valid domain name";
    }
    schedules.forEach((s, i) => {
      if (s.startTime && s.endTime && s.startTime >= s.endTime) {
        newErrors[`schedule-${i}`] = "Start time must be before end time";
      }
    });
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = () => {
    if (!validate()) return;
    const validSchedules = schedules.filter((s) => s.startTime && s.endTime && s.weekday);
    mutation.mutate(
      { domain: domain.trim(), schedules: validSchedules },
      {
        onSuccess: () => {
          toast.success(`Added ${domain.trim()}`);
          setDomain("");
          setSchedules([]);
          setErrors({});
          onOpenChange(false);
        },
        onError: () => toast.error("Failed to add domain"),
      }
    );
  };

  const reset = () => {
    setDomain("");
    setSchedules([]);
    setErrors({});
  };

  return (
    <Dialog
      open={open}
      onOpenChange={(o) => {
        if (!o) reset();
        onOpenChange(o);
      }}
    >
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Add Blocked Domain</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="domain">Domain name</Label>
            <Input
              id="domain"
              placeholder="example.com"
              value={domain}
              onChange={(e) => setDomain(e.target.value)}
              className="font-mono text-sm"
            />
            {errors.domain && <p className="text-xs text-destructive">{errors.domain}</p>}
          </div>

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label>Schedules (optional)</Label>
              <Button type="button" variant="ghost" size="sm" onClick={addScheduleRow}>
                <Plus className="h-3.5 w-3.5 mr-1" /> Add schedule
              </Button>
            </div>

            {schedules.map((s, i) => (
              <div key={i} className="flex items-start gap-2 p-3 rounded-md bg-muted/50 border">
                <div className="flex-1 grid grid-cols-3 gap-2">
                  <div>
                    <Label className="text-xs text-muted-foreground">Start</Label>
                    <Input
                      type="time"
                      value={s.startTime}
                      onChange={(e) => updateSchedule(i, "startTime", e.target.value)}
                      className="text-sm"
                    />
                  </div>
                  <div>
                    <Label className="text-xs text-muted-foreground">End</Label>
                    <Input
                      type="time"
                      value={s.endTime}
                      onChange={(e) => updateSchedule(i, "endTime", e.target.value)}
                      className="text-sm"
                    />
                  </div>
                  <div>
                    <Label className="text-xs text-muted-foreground">Day</Label>
                    <Select
                      value={s.weekday}
                      onValueChange={(v) => updateSchedule(i, "weekday", v)}
                    >
                      <SelectTrigger className="text-sm">
                        <SelectValue placeholder="Day" />
                      </SelectTrigger>
                      <SelectContent>
                        {WEEKDAYS.map((d) => (
                          <SelectItem key={d} value={d}>{d}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <Button variant="ghost" size="icon" className="shrink-0 mt-5" onClick={() => removeScheduleRow(i)}>
                  <Trash2 className="h-3.5 w-3.5 text-destructive" />
                </Button>
                {errors[`schedule-${i}`] && (
                  <p className="text-xs text-destructive col-span-full">{errors[`schedule-${i}`]}</p>
                )}
              </div>
            ))}
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>Cancel</Button>
          <Button onClick={handleSubmit} disabled={mutation.isPending}>
            {mutation.isPending ? "Adding…" : "Add Domain"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
