// Viserai — Runeblade Hero, Young. Health 20, Intelligence 4.
// Text: "Whenever you play a Runeblade card, if you have played another 'non-attack' action card
// this turn, create a Runechant token."

package hero

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var viseraiTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeHero, card.TypeYoung)

// Viserai is Young Viserai.
type Viserai struct{}

func (Viserai) ID() ids.HeroID      { return ids.ViseraiID }
func (Viserai) Name() string        { return "Viserai" }
func (Viserai) Health() int         { return 20 }
func (Viserai) Intelligence() int   { return 4 }
func (Viserai) Types() card.TypeSet { return viseraiTypes }

// OnCardPlayed implements Viserai's hero ability: whenever a Runeblade card is played, if a
// non-attack action (Action without Attack) has been played this turn, create a Runechant
// token.
func (Viserai) OnCardPlayed(played card.Card, s *card.TurnState) int {
	t := played.Types()
	// Weapon swings aren't "playing a card" and don't trigger Viserai.
	if !t.Has(card.TypeRuneblade) || t.Has(card.TypeWeapon) {
		return 0
	}
	if s.NonAttackActionPlayed {
		return s.CreateAndLogRunechants("Viserai", card.DisplayName(played), 1)
	}
	return 0
}
