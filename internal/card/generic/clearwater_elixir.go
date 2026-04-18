// Clearwater Elixir — Generic Action. Cost 1, Pitch 1, Defense 3. Only printed in Red.
//
// Text: "Your next attack this turn gets +3{p}. You may destroy a Bloodrot Pox token you control.
// If you do, gain 1{h}. **Go again**"
//
// Simplification: Bloodrot Pox health-gain rider dropped. Scans TurnState.CardsRemaining for the
// first matching attack action card and credits the bonus assuming it will be played; if none is
// scheduled after this card, the bonus fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var clearwaterElixirTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// clearwaterElixirPlay returns 3 when a matching attack action card is scheduled later this turn.
func clearwaterElixirPlay(s *card.TurnState) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		return 3
	}
	return 0
}

type ClearwaterElixirRed struct{}

func (ClearwaterElixirRed) ID() card.ID                 { return card.ClearwaterElixirRed }
func (ClearwaterElixirRed) Name() string                { return "Clearwater Elixir (Red)" }
func (ClearwaterElixirRed) Cost() int                   { return 1 }
func (ClearwaterElixirRed) Pitch() int                  { return 1 }
func (ClearwaterElixirRed) Attack() int                 { return 0 }
func (ClearwaterElixirRed) Defense() int                { return 3 }
func (ClearwaterElixirRed) Types() card.TypeSet         { return clearwaterElixirTypes }
func (ClearwaterElixirRed) GoAgain() bool               { return true }
func (ClearwaterElixirRed) Play(s *card.TurnState) int { return clearwaterElixirPlay(s) }
