package sim

// An Aura is a persistent hook attached to a card in play. The sim walks each TurnState's
// Auras list on every Trigger condition and fires the matching handlers; lifecycle (when
// to decrement Count, when to send Self to the graveyard, when to deregister) belongs to
// the handler — sim only keeps Auras alive between fires and drops them once the handler
// flips Destroyed. Used today for start-of-turn upkeep auras (Sigil of Deadwood, Sigil
// of Fyendal, Blessing of Occult, Runeblood Incantation, Sigil of the Arknight, Sigil
// of Silphidae) and per-attack-action triggers (Malefic Incantation).

// TriggerType categorises when an Aura's Handler fires.
type TriggerType int

const (
	// TriggerStartOfTurn fires at the start of the owning player's action phase, before the
	// best-line search. The classic upkeep trigger for "at the beginning of your action phase
	// …" auras.
	TriggerStartOfTurn TriggerType = iota
	// TriggerAttackAction fires each time an attack action card resolves during the attack
	// chain. Auras that set OncePerTurn cap themselves at one fire per turn regardless of
	// how many attack actions resolve ("once per turn, when you play an attack action card
	// …" clauses).
	TriggerAttackAction
)

// AuraHandler is the business-logic callback attached to an Aura. Called when the Aura's
// TriggerType condition fires — it's where the printed "create a runechant", "gain 1{h}",
// "reveal top of deck" effect lives. Handlers mutate the passed TurnState directly
// (e.g. s.CreateRunechants, s.AddToGraveyard) and return the damage-equivalent that folds
// 1-to-1 into Value. Lifecycle is the handler's responsibility: a one-shot aura calls
// s.DestroyAura(t, addToGraveyard) at the end of its body; a counter-based aura decrements
// t.Count and calls s.DestroyAura when the count expires.
type AuraHandler func(s *TurnState, t *Aura) int

// Aura is one persistent hook attached to a card in play. Each time TriggerType's
// condition fires — and, when OncePerTurn is set, at most once per turn — the sim calls
// Handler. The Aura survives until its handler calls s.DestroyAura.
type Aura struct {
	// Self is the card this Aura belongs to. Surfaced in per-turn summaries (e.g. the
	// "(from previous turn)" formatter line naming the Aura that fired). Handlers that
	// want the underlying card to land in the graveyard call s.DestroyAura(t, true).
	Self Card
	// TriggerType is the trigger condition that fires this Aura's Handler.
	TriggerType TriggerType
	// Count is a per-Aura counter. Its meaning is card-specific: Malefic Incantation and
	// Runeblood Incantation read it as fires remaining and decrement themselves; one-shot
	// sigils ignore it; future Auras may use it for other things (e.g. tokens in play).
	// The sim treats Count as opaque storage and never mutates it.
	Count int
	// Handler runs when TriggerType fires.
	Handler AuraHandler
	// OncePerTurn caps the Handler at a single fire per turn regardless of how many matching
	// events occur. The sim sets FiredThisTurn the first time Handler runs each turn and
	// clears it at the next turn boundary.
	OncePerTurn bool
	// FiredThisTurn is sim-managed bookkeeping for OncePerTurn. Cards must not set it.
	FiredThisTurn bool
	// N is an optional small-integer payload available to Handler. Lets per-variant Auras
	// (e.g. Malefic Incantation's per-color counter count) read their N off the Aura
	// instead of closing over it, so the handler can be a top-level function with no
	// per-Play closure allocation.
	N int
}
