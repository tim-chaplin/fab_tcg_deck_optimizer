// Sound the Alarm — Generic Action - Attack. Cost 1, Pitch 1, Power 5, Defense 3. Only printed in
// Red.
//
// Text: "When this attacks a hero, they reveal their hand. If an attack reaction card is revealed
// this way, you may search your deck for a defense reaction card, reveal it, then shuffle and put
// it on top."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var soundTheAlarmTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type SoundTheAlarmRed struct{}

func (SoundTheAlarmRed) ID() ids.CardID           { return ids.SoundTheAlarmRed }
func (SoundTheAlarmRed) Name() string             { return "Sound the Alarm" }
func (SoundTheAlarmRed) Cost(*card.TurnState) int { return 1 }
func (SoundTheAlarmRed) Pitch() int               { return 1 }
func (SoundTheAlarmRed) Attack() int              { return 5 }
func (SoundTheAlarmRed) Defense() int             { return 3 }
func (SoundTheAlarmRed) Types() card.TypeSet      { return soundTheAlarmTypes }
func (SoundTheAlarmRed) GoAgain() bool            { return false }

// not implemented: opponent hand reveal, defense-reaction deck search
func (SoundTheAlarmRed) NotImplemented() {}
func (c SoundTheAlarmRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
