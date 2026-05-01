package cards_test

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	notimpl "github.com/tim-chaplin/fab-deck-optimizer/internal/cards/notimplemented"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that the played-from-arsenal go-again rider flips self.GrantedGoAgain iff
// self.FromArsenal is true.
func TestFromArsenalGoAgain_GrantsOnArsenalCopyOnly(t *testing.T) {
	cards := []sim.Card{
		cards.FerventForerunnerRed{}, cards.FerventForerunnerYellow{}, cards.FerventForerunnerBlue{},
		notimpl.FrontlineScoutRed{}, notimpl.FrontlineScoutYellow{}, notimpl.FrontlineScoutBlue{},
		notimpl.PerformanceBonusRed{}, notimpl.PerformanceBonusYellow{}, notimpl.PerformanceBonusBlue{},
		notimpl.PromiseOfPlentyRed{}, notimpl.PromiseOfPlentyYellow{}, notimpl.PromiseOfPlentyBlue{},
		notimpl.ScourTheBattlescapeRed{}, notimpl.ScourTheBattlescapeYellow{}, notimpl.ScourTheBattlescapeBlue{},
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
