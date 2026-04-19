// Destructive Deliberation — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue
// 3. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, create a Ponder token."
//
// Simplification: Ponder token creation isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var destructiveDeliberationTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type DestructiveDeliberationRed struct{}

func (DestructiveDeliberationRed) ID() card.ID                 { return card.DestructiveDeliberationRed }
func (DestructiveDeliberationRed) Name() string                { return "Destructive Deliberation (Red)" }
func (DestructiveDeliberationRed) Cost() int                   { return 2 }
func (DestructiveDeliberationRed) Pitch() int                  { return 1 }
func (DestructiveDeliberationRed) Attack() int                 { return 5 }
func (DestructiveDeliberationRed) Defense() int                { return 2 }
func (DestructiveDeliberationRed) Types() card.TypeSet         { return destructiveDeliberationTypes }
func (DestructiveDeliberationRed) GoAgain() bool               { return false }
func (c DestructiveDeliberationRed) Play(s *card.TurnState) int { return destructiveDeliberationDamage(c.Attack()) }

type DestructiveDeliberationYellow struct{}

func (DestructiveDeliberationYellow) ID() card.ID                 { return card.DestructiveDeliberationYellow }
func (DestructiveDeliberationYellow) Name() string                { return "Destructive Deliberation (Yellow)" }
func (DestructiveDeliberationYellow) Cost() int                   { return 2 }
func (DestructiveDeliberationYellow) Pitch() int                  { return 2 }
func (DestructiveDeliberationYellow) Attack() int                 { return 4 }
func (DestructiveDeliberationYellow) Defense() int                { return 2 }
func (DestructiveDeliberationYellow) Types() card.TypeSet         { return destructiveDeliberationTypes }
func (DestructiveDeliberationYellow) GoAgain() bool               { return false }
func (c DestructiveDeliberationYellow) Play(s *card.TurnState) int { return destructiveDeliberationDamage(c.Attack()) }

type DestructiveDeliberationBlue struct{}

func (DestructiveDeliberationBlue) ID() card.ID                 { return card.DestructiveDeliberationBlue }
func (DestructiveDeliberationBlue) Name() string                { return "Destructive Deliberation (Blue)" }
func (DestructiveDeliberationBlue) Cost() int                   { return 2 }
func (DestructiveDeliberationBlue) Pitch() int                  { return 3 }
func (DestructiveDeliberationBlue) Attack() int                 { return 3 }
func (DestructiveDeliberationBlue) Defense() int                { return 2 }
func (DestructiveDeliberationBlue) Types() card.TypeSet         { return destructiveDeliberationTypes }
func (DestructiveDeliberationBlue) GoAgain() bool               { return false }
func (c DestructiveDeliberationBlue) Play(s *card.TurnState) int { return destructiveDeliberationDamage(c.Attack()) }

// destructiveDeliberationDamage is a breadcrumb for the on-hit "create a Ponder token" rider —
// Ponder tokens aren't tracked (see TODO.md).
func destructiveDeliberationDamage(attack int) int {
	if card.LikelyToHit(attack) {
		// TODO: model on-hit Ponder token creation rider.
	}
	return attack
}
