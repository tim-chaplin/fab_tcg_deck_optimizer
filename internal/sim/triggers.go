package sim

// Trigger types registered on a TurnState: AuraTrigger is the counter-tracked, cross-turn-
// capable hook firing on start-of-turn or attack-action events, and EphemeralAttackTrigger is
// the same-turn, fire-once "next attack" hook a card's Play can register.

// AuraTriggerType categorizes when an AuraTrigger's Handler fires. The sim walks the
// TurnState's AuraTriggers list on each matching condition and invokes every applicable
// handler.
type AuraTriggerType int

const (
	// TriggerStartOfTurn fires at the start of the owning player's action phase, before the
	// best-line search. The classic upkeep trigger for "at the beginning of your action phase
	// …" auras.
	TriggerStartOfTurn AuraTriggerType = iota
	// TriggerAttackAction fires each time an attack action card resolves during the attack
	// chain. Triggers that set OncePerTurn cap themselves at one fire per turn regardless of
	// how many attack actions resolve ("once per turn, when you play an attack action card
	// …" clauses).
	TriggerAttackAction
)

// OnAuraTrigger is the business-logic callback attached to an AuraTrigger. Called when the
// trigger's Type condition fires — it's where the printed "create a runechant", "gain 1{h}",
// "reveal top of deck" effect lives. Handlers mutate the passed TurnState directly
// (e.g. s.CreateRunechants, s.AddToGraveyard) and return the damage-equivalent that folds
// 1-to-1 into Value. The sim handles the counter bookkeeping (decrementing Count,
// graveyarding the aura when Count hits zero); the handler does not.
type OnAuraTrigger func(s *TurnState, t *AuraTrigger) int

// AuraTrigger is a counter-tracked handler attached to an aura in play. Each time Type's
// condition fires — and, when OncePerTurn is set, at most once per turn — the sim calls
// Handler and decrements Count. When Count reaches zero the sim sends Self to the graveyard
// and drops the trigger from TurnState.AuraTriggers. Self is the aura card itself so the
// sim can graveyard it without needing a back-reference.
type AuraTrigger struct {
	// Self is the aura card this trigger belongs to. Used by the sim to graveyard the aura
	// when Count reaches zero; also surfaced in per-turn summaries (e.g. the "(from previous
	// turn)" formatter line naming the aura that fired).
	Self Card
	// Type is the condition that fires this trigger.
	Type AuraTriggerType
	// Count is the number of times this trigger will still fire before the aura is destroyed.
	Count int
	// Handler runs when Type fires.
	Handler OnAuraTrigger
	// OncePerTurn caps the trigger at a single fire per turn regardless of how many matching
	// events occur. The sim sets FiredThisTurn the first time Handler runs each turn and
	// clears it at the next turn boundary.
	OncePerTurn bool
	// FiredThisTurn is sim-managed bookkeeping for OncePerTurn. Cards must not set it.
	FiredThisTurn bool
	// N is an optional small-integer payload available to Handler. Lets per-variant trigger
	// handlers (e.g. Malefic Incantation's per-color counter count) read their N off the
	// trigger instead of closing over it, so the handler can be a top-level function with no
	// per-Play closure allocation.
	N int
}

// OnEphemeralAttackTrigger is the business-logic callback attached to an
// EphemeralAttackTrigger. t is the trigger struct that owns the handler, exposing Source
// and the optional N payload to top-level handler implementations so cards don't need to
// allocate a closure per registration. target is the CardState of the attacker whose
// resolution triggered the fire; the handler may read target.Card.Attack(),
// target.EffectiveDominate(), etc. to decide whether a rider effect fires and what
// damage-equivalent to credit. Handlers mutate the passed TurnState directly
// (e.g. s.CreateRunechants) and return the damage-equivalent.
type OnEphemeralAttackTrigger func(s *TurnState, t *EphemeralAttackTrigger, target *CardState) int

// EphemeralAttackTrigger is a same-turn, fire-once "next attack" trigger registered by a
// card's Play (via TurnState.AddEphemeralAttackTrigger). Fires on the next attack action
// whose resolution matches its Matches predicate, AFTER the attacker's Play, hero
// OnCardPlayed, and AuraTriggers have all settled — so the Handler sees the fully-resolved
// attacker state (incl. any Dominate grants and hero-created auras).
//
// Distinct from AuraTrigger on three axes:
//   - Fire-once. No Count / OncePerTurn bookkeeping; the trigger resolves and drops out.
//   - Doesn't graveyard a source when it fires or fizzles — the registering card was
//     already graveyarded on its own resolution; only the trigger effect "stays in play."
//   - Doesn't persist across turns. Non-matching attack actions leave the trigger in place
//     for a later match, but anything unresolved at end of turn fizzles silently.
//
// The Source card keeps damage attribution clean: Handler's return is credited to Source's
// position in the chain (via SourceIndex), so a trigger fired by Mauvrion Skies during
// Drowning Dire's attack surfaces as damage on Mauvrion's BestLine entry rather than DD's.
type EphemeralAttackTrigger struct {
	// Source is the card that registered the trigger. Damage the handler returns accrues to
	// Source's per-card attribution; also surfaces in per-turn debug output.
	Source Card
	// Matches decides whether the trigger fires on a given attack action card's
	// resolution. Nil matches any attack action — non-matching resolutions leave the
	// trigger in place for a later attack.
	Matches func(target *CardState) bool
	// Handler runs when the trigger fires. Returns damage-equivalent credited to Source.
	Handler OnEphemeralAttackTrigger
	// SourceIndex is sim-managed bookkeeping: the position of Source in the played-chain
	// permutation, used to route Handler's damage back to Source's perCardOut slot. Cards
	// must not set it — the solver stamps it on registration.
	SourceIndex int
	// N is an optional small-integer payload available to Handler. Lets per-variant trigger
	// handlers (Mauvrion Skies's "create N runechants on hit") read their N off the trigger
	// instead of closing over it, so the handler can be a top-level function with no per-Play
	// closure allocation.
	N int
}
