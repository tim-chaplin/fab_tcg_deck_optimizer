// Sound the Alarm — Generic Action - Attack. Cost 1, Pitch 1, Power 5, Defense 3. Only printed in
// Red.
//
// Text: "When this attacks a hero, they reveal their hand. If an attack reaction card is revealed
// this way, you may search your deck for a defense reaction card, reveal it, then shuffle and put
// it on top."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: opponent hand reveal, defense-reaction deck search

func (c SoundTheAlarmRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
}
