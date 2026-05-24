// Sound the Alarm — Generic Action - Attack. Cost 1, Pitch 1, Power 5, Defense 3. Only printed in
// Red.
//
// Text: "When this attacks a hero, they reveal their hand. If an attack reaction card is revealed
// this way, you may search your deck for a defense reaction card, reveal it, then shuffle and put
// it on top."
//
// Neither clause is modelled: the single-turn solver doesn't track the opponent's hand, so the
// attack-reaction gate is never satisfiable and the defense-reaction tutor never fires. Modelled
// as the printed attack.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (SoundTheAlarmRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
