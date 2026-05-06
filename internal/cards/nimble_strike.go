// Nimble Strike — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Nimble Strike, you may banish a card named Nimblism from
// your graveyard. If you do, Nimble Strike gain +1{p} and **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var nimbleStrikeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func nimbleStrikePlay(s *sim.TurnState, self *sim.CardState) {
	if _, ok := s.BanishFromGraveyard(isNimblism); ok {
		self.BonusAttack++
		self.GrantedGoAgain = true
		s.LogRider(self, 1, "Banished a Nimblism, +1{p} and go again")
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

func isNimblism(c sim.Card) bool { return c.Name() == "Nimblism" }

type NimbleStrikeRed struct{}

func (NimbleStrikeRed) ID() ids.CardID          { return ids.NimbleStrikeRed }
func (NimbleStrikeRed) Name() string            { return "Nimble Strike" }
func (NimbleStrikeRed) Cost(*sim.TurnState) int { return 1 }
func (NimbleStrikeRed) Pitch() int              { return 1 }
func (NimbleStrikeRed) Attack() int             { return 4 }
func (NimbleStrikeRed) Defense() int            { return 2 }
func (NimbleStrikeRed) Types() card.TypeSet     { return nimbleStrikeTypes }
func (NimbleStrikeRed) GoAgain() bool           { return false }
func (NimbleStrikeRed) Play(s *sim.TurnState, self *sim.CardState) {
	nimbleStrikePlay(s, self)
}

type NimbleStrikeYellow struct{}

func (NimbleStrikeYellow) ID() ids.CardID          { return ids.NimbleStrikeYellow }
func (NimbleStrikeYellow) Name() string            { return "Nimble Strike" }
func (NimbleStrikeYellow) Cost(*sim.TurnState) int { return 1 }
func (NimbleStrikeYellow) Pitch() int              { return 2 }
func (NimbleStrikeYellow) Attack() int             { return 3 }
func (NimbleStrikeYellow) Defense() int            { return 2 }
func (NimbleStrikeYellow) Types() card.TypeSet     { return nimbleStrikeTypes }
func (NimbleStrikeYellow) GoAgain() bool           { return false }
func (NimbleStrikeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	nimbleStrikePlay(s, self)
}

type NimbleStrikeBlue struct{}

func (NimbleStrikeBlue) ID() ids.CardID          { return ids.NimbleStrikeBlue }
func (NimbleStrikeBlue) Name() string            { return "Nimble Strike" }
func (NimbleStrikeBlue) Cost(*sim.TurnState) int { return 1 }
func (NimbleStrikeBlue) Pitch() int              { return 3 }
func (NimbleStrikeBlue) Attack() int             { return 2 }
func (NimbleStrikeBlue) Defense() int            { return 2 }
func (NimbleStrikeBlue) Types() card.TypeSet     { return nimbleStrikeTypes }
func (NimbleStrikeBlue) GoAgain() bool           { return false }
func (NimbleStrikeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	nimbleStrikePlay(s, self)
}
