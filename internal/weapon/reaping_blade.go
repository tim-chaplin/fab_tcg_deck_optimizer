// Reaping Blade — Runeblade Weapon - Sword (2H). Power 3.
// Text: "Once per Turn Action - {r}: Attack. If a hero has more {h} than any other hero, they can't
// gain {h}."
//
// Simulation: modelled as an attack source costing 1 resource, dealing 3 damage. The
// health-symmetry rider is ignored (irrelevant to single-turn damage evaluation).
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package weapon

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var reapingBladeTypes = map[string]bool{
	"Runeblade": true,
	"Weapon":    true,
	"Sword":     true,
	"2H":        true,
}

type ReapingBlade struct{}

func (ReapingBlade) Name() string                 { return "Reaping Blade" }
func (ReapingBlade) Cost() int                    { return 1 }
func (ReapingBlade) Pitch() int                   { return 0 }
func (ReapingBlade) Attack() int                  { return 3 }
func (ReapingBlade) Defense() int                 { return 0 }
func (ReapingBlade) Types() map[string]bool       { return reapingBladeTypes }
func (ReapingBlade) GoAgain() bool                { return false }
func (ReapingBlade) Hands() int                   { return 2 }
func (c ReapingBlade) Play(*card.TurnState) int   { return c.Attack() }
