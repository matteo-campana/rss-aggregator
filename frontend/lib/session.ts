import { cookies } from "next/headers";

export const SESSION_COOKIE = "rss_api_key";

/**
 * Reads the API key from the session cookie. The cookie is httpOnly, so the
 * key never reaches client-side JavaScript: every call to the Go API is made
 * from the server.
 */
export async function getApiKey(): Promise<string | null> {
  const store = await cookies();

  return store.get(SESSION_COOKIE)?.value ?? null;
}
