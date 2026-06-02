# cmd/parsecarddb

## Purpose

A helper command-line tool for exploring the upstream Flesh and Blood card database. It
parses the `card.csv` from `github.com/the-fab-cube/flesh-and-blood-cards` and prints
matching cards in either a pretty human-readable form or JSON. It is the preferred way to
look up card data and printed text when implementing a card — use it instead of grepping
the raw CSV.

## How it works

`card.csv` is a tab-separated export, read by `internal/textio`'s `LoadCardCSV` (a header-name-
to-field mapping, so an upstream column reorder doesn't break parsing). `poolFilter.matches`
then narrows the rows. The class / talent filters express a hero's legal pool the way `registry.heroCanPlay`
does: Generic and classless cards are always class-legal, and a card's talents must be a subset
of the hero's. Class / talent / non-deck words are classified via `internal/fabtype`.

## How to use

```
go run ./cmd/parsecarddb [flags]
```

Flags:

- `-in` — path to `card.csv` (default `data_sources/card.csv`).
- `-name` — only print cards whose name contains this substring (case-insensitive).
- `-type` — only print cards whose Types field contains this substring (case-insensitive).
- `-classes` — comma-separated hero classes; a card must be Generic, classless, or one of
  these. Empty (default) = any class.
- `-talents` — comma-separated allowed talents; a card's talents must be a subset, so
  `Lightning` admits no-talent and Lightning cards. Empty (default) = no-talent cards only.
- `-any-talent` — ignore `-talents` and admit cards carrying any talent.
- `-require-talents` — comma-separated talents; admit only cards carrying at least one (e.g.
  to list just the Lightning cards in a pool).
- `-modeled` — keep only card types the optimizer models — deck cards **and weapons** —
  dropping equipment, heroes, and tokens.
- `-silver-age` — require Silver Age legality via the card.csv `Silver Age Legal` column
  (default true). This column lags new sets — pair it with `-rarity-legal` for an accurate pool.
- `-exclude-banned` — drop cards on our `internal/format` Silver Age banlist.
- `-rarity-legal` — keep only cards with a Basic/Common/Rare printing — the Silver Age rarity
  rule, which is the reliable signal (the `Silver Age Legal` column lags new sets). Reads
  `-printings` (`card-printing.csv`) and `-rarity-file` (`rarity.csv`), joining by card Unique
  ID. Run `fetch.sh` first to have those files locally.
- `-format` — `pretty` (default) or `json`.
- `-names_only` — print only the distinct card names, one per line.

### Examples

```
# Viserai's pool (Runeblade + Generic, no talent, incl. weapons, B/C/R only):
go run ./cmd/parsecarddb -classes Runeblade -modeled -rarity-legal -names_only

# Aurora's Lightning cards worth implementing (incl. weapons, drops banned + non-B/C/R):
go run ./cmd/parsecarddb -classes Runeblade -talents Lightning -require-talents Lightning \
  -modeled -exclude-banned -rarity-legal -names_only
```

## Important files

- `main.go` — the entire tool: the `Card` struct mirroring the CSV columns, the loader, the
  `poolFilter`, and the output formatting.

## Gotchas

- `-talents` is subset semantics, so the default (empty) admits **only** no-talent cards. Pass
  the hero's talents (or `-any-talent`) to widen it — a Lightning card won't appear under the
  default filter even with a matching `-name`.
