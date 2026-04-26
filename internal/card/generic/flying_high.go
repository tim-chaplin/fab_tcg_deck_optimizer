// Flying High — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "Your next attack this turn gets **go again**. If it's <matching color>, it gets +1{p}.
// **Go again**" (Red checks for a red attack, Yellow for a yellow attack, Blue for a blue attack.)

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var flyingHighTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// flyingHighPlay grants go again to the next attack action card scheduled later this turn.
// If that target's pitch matches matchPitch (this card's own pitch), we also add +1 to its
// BonusAttack — the "+1{p} if it's <matching color>" rider — so EffectiveAttack picks the
// buff up in any LikelyToHit check on the buffed attack. Returns 0; the +1 attributes to
// the target's slot, not Flying High's.
func flyingHighPlay(s *card.TurnState, matchPitch int) int {
	for _, pc := range s.CardsRemaining {
		if !pc.Card.Types().IsAttackAction() {
			continue
		}
		pc.GrantedGoAgain = true
		if pc.Card.Pitch() == matchPitch {
			pc.BonusAttack++
		}
		return 0
	}
	return 0
}

type FlyingHighRed struct{}

func (FlyingHighRed) ID() card.ID                 { return card.FlyingHighRed }
func (FlyingHighRed) Name() string                { return "Flying High" }
func (FlyingHighRed) Cost(*card.TurnState) int                   { return 0 }
func (FlyingHighRed) Pitch() int                  { return 1 }
func (FlyingHighRed) Attack() int                 { return 0 }
func (FlyingHighRed) Defense() int                { return 2 }
func (FlyingHighRed) Types() card.TypeSet         { return flyingHighTypes }
func (FlyingHighRed) GoAgain() bool               { return true }
func (FlyingHighRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, flyingHighPlay(s, 1))
}
type FlyingHighYellow struct{}

func (FlyingHighYellow) ID() card.ID                 { return card.FlyingHighYellow }
func (FlyingHighYellow) Name() string                { return "Flying High" }
func (FlyingHighYellow) Cost(*card.TurnState) int                   { return 0 }
func (FlyingHighYellow) Pitch() int                  { return 2 }
func (FlyingHighYellow) Attack() int                 { return 0 }
func (FlyingHighYellow) Defense() int                { return 2 }
func (FlyingHighYellow) Types() card.TypeSet         { return flyingHighTypes }
func (FlyingHighYellow) GoAgain() bool               { return true }
func (FlyingHighYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, flyingHighPlay(s, 2))
}
type FlyingHighBlue struct{}

func (FlyingHighBlue) ID() card.ID                 { return card.FlyingHighBlue }
func (FlyingHighBlue) Name() string                { return "Flying High" }
func (FlyingHighBlue) Cost(*card.TurnState) int                   { return 0 }
func (FlyingHighBlue) Pitch() int                  { return 3 }
func (FlyingHighBlue) Attack() int                 { return 0 }
func (FlyingHighBlue) Defense() int                { return 2 }
func (FlyingHighBlue) Types() card.TypeSet         { return flyingHighTypes }
func (FlyingHighBlue) GoAgain() bool               { return true }
func (FlyingHighBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, flyingHighPlay(s, 3))
}