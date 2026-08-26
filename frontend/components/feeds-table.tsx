import { Feed } from "@/lib/types";

function formatDate(value: string | null): string {
  if (!value) {
    return "never";
  }

  const parsed = new Date(value);

  if (Number.isNaN(parsed.getTime())) {
    return "never";
  }

  return parsed.toISOString().slice(0, 16).replace("T", " ");
}

export function FeedsTable({ feeds }: { feeds: Feed[] }) {
  return (
    <div className="overflow-x-auto rounded border border-neutral-200 dark:border-neutral-800">
      <table className="w-full min-w-[44rem] border-collapse text-sm">
        <thead className="bg-neutral-50 text-left dark:bg-neutral-900">
          <tr>
            <th className="px-3 py-2 font-medium">Name</th>
            <th className="px-3 py-2 font-medium">URL</th>
            <th className="px-3 py-2 font-medium">Last fetched</th>
          </tr>
        </thead>

        <tbody>
          {feeds.map((feed) => (
            <tr key={feed.id} className="border-t border-neutral-200 dark:border-neutral-800">
              <td className="px-3 py-2 font-medium">{feed.name}</td>
              <td className="max-w-md truncate px-3 py-2 text-neutral-600 dark:text-neutral-400">
                {feed.url}
              </td>
              <td className="px-3 py-2 whitespace-nowrap text-neutral-600 dark:text-neutral-400">
                {formatDate(feed.last_fetched_at)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
