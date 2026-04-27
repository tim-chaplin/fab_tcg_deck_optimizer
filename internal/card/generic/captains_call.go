// Captain's Call — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 2.
//
// Text: "Choose 1; The next attack action card with cost N or less you play this turn gains +2{p}.
// The next attack action card with cost N or less you play this turn gains **go again**. **Go
// again**" (Red N=2, Yellow N=1, Blue N=0.)

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var captainsCallTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// captainsCallApplySideEffect grants +2 to the first scheduled attack action card whose cost is
// at most maxCost, by adding to its BonusAttack. The +2 attributes to the buffed attack (so
// EffectiveAttack picks it up in LikelyToHit) rather than to Captain's Call itself. The
// alternative "go again" mode is dropped (advertised on each variant's NotImplemented marker).
func captainsCallApplySideEffect(s *card.TurnState, maxCost int) {
	for _, pc := range s.CardsRemaining {
		if !pc.Card.Types().IsAttackAction() {
			continue
		}
		if pc.Card.Cost(s) <= maxCost {
			pc.BonusAttack += 2
			return
		}
	}
}

type CaptainsCallRed struct{}

func (CaptainsCallRed) ID() card.ID              { return card.CaptainsCallRed }
func (CaptainsCallRed) Name() string             { return "Captain's Call" }
func (CaptainsCallRed) Cost(*card.TurnState) int { return 0 }
func (CaptainsCallRed) Pitch() int               { return 1 }
func (CaptainsCallRed) Attack() int              { return 0 }
func (CaptainsCallRed) Defense() int             { return 2 }
func (CaptainsCallRed) Types() card.TypeSet      { return captainsCallTypes }
func (CaptainsCallRed) GoAgain() bool            { return true }

// not implemented: modal pick hard-coded to +2{p}; 'go again' mode is dropped
func (CaptainsCallRed) NotImplemented() {}
func (CaptainsCallRed) Play(s *card.TurnState, self *card.CardState) {
	captainsCallApplySideEffect(s, 2)
	s.ApplyAndLogEffectiveAttack(self)
}

type CaptainsCallYellow struct{}

func (CaptainsCallYellow) ID() card.ID              { return card.CaptainsCallYellow }
func (CaptainsCallYellow) Name() string             { return "Captain's Call" }
func (CaptainsCallYellow) Cost(*card.TurnState) int { return 0 }
func (CaptainsCallYellow) Pitch() int               { return 2 }
func (CaptainsCallYellow) Attack() int              { return 0 }
func (CaptainsCallYellow) Defense() int             { return 2 }
func (CaptainsCallYellow) Types() card.TypeSet      { return captainsCallTypes }
func (CaptainsCallYellow) GoAgain() bool            { return true }

// not implemented: modal pick hard-coded to +2{p}; 'go again' mode is dropped
func (CaptainsCallYellow) NotImplemented() {}
func (CaptainsCallYellow) Play(s *card.TurnState, self *card.CardState) {
	captainsCallApplySideEffect(s, 1)
	s.ApplyAndLogEffectiveAttack(self)
}

type CaptainsCallBlue struct{}

func (CaptainsCallBlue) ID() card.ID              { return card.CaptainsCallBlue }
func (CaptainsCallBlue) Name() string             { return "Captain's Call" }
func (CaptainsCallBlue) Cost(*card.TurnState) int { return 0 }
func (CaptainsCallBlue) Pitch() int               { return 3 }
func (CaptainsCallBlue) Attack() int              { return 0 }
func (CaptainsCallBlue) Defense() int             { return 2 }
func (CaptainsCallBlue) Types() card.TypeSet      { return captainsCallTypes }
func (CaptainsCallBlue) GoAgain() bool            { return true }

// not implemented: modal pick hard-coded to +2{p}; 'go again' mode is dropped
func (CaptainsCallBlue) NotImplemented() {}
func (CaptainsCallBlue) Play(s *card.TurnState, self *card.CardState) {
	captainsCallApplySideEffect(s, 0)
	s.ApplyAndLogEffectiveAttack(self)
}
