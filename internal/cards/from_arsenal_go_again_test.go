package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that the played-from-arsenal go-again rider flips self.GrantedGoAgain iff
// self.FromArsenal is true.
func TestFromArsenalGoAgain_GrantsOnArsenalCopyOnly(t *testing.T) {
	cards := []sim.Card{
		FerventForerunnerRed{}, FerventForerunnerYellow{}, FerventForerunnerBlue{},
		FrontlineScoutRed{}, FrontlineScoutYellow{}, FrontlineScoutBlue{},
		PerformanceBonusRed{}, PerformanceBonusYellow{}, PerformanceBonusBlue{},
		PromiseOfPlentyRed{}, PromiseOfPlentyYellow{}, PromiseOfPlentyBlue{},
		ScourTheBattlescapeRed{}, ScourTheBattlescapeYellow{}, ScourTheBattlescapeBlue{},
	}
	for _, c := range cards {
		hand := &sim.CardState{Card: c}
		c.Play(&sim.TurnState{}, hand)
		if hand.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = true with FromArsenal=false, want false", c.Name())
		}
		arsenal := &sim.CardState{Card: c, FromArsenal: true}
		c.Play(&sim.TurnState{}, arsenal)
		if !arsenal.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false with FromArsenal=true, want true", c.Name())
		}
	}
}
