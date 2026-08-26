# RSS Aggregator — frontend

Next.js reader for the items collected by the Go backend. See the repository README for
the full setup; in short:

```bash
npm install
cp .env.example .env.local     # API_BASE_URL points at the Go API
npm run dev
```

Sign in with the API key returned by `POST /api/v1/users/`.

## How it is put together

- **The API key never reaches the browser.** `app/api/session/route.ts` validates it against
  `GET /users/me` and stores it in an httpOnly cookie; every call to the Go API is made from a
  server component through `lib/api.ts`. The browser only talks to this origin, so the Go API
  needs no CORS configuration.
- **Filter state lives in the URL.** `app/page.tsx` is a server component that reads its
  `searchParams`, so a filtered view can be linked and shared and no client-side store is needed.
  `lib/search-params.ts` holds the parsing and is covered by unit tests.
- `middleware.ts` redirects to `/login` when the session cookie is missing.
- `/feeds` lists the configured feeds and when the scraper last refreshed each one. It is
  read-only: creating and editing feeds still happens through the API.

## Scripts

| Command | What it does |
| --- | --- |
| `npm run dev` | development server |
| `npm run build` | production build |
| `npm run lint` | ESLint |
| `npm run typecheck` | `tsc --noEmit` |
| `npm test` | vitest |
