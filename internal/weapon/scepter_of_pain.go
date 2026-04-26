// Scepter of Pain — Runeblade Weapon - Scepter (1H). Cost 2, Arcane 1.
// Text: "Once per Turn Action - {r}{r}: Deal 1 arcane damage to any opposing target. Create a
// Runechant token for each damage dealt this way."

package weapon

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var scepterOfPainTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeScepter, card.TypeOneHand)

type ScepterOfPain struct{}

func (ScepterOfPain) ID() card.ID              { return card.ScepterOfPainID }
func (ScepterOfPain) Name() string             { return "Scepter of Pain" }
func (ScepterOfPain) Cost(*card.TurnState) int { return 2 }
func (ScepterOfPain) Pitch() int               { return 0 }
func (ScepterOfPain) Attack() int              { return 1 }
func (ScepterOfPain) Defense() int             { return 0 }
func (ScepterOfPain) Types() card.TypeSet      { return scepterOfPainTypes }
func (ScepterOfPain) GoAgain() bool            { return false }
func (ScepterOfPain) Hands() int               { return 1 }
func (c ScepterOfPain) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.CreateAndLogRunechantsOnPlay(self, 1)
}
