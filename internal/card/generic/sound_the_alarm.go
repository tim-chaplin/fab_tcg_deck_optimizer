// Sound the Alarm — Generic Action - Attack. Cost 1, Pitch 1, Power 5, Defense 3. Only printed in
// Red.
//
// Text: "When this attacks a hero, they reveal their hand. If an attack reaction card is revealed
// this way, you may search your deck for a defense reaction card, reveal it, then shuffle and put
// it on top."
//
// Simplification: Deck search rider isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var soundTheAlarmTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type SoundTheAlarmRed struct{}

func (SoundTheAlarmRed) ID() card.ID                 { return card.SoundTheAlarmRed }
func (SoundTheAlarmRed) Name() string                { return "Sound the Alarm (Red)" }
func (SoundTheAlarmRed) Cost(*card.TurnState) int                   { return 1 }
func (SoundTheAlarmRed) Pitch() int                  { return 1 }
func (SoundTheAlarmRed) Attack() int                 { return 5 }
func (SoundTheAlarmRed) Defense() int                { return 3 }
func (SoundTheAlarmRed) Types() card.TypeSet         { return soundTheAlarmTypes }
func (SoundTheAlarmRed) GoAgain() bool               { return false }
func (c SoundTheAlarmRed) Play(s *card.TurnState) int { return c.Attack() }
