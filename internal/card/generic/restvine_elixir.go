// Restvine Elixir — Generic Action. Cost 1, Pitch 1, Defense 3. Only printed in Red.
//
// Text: "Your next attack this turn gets +3{p}. You may destroy an Inertia token you control. If
// you do, gain 1{h}. **Go again**"
//
// Simplification: Inertia health-gain rider dropped. Scans TurnState.CardsRemaining for the first
// matching attack action card and credits the bonus assuming it will be played; if none is
// scheduled after this card, the bonus fizzles.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var restvineElixirTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type RestvineElixirRed struct{}

func (RestvineElixirRed) ID() card.ID                 { return card.RestvineElixirRed }
func (RestvineElixirRed) Name() string                { return "Restvine Elixir (Red)" }
func (RestvineElixirRed) Cost(*card.TurnState) int                   { return 1 }
func (RestvineElixirRed) Pitch() int                  { return 1 }
func (RestvineElixirRed) Attack() int                 { return 0 }
func (RestvineElixirRed) Defense() int                { return 3 }
func (RestvineElixirRed) Types() card.TypeSet         { return restvineElixirTypes }
func (RestvineElixirRed) GoAgain() bool               { return true }
func (RestvineElixirRed) Play(s *card.TurnState, _ *card.CardState) int { return grantNextAttackActionBonus(s, 3) }
