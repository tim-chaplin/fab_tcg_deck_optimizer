# cmd/parsecarddb

## Purpose

A helper command-line tool for exploring the upstream Flesh and Blood card database. It
parses the `card.csv` from `github.com/the-fab-cube/flesh-and-blood-cards` and prints
matching cards in either a pretty human-readable form or JSON. It is the preferred way to
look up card data and printed text when implementing a card — use it instead of grepping
the raw CSV.

## How it works

`card.csv` is a tab-separated export. `loadCards` reads it with a header-name-to-field
mapping (`cardCSVColumns`), so an upstream column reorder doesn't break parsing — a new
column just needs a mapping entry. `poolFilter.matches` then narrows the rows. The class /
talent filters express a hero's legal pool the way `registry.heroCanPlay` does: Generic and
classless cards are always class-legal, and a card's talents must be a subset of the hero's.
Class / talent / non-deck words are classified via `internal/fabtype`.

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
  dropping equipment, heroes, tokens, and landmarks.
- `-silver-age` — require Silver Age legality (default true).
- `-exclude-banned` — drop cards on our `internal/format` Silver Age banlist.
- `-format` — `pretty` (default) or `json`.
- `-names_only` — print only the distinct card names, one per line.

### Examples

```
# Viserai's pool (Runeblade + Generic, no talent, incl. weapons):
go run ./cmd/parsecarddb -classes Runeblade -modeled -names_only

# Aurora's Lightning cards worth implementing (incl. weapons, drops banned):
go run ./cmd/parsecarddb -classes Runeblade -talents Lightning -require-talents Lightning \
  -modeled -exclude-banned -names_only
```

## Important files

- `main.go` — the entire tool: the `Card` struct mirroring the CSV columns, the loader, the
  `poolFilter`, and the output formatting.

## Gotchas

- `-talents` is subset semantics, so the default (empty) admits **only** no-talent cards. Pass
  the hero's talents (or `-any-talent`) to widen it — a Lightning card won't appear under the
  default filter even with a matching `-name`.
