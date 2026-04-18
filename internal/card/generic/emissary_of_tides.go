// Emissary of Tides — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 2. Only printed in
// Red.
//
// Text: "When this attacks, you may put a card from your hand on the bottom of your deck. If you
// do, this gets +2{p}."
//
// Simplification: Hand-cycle-for-+2{p} rider isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var emissaryOfTidesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type EmissaryOfTidesRed struct{}

func (EmissaryOfTidesRed) ID() card.ID                 { return card.EmissaryOfTidesRed }
func (EmissaryOfTidesRed) Name() string                { return "Emissary of Tides (Red)" }
func (EmissaryOfTidesRed) Cost() int                   { return 0 }
func (EmissaryOfTidesRed) Pitch() int                  { return 1 }
func (EmissaryOfTidesRed) Attack() int                 { return 4 }
func (EmissaryOfTidesRed) Defense() int                { return 2 }
func (EmissaryOfTidesRed) Types() card.TypeSet         { return emissaryOfTidesTypes }
func (EmissaryOfTidesRed) GoAgain() bool               { return false }
func (c EmissaryOfTidesRed) Play(s *card.TurnState) int { return c.Attack() }
