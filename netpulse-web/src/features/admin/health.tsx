import { ObservationCard } from "@/features/developer/metric-cards";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { AdminSystemView, HealthComponentView } from "@/domain/admin";

export function SystemHealth({ system }: { system: AdminSystemView }) {
  return (
    <div className="space-y-4">
      <p className="text-sm text-muted-foreground">{system.summary}</p>
      <div className="grid gap-4 md:grid-cols-2">
        <HealthCard item={system.api} />
        <HealthCard item={system.worker} />
        <HealthCard item={system.queue} />
        <HealthCard item={system.database} />
        <HealthCard item={system.cache} />
        <ObservationCard title="Measurement failures" item={system.measurementFailures} />
        <ObservationCard title="Error rates" item={system.errorRates} />
        <Card>
          <CardHeader>
            <CardTitle>Latency</CardTitle>
          </CardHeader>
          <CardContent className="space-y-1 text-sm text-muted-foreground">
            <p>{system.latency.summary}</p>
            <p>p50 {system.latency.p50 ?? "unmeasured"} · p95 {system.latency.p95 ?? "unmeasured"} · p99 {system.latency.p99 ?? "unmeasured"}</p>
            <p>Samples {system.latency.sampleCount}</p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function HealthCard({ item }: { item: HealthComponentView }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>
          {item.name} · {item.status}
        </CardTitle>
      </CardHeader>
      <CardContent className="text-sm text-muted-foreground">
        <p>{item.detail}</p>
        <p>{item.measured ? "Measured from this process or stored state." : "Unmeasured."}</p>
      </CardContent>
    </Card>
  );
}
