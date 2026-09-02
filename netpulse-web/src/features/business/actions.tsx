"use client";

import { useRouter } from "next/navigation";
import { useState, type FormEvent, type ReactNode } from "react";

import { BUSINESS_ROLES, REPORT_KINDS } from "@/domain/business";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function CreateOrgForm() {
  return (
    <BizForm action="/api/orgs" method="POST" success="Organization stored. You are the owner." selectCreatedOrg>
      <Field id="name" label="Organization name" required />
    </BizForm>
  );
}

export function RenameOrgForm({ orgId, name }: { orgId: string; name: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}`} method="PATCH" success="Organization name updated.">
      <Field id="name" label="Name" required defaultValue={name} />
    </BizForm>
  );
}

export function InviteMemberForm({ orgId }: { orgId: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/members`} method="POST" success="Invite stored. Existing accounts are added immediately.">
      <Field id="email" label="Email" type="email" required />
      <RoleSelect defaultValue="viewer" />
    </BizForm>
  );
}

export function ChangeRoleForm({ orgId, memberId, role }: { orgId: string; memberId: string; role: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/members/${memberId}`} method="PATCH" success="Role updated." hideSave>
      <RoleSelect defaultValue={role} id={`role-${memberId}`} />
      <Button type="submit" variant="outline" size="sm">
        Change role
      </Button>
    </BizForm>
  );
}

export function RemoveMemberButton({ orgId, memberId }: { orgId: string; memberId: string }) {
  return (
    <BizForm
      action={`/api/orgs/${orgId}/members/${memberId}`}
      method="DELETE"
      success="Member removed."
      confirm="Remove this member?"
      hideSave
    >
      <Button type="submit" variant="destructive" size="sm">
        Remove
      </Button>
    </BizForm>
  );
}

export function CreateTeamForm({ orgId }: { orgId: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/teams`} method="POST" success="Team stored.">
      <Field id="name" label="Team name" required />
      <Field id="memberIds" label="Member ids (comma-separated)" />
    </BizForm>
  );
}

export function DeleteTeamButton({ orgId, id }: { orgId: string; id: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/teams/${id}`} method="DELETE" success="Team deleted." confirm="Delete this team?" hideSave>
      <Button type="submit" variant="destructive" size="sm">
        Delete
      </Button>
    </BizForm>
  );
}

export function CreateDeviceForm({ orgId }: { orgId: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/devices`} method="POST" success="Device label stored. This is not a discovered host.">
      <Field id="name" label="Name" required />
      <Field id="label" label="Label" />
      <Field id="region" label="Region" />
    </BizForm>
  );
}

export function DeleteDeviceButton({ orgId, id }: { orgId: string; id: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/devices/${id}`} method="DELETE" success="Device removed." confirm="Remove this device?" hideSave>
      <Button type="submit" variant="destructive" size="sm">
        Remove
      </Button>
    </BizForm>
  );
}

export function CreateNetworkForm({ orgId }: { orgId: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/networks`} method="POST" success="Network record stored.">
      <Field id="name" label="Name" required />
      <Field id="asn" label="ASN" />
      <Field id="region" label="Region" />
    </BizForm>
  );
}

export function DeleteNetworkButton({ orgId, id }: { orgId: string; id: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/networks/${id}`} method="DELETE" success="Network removed." confirm="Remove this network?" hideSave>
      <Button type="submit" variant="destructive" size="sm">
        Remove
      </Button>
    </BizForm>
  );
}

export function CreateOrgServiceForm({ orgId }: { orgId: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/services`} method="POST" success="Service record stored.">
      <Field id="name" label="Name" required />
      <Field id="slug" label="Slug" />
      <Field id="endpoint" label="Endpoint" />
    </BizForm>
  );
}

export function DeleteOrgServiceButton({ orgId, id }: { orgId: string; id: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/services/${id}`} method="DELETE" success="Service removed." confirm="Remove this service?" hideSave>
      <Button type="submit" variant="destructive" size="sm">
        Remove
      </Button>
    </BizForm>
  );
}

export function CreateOrgMonitorForm({ orgId }: { orgId: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/monitors`} method="POST" success="Monitor stored. Status stays unmeasured until a check is recorded.">
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
    </BizForm>
  );
}

export function DeleteOrgMonitorButton({ orgId, id }: { orgId: string; id: string }) {
  return (
    <BizForm
      action={`/api/orgs/${orgId}/monitors/${id}`}
      method="DELETE"
      success="Monitor deleted."
      confirm="Delete this monitor?"
      hideSave
    >
      <Button type="submit" variant="destructive" size="sm">
        Delete
      </Button>
    </BizForm>
  );
}

