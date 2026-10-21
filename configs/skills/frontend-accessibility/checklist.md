# Frontend Accessibility Checklist

## Semantic HTML
- [ ] Heading hierarchy correct (h1-h6, no skips)
- [ ] Landmarks used (`<nav>`, `<main>`, `<aside>`, `<section>`, etc.)
- [ ] Form elements proper (`<label>`, `<input>`, `<fieldset>`, `<legend>`)
- [ ] Links are `<a>` (navigation) or `<button>` (action), not divs
- [ ] Interactive elements native (`<button>`, `<input>`) or proper ARIA

## Keyboard Navigation
- [ ] All interactive elements reachable by Tab
- [ ] Tab order follows visual flow (left-to-right, top-to-bottom)
- [ ] Focus indicator visible on all focusable elements
- [ ] Focus indicator contrast ≥ 3:1
- [ ] Can interact with Enter/Space
- [ ] No keyboard traps (can always move focus out)
- [ ] Escape closes modals/menus
- [ ] Arrow keys work in custom components (tabs, combobox, tree)
- [ ] Skip links present (skip navigation)
- [ ] Focus moves to opened modal, returns on close

## Color and Contrast
- [ ] Normal text contrast ≥ 4.5:1
- [ ] Large text (18pt+) contrast ≥ 3:1
- [ ] UI components (active/inactive) contrast ≥ 3:1
- [ ] Graphical elements (if meaningful) contrast ≥ 3:1
- [ ] Meaning not conveyed by color alone (icon + label, etc.)
- [ ] Color combinations avoid red/green, blue/yellow
- [ ] Tested with contrast checker tool

## Images and Alternative Text
- [ ] All images have alt text or `alt=""`
- [ ] Alt text describes purpose/content (not just "image")
- [ ] Decorative images have `alt=""`
- [ ] Complex images have longer description (link or adjacent text)
- [ ] Text in images included in alt or adjacent content

## Forms
- [ ] All inputs have associated `<label>`
- [ ] Label text clearly describes input
- [ ] Required fields marked
- [ ] Format instructions provided (if relevant)
- [ ] Error messages clear and associated with field
- [ ] Related inputs grouped with `<fieldset>` and `<legend>`
- [ ] Input types semantic (`type="email"`, `type="tel"`, etc.)
- [ ] Autocomplete attributes set for password managers

## ARIA (if used)
- [ ] ARIA used only when semantic HTML insufficient
- [ ] `role` correct and justified
- [ ] `aria-label` or `aria-labelledby` for unlabeled elements
- [ ] `aria-describedby` for descriptions
- [ ] `aria-expanded`, `aria-selected`, `aria-hidden` accurate
- [ ] `aria-live` for dynamic content
- [ ] Live regions tested (announcement verified)

## Motion and Animation
- [ ] No autoplaying video/audio
- [ ] Animation > 5 seconds has pause/stop control
- [ ] `prefers-reduced-motion` respected (animations reduced/disabled)
- [ ] No flashing content (>3x per second)
- [ ] Transitions smooth and purpose clear

## Mobile and Touch
- [ ] Touch targets ≥ 24x24 CSS pixels or adequately spaced (WCAG 2.2 AA, 2.5.8); 44x44 for AAA (2.5.5)
- [ ] Touch targets spaced apart (not cramped)
- [ ] Zoom supported (no `user-scalable=no`)
- [ ] Responsive at 200% zoom
- [ ] Text readable without horizontal scroll at 200% zoom
- [ ] Mobile touch-friendly (not hover-only)

## Testing
- [ ] Automated scan run (Axe, WAVE, Lighthouse)
- [ ] No critical/serious issues in scan results
- [ ] Screen reader tested (NVDA, JAWS, VoiceOver)
- [ ] Keyboard-only navigation tested
- [ ] Color contrast verified (manually or tool)
- [ ] Zoom to 200% tested (readable, no horizontal scroll)
- [ ] Mobile/touch tested
- [ ] Focus and tab order verified

## Component Patterns
- [ ] Custom components (tabs, combobox, modal) implement expected keyboard behavior
- [ ] Dialog/modal focus management working
- [ ] Tooltips keyboard accessible
- [ ] Notifications announce to screen readers
- [ ] Form validation errors clear and accessible

## Documentation
- [ ] Accessibility approach documented
- [ ] WCAG level targeted (A, AA, AAA)
- [ ] Known accessibility limitations documented
- [ ] Future accessibility improvements planned
- [ ] Component accessibility patterns documented
- [ ] Testing procedures documented

