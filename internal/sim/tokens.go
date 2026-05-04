package sim

// Aura tokens — auras that don't belong to a card and just disappear when destroyed
// (no graveyard, no card to attribute). Each token type has a fixed handler defined here
// rather than per-card, since a Runechant on the table behaves the same regardless of who
// created it.
//
// One Aura per token type per TurnState: helpers like CreateRunechants find the existing
// entry and bump its Count instead of appending duplicates. The single-aura invariant
// keeps cache keys (and the trigger-fire loop) compact — n runechants is "Count = n", not
// n separate entries.

// TokenType identifies an aura token kind. TokenTypeNone is the zero value used by card
// auras (which set Aura.Self.Card instead).
type TokenType int

const (
	// TokenTypeNone marks a non-token aura (Aura.Self.Card is set instead).
	TokenTypeNone TokenType = iota
	// TokenTypeRunechant is the runechant aura token. Created by Viserai's hero ability,
	// Sigil of Deadwood, Malefic Incantation, Read the Runes, …; consumed by the next
	// attack or weapon swing the controller resolves (see runechantAuraHandler).
	TokenTypeRunechant
)

// tokenDisplayName returns the printed name shown in logs and "(from previous turn)"
// summaries for the given token type. Mirrors DisplayName(Card) for card auras.
func tokenDisplayName(t TokenType) string {
	switch t {
	case TokenTypeRunechant:
		return "Runechant"
	}
	return ""
}

// runechantAuraHandler is the TriggerAttack handler shared by every Runechant aura. The
// chain runner fires it before each attack / weapon swing resolves; the handler flips
// ArcaneDamageDealt when t.Count hits the LikelyDamageHits window and destroys the aura
// (token, no graveyard). Damage was credited at creation time inside CreateRunechants —
// this handler does not re-credit, it's pure state cleanup.
func runechantAuraHandler(s *TurnState, t *Aura) {
	if LikelyDamageHits(t.Count, false) {
		s.ArcaneDamageDealt = true
	}
	s.DestroyAura(t, false)
}

// NewRunechantAura returns a runechant token aura at count n. Production code calls
// s.CreateRunechants which both creates the aura (or bumps an existing one) and credits
// +n damage at creation time; tests that want to seed a runechant aura without crediting
// damage build the aura via this factory and append directly to TurnState.Auras.
func NewRunechantAura(n int) Aura {
	return Aura{
		Self:        CardOrTokenType{TokenType: TokenTypeRunechant},
		TriggerType: TriggerAttack,
		Count:       n,
		Handler:     runechantAuraHandler,
	}
}

// priorRunechantCount returns the runechant count in a priorAuras carryover slice. The
// chain runner uses this to seed bestLeftoverRunechants — the tiebreaker for "no chain
// played, fall back to whatever was already there".
func priorRunechantCount(priorAuras []Aura) int {
	for i := range priorAuras {
		if priorAuras[i].Self.TokenType == TokenTypeRunechant {
			return priorAuras[i].Count
		}
	}
	return 0
}
