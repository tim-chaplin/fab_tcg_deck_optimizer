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

	"github.com/tim-chaplin/fab-deck-optimizer/v2/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// mauvrionTargetMatches accepts Runeblade attack action cards (weapons don't qualify).
func mauvrionTargetMatches(ge card.GameEngine, target *card.CardState) bool {
	t := target.Card.Types(ge)
	return t.Has(card.TypeRuneblade) && t.IsAttackAction()
}

// mauvrionSkiesPlay grants the next matching attack go-again and an on-hit n-runechant
// rider.
func mauvrionSkiesPlay(ge card.GameEngine, l card.Logger, selfState *card.CardState, source card.Card, n int) {
	for _, pc := range ge.CardsRemaining() {
		if mauvrionTargetMatches(ge, pc) {
			pc.GrantedGoAgain = true
			text := onHitRunechantText[source.ID()]
			pc.OnHit = append(pc.OnHit, card.OnHitHandler{
				Fire:    onHitCreateRunechants,
				Source:  source,
				LogText: text,
				N:       n,
			})
			break
		}
	}
}

// onHitCreateRunechants fires the on-hit "create N runechants" rider attached to the
// targeted attack. Reads N and the precomputed log line off the handler; self is the
// targeted attack (whose name credits the trigger line).
func onHitCreateRunechants(ge card.GameEngine, l card.Logger, self *card.CardState, h *card.OnHitHandler) {
	ge.CreateRunechants(h.N)
	l.AppendPostTrigger(self.Card.DisplayName(), h.LogText, h.N)
}

// onHitRunechantText is the precomputed rider line for each Mauvrion Skies / Runic Reaping
// printing.
var onHitRunechantText = func() map[ids.CardID]string {
	out := make(map[ids.CardID]string, 6)
	for _, p := range []struct {
		c card.Card
		n int
	}{
		{MauvrionSkiesRed{}, 3},
		{MauvrionSkiesYellow{}, 2},
		{MauvrionSkiesBlue{}, 1},
		{RunicReapingRed{}, 1},
		{RunicReapingYellow{}, 2},
		{RunicReapingBlue{}, 3},
	} {
		out[p.c.ID()] = fmt.Sprintf("%s created %d runechants on hit", p.c.DisplayName(), p.n)
	}
	return out
}()

func (c MauvrionSkiesRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	mauvrionSkiesPlay(ge, l, self, c, 3)
}

func (c MauvrionSkiesYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	mauvrionSkiesPlay(ge, l, self, c, 2)
}

func (c MauvrionSkiesBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	mauvrionSkiesPlay(ge, l, self, c, 1)
}
