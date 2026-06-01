# internal/fabtype

The text-side Flesh and Blood vocabulary: the sets of **class**, **talent**, and **non-deck
card-type** names as they appear in printed type lines and the upstream `card.csv` Types
column.

## Key values

- `ClassNames` — every hero class word (Runeblade, Wizard, …; Generic included).
- `TalentNames` — every talent word (Lightning, Shadow, Mystic, Royal, Revered, Reviled, …).
- `ClassMatches(word, heroClasses)` — class-legality for a single type word: Generic and
  non-class words always pass; a real class passes only for a hero that plays it. Generic is
  built in, so callers pass only the hero's own class(es).
- `NonDeckTypes` — printed card types that never enter the deck (Weapon, Equipment, Hero,
  Token, Landmark, …).
- `UnmodeledTypes` — `NonDeckTypes` minus `Weapon`: the types the optimizer doesn't model
  (it does model the weapon slot).

## Why it's separate from the engine

`internal/registry` decides what a hero can legally play with `classMask` / `talentMask` —
`TypeSet` bitmasks over the implemented `card.CardType` constants. That only covers types we
have cards for. This package is the broader **name** list (including classes / talents with no
implemented cards yet), used by the CSV / typebox tools — `cmd/parsecarddb` and
`cmd/crosscheck` — to classify a raw type word as a class, talent, or non-deck type.

Keep the two in sync: when a new class or talent ships, add its name here and (once a card
needs it) its `CardType` constant + mask entry in `internal/card` / `internal/registry`.
