package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	notimpl "github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards/notimplemented"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Tests that the played-from-arsenal go-again rider flips self.GrantedGoAgain iff
// self.FromArsenal is true.
func TestFromArsenalGoAgain_GrantsOnArsenalCopyOnly(t *testing.T) {
	cards := []card.Card{
		cards.FerventForerunnerRed{}, cards.FerventForerunnerYellow{}, cards.FerventForerunnerBlue{},
		cards.FrontlineScoutRed{}, cards.FrontlineScoutYellow{}, cards.FrontlineScoutBlue{},
		cards.PerformanceBonusRed{}, cards.PerformanceBonusYellow{}, cards.PerformanceBonusBlue{},
		notimpl.PromiseOfPlentyRed{}, notimpl.PromiseOfPlentyYellow{}, notimpl.PromiseOfPlentyBlue{},
		cards.ScourTheBattlescapeRed{}, cards.ScourTheBattlescapeYellow{}, cards.ScourTheBattlescapeBlue{},
	}
	for _, c := range cards {
		hand := &card.CardState{Card: c}
		ge := gameengine.New()
		ge.ResolveAttackStep(ge.Logger(), hand)
		if hand.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = true with FromArsenal=false, want false", c.Name())
		}
		arsenal := &card.CardState{Card: c, FromArsenal: true}
		s2 := gameengine.New()
		s2.ResolveAttackStep(s2.Logger(), arsenal)
		if !arsenal.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false with FromArsenal=true, want true", c.Name())
		}
	}
}
