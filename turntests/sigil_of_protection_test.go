package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// sigilOfProtectionSurvived reports whether any aura in auras carries c's display name.
func sigilOfProtectionSurvived(auras []gameengine.Aura, c card.Card) bool {
	for _, a := range auras {
		if a.CardName() == c.DisplayName() {
			return true
		}
	}
	return false
}

// Tests that each Sigil of Protection variant flips AuraCreated on Play and credits no
// direct damage.
func TestSigilOfProtection_SetsAuraCreated(t *testing.T) {
	cases := []card.Card{cards.SigilOfProtectionRed{}, cards.SigilOfProtectionYellow{}, cards.SigilOfProtectionBlue{}}
	for _, c := range cases {
		ge := gameengine.New()
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
		if !ge.AuraCreated() {
			t.Errorf("%s: AuraCreated = false, want true", c.Name())
		}
	}
}

// Tests the DamageAboutToBeTaken absorb clause directly: with the aura registered and
// incoming damage seeded, firing DamageAboutToBeTaken caps prevention at Ward 4, credits
// the prevented amount, and destroys the aura. In normal turn flow Sigil's StartOfTurn
// destroy fires before defense resolves so the absorb path is hard to reach via the public
// EvalOneTurnForTesting entry point; the file's other tests cover the start-of-turn destroy
// and AuraCreated wiring via the public path.
func TestSigilOfProtection_AbsorbsUpToFourThenSelfDestructs(t *testing.T) {
	c := cards.SigilOfProtectionRed{}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(5).Build()}
	ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c})
	if len(ge.Auras()) != 1 {
		t.Fatalf("Auras = %d, want 1 (sigil aura registered)", len(ge.Auras()))
	}

	ge.FireTriggers(triggertype.DamageAboutToBeTaken, nil)

	if got := ge.Value(); got != 4 {
		t.Errorf("Value = %d, want 4 (Ward 4 caps prevention against 5 incoming)", got)
	}
	if rem := ge.RemainingUnblockedDamage(); rem != 1 {
		t.Errorf("RemainingUnblockedDamage = %d, want 1 (5 incoming - 4 absorbed)", rem)
	}
	if len(ge.Auras()) != 0 {
		t.Errorf("Auras = %d, want 0 (sigil self-destructs after absorbing)", len(ge.Auras()))
	}
	if !graveyardContains(ge.Graveyard(), c) {
		t.Errorf("graveyard = %v, want sigil destroyed (added to graveyard)", ge.Graveyard())
	}
}

// Tests that Sigil of Protection self-destructs at the start of the owner's next action
// phase. Played on turn 1, the aura should survive turn 1 (StartOfTurn only fires on
// subsequent turns) and be destroyed at the start of turn 2 by the printed "destroy this"
// clause — landing in the graveyard.
func TestSigilOfProtection_DestroysAtStartOfNextTurn(t *testing.T) {
	sigil := cards.SigilOfProtectionRed{}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	prior := gameengine.GameStateBuilder().SetIncomingDamage(0).Build()

	// FakeBlueResource pitches for the cost-1 sigil; sigil's Ward absorb doesn't fire
	// with no incoming damage, so it must survive turn 1 on the play-this-turn path.
	hand := []card.Card{sigil, testutils.FakeBlueResource()}
	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, prior, hand)

	if !sigilOfProtectionSurvived(turn1.State.Auras(), sigil) {
		t.Fatalf("turn 1 auras = %v, graveyard = %v, BestLine = %s; want Sigil of Protection still present (StartOfTurn destroy fires on turn 2, not the turn it was played)",
			turn1.State.Auras(), turn1.State.Graveyard(), formatBestLine(turn1.BestLine))
	}
	if sigilOfProtectionSurvived(turn2.State.Auras(), sigil) {
		t.Errorf("turn 2 auras = %v, want Sigil of Protection gone (StartOfTurn destroy should fire at the start of the action phase)",
			turn2.State.Auras())
	}
	if !graveyardContains(turn2.State.Graveyard(), sigil) {
		t.Errorf("turn 2 graveyard = %v, want Sigil of Protection destroyed",
			turn2.State.Graveyard())
	}
}
