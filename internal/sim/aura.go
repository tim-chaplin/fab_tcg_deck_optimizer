package sim

import "github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"

// Aura is a persistent hook attached to a card or a token in play. The sim walks each
// TurnState's Auras list on every Trigger condition and fires the matching handlers;
// lifecycle (when to decrement Count, when to send Self to the graveyard, when to
// deregister) belongs to the handler. Auras compose a Trigger by embed: an Aura is a
// Trigger plus persistent-instance state. Standalone (one-shot) Triggers live on
// TurnState.Triggers and are removed by the sim after firing — see trigger.go.

// TriggerType categorises when a Trigger / Aura's Handler fires.
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
	// TriggerEndOfTurn fires after the chain finishes, before the carry-state snapshot
	// and the post-hoc arsenal-promotion step. Ponder uses this to draw a card into the
	// held cards so the existing arsenal-promotion logic fills an otherwise-empty slot.
	TriggerEndOfTurn
	// TriggerHit fires when an attack hits (LikelyToHit on the post-AR-buff EffectiveAttack).
	// Used by "the next time an X you control hits this turn, do Y" printed text (Plunder
	// Run, High Striker). Trigger.TypeFilter narrows qualifying hits to a card-type
	// predicate. Standalone-trigger only — auras don't currently use this trigger type.
	TriggerHit
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
		return c.Card.DisplayName()
	}
	return tokenDisplayName(c.TokenType)
}

// Aura is one persistent hook attached to a card or a token in play. The embedded
// Trigger carries the firing-data (Source, TriggerType, TypeFilter, Handler); the
// fields below add per-instance state that distinguishes a long-lived aura from a
// one-shot trigger. Each time TriggerType's condition fires — and, when OncePerTurn
// is set, at most once per turn — the sim invokes Handler. The Aura survives until
// its handler calls s.DestroyAura.
type Aura struct {
	// Trigger embeds the firing-data (Source, TriggerType, TypeFilter, Handler). Aura
	// handlers receive (s, t, a); the aura-specific fields below are reached via a.
	Trigger
	// Self identifies what this Aura belongs to — a card or a token type. Surfaced in
	// per-turn summaries via CardOrTokenType.DisplayName.
	Self CardOrTokenType
	// Count is a per-aura counter the handler defines: Malefic / Runeblood read it as
	// fires remaining and decrement themselves; Runechant / Ponder read it as a stack
	// count; Blessing reads it as runes to create on the only fire. The partition
	// tiebreak (pendingFutureValue) also sums Count across all live auras as an
	// estimate of unrealised next-turn value, so even single-fire auras with no
	// counter semantics set Count: 1 to advertise their pending payoff.
	Count int
	// OncePerTurn caps the Handler at a single fire per turn regardless of how many matching
	// events occur. The sim sets FiredThisTurn the first time Handler runs each turn and
	// clears it at the next turn boundary.
	OncePerTurn bool
	// FiredThisTurn is sim-managed bookkeeping for OncePerTurn. Cards must not set it.
	FiredThisTurn bool
}
