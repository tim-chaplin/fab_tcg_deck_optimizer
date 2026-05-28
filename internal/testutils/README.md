# internal/testutils

## Purpose

Fake Card, Hero, and Weapon implementations shared by tests in multiple packages (card,
cards, deck, hand, sim, turntests). Predictable, controllable fakes so predicate,
lookahead, partition, and attack-turn runner tests have known inputs without pulling in real
cards whose printed effects would perturb the measured value.

## Key types

- `Fake` — single configurable `card.Card` implementation. Constructed via one of the
  colour-and-shape constructors below plus follow-up `With...` methods so every attribute
  the test cares about is visible at the call site.
- Colour-and-shape constructors — pitch is encoded in the constructor name (Red=1,
  Yellow=2, Blue=3) and shape (Attack / Action / Resource / DR / Instant / Aura) seeds
  the right base TypeSet. Defaults: cost=0, power=0, defense=0, no Go again, no Play
  side effect.
  - `FakeRedAttack` / `FakeYellowAttack` / `FakeBlueAttack` — generic Action-Attack.
  - `FakeRedAction` / `FakeYellowAction` / `FakeBlueAction` — generic non-attack action.
  - `FakeRedResource` / `FakeYellowResource` / `FakeBlueResource` — TypeResource (can
    only pitch, never plays).
  - `FakeRedDR` / `FakeYellowDR` / `FakeBlueDR` — defense reactions.
  - `FakeRedInstant` / `FakeYellowInstant` / `FakeBlueInstant` — generic instants.
  - `FakeRedAura` / `FakeYellowAura` / `FakeBlueAura` — auras.
  - `FakeRedBlocker` / `FakeYellowBlocker` / `FakeBlueBlocker` — Generic + Block, no
    Action / Attack / DefenseReaction (used by Opt-strategy tests that need a TypeBlock
    card competing with DRs for the defender slot).
  - `FakeWeaponSwing` — Generic Weapon+Attack base, no pitch; subtypes via `WithTypes`
    (Club, Hammer, Sword, …).
- Builder methods on `Fake` — `WithPower`, `WithCost`, `WithDefense`, `WithGoAgain`,
  `WithTypes(types ...card.CardType)` (variadic; ORs the given subtypes into the
  existing TypeSet — additive, for tagging on Runeblade / Hammer / Sword / …),
  `WithName` (override default name for log-assertion disambiguation), `WithPitch`
  (override the constructor's pitch; only for tests that iterate pitch dynamically),
  `WithPlay` (custom Play body), `WithDrawOne` (convenience for "this card draws on
  resolution").
- `Hero` — a minimal no-op hero with configurable `Intel` (hand-draw size) and an
  injectable `OptStrategy`.
- `ClubWeapon` — a 1-handed Club weapon (Weapon interface) whose `Ability()` returns a
  `FakeWeaponSwing` with Club + OneHand types. Used by turn-level tests that need an
  equipped weapon of a Club/Hammer type the card pool doesn't print.
- `GrantAll` / `GrantSpy` — a paired probe for detecting cross-permutation `CardState`
  wrapper leakage in the attack-turn runner. Build via `NewGrantAll()` / `NewGrantSpy(&saw)`.
- `FireOnHitIfLikely` — fires every `OnHit` handler on a card when `LikelyToHit`, so a
  unit test that calls `Play` directly can exercise on-hit riders without the full
  attack-turn runner.

## How to use / extend

Import the package from a test, pick the constructor matching the card shape under
test, and chain `With...` to set the attributes the assertion depends on. To add a new
shape (e.g. a printed-keyword card category), add one constructor per colour to
`fakes.go` seeded with the right base TypeSet.

## Important files

- `fakes.go` — the `Fake` builder type and every colour-and-shape constructor.
- `cards.go` — `GrantAll` / `GrantSpy` behaviour spies and shared helpers
  (`CardNamesSim`, `FireOnHitIfLikely`).
- `hero.go` — the `Hero` fake.
- `weapons.go` — the `ClubWeapon` weapon fake.

## Gotchas

- Every fake returns `ids.InvalidCard` (or `InvalidHero`) from `ID()` — the weapon fake
  `ClubWeapon` is a card now, so it too returns `InvalidCard`. Per-ID caches
  (`cardMetaCache`, `attackStepCache`) special-case `InvalidCard` so multiple fakes in one
  test don't interfere; the eval cache bails out whenever any input has an Invalid id
  (production cards always carry a unique non-zero ID).
- `Fake` is comparable (`a == b` works) — the play hook is stored behind a pointer so
  the struct itself stays comparable. The engine uses `==` in places like
  `RemoveFromHand`; if you write a fake outside the builder, keep that invariant.
- The colour constructors imply pitch (Red=1, Yellow=2, Blue=3). Real FaB cards all
  have pitch 1/2/3 so this covers the realistic space; tests that need an
  off-spectrum pitch override with `WithPitch(n)`.
