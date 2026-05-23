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
column just needs a mapping entry. `cardMatchesFilter` then narrows the rows: optional
user substring filters on Name and Types, the Silver Age legality gate (a blank legality
column is treated as legal), the class gate (Runeblade and Generic only — the optimizer's
current hero pool), and an exclusion list for types the optimizer doesn't model
(specialized elemental talents, tokens, equipment).

## How to use

```
go run ./cmd/parsecarddb [flags]
```

Flags:

- `-in` — path to `card.csv` (default `data_sources/card.csv`).
- `-name` — only print cards whose name contains this substring (case-insensitive).
- `-type` — only print cards whose Types field contains this substring (case-insensitive).
- `-format` — `pretty` (default) or `json`.
- `-names_only` — print only the distinct card names, one per line.

## Important files

- `main.go` — the entire tool: the `Card` struct mirroring the CSV columns, the loader, the
  filter, and the output formatting.

## Gotchas

- The class and type filters are deliberately narrow to the optimizer's current scope; a
  card outside Runeblade / Generic, or carrying an excluded type, will not appear even with
  a matching `-name`.
