"use client";

import { useRouter } from "next/navigation";
import { useState, type FormEvent, type ReactNode } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function AdminField({
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

export function CreateFlagForm() {
  return (
    <AdminForm action="/api/admin/flags" method="POST" success="Flag stored. Rollout counts are not estimated.">
      <AdminField id="name" label="Name" required />
      <AdminField id="environment" label="Environment" placeholder="development" />
      <AdminField id="percentage" label="Percentage (0-100)" type="number" defaultValue="0" />
      <AdminField id="userIds" label="User ids (comma-separated)" />
      <AdminField id="orgIds" label="Organization ids (comma-separated)" />
      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" name="enabled" value="true" />
        Enabled
      </label>
    </AdminForm>
  );
}

export function UpdateConfigForm() {
  return (
    <AdminForm action="/api/admin/config" method="PUT" success="Remote configuration updated.">
      <AdminField id="key" label="Key" required placeholder="diagnose.timeoutSeconds" />
      <AdminField id="value" label="Value" required />
    </AdminForm>
  );
}

export function IncidentActionForm({
  action,
  label,
  extra,
}: {
  action: "annotate" | "investigate" | "escalate" | "resolve" | "override";
  label: string;
  extra?: ReactNode;
}) {
  return (
    <AdminForm action="" method="POST" success="Stored." pathFromId={`/api/admin/incidents/{id}/${action}`}>
      <AdminField id="id" label="Incident id" required />
      {action === "override" ? <AdminField id="classification" label="Classification" required /> : null}
      <AdminField id="reason" label={action === "annotate" ? "Note" : "Reason"} required={action === "override"} />
      {extra}
      <Button type="submit">{label}</Button>
    </AdminForm>
  );
}

export function RuleLabelForm() {
  return (
    <AdminForm action="/api/admin/rules/labels" method="POST" success="Label stored. Counts stay unlabeled at 0 until this point.">
      <AdminField id="diagnosisId" label="Diagnosis id" required />
      <div>
        <label htmlFor="kind" className="text-sm font-medium">
          Kind
        </label>
        <select id="kind" name="kind" className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm">
          <option value="false_positive">False positive</option>
          <option value="false_negative">False negative</option>
        </select>
      </div>
    </AdminForm>
  );
}

function AdminForm({
  action,
  method,
  children,
  success,
  pathFromId,
}: {
  action: string;
  method: string;
  children: ReactNode;
  success?: string;
  pathFromId?: string;
}) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);
  const [pending, setPending] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPending(true);
    setError(null);
    setNotice(null);
    const data = new FormData(event.currentTarget);
    const payload: Record<string, unknown> = {};
    for (const [key, value] of data.entries()) {
      payload[key] = String(value);
    }
    if (typeof payload.userIds === "string") {
      payload.userIds = payload.userIds.split(",").map((item) => item.trim()).filter(Boolean);
    }
    if (typeof payload.orgIds === "string") {
      payload.orgIds = payload.orgIds.split(",").map((item) => item.trim()).filter(Boolean);
    }
    if (typeof payload.percentage === "string" && payload.percentage !== "") {
      payload.percentage = Number(payload.percentage);
    }
    if (typeof payload.recoveries === "string" && payload.recoveries !== "") {
      payload.recoveries = Number(payload.recoveries);
    }
    payload.enabled = payload.enabled === "true";
    payload.override = payload.override === "true";
    payload.identifiedCause = payload.identifiedCause === "true";
    if (typeof payload.body !== "string" && typeof payload.reason === "string") {
      payload.body = payload.reason;
    }
    const id = typeof payload.id === "string" ? payload.id : "";
    const target = pathFromId ? pathFromId.replace("{id}", encodeURIComponent(id)) : action;
    try {
      const response = await fetch(target, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      const body = (await response.json()) as { ok?: boolean; error?: { message?: string } };
      if (!response.ok || body.ok === false) {
        setError(body.error?.message ?? "The request failed.");
        return;
      }
      if (success) {
        setNotice(success);
      }
      router.refresh();
    } catch {
      setError("The operations service is unavailable.");
    } finally {
      setPending(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-3">
      {children}
      {error ? <p className="text-sm text-destructive">{error}</p> : null}
      {notice ? <p className="text-sm text-muted-foreground">{notice}</p> : null}
      {pathFromId ? null : (
        <Button type="submit" disabled={pending}>
          {pending ? "Saving…" : "Save"}
        </Button>
      )}
    </form>
  );
}
