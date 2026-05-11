// Scepter of Pain — Runeblade Weapon - Scepter (1H). Cost 2, Arcane 1.
// Text: "Once per Turn Action - {r}{r}: Deal 1 arcane damage to any opposing target. Create a
// Runechant token for each damage dealt this way."

package weapons

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var scepterOfPainTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeScepter, card.TypeOneHand)

type ScepterOfPain struct{}

func (ScepterOfPain) ID() ids.WeaponID    { return ids.ScepterOfPainID }
func (ScepterOfPain) Name() string        { return "Scepter of Pain" }
func (ScepterOfPain) Types() card.TypeSet { return scepterOfPainTypes }
func (ScepterOfPain) Hands() int          { return 1 }
func (ScepterOfPain) Ability() sim.Card   { return scepterOfPainAbility }

// Cached at package init — see nebula_blade.go for the alloc-free rationale.
var scepterOfPainAbility sim.Card = ScepterOfPainAbility{}

var scepterOfPainAbilityTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeScepter, card.TypeOneHand, card.TypeAttack)

type ScepterOfPainAbility struct{}

func (ScepterOfPainAbility) ID() ids.CardID          { return ids.ScepterOfPainAbilityID }
func (ScepterOfPainAbility) Name() string            { return "Scepter of Pain" }
func (ScepterOfPainAbility) DisplayName() string     { return "Scepter of Pain" }
func (ScepterOfPainAbility) Cost(sim.GameEngine) int { return 2 }
func (ScepterOfPainAbility) Pitch() int              { return 0 }
func (ScepterOfPainAbility) Attack() int             { return 1 }
func (ScepterOfPainAbility) Defense() int            { return 0 }
func (ScepterOfPainAbility) Types() card.TypeSet     { return scepterOfPainAbilityTypes }
func (ScepterOfPainAbility) GoAgain() bool           { return false }
func (ScepterOfPainAbility) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}
