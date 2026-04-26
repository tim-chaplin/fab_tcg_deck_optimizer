// On a Knife Edge — Generic Action. Cost 0, Pitch 2, Defense 2. Only printed in Yellow.
//
// Text: "Your next sword attack this turn gains **go again**. **Go again**"

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var onAKnifeEdgeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type OnAKnifeEdgeYellow struct{}

func (OnAKnifeEdgeYellow) ID() card.ID                 { return card.OnAKnifeEdgeYellow }
func (OnAKnifeEdgeYellow) Name() string                { return "On a Knife Edge (Yellow)" }
func (OnAKnifeEdgeYellow) Cost(*card.TurnState) int                   { return 0 }
func (OnAKnifeEdgeYellow) Pitch() int                  { return 2 }
func (OnAKnifeEdgeYellow) Attack() int                 { return 0 }
func (OnAKnifeEdgeYellow) Defense() int                { return 2 }
func (OnAKnifeEdgeYellow) Types() card.TypeSet         { return onAKnifeEdgeTypes }
func (OnAKnifeEdgeYellow) GoAgain() bool               { return true }
// not implemented: next-sword-attack go-again grant (weapon chain not scanned)
func (OnAKnifeEdgeYellow) NotImplemented()             {}
func (OnAKnifeEdgeYellow) Play(s *card.TurnState, _ *card.CardState) int { return 0 }
