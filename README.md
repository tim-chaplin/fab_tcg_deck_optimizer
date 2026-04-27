# fab-deck-optimizer

A deck-finding tool for the Flesh and Blood TCG, written in Go.

Built with Claude Code Opus.

## Goal

Find optimal deck lists under a given set of assumptions. I model a deck's value as the average
value produced by each hand over many plays through the deck. A hand's value is the sum of damage
dealt, and damage prevented.

## FAQ

### So are these AI-generated decks?

No. This is just a computer program that implements an evaluation function (the simulator) and
uses known optimization techniques to look for decks that optimize the evaluation. You can
compile and run the program on your own computer without having invoked AI at all. AI was used
to write the program (much) faster than I could have by hand; but the program could have been
written without AI.

### Is this the best possible deck in the format?

No.

1) The search space is too big to exhaust; for any given set of modelling assumptions I don't
   know that it'll ever be practical to find the true global optimum. This program finds local
   maxima, and I don't have a way to prove how far any given local max is from the global one.

2) "Best" depends on your assumptions, and especially on the specific matchup; there's no such
   thing as a single deck that's optimal across every matchup and every assumption (assuming
   even a moderately well-balanced format). The tool outputs a deck that's strong for
   goldfishing under the assumptions listed in "Scope & limitations".

## Observations

- Nebula Blade is probably the strongest weapon for Viserai in the format. Even at very early
  stages of development (when the card pool and card-effect modelling were much rougher than
  they are today) the optimizer converged on Nebula Blade over every other weapon loadout
  almost immediately.
- Mauvrion Skies wants all six copies. Across seeds and starting decks the optimizer tends to
  fill all six legal slots (two each of red / yellow / blue).
- Even with card draw being undervalued in the simulation (because we never even consider lines
  where the card is played the same turn it's drawn), the optimizer keeps converging on Drawn to
  the Dark Dimension and Snatch.

## Scope & limitations

This is a work in progress. The current model is deliberately narrow:

- **Hero pool.** Only cards legal for Viserai in Silver Age are modeled.
- **No opponent counterplay.** The opponent is represented by a single configurable `-incoming`
  value: a static amount of damage per turn that the hand can defend against. There are no
  blocks from hand, no disruption, no reaction windows. This is the goldfishing assumption.
- **Card coverage is incomplete.** Most Runeblade and Generic Silver Age cards are implemented, but
    many are still stubs, or simplified.

## How it works

- Loads a 40-card deck from `mydecks/<name>.json`.
- Shuffles and repeatedly draws hands of 4 cards.
- For each hand, brute-forces the optimal play: every partition of the hand into Pitch / Attack
  / Defend, every weapon-swing subset, every legal attack ordering (respecting Go again). Hand
  value = damage dealt + damage prevented (capped at `-incoming`).
- After shuffling and drawing through the deck a certain number of times, assigns it a score.
- Generates new decks by randomly swapping out cards, and saving the deck with the higher score.

## Usage

`fabsim` takes a subcommand as its first argument. Running `fabsim` with no subcommand prints
the catalogue.

Deck names are resolved to `mydecks/<name>.json`; the `.json` suffix is optional. Subcommands
that always operate on a specific deck (`eval`, `compare`) take the deck name(s) as positional
arguments. `anneal` uses a `-deck` flag instead because the name can be omitted (the default
is `<hero>_<format>_<incoming>_incoming`, e.g. `viserai_silver_age_0_incoming`, so different
(hero, format, `-incoming`) regimes keep separate checkpoints) and the named file doubles as
a resume point when it already exists.

Each subcommand parses its own flag set, so `fabsim <subcommand> -help` lists exactly the flags
that apply.

- **`anneal`** — simulated-annealing search on the deck at `-deck`, or on a fresh random deck
  if the file doesn't exist yet. Each round enumerates every single-slot mutation (every
  alternative weapon loadout + every (card-in-deck, card-out-of-deck) swap) and evaluates each
  one at the current `-shuffles` budget. The acceptance gate is the Metropolis rule: strict
  improvements are always accepted, worse mutations are accepted with probability
  `exp((avg - baseline) / T)` when `-start-temp > 0`. Temperature decays geometrically per
  acceptance (`-temp-decay`, floored at `-min-temp`). At `-start-temp 0` the gate is strictly
  `> baseline` and anneal degenerates to a classical hill climb. A round with zero acceptances
  is treated as a local maximum and anneal exits. Press Enter to abort mid-round (exits 130 so
  wrapper scripts can tell this apart from natural convergence). Only the best-ever deck is
  persisted to disk — walks through worse states under annealing don't regress the JSON.
