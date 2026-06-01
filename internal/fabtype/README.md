# internal/fabtype

The text-side Flesh and Blood vocabulary: the sets of **class**, **talent**, and **non-deck
card-type** names as they appear in printed type lines and the upstream `card.csv` Types
column.

## Key values

- `ClassNames` — every hero class word (Runeblade, Wizard, …; Generic included).
- `TalentNames` — every talent word (Lightning, Shadow, Earth, …). Royal is intentionally
  absent — it's a supertype, not a talent.
- `NonDeckTypes` — printed card types that never enter the deck (Weapon, Equipment, Hero,
  Token, Landmark, …).

## Why it's separate from the engine

`internal/registry` decides what a hero can legally play with `classMask` / `talentMask` —
`TypeSet` bitmasks over the implemented `card.CardType` constants. That only covers types we
have cards for. This package is the broader **name** list (including classes / talents with no
implemented cards yet), used by the CSV / typebox tools — `cmd/parsecarddb` and
`cmd/crosscheck` — to classify a raw type word as a class, talent, or non-deck type.

Keep the two in sync: when a new class or talent ships, add its name here and (once a card
needs it) its `CardType` constant + mask entry in `internal/card` / `internal/registry`.
