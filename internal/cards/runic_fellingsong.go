// Runic Fellingsong — Runeblade Action - Attack. Cost 3, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 7, Yellow 6, Blue 5.
// Text: "When this attacks, you may banish an aura from your graveyard. If you do, deal 1 arcane
// damage to target hero."
//
// Play credits Attack() plus 1 arcane when banishAuraFromGraveyard lands an aura in s.Banish.
// No aura in the graveyard → the banish rider fizzles and Play returns just Attack().

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var runicFellingsongTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// runicFellingsongPlay emits the chain step at printed power, then writes the banish-for-
// arcane rider as a sub-line under self when an aura was successfully banished from the
// graveyard. banishAuraFromGraveyard flips ArcaneDamageDealt internally as part of its
// arcane-damage payload.
func runicFellingsongPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	if n := banishAuraFromGraveyard(s); n > 0 {
		s.AddValue(n)
		s.LogRider(self, n, "Banished an aura, dealt 1 arcane damage")
	}
}

type RunicFellingsongRed struct{}

func (RunicFellingsongRed) ID() ids.CardID          { return ids.RunicFellingsongRed }
func (RunicFellingsongRed) Name() string            { return "Runic Fellingsong" }
func (RunicFellingsongRed) Cost(*sim.TurnState) int { return 3 }
func (RunicFellingsongRed) Pitch() int              { return 1 }
func (RunicFellingsongRed) Attack() int             { return 7 }
func (RunicFellingsongRed) Defense() int            { return 3 }
func (RunicFellingsongRed) Types() card.TypeSet     { return runicFellingsongTypes }
func (RunicFellingsongRed) GoAgain() bool           { return false }
func (RunicFellingsongRed) Play(s *sim.TurnState, self *sim.CardState) {
	runicFellingsongPlay(s, self)
}

type RunicFellingsongYellow struct{}

func (RunicFellingsongYellow) ID() ids.CardID          { return ids.RunicFellingsongYellow }
func (RunicFellingsongYellow) Name() string            { return "Runic Fellingsong" }
func (RunicFellingsongYellow) Cost(*sim.TurnState) int { return 3 }
func (RunicFellingsongYellow) Pitch() int              { return 2 }
func (RunicFellingsongYellow) Attack() int             { return 6 }
func (RunicFellingsongYellow) Defense() int            { return 3 }
func (RunicFellingsongYellow) Types() card.TypeSet     { return runicFellingsongTypes }
func (RunicFellingsongYellow) GoAgain() bool           { return false }
func (RunicFellingsongYellow) Play(s *sim.TurnState, self *sim.CardState) {
	runicFellingsongPlay(s, self)
}

type RunicFellingsongBlue struct{}

func (RunicFellingsongBlue) ID() ids.CardID          { return ids.RunicFellingsongBlue }
func (RunicFellingsongBlue) Name() string            { return "Runic Fellingsong" }
func (RunicFellingsongBlue) Cost(*sim.TurnState) int { return 3 }
func (RunicFellingsongBlue) Pitch() int              { return 3 }
func (RunicFellingsongBlue) Attack() int             { return 5 }
func (RunicFellingsongBlue) Defense() int            { return 3 }
func (RunicFellingsongBlue) Types() card.TypeSet     { return runicFellingsongTypes }
func (RunicFellingsongBlue) GoAgain() bool           { return false }
func (RunicFellingsongBlue) Play(s *sim.TurnState, self *sim.CardState) {
	runicFellingsongPlay(s, self)
}
