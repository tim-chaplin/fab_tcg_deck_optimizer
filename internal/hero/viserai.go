// Viserai — Runeblade Hero, Young. Health 20, Intelligence 4.
// Text: "Whenever you play a Runeblade card, if you have played another
// 'non-attack' action card this turn, create a Runechant token."
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package hero

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var viseraiTypes = map[string]bool{"Runeblade": true, "Hero": true, "Young": true}

// Viserai is Young Viserai.
type Viserai struct{}

func (Viserai) Name() string           { return "Viserai" }
func (Viserai) Health() int            { return 20 }
func (Viserai) Intelligence() int      { return 4 }
func (Viserai) Types() map[string]bool { return viseraiTypes }

// OnCardPlayed implements Viserai's hero ability: whenever a Runeblade
// card is played, if another "non-attack action" (an Action that is not
// also an Attack) has already been played this turn, create a Runechant
// token (modelled as +1 damage).
func (Viserai) OnCardPlayed(played card.Card, s *card.TurnState) int {
	if !played.Types()["Runeblade"] {
		return 0
	}
	for _, c := range s.CardsPlayed {
		t := c.Types()
		if t["Action"] && !t["Attack"] {
			return 1
		}
	}
	return 0
}
