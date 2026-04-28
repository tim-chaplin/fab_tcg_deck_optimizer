// Vexing Malice — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 2.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Deal 2 arcane damage to target hero."
//
// The printed 2 arcane is added to combat damage (both hit the same target). Play also sets
// ArcaneDamageDealt so same-turn triggers keyed on that flag fire.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var vexingMaliceTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type VexingMaliceRed struct{}

func (VexingMaliceRed) ID() ids.CardID           { return ids.VexingMaliceRed }
func (VexingMaliceRed) Name() string             { return "Vexing Malice" }
func (VexingMaliceRed) Cost(*card.TurnState) int { return 1 }
func (VexingMaliceRed) Pitch() int               { return 1 }
func (VexingMaliceRed) Attack() int              { return 3 }
func (VexingMaliceRed) Defense() int             { return 3 }
func (VexingMaliceRed) Types() card.TypeSet      { return vexingMaliceTypes }
func (VexingMaliceRed) GoAgain() bool            { return false }
func (VexingMaliceRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.DealAndLogArcaneDamage(self, 2)
}

type VexingMaliceYellow struct{}

func (VexingMaliceYellow) ID() ids.CardID           { return ids.VexingMaliceYellow }
func (VexingMaliceYellow) Name() string             { return "Vexing Malice" }
func (VexingMaliceYellow) Cost(*card.TurnState) int { return 1 }
func (VexingMaliceYellow) Pitch() int               { return 2 }
func (VexingMaliceYellow) Attack() int              { return 2 }
func (VexingMaliceYellow) Defense() int             { return 3 }
func (VexingMaliceYellow) Types() card.TypeSet      { return vexingMaliceTypes }
func (VexingMaliceYellow) GoAgain() bool            { return false }
func (VexingMaliceYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.DealAndLogArcaneDamage(self, 2)
}

type VexingMaliceBlue struct{}

func (VexingMaliceBlue) ID() ids.CardID           { return ids.VexingMaliceBlue }
func (VexingMaliceBlue) Name() string             { return "Vexing Malice" }
func (VexingMaliceBlue) Cost(*card.TurnState) int { return 1 }
func (VexingMaliceBlue) Pitch() int               { return 3 }
func (VexingMaliceBlue) Attack() int              { return 1 }
func (VexingMaliceBlue) Defense() int             { return 3 }
func (VexingMaliceBlue) Types() card.TypeSet      { return vexingMaliceTypes }
func (VexingMaliceBlue) GoAgain() bool            { return false }
func (VexingMaliceBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.DealAndLogArcaneDamage(self, 2)
}
