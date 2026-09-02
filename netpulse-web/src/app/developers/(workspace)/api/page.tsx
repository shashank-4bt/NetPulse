import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CreateKeyForm, RevokeKeyButton, RotateKeyButton } from "@/features/developer/developer-actions";
import { DeveloperUnavailable } from "@/features/developer/unavailable";
import { requireAccount } from "@/lib/account/guard";
import { getDevKeys } from "@/lib/api/developer";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "API keys",
  robots: { index: false, follow: false },
};

export default async function DeveloperAPIPage() {
  const { unavailable } = await requireAccount("/developers/api");
  if (unavailable) {
    return <DeveloperUnavailable title="API keys" description="Key management needs the API." />;
  }
  const token = await readSessionToken();
  const loaded = token ? await getDevKeys(token) : null;
  const keys = loaded && loaded.ok ? loaded.keys : [];

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Developers"
        title="API keys"
        description="Raw secrets are shown once on create or rotate. The server stores a hash. Keys are not written to localStorage."
      />
      <PageContainer className="grid gap-6 py-10 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Stored keys</CardTitle>
          </CardHeader>
          <CardContent>
            {!loaded || !loaded.ok ? (
              <DevelopmentBanner description={loaded && !loaded.ok ? loaded.message : "Keys unavailable."} />
            ) : keys.length === 0 ? (
              <p className="text-sm text-muted-foreground">No API keys are stored.</p>
            ) : (
              <ul className="space-y-4 text-sm">
                {keys.map((key) => (
                  <li key={key.id} className="space-y-2">
                    <p className="font-medium">
                      {key.name} · {key.prefix}…{key.last4}
                    </p>
                    <p className="text-muted-foreground">
                      {key.revoked ? "Revoked" : "Active"} · {key.rateLimitPerMin}/min · {key.scopes.join(", ")}
                    </p>
                    {key.revoked ? null : (
                      <div className="flex flex-wrap gap-2">
                        <RotateKeyButton id={key.id} />
                        <RevokeKeyButton id={key.id} />
                      </div>
                    )}
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Create key</CardTitle>
          </CardHeader>
          <CardContent>
            <CreateKeyForm />
          </CardContent>
        </Card>
      </PageContainer>
    </main>
  );
}
