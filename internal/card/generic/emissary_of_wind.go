// Emissary of Wind — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 2. Only printed in
// Red.
//
// Text: "When this attacks, you may put a card from your hand on the bottom of your deck. If you
// do, this gets **go again**."
//
// Simplification: Hand-cycle-for-go-again rider isn't modelled.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var emissaryOfWindTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type EmissaryOfWindRed struct{}

func (EmissaryOfWindRed) ID() card.ID                 { return card.EmissaryOfWindRed }
func (EmissaryOfWindRed) Name() string                { return "Emissary of Wind (Red)" }
func (EmissaryOfWindRed) Cost(*card.TurnState) int                   { return 0 }
func (EmissaryOfWindRed) Pitch() int                  { return 1 }
func (EmissaryOfWindRed) Attack() int                 { return 4 }
func (EmissaryOfWindRed) Defense() int                { return 2 }
func (EmissaryOfWindRed) Types() card.TypeSet         { return emissaryOfWindTypes }
func (EmissaryOfWindRed) GoAgain() bool               { return false }
func (c EmissaryOfWindRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
