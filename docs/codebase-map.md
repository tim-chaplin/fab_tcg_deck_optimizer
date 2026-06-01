# Codebase map

A high-level map of `fab_tcg_deck_optimizer` — a deck-finding tool for the Flesh and Blood
TCG. Start here when bootstrapping context for a new task.

The tool loads a 40-card deck, simulates many hands of solitaire play, scores each hand
(damage dealt + damage prevented), and runs a simulated-annealing search over single-card
swaps to find high-value decks. The root `README.md` has the user-facing overview and the
full `fabsim` CLI flag reference.

## How to read this repo

1. This file — the directory map and a one-line description of every package.
2. Each package directory has its own `README.md` with purpose, key types, how to extend it,
   and gotchas. Open the README for whatever package a task touches.
3. `docs/dev-standards.md` — code-quality and convention rules; the audience is the
   `pr-standards-reviewer` agent. It is a quality checklist, not a how-to guide.

How-to-implement material (adding a card, wiring a rider, the aura/weapon/item lifecycles)
lives in the package READMEs — `internal/card/cards/README.md` is the main "how to add a
card" reference.

## A run, end to end

`cmd/fabsim` parses a subcommand, loads a deck via `internal/textio`, and hands it to
`internal/sim`. For each shuffled hand, `sim` brute-forces the optimal turn with a two-layer
search: enumerate every Pitch/Attack/Defend role partition, then permute the attack-turn search,
replaying each ordering through the `internal/gameengine` state engine. Per-hand value folds
into deck `Stats`. `internal/deck` enumerates single-slot mutations; `anneal` accepts or
rejects each via the Metropolis rule and keeps the best deck found.

## Directory tree

```
cmd/
  cardgen/         Generator CLI: regenerates per-card _gen.go from yaml
  fabsim/          Main CLI: anneal / eval / compare / import
  parsecarddb/     Helper: parse and filter the upstream card database
  crosscheck/      Diagnostic: diff our legal pool against the official cardvault API
  cardaudit/       Diagnostic: cross-check implemented card stats vs card.csv + the official API
internal/
  card/            Card interface, CardState, TypeSet, engine-facing contracts
    cards/         Every implemented card (yaml + _gen.go + .go triples)
  cardgen/         Code generator: yaml -> _gen.go + registry map
  deck/            Deck type, mutation enumeration, per-deck Stats
  gameengine/      Per-turn game-state engine the attack-turn runner drives
  aura/            Aura: a persistent hook firing a handler at a scheduled event
  trigger/         Shared trigger machinery (generic core + one-shot triggers)
  triggertype/     Bitmask enum for when an aura/trigger fires
  item/            In-play permanent items (token abilities, card-sourced triggers)
  token/           Factories for the five built-in tokens
  hero/            Hero interface; concrete heroes in heroes/
  weapon/          Weapon interface; concrete weapons in weapons/
  ids/             Stable integer ID allocation for cards / heroes / weapons
  fabtype/         FaB class / talent / non-deck type NAME lists for the CSV + typebox tools
  format/          Constructed gameplay formats (Silver Age banlist + legality predicate)
  registry/        Master roster of every card / weapon / hero
  sim/             Hand-and-deck evaluator and attack-turn search
  textio/          On-disk deck encodings (JSON deck+stats, fabrary text)
  lint/            Repo-wide convention tests
  testutils/       Card / hero / weapon fakes shared across tests
turntests/         Turn-level tests driving full turns through public entry points
data_sources/      Upstream card.csv, comprehensive rules, banlist
docs/              This map + dev-standards.md
scripts/           PowerShell anneal wrappers + Python card-stub generators
mydecks/           Local working deck files (untracked by git)
```

## Packages

### Card model and data

- **internal/card** — The contract layer: the `Card` interface every card implements, the
  per-attack-step `CardState` wrapper, the `TypeSet` type-line bitfield, and the narrow
  `GameEngine` / `Logger` / `Aura` interfaces cards consume from the simulator without
  importing it.
- **internal/card/cards** — Every implemented card as a three-file group (`.yaml` data,
  generated `_gen.go` statics, hand-written `.go` carrying `Play` and riders); pool-excluded
  cards are quarantined in the `notimplemented/` and `unplayable/` subpackages.
- **internal/cardgen** — The code generator that turns `<card>.yaml` data files into per-card
  `_gen.go` struct/method declarations and the registry's `cardsByID` map.
- **cmd/cardgen** — The thin CLI wrapper that runs `internal/cardgen` and writes the
  generated files, normally invoked via `go generate`.
- **internal/ids** — The central allocation of stable integer identifiers (`CardID`,
  `HeroID`) anchored into contiguous non-overlapping ranges so per-entity caches can be
  slice-indexed. Weapons are cards, so they take `CardID`s in the shared range.
