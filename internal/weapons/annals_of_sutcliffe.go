// Annals of Sutcliffe — Runeblade Weapon - Book (2H). Activation cost {r}{r}{r}.
//
// Text: "**Once per Turn Action** - {r}{r}{r}: Draw a card. If an attack action card and a
// 'non-attack' action card were pitched this way, create a Runechant token."

package weapons

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

var annalsOfSutcliffeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeBook, card.TypeTwoHand)

type AnnalsOfSutcliffe struct{}

func (AnnalsOfSutcliffe) ID() ids.WeaponID    { return ids.AnnalsOfSutcliffeID }
func (AnnalsOfSutcliffe) Name() string        { return "Annals of Sutcliffe" }
func (AnnalsOfSutcliffe) Types() card.TypeSet { return annalsOfSutcliffeTypes }
func (AnnalsOfSutcliffe) Hands() int          { return 2 }
func (AnnalsOfSutcliffe) Ability() card.Card  { return annalsOfSutcliffeAbility }

// Cached at package init — see nebula_blade.go for the alloc-free rationale.
var annalsOfSutcliffeAbility card.Card = AnnalsOfSutcliffeAbility{}

// not implemented: draw rider and conditional Runechant rider; activation pays 3 resources
// for zero modelled value, so the optimizer naturally avoids equipping it
func (AnnalsOfSutcliffe) NotImplemented() {}

var annalsOfSutcliffeAbilityTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeBook, card.TypeTwoHand, card.TypeAttack)

type AnnalsOfSutcliffeAbility struct{}

func (AnnalsOfSutcliffeAbility) ID() ids.CardID           { return ids.AnnalsOfSutcliffeAbilityID }
func (AnnalsOfSutcliffeAbility) Name() string             { return "Annals of Sutcliffe" }
func (AnnalsOfSutcliffeAbility) DisplayName() string      { return "Annals of Sutcliffe" }
func (AnnalsOfSutcliffeAbility) Cost(card.GameEngine) int { return 3 }
func (AnnalsOfSutcliffeAbility) Pitch() int               { return 0 }
func (AnnalsOfSutcliffeAbility) Attack() int              { return 0 }
func (AnnalsOfSutcliffeAbility) Defense() int             { return 0 }
func (AnnalsOfSutcliffeAbility) Types() card.TypeSet      { return annalsOfSutcliffeAbilityTypes }
func (AnnalsOfSutcliffeAbility) GoAgain() bool            { return false }
func (AnnalsOfSutcliffeAbility) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
}
