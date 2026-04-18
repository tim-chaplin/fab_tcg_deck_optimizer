// Pursue to the Pits of Despair — Generic Action - Attack. Cost 1, Pitch 1, Power 5, Defense 3.
// Only printed in Red.
//
// Text: "When this hits a hero, **mark** them."
//
// Simplification: Mark on hit isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var pursueToThePitsOfDespairTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type PursueToThePitsOfDespairRed struct{}

func (PursueToThePitsOfDespairRed) ID() card.ID                 { return card.PursueToThePitsOfDespairRed }
func (PursueToThePitsOfDespairRed) Name() string                { return "Pursue to the Pits of Despair (Red)" }
func (PursueToThePitsOfDespairRed) Cost() int                   { return 1 }
func (PursueToThePitsOfDespairRed) Pitch() int                  { return 1 }
func (PursueToThePitsOfDespairRed) Attack() int                 { return 5 }
func (PursueToThePitsOfDespairRed) Defense() int                { return 3 }
func (PursueToThePitsOfDespairRed) Types() card.TypeSet         { return pursueToThePitsOfDespairTypes }
func (PursueToThePitsOfDespairRed) GoAgain() bool               { return false }
func (c PursueToThePitsOfDespairRed) Play(s *card.TurnState) int { return c.Attack() }