- **`eval`** — `fabsim eval <deck>`. Loads the deck file, simulates it for `-shuffles` hands
  against `-incoming` damage, prints the resulting stats, and rewrites both the `.json` and
  sibling fabrary `.txt` so the saved copy stays in sync with the current binary's modelling.
  Pass `-print-only` to skip the sim and just print the last run's stats without touching the
  file.
- **`import`** — interactively imports a deck from fabrary.net. Prompts for a deck name, then
  asks you to paste the plain-text export; input ends automatically at fabrary's
  `See the full deck @ …` footer. Saves the result as `mydecks/<name>.json`. Cards the
  optimizer hasn't implemented yet are skipped with a warning rather than blocking the import.
- **`compare`** — `fabsim compare <deck1> <deck2> -incoming N`. Re-scores both decks at the
  same fixed `-shuffles` / `-incoming` so the comparison is apples-to-apples (the .json files
  are rewritten with the fresh stats — card lists are unchanged), then prints a stat-by-stat
  side-by-side report: pitch counts, mean hand value, per-cycle means, the two hand-value
  histograms, and the per-card count delta. The (shuffles, incoming) settings ride at the top
  so the per-section rows don't repeat them.

#### Adaptive vs fixed shuffles

`-shuffles` controls the per-eval shuffle budget across all subcommands:

- `-shuffles -1` (default for `anneal` and `eval`) — adaptive. Each eval keeps shuffling until
  the per-turn mean's standard error drops below an internal target (~±0.05), then stops.
  Typical Viserai decks converge in 200–400 shuffles; an internal cap stops a pathological
  high-variance regime that doesn't converge. Use this for everyday hill-climbs and one-off
  re-scores where ±0.05 precision on the mean is plenty.
- `-shuffles N` (any non-negative value) — fixed. Every eval runs exactly N shuffles, giving
  apples-to-apples comparisons across mutations. Use this for repro flows and any time you
  want every eval scored at the same budget.
- `compare` always uses fixed `-shuffles` (default 10000); adaptive isn't allowed there because
  the side-by-side comparison needs matched conditions on both decks.
- `anneal -finalize` pins `-shuffles` to 100000 and tightens `-min-improvement` to 0.01 — a
  high-precision pass for decks that have already converged.

### Suggested workflow

Choose an amount of incoming damage per turn (basically your tuning knob for how aggressive vs.
defensive the deck will be). Run in continuous annealing mode to look for the best deck for the
chosen assumption:

```
./scripts/anneal-reanneal.ps1 -Deck viserai_silver_age_7_incoming -StartTemp 1 -Incoming 7
```

Or fan out across several independent starts for potentially better coverage of the solution space,
and rank the results:

```
./scripts/anneal-restarts.ps1 -N 10 -DeckTemplate 'viserai_*' -Incoming 7 -StartTemp 1
```

Each run only overwrites its deck file when a new best-ever avg is found, so it's safe to
re-run indefinitely.

### Flags

Each subcommand owns its own flag set; `fabsim <subcommand> -help` is the authoritative list.
The summary below groups the flags by subcommand.

**`anneal`** (search + resume + re-score the baseline on load):

