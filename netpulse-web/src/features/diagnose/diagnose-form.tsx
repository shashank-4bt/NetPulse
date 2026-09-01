"use client";

import { useRouter } from "next/navigation";
import { useState, type FormEvent } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { startDiagnosis } from "@/features/diagnose/start-diagnosis";

const EXAMPLES = ["youtube.com", "google.com", "instagram.com"] as const;

type DiagnoseFormProps = {
  initialValue?: string;
};

export function DiagnoseForm({ initialValue = "" }: DiagnoseFormProps) {
  const router = useRouter();
  const [value, setValue] = useState(initialValue);
  const [error, setError] = useState<string | null>(null);
  const [pending, setPending] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPending(true);
    setError(null);
    try {
      const result = await startDiagnosis(value);
      if (!result.ok) {
        setError(result.error);
        return;
      }
      router.push(`/diagnose?run=${result.reportId}`);
    } catch {
      setError(
        "The diagnose service is unavailable. No measurement was started."
      );
    } finally {
      setPending(false);
    }
  }

  return (
    <div className="space-y-4">
      <form
        onSubmit={onSubmit}
        className="flex flex-col gap-3 sm:flex-row sm:items-end"
      >
        <div className="min-w-0 flex-1">
          <label htmlFor="diagnose-target" className="text-sm font-medium">
            Domain, URL, or known service
          </label>
          <Input
            id="diagnose-target"
            name="target"
            value={value}
            onChange={(event) => setValue(event.target.value)}
            placeholder="youtube.com"
            autoComplete="off"
            spellCheck={false}
            aria-invalid={error ? true : undefined}
            aria-describedby={error ? "diagnose-error" : "diagnose-help"}
            className="mt-1"
          />
        </div>
        <Button type="submit" disabled={pending}>
          {pending ? "Starting…" : "Check My Internet"}
        </Button>
      </form>
      <div className="flex flex-wrap gap-2">
        {EXAMPLES.map((hostname) => (
          <Button
            key={hostname}
            type="button"
            variant="outline"
            onClick={() => setValue(hostname)}
          >
            {hostname}
          </Button>
        ))}
      </div>
      <p id="diagnose-help" className="text-sm text-muted-foreground">
        Accepts a public hostname, http(s) URL, or a catalog service. Localhost,
        private networks, credentials, and non-http schemes are rejected. No
        probe is sent until workers exist.
      </p>
      {error ? (
        <p id="diagnose-error" className="text-sm text-destructive" role="alert">
          {error}
        </p>
      ) : null}
    </div>
  );
}
