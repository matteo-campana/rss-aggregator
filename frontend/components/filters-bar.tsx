"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";

import { ItemFilters, SORTS } from "@/lib/search-params";
import { Channel } from "@/lib/types";

type Props = {
  filters: ItemFilters;
  categories: string[];
  channels: Channel[];
};

// Submitting rewrites the URL: the page is a server component, so the new
// search params are what triggers the refetch. Page resets to 1 on any change.
export function FiltersBar({ filters, categories, channels }: Props) {
  const router = useRouter();
  const searchParams = useSearchParams();

  function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const form = new FormData(event.currentTarget);
    const query = new URLSearchParams();

    for (const [key, value] of form.entries()) {
      const trimmed = String(value).trim();

      if (trimmed !== "") {
        query.set(key, trimmed);
      }
    }

    router.push(query.toString() ? `/?${query.toString()}` : "/");
  }

  const hasFilters = searchParams.toString() !== "";

  return (
    <form
      onSubmit={onSubmit}
      className="flex flex-wrap items-end gap-3 rounded border border-neutral-200 p-4 dark:border-neutral-800"
    >
      <div className="flex flex-col gap-1">
        <label htmlFor="search" className="text-xs font-medium text-neutral-600 dark:text-neutral-400">
          Title contains
        </label>
        <input
          id="search"
          name="search"
          defaultValue={filters.search}
          placeholder="SubsPlease"
          className="w-52 rounded border border-neutral-300 px-2 py-1.5 text-sm dark:border-neutral-700 dark:bg-neutral-900"
        />
      </div>

      <div className="flex flex-col gap-1">
        <label htmlFor="category" className="text-xs font-medium text-neutral-600 dark:text-neutral-400">
          Category
        </label>
        <select
          id="category"
          name="category"
          defaultValue={filters.category}
          className="w-56 rounded border border-neutral-300 px-2 py-1.5 text-sm dark:border-neutral-700 dark:bg-neutral-900"
        >
          <option value="">All</option>
          {categories.map((category) => (
            <option key={category} value={category}>
              {category}
            </option>
          ))}
        </select>
      </div>

      <div className="flex flex-col gap-1">
        <label htmlFor="channelId" className="text-xs font-medium text-neutral-600 dark:text-neutral-400">
          Source
        </label>
        <select
          id="channelId"
          name="channelId"
          defaultValue={filters.channelId}
          className="w-56 rounded border border-neutral-300 px-2 py-1.5 text-sm dark:border-neutral-700 dark:bg-neutral-900"
        >
          <option value="">All</option>
          {channels.map((channel) => (
            <option key={channel.id} value={channel.id}>
              {channel.title ?? channel.id}
            </option>
          ))}
        </select>
      </div>

      <div className="flex flex-col gap-1">
        <label htmlFor="minSeeders" className="text-xs font-medium text-neutral-600 dark:text-neutral-400">
          Min seeders
        </label>
        <input
          id="minSeeders"
          name="minSeeders"
          type="number"
          min={0}
          defaultValue={filters.minSeeders ?? ""}
          className="w-28 rounded border border-neutral-300 px-2 py-1.5 text-sm dark:border-neutral-700 dark:bg-neutral-900"
        />
      </div>

      <div className="flex flex-col gap-1">
        <label htmlFor="sort" className="text-xs font-medium text-neutral-600 dark:text-neutral-400">
          Sort by
        </label>
        <select
          id="sort"
          name="sort"
          defaultValue={filters.sort}
          className="w-36 rounded border border-neutral-300 px-2 py-1.5 text-sm dark:border-neutral-700 dark:bg-neutral-900"
        >
          {SORTS.map((sort) => (
            <option key={sort} value={sort}>
              {sort === "recent" ? "Newest" : sort === "oldest" ? "Oldest" : "Most seeders"}
            </option>
          ))}
        </select>
      </div>

      <button
        type="submit"
        className="rounded bg-neutral-900 px-3 py-1.5 text-sm font-medium text-white dark:bg-white dark:text-neutral-900"
      >
        Apply
      </button>

      {hasFilters ? (
        <Link href="/" className="px-1 py-1.5 text-sm text-neutral-600 underline dark:text-neutral-400">
          Reset
        </Link>
      ) : null}
    </form>
  );
}
