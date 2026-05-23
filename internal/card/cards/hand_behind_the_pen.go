// Hand Behind the Pen — Generic Action - Attack. Cost 2, Pitch 1, Power 6, Defense 2. Only
// printed in Red.
//
// Text: "When this hits a hero, turn a card in their arsenal face-up, then banish a non-attack
// action card from their arsenal."
//
// The on-hit arsenal rider isn't modelled: the single-turn solver doesn't track the opponent's
// arsenal contents, so flipping a face-down card and banishing a non-attack action card from it
// has no measurable single-turn effect. Modelled as the printed attack.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (HandBehindThePenRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
