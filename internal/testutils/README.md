# internal/testutils

## Purpose

Card, hero, and weapon fakes shared by tests in multiple packages (card, cards, deck, hand,
sim, turntests). It provides predictable, controllable implementations so predicate,
lookahead, partition, and chain-runner tests have known inputs without pulling in real
cards whose printed effects would perturb the measured value.

## Key types

- `Card` — a configurable `card.Card` fake. Zero-value fields mean "don't care"; tests set
  only the type / cost / power / pitch the helper under test predicates on. `GenericAttack`,
  `GenericAttackPitch`, `GenericAction` build common shapes.
- `FakeCard` — a fluent-builder card fake: `NewFakeCard(name)` plus `WithTypes` /
  `WithAttack` / `WithPitch` / `WithDefense` / `WithGoAgain`.
- Fixed-stat-line attack fakes — `RedAttack`, `YellowAttack`, `BlueAttack` (and pitch-only
  `RedPitch`, `BluePitch`) — deliberately simple cards used as deck contents when partition
  / ordering assertions need known optimal values.
- Typed predicate fakes — `RunebladeAttack`, `RunebladeWeapon`, `NonAttack`,
  `NonRunebladeAttack`, `AttackWithPower`, `Aura`, `FakeInstant`, `FakeNoGoAgainAttack` —
  each a fixed type line for one lookahead / predicate case.
- `Hero` — a minimal no-op hero with a configurable `Intel` (hand-draw size) and an
  injectable `OptStrategy`.
- `ClubWeapon` / `ClubWeaponAbility` / `HammerWeaponAbility` — weapon and weapon-ability
  fakes for types the real card pool lacks (Club, Hammer), so AR predicates that gate on
  those types can be exercised end-to-end.
- `GrantAll` / `GrantSpy` — a paired probe for detecting cross-permutation `CardState`
  wrapper leakage in the chain runner.
- `FireOnHitIfLikely` — fires every `OnHit` handler on a card when `LikelyToHit`, so a unit
  test that calls `Play` directly can exercise on-hit riders without the full chain runner.

## How to use / extend

Import the package from a test, pick the fake matching the case under test, and set only
the fields the assertion depends on. To add a fake, follow the existing pattern: a struct
implementing `card.Card` (or `sim.Weapon`) whose `ID()` returns `ids.InvalidCard`.

## Important files

- `cards.go` — the `Card` / `FakeCard` framework and every card fake.
- `hero.go` — the `Hero` fake.
- `weapons.go` — the Club / Hammer weapon and ability fakes.

## Gotchas

- Every fake returns `ids.InvalidCard` (or `InvalidWeapon` / `InvalidHero`) from `ID()`.
  Per-ID caches (`cardMetaCache`, `chainStepCache`) special-case `InvalidCard` so multiple
  fakes in one test don't interfere; the eval cache bails out whenever any input has an
  Invalid id (production cards always carry a unique non-zero ID).
