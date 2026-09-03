import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { IncidentActionForm } from "@/features/admin/actions";
import { AdminForbidden, AdminScreen, AdminUnavailable } from "@/features/admin/screen";
import { requireAdmin } from "@/lib/admin/page";
import { getAdminIncidents } from "@/lib/api/admin";

export const dynamic = "force-dynamic";

export const metadata: Metadata = { title: "Incidents", robots: { index: false, follow: false } };

export default async function AdminIncidentsPage() {
  const ctx = await requireAdmin("/admin/incidents");
  if (ctx.unavailable) {
    return <AdminUnavailable title="Incidents" description="Stored incidents stay unavailable until the API is connected." />;
  }
  if (ctx.forbidden) {
    return <AdminForbidden />;
  }
  const loaded = await getAdminIncidents(ctx.token);
  return (
    <AdminScreen
      title="Incident operations"
      description="Investigate, annotate, escalate, resolve, or override automated classification. Every override is audited."
    >
      {!loaded.ok ? (
        <DevelopmentBanner description={loaded.message} />
      ) : loaded.incidents.length === 0 ? (
        <p className="text-sm text-muted-foreground">No stored incidents. Operators cannot invent an outage from this page.</p>
      ) : (
        <ul className="space-y-4 text-sm">
          {loaded.incidents.map((item) => (
            <li key={item.id} className="rounded-md border border-border p-4">
              <p className="font-medium">{item.title}</p>
              <p className="text-muted-foreground">
                {item.id} · {item.status} · {item.severity} · samples {item.sampleCount}
              </p>
              {item.overrideClassification ? (
                <p className="text-muted-foreground">
                  Override: {item.overrideClassification}. {item.overrideReason}
                </p>
              ) : null}
              {item.notes.length ? (
                <ul className="mt-2 space-y-1 text-muted-foreground">
                  {item.notes.map((note) => (
                    <li key={note.id}>
                      {note.at} · {note.kind} · {note.body}
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="mt-2 text-muted-foreground">No operator notes are stored.</p>
              )}
            </li>
          ))}
        </ul>
      )}
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Investigate</CardTitle>
          </CardHeader>
          <CardContent>
            <IncidentActionForm action="investigate" label="Investigate" />
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Annotate</CardTitle>
          </CardHeader>
          <CardContent>
            <IncidentActionForm action="annotate" label="Annotate" />
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Escalate</CardTitle>
          </CardHeader>
          <CardContent>
            <IncidentActionForm action="escalate" label="Escalate" />
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Resolve</CardTitle>
          </CardHeader>
          <CardContent>
            <IncidentActionForm
              action="resolve"
              label="Resolve"
              extra={
                <div className="space-y-2 text-sm">
                  <label className="flex items-center gap-2">
                    <input type="checkbox" name="override" value="true" />
                    Override automated gates (audited)
                  </label>
                  <label className="flex items-center gap-2">
                    <input type="checkbox" name="identifiedCause" value="true" />
                    Identified cause is stored
                  </label>
                </div>
              }
            />
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Override classification</CardTitle>
          </CardHeader>
          <CardContent>
            <IncidentActionForm action="override" label="Override" />
          </CardContent>
        </Card>
      </div>
    </AdminScreen>
  );
}
