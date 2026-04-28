// Warmonger's Recital — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next attack action card you play this turn gains +N{p} and "When this hits, put it on
// the bottom of its owner's deck." **Go again**" (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var warmongersRecitalTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type WarmongersRecitalRed struct{}

func (WarmongersRecitalRed) ID() ids.CardID           { return ids.WarmongersRecitalRed }
func (WarmongersRecitalRed) Name() string             { return "Warmonger's Recital" }
func (WarmongersRecitalRed) Cost(*card.TurnState) int { return 1 }
func (WarmongersRecitalRed) Pitch() int               { return 1 }
func (WarmongersRecitalRed) Attack() int              { return 0 }
func (WarmongersRecitalRed) Defense() int             { return 2 }
func (WarmongersRecitalRed) Types() card.TypeSet      { return warmongersRecitalTypes }
func (WarmongersRecitalRed) GoAgain() bool            { return true }

// not implemented: bottom-of-deck rider on next-attack-action target
func (WarmongersRecitalRed) NotImplemented() {}
func (WarmongersRecitalRed) Play(s *card.TurnState, self *card.CardState) {
	grantNextAttackActionBonus(s, 3)
	s.ApplyAndLogEffectiveAttack(self)
}

type WarmongersRecitalYellow struct{}

func (WarmongersRecitalYellow) ID() ids.CardID           { return ids.WarmongersRecitalYellow }
func (WarmongersRecitalYellow) Name() string             { return "Warmonger's Recital" }
func (WarmongersRecitalYellow) Cost(*card.TurnState) int { return 1 }
func (WarmongersRecitalYellow) Pitch() int               { return 2 }
func (WarmongersRecitalYellow) Attack() int              { return 0 }
func (WarmongersRecitalYellow) Defense() int             { return 2 }
func (WarmongersRecitalYellow) Types() card.TypeSet      { return warmongersRecitalTypes }
func (WarmongersRecitalYellow) GoAgain() bool            { return true }

// not implemented: bottom-of-deck rider on next-attack-action target
func (WarmongersRecitalYellow) NotImplemented() {}
func (WarmongersRecitalYellow) Play(s *card.TurnState, self *card.CardState) {
	grantNextAttackActionBonus(s, 2)
	s.ApplyAndLogEffectiveAttack(self)
}

type WarmongersRecitalBlue struct{}

func (WarmongersRecitalBlue) ID() ids.CardID           { return ids.WarmongersRecitalBlue }
func (WarmongersRecitalBlue) Name() string             { return "Warmonger's Recital" }
func (WarmongersRecitalBlue) Cost(*card.TurnState) int { return 1 }
func (WarmongersRecitalBlue) Pitch() int               { return 3 }
func (WarmongersRecitalBlue) Attack() int              { return 0 }
func (WarmongersRecitalBlue) Defense() int             { return 2 }
func (WarmongersRecitalBlue) Types() card.TypeSet      { return warmongersRecitalTypes }
func (WarmongersRecitalBlue) GoAgain() bool            { return true }

// not implemented: bottom-of-deck rider on next-attack-action target
func (WarmongersRecitalBlue) NotImplemented() {}
func (WarmongersRecitalBlue) Play(s *card.TurnState, self *card.CardState) {
	grantNextAttackActionBonus(s, 1)
	s.ApplyAndLogEffectiveAttack(self)
}
