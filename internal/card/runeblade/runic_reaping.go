// Runic Reaping — Runeblade Action. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "The next Runeblade attack action card you play this turn gains 'When this hits, create N
// Runechant tokens'. If an attack card was pitched to play Runic Reaping, the next Runeblade
// attack action card you play this turn gains +1{p}. Go again."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling: The N Runechants are credited only when the target's printed Attack() satisfies
// card.LikelyToHit — mirrors Mauvrion Skies' gate, since both riders are on-hit clauses that
// fizzle when the target gets blocked. The +1{p} pitched-attack rider isn't gated on hitting
// and fires if any attack-typed card was pitched this turn; pitch-to-play attribution isn't
// tracked. Both riders target only Runeblade attack action cards — weapon swings don't qualify.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runicReapingTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

type RunicReapingRed struct{}

func (RunicReapingRed) ID() card.ID                 { return card.RunicReapingRed }
func (RunicReapingRed) Name() string               { return "Runic Reaping (Red)" }
func (RunicReapingRed) Cost(*card.TurnState) int                  { return 1 }
func (RunicReapingRed) Pitch() int                 { return 1 }
func (RunicReapingRed) Attack() int                { return 0 }
func (RunicReapingRed) Defense() int               { return 2 }
func (RunicReapingRed) Types() card.TypeSet        { return runicReapingTypes }
func (RunicReapingRed) GoAgain() bool              { return true }
func (RunicReapingRed) Play(s *card.TurnState, _ *card.CardState) int { return runicReapingPlay(s, 3) }

type RunicReapingYellow struct{}

func (RunicReapingYellow) ID() card.ID                 { return card.RunicReapingYellow }
func (RunicReapingYellow) Name() string               { return "Runic Reaping (Yellow)" }
func (RunicReapingYellow) Cost(*card.TurnState) int                  { return 1 }
func (RunicReapingYellow) Pitch() int                 { return 2 }
func (RunicReapingYellow) Attack() int                { return 0 }
func (RunicReapingYellow) Defense() int               { return 2 }
func (RunicReapingYellow) Types() card.TypeSet        { return runicReapingTypes }
func (RunicReapingYellow) GoAgain() bool              { return true }
func (RunicReapingYellow) Play(s *card.TurnState, _ *card.CardState) int { return runicReapingPlay(s, 2) }

type RunicReapingBlue struct{}

func (RunicReapingBlue) ID() card.ID                 { return card.RunicReapingBlue }
func (RunicReapingBlue) Name() string               { return "Runic Reaping (Blue)" }
func (RunicReapingBlue) Cost(*card.TurnState) int                  { return 1 }
func (RunicReapingBlue) Pitch() int                 { return 3 }
func (RunicReapingBlue) Attack() int                { return 0 }
func (RunicReapingBlue) Defense() int               { return 2 }
func (RunicReapingBlue) Types() card.TypeSet        { return runicReapingTypes }
func (RunicReapingBlue) GoAgain() bool              { return true }
func (RunicReapingBlue) Play(s *card.TurnState, _ *card.CardState) int { return runicReapingPlay(s, 1) }

func runicReapingPlay(s *card.TurnState, n int) int {
	var target card.Card
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if t.Has(card.TypeRuneblade) && t.IsAttackAction() {
			target = pc.Card
			break
		}
	}
	if target == nil {
		return 0
	}
	bonus := 0
	for _, p := range s.Pitched {
		if p.Types().Has(card.TypeAttack) {
			bonus = 1
			break
		}
	}
	if card.LikelyToHit(target.Attack()) {
		return s.CreateRunechants(n) + bonus
	}
	return bonus
}
