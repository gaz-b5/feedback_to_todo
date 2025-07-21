"use client";
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useRouter } from "next/navigation";


const API_BASE_URL: string = process.env.NEXT_PUBLIC_API_BASE_URL as string;

export function RegisterForm({
  className,
  ...props
}: React.ComponentProps<"form">) {

  const router  = useRouter();

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const form = event.currentTarget;
    const formData = new FormData(form);

    const obj = {name:formData.get("name"), email:formData.get("identity"), password:formData.get("password"), passwordConfirm:formData.get("passwordConfirm")};

    try {
      const response = await fetch(`${API_BASE_URL}/collections/users/records`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(obj)
      });
      if (!response.ok) throw new Error("Login failed");
      // handle login success here
      router.push("/login");
      alert("Register successful");
    } catch (err) {
      alert((err as Error).message);
    }
  }

  return (
    <form className={cn("flex flex-col gap-6", className)} {...props} onSubmit={handleSubmit}>
      <div className="flex flex-col items-center gap-2 text-center">
        <h1 className="text-2xl font-bold">Register new account</h1>
        <p className="text-muted-foreground text-sm text-balance">
          Enter your email below to register
        </p>
      </div>
      <div className="grid gap-6">
        <div className="grid gap-3">
          <Label htmlFor="email">Name</Label>
          <Input id="name" type="name" name = "name" placeholder="Johnny Depp" required />
        </div>
        <div className="grid gap-3">
          <Label htmlFor="email">Email</Label>
          <Input id="email" type="email" name = "identity" placeholder="m@example.com" required />
        </div>
        <div className="grid gap-3">
          <div className="flex items-center">
            <Label htmlFor="password" >Password</Label>
          </div>
          <Input id="password" type="password" name = "password" required />
        </div>
        <div className="grid gap-3">
          <div className="flex items-center">
            <Label htmlFor="password" >Confirm password</Label>
          </div>
          <Input id="passwordConfirm" type="password" name = "passwordConfirm" required />
        </div>
        <Button type="submit" className="w-full">
          Register
        </Button>
      </div>
    </form>
  )
}
