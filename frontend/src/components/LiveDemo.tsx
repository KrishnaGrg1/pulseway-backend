import { useEffect, useState } from "react";
import { Badge } from "./ui/badge";

type LiveMonitor = {
  name: string;
  url: string;
  status: string;
  latency: number;
};

export default function LiveDemo() {
  const [monitor, setMonitor] = useState<LiveMonitor>({
    name: "LevelUp",
    url: "https://www.melevelup.me",
    status: "unknown",
    latency: 0,
  });

  useEffect(() => {
    const apiUrl = import.meta.env.VITE_API_URL;
    let source: EventSource;

    const connect = () => {
      source = new EventSource(`${apiUrl}/sse`);

      source.onmessage = (event) => {
        const data = JSON.parse(event.data);
        // Only update if it's the LevelUp monitor
        if (data.monitor_id) {
          setMonitor((prev) => ({
            ...prev,
            status: data.status,
            latency: data.latency_ms,
          }));
        }
      };

      source.onerror = () => {
        source.close();
        setTimeout(connect, 3000);
      };
    };

    connect();
    return () => source?.close();
  }, []);

  return (
    <div className="flex items-center justify-between px-4 py-3 border rounded-lg">
      <div className="flex items-center gap-3">
        <Badge
          variant={
            monitor.status === "up"
              ? "default"
              : monitor.status === "down"
                ? "destructive"
                : "secondary"
          }
          className="text-xs"
        >
          {monitor.status === "up"
            ? "● UP"
            : monitor.status === "down"
              ? "● DOWN"
              : "○ CHECKING"}
        </Badge>
        <div>
          <p className="text-sm font-medium">{monitor.name}</p>
          <p className="text-xs text-muted-foreground">{monitor.url}</p>
        </div>
      </div>
      <span className="text-xs text-muted-foreground">
        {monitor.latency > 0 ? `${monitor.latency}ms` : "..."} · live
      </span>
    </div>
  );
}
