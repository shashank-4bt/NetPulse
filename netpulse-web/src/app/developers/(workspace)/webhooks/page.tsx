import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  CreateWebhookForm,
  DeleteWebhookButton,
  RetryWebhookButton,
  RotateWebhookButton,
} from "@/features/developer/developer-actions";
import { DeveloperUnavailable } from "@/features/developer/unavailable";
import { requireAccount } from "@/lib/account/guard";
import { getDevDeliveries, getDevWebhooks } from "@/lib/api/developer";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Webhooks",
  robots: { index: false, follow: false },
};

export default async function DeveloperWebhooksPage() {
  const { unavailable } = await requireAccount("/developers/webhooks");
  if (unavailable) {
    return <DeveloperUnavailable title="Webhooks" description="Webhook endpoints need the API." />;
  }
  const token = await readSessionToken();
  const loaded = token ? await getDevWebhooks(token) : null;
  const webhooks = loaded && loaded.ok ? loaded.webhooks : [];
  const deliveries = token
    ? await Promise.all(webhooks.map(async (hook) => ({ id: hook.id, result: await getDevDeliveries(token, hook.id) })))
    : [];

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Developers"
        title="Webhooks"
        description="HTTPS endpoints only. Local and private receivers are rejected. Signing secrets are shown once."
      />
      <PageContainer className="grid gap-6 py-10 lg:grid-cols-2">
        <div className="space-y-4">
          {!loaded || !loaded.ok ? (
            <DevelopmentBanner description={loaded && !loaded.ok ? loaded.message : "Webhooks unavailable."} />
          ) : webhooks.length === 0 ? (
            <p className="text-sm text-muted-foreground">No webhooks are stored.</p>
          ) : (
            webhooks.map((hook) => {
              const found = deliveries.find((item) => item.id === hook.id);
              const items = found && found.result.ok ? found.result.deliveries : [];
              return (
                <Card key={hook.id}>
                  <CardHeader>
                    <CardTitle>{hook.url}</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-3 text-sm">
                    <p className="text-muted-foreground">
                      {hook.events.join(", ")} · hint {hook.secretHint}…
                    </p>
                    <div className="flex flex-wrap gap-2">
                      <RotateWebhookButton id={hook.id} />
                      <RetryWebhookButton id={hook.id} />
                      <DeleteWebhookButton id={hook.id} />
                    </div>
                    {items.length === 0 ? (
                      <p className="text-muted-foreground">No deliveries are stored.</p>
                    ) : (
                      <ul className="space-y-2">
                        {items.map((item) => (
                          <li key={item.id}>
                            <p className="font-medium">
                              {item.event} · {item.status} · attempt {item.attempt}
                            </p>
                            <p className="break-all text-muted-foreground">
                              {item.eventId} · {item.idempotencyKey}
                            </p>
                            <p className="text-muted-foreground">{item.summary}</p>
                          </li>
                        ))}
                      </ul>
                    )}
                  </CardContent>
                </Card>
              );
            })
          )}
        </div>
        <Card>
          <CardHeader>
            <CardTitle>Add webhook</CardTitle>
          </CardHeader>
          <CardContent>
            <CreateWebhookForm />
          </CardContent>
        </Card>
      </PageContainer>
    </main>
  );
}
