# internal/textio

## Purpose

The durable on-disk text encodings the optimizer reads and writes. It converts a
`*deck.Deck` (plus its `deck.Stats`) to and from two formats and resolves the local
`mydecks/` directory path. Both encodings are pure data formatters with no sim or
gameengine dependency, so any consumer that needs to read or write a deck file can import
this package.

## Key types

- `DeckJSON` / `StatsJSON` / `PitchCountsJSON` / `BestTurnJSON` /
  `CardMarginalStatsJSON` — the on-disk JSON shapes. Every field trades a runtime interface
  value for a display-name string so files are human-readable and don't depend on
  card-registry indexing.

## How to use / extend

- JSON persistence (the canonical `mydecks/*.json` format): `MarshalDeck(d, stats)` encodes,
  `UnmarshalDeck(data)` decodes. Round-trips the recorded `Best` turn so saved decks render
  the same headline value as live runs.
- Fabrary plain text (fabrary.net import / export): `MarshalFabrary(d)` writes the text,
  `UnmarshalFabrary(text)` parses it. Stats are not part of the fabrary format — they
  round-trip through the JSON encoding only.
- Upstream card database (read-only): `LoadCardCSV(path)` parses the the-fab-cube `card.csv`
  into `CardCSV` rows (`cardcsv.go`). Unlike the deck encodings it doesn't touch the registry;
  it's the shared reader for the card-data tools (`cmd/parsecarddb`, `cmd/cardaudit`).
- Path resolution: `MydecksPath(name)` returns `mydecks/<name>.json` (a trailing `.json`
  on the name is stripped); `ValidateMydecksName` rejects path-traversal and
  Windows-reserved characters. `MydecksDir` is the relative directory constant.
- To add a new format, add a `<format>_marshal.go` / `<format>_unmarshal.go` pair that
  resolves card / weapon / hero names through `internal/registry` lookups.

## Important files

- `doc.go` — package overview of the two formats.
- `json_marshal.go` / `json_unmarshal.go` / `json_types.go` — the JSON encoding.
- `fabrary_marshal.go` / `fabrary_unmarshal.go` / `fabrary_names.go` — the fabrary text
  encoding and its pitch-suffix name canonicalisation.
- `mydecks.go` — `mydecks/` path resolution and name validation.
- `cardcsv.go` — `CardCSV` row type + `LoadCardCSV` for the upstream the-fab-cube `card.csv`.

## Gotchas

- Both encodings resolve every card / weapon / hero name through `internal/registry`.
  `UnmarshalDeck` fails loudly on an unrecognised name so a corrupted file doesn't produce
  silent nil-entry crashes downstream.
- `UnmarshalFabrary` is more lenient: deck cards the registry doesn't know are returned to
  the caller in a count-keyed `skipped` map rather than failing, so a silently-reduced deck
  surfaces to the user. A missing hero still aborts.
- `Sideboard` and `Equipment` are name-only string lists. The registry is not consulted for
  them, so the user can list equipment pieces and other items the simulator doesn't model.
- Card name lists are sorted on marshal so diffs across runs stay stable. `Avg` and the
  `Pitch` block in JSON are emitted for human readability only and ignored on unmarshal.
