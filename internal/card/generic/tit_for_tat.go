// Tit for Tat — Generic Action. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "{t} target hero. {u} another target hero. **Go again**"
//
// Simplification: Freeze/unfreeze (tap/untap) isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var titForTatTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type TitForTatBlue struct{}

func (TitForTatBlue) ID() card.ID                 { return card.TitForTatBlue }
func (TitForTatBlue) Name() string                { return "Tit for Tat (Blue)" }
func (TitForTatBlue) Cost() int                   { return 0 }
func (TitForTatBlue) Pitch() int                  { return 3 }
func (TitForTatBlue) Attack() int                 { return 0 }
func (TitForTatBlue) Defense() int                { return 2 }
func (TitForTatBlue) Types() card.TypeSet         { return titForTatTypes }
func (TitForTatBlue) GoAgain() bool               { return true }
func (TitForTatBlue) Play(s *card.TurnState) int { return 0 }
