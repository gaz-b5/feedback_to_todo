"use client";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import React, { useState, useTransition } from "react";
import { useRouter } from "next/navigation";
import { Alert, AlertDescription } from "@/components/ui/alert";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL as string;

export function PassResetForm({
  className,
  ...props
}: React.ComponentProps<"form">) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [pending, startTransition] = useTransition();
  const [email, setEmail] = useState("");
  const [disabled, setDisabled] = useState(false);

  async function onSendReset(formData?: FormData) {
    setError(null);
    setSuccess(null);

    let targetEmail = email;
    if (!disabled && formData) {
      targetEmail = (formData.get("email") as string) || "";
      setEmail(targetEmail);
    }

    if (!targetEmail) {
      setError("Email is required");
      return;
    }

    try {
      const res = await fetch(
        `${API_BASE_URL}/collections/users/request-password-reset`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email: targetEmail }),
        }
      );

      let data: any = {};
      if (res.headers.get("content-type")?.includes("application/json")) {
        data = await res.json().catch(() => ({})); // ignore parse error if no body
      }

      if (!res.ok) {
        setError(data.message || "Failed to send password reset email");
        return;
      }

      setSuccess("Password reset email sent! Please check your inbox.");
      setDisabled(true);
    } catch (err) {
      // This only runs if fetch completely fails (network error)
      console.error("Network or fetch error:", err);
      setError("Network error while sending password reset email");
    }
  }


  return (
    <>
      <form
        className={cn("flex flex-col gap-6 mb-8", className)}
        action={(formData: FormData) =>
          startTransition(() => onSendReset(formData))
        }
        {...props}
      >
        <div className="flex flex-col items-center gap-2 text-center">
          <h1 className="text-2xl font-bold">Reset password</h1>
          <p className="text-muted-foreground text-sm">
            Enter your email below to receive a password reset link.
          </p>
        </div>
        <div className="grid gap-6">
          <div className="grid gap-3">
            <Label htmlFor="email">Email</Label>
            <div className="flex flex-flow-row gap-2">
              <Input
                id="email"
                type="email"
                name="email"
                placeholder="m@example.com"
                className="flex-4"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={disabled}
              />
              <Button
                type="submit"
                className="flex-1"
                onClick={(e) => {
                  // If already locked, prevent form from recreating FormData
                  if (disabled) {
                    e.preventDefault();
                    startTransition(() => onSendReset());
                  }
                }}
                disabled={pending}
              >
                {disabled ? "Resend" : "Send"}
              </Button>
            </div>
          </div>
          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
          {success && (
            <Alert variant="default">
              <AlertDescription>{success}</AlertDescription>
            </Alert>
          )}
        </div>
      </form>
    </>
  );
}
