# internal/weapon

## Purpose

Owns the Flesh and Blood weapon model, mirroring the `internal/aura` / `internal/item`
split: a platonic weapon **card** plus the mutable engine-side **object** built when that
card is equipped. Weapons are equipment that sits in the arena, not deck cards. A deck
equips 0-2 weapons under the "2 x 1H or 1 x 2H" rule. Concrete weapons live in the
`internal/weapon/weapons` subpackage.

This package must NOT import `gameengine` (same layering rule as `internal/aura` /
`internal/item`); the concrete object's `Destroy` routes back through `card.GameEngine`,
which `gameengine.GameEngine` satisfies structurally.

## Key types

- `Card` — the platonic weapon card: a full `card.Card` (so it carries a `CardID`,
  `DisplayName`, `Types(GameEngine)`, lives in the card universe, and goes to the graveyard
  when destroyed) plus `Hands()` and `Ability()`.
- `Weapon` — the concrete mutable object the engine stores in `GameState.weapons`. Embeds
  `trigger.Trigger` (so an end-phase handler can subscribe — Talishar's self-destruct),
  carries the per-turn counter total (`Count` / `SetCount`), the source card (`SourceCard`),
  and caches the card's `Name` / `Hands` / `Ability`. Built at game start via
  `gameengine.EquipFromCards`, deep-copied per permutation via `Copy`, destroyed via
  `GameEngine.DestroyWeapon`. Satisfies `gameengine.Weapon` and `card.Weapon`.

## Weapon activated abilities

Each weapon is two paired Go types in one `<weapon>.go` file:

- A platonic weapon `Card` (`ID`, `DisplayName`, `Cost`/`Pitch`/`Attack`/`Defense` all 0,
  `Types`, `Hands`, `Ability`, plus a `Play` that registers the equipped object via
  `ge.CreateWeapon` at equip time) — the equipped permanent that sits in the arena and never
  enters the attack turn.
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

IDs are anchored in `internal/ids/weapon_ids.go`: weapon permanents take `CardID` values
past the last deck card (weapons are cards, so they live in the shared `CardID` space);
ability `CardID` values anchor past the last weapon-permanent ID; token activated abilities
anchor past the last weapon-ability ID; test fakes anchor past the last token-ability ID.

Card-attack predicates (`internal/card/types.go`) gate purely on `TypeAttack` (and
`TypeWeapon` / `TypeRuneblade` where the rider needs the subtype). A bare weapon permanent
— types only, no `TypeAttack` — never matches.

## Equip and end-of-turn plumbing

Weapons are equipped at game start, not played from hand. A weapon card's `Play` is its
equip-time registration hook: it calls `card.GameEngine.CreateWeapon` — the `CreateAura` /
`CreateItem` counterpart — to put its mutable object into `GameState.weapons`, passing a
trigger (firing event + handler) or `(0, nil)` for an untriggered weapon. `EquipFromCards`
drives this by playing each platonic weapon card (the sim equips a deck's weapons in
`weaponsFromDeck`; tests reach for `StateBuilder.EquipWeapons`). `FireTriggers` walks equipped
weapons like auras / items; Talishar registers an `EndOfTurn` handler (its rust-counter
self-destruct), so the walk fires for any loadout that includes it.

## Weapon markers

Unlike cards, which keep their `NotImplemented` / `Unplayable` markers in dedicated
`notimplemented/` and `unplayable/` subdirectories, an unimplemented or unplayable weapon
stays as a normal file in `internal/weapon/weapons/` and carries a `NotImplemented()` /
`Unplayable()` marker method directly on the platonic `Card` type. There is no weapon
subpackage for markers. Pair a `NotImplemented()` with a `// not implemented: <quirk>`
comment; Rosetta Thorn is the model.

## How to use / extend

To add a weapon, drop a `<weapon>.go` file in `internal/weapon/weapons/` with the two
paired types described above, anchor the new IDs in `internal/ids/weapon_ids.go`, and add
it to the weapon registry. `TestRegistry_CoversEveryWeapon` (in `internal/lint`) fails if a
weapon is left out of `registry.AllWeapons`; it discovers weapon types by their `Hands`
method.

## Important files

- `weapon.go` — the `Card` and `Weapon` types, the equip-time builders, and the package doc.
- `weapons/nebula_blade.go` — a fully-modelled weapon (on-hit Runechant rider, conditional
  +3 power); the model for new weapon files.
- `weapons/talishar.go` — a self-triggering weapon: its `Play` registers an `EndOfTurn`
  self-destruct via `ge.CreateWeapon`; the model for weapons that subscribe a trigger.
- `weapons/rosetta_thorn.go` — a `NotImplemented` / `NotSilverAgeLegal`-marked weapon, the
  marker model.

## Gotchas

- The weapon permanent itself never enters the attack turn; only its `Ability()` Card does.
- Weapons are mutable once equipped (counters), so the engine deep-copies them per
  permutation — never share a `Weapon` object across two GameStates.
- Each weapon's top docstring must carry the printed `Text: "..."` block and any modelling
  fudge (e.g. Reaping Blade ignores its health-symmetry rider).
