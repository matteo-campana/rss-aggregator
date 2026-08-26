"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

export function SignOutButton() {
  const router = useRouter();
  const [pending, setPending] = useState(false);

  async function onClick() {
    setPending(true);

    await fetch("/api/session", { method: "DELETE" });

    router.replace("/login");
    router.refresh();
  }

  return (
    <button
      type="button"
      onClick={onClick}
      disabled={pending}
      className="text-sm text-neutral-600 underline underline-offset-2 disabled:opacity-50 dark:text-neutral-400"
    >
      Sign out
    </button>
  );
}
