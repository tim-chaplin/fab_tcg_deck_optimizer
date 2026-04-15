// Vantage Point — Runeblade Action - Attack. Cost 3, Defense 3.
// Printed power: Red 7, Yellow 6, Blue 5.
// Text: "If you've played or created an aura this turn, this gets **overpower**."
//
// Sets TurnState.Overpower when the aura condition is met.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var vantagePointTypes = map[string]bool{"Runeblade": true, "Action": true, "Attack": true}

type VantagePointRed struct{}

func (VantagePointRed) Name() string               { return "Vantage Point (Red)" }
func (VantagePointRed) Cost() int                  { return 3 }
func (VantagePointRed) Pitch() int                 { return 1 }
func (VantagePointRed) Attack() int                { return 7 }
func (VantagePointRed) Defense() int               { return 3 }
func (VantagePointRed) Types() map[string]bool     { return vantagePointTypes }
func (VantagePointRed) GoAgain() bool              { return false }
func (c VantagePointRed) Play(s *card.TurnState) int { return vantagePointPlay(c.Attack(), s) }

type VantagePointYellow struct{}

func (VantagePointYellow) Name() string                 { return "Vantage Point (Yellow)" }
func (VantagePointYellow) Cost() int                    { return 3 }
func (VantagePointYellow) Pitch() int                   { return 2 }
func (VantagePointYellow) Attack() int                  { return 6 }
func (VantagePointYellow) Defense() int                 { return 3 }
func (VantagePointYellow) Types() map[string]bool       { return vantagePointTypes }
func (VantagePointYellow) GoAgain() bool                { return false }
func (c VantagePointYellow) Play(s *card.TurnState) int { return vantagePointPlay(c.Attack(), s) }

type VantagePointBlue struct{}

func (VantagePointBlue) Name() string                 { return "Vantage Point (Blue)" }
func (VantagePointBlue) Cost() int                    { return 3 }
func (VantagePointBlue) Pitch() int                   { return 3 }
func (VantagePointBlue) Attack() int                  { return 5 }
func (VantagePointBlue) Defense() int                 { return 3 }
func (VantagePointBlue) Types() map[string]bool       { return vantagePointTypes }
func (VantagePointBlue) GoAgain() bool                { return false }
func (c VantagePointBlue) Play(s *card.TurnState) int { return vantagePointPlay(c.Attack(), s) }

func vantagePointPlay(base int, s *card.TurnState) int {
	if s.AuraCreated || s.HasPlayedType("Aura") {
		s.Overpower = true
	}
	return base
}
