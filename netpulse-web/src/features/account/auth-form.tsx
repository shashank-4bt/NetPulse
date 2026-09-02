"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, type ComponentProps, type FormEvent } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

type AuthFormProps = {
  mode: "login" | "register" | "forgot" | "reset" | "verify";
  nextPath?: string;
  token?: string;
};

export function AuthForm({ mode, nextPath = "/dashboard", token = "" }: AuthFormProps) {
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
    const payload: Record<string, string> = {};
    for (const [key, value] of data.entries()) {
      payload[key] = String(value);
    }
    if (token) {
      payload.token = token;
    }
    try {
      const path = endpoint(mode);
      const response = await fetch(path, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      const body = (await response.json()) as {
        ok?: boolean;
        error?: { message?: string };
        auth?: { emailReason?: string; devToken?: string };
      };
      if (!response.ok || body.ok === false) {
        setError(body.error?.message ?? "The request failed.");
        return;
      }
      if (mode === "login" || mode === "register") {
        router.push(nextPath);
        router.refresh();
        return;
      }
      if (mode === "forgot") {
        const extra = body.auth?.devToken
          ? ` Local reset token: ${body.auth.devToken}`
          : "";
        setNotice((body.auth?.emailReason ?? "If an account exists, a reset path was created.") + extra);
        return;
      }
      if (mode === "reset") {
        router.push("/login");
        return;
      }
      setNotice("Email verified.");
      router.push("/dashboard");
    } catch {
      setError("The account service is unavailable.");
    } finally {
      setPending(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-4">
      {mode === "register" ? (
        <Field id="displayName" label="Display name" name="displayName" autoComplete="nickname" />
      ) : null}
      {mode === "login" || mode === "register" || mode === "forgot" ? (
        <Field id="email" label="Email" name="email" type="email" autoComplete="email" required />
      ) : null}
      {mode === "login" || mode === "register" || mode === "reset" ? (
        <Field
          id="password"
          label={mode === "reset" ? "New password" : "Password"}
          name="password"
          type="password"
          autoComplete={mode === "login" ? "current-password" : "new-password"}
          required
          minLength={10}
        />
      ) : null}
      {mode === "verify" && !token ? (
        <Field id="token" label="Verification token" name="token" required />
      ) : null}
      {error ? (
        <p id="auth-error" className="text-sm text-destructive" role="alert">
          {error}
        </p>
      ) : null}
      {notice ? (
        <p className="text-sm text-muted-foreground" role="status">
          {notice}
        </p>
      ) : null}
      <Button type="submit" disabled={pending}>
        {pending ? "Working…" : actionLabel(mode)}
      </Button>
      <p className="text-sm text-muted-foreground">
        {mode === "login" ? (
          <>
            <Link href="/register" className="underline-offset-4 hover:underline">
              Create an account
            </Link>
            {" · "}
            <Link href="/forgot-password" className="underline-offset-4 hover:underline">
              Forgot password
            </Link>
          </>
        ) : (
          <Link href="/login" className="underline-offset-4 hover:underline">
            Sign in
          </Link>
        )}
      </p>
    </form>
  );
}

function Field({
  id,
  label,
  ...props
}: ComponentProps<typeof Input> & { id: string; label: string }) {
  return (
    <div>
      <label htmlFor={id} className="text-sm font-medium">
        {label}
      </label>
      <Input id={id} className="mt-1" {...props} />
    </div>
  );
}

function endpoint(mode: AuthFormProps["mode"]): string {
  switch (mode) {
    case "login":
      return "/api/auth/login";
    case "register":
      return "/api/auth/register";
    case "forgot":
      return "/api/auth/forgot-password";
    case "reset":
      return "/api/auth/reset-password";
    case "verify":
      return "/api/auth/verify-email";
  }
}

function actionLabel(mode: AuthFormProps["mode"]): string {
  switch (mode) {
    case "login":
      return "Sign in";
    case "register":
      return "Create account";
    case "forgot":
      return "Request reset";
    case "reset":
      return "Set new password";
    case "verify":
      return "Verify email";
  }
}
