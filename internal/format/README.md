# internal/format

The constructed gameplay formats — the deck-construction rulesets that decide which cards are
legal to include in a deck. Today only Silver Age exists.

## Key types

- `Format` (`format.go`) — an interface: `Name()` (stable id for CLI flags / deck filenames,
  e.g. `"silver_age"`), `DisplayName()` (`"Silver Age"`, for deck headers), and
  `IsCardLegal(Card) bool` (the per-card legality predicate). A `Format` is **hero-agnostic**:
  class/talent legality is a universal deckbuilding rule owned by `internal/registry`, not a
  per-format concern. `Format` knows only about banlists (and, eventually, copy caps / a legal
  pool).
- `SilverAge` (`format.go`) — the one live format. `IsCardLegal` returns "name not on the
  banlist".
- `Card` (`format.go`) — the minimal view a format needs of a card: just `Name()`. Declared
  locally so the package stays decoupled from `internal/card` (concrete `card.Card` values
  satisfy it).
- `Parse(string) (Format, error)` (`format.go`) — CLI flag value → `Format`, with a
  self-describing error for unknown formats.

## The banlist

`banlist.go` holds `silverAgeBanlist`, a hardcoded `map[string]struct{}` of banned card names
— the single source of truth for Silver Age legality. Update it manually (the official banlist
changes infrequently).

**Names must match `card.Card.Name()` exactly**: base printed name, no pitch suffix, straight
ASCII apostrophe (`'`). A mismatch silently leaks a banned card back into the pool.
`registry.TestLegalCards_ExcludesFormatBanned` and this package's `format_test.go` guard the
implemented banned cards.

## How it's consumed

- `internal/registry` — `LegalCards` / `LegalWeapons` apply `format.SilverAge.IsCardLegal`
  (ANDed with the marker filter, and — in a later change — the hero-pool class/talent check).
- `cmd/fabsim` — the `-format` flag parses via `Parse`; the format scopes a run and names its
  output decks (`<hero>_<format>_<incoming>`).
