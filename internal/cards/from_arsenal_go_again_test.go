package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestFromArsenalGoAgain_GrantsOnArsenalCopyOnly checks the shared rider on cards whose
// printed "If <name> is played from arsenal, it gains **go again**" clause flips
// self.GrantedGoAgain iff self.FromArsenal is true. Hand-played copies must not gain go again
// (printed GoAgain() is already false; the grant is the only path to a chain step that
// restocks ActionPoints), so a missed wiring would let every copy go again and over-credit
// arsenal-gated decks.
func TestFromArsenalGoAgain_GrantsOnArsenalCopyOnly(t *testing.T) {
	cards := []sim.Card{
		FerventForerunnerRed{}, FerventForerunnerYellow{}, FerventForerunnerBlue{},
		FrontlineScoutRed{}, FrontlineScoutYellow{}, FrontlineScoutBlue{},
		PerformanceBonusRed{}, PerformanceBonusYellow{}, PerformanceBonusBlue{},
		PromiseOfPlentyRed{}, PromiseOfPlentyYellow{}, PromiseOfPlentyBlue{},
		ScourTheBattlescapeRed{}, ScourTheBattlescapeYellow{}, ScourTheBattlescapeBlue{},
	}
	for _, c := range cards {
		// Hand-played copy: no grant — EffectiveGoAgain must stay false.
		hand := &sim.CardState{Card: c}
		c.Play(&sim.TurnState{}, hand)
		if hand.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = true with FromArsenal=false, want false", c.Name())
		}
		if hand.EffectiveGoAgain() {
			t.Errorf("%s: EffectiveGoAgain = true with FromArsenal=false, want false", c.Name())
		}
		// Arsenal-played copy: rider fires — GrantedGoAgain must flip and EffectiveGoAgain
		// must report true so the chain dispatcher restocks an action point.
		arsenal := &sim.CardState{Card: c, FromArsenal: true}
		c.Play(&sim.TurnState{}, arsenal)
		if !arsenal.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false with FromArsenal=true, want true", c.Name())
		}
		if !arsenal.EffectiveGoAgain() {
			t.Errorf("%s: EffectiveGoAgain = false with FromArsenal=true, want true", c.Name())
		}
	}
}
