import { ISOLATION_LAYERS } from "@/lib/content/isolation-layers";

export function IsolationPath() {
  return (
    <ol className="grid gap-3 sm:grid-cols-2 lg:grid-cols-5">
      {ISOLATION_LAYERS.map((layer, index) => (
        <li
          key={layer.id}
          className="flex flex-col gap-2 rounded-lg border border-border bg-card p-4"
        >
          <p className="font-mono text-xs text-muted-foreground">
            {String(index + 1).padStart(2, "0")}
          </p>
          <h3 className="text-sm font-medium">{layer.name}</h3>
          <p className="text-sm text-muted-foreground">{layer.summary}</p>
        </li>
      ))}
    </ol>
  );
}
