// Scepter of Pain — Runeblade Weapon - Scepter (1H). Cost 2, Arcane 1.
// Text: "Once per Turn Action - {r}{r}: Deal 1 arcane damage to any opposing target. Create a
// Runechant token for each damage dealt this way."

package weapons

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

var scepterOfPainTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeScepter, card.TypeOneHand)

// ScepterOfPain is the platonic weapon card — the equipped permanent. Its activated ability
// is what the attack-turn runner enqueues each turn.
type ScepterOfPain struct{}

func (ScepterOfPain) ID() ids.CardID                     { return ids.ScepterOfPainID }
func (ScepterOfPain) Name() string                       { return "Scepter of Pain" }
func (ScepterOfPain) DisplayName() string                { return "Scepter of Pain" }
func (ScepterOfPain) Cost() int                          { return 0 }
func (ScepterOfPain) Pitch() int                         { return 0 }
func (ScepterOfPain) Attack() int                        { return 0 }
func (ScepterOfPain) Defense() int                       { return 0 }
func (ScepterOfPain) Types(card.GameEngine) card.TypeSet { return scepterOfPainTypes }
func (ScepterOfPain) GoAgain(card.GameEngine) bool       { return false }
func (ScepterOfPain) Hands() int                         { return 1 }
func (ScepterOfPain) Ability() card.Card                 { return scepterOfPainAbility }

func (ScepterOfPain) Play(ge card.GameEngine, _ card.Logger, self *card.CardState) {
	ge.CreateWeapon(self.Card, 0, nil, false, nil)
}

var scepterOfPainAbility card.Card = ScepterOfPainAbility{}

var scepterOfPainAbilityTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeScepter, card.TypeOneHand, card.TypeAttack)

type ScepterOfPainAbility struct{}

func (ScepterOfPainAbility) ID() ids.CardID                     { return ids.ScepterOfPainAbilityID }
func (ScepterOfPainAbility) Name() string                       { return "Scepter of Pain" }
func (ScepterOfPainAbility) DisplayName() string                { return "Scepter of Pain" }
func (ScepterOfPainAbility) Cost() int                          { return 2 }
func (ScepterOfPainAbility) Pitch() int                         { return 0 }
func (ScepterOfPainAbility) Attack() int                        { return 1 }
func (ScepterOfPainAbility) Defense() int                       { return 0 }
func (ScepterOfPainAbility) Types(card.GameEngine) card.TypeSet { return scepterOfPainAbilityTypes }
func (ScepterOfPainAbility) GoAgain(card.GameEngine) bool       { return false }
func (ScepterOfPainAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}
