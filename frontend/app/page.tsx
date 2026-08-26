import { redirect } from "next/navigation";

import { EmptyState } from "@/components/empty-state";
import { FiltersBar } from "@/components/filters-bar";
import { ItemsTable } from "@/components/items-table";
import { Pagination } from "@/components/pagination";
import { SiteHeader } from "@/components/site-header";
import { ApiError, listCategories, listChannels, listItems } from "@/lib/api";
import { parseFilters, type RawSearchParams } from "@/lib/search-params";
import { getApiKey } from "@/lib/session";

export const metadata = {
  title: "Items · RSS Aggregator",
};

export default async function ItemsPage({
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
  let categories: string[] = [];
  let channels = [];

  try {
    [page, categories, channels] = await Promise.all([
      listItems(apiKey, filters),
      listCategories(apiKey),
      listChannels(apiKey),
    ]);
  } catch (error) {
    // A revoked or edited cookie must send the visitor back to the login page
    // rather than showing a stack trace.
    if (error instanceof ApiError && error.status === 401) {
      redirect("/login");
    }

    throw error;
  }

  const isFiltered =
    filters.search !== "" ||
    filters.category !== "" ||
    filters.channelId !== "" ||
    filters.minSeeders !== null;

  return (
    <main className="mx-auto flex max-w-6xl flex-col gap-6 px-6 py-10">
      <SiteHeader
        title="Aggregated items"
        subtitle={`${page.total.toLocaleString("en-US")} item${page.total === 1 ? "" : "s"} stored${
          isFiltered ? " matching these filters" : ""
        }.`}
        current="items"
      />

      <FiltersBar filters={filters} categories={categories} channels={channels} />

      {page.items.length === 0 ? (
        <EmptyState filtered={isFiltered} />
      ) : (
        <>
          <ItemsTable items={page.items} />
          <Pagination filters={filters} total={page.total} />
        </>
      )}
    </main>
  );
}
