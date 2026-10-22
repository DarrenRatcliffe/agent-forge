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

## Accessibility (WCAG 2.2 AA)

### Semantic HTML and ARIA
- Semantic HTML first: `button`, `nav`, `main`, `section`, `article`, `header`, `footer`, `label`, `form`. A clickable `div` with `onClick` is a finding.
- Use `role` and `aria-*` attributes only when semantic HTML is insufficient. Example: custom combobox uses `role="combobox"`, `aria-expanded`, `aria-controls` (pointing at the listbox) and `aria-activedescendant` (the highlighted option), per the ARIA 1.2 combobox pattern.
- Landmarks: use `<nav>`, `<main>`, `<aside>` to structure page regions. Screen reader users navigate by landmark.
- Form inputs always have associated `<label>` or `aria-label`. Never orphan form controls.

### Keyboard navigation
- Every interactive element (button, link, input, custom control) is keyboard accessible.
- Tab order follows visual flow (left-to-right, top-to-bottom). Use `tabindex` only to repair broken order (never positive values).
- Visible focus indicator on all focusable elements. CSS outline must not be removed; if styling, keep contrast ≥ 3:1.
- Keyboard traps: user can't get focus stuck (e.g., modal dialog must allow Escape to close).
- Skip links: allow keyboard users to skip repetitive content (e.g., skip navigation to main content).

### Color and contrast
- Don't convey information by color alone. Examples: "Click red button" should be "Click the Delete button", error messages shown with icon+color+text.
- Text contrast: normal text ≥ 4.5:1, large text (18pt+) ≥ 3:1. Use a contrast checker (WebAIM, axe DevTools).
- UI component contrast: active/inactive states ≥ 3:1.
- Avoid color combinations that confuse red-blind or green-blind users (red/green, blue/yellow combinations).

### Images and media
- Images have `alt` text. Decorative images use `alt=""`. Alt text describes content, not "image of" (screen readers add "image").
- Complex images (charts, diagrams) have long description: `aria-describedby`, linked caption, or adjacent text.
- Avoid text in images (not readable by screen readers, fails on zoom). If text must be in image, include in `alt`.
- Video/audio: captions required (speech and relevant sounds), transcripts recommended.

### Motion and animations
- No auto-playing video/audio. User must start it.
- Animation that lasts > 5 seconds has pause/stop control.
- Avoid flashing content (>3x per second can trigger seizures). Use `prefers-reduced-motion` media query to disable or reduce animations for users who request it.

### Mobile and touch
- Touch targets: at least 24x24 CSS pixels, or enough spacing that a 24px circle around each target doesn't overlap another (WCAG 2.2 AA, 2.5.8). Aim for 44x44 where layout allows (good practice, and the AAA level in 2.5.5).
- Touch spacing: targets spaced apart (not cramped, easy to hit).
- Zoom support: text should enlarge when browser zoomed to 200%; don't disable zoom with `user-scalable=no`.

### Testing and tools
- Axe DevTools browser extension (quick automated scan for violations).
- WAVE or Lighthouse (built into Chrome DevTools, Accessibility tab).
- Screen reader testing: NVDA (Windows, free), JAWS (Windows, paid), VoiceOver (Mac, free). Test with actual screen reader, not just automated tools.
- Keyboard-only testing: navigate site with keyboard only, no mouse. Identify stuck focus, missing skip links.

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
