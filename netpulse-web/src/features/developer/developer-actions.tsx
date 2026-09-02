"use client";

import { useRouter } from "next/navigation";
import { useState, type FormEvent, type ReactNode } from "react";

import { WEBHOOK_EVENTS } from "@/domain/developer";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function CreateMonitorForm() {
  return (
    <DevForm action="/api/dev/monitors" method="POST" success="Monitor stored. No checks were invented.">
      <Field id="name" label="Name" required />
      <Field id="target" label="Target" required placeholder="example.com" />
      <label htmlFor="type" className="text-sm font-medium">
        Type
      </label>
      <select id="type" name="type" required className="mt-1 h-11 w-full rounded-lg border border-input bg-transparent px-3 text-sm">
        <option value="http">HTTP</option>
        <option value="dns">DNS</option>
        <option value="tls">TLS</option>
      </select>
      <Field id="regions" label="Requested regions" placeholder="us-east, eu-west" />
      <Field id="frequencySeconds" label="Frequency (seconds)" type="number" defaultValue="300" />
      <Field id="timeoutSeconds" label="Timeout (seconds)" type="number" defaultValue="10" />
      <Field id="availabilityBelow" label="Availability threshold (0–1)" />
      <Field id="latencyMsAbove" label="Latency threshold (ms)" />
      <Field id="errorRateAbove" label="Error threshold (0–1)" />
    </DevForm>
  );
}

export function UpdateMonitorForm({
  id,
  name,
  target,
  type,
  regions,
  frequencySeconds,
  timeoutSeconds,
}: {
  id: string;
  name: string;
  target: string;
  type: string;
  regions: string[];
  frequencySeconds: number;
  timeoutSeconds: number;
}) {
  return (
    <DevForm action={`/api/dev/monitors/${id}`} method="PATCH" success="Monitor updated.">
      <Field id="name" label="Name" required defaultValue={name} />
      <Field id="target" label="Target" required defaultValue={target} />
      <label htmlFor="type" className="text-sm font-medium">
        Type
      </label>
      <select id="type" name="type" defaultValue={type} className="mt-1 h-11 w-full rounded-lg border border-input bg-transparent px-3 text-sm">
        <option value="http">HTTP</option>
        <option value="dns">DNS</option>
        <option value="tls">TLS</option>
      </select>
      <Field id="regions" label="Requested regions" defaultValue={regions.join(", ")} />
      <Field id="frequencySeconds" label="Frequency (seconds)" type="number" defaultValue={String(frequencySeconds)} />
      <Field id="timeoutSeconds" label="Timeout (seconds)" type="number" defaultValue={String(timeoutSeconds)} />
    </DevForm>
  );
}

export function RunMonitorButton({ id }: { id: string }) {
  return (
    <DevForm action={`/api/dev/monitors/${id}/run`} method="POST" success="A worker vantage check was requested." hideSave>
      <Button type="submit">Run check</Button>
    </DevForm>
  );
}

export function DeleteMonitorButton({ id }: { id: string }) {
  return (
    <DevForm
      action={`/api/dev/monitors/${id}`}
      method="DELETE"
      success="Monitor deleted."
      redirectTo="/developers/monitors"
      confirm="Delete this monitor?"
      hideSave
    >
      <Button type="submit" variant="destructive">
        Delete monitor
      </Button>
    </DevForm>
  );
}

export function CreateKeyForm() {
  return (
    <DevForm action="/api/dev/keys" method="POST" success="Key created. Copy the secret now; it is not stored in the browser." secretField="keySecret">
      <Field id="name" label="Name" defaultValue="API key" />
      <Field id="scopes" label="Scopes" defaultValue="monitors:read,monitors:write,incidents:read,webhooks:read,webhooks:write,usage:read,sla:read,alerts:read,alerts:write,dashboard:read" />
      <Field id="rateLimitPerMin" label="Rate limit per minute" type="number" defaultValue="60" />
    </DevForm>
  );
}

export function RotateKeyButton({ id }: { id: string }) {
  return (
    <DevForm action={`/api/dev/keys/${id}/rotate`} method="POST" success="Key rotated. Copy the new secret now." secretField="keySecret" hideSave>
      <Button type="submit" variant="outline" size="sm">
        Rotate
      </Button>
    </DevForm>
  );
}

export function RevokeKeyButton({ id }: { id: string }) {
  return (
    <DevForm action={`/api/dev/keys/${id}/revoke`} method="POST" success="Key revoked." confirm="Revoke this API key?" hideSave>
      <Button type="submit" variant="destructive" size="sm">
        Revoke
      </Button>
    </DevForm>
  );
}

export function CreateWebhookForm() {
  return (
    <DevForm action="/api/dev/webhooks" method="POST" success="Webhook stored. Copy the signing secret now." secretField="webhookSecret">
      <Field id="url" label="HTTPS URL" required placeholder="https://example.com/hooks" />
      <fieldset className="space-y-2">
        <legend className="text-sm font-medium">Events</legend>
        {WEBHOOK_EVENTS.map((event) => (
          <label key={event} className="flex items-center gap-2 text-sm">
            <input type="checkbox" name="events" value={event} />
            {event}
          </label>
        ))}
      </fieldset>
    </DevForm>
  );
}

export function RotateWebhookButton({ id }: { id: string }) {
  return (
    <DevForm action={`/api/dev/webhooks/${id}/rotate`} method="POST" success="Signing secret rotated." secretField="webhookSecret" hideSave>
      <Button type="submit" variant="outline" size="sm">
        Rotate secret
      </Button>
    </DevForm>
  );
}

