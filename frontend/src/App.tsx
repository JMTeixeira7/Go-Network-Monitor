import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { Toaster } from "@/components/ui/toaster";
import { TooltipProvider } from "@/components/ui/tooltip";
import { AppLayout } from "@/components/layout/AppLayout";
import DashboardPage from "@/pages/DashboardPage";
import BlockedDomainsPage from "@/pages/BlockedDomainsPage";
import DomainDetailsPage from "@/pages/DomainDetailsPage";
import ActivityPage from "@/pages/ActivityPage";
import SystemPage from "@/pages/SystemPage";
import VisitedHostsPage from "@/pages/VisitedHostsPage";
import NotFound from "@/pages/NotFound";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { staleTime: 30_000, retry: 1 },
  },
});

const App = () => (
  <QueryClientProvider client={queryClient}>
    <TooltipProvider>
      <Toaster />
      <Sonner />
      <BrowserRouter>
        <Routes>
          <Route element={<AppLayout />}>
            <Route path="/" element={<DashboardPage />} />
            <Route path="/blocked-domains" element={<BlockedDomainsPage />} />
            <Route path="/blocked-domains/:domain" element={<DomainDetailsPage />} />
            <Route path="/activity" element={<ActivityPage />} />
            <Route path="/system" element={<SystemPage />} />
            <Route path="/visited-hosts" element={<VisitedHostsPage />} />
          </Route>
          <Route path="*" element={<NotFound />} />
        </Routes>
      </BrowserRouter>
    </TooltipProvider>
  </QueryClientProvider>
);

export default App;
