import { GalleryVerticalEnd } from "lucide-react"

import { RegisterForm } from "@/app/register/register-form"

export default function LoginPage() {
  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      <div className="flex flex-col gap-4 p-6 md:p-10">
        <div className="flex justify-center gap-2 md:justify-start">
          <a href="#" className="flex items-center gap-2 font-medium">
            <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
              <GalleryVerticalEnd className="size-4" />
            </div>
            Acme Inc.
          </a>
        </div>
        <div className="flex flex-1 items-center justify-center flex-col gap-6">
          <div className="w-full max-w-xs">
            <RegisterForm />
          </div>
          <div className="text-sm text-muted-foreground text-center">
            When you click the "Register" button, an email will be sent to the provided address with a <span className="underline underline-offset-4">verification link. Please check your inbox and click the link to verify your account</span> before Logging in.
          </div>

        </div>
      </div>
      <div className="bg-muted relative hidden lg:block">
        <img
          src="https://picsum.photos/1080/1080"
          alt="Image"
          className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
        />
      </div>
    </div>
  )
}
