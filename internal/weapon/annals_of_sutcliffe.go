// Annals of Sutcliffe — Runeblade Weapon - Book (2H). Activation cost {r}{r}{r}.
//
// Text: "**Once per Turn Action** - {r}{r}{r}: Draw a card. If an attack action card and a
// 'non-attack' action card were pitched this way, create a Runechant token."

package weapon

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var annalsOfSutcliffeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeBook, card.TypeTwoHand)

type AnnalsOfSutcliffe struct{}

func (AnnalsOfSutcliffe) ID() ids.CardID           { return ids.AnnalsOfSutcliffeID }
func (AnnalsOfSutcliffe) Name() string             { return "Annals of Sutcliffe" }
func (AnnalsOfSutcliffe) Cost(*card.TurnState) int { return 3 }
func (AnnalsOfSutcliffe) Pitch() int               { return 0 }
func (AnnalsOfSutcliffe) Attack() int              { return 0 }
func (AnnalsOfSutcliffe) Defense() int             { return 0 }
func (AnnalsOfSutcliffe) Types() card.TypeSet      { return annalsOfSutcliffeTypes }
func (AnnalsOfSutcliffe) GoAgain() bool            { return false }
func (AnnalsOfSutcliffe) Hands() int               { return 2 }

// not implemented: draw rider and conditional Runechant rider; activation pays 3 resources
// for zero modelled value, so the optimizer naturally avoids equipping it
func (AnnalsOfSutcliffe) NotImplemented() {}
func (AnnalsOfSutcliffe) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
