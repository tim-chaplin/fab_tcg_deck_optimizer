# internal/hero

## Purpose

Defines the `Hero` interface a Flesh and Blood hero card satisfies. A hero with a printed
"whenever X" ability embeds `trigger.Trigger[card.Hero]` and is walked by the engine's
`FireTriggers` alongside auras, items, and one-shot triggers. Concrete heroes live in the
`internal/hero/heroes` subpackage.

## Key types

- `Hero` — `ID`, `Name`, `Class`, `Types`, `Intelligence`, `Opt`, plus the trigger-dispatch
  surface (`TriggerType`, `OncePerTurn`, `FiredThisTurn`, `SetFiredThisTurn`, `Matches`,
  `Fire`). The trigger methods mirror what auras and items expose through their embedded
  `trigger.Trigger`.

## How to use / extend

To add a hero, drop a `<hero>.go` file in `internal/hero/heroes/`:

- Define an unexported `<hero>Hero` struct embedding `trigger.Trigger[card.Hero]`.
- Expose the hero as a package-level pointer var (e.g. `var Viserai =
  &viseraiHero{Trigger: trigger.FromHero[card.Hero](...)}`) so callers reference the
  singleton with `heroes.<Hero>`.
- Implement `ID`, `Name`, `Class`, `Types`, `Intelligence`, `Opt`, and `Fire` on
  `*<hero>Hero`. `Fire` calls `Invoke(engine, logger, self)` — symmetric with
  `*aura.Aura.Fire` and `*item.Item.Fire`.
- The trigger's `TypeFilter` narrows the firing site (e.g. Viserai filters to Runeblade
  non-Weapon). The handler reads further context off the engine —
  `NonAttackActionPlayed`, `TriggeringCard`, etc.

Heroes with no triggered ability (`defaultHero`, test stubs) provide no-op trigger
methods returning zero values; `TriggerType == 0` short-circuits dispatch.

## Important files

- `hero.go` — the `Hero` interface and package doc.
- `heroes/viserai.go` — Young Viserai: a `CardOrAbility` trigger that creates a Runechant
  on Runeblade-card-after-non-attack-action; the model for new hero files.

## Gotchas

- A Runeblade weapon swing must not trigger Viserai. The trigger's `TypeFilter`
  (`viseraiTypeFilter`) excludes `TypeWeapon` so equipping or swinging doesn't qualify.
- The `Opt` heuristic reads printed values only — `Types(nil)` and `GoAgain(nil)` skip
  any runtime-state fold, which is the right thing for an Opt-time decision.
