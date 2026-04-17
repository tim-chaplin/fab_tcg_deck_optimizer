// Consuming Volition — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you've dealt arcane damage this turn, this gets 'When this hits a hero, they discard
// a card.'"
//
// Simplifications:
//   - "Dealt arcane damage this turn" is approximated as "a Runechant exists in play at the
//     moment Consuming Volition is played" — the same live-count check DiscountPerRunechant uses.
//     Runechants left in play will fire on this attack's damage step, so if any exist we know
//     arcane damage will occur this turn. We don't track whether arcane damage was already dealt
//     earlier in the chain; the common Runeblade play pattern is create-then-attack, where the
//     tokens are still live when Consuming Volition resolves.
//   - Assume the attack hits and the opponent discards when the rider is active. A discarded
//     card is valued at 3, mirroring the value we assign a drawn card for Drawn to the Dark
//     Dimension.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var consumingVolitionTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// discardRiderValue is the damage-equivalent we credit when the "discard a card" rider fires.
// Matches the "draw a card" value used by Drawn to the Dark Dimension.
const discardRiderValue = 3

// consumingVolitionDamage returns the base attack plus the discard rider when there's a live
// Runechant to fire on this attack. Extracted so all three printings share one implementation.
func consumingVolitionDamage(attack int, s *card.TurnState) int {
	if s != nil && s.Runechants > 0 {
		return attack + discardRiderValue
	}
	return attack
}

type ConsumingVolitionRed struct{}

func (ConsumingVolitionRed) ID() card.ID                  { return card.ConsumingVolitionRed }
func (ConsumingVolitionRed) Name() string                 { return "Consuming Volition (Red)" }
func (ConsumingVolitionRed) Cost() int                    { return 1 }
func (ConsumingVolitionRed) Pitch() int                   { return 1 }
func (ConsumingVolitionRed) Attack() int                  { return 4 }
func (ConsumingVolitionRed) Defense() int                 { return 3 }
func (ConsumingVolitionRed) Types() card.TypeSet          { return consumingVolitionTypes }
func (ConsumingVolitionRed) GoAgain() bool                { return false }
func (c ConsumingVolitionRed) Play(s *card.TurnState) int { return consumingVolitionDamage(c.Attack(), s) }

type ConsumingVolitionYellow struct{}

func (ConsumingVolitionYellow) ID() card.ID                  { return card.ConsumingVolitionYellow }
func (ConsumingVolitionYellow) Name() string                 { return "Consuming Volition (Yellow)" }
func (ConsumingVolitionYellow) Cost() int                    { return 1 }
func (ConsumingVolitionYellow) Pitch() int                   { return 2 }
func (ConsumingVolitionYellow) Attack() int                  { return 3 }
func (ConsumingVolitionYellow) Defense() int                 { return 3 }
func (ConsumingVolitionYellow) Types() card.TypeSet          { return consumingVolitionTypes }
func (ConsumingVolitionYellow) GoAgain() bool                { return false }
func (c ConsumingVolitionYellow) Play(s *card.TurnState) int { return consumingVolitionDamage(c.Attack(), s) }

type ConsumingVolitionBlue struct{}

func (ConsumingVolitionBlue) ID() card.ID                  { return card.ConsumingVolitionBlue }
func (ConsumingVolitionBlue) Name() string                 { return "Consuming Volition (Blue)" }
func (ConsumingVolitionBlue) Cost() int                    { return 1 }
func (ConsumingVolitionBlue) Pitch() int                   { return 3 }
func (ConsumingVolitionBlue) Attack() int                  { return 2 }
func (ConsumingVolitionBlue) Defense() int                 { return 3 }
func (ConsumingVolitionBlue) Types() card.TypeSet          { return consumingVolitionTypes }
func (ConsumingVolitionBlue) GoAgain() bool                { return false }
func (c ConsumingVolitionBlue) Play(s *card.TurnState) int { return consumingVolitionDamage(c.Attack(), s) }
