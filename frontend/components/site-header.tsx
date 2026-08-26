import Link from "next/link";

import { SignOutButton } from "@/components/sign-out-button";

type Props = {
  title: string;
  subtitle: string;
  current: "items" | "feeds";
};

export function SiteHeader({ title, subtitle, current }: Props) {
  const linkClass = "text-sm underline-offset-4";
  const active = "font-medium underline";
  const inactive = "text-neutral-600 hover:underline dark:text-neutral-400";

  return (
    <header className="flex flex-wrap items-baseline justify-between gap-4">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">{title}</h1>
        <p className="mt-1 text-sm text-neutral-600 dark:text-neutral-400">{subtitle}</p>
      </div>

      <nav className="flex items-center gap-4">
        <Link href="/" className={`${linkClass} ${current === "items" ? active : inactive}`}>
          Items
        </Link>
        <Link href="/feeds" className={`${linkClass} ${current === "feeds" ? active : inactive}`}>
          Feeds
        </Link>
        <SignOutButton />
      </nav>
    </header>
  );
}
