// Sirens of Safe Harbor — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this is put into your graveyard from anywhere, gain 1{h}."
//
// Modelling: the card hits the graveyard after resolving as an attack, so the 1{h} gain fires
// on every Play — credited as +1 damage equivalent. Pitched copies go to the bottom of the
// deck instead of the graveyard, so they don't trigger the rider.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var sirensOfSafeHarborTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type SirensOfSafeHarborRed struct{}

func (SirensOfSafeHarborRed) ID() ids.CardID          { return ids.SirensOfSafeHarborRed }
func (SirensOfSafeHarborRed) Name() string            { return "Sirens of Safe Harbor" }
func (SirensOfSafeHarborRed) Cost(*sim.TurnState) int { return 2 }
func (SirensOfSafeHarborRed) Pitch() int              { return 1 }
func (SirensOfSafeHarborRed) Attack() int             { return 6 }
func (SirensOfSafeHarborRed) Defense() int            { return 2 }
func (SirensOfSafeHarborRed) Types() card.TypeSet     { return sirensOfSafeHarborTypes }
func (SirensOfSafeHarborRed) GoAgain() bool           { return false }
func (SirensOfSafeHarborRed) NotSilverAgeLegal()      {}
func (SirensOfSafeHarborRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
	s.LogRider(self, s.AddValue(1), "Gained 1 health (graveyard trigger)")
}

type SirensOfSafeHarborYellow struct{}

func (SirensOfSafeHarborYellow) ID() ids.CardID          { return ids.SirensOfSafeHarborYellow }
func (SirensOfSafeHarborYellow) Name() string            { return "Sirens of Safe Harbor" }
func (SirensOfSafeHarborYellow) Cost(*sim.TurnState) int { return 2 }
func (SirensOfSafeHarborYellow) Pitch() int              { return 2 }
func (SirensOfSafeHarborYellow) Attack() int             { return 5 }
func (SirensOfSafeHarborYellow) Defense() int            { return 2 }
func (SirensOfSafeHarborYellow) Types() card.TypeSet     { return sirensOfSafeHarborTypes }
func (SirensOfSafeHarborYellow) GoAgain() bool           { return false }
func (SirensOfSafeHarborYellow) NotSilverAgeLegal()      {}
func (SirensOfSafeHarborYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
	s.LogRider(self, s.AddValue(1), "Gained 1 health (graveyard trigger)")
}

type SirensOfSafeHarborBlue struct{}

func (SirensOfSafeHarborBlue) ID() ids.CardID          { return ids.SirensOfSafeHarborBlue }
func (SirensOfSafeHarborBlue) Name() string            { return "Sirens of Safe Harbor" }
func (SirensOfSafeHarborBlue) Cost(*sim.TurnState) int { return 2 }
func (SirensOfSafeHarborBlue) Pitch() int              { return 3 }
func (SirensOfSafeHarborBlue) Attack() int             { return 4 }
func (SirensOfSafeHarborBlue) Defense() int            { return 2 }
func (SirensOfSafeHarborBlue) Types() card.TypeSet     { return sirensOfSafeHarborTypes }
func (SirensOfSafeHarborBlue) GoAgain() bool           { return false }
func (SirensOfSafeHarborBlue) NotSilverAgeLegal()      {}
func (SirensOfSafeHarborBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
	s.LogRider(self, s.AddValue(1), "Gained 1 health (graveyard trigger)")
}
