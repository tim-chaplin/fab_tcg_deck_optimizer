// Spellblade Assault — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When you attack with Spellblade Assault, create 2 Runechant tokens."
//
// Simplification: +1 damage per Runechant. Play returns power + 2.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var spellbladeAssaultTypes = map[string]bool{"Runeblade": true, "Action": true, "Attack": true}

type SpellbladeAssaultRed struct{}

func (SpellbladeAssaultRed) Name() string               { return "Spellblade Assault (Red)" }
func (SpellbladeAssaultRed) Cost() int                  { return 2 }
func (SpellbladeAssaultRed) Pitch() int                 { return 1 }
func (SpellbladeAssaultRed) Attack() int                { return 4 }
func (SpellbladeAssaultRed) Defense() int               { return 3 }
func (SpellbladeAssaultRed) Types() map[string]bool     { return spellbladeAssaultTypes }
func (SpellbladeAssaultRed) GoAgain() bool              { return false }
func (c SpellbladeAssaultRed) Play(*card.TurnState) int { return c.Attack() + 2 }

type SpellbladeAssaultYellow struct{}

func (SpellbladeAssaultYellow) Name() string               { return "Spellblade Assault (Yellow)" }
func (SpellbladeAssaultYellow) Cost() int                  { return 2 }
func (SpellbladeAssaultYellow) Pitch() int                 { return 2 }
func (SpellbladeAssaultYellow) Attack() int                { return 3 }
func (SpellbladeAssaultYellow) Defense() int               { return 3 }
func (SpellbladeAssaultYellow) Types() map[string]bool     { return spellbladeAssaultTypes }
func (SpellbladeAssaultYellow) GoAgain() bool              { return false }
func (c SpellbladeAssaultYellow) Play(*card.TurnState) int { return c.Attack() + 2 }

type SpellbladeAssaultBlue struct{}

func (SpellbladeAssaultBlue) Name() string               { return "Spellblade Assault (Blue)" }
func (SpellbladeAssaultBlue) Cost() int                  { return 2 }
func (SpellbladeAssaultBlue) Pitch() int                 { return 3 }
func (SpellbladeAssaultBlue) Attack() int                { return 2 }
func (SpellbladeAssaultBlue) Defense() int               { return 3 }
func (SpellbladeAssaultBlue) Types() map[string]bool     { return spellbladeAssaultTypes }
func (SpellbladeAssaultBlue) GoAgain() bool              { return false }
func (c SpellbladeAssaultBlue) Play(*card.TurnState) int { return c.Attack() + 2 }
