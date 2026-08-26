import { redirect } from "next/navigation";

import { FeedsTable } from "@/components/feeds-table";
import { Pagination } from "@/components/pagination";
import { SiteHeader } from "@/components/site-header";
import { ApiError, listFeeds } from "@/lib/api";
import { parseFilters, type RawSearchParams } from "@/lib/search-params";
import { getApiKey } from "@/lib/session";

export const metadata = {
  title: "Feeds · RSS Aggregator",
};

// Read-only on purpose: creating and editing feeds still happens through the API.
export default async function FeedsPage({
  searchParams,
}: {
  searchParams: Promise<RawSearchParams>;
}) {
  const apiKey = await getApiKey();

  if (!apiKey) {
    redirect("/login");
  }

  const filters = parseFilters(await searchParams);

  let page;

  try {
    page = await listFeeds(apiKey, filters.page, filters.perPage);
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      redirect("/login");
    }

    throw error;
  }

  return (
    <main className="mx-auto flex max-w-6xl flex-col gap-6 px-6 py-10">
      <SiteHeader
        title="Feeds"
        subtitle={`${page.total.toLocaleString("en-US")} feed${page.total === 1 ? "" : "s"} configured.`}
        current="feeds"
      />

      {page.feeds.length === 0 ? (
        <div className="rounded border border-dashed border-neutral-300 p-10 text-center dark:border-neutral-700">
          <p className="text-sm font-medium">No feeds yet.</p>
          <p className="mx-auto mt-2 max-w-prose text-sm text-neutral-600 dark:text-neutral-400">
            Add one with <code>POST /api/v1/feeds/</code>, then run the API with{" "}
            <code>SCRAPER_ENABLED=true</code> to start collecting its items.
          </p>
        </div>
      ) : (
        <>
          <FeedsTable feeds={page.feeds} />
          <Pagination filters={filters} total={page.total} basePath="/feeds" />
        </>
      )}
    </main>
  );
}
