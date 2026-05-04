// Yinti Yanti — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While Yinti Yanti is attacking and you control an aura, it has +1{p}. While Yinti
// Yanti is defending and you control an aura, it has +1{d}."
//
// Both bonuses gate on len(s.Auras) > 0 — any aura in play (sigil, incantation, runechant,
// ponder, …) qualifies. Per the FaB-correct turn order (defense → action), DR-created
// auras on the same turn are visible to Yinti Yanti during action; per-DR aura
// accumulation in defendersDamage means plain-block Yinti Yanti also sees auras from
// earlier defenders.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var yintiYantiTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// yintiYantiBonus returns +1 when any aura is in play, else 0. Used by both Play (attack
// side) and Block (defense side) so the gate stays in one place.
func yintiYantiBonus(s *sim.TurnState) int {
	if s != nil && len(s.Auras) > 0 {
		return 1
	}
	return 0
}

func yintiYantiPlay(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += yintiYantiBonus(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

func yintiYantiBlock(s *sim.TurnState, self *sim.CardState) {
	self.BonusDefense += yintiYantiBonus(s)
}

type YintiYantiRed struct{}

func (YintiYantiRed) ID() ids.CardID          { return ids.YintiYantiRed }
func (YintiYantiRed) Name() string            { return "Yinti Yanti" }
func (YintiYantiRed) Cost(*sim.TurnState) int { return 0 }
func (YintiYantiRed) Pitch() int              { return 1 }
func (YintiYantiRed) Attack() int             { return 3 }
func (YintiYantiRed) Defense() int            { return 2 }
func (YintiYantiRed) Types() card.TypeSet     { return yintiYantiTypes }
func (YintiYantiRed) GoAgain() bool           { return false }
func (YintiYantiRed) Play(s *sim.TurnState, self *sim.CardState) {
	yintiYantiPlay(s, self)
}
func (YintiYantiRed) Block(s *sim.TurnState, self *sim.CardState) {
	yintiYantiBlock(s, self)
}

type YintiYantiYellow struct{}

func (YintiYantiYellow) ID() ids.CardID          { return ids.YintiYantiYellow }
func (YintiYantiYellow) Name() string            { return "Yinti Yanti" }
func (YintiYantiYellow) Cost(*sim.TurnState) int { return 0 }
func (YintiYantiYellow) Pitch() int              { return 2 }
func (YintiYantiYellow) Attack() int             { return 2 }
func (YintiYantiYellow) Defense() int            { return 2 }
func (YintiYantiYellow) Types() card.TypeSet     { return yintiYantiTypes }
func (YintiYantiYellow) GoAgain() bool           { return false }
func (YintiYantiYellow) Play(s *sim.TurnState, self *sim.CardState) {
	yintiYantiPlay(s, self)
}
func (YintiYantiYellow) Block(s *sim.TurnState, self *sim.CardState) {
	yintiYantiBlock(s, self)
}

type YintiYantiBlue struct{}

func (YintiYantiBlue) ID() ids.CardID          { return ids.YintiYantiBlue }
func (YintiYantiBlue) Name() string            { return "Yinti Yanti" }
func (YintiYantiBlue) Cost(*sim.TurnState) int { return 0 }
func (YintiYantiBlue) Pitch() int              { return 3 }
func (YintiYantiBlue) Attack() int             { return 1 }
func (YintiYantiBlue) Defense() int            { return 2 }
func (YintiYantiBlue) Types() card.TypeSet     { return yintiYantiTypes }
func (YintiYantiBlue) GoAgain() bool           { return false }
func (YintiYantiBlue) Play(s *sim.TurnState, self *sim.CardState) {
	yintiYantiPlay(s, self)
}
func (YintiYantiBlue) Block(s *sim.TurnState, self *sim.CardState) {
	yintiYantiBlock(s, self)
}
