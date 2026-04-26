// Emissary of Moon — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 2. Only printed in
// Red.
//
// Text: "When this attacks, you may put a card from your hand on the bottom of your deck. If you
// do, draw a card."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var emissaryOfMoonTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type EmissaryOfMoonRed struct{}

func (EmissaryOfMoonRed) ID() card.ID                 { return card.EmissaryOfMoonRed }
func (EmissaryOfMoonRed) Name() string                { return "Emissary of Moon" }
func (EmissaryOfMoonRed) Cost(*card.TurnState) int                   { return 0 }
func (EmissaryOfMoonRed) Pitch() int                  { return 1 }
func (EmissaryOfMoonRed) Attack() int                 { return 4 }
func (EmissaryOfMoonRed) Defense() int                { return 2 }
func (EmissaryOfMoonRed) Types() card.TypeSet         { return emissaryOfMoonTypes }
func (EmissaryOfMoonRed) GoAgain() bool               { return false }
// not implemented: hand-cycle draw rider
func (EmissaryOfMoonRed) NotImplemented()             {}
func (c EmissaryOfMoonRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