- **internal/fabtype** — The text-side FaB vocabulary: the sets of class, talent, and non-deck
  card-type NAMES as they appear in printed type lines / the upstream `card.csv` Types column.
  Used by the CSV / typebox tools (`cmd/parsecarddb`, `cmd/crosscheck`) to classify a type word;
  the in-engine legality check uses `registry`'s `classMask` / `talentMask` instead.
- **internal/format** — The constructed gameplay formats (Silver Age today): a `Format`
  interface whose `IsCardLegal` predicate filters the card pool, backed by the hardcoded
  banlist in `banlist.go`. Hero-agnostic — class/talent legality is the registry's job.
- **internal/registry** — The master roster of every implemented card, weapon, and hero,
  providing ID/name lookups and the deck-construction pools filtered by marker, format
  banlist, and the hero's class / talents (`LegalCardsFor` / `LegalWeaponsFor`).

### Simulation engine

- **internal/sim** — The hand-and-deck evaluator: brute-forces each hand's optimal turn line
  via a two-layer search (role-partition enumeration then attack-turn permutation) and
  folds per-hand value into deck stats for the annealing optimizer.
- **internal/gameengine** — The per-turn game-state engine the attack-turn runner drives, splitting
  raw turn data (`GameState`) from the rules-engine API (`GameEngine`: trigger dispatch,
  attack-step resolution, token economy, deck manipulation).
- **internal/aura** — The concrete `Aura` type: a persistent hook that fires a typed handler
  at a scheduled lifecycle event, used to model any card that creates something firing later.
- **internal/trigger** — The shared trigger machinery: the embeddable generic `Trigger`
  core that auras and items carry, plus the one-shot `EphemeralTrigger` the engine fires once
  and drops.
- **internal/triggertype** — A dependency-free micro-package holding the bitmask enum that
  categorizes when an aura or trigger fires (start of turn, card played, hit, end of turn).
- **internal/item** — The concrete `Item` type for in-play permanents: token items carrying
  an activated ability the attack-turn runner enqueues, or card-sourced items carrying a trigger
  that fires like an aura.
- **internal/token** — Factories for Flesh and Blood's five built-in tokens (Gold / Silver /
  Copper items, Runechant / Ponder auras), pre-wired with identity and fire behavior.

### Deck, heroes, weapons

- **internal/deck** — A candidate deck (hero, format, weapons, card list, sideboard) plus the
  single-slot and synergy-pair mutation enumeration the anneal search drives over, and the
  per-deck `Stats` result types. Hero and format are durable attributes that scope the legal
  pool the mutation enumeration draws from.
- **internal/hero** — The `Hero` interface a hero card satisfies and the narrow engine
  surfaces its abilities consume; concrete heroes (currently Viserai) live in `heroes/`.
- **internal/weapon** — The platonic weapon `Card` plus the mutable engine-side `Weapon`
  object built when it's equipped (mirroring the aura / item model), where each weapon pairs a
  permanent type with a `card.Card` activated ability; concrete weapons (including
  `NotImplemented`-marked ones) live in `weapons/`.

### Persistence and performance

- **internal/textio** — Durable on-disk text encodings: the canonical `mydecks/*.json`
  deck-plus-stats format and the fabrary.net plain-text import/export format, plus `mydecks/`
  path resolution and a read-only loader for the upstream `card.csv` (`LoadCardCSV`, shared by
  the card-data tools).

### CLI and tooling

- **cmd/fabsim** — The user-facing CLI that searches (`anneal`), evaluates, compares, and
  imports decks, with one `mode_*.go` file per subcommand.
- **cmd/parsecarddb** — A helper tool that parses and filters the upstream the-fab-cube
  `card.csv`; the preferred way to look up card data and printed text. Its class / talent flags
  enumerate a hero's legal pool (e.g. Aurora's Lightning cards) using `internal/fabtype`.
- **cmd/crosscheck** — A diagnostic that pulls the official cardvault.fabtcg.com card list and
  diffs it against `registry.LegalCardsFor(SilverAge, Viserai)`, flagging cards the official
  site lists that we're missing (and any in our pool it no longer lists). Needs network.
- **cmd/cardaudit** — A report-only diagnostic that cross-checks every implemented card's yaml
  stats (cost / power / defense) against both `card.csv` and the official cardvault API,
  bucketing discrepancies by confidence and suppressing an allowlist of intentional modeling
  deviations. Needs network.

### Testing

- **turntests** — Top-level turn-level tests that drive a single turn through the public
  `Eval` entry points and assert per-turn `Value`; the project's primary card-behaviour
  verification.
- **internal/lint** — Repo-wide convention tests that each walk the whole tree and assert one
  structural rule (marker placement, generated-file staleness, registry coverage, turntests
  entry-point discipline).
- **internal/testutils** — Configurable card, hero, and weapon fakes shared across package
  tests so predicate, partition, and attack-turn runner assertions have controllable inputs.
