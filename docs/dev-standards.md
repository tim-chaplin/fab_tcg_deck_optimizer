# Developer Standards

Code-quality and convention rules for this repo. The audience is the `pr-standards-reviewer`
agent, which audits every PR against this document before it is surfaced for review.

This file is deliberately scoped to *quality* rules. How-to-implement material — adding a
card, wiring a rider, the aura/weapon/item lifecycles, registry layout — lives in the
per-package `README.md` files; `docs/codebase-map.md` is the index into them.

## Comments

- Wrap at 100 chars. Wrap only at word boundaries; never reflow code to meet the limit.
- Describe CURRENT behavior and its motivation. No history references: drop "replaces X",
  "was Y before", "now does Z", and "previously / formerly / legacy / deprecated" framing.
  The git log is the authoritative history.
- Delete completed TODOs instead of rewording them. If that empties a section, delete the
  heading too.
- Card docstrings cover card-SPECIFIC quirks — the printed rules text, any modelling fudge,
  and anything surprising about how the card interacts with the solver. Generic framework
  plumbing belongs in package docs, not repeated per card.
- Every card file's top-of-file docstring MUST include the printed `Text: "..."` block, even
  when `Play` wires it trivially. This is required — do not flag it as trimmable
  transcription. The "card-specific quirks only" rule covers the prose around the printed
  text, not the text itself.
- Don't restate behavior already documented by an external function, type, or marker the
  card uses. A card carrying `card.Dominator` doesn't re-explain Dominate; a `NotImplemented`
  card with a `// not implemented: <quirk>` line doesn't repeat "rider isn't modelled".
- A function's docstring describes what the function does, not how callers in other packages
  use it. Cross-file implementation references ("used by the chain runner's per-permutation
  reset to drop millions of allocs") couple two sites that rot independently — that belongs
  in the PR description.
- Default to no comment; add one only when the *why* is non-obvious.

## Cross-file references

If a comment's rationale would otherwise cite "matches the pattern in foo.go, bar.go,
baz.go", factor the shared rule into the relevant package `README.md` and cite only the
local behaviour at the call site.

## Test docstrings

A test's doc comment is a single brief sentence stating the behavior under test, e.g.
`// Tests that a single pitch paying for multiple Aether Slashes activates the bonus on each.`
Inputs, expected values, and chain shape are visible in the test body. Same rule for unit and
turn-level tests.

## Test scope

Tests exercise card and engine BEHAVIOR, not printed-attribute transcription. Don't assert
`Card.GoAgain()` / `Attack()` / `Cost()` / `Pitch()` for a specific card — the solver pulls
those numbers through on every hand, so real miswiring surfaces in the behavioural tests. A
narrow "stat == X for this card" pin adds maintenance cost and catches nothing.
