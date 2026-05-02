// Mauvrion Skies — Runeblade Action. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "The next Runeblade attack action card you play this turn gains go again and 'If
// this hits, create N Runechant tokens.'" (Red N=3, Yellow N=2, Blue N=1.)
//
// Targets the next Runeblade attack action card in CardsRemaining (weapons don't qualify);
// flips its GrantedGoAgain and appends to its OnHit.

package cards

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var mauvrionSkiesTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

// mauvrionTargetMatches accepts Runeblade attack action cards (weapons don't qualify).
func mauvrionTargetMatches(target *sim.CardState) bool {
	t := target.Card.Types()
	return t.Has(card.TypeRuneblade) && t.IsAttackAction()
}

// mauvrionSkiesPlay grants the next matching attack go-again and an on-hit n-runechant
// rider.
func mauvrionSkiesPlay(s *sim.TurnState, selfState *sim.CardState, source sim.Card, n int) {
	for _, pc := range s.CardsRemaining {
		if mauvrionTargetMatches(pc) {
			pc.GrantedGoAgain = true
			text := onHitRunechantText[source.ID()]
			target := pc
			pc.OnHit = append(pc.OnHit, func(state *sim.TurnState) {
				created := state.CreateRunechants(n)
				state.AddValue(created)
				state.LogPostTrigger(sim.DisplayName(target.Card), text, created)
			})
			break
		}
	}
	s.Log(selfState, 0)
}

// onHitRunechantText is the precomputed rider line for each Mauvrion Skies / Runic Reaping
// printing — built once so the OnHit closure doesn't fmt.Sprintf on the hot path.
var onHitRunechantText = func() map[ids.CardID]string {
	out := make(map[ids.CardID]string, 6)
	for _, p := range []struct {
		c sim.Card
		n int
	}{
		{MauvrionSkiesRed{}, 3},
		{MauvrionSkiesYellow{}, 2},
		{MauvrionSkiesBlue{}, 1},
		{RunicReapingRed{}, 1},
		{RunicReapingYellow{}, 2},
		{RunicReapingBlue{}, 3},
	} {
		out[p.c.ID()] = fmt.Sprintf("%s created %d runechants on hit", sim.DisplayName(p.c), p.n)
	}
	return out
}()

type MauvrionSkiesRed struct{}

func (MauvrionSkiesRed) ID() ids.CardID          { return ids.MauvrionSkiesRed }
func (MauvrionSkiesRed) Name() string            { return "Mauvrion Skies" }
func (MauvrionSkiesRed) Cost(*sim.TurnState) int { return 0 }
func (MauvrionSkiesRed) Pitch() int              { return 1 }
func (MauvrionSkiesRed) Attack() int             { return 0 }
func (MauvrionSkiesRed) Defense() int            { return 2 }
func (MauvrionSkiesRed) Types() card.TypeSet     { return mauvrionSkiesTypes }
func (MauvrionSkiesRed) GoAgain() bool           { return true }
func (c MauvrionSkiesRed) Play(s *sim.TurnState, self *sim.CardState) {
	mauvrionSkiesPlay(s, self, c, 3)
}

type MauvrionSkiesYellow struct{}

func (MauvrionSkiesYellow) ID() ids.CardID          { return ids.MauvrionSkiesYellow }
func (MauvrionSkiesYellow) Name() string            { return "Mauvrion Skies" }
func (MauvrionSkiesYellow) Cost(*sim.TurnState) int { return 0 }
func (MauvrionSkiesYellow) Pitch() int              { return 2 }
func (MauvrionSkiesYellow) Attack() int             { return 0 }
func (MauvrionSkiesYellow) Defense() int            { return 2 }
func (MauvrionSkiesYellow) Types() card.TypeSet     { return mauvrionSkiesTypes }
func (MauvrionSkiesYellow) GoAgain() bool           { return true }
func (c MauvrionSkiesYellow) Play(s *sim.TurnState, self *sim.CardState) {
	mauvrionSkiesPlay(s, self, c, 2)
}

type MauvrionSkiesBlue struct{}

func (MauvrionSkiesBlue) ID() ids.CardID          { return ids.MauvrionSkiesBlue }
func (MauvrionSkiesBlue) Name() string            { return "Mauvrion Skies" }
func (MauvrionSkiesBlue) Cost(*sim.TurnState) int { return 0 }
func (MauvrionSkiesBlue) Pitch() int              { return 3 }
func (MauvrionSkiesBlue) Attack() int             { return 0 }
func (MauvrionSkiesBlue) Defense() int            { return 2 }
func (MauvrionSkiesBlue) Types() card.TypeSet     { return mauvrionSkiesTypes }
func (MauvrionSkiesBlue) GoAgain() bool           { return true }
func (c MauvrionSkiesBlue) Play(s *sim.TurnState, self *sim.CardState) {
	mauvrionSkiesPlay(s, self, c, 1)
}
