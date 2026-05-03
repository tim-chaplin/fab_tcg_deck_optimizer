package sim

// The Aura framework: a persistent, counter-tracked hook attached to a card in play.
// Each Aura fires its Handler when the matching Trigger condition occurs and decrements
// Count; when Count hits zero the owning card heads to the graveyard. Used today for
// start-of-turn upkeep auras (Sigil of Deadwood, Sigil of Fyendal, Blessing of Occult,
// Runeblood Incantation, Sigil of the Arknight, Sigil of Silphidae) and per-attack-
// action triggers (Malefic Incantation).

// TriggerType categorises when an Aura's Handler fires. The sim walks the TurnState's
// Auras list on each matching condition and invokes every applicable handler.
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
// Type condition fires — it's where the printed "create a runechant", "gain 1{h}",
// "reveal top of deck" effect lives. Handlers mutate the passed TurnState directly
// (e.g. s.CreateRunechants, s.AddToGraveyard) and return the damage-equivalent that folds
// 1-to-1 into Value. The sim handles the counter bookkeeping (decrementing Count,
// graveyarding the aura when Count hits zero); the handler does not.
type AuraHandler func(s *TurnState, t *Aura) int

// Aura is a persistent hook attached to a card in play. Each time Type's condition fires
// — and, when OncePerTurn is set, at most once per turn — the sim calls Handler and
// decrements Count. When Count reaches zero the sim sends Self to the graveyard and drops
// the Aura from TurnState.Auras. Self is the originating card so the sim can graveyard it
// without needing a back-reference.
type Aura struct {
	// Self is the card this Aura belongs to. Used by the sim to graveyard the card when
	// Count reaches zero; also surfaced in per-turn summaries (e.g. the "(from previous
	// turn)" formatter line naming the Aura that fired).
	Self Card
	// Type is the trigger condition that fires this Aura's Handler.
	Type TriggerType
	// Count is how many fires remain before the Aura is destroyed. Card-specific meaning:
	// for Malefic Incantation it's the verse counter; for runechant-creating sigils it's
	// always 1 (re-armed each turn).
	Count int
	// Handler runs when Type fires.
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
	// LogText is the optional pre-built rider-line text the handler emits via LogPostTrigger
	// / LogPreTrigger. Card.Play computes it once at registration (typically
	// `"<DisplayName> <verb phrase>"`); the handler reads t.LogText directly so the hot fire
	// path runs zero string allocations even when the chain is materialising the log.
	// Empty-string means the handler authors its own text dynamically.
	LogText string
}
