// Tongue Tied — Generic Action - Attack. Cost 3, Pitch 1, Power 7, Defense 2. Only printed in
// Red.
//
// Text: "When this hits a hero, turn a card in their arsenal face-up, then banish an instant
// card from their arsenal."
//
// The on-hit arsenal rider isn't modelled: the single-turn solver doesn't track the opponent's
// arsenal contents, so flipping a face-down card and banishing an instant card from it has no
// measurable single-turn effect. Modelled as the printed attack.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (TongueTiedRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
