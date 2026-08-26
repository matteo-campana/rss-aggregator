import { NextResponse } from "next/server";

import { getCurrentUser } from "@/lib/api";
import { SESSION_COOKIE } from "@/lib/session";

// The API key is validated here and then stored in an httpOnly cookie. It is
// never handed to the browser, so an XSS on the page cannot read it, and the
// browser only ever talks to this origin — no CORS involved.
export async function POST(request: Request) {
  let apiKey = "";

  try {
    const body = (await request.json()) as { apiKey?: unknown };

    apiKey = typeof body.apiKey === "string" ? body.apiKey.trim() : "";
  } catch {
    return NextResponse.json({ error: "malformed request body" }, { status: 400 });
  }

  if (!apiKey) {
    return NextResponse.json({ error: "the API key is required" }, { status: 400 });
  }

  try {
    await getCurrentUser(apiKey);
  } catch {
    return NextResponse.json({ error: "the API key was rejected by the backend" }, { status: 401 });
  }

  const response = NextResponse.json({ ok: true });

  response.cookies.set({
    name: SESSION_COOKIE,
    value: apiKey,
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: 60 * 60 * 24 * 30,
  });

  return response;
}

export async function DELETE() {
  const response = NextResponse.json({ ok: true });

  response.cookies.delete(SESSION_COOKIE);

  return response;
}
