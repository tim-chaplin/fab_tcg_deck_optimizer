package sim

// Aura tokens are auras with no originating card: when destroyed they just disappear
// (no graveyard append). Each token type has one fixed handler defined here, since
// behaviour is independent of the card that created the token.
//
// Invariant: at most one Aura per token type per TurnState — helpers bump Count on the
// existing entry rather than appending duplicates. Keeps cache keys and the trigger-fire
// loop compact.

// TokenType identifies an aura token kind. TokenTypeNone is the zero value used by card
// auras (which set Aura.Self.Card instead).
type TokenType int

const (
	// TokenTypeNone marks a non-token aura (Aura.Self.Card is set instead).
	TokenTypeNone TokenType = iota
	// TokenTypeRunechant is the runechant aura token. Consumed by the next attack or
	// weapon swing the controller resolves (see runechantAuraHandler).
	TokenTypeRunechant
	// TokenTypePonder is the ponder aura token. Destroys itself at end of the turn it
	// was created, drawing a card before the arsenal-promotion step (see
	// ponderAuraHandler).
	TokenTypePonder
)

// tokenDisplayName returns the printed name shown in logs and "(from previous turn)"
// summaries for the given token type. Mirrors DisplayName(Card) for card auras.
func tokenDisplayName(t TokenType) string {
	switch t {
	case TokenTypeRunechant:
		return "Runechant"
	case TokenTypePonder:
		return "Ponder"
	}
	return ""
}

// runechantAuraHandler is the TriggerAttack handler shared by every Runechant aura.
// Fires before each attack / weapon swing resolves: flips ArcaneDamageDealt when
// t.Count clears the LikelyDamageHits window and destroys the aura. Damage was credited
// at creation time in CreateRunechants — this handler is pure state cleanup.
func runechantAuraHandler(s *TurnState, t *Aura) {
	if LikelyDamageHits(t.Count, false) {
		s.ArcaneDamageDealt = true
	}
	s.DestroyAura(t, false)
}

// NewRunechantAura returns a runechant token aura at count n. Production code calls
// s.CreateRunechants instead — it bumps an existing aura and credits +n damage. This
// factory is for tests that need to seed a runechant aura without the damage credit.
func NewRunechantAura(n int) Aura {
	return Aura{
		Self:        CardOrTokenType{TokenType: TokenTypeRunechant},
		TriggerType: TriggerAttack,
		Count:       n,
		Handler:     runechantAuraHandler,
	}
}

// ponderAuraHandler is the TriggerEndOfTurn handler shared by every Ponder aura. For
// each token in play it pops the deck top into the held pile (s.Hand), letting the
// post-hoc arsenal-promotion step fill an otherwise-empty arsenal slot. Pops past
// deck-end are silently skipped — empty deck just means no draw. Reading the deck top
// flips s.cacheable (PopDeckTop's contract).
func ponderAuraHandler(s *TurnState, t *Aura) {
	for i := 0; i < t.Count; i++ {
		c, ok := s.PopDeckTop()
		if !ok {
			break
		}
		s.Hand = append(s.Hand, c)
	}
	s.DestroyAura(t, false)
}

// NewPonderAura returns a ponder token aura at count n. Production code calls
// s.CreatePonder instead; this factory is for tests that need to seed the aura directly.
func NewPonderAura(n int) Aura {
	return Aura{
		Self:        CardOrTokenType{TokenType: TokenTypePonder},
		TriggerType: TriggerEndOfTurn,
		Count:       n,
		Handler:     ponderAuraHandler,
	}
}

// tokenCountIn scans an aura slice for a token aura of the given type and returns its
// count. Single read site for the at-most-one-aura-per-token-type invariant.
func tokenCountIn(auras []Aura, t TokenType) int {
	for i := range auras {
		if auras[i].Self.TokenType == t {
			return auras[i].Count
		}
	}
	return 0
}
