// Viserai — Runeblade Hero, Young. Health 20, Intelligence 4.
// Text: "Whenever you play a Runeblade card, if you have played another 'non-attack' action card
// this turn, create a Runechant token."
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package hero

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var viseraiTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeHero, card.TypeYoung)

// Viserai is Young Viserai.
type Viserai struct{}

func (Viserai) Name() string           { return "Viserai" }
func (Viserai) Health() int            { return 20 }
func (Viserai) Intelligence() int      { return 4 }
func (Viserai) Types() card.TypeSet    { return viseraiTypes }

// OnCardPlayed implements Viserai's hero ability: whenever a Runeblade card is played, if another
// "non-attack action" (an Action that is not also an Attack) has already been played this turn,
// create a Runechant token (modelled as +1 damage).
func (Viserai) OnCardPlayed(played card.Card, s *card.TurnState) int {
	t := played.Types()
	// Weapon swings are not "playing a card" and don't trigger Viserai.
	if !t.Has(card.TypeRuneblade) || t.Has(card.TypeWeapon) {
		return 0
	}
	for _, c := range s.CardsPlayed {
		ct := c.Types()
		if ct.Has(card.TypeAction) && !ct.Has(card.TypeAttack) {
			s.AuraCreated = true
			return 1
		}
	}
	return 0
}
