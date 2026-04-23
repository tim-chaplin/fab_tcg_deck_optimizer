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

// waterTheSeedsPlay returns basePower plus 1 if the next attack action card with base {p} <= 1 is
// scheduled later this turn, otherwise just basePower.
func waterTheSeedsPlay(basePower int, s *card.TurnState) int {
	for _, pc := range s.CardsRemaining {
		if !pc.Card.Types().IsAttackAction() {
			continue
		}
		if pc.Card.Attack() <= 1 {
			return basePower + 1
		}
	}
	return basePower
}

type WaterTheSeedsRed struct{}

func (WaterTheSeedsRed) ID() card.ID                 { return card.WaterTheSeedsRed }
func (WaterTheSeedsRed) Name() string                { return "Water the Seeds (Red)" }
func (WaterTheSeedsRed) Cost(*card.TurnState) int                   { return 1 }
func (WaterTheSeedsRed) Pitch() int                  { return 1 }
func (WaterTheSeedsRed) Attack() int                 { return 3 }
func (WaterTheSeedsRed) Defense() int                { return 2 }
func (WaterTheSeedsRed) Types() card.TypeSet         { return waterTheSeedsTypes }
func (WaterTheSeedsRed) GoAgain() bool               { return true }
func (c WaterTheSeedsRed) Play(s *card.TurnState, _ *card.CardState) int { return waterTheSeedsPlay(c.Attack(), s) }

type WaterTheSeedsYellow struct{}

func (WaterTheSeedsYellow) ID() card.ID                 { return card.WaterTheSeedsYellow }
func (WaterTheSeedsYellow) Name() string                { return "Water the Seeds (Yellow)" }
func (WaterTheSeedsYellow) Cost(*card.TurnState) int                   { return 1 }
func (WaterTheSeedsYellow) Pitch() int                  { return 2 }
func (WaterTheSeedsYellow) Attack() int                 { return 2 }
func (WaterTheSeedsYellow) Defense() int                { return 2 }
func (WaterTheSeedsYellow) Types() card.TypeSet         { return waterTheSeedsTypes }
func (WaterTheSeedsYellow) GoAgain() bool               { return true }
func (c WaterTheSeedsYellow) Play(s *card.TurnState, _ *card.CardState) int { return waterTheSeedsPlay(c.Attack(), s) }

type WaterTheSeedsBlue struct{}

func (WaterTheSeedsBlue) ID() card.ID                 { return card.WaterTheSeedsBlue }
func (WaterTheSeedsBlue) Name() string                { return "Water the Seeds (Blue)" }
func (WaterTheSeedsBlue) Cost(*card.TurnState) int                   { return 1 }
func (WaterTheSeedsBlue) Pitch() int                  { return 3 }
func (WaterTheSeedsBlue) Attack() int                 { return 1 }
func (WaterTheSeedsBlue) Defense() int                { return 2 }
func (WaterTheSeedsBlue) Types() card.TypeSet         { return waterTheSeedsTypes }
func (WaterTheSeedsBlue) GoAgain() bool               { return true }
func (c WaterTheSeedsBlue) Play(s *card.TurnState, _ *card.CardState) int { return waterTheSeedsPlay(c.Attack(), s) }
