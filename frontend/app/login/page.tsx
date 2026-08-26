import { LoginForm } from "@/components/login-form";

export const metadata = {
  title: "Sign in · RSS Aggregator",
};

export default function LoginPage() {
  return (
    <main className="mx-auto flex min-h-screen max-w-md flex-col justify-center px-6">
      <h1 className="text-2xl font-semibold tracking-tight">RSS Aggregator</h1>
      <p className="mt-2 text-sm text-neutral-600 dark:text-neutral-400">
        Sign in with the API key returned when your user was created.
      </p>

      <LoginForm />

      <p className="mt-8 text-xs text-neutral-500">
        The key is kept in an httpOnly cookie and used only server-side: it is never exposed to the
        page.
      </p>
    </main>
  );
}
