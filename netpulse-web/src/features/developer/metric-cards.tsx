import type { DeveloperIncidentView, ObservationView, PercentileView, RegionalView } from "@/domain/developer";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function ObservationCard({ title, item }: { title: string; item: ObservationView }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-1 text-sm text-muted-foreground">
        <p>{item.measured && item.value !== null ? `Stored ratio: ${item.value}` : "Not measured"}</p>
        <p>Samples: {item.sampleCount}</p>
        <p>{item.summary}</p>
      </CardContent>
    </Card>
  );
}

export function LatencyCard({ item }: { item: PercentileView }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Latency</CardTitle>
      </CardHeader>
      <CardContent className="space-y-1 text-sm text-muted-foreground">
        <p>P50: {item.p50 === null ? "not estimated" : `${item.p50} ms`}</p>
        <p>P95: {item.p95 === null ? "not estimated" : `${item.p95} ms`}</p>
        <p>P99: {item.p99 === null ? "not estimated" : `${item.p99} ms`}</p>
        <p>Samples: {item.sampleCount}</p>
        <p>{item.summary}</p>
      </CardContent>
    </Card>
  );
}

export function IncidentList({ items }: { items: DeveloperIncidentView[] }) {
  if (!items.length) {
    return <p className="text-sm text-muted-foreground">No tenant monitor incidents are stored.</p>;
  }
  return (
    <ul className="space-y-2 text-sm">
      {items.map((item) => (
        <li key={item.id}>
          <p className="font-medium">{item.title}</p>
          <p className="text-muted-foreground">
            {item.status} · samples {item.sampleCount} · {item.startedAt}
          </p>
          <p className="text-muted-foreground">{item.summary}</p>
        </li>
      ))}
    </ul>
  );
}

export function RegionalList({ items }: { items: RegionalView[] }) {
  if (!items.length) {
    return <p className="text-sm text-muted-foreground">No regional check samples are stored.</p>;
  }
  return (
    <ul className="space-y-2 text-sm">
      {items.map((item) => (
        <li key={item.region}>
          <p className="font-medium">
            {item.region} · {item.status} · {item.sampleCount} samples
          </p>
          <p className="text-muted-foreground">{item.summary}</p>
        </li>
      ))}
    </ul>
  );
}
