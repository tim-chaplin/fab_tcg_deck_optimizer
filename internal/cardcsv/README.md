# internal/cardcsv

The shared loader for the upstream the-fab-cube `card.csv` — a tab-separated export of every
Flesh and Blood card printing. Used by `cmd/parsecarddb` and `cmd/cardaudit` so the column
mapping lives in one place.

## Key types

- `Card` — one CSV row (one printing; a card with three pitch colours is three rows). Mirrors
  the CSV columns as string fields, plus a `String()` Stringer that pretty-prints non-blank
  fields.
- `Load(path) ([]Card, error)` — reads the TSV, mapping columns by **header name** (not
  position), so an upstream column reorder doesn't break parsing. A new upstream column just
  needs an entry in the unexported `columns` table.

## Gotchas

- Values are raw strings as they appear in the CSV — `Pitch`, `Cost`, `Power`, etc. are
  `"1"` / `"0"` / `""`, not ints. Callers parse as needed.
- `Load` is lenient (`LazyQuotes`, variable field count) to tolerate the upstream export's
  quoting quirks; it errors only if a mapped header column is missing entirely.
