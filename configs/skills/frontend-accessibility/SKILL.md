---
name: frontend-accessibility
description: Guides comprehensive web accessibility (WCAG 2.2 AA) implementation, testing, and auditing for frontend applications. Use when building accessible UI, conducting accessibility audit, or improving accessibility of existing site.
---

# Frontend Accessibility

## Scope

This skill provides comprehensive guidance for building accessible web applications that work for all users, including those with disabilities. It complements `frontend-dev` (which now includes accessibility in rules) and `frontend-reviewer` (which checks accessibility compliance). This skill covers WCAG 2.2 AA standards in depth and testing strategies.

Use this skill when:
- Building new UI component or page
- Accessibility audit of existing site
- Fixing accessibility issues
- Choosing accessible UI library or pattern
- Planning accessibility testing strategy
- Evaluating tool/library for accessibility support

## WCAG 2.2 AA Principles

**Four principles (POUR):**
1. **Perceivable:** Information presented in ways users can perceive (not invisible)
2. **Operable:** Users can navigate and use functionality (keyboard accessible)
3. **Understandable:** Content and interfaces are clear and predictable
4. **Robust:** Works with current and future assistive technology

Each principle has guidelines and success criteria at three levels: A, AA (recommended), AAA (enhanced).

## Methodology

### 1. Semantic HTML foundation

**Use semantic elements, not generic divs:**
- Heading hierarchy: `<h1>` page title, `<h2>`/`<h3>` sections, never skip levels
- Landmarks: `<nav>`, `<main>`, `<aside>`, `<section>`, `<article>`, `<header>`, `<footer>`
- Form: `<label>` associated with `<input>`, `<fieldset>` for groups
- Lists: `<ul>`/`<ol>` for lists, not styled divs
- Links: `<a>` with href (navigates) vs `<button>` (action)
- Interactive: `<button>` for actions, not `<div onclick>`

**Benefits:**
- Assistive technology understands structure
- Keyboard navigation works automatically
- Styling/behavior predictable to users
- SEO improvement

**Screen reader experience:**
- User hears "navigation", then links (because of `<nav>`)
- User hears heading level (h1, h2, etc.) and can skip between
- Form labels read before input, not separately

### 2. Keyboard navigation

**Requirements:**
- Every interactive element reachable by Tab key
- Tab order follows visual flow (left-to-right, top-to-bottom)
- Visible focus indicator on all elements
- Can interact with Enter or Space (never mouse-only)
- Escape closes modals/menus
- Arrow keys in custom components (combobox, tree, tabs)

**Implementation:**
- Native elements keyboard accessible by default (button, link, input)
- Custom components need `tabindex`, keyboard handlers
- `tabindex="0"`: natural tab order
- `tabindex="-1"`: focusable by script, not Tab
- Never `tabindex="+1"` or higher (breaks natural order)

