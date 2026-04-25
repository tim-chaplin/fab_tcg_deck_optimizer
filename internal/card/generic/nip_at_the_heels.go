// Nip at the Heels — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3. Defense 3.
//
// Text: "Target attack with 3 or less base {p} gets +1{p}."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var nipAtTheHeelsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type NipAtTheHeelsBlue struct{}

func (NipAtTheHeelsBlue) ID() card.ID                               { return card.NipAtTheHeelsBlue }
func (NipAtTheHeelsBlue) Name() string                              { return "Nip at the Heels (Blue)" }
func (NipAtTheHeelsBlue) Cost(*card.TurnState) int                  { return 0 }
func (NipAtTheHeelsBlue) Pitch() int                                { return 3 }
func (NipAtTheHeelsBlue) Attack() int                               { return 0 }
func (NipAtTheHeelsBlue) Defense() int                              { return 3 }
func (NipAtTheHeelsBlue) Types() card.TypeSet                       { return nipAtTheHeelsTypes }
func (NipAtTheHeelsBlue) GoAgain() bool                             { return false }
// not implemented: AR +1{p} buff and on-hit draw
func (NipAtTheHeelsBlue) NotImplemented()                           {}
func (NipAtTheHeelsBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
