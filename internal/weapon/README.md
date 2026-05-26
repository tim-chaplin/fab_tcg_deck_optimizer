# internal/weapon

## Purpose

Defines the `Weapon` interface for Flesh and Blood weapons — equipment that sits in the
arena, not deck cards. A deck equips 0-2 weapons under the "2 x 1H or 1 x 2H" rule.
Concrete weapons live in the `internal/weapon/weapons` subpackage.

## Key types

- `Weapon` — `ID`, `Name`, `Types`, `Hands`, and `Ability`. The equipped permanent.

## Weapon activated abilities

Each weapon is two paired Go types in one `<weapon>.go` file:

- A `Weapon` permanent (`ID`, `Name`, `Types`, `Hands`, `Ability`) — the equipped
  permanent that sits in the arena and never enters the attack turn.
- A `card.Card` activated ability (`<Weapon>Ability`) — what the attack-turn runner enqueues
  each turn when the player swings (1 AP, pays the ability's printed activation cost). It
  carries the printed Cost / Pitch / Attack / Defense / GoAgain / Play plus the parent's
  `card.TypeSet` with `card.TypeAttack` added, so `IsAttack`, `IsWeaponAttack`, and — for
  Runeblade weapons — `IsRunebladeAttack` all fire.

`Ability()` returns a package-level cached `card.Card`: each weapon file declares
`var <weapon>Ability card.Card = <Weapon>Ability{}` and returns it. The cache keeps the
lookup allocation-free, since Go's interface boxing of a zero-size struct allocates per
call when the result escapes. The ability's `Name()` matches the weapon's so attack-turn logs
read `<Weapon>: WEAPON ATTACK`.

IDs are anchored in `internal/ids/weapon_ids.go`: weapon permanents take `WeaponID`
values past the last card ID; ability `CardID` values anchor past the last weapon-permanent
ID; token activated abilities anchor past the last weapon-ability ID; test fakes anchor
past the last token-ability ID. `WeaponID` is an alias of `CardID` because every weapon
swing flows through the same attack-turn runner pipeline and per-card caches as deck cards.

Card-attack predicates (`internal/card/types.go`) gate purely on `TypeAttack` (and
`TypeWeapon` / `TypeRuneblade` where the rider needs the subtype). A bare weapon permanent
— types only, no `TypeAttack` — never matches.

## Weapon markers

Unlike cards, which keep their `NotImplemented` / `Unplayable` markers in dedicated
`notimplemented/` and `unplayable/` subdirectories, an unimplemented or unplayable weapon
stays as a normal file in `internal/weapon/weapons/` and carries a `NotImplemented()` /
`Unplayable()` marker method directly on the `Weapon` type. There is no weapon subpackage
for markers. Pair a `NotImplemented()` with a `// not implemented: <quirk>` comment;
Annals of Sutcliffe is the model.

## How to use / extend

To add a weapon, drop a `<weapon>.go` file in `internal/weapon/weapons/` with the two
paired types described above, anchor the new IDs in `internal/ids/weapon_ids.go`, and add
it to the weapon registry. `TestRegistry_CoversEveryWeapon` (in `internal/lint`) fails if a
weapon is left out of `registry.AllWeapons`.

## Important files

- `weapon.go` — the `Weapon` interface and package doc.
- `weapons/nebula_blade.go` — a fully-modelled weapon (on-hit Runechant rider, conditional
  +3 power); the model for new weapon files.
- `weapons/annals_of_sutcliffe.go` — a `NotImplemented`-marked weapon, the marker model.

## Gotchas

- The weapon permanent itself never enters the attack turn; only its `Ability()` Card does.
- Each weapon's top docstring must carry the printed `Text: "..."` block and any modelling
  fudge (e.g. Reaping Blade ignores its health-symmetry rider).
