# internal/lint

## Purpose

Repo-wide convention tests that don't belong to any single production package. Each test
walks the whole repository tree and asserts one structural rule. The package has almost no
production code — it exists to be run by `go test ./...`.

## How it works

- `RepoRoot(t)` (reporoot.go) walks up from the test's working directory to the directory
  holding `go.mod`, so a test can anchor a `filepath.WalkDir` over the entire tree
  regardless of which package's test binary is running.
- Each `*_lint_test.go` file holds one `Test…` function that scans files and fails with a
  list of offenders.

## The conventions enforced

- `turntests_lint_test.go` — new `turntests/` test files must not drive the attack turn via
  `ge.ResolveAttackStep(...)` directly; they must use the public Eval entry points
  (`sim.EvalOneTurnForTesting` / `sim.EvalTwoTurnsForTesting`). A
  `grandfatheredResolveAttackStepFiles` allowlist freezes the v2-migration backlog of files
  that still call `ResolveAttackStep`; migrating one means rewriting it against the public API
  (or moving it to a same-package unit test) and removing its allowlist entry, after which
  the lint rejects reintroduction.
- `card_markers_lint_test.go` — `NotImplemented` / `Unplayable` markers on cards appear
  only in the dedicated `notimplemented/` and `unplayable/` subpackages.
- `test_package_lint_test.go` — test files use the production package name, not an
  external `package foo_test`, so tests can reach unexported helpers.
- `builder_lint_test.go` — a fluent `GameStateBuilder()` chain with two or more `Set*`
  calls is broken one call per line.
- `cardgen_staleness_test.go` — every cardgen-generated `_gen.go` file matches what
  `cardgen` produces from the yamls today.
- `registry_lint_test.go` / `weapon_registry_lint_test.go` — every implemented card
  variant and every implemented weapon is present in the registry.
- `resource_source_lint_test.go` — a card calling `GameEngine.AddResourcePoints` also
  declares `MaxResourcePoints`, the `card.ResourceSource` contract the attack-budget prune
  depends on.

## How to extend

To add a convention, drop a `<rule>_lint_test.go` file with a single `Test…` function that
uses `RepoRoot(t)` plus `filepath.WalkDir`, collects offenders, and `t.Errorf`s each one
with a message pointing at the rule.

## Gotchas

- These tests read the repo from disk, so they fail outside a checkout (the `go.mod` walk
  panics). They are meant to run via `go test ./...` inside the repo.
