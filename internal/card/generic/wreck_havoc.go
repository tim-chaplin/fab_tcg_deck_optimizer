// Wreck Havoc — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "Defense reactions can't be played to this chain link. When this hits a hero, you may turn
// a card in their arsenal face up, then destroy a defense reaction in their arsenal."
//
// Simplification: Defense-reaction lockout and arsenal-banish aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var wreckHavocTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type WreckHavocRed struct{}

func (WreckHavocRed) ID() card.ID                 { return card.WreckHavocRed }
func (WreckHavocRed) Name() string                { return "Wreck Havoc (Red)" }
func (WreckHavocRed) Cost(*card.TurnState) int                   { return 2 }
func (WreckHavocRed) Pitch() int                  { return 1 }
func (WreckHavocRed) Attack() int                 { return 6 }
func (WreckHavocRed) Defense() int                { return 2 }
func (WreckHavocRed) Types() card.TypeSet         { return wreckHavocTypes }
func (WreckHavocRed) GoAgain() bool               { return false }
func (c WreckHavocRed) Play(s *card.TurnState) int { return wreckHavocDamage(c.Attack()) }

type WreckHavocYellow struct{}

func (WreckHavocYellow) ID() card.ID                 { return card.WreckHavocYellow }
func (WreckHavocYellow) Name() string                { return "Wreck Havoc (Yellow)" }
func (WreckHavocYellow) Cost(*card.TurnState) int                   { return 2 }
func (WreckHavocYellow) Pitch() int                  { return 2 }
func (WreckHavocYellow) Attack() int                 { return 5 }
func (WreckHavocYellow) Defense() int                { return 2 }
func (WreckHavocYellow) Types() card.TypeSet         { return wreckHavocTypes }
func (WreckHavocYellow) GoAgain() bool               { return false }
func (c WreckHavocYellow) Play(s *card.TurnState) int { return wreckHavocDamage(c.Attack()) }

type WreckHavocBlue struct{}

func (WreckHavocBlue) ID() card.ID                 { return card.WreckHavocBlue }
func (WreckHavocBlue) Name() string                { return "Wreck Havoc (Blue)" }
func (WreckHavocBlue) Cost(*card.TurnState) int                   { return 2 }
func (WreckHavocBlue) Pitch() int                  { return 3 }
func (WreckHavocBlue) Attack() int                 { return 4 }
func (WreckHavocBlue) Defense() int                { return 2 }
func (WreckHavocBlue) Types() card.TypeSet         { return wreckHavocTypes }
func (WreckHavocBlue) GoAgain() bool               { return false }
func (c WreckHavocBlue) Play(s *card.TurnState) int { return wreckHavocDamage(c.Attack()) }

// wreckHavocDamage is a breadcrumb for the on-hit "DR lockout + arsenal-face-up / banish DR"
// rider — not modelled yet (see TODO.md). LikelyToHit marks where the rider value would plug in.
func wreckHavocDamage(attack int) int {
	if card.LikelyToHit(attack) {
		// TODO: model on-hit arsenal manipulation rider.
	}
	return attack
}
