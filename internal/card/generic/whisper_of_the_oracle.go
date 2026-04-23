// Whisper of the Oracle — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 3.
//
// Text: "**Opt 4** **Go again**"
//
// Simplification: Opt isn't modelled.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var whisperOfTheOracleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type WhisperOfTheOracleRed struct{}

func (WhisperOfTheOracleRed) ID() card.ID                 { return card.WhisperOfTheOracleRed }
func (WhisperOfTheOracleRed) Name() string                { return "Whisper of the Oracle (Red)" }
func (WhisperOfTheOracleRed) Cost(*card.TurnState) int                   { return 0 }
func (WhisperOfTheOracleRed) Pitch() int                  { return 1 }
func (WhisperOfTheOracleRed) Attack() int                 { return 0 }
func (WhisperOfTheOracleRed) Defense() int                { return 3 }
func (WhisperOfTheOracleRed) Types() card.TypeSet         { return whisperOfTheOracleTypes }
func (WhisperOfTheOracleRed) GoAgain() bool               { return true }
func (WhisperOfTheOracleRed) Play(s *card.TurnState, _ *card.CardState) int { return 0 }

type WhisperOfTheOracleYellow struct{}

func (WhisperOfTheOracleYellow) ID() card.ID                 { return card.WhisperOfTheOracleYellow }
func (WhisperOfTheOracleYellow) Name() string                { return "Whisper of the Oracle (Yellow)" }
func (WhisperOfTheOracleYellow) Cost(*card.TurnState) int                   { return 0 }
func (WhisperOfTheOracleYellow) Pitch() int                  { return 2 }
func (WhisperOfTheOracleYellow) Attack() int                 { return 0 }
func (WhisperOfTheOracleYellow) Defense() int                { return 3 }
func (WhisperOfTheOracleYellow) Types() card.TypeSet         { return whisperOfTheOracleTypes }
func (WhisperOfTheOracleYellow) GoAgain() bool               { return true }
func (WhisperOfTheOracleYellow) Play(s *card.TurnState, _ *card.CardState) int { return 0 }

type WhisperOfTheOracleBlue struct{}

func (WhisperOfTheOracleBlue) ID() card.ID                 { return card.WhisperOfTheOracleBlue }
func (WhisperOfTheOracleBlue) Name() string                { return "Whisper of the Oracle (Blue)" }
func (WhisperOfTheOracleBlue) Cost(*card.TurnState) int                   { return 0 }
func (WhisperOfTheOracleBlue) Pitch() int                  { return 3 }
func (WhisperOfTheOracleBlue) Attack() int                 { return 0 }
func (WhisperOfTheOracleBlue) Defense() int                { return 3 }
func (WhisperOfTheOracleBlue) Types() card.TypeSet         { return whisperOfTheOracleTypes }
func (WhisperOfTheOracleBlue) GoAgain() bool               { return true }
func (WhisperOfTheOracleBlue) Play(s *card.TurnState, _ *card.CardState) int { return 0 }
