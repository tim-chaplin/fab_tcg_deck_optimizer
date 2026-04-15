// Arcanic Spike — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 5, Yellow 4, Blue 3.
// Text: "If you've dealt arcane damage this turn, this gets +2{p}."
//
// Simplification: assume any attack triggers some Runechant arcane damage when played, so the
// "dealt arcane damage this turn" clause is always satisfied. The +2{p} is baked into Attack().
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var arcanicSpikeTypes = map[string]bool{"Runeblade": true, "Action": true, "Attack": true}

type ArcanicSpikeRed struct{}

func (ArcanicSpikeRed) Name() string                 { return "Arcanic Spike (Red)" }
func (ArcanicSpikeRed) Cost() int                    { return 2 }
func (ArcanicSpikeRed) Pitch() int                   { return 1 }
func (ArcanicSpikeRed) Attack() int                  { return 7 }
func (ArcanicSpikeRed) Defense() int                 { return 3 }
func (ArcanicSpikeRed) Types() map[string]bool       { return arcanicSpikeTypes }
func (ArcanicSpikeRed) GoAgain() bool                { return false }
func (c ArcanicSpikeRed) Play(*card.TurnState) int   { return c.Attack() }

type ArcanicSpikeYellow struct{}

func (ArcanicSpikeYellow) Name() string                 { return "Arcanic Spike (Yellow)" }
func (ArcanicSpikeYellow) Cost() int                    { return 2 }
func (ArcanicSpikeYellow) Pitch() int                   { return 2 }
func (ArcanicSpikeYellow) Attack() int                  { return 6 }
func (ArcanicSpikeYellow) Defense() int                 { return 3 }
func (ArcanicSpikeYellow) Types() map[string]bool       { return arcanicSpikeTypes }
func (ArcanicSpikeYellow) GoAgain() bool                { return false }
func (c ArcanicSpikeYellow) Play(*card.TurnState) int   { return c.Attack() }

type ArcanicSpikeBlue struct{}

func (ArcanicSpikeBlue) Name() string                 { return "Arcanic Spike (Blue)" }
func (ArcanicSpikeBlue) Cost() int                    { return 2 }
func (ArcanicSpikeBlue) Pitch() int                   { return 3 }
func (ArcanicSpikeBlue) Attack() int                  { return 5 }
func (ArcanicSpikeBlue) Defense() int                 { return 3 }
func (ArcanicSpikeBlue) Types() map[string]bool       { return arcanicSpikeTypes }
func (ArcanicSpikeBlue) GoAgain() bool                { return false }
func (c ArcanicSpikeBlue) Play(*card.TurnState) int   { return c.Attack() }
