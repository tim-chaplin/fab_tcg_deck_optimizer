# cmd/fabsim

## Purpose

`fabsim` is the command-line entry point for the deck optimizer: it searches for, evaluates,
compares, and imports Flesh and Blood decks. It is the user-facing tool the project is
built around.

## Subcommands

`fabsim` takes a subcommand as its first argument; running it with none prints the
catalogue. Each subcommand parses its own `flag.FlagSet`, so `fabsim <subcommand> -help`
lists exactly the flags that apply.

- `anneal` — simulated-annealing (or, at temperature 0, classical hill-climb) search on the
  deck at `-deck`, or on a fresh random deck if the file doesn't exist. Each round
  enumerates every single-slot mutation and accepts via the Metropolis rule.
- `eval` — re-score a saved deck and rewrite both its `.json` and sibling fabrary `.txt`;
  `-print-only` skips the sim and just prints the last run's stats.
- `compare` — re-score two saved decks at matched fixed `-shuffles` / `-incoming` and print
  a side-by-side stat / histogram / card-delta report.
- `import` — interactively paste a fabrary.net plain-text export and save it as
  `mydecks/<name>.json`.

The root `README.md` documents the full flag set, the adaptive-vs-fixed `-shuffles`
behaviour, the suggested workflow, and the wrapper scripts — refer to it rather than
duplicating flag docs here.

## Layout

- `main.go` — subcommand dispatch; an `init` warms the attack-step text cache from the
  registry.
- `flags.go` — shared flag-parsing helpers (`parseFlagsAnywhere`, `requireFlag`).
- `mode_anneal.go` / `mode_eval.go` / `mode_compare.go` / `mode_import.go` — one file per
  subcommand: flag set, dispatch, and run logic.
- `deckio_helpers.go` — deck load / save plumbing (`writeDeck` applies `ViseraiDefaults`
  before serialising so the persisted `.json` and `.txt` carry the full loadout).
- `gameplay_format.go` — the constructed-format flag value and parsing.
- `print.go` — human-readable deck / stats renderings.

## How to extend

To add a subcommand, add a `mode_<name>.go` with a `run<Name>Cmd(args)` function that owns
its `flag.FlagSet`, register it in `main`'s switch, and add a line to `printSubcommands`.

## Gotchas

- `anneal` exits 130 when the user aborts mid-run by pressing Enter, so wrapper scripts can
  distinguish an abort from natural convergence.
- Only the best-ever deck is persisted; annealing's walks through worse states never
  regress the on-disk JSON.
- `main` creates `mydecks/` up front so a long run can't fail a late `WriteFile` on a
  missing directory.
