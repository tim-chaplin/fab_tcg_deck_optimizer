// Water the Seeds — Generic Action - Attack. Cost 1. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, your next attack this combat chain with 1 or less base {p} gets +1{p}.
// **Go again**"
//
// Scans TurnState.CardsRemaining for the first attack action card with base power 1 or less and
// credits the +1 assuming it will be played; if no matching attack follows, the rider fizzles.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var waterTheSeedsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// grantWaterTheSeedsBonus adds +1 to the first scheduled attack action card with base {p}
// of 1 or less — the buff lands on that card's BonusAttack so EffectiveAttack and
// LikelyToHit see the buffed power on the buffed card, not on the granter. Fizzles
// silently when no qualifying target follows.
func grantWaterTheSeedsBonus(s *card.TurnState) {
	for _, pc := range s.CardsRemaining {
		if !pc.Card.Types().IsAttackAction() {
			continue
		}
		if pc.Card.Attack() <= 1 {
			pc.BonusAttack += 1
			return
		}
	}
}

type WaterTheSeedsRed struct{}

func (WaterTheSeedsRed) ID() card.ID              { return card.WaterTheSeedsRed }
func (WaterTheSeedsRed) Name() string             { return "Water the Seeds" }
func (WaterTheSeedsRed) Cost(*card.TurnState) int { return 1 }
func (WaterTheSeedsRed) Pitch() int               { return 1 }
func (WaterTheSeedsRed) Attack() int              { return 3 }
func (WaterTheSeedsRed) Defense() int             { return 2 }
func (WaterTheSeedsRed) Types() card.TypeSet      { return waterTheSeedsTypes }
func (WaterTheSeedsRed) GoAgain() bool            { return true }
func (WaterTheSeedsRed) Play(s *card.TurnState, self *card.CardState) {
	grantWaterTheSeedsBonus(s)
	s.ApplyAndLogEffectiveAttack(self)
}

type WaterTheSeedsYellow struct{}

func (WaterTheSeedsYellow) ID() card.ID              { return card.WaterTheSeedsYellow }
func (WaterTheSeedsYellow) Name() string             { return "Water the Seeds" }
func (WaterTheSeedsYellow) Cost(*card.TurnState) int { return 1 }
func (WaterTheSeedsYellow) Pitch() int               { return 2 }
func (WaterTheSeedsYellow) Attack() int              { return 2 }
func (WaterTheSeedsYellow) Defense() int             { return 2 }
func (WaterTheSeedsYellow) Types() card.TypeSet      { return waterTheSeedsTypes }
func (WaterTheSeedsYellow) GoAgain() bool            { return true }
func (WaterTheSeedsYellow) Play(s *card.TurnState, self *card.CardState) {
	grantWaterTheSeedsBonus(s)
	s.ApplyAndLogEffectiveAttack(self)
}

type WaterTheSeedsBlue struct{}

func (WaterTheSeedsBlue) ID() card.ID              { return card.WaterTheSeedsBlue }
func (WaterTheSeedsBlue) Name() string             { return "Water the Seeds" }
func (WaterTheSeedsBlue) Cost(*card.TurnState) int { return 1 }
func (WaterTheSeedsBlue) Pitch() int               { return 3 }
func (WaterTheSeedsBlue) Attack() int              { return 1 }
func (WaterTheSeedsBlue) Defense() int             { return 2 }
func (WaterTheSeedsBlue) Types() card.TypeSet      { return waterTheSeedsTypes }
func (WaterTheSeedsBlue) GoAgain() bool            { return true }
func (WaterTheSeedsBlue) Play(s *card.TurnState, self *card.CardState) {
	grantWaterTheSeedsBonus(s)
	s.ApplyAndLogEffectiveAttack(self)
}
