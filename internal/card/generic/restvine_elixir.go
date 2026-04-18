// Restvine Elixir — Generic Action. Cost 1, Pitch 1, Defense 3. Only printed in Red.
//
// Text: "Your next attack this turn gets +3{p}. You may destroy an Inertia token you control. If
// you do, gain 1{h}. **Go again**"
//
// Simplification: Inertia health-gain rider dropped. Scans TurnState.CardsRemaining for the first
// matching attack action card and credits the bonus assuming it will be played; if none is
// scheduled after this card, the bonus fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var restvineElixirTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// restvineElixirPlay returns 3 when a matching attack action card is scheduled later this turn.
func restvineElixirPlay(s *card.TurnState) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		return 3
	}
	return 0
}

type RestvineElixirRed struct{}

func (RestvineElixirRed) ID() card.ID                 { return card.RestvineElixirRed }
func (RestvineElixirRed) Name() string                { return "Restvine Elixir (Red)" }
func (RestvineElixirRed) Cost() int                   { return 1 }
func (RestvineElixirRed) Pitch() int                  { return 1 }
func (RestvineElixirRed) Attack() int                 { return 0 }
func (RestvineElixirRed) Defense() int                { return 3 }
func (RestvineElixirRed) Types() card.TypeSet         { return restvineElixirTypes }
func (RestvineElixirRed) GoAgain() bool               { return true }
func (RestvineElixirRed) Play(s *card.TurnState) int { return restvineElixirPlay(s) }
