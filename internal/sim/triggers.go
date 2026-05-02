package sim

// Trigger types registered on a TurnState. AuraTrigger is the counter-tracked, cross-turn-
// capable hook firing on start-of-turn or attack-action events.

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
	// LogText is the optional pre-built rider-line text the handler emits via LogPostTrigger
	// / LogPreTrigger. Card.Play computes it once at registration (typically
	// `"<DisplayName> <verb phrase>"`); the handler reads t.LogText directly so the hot fire
	// path runs zero string allocations even when the chain is materialising the log.
	// Empty-string means the handler authors its own text dynamically.
	LogText string
}
