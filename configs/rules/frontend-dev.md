---
paths:
  - "**/*.{tsx,jsx}"
  - "**/{components,pages,app}/**/*.{ts,js,css,scss}"
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

## Accessibility (WCAG 2.2 AA)

The target is WCAG 2.2 AA. These rules are non-negotiable; the `frontend-accessibility` skill (`~/.claude/skills/frontend-accessibility/`) has the detail, component patterns, and testing checklist.
- Semantic HTML first (`button`, `nav`, `main`, `label`, `form`); `role` and `aria-*` only when HTML can't express it. A clickable `div` with `onClick` is a finding.
- Every interactive element works with the keyboard alone, in visual order, with no keyboard traps (Escape closes dialogs). No positive `tabindex` values.
- Every form control has an associated `<label>` or `aria-label`. Images have `alt` text, or `alt=""` if decorative.
- Contrast: normal text ≥ 4.5:1, large text ≥ 3:1, UI components and focus indicators ≥ 3:1. Never convey information by colour alone.
- Visible focus indicator on every focusable element. Don't remove the outline without a replacement.

## Security
- No `dangerouslySetInnerHTML` with unsanitised input.
- No secrets in client code. Only variables intentionally prefixed `NEXT_PUBLIC_` reach the browser.
- Sanitise user input in forms (use libraries like DOMPurify if rendering user HTML).
- CSP (Content Security Policy) headers configured (if available through server).

## Performance
- Use `next/image` and `next/link` rather than raw `img` and `a` for internal content.
- Lists need stable `key` props, never array indexes for reorderable data.
- Don't memoise by default. Add `useMemo`/`useCallback` when a measured problem calls for it.
- Core Web Vitals targets: LCP (≤ 2.5s), INP (≤ 200ms), CLS (≤ 0.1). INP replaced FID in 2024. Measure with Lighthouse or Web Vitals library.

## Testing

### Test approach
- React Testing Library: test what the user sees and does, not implementation details. Don't test internal state or component structure.
- Behavior over implementation: test that button click shows modal, not that onClick handler was called.
- User perspective: query by label, placeholder, role, text (in that order). Avoid query by test ID unless text/role/label insufficient.

### Coverage expectations
- Cover main happy path and primary user workflow.
- Cover loading states (skeleton, spinner, etc.) and empty states (no results).
- Cover error states (error message, retry option).
- Cover edge cases specific to component (disabled, focused, selected).
- Don't test third-party library internals (assume library works; test your usage).

### Test organization
- One test file per component (co-locate with component, e.g., `Button.tsx` and `Button.test.tsx`).
- Test naming: describe what a user can do. Examples: "displays button with text", "submits form on click", "shows error when API fails".
- Use `describe()` blocks to group related tests.
- Setup/teardown: mocks, fixtures in beforeEach; clean up in afterEach.

### Fixtures and mocks
- Mock external dependencies (API calls, third-party libraries).
- Use real component rendering, not shallow/mount (React Testing Library enforces this).
- Fixtures for repeated data (mock user object, mock API response).
- Don't over-mock: if component uses native browser APIs correctly, don't mock the browser (render with jsdom and test real behavior).

### Accessibility in tests
- Use `screen.getByRole()` to query elements (encourages accessible markup).
- Test keyboard navigation: `userEvent.tab()` to verify focus moves correctly.
- Test screen reader labels: `screen.getByLabelText()`, `screen.getByAltText()` for images.
- Avoid brittle selectors: don't query by CSS class, test ID unless necessary.

### Performance testing
- Measure component render time for performance-critical components (large lists, complex rendering).
- Use React Testing Library's `render()` and browser devtools `Performance` tab to identify slow renders.
- Don't over-optimize: profile first, optimize only slow components.

### Test tools
- Jest (test runner, assertions, mocking).
- React Testing Library (component testing, realistic user interactions).
- @testing-library/user-event (simulate user actions: type, click, tab).
- Mock Service Worker (MSW) for mocking API calls (preferred over fetch mocks).
- Vitest (modern Jest alternative, faster, better ESM support).