export function DeleteWebhookButton({ id }: { id: string }) {
  return (
    <DevForm action={`/api/dev/webhooks/${id}`} method="DELETE" success="Webhook deleted." confirm="Delete this webhook?" hideSave>
      <Button type="submit" variant="destructive" size="sm">
        Delete
      </Button>
    </DevForm>
  );
}

export function RetryWebhookButton({ id }: { id: string }) {
  return (
    <DevForm action={`/api/dev/webhooks/${id}/retry`} method="POST" success="Retryable deliveries were attempted again." hideSave>
      <Button type="submit" variant="outline" size="sm">
        Retry deliveries
      </Button>
    </DevForm>
  );
}

export function CreateAlertForm() {
  return (
    <DevForm action="/api/dev/alerts" method="POST" success="Alert rule stored. Email is not sent.">
      <label htmlFor="kind" className="text-sm font-medium">
        Kind
      </label>
      <select id="kind" name="kind" required className="mt-1 h-11 w-full rounded-lg border border-input bg-transparent px-3 text-sm">
        <option value="availability">Availability threshold</option>
        <option value="latency">Latency threshold</option>
        <option value="error">Error threshold</option>
        <option value="incident">Incident events</option>
      </select>
      <Field id="threshold" label="Threshold" defaultValue="0.99" />
      <Field id="monitorId" label="Monitor id (optional)" />
    </DevForm>
  );
}

export function DeleteAlertButton({ id }: { id: string }) {
  return (
    <DevForm action={`/api/dev/alerts/${id}`} method="DELETE" success="Alert rule deleted." hideSave>
      <Button type="submit" variant="outline" size="sm">
        Delete
      </Button>
    </DevForm>
  );
}

function Field({
  id,
  label,
  type = "text",
  required,
  defaultValue,
  placeholder,
}: {
  id: string;
  label: string;
  type?: string;
  required?: boolean;
  defaultValue?: string;
  placeholder?: string;
}) {
  return (
    <div>
      <label htmlFor={id} className="text-sm font-medium">
        {label}
      </label>
      <Input id={id} name={id} type={type} required={required} defaultValue={defaultValue} placeholder={placeholder} className="mt-1" />
    </div>
  );
}

function DevForm({
  action,
  method,
  children,
  success,
  redirectTo,
  confirm,
  hideSave,
  secretField,
}: {
  action: string;
  method: string;
  children: ReactNode;
  success?: string;
  redirectTo?: string;
  confirm?: string;
  hideSave?: boolean;
  secretField?: "keySecret" | "webhookSecret";
}) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);
  const [secret, setSecret] = useState<string | null>(null);
  const [pending, setPending] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (confirm && !window.confirm(confirm)) {
      return;
    }
    setPending(true);
    setError(null);
    setNotice(null);
    setSecret(null);
    const form = event.currentTarget;
    const data = new FormData(form);
    const payload: Record<string, unknown> = {};
    const events: string[] = [];
    for (const [key, value] of data.entries()) {
      if (key === "events") {
        events.push(String(value));
        continue;
      }
      payload[key] = String(value);
    }
    if (events.length) {
      payload.events = events;
    }
    if (typeof payload.regions === "string") {
      payload.regions = payload.regions
        .split(",")
        .map((item) => item.trim())
        .filter(Boolean);
    }
    if (typeof payload.scopes === "string") {
      payload.scopes = payload.scopes
        .split(",")
        .map((item) => item.trim())
        .filter(Boolean);
    }
    for (const key of ["frequencySeconds", "timeoutSeconds", "rateLimitPerMin", "latencyMsAbove"]) {
      if (typeof payload[key] === "string" && payload[key] !== "") {
        payload[key] = Number(payload[key]);
      }
    }
    for (const key of ["availabilityBelow", "errorRateAbove", "threshold"]) {
      if (typeof payload[key] === "string" && payload[key] !== "") {
        payload[key] = Number(payload[key]);
      } else if (payload[key] === "") {
        delete payload[key];
      }
    }
    if (payload.monitorId === "") {
      delete payload.monitorId;
    }
    try {
      const response = await fetch(action, {
        method,
        headers: { "Content-Type": "application/json" },
        body: method === "GET" ? undefined : JSON.stringify(payload),
      });
      const body = (await response.json()) as Record<string, unknown> & {
        ok?: boolean;
        error?: { message?: string };
      };
      if (!response.ok || body.ok === false) {
        setError(body.error?.message ?? "The request failed.");
        return;
      }
      if (secretField && typeof body[secretField] === "string") {
        setSecret(String(body[secretField]));
      }
      if (success) {
        setNotice(success);
      }
      if (redirectTo) {
        router.push(redirectTo);
      }
      router.refresh();
    } catch {
      setError("The developer service is unavailable.");
    } finally {
      setPending(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-3">
      {children}
      {error ? <p className="text-sm text-destructive">{error}</p> : null}
      {notice ? <p className="text-sm text-muted-foreground">{notice}</p> : null}
      {secret ? (
        <p className="break-all rounded-md border border-border bg-muted p-3 text-sm">
          Secret (shown once): {secret}
        </p>
      ) : null}
      {hideSave ? null : (
        <Button type="submit" disabled={pending}>
          {pending ? "Saving…" : "Save"}
        </Button>
      )}
    </form>
  );
}
