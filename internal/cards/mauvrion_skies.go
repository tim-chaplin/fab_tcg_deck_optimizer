// Mauvrion Skies — Runeblade Action. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "The next Runeblade attack action card you play this turn gains go again and 'If this
// hits, create N Runechant tokens.'"
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling splits the two grants by when they need to resolve:
//   - Go again is a static property of the target card (not gated on how the attack plays
//     out), and must be visible before the target's chain-legality check. Play scans
//     CardsRemaining for the first matching Runeblade attack action card and flips its
//     GrantedGoAgain immediately.
//   - The "if this hits" Runechant rider depends on the target's fully-resolved attack
//     state — the target's own Play effects, hero triggers, and aura triggers can shift
//     what LikelyToHit sees (e.g. Drowning Dire gaining Dominate from an aura created
//     mid-chain). Play registers an EphemeralAttackTrigger that fires after those settle;
//     the Handler reads target.EffectiveDominate() to pick up printed + granted + conditional
//     Dominate in one probe. Damage attributes back to Mauvrion via SourceIndex.
//
// "Attack action card" excludes weapons; the Matches predicate requires both
// TypeRuneblade and the attack-action bit.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var mauvrionSkiesTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

// mauvrionTargetMatches is the shared predicate for Mauvrion's grants: the next Runeblade
// attack action card (weapons don't qualify).
func mauvrionTargetMatches(target *sim.CardState) bool {
	t := target.Card.Types()
	return t.Has(card.TypeRuneblade) && t.IsAttackAction()
}

// mauvrionSkiesPlay applies the go-again grant via look-ahead, registers the ephemeral
// "if hits, create n Runechants" trigger for the same target's fully-resolved attack,
// and emits the chain step (no value contribution; rider damage is credited via the
// trigger handler).
func mauvrionSkiesPlay(s *sim.TurnState, selfState *sim.CardState, source sim.Card, n int) {
	// Go again is static — flip it on the first matching target so its chain-legality check
	// sees the grant before the target's Play runs.
	for _, pc := range s.CardsRemaining {
		if mauvrionTargetMatches(pc) {
			pc.GrantedGoAgain = true
			break
		}
	}
	// The Runechant rider depends on whether the attack actually hits, so defer until the
	// target's full resolution. The trigger fires on the first matching attack action's
	// post-resolution; non-matching attacks (e.g. a Generic attack played before the
	// Runeblade one) leave it in place. Unconsumed triggers fizzle silently at end of turn.
	s.AddEphemeralAttackTrigger(sim.EphemeralAttackTrigger{
		Source:  source,
		Matches: mauvrionTargetMatches,
		Handler: func(s *sim.TurnState, target *sim.CardState) int {
			if !sim.LikelyToHit(target) {
				return 0
			}
			return s.CreateAndLogRunechantsOnHit(sim.DisplayName(source), sim.DisplayName(target.Card), n)
		},
	})
	s.LogPlay(selfState)
}

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
