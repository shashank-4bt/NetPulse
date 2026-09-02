import type { Metadata } from "next";

import { EmptyState } from "@/components/feedback/empty-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import {
  PasswordForm,
  RevokeOthersButton,
  RevokeSessionButton,
} from "@/features/account/account-actions";
import { requireAccount } from "@/lib/account/guard";
import { getEvents, getSessions } from "@/lib/api/account";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Security",
  robots: { index: false, follow: false },
};

export default async function SecurityPage() {
  const { unavailable } = await requireAccount("/account/security");
  if (unavailable) {
    return (
      <main id="main-content" className="flex-1">
        <PageHero eyebrow="Account" title="Security" description="Sessions and password changes need the API." />
        <PageContainer className="py-10">
          <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set." />
        </PageContainer>
      </main>
    );
  }
  const token = await readSessionToken();
  const sessions = token ? await getSessions(token) : null;
  const events = token ? await getEvents(token) : null;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Account"
        title="Security"
        description="Active sessions, password change, and security events. Session secrets stay in the HTTP-only cookie."
      />
      <PageContainer className="max-w-2xl space-y-8 py-10">
        <section className="space-y-3">
          <h2 className="text-lg font-semibold">Password</h2>
          <PasswordForm />
        </section>
        <section className="space-y-3">
          <h2 className="text-lg font-semibold">Active sessions</h2>
          <RevokeOthersButton />
          {!sessions || !sessions.ok ? (
            <EmptyState title="Sessions unavailable" description={sessions && !sessions.ok ? sessions.message : ""} />
          ) : (
            <ul className="space-y-3">
              {sessions.sessions.map((session) => (
                <li key={session.id} className="rounded-lg border border-border p-3 text-sm">
                  <p className="font-medium">
                    {session.label}
                    {session.current ? " · this session" : ""}
                  </p>
                  <p className="text-muted-foreground">
                    {session.ip || "IP not stored"} · last seen {session.lastSeenAt}
                  </p>
                  {session.current ? null : <RevokeSessionButton sessionId={session.id} />}
                </li>
              ))}
            </ul>
          )}
        </section>
        <section className="space-y-3">
          <h2 className="text-lg font-semibold">Security events</h2>
          {!events || !events.ok ? (
            <EmptyState title="Events unavailable" description={events && !events.ok ? events.message : ""} />
          ) : events.events.length === 0 ? (
            <p className="text-sm text-muted-foreground">No security events are stored yet.</p>
          ) : (
            <ul className="space-y-2 text-sm">
              {events.events.map((event) => (
                <li key={event.id}>
                  <span className="font-medium">{event.kind}</span>
                  <span className="text-muted-foreground">
                    {" "}
                    · {event.at} · {event.summary}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </section>
      </PageContainer>
    </main>
  );
}
