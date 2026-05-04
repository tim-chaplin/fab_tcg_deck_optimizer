package sim

import "github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"

// Aura is a persistent hook attached to a card or a token in play. The sim walks each
// TurnState's Auras list on every Trigger condition and fires the matching handlers;
// lifecycle (when to decrement Count, when to send Self to the graveyard, when to
// deregister) belongs to the handler. Used for start-of-turn upkeep auras (sigils,
// Blessing of Occult, Runeblood Incantation), per-attack-action triggers (Malefic
// Incantation), and aura tokens (Runechants — see tokens.go).

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
	// TriggerAttack fires when ANY attack resolves — attack action card or weapon swing.
	// The runechant aura uses this so each attack consumes the runechants in play.
	TriggerAttack
)

// CardOrTokenType identifies what an Aura belongs to: a specific card in play, or an
// aura token. Exactly one of Card / TokenType is set; the other carries its zero value.
type CardOrTokenType struct {
	// Card is the originating card for non-token auras. nil for token auras.
	Card Card
	// TokenType identifies the token kind for token auras. TokenTypeNone for card auras.
	TokenType TokenType
}

// IsToken reports whether this Aura belongs to a token (no originating card).
func (c CardOrTokenType) IsToken() bool { return c.TokenType != TokenTypeNone }

// CardID returns the originating card's ID when this is a card aura, or ids.InvalidCard
// for tokens. Tokens are distinguished in cache keys by their TokenType, not by CardID.
func (c CardOrTokenType) CardID() ids.CardID {
	if c.Card != nil {
		return c.Card.ID()
	}
	return ids.InvalidCard
}

// DisplayName returns the human-readable name — the card's DisplayName for card auras,
// or the token's printed name (e.g. "Runechant") for token auras.
func (c CardOrTokenType) DisplayName() string {
	if c.Card != nil {
		return DisplayName(c.Card)
	}
	return tokenDisplayName(c.TokenType)
}

// AuraHandler is the business-logic callback attached to an Aura. Called when the Aura's
// TriggerType condition fires — it's where the printed effect lives. Handlers mutate the
// passed TurnState directly, crediting damage / life gain via s.AddValue. Same shape as
// Card.Play — no return. Lifecycle is the handler's responsibility: call
// s.DestroyAura(t, addToGraveyard) when done. Token-style auras pass false for
// addToGraveyard since the underlying token isn't a card and just disappears.
type AuraHandler func(s *TurnState, t *Aura)

// Aura is one persistent hook attached to a card or a token in play. Each time
// TriggerType's condition fires — and, when OncePerTurn is set, at most once per turn —
// the sim calls Handler. The Aura survives until its handler calls s.DestroyAura.
type Aura struct {
	// Self identifies what this Aura belongs to — a card or a token type. Surfaced in
	// per-turn summaries via CardOrTokenType.DisplayName.
	Self CardOrTokenType
	// TriggerType is the trigger condition that fires this Aura's Handler.
	TriggerType TriggerType
	// Count is a per-Aura counter. Its meaning is card-specific: Malefic Incantation and
	// Runeblood Incantation read it as fires remaining and decrement themselves; one-shot
	// sigils ignore it; token auras (Runechant) use it as a copy count. The sim treats
	// Count as opaque storage and never mutates it.
	Count int
	// Handler runs when TriggerType fires.
	Handler AuraHandler
	// OncePerTurn caps the Handler at a single fire per turn regardless of how many matching
	// events occur. The sim sets FiredThisTurn the first time Handler runs each turn and
	// clears it at the next turn boundary.
	OncePerTurn bool
	// FiredThisTurn is sim-managed bookkeeping for OncePerTurn. Cards must not set it.
	FiredThisTurn bool
}
