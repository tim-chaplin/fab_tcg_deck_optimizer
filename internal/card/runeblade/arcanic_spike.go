// Arcanic Spike — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 5, Yellow 4, Blue 3.
// Text: "If you've dealt arcane damage this turn, this gets +2{p}."
//
// Rider reads TurnState.ArcaneDamageDealt: when set at Play time, +2{p}; otherwise printed
// attack alone.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var arcanicSpikeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// arcaneDamageBonus is the +2{p} gained when the "dealt arcane damage this turn" clause is live.
const arcaneDamageBonus = 2

// arcanicSpikeDamage returns the base attack plus the +2{p} rider when ArcaneDamageDealt is set.
func arcanicSpikeDamage(attack int, s *card.TurnState) int {
	if s != nil && s.ArcaneDamageDealt {
		return attack + arcaneDamageBonus
	}
	return attack
}

type ArcanicSpikeRed struct{}

func (ArcanicSpikeRed) ID() card.ID                    { return card.ArcanicSpikeRed }
func (ArcanicSpikeRed) Name() string                   { return "Arcanic Spike (Red)" }
func (ArcanicSpikeRed) Cost(*card.TurnState) int                      { return 2 }
func (ArcanicSpikeRed) Pitch() int                     { return 1 }
func (ArcanicSpikeRed) Attack() int                    { return 5 }
func (ArcanicSpikeRed) Defense() int                   { return 3 }
func (ArcanicSpikeRed) Types() card.TypeSet            { return arcanicSpikeTypes }
func (ArcanicSpikeRed) GoAgain() bool                  { return false }
func (c ArcanicSpikeRed) Play(s *card.TurnState) int   { return arcanicSpikeDamage(c.Attack(), s) }

type ArcanicSpikeYellow struct{}

func (ArcanicSpikeYellow) ID() card.ID                    { return card.ArcanicSpikeYellow }
func (ArcanicSpikeYellow) Name() string                   { return "Arcanic Spike (Yellow)" }
func (ArcanicSpikeYellow) Cost(*card.TurnState) int                      { return 2 }
func (ArcanicSpikeYellow) Pitch() int                     { return 2 }
func (ArcanicSpikeYellow) Attack() int                    { return 4 }
func (ArcanicSpikeYellow) Defense() int                   { return 3 }
func (ArcanicSpikeYellow) Types() card.TypeSet            { return arcanicSpikeTypes }
func (ArcanicSpikeYellow) GoAgain() bool                  { return false }
func (c ArcanicSpikeYellow) Play(s *card.TurnState) int   { return arcanicSpikeDamage(c.Attack(), s) }

type ArcanicSpikeBlue struct{}

func (ArcanicSpikeBlue) ID() card.ID                    { return card.ArcanicSpikeBlue }
func (ArcanicSpikeBlue) Name() string                   { return "Arcanic Spike (Blue)" }
func (ArcanicSpikeBlue) Cost(*card.TurnState) int                      { return 2 }
func (ArcanicSpikeBlue) Pitch() int                     { return 3 }
func (ArcanicSpikeBlue) Attack() int                    { return 3 }
func (ArcanicSpikeBlue) Defense() int                   { return 3 }
func (ArcanicSpikeBlue) Types() card.TypeSet            { return arcanicSpikeTypes }
func (ArcanicSpikeBlue) GoAgain() bool                  { return false }
func (c ArcanicSpikeBlue) Play(s *card.TurnState) int   { return arcanicSpikeDamage(c.Attack(), s) }
