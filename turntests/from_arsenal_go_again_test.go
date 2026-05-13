package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	notimpl "github.com/tim-chaplin/fab-deck-optimizer/internal/cards/notimplemented"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that the played-from-arsenal go-again rider flips self.GrantedGoAgain iff
// self.FromArsenal is true.
func TestFromArsenalGoAgain_GrantsOnArsenalCopyOnly(t *testing.T) {
	cards := []card.Card{
		cards.FerventForerunnerRed{}, cards.FerventForerunnerYellow{}, cards.FerventForerunnerBlue{},
		notimpl.FrontlineScoutRed{}, notimpl.FrontlineScoutYellow{}, notimpl.FrontlineScoutBlue{},
		cards.PerformanceBonusRed{}, cards.PerformanceBonusYellow{}, cards.PerformanceBonusBlue{},
		notimpl.PromiseOfPlentyRed{}, notimpl.PromiseOfPlentyYellow{}, notimpl.PromiseOfPlentyBlue{},
		notimpl.ScourTheBattlescapeRed{}, notimpl.ScourTheBattlescapeYellow{}, notimpl.ScourTheBattlescapeBlue{},
	}
	for _, c := range cards {
		hand := &card.CardState{Card: c}
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
		s.ResolveChainStep(s.Logger(), hand)
		if hand.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = true with FromArsenal=false, want false", c.Name())
		}
		arsenal := &card.CardState{Card: c, FromArsenal: true}
		s2 := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
		s2.ResolveChainStep(s2.Logger(), arsenal)
		if !arsenal.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false with FromArsenal=true, want true", c.Name())
		}
	}
}
