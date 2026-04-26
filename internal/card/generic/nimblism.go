// Nimblism — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card with cost 1 or less you play this turn gains +N{p}. **Go
// again**" (Red N=3, Yellow N=2, Blue N=1.)

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var nimblismTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// nimblismPlay grants +n to the first scheduled attack action card whose cost is 1 or less,
// by adding to its BonusAttack. Returns 0 — the +n attributes to the buffed attack (so
// EffectiveAttack picks it up in LikelyToHit) rather than to Nimblism itself.
func nimblismPlay(s *card.TurnState, n int) int {
	for _, pc := range s.CardsRemaining {
		if !pc.Card.Types().IsAttackAction() {
			continue
		}
		if pc.Card.Cost(s) <= 1 {
			pc.BonusAttack += n
			return 0
		}
	}
	return 0
}

type NimblismRed struct{}

func (NimblismRed) ID() card.ID                 { return card.NimblismRed }
func (NimblismRed) Name() string                { return "Nimblism" }
func (NimblismRed) Cost(*card.TurnState) int                   { return 0 }
func (NimblismRed) Pitch() int                  { return 1 }
func (NimblismRed) Attack() int                 { return 0 }
func (NimblismRed) Defense() int                { return 2 }
func (NimblismRed) Types() card.TypeSet         { return nimblismTypes }
func (NimblismRed) GoAgain() bool               { return true }
func (NimblismRed) Play(s *card.TurnState, _ *card.CardState) int { return nimblismPlay(s, 3) }

type NimblismYellow struct{}

func (NimblismYellow) ID() card.ID                 { return card.NimblismYellow }
func (NimblismYellow) Name() string                { return "Nimblism" }
func (NimblismYellow) Cost(*card.TurnState) int                   { return 0 }
func (NimblismYellow) Pitch() int                  { return 2 }
func (NimblismYellow) Attack() int                 { return 0 }
func (NimblismYellow) Defense() int                { return 2 }
func (NimblismYellow) Types() card.TypeSet         { return nimblismTypes }
func (NimblismYellow) GoAgain() bool               { return true }
func (NimblismYellow) Play(s *card.TurnState, _ *card.CardState) int { return nimblismPlay(s, 2) }

type NimblismBlue struct{}

func (NimblismBlue) ID() card.ID                 { return card.NimblismBlue }
func (NimblismBlue) Name() string                { return "Nimblism" }
func (NimblismBlue) Cost(*card.TurnState) int                   { return 0 }
func (NimblismBlue) Pitch() int                  { return 3 }
func (NimblismBlue) Attack() int                 { return 0 }
func (NimblismBlue) Defense() int                { return 2 }
func (NimblismBlue) Types() card.TypeSet         { return nimblismTypes }
func (NimblismBlue) GoAgain() bool               { return true }
func (NimblismBlue) Play(s *card.TurnState, _ *card.CardState) int { return nimblismPlay(s, 1) }
