// Runic Reaping — Runeblade Action. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "The next Runeblade attack action card you play this turn gains 'When this hits, create N
// Runechant tokens'. If an attack card was pitched to play Runic Reaping, the next Runeblade
// attack action card you play this turn gains +1{p}. Go again."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Both riders target Runeblade attack action cards only — weapon swings don't qualify.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var runicReapingTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

// runicReapingTargetMatches accepts Runeblade attack action cards (weapons don't qualify).
func runicReapingTargetMatches(target *sim.CardState) bool {
	t := target.Card.Types()
	return t.Has(card.TypeRuneblade) && t.IsAttackAction()
}

type RunicReapingRed struct{}

func (RunicReapingRed) ID() ids.CardID          { return ids.RunicReapingRed }
func (RunicReapingRed) Name() string            { return "Runic Reaping" }
func (RunicReapingRed) Cost(*sim.TurnState) int { return 1 }
func (RunicReapingRed) Pitch() int              { return 1 }
func (RunicReapingRed) Attack() int             { return 0 }
func (RunicReapingRed) Defense() int            { return 2 }
func (RunicReapingRed) Types() card.TypeSet     { return runicReapingTypes }
func (RunicReapingRed) GoAgain() bool           { return true }
func (c RunicReapingRed) Play(s *sim.TurnState, self *sim.CardState) {
	runicReapingPlay(s, self, c, 3)
}

type RunicReapingYellow struct{}

func (RunicReapingYellow) ID() ids.CardID          { return ids.RunicReapingYellow }
func (RunicReapingYellow) Name() string            { return "Runic Reaping" }
func (RunicReapingYellow) Cost(*sim.TurnState) int { return 1 }
func (RunicReapingYellow) Pitch() int              { return 2 }
func (RunicReapingYellow) Attack() int             { return 0 }
func (RunicReapingYellow) Defense() int            { return 2 }
func (RunicReapingYellow) Types() card.TypeSet     { return runicReapingTypes }
func (RunicReapingYellow) GoAgain() bool           { return true }
func (c RunicReapingYellow) Play(s *sim.TurnState, self *sim.CardState) {
	runicReapingPlay(s, self, c, 2)
}

type RunicReapingBlue struct{}

func (RunicReapingBlue) ID() ids.CardID          { return ids.RunicReapingBlue }
func (RunicReapingBlue) Name() string            { return "Runic Reaping" }
func (RunicReapingBlue) Cost(*sim.TurnState) int { return 1 }
func (RunicReapingBlue) Pitch() int              { return 3 }
func (RunicReapingBlue) Attack() int             { return 0 }
func (RunicReapingBlue) Defense() int            { return 2 }
func (RunicReapingBlue) Types() card.TypeSet     { return runicReapingTypes }
func (RunicReapingBlue) GoAgain() bool           { return true }
func (c RunicReapingBlue) Play(s *sim.TurnState, self *sim.CardState) {
	runicReapingPlay(s, self, c, 1)
}

// runicReapingPlay buffs the next matching attack +1{p} when an attack card was pitched
// and appends an on-hit n-runechant rider.
func runicReapingPlay(s *sim.TurnState, selfState *sim.CardState, source sim.Card, n int) {
	var target *sim.CardState
	for _, pc := range s.CardsRemaining {
		if runicReapingTargetMatches(pc) {
			target = pc
			break
		}
	}
	if target == nil {
		s.Log(selfState, 0)
		return
	}
	for _, p := range selfState.PitchedToPlay {
		if p.Types().Has(card.TypeAttack) {
			target.BonusAttack++
			break
		}
	}
	text := onHitRunechantText[source.ID()]
	target.OnHit = append(target.OnHit, func(state *sim.TurnState) {
		created := state.CreateRunechants(n)
		state.AddValue(created)
		state.LogPostTrigger(sim.DisplayName(target.Card), text, created)
	})
	s.Log(selfState, 0)
}
