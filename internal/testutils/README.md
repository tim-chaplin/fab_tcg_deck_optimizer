# internal/testutils

## Purpose

Card, hero, and weapon stubs shared by tests in multiple packages (card, cards, deck, hand,
sim, turntests). It provides predictable, controllable fakes so predicate, lookahead,
partition, and chain-runner tests have known inputs without pulling in real cards whose
printed effects would perturb the measured value.

## Key types

- `Card` — a configurable `card.Card` stub. Zero-value fields mean "don't care"; tests set
  only the type / cost / power / pitch the helper under test predicates on. `GenericAttack`,
  `GenericAttackPitch`, `GenericAction` build common shapes.
- `StubCard` — a fluent-builder card stub: `NewStubCard(name)` plus `WithID` / `WithTypes` /
  `WithAttack` / `WithPitch` / `WithDefense` / `WithGoAgain`.
- Fixed-stat-line attack fakes — `RedAttack`, `YellowAttack`, `BlueAttack` (and pitch-only
  `RedPitch`, `BluePitch`) — deliberately simple cards used as deck contents when partition
  / ordering assertions need known optimal values.
- Typed predicate fakes — `RunebladeAttack`, `RunebladeWeapon`, `NonAttack`,
  `NonRunebladeAttack`, `AttackWithPower`, `Aura`, `InstantStub`, `NoGoAgainAttackStub` —
  each a fixed type line for one lookahead / predicate case.
- `Hero` — a minimal no-op hero with a configurable `Intel` (hand-draw size) and an
  injectable `OptStrategy`.
- `ClubWeapon` / `ClubWeaponAbility` / `HammerWeaponAbility` — weapon and weapon-ability
  stubs for types the real card pool lacks (Club, Hammer), so AR predicates that gate on
  those types can be exercised end-to-end.
- `GrantAll` / `GrantSpy` — a paired probe for detecting cross-permutation `CardState`
  wrapper leakage in the chain runner.
- `FireOnHitIfLikely` — fires every `OnHit` handler on a card when `LikelyToHit`, so a unit
  test that calls `Play` directly can exercise on-hit riders without the full chain runner.

## How to use / extend

Import the package from a test, pick the stub matching the case under test, and set only
the fields the assertion depends on. To add a fake, follow the existing pattern: a struct
implementing `card.Card` (or `sim.Weapon`), with a synthetic ID from `ids.go` when it needs
a distinct cache slot.

## Important files

- `cards.go` — the `Card` / `StubCard` framework and every card fake.
- `hero.go` — the `Hero` stub.
- `weapons.go` — the Club / Hammer weapon and ability stubs.
- `ids.go` — the `Fake*` synthetic card IDs.

## Gotchas

- Synthetic IDs (`FakeRedAttack`, …) are anchored in `ids.go` past every real card, weapon,
  and token-ability ID so a test fake never shares an ID-keyed cache slot (`cardMetaCache`,
  `chainStepCache`) with a real printing or another fake.
- Many fakes return `ids.InvalidCard` from `ID()`. A test that reaches into an ID-keyed
  cache must attach a distinct ID with `WithID(testutils.FakeX)` to avoid sharing slot 0.
