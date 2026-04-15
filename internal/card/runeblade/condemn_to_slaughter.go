// Condemn to Slaughter — Runeblade Action. Cost 1, Defense 3, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Your next Runeblade attack this turn gets +N{p}. You may destroy an aura you control. If
// you do, each opponent destroys an aura permanent they control. Go again."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: the aura-trade rider is ignored. The +N{p} fires only if a Runeblade attack
// (an attack action card OR a weapon swing) follows later in this turn's ordering (peeking
// TurnState.CardsRemaining); in that case Play returns N as the bonus damage it will eventually
// confer. The opponent-aura-destruction clause is ignored.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var condemnToSlaughterTypes = map[string]bool{"Runeblade": true, "Action": true}

type CondemnToSlaughterRed struct{}

func (CondemnToSlaughterRed) Name() string               { return "Condemn to Slaughter (Red)" }
func (CondemnToSlaughterRed) Cost() int                  { return 1 }
func (CondemnToSlaughterRed) Pitch() int                 { return 1 }
func (CondemnToSlaughterRed) Attack() int                { return 0 }
func (CondemnToSlaughterRed) Defense() int               { return 3 }
func (CondemnToSlaughterRed) Types() map[string]bool     { return condemnToSlaughterTypes }
func (CondemnToSlaughterRed) GoAgain() bool              { return true }
func (CondemnToSlaughterRed) Play(s *card.TurnState) int { return condemnToSlaughterBonus(s, 3) }

type CondemnToSlaughterYellow struct{}

func (CondemnToSlaughterYellow) Name() string               { return "Condemn to Slaughter (Yellow)" }
func (CondemnToSlaughterYellow) Cost() int                  { return 1 }
func (CondemnToSlaughterYellow) Pitch() int                 { return 2 }
func (CondemnToSlaughterYellow) Attack() int                { return 0 }
func (CondemnToSlaughterYellow) Defense() int               { return 3 }
func (CondemnToSlaughterYellow) Types() map[string]bool     { return condemnToSlaughterTypes }
func (CondemnToSlaughterYellow) GoAgain() bool              { return true }
func (CondemnToSlaughterYellow) Play(s *card.TurnState) int { return condemnToSlaughterBonus(s, 2) }

type CondemnToSlaughterBlue struct{}

func (CondemnToSlaughterBlue) Name() string               { return "Condemn to Slaughter (Blue)" }
func (CondemnToSlaughterBlue) Cost() int                  { return 1 }
func (CondemnToSlaughterBlue) Pitch() int                 { return 3 }
func (CondemnToSlaughterBlue) Attack() int                { return 0 }
func (CondemnToSlaughterBlue) Defense() int               { return 3 }
func (CondemnToSlaughterBlue) Types() map[string]bool     { return condemnToSlaughterTypes }
func (CondemnToSlaughterBlue) GoAgain() bool              { return true }
func (CondemnToSlaughterBlue) Play(s *card.TurnState) int { return condemnToSlaughterBonus(s, 1) }

// condemnToSlaughterBonus returns n if some Runeblade attack (attack action card or weapon swing)
// is scheduled later this turn, otherwise 0.
func condemnToSlaughterBonus(s *card.TurnState, n int) int {
	for _, c := range s.CardsRemaining {
		t := c.Types()
		if !t["Runeblade"] {
			continue
		}
		if t["Attack"] || t["Weapon"] {
			return n
		}
	}
	return 0
}
