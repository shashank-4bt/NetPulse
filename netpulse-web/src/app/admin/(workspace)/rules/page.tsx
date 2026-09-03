import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { RuleLabelForm } from "@/features/admin/actions";
import { AdminForbidden, AdminScreen, AdminUnavailable } from "@/features/admin/screen";
import { requireAdmin } from "@/lib/admin/page";
import { getAdminRules } from "@/lib/api/admin";

export const dynamic = "force-dynamic";

export const metadata: Metadata = { title: "Rules", robots: { index: false, follow: false } };

export default async function AdminRulesPage() {
  const ctx = await requireAdmin("/admin/rules");
  if (ctx.unavailable) {
    return <AdminUnavailable title="Rules" description="Diagnostic rules stay unavailable until the API is connected." />;
  }
  if (ctx.forbidden) {
    return <AdminForbidden />;
  }
  const loaded = await getAdminRules(ctx.token);
  return (
    <AdminScreen
      title="Diagnostic rules"
      description="Rules, versions, and thresholds come from the shipped engine and remote configuration. False positives and false negatives stay 0 until labeled."
    >
      {!loaded.ok ? (
        <DevelopmentBanner description={loaded.message} />
      ) : (
        <>
          <Card>
            <CardHeader>
              <CardTitle>Outcomes</CardTitle>
            </CardHeader>
            <CardContent className="space-y-1 text-sm text-muted-foreground">
              <p>{loaded.outcomes.summary}</p>
              <p>
                False positives {loaded.outcomes.falsePositives} · false negatives {loaded.outcomes.falseNegatives} ·
                diagnoses {loaded.outcomes.sampleCount}
              </p>
              {Object.keys(loaded.outcomes.statuses).length ? (
                <ul>
                  {Object.entries(loaded.outcomes.statuses).map(([status, count]) => (
                    <li key={status}>
                      {status}: {count}
                    </li>
                  ))}
                </ul>
              ) : (
                <p>No stored diagnosis statuses.</p>
              )}
            </CardContent>
          </Card>
          <ul className="space-y-3 text-sm">
            {loaded.rules.map((item) => (
              <li key={item.id} className="rounded-md border border-border p-4">
                <p className="font-medium">
                  {item.name} · {item.version}
                </p>
                <p className="text-muted-foreground">{item.layer} · {item.summary}</p>
                <ul className="mt-2 text-muted-foreground">
                  {Object.entries(item.thresholds).map(([key, value]) => (
                    <li key={key}>
                      {key}: {value}
                    </li>
                  ))}
                </ul>
              </li>
            ))}
          </ul>
          <Card>
            <CardHeader>
              <CardTitle>Label a diagnosis</CardTitle>
            </CardHeader>
            <CardContent>
              <RuleLabelForm />
            </CardContent>
          </Card>
        </>
      )}
    </AdminScreen>
  );
}
