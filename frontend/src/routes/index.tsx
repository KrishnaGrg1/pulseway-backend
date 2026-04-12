import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Button } from "#/components/ui/button";
import { Badge } from "#/components/ui/badge";
import { Card, CardContent } from "#/components/ui/card";
import LiveDemo from "#/components/LiveDemo";

export const Route = createFileRoute("/")({
  component: HomePage,
});

function HomePage() {
  const navigate = useNavigate();
  const [latency, setLatency] = useState(124);

  useEffect(() => {
    const latencies = [98, 45, 201, 67, 134, 89, 156];
    let i = 0;
    const interval = setInterval(() => {
      i++;
      setLatency(latencies[i % latencies.length]);
    }, 2000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-background">
      {/* Nav */}
      <nav className="flex items-center justify-between px-8 py-4 border-b">
        <span className="text-base font-medium">Pulseway</span>
        <div className="flex items-center gap-6">
          <a
            href="#features"
            className="text-sm text-muted-foreground hover:text-foreground"
          >
            Features
          </a>
          <a
            href="#pricing"
            className="text-sm text-muted-foreground hover:text-foreground"
          >
            Pricing
          </a>
          <a
            href="#"
            className="text-sm text-muted-foreground hover:text-foreground"
          >
            Docs
          </a>
          <Button size="sm" onClick={() => navigate({ to: "/register" })}>
            Get started
          </Button>
        </div>
      </nav>

      {/* Hero */}
      <section className="text-center py-24 px-4 max-w-3xl mx-auto">
        <div className="inline-block bg-secondary border text-xs text-muted-foreground px-4 py-1 rounded-full mb-6">
          Built with Go — 10k monitors on a $12 server
        </div>
        <h1 className="text-5xl font-medium leading-tight mb-5">
          Know when your APIs go down
          <br />
          before your users do
        </h1>
        <p className="text-lg text-muted-foreground max-w-md mx-auto mb-8 leading-relaxed">
          Pulseway monitors your endpoints every 30 seconds and alerts you
          instantly when something breaks.
        </p>
        <div className="flex gap-3 justify-center flex-wrap">
          <Button size="lg" onClick={() => navigate({ to: "/register" })}>
            Start monitoring free
          </Button>
          <Button
            size="lg"
            variant="outline"
            onClick={() => navigate({ to: "/login" })}
          >
            Sign in
          </Button>
        </div>
      </section>

      {/* Live Demo */}
      <section className="max-w-4xl mx-auto px-4 mb-24">
        <div className="border rounded-xl overflow-hidden">
          {/* Browser bar */}
          <div className="bg-secondary px-4 py-2 border-b flex items-center gap-2">
            <div className="w-3 h-3 rounded-full bg-border" />
            <div className="w-3 h-3 rounded-full bg-border" />
            <div className="w-3 h-3 rounded-full bg-border" />
            <span className="text-xs text-muted-foreground ml-2">
              app.pulseway.tech/dashboard
            </span>
          </div>
          {/* Dashboard preview */}
          <div className="p-6 bg-background">
            {/* Stats */}

            <div className="grid grid-cols-4 gap-3 mb-6">
              {[
                { label: "Total monitors", value: "4" },
                { label: "Active", value: "4" },
                { label: "Uptime (24h)", value: "99.8%" },
                { label: "Avg latency", value: `${latency}ms` },
              ].map((stat) => (
                <div key={stat.label} className="bg-secondary rounded-lg p-4">
                  <p className="text-xs text-muted-foreground mb-1">
                    {stat.label}
                  </p>
                  <p className="text-2xl font-medium">{stat.value}</p>
                </div>
              ))}
            </div>
            <LiveDemo />
            {/* Monitor rows */}
            <div className="flex items-center justify-between px-4 py-3 border rounded-lg opacity-40">
              <div className="flex items-center gap-3">
                <Badge variant="default" className="text-xs">
                  ● UP
                </Badge>
                <div>
                  <p className="text-sm font-medium">Auth Service</p>
                  <p className="text-xs text-muted-foreground">
                    https://auth.yourapp.com/ping
                  </p>
                </div>
              </div>
              <span className="text-xs text-muted-foreground">
                45ms · every 60s
              </span>
            </div>
            <div className="flex items-center justify-between px-4 py-3 border rounded-lg opacity-40">
              <div className="flex items-center gap-3">
                <Badge variant="destructive" className="text-xs">
                  ● DOWN
                </Badge>
                <div>
                  <p className="text-sm font-medium">Payment Webhook</p>
                  <p className="text-xs text-muted-foreground">
                    https://pay.yourapp.com/health
                  </p>
                </div>
              </div>
              <span className="text-xs text-muted-foreground">
                timeout · every 30s
              </span>
            </div>
          </div>
        </div>
      </section>

      {/* Features */}
      <section id="features" className="max-w-4xl mx-auto px-4 mb-24">
        <p className="text-sm text-muted-foreground text-center mb-3">
          Features
        </p>
        <h2 className="text-3xl font-medium text-center mb-3">
          Everything you need to stay online
        </h2>
        <p className="text-base text-muted-foreground text-center max-w-md mx-auto mb-12">
          Built for developers who care about reliability.
        </p>
        <div className="grid grid-cols-3 gap-4">
          {[
            {
              icon: "⚡",
              title: "30-second checks",
              desc: "Your endpoints are checked every 30 seconds. Catch outages in under a minute.",
            },
            {
              icon: "🔔",
              title: "Instant alerts",
              desc: "Get notified via email the moment a monitor goes down. Recovery alerts too.",
            },
            {
              icon: "📡",
              title: "Live dashboard",
              desc: "Real-time status updates via SSE. Dashboard updates as checks come in.",
            },
            {
              icon: "📈",
              title: "Latency tracking",
              desc: "Track response times over 24 hours. Spot degradation before outages.",
            },
            {
              icon: "🔁",
              title: "Incident history",
              desc: "Full log of every incident — when it started, resolved, how long it lasted.",
            },
            {
              icon: "🔐",
              title: "Secure by default",
              desc: "JWT auth, bcrypt passwords, HTTPS everywhere. Your data stays private.",
            },
          ].map((f) => (
            <Card key={f.title}>
              <CardContent className="pt-5">
                <div className="w-8 h-8 rounded-md bg-secondary border flex items-center justify-center text-base mb-4">
                  {f.icon}
                </div>
                <h3 className="text-sm font-medium mb-2">{f.title}</h3>
                <p className="text-xs text-muted-foreground leading-relaxed">
                  {f.desc}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      {/* Pricing */}
      <section id="pricing" className="max-w-4xl mx-auto px-4 mb-24">
        <p className="text-sm text-muted-foreground text-center mb-3">
          Pricing
        </p>
        <h2 className="text-3xl font-medium text-center mb-3">
          Simple, honest pricing
        </h2>
        <p className="text-base text-muted-foreground text-center max-w-md mx-auto mb-12">
          No hidden fees. Cancel anytime.
        </p>
        <div className="grid grid-cols-3 gap-4">
          {[
            {
              name: "Hobby",
              price: "$0",
              desc: "For personal projects and side hustles.",
              featured: false,
              features: [
                "5 monitors",
                "5-minute checks",
                "Email alerts",
                "7-day history",
              ],
              cta: "Get started",
            },
            {
              name: "Pro",
              price: "$12",
              desc: "For teams that cannot afford downtime.",
              featured: true,
              features: [
                "50 monitors",
                "30-second checks",
                "Email + webhook alerts",
                "90-day history",
                "Incident reports",
              ],
              cta: "Start free trial",
            },
            {
              name: "Business",
              price: "$49",
              desc: "For companies running critical infrastructure.",
              featured: false,
              features: [
                "Unlimited monitors",
                "10-second checks",
                "All alert channels",
                "1-year history",
                "Priority support",
              ],
              cta: "Contact us",
            },
          ].map((plan) => (
            <div
              key={plan.name}
              className={`border rounded-xl p-6 flex flex-col ${
                plan.featured ? "border-blue-500 border-2" : ""
              }`}
            >
              {plan.featured && (
                <span className="text-xs bg-blue-50 text-blue-800 dark:bg-blue-900 dark:text-blue-200 px-3 py-1 rounded-full self-start mb-3">
                  Most popular
                </span>
              )}
              <p className="text-base font-medium">{plan.name}</p>
              <p className="text-3xl font-medium mt-3 mb-1">
                {plan.price}
                <span className="text-sm font-normal text-muted-foreground">
                  {" "}
                  / month
                </span>
              </p>
              <p className="text-xs text-muted-foreground mb-4">{plan.desc}</p>
              <ul className="space-y-2 mb-6 flex-1">
                {plan.features.map((f) => (
                  <li
                    key={f}
                    className="text-xs text-muted-foreground flex gap-2 items-center border-b pb-2 last:border-0"
                  >
                    <span className="text-green-600">✓</span> {f}
                  </li>
                ))}
              </ul>
              <Button
                variant={plan.featured ? "default" : "outline"}
                className="w-full"
                onClick={() => navigate({ to: "/register" })}
              >
                {plan.cta}
              </Button>
            </div>
          ))}
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t px-8 py-6 text-center">
        <p className="text-sm text-muted-foreground">
          Pulseway · Built with Go + TanStack Start · Deployed on DigitalOcean
        </p>
      </footer>
    </div>
  );
}
