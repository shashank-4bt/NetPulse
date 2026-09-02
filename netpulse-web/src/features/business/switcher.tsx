"use client";

import { useRouter } from "next/navigation";
import { useState, type FormEvent } from "react";

import type { OrganizationView } from "@/domain/business";
import { roleLabel } from "@/domain/business";
import { Button } from "@/components/ui/button";

export function OrgSwitcher({
  organizations,
  selected,
}: {
  organizations: OrganizationView[];
  selected: string;
}) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const orgId = String(new FormData(event.currentTarget).get("orgId") ?? "");
    const response = await fetch("/api/orgs/select", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ orgId }),
    });
    if (!response.ok) {
      setError("That organization is not available.");
      return;
    }
    setError(null);
    router.refresh();
  }

  if (!organizations.length) {
    return null;
  }

  return (
    <form onSubmit={onSubmit} className="flex flex-wrap items-end gap-2">
      <div>
        <label htmlFor="orgId" className="text-sm font-medium">
          Organization
        </label>
        <select
          id="orgId"
          name="orgId"
          defaultValue={selected}
          className="mt-1 h-11 min-w-56 rounded-lg border border-input bg-transparent px-3 text-sm"
        >
          {organizations.map((item) => (
            <option key={item.id} value={item.id}>
              {item.name} · {roleLabel(item.role)}
            </option>
          ))}
        </select>
      </div>
      <Button type="submit" variant="outline">
        Use organization
      </Button>
      {error ? <p className="text-sm text-destructive">{error}</p> : null}
    </form>
  );
}
