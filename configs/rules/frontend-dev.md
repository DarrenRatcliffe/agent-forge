---
paths:
  - "src/**/*.tsx"
  - "src/**/*.jsx"
  - "app/**/*.tsx"
  - "components/**/*.tsx"
---

# Frontend Rules — React / Next.js

## Stack assumptions
- TypeScript, strict mode. No `any` without a comment saying why.
- Function components and hooks only. No class components in new code.
- Match the framework version and router (App Router vs Pages) the project already uses. Don't mix them.

## Server vs Client Components
- Server Components by default. Add `"use client"` only for state, effects, event handlers, or browser APIs.
- Push `"use client"` as far down the tree as possible. A client boundary on a page-level component is a finding.
- Don't pass non-serialisable props (functions, class instances) from a Server to a Client Component.

## Structure
- One component per file, named in PascalCase to match the file.
- Keep components small. If one mixes data fetching, layout, and logic, split it.
- Shared logic goes in hooks or utilities, not copied between components.

## Styling
- Use the styling approach already in the codebase (Tailwind, CSS Modules, etc.). Don't introduce a second one.
- No inline styles except for values computed at runtime.

## State and data
- Prefer server data fetching and URL state over client state. Reach for global state only when two distant components genuinely share it.
- Handle loading, empty, and error states for every async view. A component that only renders the happy path is incomplete.
- No data fetching inside `useEffect` when a Server Component or the project's data library can do it.

## Accessibility
- Semantic HTML first (`button`, `nav`, `main`, `label`). A clickable `div` is a finding.
- Every interactive element is keyboard reachable with a visible focus state.
- Images have `alt` text; decorative ones use `alt=""`. Form inputs have associated labels.
- Don't convey meaning by colour alone. Check contrast.

## Security
- No `dangerouslySetInnerHTML` with unsanitised input.
- No secrets in client code. Only variables intentionally prefixed `NEXT_PUBLIC_` reach the browser.

## Performance
- Use `next/image` and `next/link` rather than raw `img` and `a` for internal content.
- Lists need stable `key` props, never array indexes for reorderable data.
- Don't memoise by default. Add `useMemo`/`useCallback` when a measured problem calls for it.

## Testing
- React Testing Library for component behaviour: test what the user sees and does, not implementation details.
- Cover loading, error, and empty states, not just the populated one.
- Tests run in CI without manual setup.