**Focus management:**
- Visible focus indicator (don't remove outline, if styling keep ≥3:1 contrast)
- Focus order matches visual order
- Move focus when appropriate (open modal → focus first input, close → return focus)
- Don't trap focus (modal must allow Escape to close)
- Skip links: help keyboard users skip repetitive content

**Testing keyboard navigation:**
- Close browser DevTools (focus visible)
- Navigate entire site with Tab/Shift+Tab
- Interact with Enter/Space
- Test special keys (Escape, arrows in combobox)
- No traps (can always move focus out)

### 3. Color and contrast

**Text contrast:**
- Normal text: 4.5:1 (black #000 on white #fff = 21:1, OK)
- Large text (18pt+ or 14pt bold): 3:1
- UI components (active/inactive states): 3:1
- Graphical elements: 3:1 if important to understand

**Tool:** WebAIM Contrast Checker, or browser DevTools (Chrome has contrast check)

**Principle:** Don't use color alone to convey meaning
- Bad: "Click red button" (what if colorblind?)
- Good: "Click Delete button" (label provides meaning)
- Visual: show disabled state with label + icon, not just grayed out

**Color combinations to avoid:**
- Red/green: hard for red-blind users
- Blue/yellow: hard for some color blindness
- Light gray on white: hard for all users
- Test with: Color Oracle, Coblis simulator

### 4. Images and alternative text

**Alt text:**
- Describe purpose/content, not just "image"
- Bad: "Image of a user"
- Good: "Sarah Johnson, Product Manager"
- Decorative images: `alt=""` (skip in screen reader)
- Complex images (chart, diagram): alt + link to data table or longer description

**Text in images:**
- Avoid when possible (not readable by screen readers, fails at zoom)
- If necessary: include text in alt or adjacent content
- Example: chart image with alt + data table below

**Image maps, SVG, Canvas:**
- Image maps: provide text alternative for each hotspot
- SVG: use `<title>` for accessible name
- Canvas: provide fallback text or alternative content

### 5. Form accessibility

**Label association:**
- Every input has associated `<label>` (screen reader reads it)
- Associate via `<label for="input-id">` and `<input id="input-id">`
- Or nest input in label: `<label>Email: <input type="email"></label>`
- Placeholder is not a label (disappears, not always read by screen readers)

**Form instructions:**
- Error messages clearly indicate which field has error
- Provide error message in `aria-describedby`, not just color
- Instructions for required fields, format (e.g., "Date format: MM/DD/YYYY")

**Fieldsets and legends:**
- Group related inputs with `<fieldset>` and `<legend>`
- Example: radio buttons for gender, checkboxes for interests
- Legend read by screen reader before each field

**Input type and autocomplete:**
- Use semantic input types: `type="email"`, `type="tel"`, `type="date"`
- Enables validation, mobile keyboards, browser autofill
- `autocomplete="email"`, `autocomplete="name"` helps password managers

### 6. ARIA (Accessible Rich Internet Applications)

**Use ARIA only when semantic HTML insufficient:**
- ARIA supplements HTML, doesn't replace
- Every ARIA attribute adds implementation burden

**Common ARIA:**
- `aria-label="Close"` on button with icon only
- `aria-describedby="error-message"` links field to error
- `aria-expanded="true"` on collapsible (expanded state)
- `aria-hidden="true"` hides decorative elements from screen readers
- `role="button"` on div acting as button (prefer `<button>` element)
- `role="tab"`, `aria-selected` for custom tabs
- `aria-live="polite"` for dynamic content updates

**Live regions:**
- `aria-live="polite"`: announce when content changes (search results, notifications)
- `aria-live="assertive"`: interrupt for urgent (errors, alerts)
- Useful for dynamic content that screen reader doesn't detect

### 7. Motion and animation

**Autoplay:**
- No autoplaying video/audio (disorienting, unexpected)
- User must play (press button)

**Animation duration:**
- If animation > 5 seconds, provide pause/stop control
- Respect `prefers-reduced-motion` (user preference for reduced animation)

**Prefers-reduced-motion:**
```css
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

**Flashing content:**
- Avoid flashing > 3 times per second (risk of seizures)
- Use tools: WAVE, Lighthouse to detect

### 8. Testing and tools

**Automated scanning (catches 25-30% of issues):**
- Axe DevTools (browser extension, in-depth reports)
- WAVE (Firefox/Chrome, visual annotations)
- Lighthouse (Chrome DevTools, accessibility audit)
- Run regularly during development

**Screen reader testing (catches most issues):**
- Windows: NVDA (free, open-source)
- Windows: JAWS (paid, industry standard)
- Mac: VoiceOver (built-in, free)
- iOS: VoiceOver (built-in)
- Android: TalkBack (built-in)
- Test by navigating site, hearing structure, understanding content

**Keyboard-only testing:**
- Close DevTools (focus visible)
- Navigate with Tab/Shift+Tab
- Interact with Enter/Space
- Test special keys (Escape, arrows)
- Identify traps, missing skip links

**Color contrast testing:**
- WebAIM Contrast Checker
- Chrome DevTools Inspector (Accessibility panel)
- Manual visual check with real colors

**Zoom and text sizing:**
- Browser zoom to 200%: content still readable, no horizontal scroll
- Text resize in browser settings: content adjusts
- Test on mobile (small screen)

**Manual testing checklist:**
- Page with screen reader on (does navigation make sense?)
- Site with keyboard only (can navigate and interact?)
- Site with high contrast mode (readable?)
- Site with motion reduced (animations still work?)
- Colors inverted or colorblind simulator (meaning still clear?)

### 9. Component accessibility patterns

**Combobox (autocomplete dropdown):**
- Keyboard: arrow keys navigate options, Enter selects
- `role="combobox"` on the input, `aria-expanded`, `aria-controls` pointing to the `role="listbox"` popup, and `aria-activedescendant` set to the highlighted option while focus stays on the input (ARIA 1.2 pattern; `aria-owns` is the outdated 1.1 pattern)
- Announce suggestions as user types (`aria-live="polite"`)

**Modal dialog:**
- Focus moves to first interactive element when opened
- Keyboard: Tab loops within modal, Escape closes
- Content outside modal inert: prefer a native `<dialog>` opened with `showModal()`, which makes the rest of the page inert and handles Escape; for a custom modal, set the `inert` attribute on the background content rather than `aria-hidden="true"` (which hides it from screen readers but leaves it focusable)
- Focus returns to trigger when closed

**Tabs:**
- `role="tab"`, `aria-selected`, `aria-controls`
- Keyboard: arrow keys switch tabs, Enter/Space activates
- Tab content shows/hides based on selection

**Tooltips:**
- Show on hover or focus (not just hover, keyboard users need it)
- Content matches `title` attribute or `aria-label`
- Not essential information (should duplicate elsewhere)

**Notification/Toast:**
- `aria-live="assertive"` for important notifications
- Include dismiss button
- Auto-dismiss after reasonable time (not too fast)

### 10. Accessibility in design and development process

**Design phase:**
- Include accessibility in design requirements
- Test designs with accessibility checklist
- Plan for keyboard navigation and screen readers early

**Development phase:**
- Use semantic HTML by default
- Write tests that include accessibility (getByRole in React)
- Use accessible component library or build components with accessibility
- Test with screen reader during development

**QA/testing phase:**
- Automated accessibility scanning (Axe, WAVE, Lighthouse)
- Screen reader testing (manual)
- Keyboard navigation testing
- Color contrast validation
- Zoom and text resize testing

**Accessibility audit (when hiring external help):**
- WCAG 2.2 AA audit against published standards
- Detailed report with findings, severity, remediation steps
- Ongoing monitoring and improvement plan

---

## Accessibility Maturity Levels

**Level 1: Basics (WCAG A)**
- Semantic HTML, basic keyboard navigation
- Alternative text for images
- Color not only conveyance

**Level 2: Standard (WCAG AA) ← Recommended target**
- Semantic HTML + ARIA where needed
- Full keyboard navigation with visible focus
- 4.5:1 text contrast
- Labels on all form fields
- Screen reader tested

**Level 3: Enhanced (WCAG AAA)**
- All Level 2 +
- 7:1 text contrast
- Sign language for video
- Full audio description
- Plain language

**Industry standard:** WCAG 2.2 AA (Level 2)

---

## Checklist

Work through `${CLAUDE_SKILL_DIR}/checklist.md` (`checklist.md` in this skill's directory) before signing off a page or component as accessible. Mark any item that doesn't apply as N/A and say why; don't skip it silently.

## Report

List findings in a table, most severe first, with the columns Severity | Finding | Evidence | Fix. Evidence is a file and line, plus the tool output or manual test that showed the problem. Use this scale:

- **Critical**: causes a bug, outage, data loss, or security problem now. Blocks shipping.
- **High**: a real problem with a plausible path to harm. Fix before release.
- **Medium**: adds risk or cost, but needs another condition before it causes harm.
- **Low**: an improvement or hardening opportunity, not a current problem.

End with what you didn't check and why. If an area is clean, say so in one line. An agent that defines its own report format uses that format instead.

