// Arcane Polarity — Generic Instant.
//
// Text: "Gain 1{h} If you've been dealt arcane damage this turn, instead gain N{h}."
// (Red N=4, Yellow N=3, Blue N=2.)
//
// Gates the alternate gain on s.ArcaneIncomingDamage > 0 (seeded from the matchup's
// -arcane-incoming). Life gain is credited 1-to-1 with damage.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var arcanePolarityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

// arcanePolarityPlay credits the conditional life gain as the chain step.
func arcanePolarityPlay(s *sim.TurnState, self *sim.CardState, arcaneGain int) {
	gain := 1
	if s.ArcaneIncomingDamage > 0 {
		gain = arcaneGain
	}
	s.AddValue(gain)
	s.Log(self, gain)
}

type ArcanePolarityRed struct{}

func (ArcanePolarityRed) ID() ids.CardID          { return ids.ArcanePolarityRed }
func (ArcanePolarityRed) Name() string            { return "Arcane Polarity" }
func (ArcanePolarityRed) Cost(*sim.TurnState) int { return 0 }
func (ArcanePolarityRed) Pitch() int              { return 1 }
func (ArcanePolarityRed) Attack() int             { return 0 }
func (ArcanePolarityRed) Defense() int            { return 0 }
func (ArcanePolarityRed) Types() card.TypeSet     { return arcanePolarityTypes }
func (ArcanePolarityRed) GoAgain() bool           { return false }
func (ArcanePolarityRed) Play(s *sim.TurnState, self *sim.CardState) {
	arcanePolarityPlay(s, self, 4)
}

type ArcanePolarityYellow struct{}

func (ArcanePolarityYellow) ID() ids.CardID          { return ids.ArcanePolarityYellow }
func (ArcanePolarityYellow) Name() string            { return "Arcane Polarity" }
func (ArcanePolarityYellow) Cost(*sim.TurnState) int { return 0 }
func (ArcanePolarityYellow) Pitch() int              { return 2 }
func (ArcanePolarityYellow) Attack() int             { return 0 }
func (ArcanePolarityYellow) Defense() int            { return 0 }
func (ArcanePolarityYellow) Types() card.TypeSet     { return arcanePolarityTypes }
func (ArcanePolarityYellow) GoAgain() bool           { return false }
func (ArcanePolarityYellow) Play(s *sim.TurnState, self *sim.CardState) {
	arcanePolarityPlay(s, self, 3)
}

type ArcanePolarityBlue struct{}

func (ArcanePolarityBlue) ID() ids.CardID          { return ids.ArcanePolarityBlue }
func (ArcanePolarityBlue) Name() string            { return "Arcane Polarity" }
func (ArcanePolarityBlue) Cost(*sim.TurnState) int { return 0 }
func (ArcanePolarityBlue) Pitch() int              { return 3 }
func (ArcanePolarityBlue) Attack() int             { return 0 }
func (ArcanePolarityBlue) Defense() int            { return 0 }
func (ArcanePolarityBlue) Types() card.TypeSet     { return arcanePolarityTypes }
func (ArcanePolarityBlue) GoAgain() bool           { return false }
func (ArcanePolarityBlue) Play(s *sim.TurnState, self *sim.CardState) {
	arcanePolarityPlay(s, self, 2)
}
