// Runic Reaping — Runeblade Action. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "The next Runeblade attack action card you play this turn gains 'When this hits, create N
// Runechant tokens'. If an attack card was pitched to play Runic Reaping, the next Runeblade
// attack action card you play this turn gains +1{p}. Go again."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Both riders target Runeblade attack action cards only — weapon swings don't qualify.
//
// Modelling splits the two riders by when they need to resolve:
//   - The +1{p} pitched-attack rider is a static buff on the target's printed power, so it's
//     applied via pc.BonusAttack on the look-ahead pass before the target plays. EffectiveAttack
//     folds it in and a card playing between Runic Reaping and the target sees the buff if it
//     scans target.BonusAttack.
//   - The "if this hits, create N Runechants" rider depends on the target's fully-resolved
//     attack state. A card that plays between Runic Reaping and the target may grant more
//     BonusAttack (or Dominate, etc.), and LikelyToHit needs to see those grants. Play
//     registers an EphemeralAttackTrigger; the handler runs after the target's Play and
//     reads target.EffectiveAttack / target.EffectiveDominate.
//
// Pitch-to-play attribution isn't tracked: any attack-typed card in Pitched satisfies the
// +1{p} rider.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var runicReapingTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

// runicReapingTargetMatches is the shared predicate for Runic Reaping's two riders: the next
// Runeblade attack action card (weapons don't qualify).
func runicReapingTargetMatches(target *card.CardState) bool {
	t := target.Card.Types()
	return t.Has(card.TypeRuneblade) && t.IsAttackAction()
}

type RunicReapingRed struct{}

func (RunicReapingRed) ID() ids.CardID           { return ids.RunicReapingRed }
func (RunicReapingRed) Name() string             { return "Runic Reaping" }
func (RunicReapingRed) Cost(*card.TurnState) int { return 1 }
func (RunicReapingRed) Pitch() int               { return 1 }
func (RunicReapingRed) Attack() int              { return 0 }
func (RunicReapingRed) Defense() int             { return 2 }
func (RunicReapingRed) Types() card.TypeSet      { return runicReapingTypes }
func (RunicReapingRed) GoAgain() bool            { return true }
func (c RunicReapingRed) Play(s *card.TurnState, self *card.CardState) {
	runicReapingPlay(s, self, c, 3)
}

type RunicReapingYellow struct{}

func (RunicReapingYellow) ID() ids.CardID           { return ids.RunicReapingYellow }
func (RunicReapingYellow) Name() string             { return "Runic Reaping" }
func (RunicReapingYellow) Cost(*card.TurnState) int { return 1 }
func (RunicReapingYellow) Pitch() int               { return 2 }
func (RunicReapingYellow) Attack() int              { return 0 }
func (RunicReapingYellow) Defense() int             { return 2 }
func (RunicReapingYellow) Types() card.TypeSet      { return runicReapingTypes }
func (RunicReapingYellow) GoAgain() bool            { return true }
func (c RunicReapingYellow) Play(s *card.TurnState, self *card.CardState) {
	runicReapingPlay(s, self, c, 2)
}

type RunicReapingBlue struct{}

func (RunicReapingBlue) ID() ids.CardID           { return ids.RunicReapingBlue }
func (RunicReapingBlue) Name() string             { return "Runic Reaping" }
func (RunicReapingBlue) Cost(*card.TurnState) int { return 1 }
func (RunicReapingBlue) Pitch() int               { return 3 }
func (RunicReapingBlue) Attack() int              { return 0 }
func (RunicReapingBlue) Defense() int             { return 2 }
func (RunicReapingBlue) Types() card.TypeSet      { return runicReapingTypes }
func (RunicReapingBlue) GoAgain() bool            { return true }
func (c RunicReapingBlue) Play(s *card.TurnState, self *card.CardState) {
	runicReapingPlay(s, self, c, 1)
}

// runicReapingPlay applies the pitched-attack +1{p} grant via the target's BonusAttack,
// registers the on-hit Runechant trigger, and emits Runic Reaping's chain step (no
// value contribution — Runic Reaping itself doesn't deal damage; the trigger handler
// credits whatever runechants get created after the target resolves).
func runicReapingPlay(s *card.TurnState, selfState *card.CardState, source card.Card, n int) {
	var target *card.CardState
	for _, pc := range s.CardsRemaining {
		if runicReapingTargetMatches(pc) {
			target = pc
			break
		}
	}
	if target == nil {
		s.LogPlay(selfState)
		return
	}
	for _, p := range s.Pitched {
		if p.Types().Has(card.TypeAttack) {
			target.BonusAttack++
			break
		}
	}
	s.AddEphemeralAttackTrigger(card.EphemeralAttackTrigger{
		Source:  source,
		Matches: runicReapingTargetMatches,
		Handler: func(s *card.TurnState, target *card.CardState) int {
			if !card.LikelyToHit(target) {
				return 0
			}
			return s.CreateAndLogRunechantsOnHit(card.DisplayName(source), card.DisplayName(target.Card), n)
		},
	})
	s.LogPlay(selfState)
}
