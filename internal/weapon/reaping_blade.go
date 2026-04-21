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

var reapingBladeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

type ReapingBlade struct{}

func (ReapingBlade) ID() card.ID                  { return card.ReapingBladeID }
func (ReapingBlade) Name() string                 { return "Reaping Blade" }
func (ReapingBlade) Cost(*card.TurnState) int                    { return 1 }
func (ReapingBlade) Pitch() int                   { return 0 }
func (ReapingBlade) Attack() int                  { return 3 }
func (ReapingBlade) Defense() int                 { return 0 }
func (ReapingBlade) Types() card.TypeSet           { return reapingBladeTypes }
func (ReapingBlade) GoAgain() bool                { return false }
func (ReapingBlade) Hands() int                   { return 2 }
func (c ReapingBlade) Play(*card.TurnState, *card.PlayedCard) int   { return c.Attack() }
