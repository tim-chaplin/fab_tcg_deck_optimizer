// Wage Gold — Generic Action - Attack. Cost 3.
// Text: "**Universal** When this attacks a hero, you may **wager** a Gold token with them."
//
// Wager modelling: opt in only when we hold a Gold at Play time. On hit, the wager is a
// no-op in our model — we don't track opponent's tokens, so "take theirs" is uncredited.
// On a fully-blocked attack, we credit the printed loss: the opponent gains our wagered
// Gold, modelled per the framework rule as ConsumeItem(Gold,1) + AddOpponentValue(1).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var wageGoldTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// wageGoldOnFullyBlocked fires when the attack is fully blocked: opponent collects our
// wagered Gold. Re-checks Gold > 0 because intermediate cards may have spent it between
// Play and finalize-active-attack — in that case there's nothing to give over.
func wageGoldOnFullyBlocked(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	if s.Gold() == 0 {
		return
	}
	s.ConsumeItem(sim.TokenTypeGold, 1)
	s.AddOpponentValue(1)
	s.LogRider(self, 0, "Lost wager — opponent gained 1 Gold")
}

func wageGoldPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	if s.Gold() > 0 {
		self.OnHit = append(self.OnHit, sim.OnHitHandler{FireBlocked: wageGoldOnFullyBlocked})
	}
}

type WageGoldRed struct{}

func (WageGoldRed) ID() ids.CardID                             { return ids.WageGoldRed }
func (WageGoldRed) Name() string                               { return "Wage Gold" }
func (WageGoldRed) Cost(*sim.TurnState) int                    { return 3 }
func (WageGoldRed) Pitch() int                                 { return 1 }
func (WageGoldRed) Attack() int                                { return 7 }
func (WageGoldRed) Defense() int                               { return 2 }
func (WageGoldRed) Types() card.TypeSet                        { return wageGoldTypes }
func (WageGoldRed) GoAgain() bool                              { return false }
func (WageGoldRed) Play(s *sim.TurnState, self *sim.CardState) { wageGoldPlay(s, self) }

type WageGoldYellow struct{}

func (WageGoldYellow) ID() ids.CardID                             { return ids.WageGoldYellow }
func (WageGoldYellow) Name() string                               { return "Wage Gold" }
func (WageGoldYellow) Cost(*sim.TurnState) int                    { return 3 }
func (WageGoldYellow) Pitch() int                                 { return 2 }
func (WageGoldYellow) Attack() int                                { return 6 }
func (WageGoldYellow) Defense() int                               { return 2 }
func (WageGoldYellow) Types() card.TypeSet                        { return wageGoldTypes }
func (WageGoldYellow) GoAgain() bool                              { return false }
func (WageGoldYellow) Play(s *sim.TurnState, self *sim.CardState) { wageGoldPlay(s, self) }

type WageGoldBlue struct{}

func (WageGoldBlue) ID() ids.CardID                             { return ids.WageGoldBlue }
func (WageGoldBlue) Name() string                               { return "Wage Gold" }
func (WageGoldBlue) Cost(*sim.TurnState) int                    { return 3 }
func (WageGoldBlue) Pitch() int                                 { return 3 }
func (WageGoldBlue) Attack() int                                { return 5 }
func (WageGoldBlue) Defense() int                               { return 2 }
func (WageGoldBlue) Types() card.TypeSet                        { return wageGoldTypes }
func (WageGoldBlue) GoAgain() bool                              { return false }
func (WageGoldBlue) Play(s *sim.TurnState, self *sim.CardState) { wageGoldPlay(s, self) }
