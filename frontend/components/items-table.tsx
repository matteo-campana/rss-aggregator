import { Item } from "@/lib/types";

function formatDate(value: string | null): string {
  if (!value) {
    return "—";
  }

  const parsed = new Date(value);

  if (Number.isNaN(parsed.getTime())) {
    return "—";
  }

  return parsed.toISOString().slice(0, 16).replace("T", " ");
}

function formatCount(value: number | null): string {
  // null means the column is empty, which is not the same as zero.
  return value === null ? "—" : value.toLocaleString("en-US");
}

export function ItemsTable({ items }: { items: Item[] }) {
  return (
    <div className="overflow-x-auto rounded border border-neutral-200 dark:border-neutral-800">
      <table className="w-full min-w-[56rem] border-collapse text-sm">
        <thead className="bg-neutral-50 text-left dark:bg-neutral-900">
          <tr>
            <th className="px-3 py-2 font-medium">Title</th>
            <th className="px-3 py-2 font-medium">Source</th>
            <th className="px-3 py-2 text-right font-medium">Seeders</th>
            <th className="px-3 py-2 text-right font-medium">Leechers</th>
            <th className="px-3 py-2 text-right font-medium">Downloads</th>
            <th className="px-3 py-2 font-medium">Size</th>
            <th className="px-3 py-2 font-medium">Published</th>
          </tr>
        </thead>

        <tbody>
          {items.map((item) => (
            <tr key={item.id} className="border-t border-neutral-200 dark:border-neutral-800">
              <td className="max-w-md px-3 py-2">
                {item.link ? (
                  <a
                    href={item.link}
                    className="underline decoration-neutral-400 underline-offset-2"
                    rel="noreferrer noopener"
                  >
                    {item.title ?? item.guid}
                  </a>
                ) : (
                  (item.title ?? item.guid)
                )}
                {item.category ? (
                  <span className="mt-0.5 block text-xs text-neutral-500">{item.category}</span>
                ) : null}
              </td>
              <td className="px-3 py-2 text-neutral-600 dark:text-neutral-400">{item.channel_title}</td>
              <td className="px-3 py-2 text-right tabular-nums">{formatCount(item.seeders)}</td>
              <td className="px-3 py-2 text-right tabular-nums">{formatCount(item.leechers)}</td>
              <td className="px-3 py-2 text-right tabular-nums">{formatCount(item.downloads)}</td>
              <td className="px-3 py-2 whitespace-nowrap">{item.size ?? "—"}</td>
              <td className="px-3 py-2 whitespace-nowrap text-neutral-600 dark:text-neutral-400">
                {formatDate(item.published_at)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
