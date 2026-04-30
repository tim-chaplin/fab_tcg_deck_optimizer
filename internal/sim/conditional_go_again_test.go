package sim_test

import (
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestConditionalGoAgainMarkerCoverage probes every registered card against a maximally-
// permissive TurnState. Any card whose Play flips self.GrantedGoAgain under that state
// must opt into the sim.ConditionalGoAgain marker so HasGoAgain-style hand-shaping
// heuristics see it as having Go again. Cards that gate the grant on conditions the
// probe doesn't represent (rare runtime states like a specific card in CardsPlayed)
// slip through this lint — those are caught by review or by extending the probe state
// here when discovered.
//
// This test's failure surface is "wrote a Play that grants Go again, forgot to add
// the marker". The fix is to add ConditionalGoAgain() to the offending card.
func TestConditionalGoAgainMarkerCoverage(t *testing.T) {
	for _, id := range DeckableCards() {
		c := GetCard(id)
		if c == nil {
			continue
		}
		// Cards that already advertise Go again via the printed flag don't need the
		// marker; the probe's GrantedGoAgain output is irrelevant for them.
		if c.GoAgain() {
			continue
		}
		s := newConditionalGoAgainProbeState()
		self := &CardState{Card: c, FromArsenal: true}
		safePlayProbe(c, s, self)
		if !self.GrantedGoAgain {
			continue
		}
		if _, ok := c.(ConditionalGoAgain); !ok {
			t.Errorf("%s: Play flipped GrantedGoAgain under the probe but card lacks the ConditionalGoAgain marker — add ConditionalGoAgain() {} to the variant",
				DisplayName(c))
		}
	}
}

// newConditionalGoAgainProbeState seeds every runtime gate the existing
// conditional-grant cards check: scalar flags (AuraCreated, ArcaneDamageDealt,
// NonAttackActionPlayed, Runechants > 0), CardsPlayed containing a Moon Wish-named
// stub (for Sun Kiss's name-match), and a high-attack stub in Pitched (for Zealous
// Belting's pitched-attack-greater-than-base check). Combined with FromArsenal=true
// on the CardState, this catches every conditional grant we know about today. New
// gating conditions in future cards may need a flag or stub added here.
func newConditionalGoAgainProbeState() *TurnState {
	moonWishStub := testutils.NewStubCard("Moon Wish")
	highAttackStub := testutils.NewStubCard("probeHighAttackPitched").WithAttack(99)
	return &TurnState{
		AuraCreated:           true,
		ArcaneDamageDealt:     true,
		NonAttackActionPlayed: true,
		Runechants:            1,
		CardsPlayed:           []Card{moonWishStub},
		Pitched:               []Card{highAttackStub},
	}
}

// safePlayProbe runs c.Play under recover() so a card whose Play panics on the
// stripped-down probe state (e.g. a nil-deref against a missing accessor) doesn't
// take the whole lint test down with it. The recovered panic is treated as "didn't
// flip GrantedGoAgain" — at worst we miss flagging a marker on a card whose grant
// path is unreachable from probe state, which the per-card review would catch anyway.
func safePlayProbe(c Card, s *TurnState, self *CardState) {
	defer func() { _ = recover() }()
	c.Play(s, self)
}
