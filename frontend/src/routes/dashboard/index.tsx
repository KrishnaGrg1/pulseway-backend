import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import {
  getMonitors,
  getDashboardStats,
  createMonitor,
  deleteMonitor,
} from "#/lib/queries";
import { logout, isAuthenticated } from "#/lib/auth";
import type { Monitor } from "#/lib/types";
import { Button } from "#/components/ui/button";
import { Input } from "#/components/ui/input";
import { Label } from "#/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "#/components/ui/card";
import { Badge } from "#/components/ui/badge";

export const Route = createFileRoute("/dashboard/")({
  component: DashboardPage,
});

function DashboardPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState("");
  const [url, setUrl] = useState("");
  const [interval, setInterval] = useState(60);
  const [liveStatuses, setLiveStatuses] = useState<Record<number, string>>({});

  // Redirect if not authenticated
  useEffect(() => {
    if (!isAuthenticated()) {
      navigate({ to: "/login" });
    }
  }, []);

  // SSE connection for live updates
  useEffect(() => {
    const apiUrl = import.meta.env.VITE_API_URL;
    let source: EventSource;

    const connect = () => {
      source = new EventSource(`${apiUrl}/sse`);

      source.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.monitor_id) {
          setLiveStatuses((prev) => ({
            ...prev,
            [data.monitor_id]: data.status,
          }));
          queryClient.invalidateQueries({ queryKey: ["monitors"] });
        }
      };

      source.onerror = () => {
        source.close();
        // Reconnect after 3 seconds
        setTimeout(connect, 3000);
      };
    };

    connect();

    return () => source?.close();
  }, []);

  const { data: monitors, isLoading: monitorsLoading } = useQuery({
    queryKey: ["monitors"],
    queryFn: getMonitors,
  });

  const { data: stats } = useQuery({
    queryKey: ["stats"],
    queryFn: getDashboardStats,
    refetchInterval: 30000,
  });

  const createMutation = useMutation({
    mutationFn: () => createMonitor({ name, url, interval_secs: interval }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["monitors"] });
      queryClient.invalidateQueries({ queryKey: ["stats"] });
      setShowForm(false);
      setName("");
      setUrl("");
      setInterval(60);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteMonitor(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["monitors"] });
      queryClient.invalidateQueries({ queryKey: ["stats"] });
    },
  });

  const getStatus = (monitor: Monitor) => {
    return liveStatuses[monitor.id] || "unknown";
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="border-b">
        <div className="max-w-6xl mx-auto px-4 py-4 flex items-center justify-between">
          <h1 className="text-xl font-semibold">Pulseway</h1>
          <Button variant="outline" onClick={logout}>
            Logout
          </Button>
        </div>
      </div>

      <div className="max-w-6xl mx-auto px-4 py-8 space-y-8">
        {/* Stats */}
        {stats && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm text-muted-foreground">
                  Total Monitors
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">{stats.total_monitors}</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm text-muted-foreground">
                  Active
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">{stats.active_monitors}</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm text-muted-foreground">
                  Uptime (24h)
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">
                  {stats.uptime_percentage?.toFixed(1)}%
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm text-muted-foreground">
                  Avg Latency
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">
                  {stats.avg_latency_ms?.toFixed(0)}ms
                </p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Monitors */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold">Monitors</h2>
            <Button onClick={() => setShowForm(!showForm)}>
              {showForm ? "Cancel" : "Add Monitor"}
            </Button>
          </div>

          {/* Add monitor form */}
          {showForm && (
            <Card>
              <CardContent className="pt-6 space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <Label>Name</Label>
                    <Input
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      placeholder="My API"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>URL</Label>
                    <Input
                      value={url}
                      onChange={(e) => setUrl(e.target.value)}
                      placeholder="https://api.example.com/health"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>Interval (seconds)</Label>
                    <Input
                      type="number"
                      value={interval}
                      onChange={(e) => setInterval(Number(e.target.value))}
                      placeholder="60"
                    />
                  </div>
                </div>
                <Button
                  onClick={() => createMutation.mutate()}
                  disabled={createMutation.isPending || !name || !url}
                >
                  {createMutation.isPending ? "Creating..." : "Create Monitor"}
                </Button>
              </CardContent>
            </Card>
          )}

          {/* Monitor list */}
          {monitorsLoading ? (
            <p className="text-muted-foreground">Loading...</p>
          ) : monitors?.length === 0 ? (
            <Card>
              <CardContent className="py-12 text-center text-muted-foreground">
                No monitors yet. Add one to get started.
              </CardContent>
            </Card>
          ) : (
            <div className="space-y-3">
              {monitors?.map((monitor) => (
                <Card key={monitor.id}>
                  <CardContent className="py-4 flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <Badge
                        variant={
                          getStatus(monitor) === "up"
                            ? "default"
                            : getStatus(monitor) === "down"
                              ? "destructive"
                              : "secondary"
                        }
                      >
                        {getStatus(monitor) === "up"
                          ? "● UP"
                          : getStatus(monitor) === "down"
                            ? "● DOWN"
                            : "○ PENDING"}
                      </Badge>
                      <div>
                        <p className="font-medium">{monitor.name}</p>
                        <p className="text-sm text-muted-foreground">
                          {monitor.url}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-4">
                      <p className="text-sm text-muted-foreground">
                        every {monitor.interval_secs}s
                      </p>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => deleteMutation.mutate(monitor.id)}
                        disabled={deleteMutation.isPending}
                      >
                        Delete
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
