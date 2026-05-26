# internal/card

The contract layer for every Flesh and Blood card. This package defines the `Card` interface
every card implements, the per-attack-step `CardState` wrapper that carries mutable flags
between resolution phases, the `CardType` / `TypeSet` type-line bitfield, and the narrow
`GameEngine` / `Logger` / `Aura` / `Item` / `EphemeralTrigger` interfaces cards consume from
the simulator. It owns the contract and deliberately does not import the sim — `*sim.TurnState`
and the sim's logger satisfy these interfaces structurally.

## Key types and interfaces

- `Card` (`card.go`) — the base contract. Static profile (`ID`, `Name`, `DisplayName`,
  `Cost`, `Pitch`, `Attack`, `Defense`, `Types`, `GoAgain`) plus a `Play` hook for on-resolve
  behaviour. Concrete implementations live in `cards/`.
- `CardState` (`card_state.go`) — wraps a `Card` with per-turn mutable flags (`Role`,
  `GrantedGoAgain`, `GrantedDominate`, `BonusAttack`, `BonusDefense`, `Mode`, `FromArsenal`,
  `PitchedToPlay`, `OnHit`, …). Created by the solver at the start of each attack-turn search and
  lives only for that attack turn. The card currently resolving receives its own `CardState` as the
  `self` argument to `Play`. `Effective*` methods fold printed values with granted bonuses
  (clamped at 0); helpers like `GrantGoAgainIfFromArsenal` and `GrantAttackReactionBuff`
  package the common rider plumbing.
- `GameEngine` (`interfaces.go`) — the rules-engine handle threaded through every card hook.
  Method-only, no fields: zone queries, token economy, value crediting, aura/trigger
  registration, draw/tutor verbs, hit heuristics. Card code needing raw field access
  type-asserts back to `*sim.TurnState`.
- `Logger` (`interfaces.go`) — the log sink. Cards use `AppendPostTrigger*` for self-riders
  and `AppendPreTrigger*` for hero/aura triggers; the `AppendAttackStep*` / `AmendLastAttackStepN`
  methods are sim-internal. A nil-pointer `Logger` (the find-best skip pass) silently elides
  every call.
- `Aura`, `Item`, `EphemeralTrigger` (`interfaces.go`) — minimal views the corresponding
  handler types see of the firing aura / item / trigger.
- Optional marker interfaces (`markers.go`) — `VariableCost`, `Modal`, `ModalCost`,
  `PlayPrecondition`, `Blocker`, `BlockCost`, `DefensiveInstant`, `Dominator`,
  `ArsenalDefenseBonus`, `ResourceSource`, `Universal`, `AttackReaction`, `LeavesArenaAura`.
  Cards opt into these to layer behaviour onto the base contract; the attack-turn runner
  type-asserts on them at the matching pipeline stage and skips the branch otherwise.
- `CardType` / `TypeSet` (`types.go`) — the type-line bitfield. Every type check is a single
  bitmask AND. `TypeSet` carries predicate helpers (`IsAttack`, `IsAttackAction`,
  `IsNonAttackAction`, `IsWeaponAttack`, `IsRunebladeAttack`, …) and the package exposes
  `func(GameEngine, *CardState) bool` predicate adapters of the same names for "next X" riders.
- `Role` / `CardAssignment` (`role.go`) — the partition role a card took on a turn (`Pitch`,
  `Attack`, `Defend`, `Held`, `Arsenal`).

## Important files

- `card.go` — the `Card` interface.
- `interfaces.go` — `GameEngine`, `Logger`, `Aura`, `Item`, `EphemeralTrigger`, `Hero`.
- `markers.go` — the optional marker interfaces.
- `card_state.go` — `CardState`, `OnHitHandler`, the `Effective*` helpers.
- `types.go` — `CardType`, `TypeSet`, predicate helpers.
- `role.go` — `Role`, `CardAssignment`.
- `cards/` — the implementations (see `cards/README.md`).

## Gotchas and invariants

- The package must not import the sim. The sim depends on `card`, never the reverse.
- `GameEngine` exposes no fields by design. Reaching for `*sim.TurnState` field access is the
  exception, not the rule, and pulls the sim import into the card file.
- `Cost(g)` must be deterministic in `g` for static-cost cards — non-`VariableCost` cards
  return the same value regardless of game state (enforced by a sim test).
- `Effective*` values clamp at 0; FaB attack power and defense cannot go negative.
- `Types(nil)` returns printed types only; `Universal` cards need a real `g` to fold in the
  active hero's class.
