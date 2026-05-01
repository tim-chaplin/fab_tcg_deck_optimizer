// Spellblade Assault — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When you attack with Spellblade Assault, create 2 Runechant tokens."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var spellbladeAssaultTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type SpellbladeAssaultRed struct{}

func (SpellbladeAssaultRed) ID() ids.CardID          { return ids.SpellbladeAssaultRed }
func (SpellbladeAssaultRed) Name() string            { return "Spellblade Assault" }
func (SpellbladeAssaultRed) Cost(*sim.TurnState) int { return 2 }
func (SpellbladeAssaultRed) Pitch() int              { return 1 }
func (SpellbladeAssaultRed) Attack() int             { return 4 }
func (SpellbladeAssaultRed) Defense() int            { return 3 }
func (SpellbladeAssaultRed) Types() card.TypeSet     { return spellbladeAssaultTypes }
func (SpellbladeAssaultRed) GoAgain() bool           { return false }
func (SpellbladeAssaultRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	s.AddValue(s.CreateRunechants(2))
	s.LogRider(self, 2, "Created 2 runechants")
}

type SpellbladeAssaultYellow struct{}

func (SpellbladeAssaultYellow) ID() ids.CardID          { return ids.SpellbladeAssaultYellow }
func (SpellbladeAssaultYellow) Name() string            { return "Spellblade Assault" }
func (SpellbladeAssaultYellow) Cost(*sim.TurnState) int { return 2 }
func (SpellbladeAssaultYellow) Pitch() int              { return 2 }
func (SpellbladeAssaultYellow) Attack() int             { return 3 }
func (SpellbladeAssaultYellow) Defense() int            { return 3 }
func (SpellbladeAssaultYellow) Types() card.TypeSet     { return spellbladeAssaultTypes }
func (SpellbladeAssaultYellow) GoAgain() bool           { return false }
func (SpellbladeAssaultYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	s.AddValue(s.CreateRunechants(2))
	s.LogRider(self, 2, "Created 2 runechants")
}

type SpellbladeAssaultBlue struct{}

func (SpellbladeAssaultBlue) ID() ids.CardID          { return ids.SpellbladeAssaultBlue }
func (SpellbladeAssaultBlue) Name() string            { return "Spellblade Assault" }
func (SpellbladeAssaultBlue) Cost(*sim.TurnState) int { return 2 }
func (SpellbladeAssaultBlue) Pitch() int              { return 3 }
func (SpellbladeAssaultBlue) Attack() int             { return 2 }
func (SpellbladeAssaultBlue) Defense() int            { return 3 }
func (SpellbladeAssaultBlue) Types() card.TypeSet     { return spellbladeAssaultTypes }
func (SpellbladeAssaultBlue) GoAgain() bool           { return false }
func (SpellbladeAssaultBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	s.AddValue(s.CreateRunechants(2))
	s.LogRider(self, 2, "Created 2 runechants")
}
