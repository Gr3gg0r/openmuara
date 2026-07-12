> **Initiative:** OpenMuara Accessibility & Usability Polish

# Risk Register

| # | Risk | Likelihood | Impact | Mitigation |
|---|------|------------|--------|------------|
| R1 | Changing table rows to buttons/links accidentally breaks the click-to-detail behavior. | Medium | Medium | Add tests covering both mouse click and keyboard activation; keep the existing `onClick` path where possible. |
| R2 | Focus-trap implementation interferes with Vite HMR or test environment. | Low | Low | Write a minimal hook that only activates in real browsers; tests can mock it. |
| R3 | Adding `<main>` to provider pages shifts layout or breaks Go template rendering. | Low | Low | Verify each page renders with sample data and run the smoke test. |
| R4 | Bundle size grows from new utility hooks. | Low | Low | Measure with `npm run bundle-size`; keep utilities under ~200 lines total. |
| R5 | Keyboard shortcuts change conflicts with browser/OS shortcuts. | Low | Low | Ignore events when Ctrl/Alt/Meta are pressed; document shortcuts in the Help modal. |