- `-deck` — checkpoint name; resolved to `mydecks/<name>.json` (default
  `<hero>_<format>_<incoming>_incoming`, keyed off the hero, format, and `-incoming`). The
  `mydecks/` directory is created automatically. If the file exists anneal resumes from it.
- `-shuffles` — per-eval shuffle budget. `-1` (default) runs adaptively; any non-negative value
  pins a fixed count for apples-to-apples acceptance. See "Adaptive vs fixed shuffles" above.
- `-incoming` — opponent damage per turn (default 0)
- `-deck-size` — cards per deck, used only for random starting decks (default 40)
- `-max-copies` — max copies of any single card printing (default 2)
- `-seed` — RNG seed (default: time-based)
- `-format` — constructed format whose banlist restricts the card pool during search. Defaults
  to `silver_age`, currently the only supported format. The authoritative Silver Age banlist
  lives at `data_sources/silver_age_banlist.txt`.
- `-start-temp` — starting temperature. `0` (default) runs a pure hill climb. Higher values
  probabilistically accept worse mutations early (Metropolis rule).
- `-temp-decay` — multiplicative cooling per acceptance (default 0.95).
- `-min-temp` — temperature floor (default 0).
- `-finalize` — high-precision pass — sets `-shuffles` to 100000 (fixed) and tightens
  `-min-improvement` to 0.01. Use on a deck that's already converged to squeeze out the
  remaining sub-percent improvements.
- `-reevaluate` — force re-evaluation of the loaded deck's baseline avg even if its prior run
  count already matches the current `-shuffles` budget. Use after adjusting modelling
  assumptions.
- `-quiet-load` — skip the baseline card-list dump at startup. Used by
  `scripts/anneal-reanneal.ps1` from pass 2 onward so the unchanging listing doesn't flood the
  log.
- `-debug` — force per-round logs even when annealing is on (T>0 normally hides them).

**`eval`** (re-score a deck and rewrite it; `-print-only` skips the sim and the rewrite):

- `-shuffles` — per-eval shuffle budget. `-1` (default) runs adaptively; any non-negative
  value runs exactly that many shuffles. See "Adaptive vs fixed shuffles" above.
- `-incoming` — opponent damage per turn (required unless `-print-only` is set)
- `-seed` — RNG seed (default: time-based)
- `-format` — format predicate applied to replacement picks when the loaded deck contains
  NotImplemented cards (default `silver_age`)
- `-max-copies` — max copies per printing, applied when replacing NotImplemented cards
  (default 2)
- `-print-only` — skip the sim; load the deck and print the persisted stats without
  rewriting the `.json` / `.txt`
- `-brief` — print only the score summary (no card list, per-card stats, or best turn)

**`compare`** (re-score both decks at matched settings before reporting):

- `-shuffles` — shuffles per deck used for the re-score (default 10000). compare always runs
  fixed shuffles; adaptive isn't allowed because both decks need matched conditions.
- `-incoming` — opponent damage per turn (required; both decks are re-scored against this value)
- `-seed` — RNG seed (default: time-based)
- `-format` — format predicate applied to replacement picks when a loaded deck contains
  NotImplemented cards (default `silver_age`)
- `-max-copies` — max copies per printing, applied when replacing NotImplemented cards
  (default 2)

**`import`**: no flags; see the usage line above.

Helper tool for exploring the upstream card database:

```
go run ./cmd/parsecarddb --names_only
```

## Tests

```
go test ./...
```

## Layout

```
.github/workflows/   GitHub Actions (anneal-sweep: parallel matrix run)
cmd/fabsim/          CLI entry point
cmd/parsecarddb/     Card-database parser / filter
internal/card/       Card interface, TurnState, and card implementations
internal/deck/       Deck construction and parallel anneal round driver
internal/deckio/     JSON serialisation / deserialisation
internal/hand/       Optimal-play solver for a single hand
internal/hero/       Hero definitions and on-play triggers
internal/weapon/     Weapon definitions
scripts/             PowerShell wrappers for multi-restart and reanneal sweeps
```
