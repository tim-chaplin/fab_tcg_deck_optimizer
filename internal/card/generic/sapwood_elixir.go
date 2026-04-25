// Sapwood Elixir — Generic Action. Cost 1, Pitch 1, Defense 3. Only printed in Red.
//
// Text: "Your next attack this turn gets +3{p}. You may destroy a Frailty token you control. If you
// do, gain 1{h}. **Go again**"

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sapwoodElixirTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type SapwoodElixirRed struct{}

func (SapwoodElixirRed) ID() card.ID                 { return card.SapwoodElixirRed }
func (SapwoodElixirRed) Name() string                { return "Sapwood Elixir (Red)" }
func (SapwoodElixirRed) Cost(*card.TurnState) int                   { return 1 }
func (SapwoodElixirRed) Pitch() int                  { return 1 }
func (SapwoodElixirRed) Attack() int                 { return 0 }
func (SapwoodElixirRed) Defense() int                { return 3 }
func (SapwoodElixirRed) Types() card.TypeSet         { return sapwoodElixirTypes }
func (SapwoodElixirRed) GoAgain() bool               { return true }
// not implemented: Frailty health-gain rider dropped (status tokens not tracked)
func (SapwoodElixirRed) NotImplemented()             {}
func (SapwoodElixirRed) Play(s *card.TurnState, _ *card.CardState) int { return grantNextAttackActionBonus(s, 3) }
