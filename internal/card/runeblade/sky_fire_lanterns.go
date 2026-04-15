// Sky Fire Lanterns — Runeblade Action. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Reveal the top card of your deck. If it's <same color as this variant>, create a
// Runechant token."
//
// Simplification: peek at the actual top card of the remaining deck (s.Deck[0]) and compare its
// pitch value to this card's pitch (pitch = color: 1 Red / 2 Yellow / 3 Blue). On match, credit
// +1 for the Runechant and set AuraCreated. Opts out of the hand-evaluation memo because the
// result depends on deck composition not captured by the memo key.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var skyFireLanternsTypes = map[string]bool{"Runeblade": true, "Action": true}

func skyFireLanternsPlay(selfPitch int, s *card.TurnState) int {
	if len(s.Deck) == 0 {
		return 0
	}
	if s.Deck[0].Pitch() != selfPitch {
		return 0
	}
	s.AuraCreated = true
	return 1
}

type SkyFireLanternsRed struct{}

func (SkyFireLanternsRed) Name() string                 { return "Sky Fire Lanterns (Red)" }
func (SkyFireLanternsRed) Cost() int                    { return 0 }
func (SkyFireLanternsRed) Pitch() int                   { return 1 }
func (SkyFireLanternsRed) Attack() int                  { return 0 }
func (SkyFireLanternsRed) Defense() int                 { return 2 }
func (SkyFireLanternsRed) Types() map[string]bool       { return skyFireLanternsTypes }
func (SkyFireLanternsRed) GoAgain() bool                { return true }
func (SkyFireLanternsRed) NoMemo()                      {} // value depends on top of deck
func (c SkyFireLanternsRed) Play(s *card.TurnState) int { return skyFireLanternsPlay(c.Pitch(), s) }

type SkyFireLanternsYellow struct{}

func (SkyFireLanternsYellow) Name() string                 { return "Sky Fire Lanterns (Yellow)" }
func (SkyFireLanternsYellow) Cost() int                    { return 0 }
func (SkyFireLanternsYellow) Pitch() int                   { return 2 }
func (SkyFireLanternsYellow) Attack() int                  { return 0 }
func (SkyFireLanternsYellow) Defense() int                 { return 2 }
func (SkyFireLanternsYellow) Types() map[string]bool       { return skyFireLanternsTypes }
func (SkyFireLanternsYellow) GoAgain() bool                { return true }
func (SkyFireLanternsYellow) NoMemo()                      {}
func (c SkyFireLanternsYellow) Play(s *card.TurnState) int { return skyFireLanternsPlay(c.Pitch(), s) }

type SkyFireLanternsBlue struct{}

func (SkyFireLanternsBlue) Name() string                 { return "Sky Fire Lanterns (Blue)" }
func (SkyFireLanternsBlue) Cost() int                    { return 0 }
func (SkyFireLanternsBlue) Pitch() int                   { return 3 }
func (SkyFireLanternsBlue) Attack() int                  { return 0 }
func (SkyFireLanternsBlue) Defense() int                 { return 2 }
func (SkyFireLanternsBlue) Types() map[string]bool       { return skyFireLanternsTypes }
func (SkyFireLanternsBlue) GoAgain() bool                { return true }
func (SkyFireLanternsBlue) NoMemo()                      {}
func (c SkyFireLanternsBlue) Play(s *card.TurnState) int { return skyFireLanternsPlay(c.Pitch(), s) }
