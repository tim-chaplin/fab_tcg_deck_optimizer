// Arcanic Spike — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 5, Yellow 4, Blue 3.
// Text: "If you've dealt arcane damage this turn, this gets +2{p}."
//
// Simplification: the "dealt arcane damage this turn" clause is approximated by
// `state.Runechants > 0` at Play time — the same live-count signal Consuming Volition uses for
// its discard rider (see PR #41). When a Runechant is live it'll fire on this attack, so arcane
// damage will be dealt this turn and the +2{p} bonus applies; otherwise Play returns the base
// printed attack alone. We don't track "arcane damage was already dealt earlier in the chain";
// the common Runeblade play pattern is create-then-attack, where tokens are still live when
// Arcanic Spike resolves.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var arcanicSpikeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// arcaneDamageBonus is the +2{p} gained when the "dealt arcane damage this turn" clause is live.
const arcaneDamageBonus = 2

// arcanicSpikeDamage returns the base attack plus the +2{p} rider when there's a live Runechant
// to fire on this attack. Extracted so all three printings share one implementation.
func arcanicSpikeDamage(attack int, s *card.TurnState) int {
	if s != nil && s.Runechants > 0 {
		return attack + arcaneDamageBonus
	}
	return attack
}

type ArcanicSpikeRed struct{}

func (ArcanicSpikeRed) ID() card.ID                    { return card.ArcanicSpikeRed }
func (ArcanicSpikeRed) Name() string                   { return "Arcanic Spike (Red)" }
func (ArcanicSpikeRed) Cost() int                      { return 2 }
func (ArcanicSpikeRed) Pitch() int                     { return 1 }
func (ArcanicSpikeRed) Attack() int                    { return 5 }
func (ArcanicSpikeRed) Defense() int                   { return 3 }
func (ArcanicSpikeRed) Types() card.TypeSet            { return arcanicSpikeTypes }
func (ArcanicSpikeRed) GoAgain() bool                  { return false }
func (c ArcanicSpikeRed) Play(s *card.TurnState) int   { return arcanicSpikeDamage(c.Attack(), s) }

type ArcanicSpikeYellow struct{}

func (ArcanicSpikeYellow) ID() card.ID                    { return card.ArcanicSpikeYellow }
func (ArcanicSpikeYellow) Name() string                   { return "Arcanic Spike (Yellow)" }
func (ArcanicSpikeYellow) Cost() int                      { return 2 }
func (ArcanicSpikeYellow) Pitch() int                     { return 2 }
func (ArcanicSpikeYellow) Attack() int                    { return 4 }
func (ArcanicSpikeYellow) Defense() int                   { return 3 }
func (ArcanicSpikeYellow) Types() card.TypeSet            { return arcanicSpikeTypes }
func (ArcanicSpikeYellow) GoAgain() bool                  { return false }
func (c ArcanicSpikeYellow) Play(s *card.TurnState) int   { return arcanicSpikeDamage(c.Attack(), s) }

type ArcanicSpikeBlue struct{}

func (ArcanicSpikeBlue) ID() card.ID                    { return card.ArcanicSpikeBlue }
func (ArcanicSpikeBlue) Name() string                   { return "Arcanic Spike (Blue)" }
func (ArcanicSpikeBlue) Cost() int                      { return 2 }
func (ArcanicSpikeBlue) Pitch() int                     { return 3 }
func (ArcanicSpikeBlue) Attack() int                    { return 3 }
func (ArcanicSpikeBlue) Defense() int                   { return 3 }
func (ArcanicSpikeBlue) Types() card.TypeSet            { return arcanicSpikeTypes }
func (ArcanicSpikeBlue) GoAgain() bool                  { return false }
func (c ArcanicSpikeBlue) Play(s *card.TurnState) int   { return arcanicSpikeDamage(c.Attack(), s) }
