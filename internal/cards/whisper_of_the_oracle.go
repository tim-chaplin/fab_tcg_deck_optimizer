// Whisper of the Oracle — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 3.
//
// Text: "**Opt 4** **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var whisperOfTheOracleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

func whisperOfTheOraclePlay(s *sim.TurnState, self *sim.CardState) {
	s.Log(self, 0)
	s.Opt(4)
}

type WhisperOfTheOracleRed struct{}

func (WhisperOfTheOracleRed) ID() ids.CardID          { return ids.WhisperOfTheOracleRed }
func (WhisperOfTheOracleRed) Name() string            { return "Whisper of the Oracle" }
func (WhisperOfTheOracleRed) Cost(*sim.TurnState) int { return 0 }
func (WhisperOfTheOracleRed) Pitch() int              { return 1 }
func (WhisperOfTheOracleRed) Attack() int             { return 0 }
func (WhisperOfTheOracleRed) Defense() int            { return 3 }
func (WhisperOfTheOracleRed) Types() card.TypeSet     { return whisperOfTheOracleTypes }
func (WhisperOfTheOracleRed) GoAgain() bool           { return true }
func (WhisperOfTheOracleRed) Play(s *sim.TurnState, self *sim.CardState) {
	whisperOfTheOraclePlay(s, self)
}

type WhisperOfTheOracleYellow struct{}

func (WhisperOfTheOracleYellow) ID() ids.CardID          { return ids.WhisperOfTheOracleYellow }
func (WhisperOfTheOracleYellow) Name() string            { return "Whisper of the Oracle" }
func (WhisperOfTheOracleYellow) Cost(*sim.TurnState) int { return 0 }
func (WhisperOfTheOracleYellow) Pitch() int              { return 2 }
func (WhisperOfTheOracleYellow) Attack() int             { return 0 }
func (WhisperOfTheOracleYellow) Defense() int            { return 3 }
func (WhisperOfTheOracleYellow) Types() card.TypeSet     { return whisperOfTheOracleTypes }
func (WhisperOfTheOracleYellow) GoAgain() bool           { return true }
func (WhisperOfTheOracleYellow) Play(s *sim.TurnState, self *sim.CardState) {
	whisperOfTheOraclePlay(s, self)
}

type WhisperOfTheOracleBlue struct{}

func (WhisperOfTheOracleBlue) ID() ids.CardID          { return ids.WhisperOfTheOracleBlue }
func (WhisperOfTheOracleBlue) Name() string            { return "Whisper of the Oracle" }
func (WhisperOfTheOracleBlue) Cost(*sim.TurnState) int { return 0 }
func (WhisperOfTheOracleBlue) Pitch() int              { return 3 }
func (WhisperOfTheOracleBlue) Attack() int             { return 0 }
func (WhisperOfTheOracleBlue) Defense() int            { return 3 }
func (WhisperOfTheOracleBlue) Types() card.TypeSet     { return whisperOfTheOracleTypes }
func (WhisperOfTheOracleBlue) GoAgain() bool           { return true }
func (WhisperOfTheOracleBlue) Play(s *sim.TurnState, self *sim.CardState) {
	whisperOfTheOraclePlay(s, self)
}
