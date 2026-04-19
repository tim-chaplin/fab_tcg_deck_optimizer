// Humble — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, they lose all hero card abilities until the end of their next
// turn."
//
// Simplification: Hero-ability suppression rider isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var humbleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type HumbleRed struct{}

func (HumbleRed) ID() card.ID                 { return card.HumbleRed }
func (HumbleRed) Name() string                { return "Humble (Red)" }
func (HumbleRed) Cost() int                   { return 2 }
func (HumbleRed) Pitch() int                  { return 1 }
func (HumbleRed) Attack() int                 { return 6 }
func (HumbleRed) Defense() int                { return 2 }
func (HumbleRed) Types() card.TypeSet         { return humbleTypes }
func (HumbleRed) GoAgain() bool               { return false }
func (c HumbleRed) Play(s *card.TurnState) int { return humbleDamage(c.Attack()) }

type HumbleYellow struct{}

func (HumbleYellow) ID() card.ID                 { return card.HumbleYellow }
func (HumbleYellow) Name() string                { return "Humble (Yellow)" }
func (HumbleYellow) Cost() int                   { return 2 }
func (HumbleYellow) Pitch() int                  { return 2 }
func (HumbleYellow) Attack() int                 { return 5 }
func (HumbleYellow) Defense() int                { return 2 }
func (HumbleYellow) Types() card.TypeSet         { return humbleTypes }
func (HumbleYellow) GoAgain() bool               { return false }
func (c HumbleYellow) Play(s *card.TurnState) int { return humbleDamage(c.Attack()) }

type HumbleBlue struct{}

func (HumbleBlue) ID() card.ID                 { return card.HumbleBlue }
func (HumbleBlue) Name() string                { return "Humble (Blue)" }
func (HumbleBlue) Cost() int                   { return 2 }
func (HumbleBlue) Pitch() int                  { return 3 }
func (HumbleBlue) Attack() int                 { return 4 }
func (HumbleBlue) Defense() int                { return 2 }
func (HumbleBlue) Types() card.TypeSet         { return humbleTypes }
func (HumbleBlue) GoAgain() bool               { return false }
func (c HumbleBlue) Play(s *card.TurnState) int { return humbleDamage(c.Attack()) }

// humbleDamage is a breadcrumb for the on-hit "lose all hero card abilities" rider — not
// modelled yet (see TODO.md).
func humbleDamage(attack int) int {
	if card.LikelyToHit(attack) {
		// TODO: model on-hit hero-ability suppression rider.
	}
	return attack
}
