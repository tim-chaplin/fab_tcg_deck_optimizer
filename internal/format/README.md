# internal/format

The constructed gameplay formats — the deck-construction rulesets that decide which cards are
legal to include in a deck. Today only Silver Age exists.

## Key types

- `Format` (`format.go`) — an interface: `Name()` (stable id for CLI flags / deck filenames,
  e.g. `"silver_age"`), `DisplayName()` (`"Silver Age"`, for deck headers), `IsCardLegal(Card)
  bool` (the name-only banlist predicate), and `IsRarityLegal(rarity string) bool` (the rarity
  rule). A `Format` is **hero-agnostic**: class/talent legality is a universal deckbuilding rule
  owned by `internal/registry`, not a per-format concern. `Format` knows only about banlists and
  rarities (and, eventually, copy caps / a legal pool).
- `SilverAge` (`format.go`) — the one live format. `IsCardLegal` returns "name not on the
  banlist"; `IsRarityLegal` returns "rarity is Basic, Common or Rare".
- `Card` (`format.go`) — the minimal view a format needs of a card: just `Name()`. Declared
  locally so the package stays decoupled from `internal/card` (concrete `card.Card` values
  satisfy it). Rarity is passed to `IsRarityLegal` as a plain string rather than added here, so
  name-only callers (e.g. `cmd/parsecarddb`'s banlist filter) still satisfy `Card`.
- `Parse(string) (Format, error)` (`format.go`) — CLI flag value → `Format`, with a
  self-describing error for unknown formats.

## The banlist

`banlist.go` holds `silverAgeBanlist`, a hardcoded `map[string]struct{}` of banned card names —
the source of truth for the banlist half of Silver Age legality. Update it manually (the
official banlist changes infrequently).

**Names must match `card.Card.Name()` exactly**: base printed name, no pitch suffix, straight
ASCII apostrophe (`'`). A mismatch silently leaks a banned card back into the pool.
`registry.TestLegalCards_ExcludesFormatBanned` and this package's `format_test.go` guard the
implemented banned cards.

## The rarity rule

The other half: a card is Silver Age legal only if its **lowest printed rarity** is Basic,
Common or Rare (so a card printed at both Majestic and Rare qualifies on its Rare printing).
`IsRarityLegal` is the predicate. Each card's lowest rarity is baked into the generated card as
`Rarity()` (see `internal/cardgen`), sourced from the upstream printing data; the registry reads
it off the card and feeds the string here. This catches cards the banlist misses — notably new
sets, whose `card.csv` "Silver Age Legal" column is blank.

## How it's consumed

- `internal/registry` — `LegalCardsFor` / `LegalWeaponsFor` AND `IsCardLegal` (banlist) with
  `IsRarityLegal` (the card's `Rarity()`), the marker filter, and the registry's own hero-pool
  class/talent check.
- `cmd/fabsim` — the `-format` flag parses via `Parse`; the format is a durable deck attribute
  (persisted by `internal/textio`) that scopes a run and names its output decks
  (`<hero>_<format>_<incoming>`).