export function ResolveIncidentButton({ orgId, id }: { orgId: string; id: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/incidents/${id}`} method="PATCH" success="Incident marked resolved." hideSave>
      <input type="hidden" name="status" value="resolved" />
      <Button type="submit" variant="outline" size="sm">
        Mark resolved
      </Button>
    </BizForm>
  );
}

export function GenerateReportForm({ orgId }: { orgId: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/reports`} method="POST" success="Report generated from stored records only.">
      <label htmlFor="kind" className="text-sm font-medium">
        Kind
      </label>
      <select id="kind" name="kind" required className="mt-1 h-11 w-full rounded-lg border border-input bg-transparent px-3 text-sm">
        {REPORT_KINDS.map((item) => (
          <option key={item.value} value={item.value}>
            {item.label}
          </option>
        ))}
      </select>
    </BizForm>
  );
}

export function CreateOrgKeyForm({ orgId }: { orgId: string }) {
  return (
    <BizForm
      action={`/api/orgs/${orgId}/keys`}
      method="POST"
      success="Key created. Copy the secret now; it is not stored in the browser."
      secretField="keySecret"
    >
      <Field id="name" label="Name" defaultValue="Organization key" />
      <Field id="scopes" label="Permissions" defaultValue="diagnosis.read,incident.read" />
      <Field id="rateLimitPerMin" label="Rate limit per minute" type="number" defaultValue="60" />
    </BizForm>
  );
}

export function RotateOrgKeyButton({ orgId, id }: { orgId: string; id: string }) {
  return (
    <BizForm
      action={`/api/orgs/${orgId}/keys/${id}/rotate`}
      method="POST"
      success="Key rotated. Copy the new secret now."
      secretField="keySecret"
      hideSave
    >
      <Button type="submit" variant="outline" size="sm">
        Rotate
      </Button>
    </BizForm>
  );
}

export function RevokeOrgKeyButton({ orgId, id }: { orgId: string; id: string }) {
  return (
    <BizForm action={`/api/orgs/${orgId}/keys/${id}/revoke`} method="POST" success="Key revoked." confirm="Revoke this organization API key?" hideSave>
      <Button type="submit" variant="destructive" size="sm">
        Revoke
      </Button>
    </BizForm>
  );
}

function RoleSelect({ defaultValue, id = "role" }: { defaultValue: string; id?: string }) {
  return (
    <div>
      <label htmlFor={id} className="text-sm font-medium">
        Role
      </label>
      <select id={id} name="role" defaultValue={defaultValue} className="mt-1 h-11 w-full rounded-lg border border-input bg-transparent px-3 text-sm">
        {BUSINESS_ROLES.map((item) => (
          <option key={item.value} value={item.value}>
            {item.label}
          </option>
        ))}
      </select>
    </div>
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

function BizForm({
  action,
  method,
  children,
  success,
  confirm,
  hideSave,
  secretField,
  selectCreatedOrg,
}: {
  action: string;
  method: string;
  children: ReactNode;
  success?: string;
  confirm?: string;
  hideSave?: boolean;
  secretField?: "keySecret";
  selectCreatedOrg?: boolean;
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
    const data = new FormData(event.currentTarget);
    const payload: Record<string, unknown> = {};
    for (const [key, value] of data.entries()) {
      payload[key] = String(value);
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
    if (typeof payload.memberIds === "string") {
      payload.memberIds = payload.memberIds
        .split(",")
        .map((item) => item.trim())
        .filter(Boolean);
    }
    if (typeof payload.rateLimitPerMin === "string" && payload.rateLimitPerMin !== "") {
      payload.rateLimitPerMin = Number(payload.rateLimitPerMin);
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
        organization?: { id?: string };
      };
      if (!response.ok || body.ok === false) {
        setError(body.error?.message ?? "The request failed.");
        return;
      }
      if (selectCreatedOrg && typeof body.organization?.id === "string") {
        await fetch("/api/orgs/select", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ orgId: body.organization.id }),
        });
      }
      if (secretField && typeof body[secretField] === "string") {
        setSecret(String(body[secretField]));
      }
      if (success) {
        setNotice(success);
      }
      router.refresh();
    } catch {
      setError("The organization service is unavailable.");
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
        <p className="break-all rounded-md border border-border bg-muted p-3 text-sm">Secret (shown once): {secret}</p>
      ) : null}
      {hideSave ? null : (
        <Button type="submit" disabled={pending}>
          {pending ? "Saving…" : "Save"}
        </Button>
      )}
    </form>
  );
}
