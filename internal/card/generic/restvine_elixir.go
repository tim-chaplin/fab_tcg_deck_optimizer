// Restvine Elixir — Generic Action. Cost 1, Pitch 1, Defense 3. Only printed in Red.
//
// Text: "Your next attack this turn gets +3{p}. You may destroy an Inertia token you control. If
// you do, gain 1{h}. **Go again**"

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var restvineElixirTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type RestvineElixirRed struct{}

func (RestvineElixirRed) ID() card.ID                 { return card.RestvineElixirRed }
func (RestvineElixirRed) Name() string                { return "Restvine Elixir" }
func (RestvineElixirRed) Cost(*card.TurnState) int                   { return 1 }
func (RestvineElixirRed) Pitch() int                  { return 1 }
func (RestvineElixirRed) Attack() int                 { return 0 }
func (RestvineElixirRed) Defense() int                { return 3 }
func (RestvineElixirRed) Types() card.TypeSet         { return restvineElixirTypes }
func (RestvineElixirRed) GoAgain() bool               { return true }
// not implemented: Inertia health-gain rider dropped (status tokens not tracked)
func (RestvineElixirRed) NotImplemented()             {}
func (RestvineElixirRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, grantNextAttackActionBonus(s, 3))
}