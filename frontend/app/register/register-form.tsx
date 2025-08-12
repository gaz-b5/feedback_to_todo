"use client";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import React, { useState, useTransition } from "react";
import { handleRegister } from "@/app/register/actions";
import { useRouter } from "next/navigation";

export function RegisterForm({
  className,
  ...props
}: React.ComponentProps<"form">) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [pending, startTransition] = useTransition();

  async function onAction(formData: FormData) {
    setError(null);
    const result = await handleRegister(formData);
    if (result?.error) {
      setError(result.error);
      return;
    }
    router.push("/login");
  }

  return (
    <form
      className={cn("flex flex-col gap-6", className)}
      action={(formData: FormData) => startTransition(() => onAction(formData))}
      {...props}
    >
      <div className="flex flex-col items-center gap-2 text-center">
        <h1 className="text-2xl font-bold">Register new account</h1>
        <p className="text-muted-foreground text-sm text-balance">
          Enter your email below to register
        </p>
      </div>
      <div className="grid gap-6">
        <div className="grid gap-3">
          <Label htmlFor="name">Name</Label>
          <Input id="name" type="text" name="name" placeholder="Johnny Depp" required />
        </div>
        <div className="grid gap-3">
          <Label htmlFor="email">Email</Label>
          <Input id="email" type="email" name="identity" placeholder="m@example.com" required />
        </div>
        <div className="grid gap-3">
          <Label htmlFor="password">Password</Label>
          <Input id="password" type="password" name="password" required />
        </div>
        <div className="grid gap-3">
          <Label htmlFor="passwordConfirm">Confirm password</Label>
          <Input id="passwordConfirm" type="password" name="passwordConfirm" required />
        </div>
        <Button type="submit" className="w-full" disabled={pending}>
          {pending ? "Registering..." : "Register"}
        </Button>
        {error && <div className="text-red-600 mt-2">{error}</div>}
      </div>
      <div className="text-center text-sm">
        Already have an account?{" "}
        <a href="/login" className="underline underline-offset-4">
          Log in
        </a>
      </div>

    </form>
  );
}
