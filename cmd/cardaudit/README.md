# cmd/cardaudit

A **report-only** diagnostic that cross-checks the numeric stats of every implemented card's
yaml against both the upstream `card.csv` and the official cardvault.fabtcg.com API, and prints
the discrepancies. It changes no files — use it to catch a yaml whose cost / power / defense was
transcribed wrong, or where `card.csv` itself disagrees with the printed card.

## What it compares

For each yaml card variant, matched to the two sources by **(card name, pitch)**:

- **cost** (skipped for variable-cost cards), **attack** (printed power), **defense**.

A field is only flagged when the source value is a clean integer that differs — a blank /
non-numeric source value (an aura has no printed power) is skipped. Type lines and rules text
aren't compared (their wording differs across sources); this is a stats pass.

## Output

- **[A] yaml disagrees with BOTH sources** — the high-confidence bucket; almost certainly a
  yaml transcription error.
- **[B] yaml disagrees with one source only** — usually the yaml matches `card.csv` (what it was
  generated from) but the live API differs, which points at a rebalance/errata or a card.csv
  that predates it. Needs a judgement call about which printing we model. (The API value can
  also depend on which printing matched first for cards with multiple printings.)
- **[C] expected deviations** — discrepancies that are intentional: a card whose effect we
  encode into a stat (e.g. Drag Down's attack-debuff modeled as block). These live in the
  `expected` allowlist in `main.go`; add an entry with a reason when you make such a choice, so
  it stops showing up as an error.
- **[D] name+pitch not found in a source** — a name mismatch, a double-faced card, or a card
  not in that source.

## Run

```
go run ./cmd/cardaudit            # needs outbound HTTPS to the cardvault API
go run ./cmd/cardaudit -csv data_sources/card.csv -cards-dir internal/card/cards
```

Re-run after refreshing `card.csv` (new set) or implementing a batch of cards.
