// Plunder Run — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next time an attack action card you control hits this turn, draw a card. If Plunder
// Run is played from arsenal, the next attack action card you play this turn gains +N{p}. **Go
// again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling: The on-hit draw rider isn't modelled. The +N{p} grant fires only when this copy
// was played from arsenal (self.FromArsenal); when it does, scan TurnState.CardsRemaining for
// the next attack action card and credit the bonus assuming it will be played.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var plunderRunTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type PlunderRunRed struct{}

func (PlunderRunRed) ID() card.ID                 { return card.PlunderRunRed }
func (PlunderRunRed) Name() string                { return "Plunder Run (Red)" }
func (PlunderRunRed) Cost(*card.TurnState) int                   { return 0 }
func (PlunderRunRed) Pitch() int                  { return 1 }
func (PlunderRunRed) Attack() int                 { return 0 }
func (PlunderRunRed) Defense() int                { return 2 }
func (PlunderRunRed) Types() card.TypeSet         { return plunderRunTypes }
func (PlunderRunRed) GoAgain() bool               { return true }
func (PlunderRunRed) NotSilverAgeLegal()           {}
func (PlunderRunRed) Play(s *card.TurnState, self *card.CardState) int {
	if !self.FromArsenal {
		return 0
	}
	return grantNextAttackActionBonus(s, 3)
}

type PlunderRunYellow struct{}

func (PlunderRunYellow) ID() card.ID                 { return card.PlunderRunYellow }
func (PlunderRunYellow) Name() string                { return "Plunder Run (Yellow)" }
func (PlunderRunYellow) Cost(*card.TurnState) int                   { return 0 }
func (PlunderRunYellow) Pitch() int                  { return 2 }
func (PlunderRunYellow) Attack() int                 { return 0 }
func (PlunderRunYellow) Defense() int                { return 2 }
func (PlunderRunYellow) Types() card.TypeSet         { return plunderRunTypes }
func (PlunderRunYellow) GoAgain() bool               { return true }
func (PlunderRunYellow) NotSilverAgeLegal()           {}
func (PlunderRunYellow) Play(s *card.TurnState, self *card.CardState) int {
	if !self.FromArsenal {
		return 0
	}
	return grantNextAttackActionBonus(s, 2)
}

type PlunderRunBlue struct{}

func (PlunderRunBlue) ID() card.ID                 { return card.PlunderRunBlue }
func (PlunderRunBlue) Name() string                { return "Plunder Run (Blue)" }
func (PlunderRunBlue) Cost(*card.TurnState) int                   { return 0 }
func (PlunderRunBlue) Pitch() int                  { return 3 }
func (PlunderRunBlue) Attack() int                 { return 0 }
func (PlunderRunBlue) Defense() int                { return 2 }
func (PlunderRunBlue) Types() card.TypeSet         { return plunderRunTypes }
func (PlunderRunBlue) GoAgain() bool               { return true }
func (PlunderRunBlue) NotSilverAgeLegal()           {}
func (PlunderRunBlue) Play(s *card.TurnState, self *card.CardState) int {
	if !self.FromArsenal {
		return 0
	}
	return grantNextAttackActionBonus(s, 1)
}
