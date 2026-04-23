// Regurgitating Slog — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Regurgitating Slog, you may banish a card named Sloggism
// from your graveyard. If you do, Regurgitating Slog gains **dominate**."
//
// Simplification: Riders described above aren't modelled; Play returns base power.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var regurgitatingSlogTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type RegurgitatingSlogRed struct{}

func (RegurgitatingSlogRed) ID() card.ID                 { return card.RegurgitatingSlogRed }
func (RegurgitatingSlogRed) Name() string                { return "Regurgitating Slog (Red)" }
func (RegurgitatingSlogRed) Cost(*card.TurnState) int                   { return 2 }
func (RegurgitatingSlogRed) Pitch() int                  { return 1 }
func (RegurgitatingSlogRed) Attack() int                 { return 6 }
func (RegurgitatingSlogRed) Defense() int                { return 2 }
func (RegurgitatingSlogRed) Types() card.TypeSet         { return regurgitatingSlogTypes }
func (RegurgitatingSlogRed) GoAgain() bool               { return false }
func (c RegurgitatingSlogRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type RegurgitatingSlogYellow struct{}

func (RegurgitatingSlogYellow) ID() card.ID                 { return card.RegurgitatingSlogYellow }
func (RegurgitatingSlogYellow) Name() string                { return "Regurgitating Slog (Yellow)" }
func (RegurgitatingSlogYellow) Cost(*card.TurnState) int                   { return 2 }
func (RegurgitatingSlogYellow) Pitch() int                  { return 2 }
func (RegurgitatingSlogYellow) Attack() int                 { return 5 }
func (RegurgitatingSlogYellow) Defense() int                { return 2 }
func (RegurgitatingSlogYellow) Types() card.TypeSet         { return regurgitatingSlogTypes }
func (RegurgitatingSlogYellow) GoAgain() bool               { return false }
func (c RegurgitatingSlogYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type RegurgitatingSlogBlue struct{}

func (RegurgitatingSlogBlue) ID() card.ID                 { return card.RegurgitatingSlogBlue }
func (RegurgitatingSlogBlue) Name() string                { return "Regurgitating Slog (Blue)" }
func (RegurgitatingSlogBlue) Cost(*card.TurnState) int                   { return 2 }
func (RegurgitatingSlogBlue) Pitch() int                  { return 3 }
func (RegurgitatingSlogBlue) Attack() int                 { return 4 }
func (RegurgitatingSlogBlue) Defense() int                { return 2 }
func (RegurgitatingSlogBlue) Types() card.TypeSet         { return regurgitatingSlogTypes }
func (RegurgitatingSlogBlue) GoAgain() bool               { return false }
func (c RegurgitatingSlogBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
