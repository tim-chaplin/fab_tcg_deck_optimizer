// Wreck Havoc — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "Defense reactions can't be played to this chain link. When this hits a hero, you may turn
// a card in their arsenal face up, then destroy a defense reaction in their arsenal."
//
// Neither rider is modelled: the single-turn solver doesn't model the opponent playing defense
// reactions against our attacks and doesn't track their arsenal contents, so the lockout and the
// on-hit arsenal banish have no measurable single-turn effect. Modelled as the printed attack.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (WreckHavocRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState)    {}
func (WreckHavocYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
func (WreckHavocBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState)   {}
