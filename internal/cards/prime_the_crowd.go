// Prime the Crowd — Generic Action. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next attack action card you play this turn gets +N{p}. **The crowd cheers** each
// Revered hero. **The crowd boos** each Reviled hero. **Go again**" (Red N=4, Yellow N=3, Blue
// N=2.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var primeTheCrowdTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type PrimeTheCrowdRed struct{}

func (PrimeTheCrowdRed) ID() ids.CardID           { return ids.PrimeTheCrowdRed }
func (PrimeTheCrowdRed) Name() string             { return "Prime the Crowd" }
func (PrimeTheCrowdRed) Cost(*card.TurnState) int { return 2 }
func (PrimeTheCrowdRed) Pitch() int               { return 1 }
func (PrimeTheCrowdRed) Attack() int              { return 0 }
func (PrimeTheCrowdRed) Defense() int             { return 2 }
func (PrimeTheCrowdRed) Types() card.TypeSet      { return primeTheCrowdTypes }
func (PrimeTheCrowdRed) GoAgain() bool            { return true }

// not implemented: Crowd cheers / Crowd boos keywords dropped
func (PrimeTheCrowdRed) NotImplemented() {}
func (PrimeTheCrowdRed) Play(s *card.TurnState, self *card.CardState) {
	grantNextAttackActionBonus(s, 4)
	s.ApplyAndLogEffectiveAttack(self)
}

type PrimeTheCrowdYellow struct{}

func (PrimeTheCrowdYellow) ID() ids.CardID           { return ids.PrimeTheCrowdYellow }
func (PrimeTheCrowdYellow) Name() string             { return "Prime the Crowd" }
func (PrimeTheCrowdYellow) Cost(*card.TurnState) int { return 2 }
func (PrimeTheCrowdYellow) Pitch() int               { return 2 }
func (PrimeTheCrowdYellow) Attack() int              { return 0 }
func (PrimeTheCrowdYellow) Defense() int             { return 2 }
func (PrimeTheCrowdYellow) Types() card.TypeSet      { return primeTheCrowdTypes }
func (PrimeTheCrowdYellow) GoAgain() bool            { return true }

// not implemented: Crowd cheers / Crowd boos keywords dropped
func (PrimeTheCrowdYellow) NotImplemented() {}
func (PrimeTheCrowdYellow) Play(s *card.TurnState, self *card.CardState) {
	grantNextAttackActionBonus(s, 3)
	s.ApplyAndLogEffectiveAttack(self)
}

type PrimeTheCrowdBlue struct{}

func (PrimeTheCrowdBlue) ID() ids.CardID           { return ids.PrimeTheCrowdBlue }
func (PrimeTheCrowdBlue) Name() string             { return "Prime the Crowd" }
func (PrimeTheCrowdBlue) Cost(*card.TurnState) int { return 2 }
func (PrimeTheCrowdBlue) Pitch() int               { return 3 }
func (PrimeTheCrowdBlue) Attack() int              { return 0 }
func (PrimeTheCrowdBlue) Defense() int             { return 2 }
func (PrimeTheCrowdBlue) Types() card.TypeSet      { return primeTheCrowdTypes }
func (PrimeTheCrowdBlue) GoAgain() bool            { return true }

// not implemented: Crowd cheers / Crowd boos keywords dropped
func (PrimeTheCrowdBlue) NotImplemented() {}
func (PrimeTheCrowdBlue) Play(s *card.TurnState, self *card.CardState) {
	grantNextAttackActionBonus(s, 2)
	s.ApplyAndLogEffectiveAttack(self)
}
