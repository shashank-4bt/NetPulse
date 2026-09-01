"use client";

import { useState, type FormEvent } from "react";

import { UnavailableState } from "@/components/feedback/unavailable-state";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { validateDiagnoseTarget } from "@/features/diagnose/validate-target";

export function DiagnoseForm() {
  const [value, setValue] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<string | null>(null);

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const validation = validateDiagnoseTarget(value);
    if (!validation.ok) {
      setResult(null);
      setError(validation.error);
      return;
    }
    setError(null);
    setResult(validation.hostname);
  }

  return (
    <div className="space-y-4">
      <form
        onSubmit={onSubmit}
        className="flex flex-col gap-3 sm:flex-row sm:items-end"
      >
        <div className="min-w-0 flex-1">
          <label htmlFor="diagnose-target" className="text-sm font-medium">
            Hostname or URL
          </label>
          <Input
            id="diagnose-target"
            name="target"
            value={value}
            onChange={(event) => setValue(event.target.value)}
            placeholder="example.com"
            autoComplete="off"
            spellCheck={false}
            aria-invalid={error ? true : undefined}
            aria-describedby={error ? "diagnose-error" : "diagnose-help"}
            className="mt-1"
          />
        </div>
        <Button type="submit">Check My Internet</Button>
      </form>
      <p id="diagnose-help" className="text-sm text-muted-foreground">
        No probe is sent from this form until measurement workers are connected.
        Localhost and private addresses are rejected.
      </p>
      {error ? (
        <p id="diagnose-error" className="text-sm text-destructive" role="alert">
          {error}
        </p>
      ) : null}
      {result ? (
        <UnavailableState
          title="Measurement unavailable"
          description={`Target accepted: ${result}. The diagnostic API is not connected, so NetPulse did not run DNS, TLS, or HTTP probes and did not invent a result.`}
        />
      ) : null}
    </div>
  );
}
