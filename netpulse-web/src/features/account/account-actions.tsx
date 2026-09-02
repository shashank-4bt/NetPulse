"use client";

import { useRouter } from "next/navigation";
import { useState, type FormEvent, type ReactNode } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function ProfileForm({ displayName }: { displayName: string }) {
  return (
    <ActionForm action="/api/me/profile" method="PATCH" success="Display name saved.">
      <label htmlFor="displayName" className="text-sm font-medium">
        Display name
      </label>
      <Input id="displayName" name="displayName" defaultValue={displayName} required className="mt-1" />
    </ActionForm>
  );
}

export function PasswordForm() {
  return (
    <ActionForm action="/api/auth/change-password" method="POST" success="Password changed.">
      <label htmlFor="currentPassword" className="text-sm font-medium">
        Current password
      </label>
      <Input id="currentPassword" name="currentPassword" type="password" required className="mt-1" />
      <label htmlFor="newPassword" className="mt-3 block text-sm font-medium">
        New password
      </label>
      <Input id="newPassword" name="newPassword" type="password" minLength={10} required className="mt-1" />
    </ActionForm>
  );
}

export function PrivacyForm({ telemetryOptIn }: { telemetryOptIn: boolean }) {
  return (
    <ActionForm action="/api/me/privacy" method="PUT" success="Telemetry preference saved.">
      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" name="telemetryOptIn" defaultChecked={telemetryOptIn} value="true" />
        Allow optional product telemetry
      </label>
    </ActionForm>
  );
}

export function AlertsForm({
  emailEnabled,
  incidentAlerts,
}: {
  emailEnabled: boolean;
  incidentAlerts: boolean;
}) {
  return (
    <ActionForm action="/api/me/alerts" method="PUT" success="Alert preferences saved. Delivery is not configured.">
      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" name="emailEnabled" defaultChecked={emailEnabled} value="true" />
        Email alerts
      </label>
      <label className="mt-2 flex items-center gap-2 text-sm">
        <input type="checkbox" name="incidentAlerts" defaultChecked={incidentAlerts} value="true" />
        Incident alerts
      </label>
    </ActionForm>
  );
}

export function DeleteAccountButton() {
  return (
    <ActionForm
      action="/api/me/deletion"
      method="POST"
      success="Account deleted."
      redirectTo="/"
      confirm="Delete this account and owned diagnoses? This cannot be undone."
      hideSave
    >
      <Button type="submit" variant="destructive">
        Delete account
      </Button>
    </ActionForm>
  );
}

export function LogoutButton() {
  return (
    <ActionForm action="/api/auth/logout" method="POST" redirectTo="/login" hideSave>
      <Button type="submit" variant="outline">
        Sign out
      </Button>
    </ActionForm>
  );
}

export function RevokeSessionButton({ sessionId }: { sessionId: string }) {
  return (
    <ActionForm action={`/api/auth/sessions/${sessionId}/revoke`} method="POST" success="Session revoked." hideSave>
      <Button type="submit" variant="outline" size="sm">
        Revoke
      </Button>
    </ActionForm>
  );
}

export function RevokeOthersButton() {
  return (
    <ActionForm action="/api/auth/sessions/revoke-others" method="POST" success="Other sessions were revoked." hideSave>
      <Button type="submit" variant="outline">
        Revoke other sessions
      </Button>
    </ActionForm>
  );
}

export function ShareReportButton({ reportId }: { reportId: string }) {
  const [link, setLink] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  async function share() {
    setError(null);
    const response = await fetch(`/api/me/reports/${reportId}/share`, { method: "POST" });
    const body = (await response.json()) as { ok?: boolean; share?: { path?: string }; error?: { message?: string } };
    if (!response.ok || !body.share?.path) {
      setError(body.error?.message ?? "Share link was not created.");
      return;
    }
    setLink(body.share.path);
  }
  return (
    <div className="space-y-2">
      <Button type="button" variant="outline" size="sm" onClick={() => void share()}>
        Create share link
      </Button>
      {link ? <p className="text-sm break-all text-muted-foreground">{link}</p> : null}
      {error ? <p className="text-sm text-destructive">{error}</p> : null}
    </div>
  );
}

export function DeleteReportButton({ reportId }: { reportId: string }) {
  return (
    <ActionForm
      action={`/api/me/reports/${reportId}`}
      method="DELETE"
      success="Report deleted."
      confirm="Delete this report? Other users will not keep access."
      hideSave
    >
      <Button type="submit" variant="destructive" size="sm">
        Delete
      </Button>
    </ActionForm>
  );
}

export function ExportReportButton({ reportId }: { reportId: string }) {
  const [error, setError] = useState<string | null>(null);
  async function download() {
    setError(null);
    const response = await fetch(`/reports/${reportId}`, { headers: { Accept: "text/html" } });
    if (!response.ok) {
      setError("The report JSON is not available. PDF export is not ready.");
      return;
    }
    const api = await fetch(`/api/reports/${reportId}`);
    if (!api.ok) {
      setError("Export preparation failed. PDF export is not ready.");
      return;
    }
    const blob = await api.blob();
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = `netpulse-report-${reportId}.json`;
    anchor.click();
    URL.revokeObjectURL(url);
  }
  return (
    <div>
      <Button type="button" variant="outline" size="sm" onClick={() => void download()}>
        Download JSON
      </Button>
      {error ? <p className="mt-1 text-sm text-destructive">{error}</p> : null}
    </div>
  );
}

export function SaveServiceForm() {
  return (
    <ActionForm action="/api/me/saved-services" method="PUT" success="Service saved.">
      <label htmlFor="slug" className="text-sm font-medium">
        Catalog slug
      </label>
      <Input id="slug" name="slug" placeholder="youtube" required className="mt-1" />
    </ActionForm>
  );
}

function ActionForm({
  action,
  method,
  children,
  success,
  redirectTo,
  confirm,
  hideSave,
}: {
  action: string;
  method: string;
  children: ReactNode;
  success?: string;
  redirectTo?: string;
  confirm?: string;
  hideSave?: boolean;
}) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);
  const [pending, setPending] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (confirm && !window.confirm(confirm)) {
      return;
    }
    setPending(true);
    setError(null);
    setNotice(null);
    const form = event.currentTarget;
    const data = new FormData(form);
    const payload: Record<string, string | boolean> = {};
    for (const [key, value] of data.entries()) {
      payload[key] = String(value);
    }
    for (const input of Array.from(form.querySelectorAll<HTMLInputElement>('input[type="checkbox"]'))) {
      payload[input.name] = input.checked;
    }
    try {
      const response = await fetch(action, {
        method,
        headers: { "Content-Type": "application/json" },
        body: method === "GET" ? undefined : JSON.stringify(payload),
      });
      const body = (await response.json()) as { ok?: boolean; error?: { message?: string } };
      if (!response.ok || body.ok === false) {
        setError(body.error?.message ?? "The request failed.");
        return;
      }
      if (success) {
        setNotice(success);
      }
      if (redirectTo) {
        router.push(redirectTo);
      }
      router.refresh();
    } catch {
      setError("The account service is unavailable.");
    } finally {
      setPending(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-3">
      {children}
      {error ? <p className="text-sm text-destructive">{error}</p> : null}
      {notice ? <p className="text-sm text-muted-foreground">{notice}</p> : null}
      {hideSave ? null : (
        <Button type="submit" disabled={pending}>
          {pending ? "Saving…" : "Save"}
        </Button>
      )}
    </form>
  );
}
