// Yinti Yanti — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While Yinti Yanti is attacking and you control an aura, it has +1{p}. While Yinti Yanti is
// defending and you control an aura, it has +1{d}."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var yintiYantiTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// yintiYantiPlay adds +1 when any aura is in play: either created this turn or played earlier.
func yintiYantiPlay(base int, s *card.TurnState) int {
	if s != nil && s.HasAuraInPlay() {
		return base + 1
	}
	return base
}

type YintiYantiRed struct{}

func (YintiYantiRed) ID() card.ID                 { return card.YintiYantiRed }
func (YintiYantiRed) Name() string                { return "Yinti Yanti" }
func (YintiYantiRed) Cost(*card.TurnState) int                   { return 0 }
func (YintiYantiRed) Pitch() int                  { return 1 }
func (YintiYantiRed) Attack() int                 { return 3 }
func (YintiYantiRed) Defense() int                { return 2 }
func (YintiYantiRed) Types() card.TypeSet         { return yintiYantiTypes }
func (YintiYantiRed) GoAgain() bool               { return false }
// not implemented: defending-side +1{d} buff (defence consumed before Play); aura-attack
// +1{p} is modelled
func (YintiYantiRed) NotImplemented()             {}
func (c YintiYantiRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, yintiYantiPlay(c.Attack(), s)-self.Card.Attack())
}
type YintiYantiYellow struct{}

func (YintiYantiYellow) ID() card.ID                 { return card.YintiYantiYellow }
func (YintiYantiYellow) Name() string                { return "Yinti Yanti" }
func (YintiYantiYellow) Cost(*card.TurnState) int                   { return 0 }
func (YintiYantiYellow) Pitch() int                  { return 2 }
func (YintiYantiYellow) Attack() int                 { return 2 }
func (YintiYantiYellow) Defense() int                { return 2 }
func (YintiYantiYellow) Types() card.TypeSet         { return yintiYantiTypes }
func (YintiYantiYellow) GoAgain() bool               { return false }
// not implemented: defending-side +1{d} buff (defence consumed before Play); aura-attack
// +1{p} is modelled
func (YintiYantiYellow) NotImplemented()             {}
func (c YintiYantiYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, yintiYantiPlay(c.Attack(), s)-self.Card.Attack())
}
type YintiYantiBlue struct{}

func (YintiYantiBlue) ID() card.ID                 { return card.YintiYantiBlue }
func (YintiYantiBlue) Name() string                { return "Yinti Yanti" }
func (YintiYantiBlue) Cost(*card.TurnState) int                   { return 0 }
func (YintiYantiBlue) Pitch() int                  { return 3 }
func (YintiYantiBlue) Attack() int                 { return 1 }
func (YintiYantiBlue) Defense() int                { return 2 }
func (YintiYantiBlue) Types() card.TypeSet         { return yintiYantiTypes }
func (YintiYantiBlue) GoAgain() bool               { return false }
// not implemented: defending-side +1{d} buff (defence consumed before Play); aura-attack
// +1{p} is modelled
func (YintiYantiBlue) NotImplemented()             {}
func (c YintiYantiBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, yintiYantiPlay(c.Attack(), s)-self.Card.Attack())
}