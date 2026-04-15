// Spellblade Strike — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "Create a Runechant token."
//
// Simplification: +1 damage per Runechant. Play returns power + 1.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var spellbladeStrikeTypes = map[string]bool{"Runeblade": true, "Action": true, "Attack": true}

type SpellbladeStrikeRed struct{}

func (SpellbladeStrikeRed) Name() string               { return "Spellblade Strike (Red)" }
func (SpellbladeStrikeRed) Cost() int                  { return 1 }
func (SpellbladeStrikeRed) Pitch() int                 { return 1 }
func (SpellbladeStrikeRed) Attack() int                { return 4 }
func (SpellbladeStrikeRed) Defense() int               { return 3 }
func (SpellbladeStrikeRed) Types() map[string]bool     { return spellbladeStrikeTypes }
func (SpellbladeStrikeRed) GoAgain() bool              { return false }
func (c SpellbladeStrikeRed) Play(s *card.TurnState) int { s.AuraCreated = true; return c.Attack() + 1 }

type SpellbladeStrikeYellow struct{}

func (SpellbladeStrikeYellow) Name() string               { return "Spellblade Strike (Yellow)" }
func (SpellbladeStrikeYellow) Cost() int                  { return 1 }
func (SpellbladeStrikeYellow) Pitch() int                 { return 2 }
func (SpellbladeStrikeYellow) Attack() int                { return 3 }
func (SpellbladeStrikeYellow) Defense() int               { return 3 }
func (SpellbladeStrikeYellow) Types() map[string]bool     { return spellbladeStrikeTypes }
func (SpellbladeStrikeYellow) GoAgain() bool              { return false }
func (c SpellbladeStrikeYellow) Play(s *card.TurnState) int { s.AuraCreated = true; return c.Attack() + 1 }

type SpellbladeStrikeBlue struct{}

func (SpellbladeStrikeBlue) Name() string               { return "Spellblade Strike (Blue)" }
func (SpellbladeStrikeBlue) Cost() int                  { return 1 }
func (SpellbladeStrikeBlue) Pitch() int                 { return 3 }
func (SpellbladeStrikeBlue) Attack() int                { return 2 }
func (SpellbladeStrikeBlue) Defense() int               { return 3 }
func (SpellbladeStrikeBlue) Types() map[string]bool     { return spellbladeStrikeTypes }
func (SpellbladeStrikeBlue) GoAgain() bool              { return false }
func (c SpellbladeStrikeBlue) Play(s *card.TurnState) int { s.AuraCreated = true; return c.Attack() + 1 }
