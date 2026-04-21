// Blanch — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, cards they own lose all colors until the end of their next turn."
//
// Simplification: Opponent 'lose all colors' debuff isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var blanchTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BlanchRed struct{}

func (BlanchRed) ID() card.ID                 { return card.BlanchRed }
func (BlanchRed) Name() string                { return "Blanch (Red)" }
func (BlanchRed) Cost(*card.TurnState) int                   { return 3 }
func (BlanchRed) Pitch() int                  { return 1 }
func (BlanchRed) Attack() int                 { return 7 }
func (BlanchRed) Defense() int                { return 2 }
func (BlanchRed) Types() card.TypeSet         { return blanchTypes }
func (BlanchRed) GoAgain() bool               { return false }
func (c BlanchRed) Play(s *card.TurnState, _ *card.CardState) int { return blanchDamage(c.Attack()) }

type BlanchYellow struct{}

func (BlanchYellow) ID() card.ID                 { return card.BlanchYellow }
func (BlanchYellow) Name() string                { return "Blanch (Yellow)" }
func (BlanchYellow) Cost(*card.TurnState) int                   { return 3 }
func (BlanchYellow) Pitch() int                  { return 2 }
func (BlanchYellow) Attack() int                 { return 6 }
func (BlanchYellow) Defense() int                { return 2 }
func (BlanchYellow) Types() card.TypeSet         { return blanchTypes }
func (BlanchYellow) GoAgain() bool               { return false }
func (c BlanchYellow) Play(s *card.TurnState, _ *card.CardState) int { return blanchDamage(c.Attack()) }

type BlanchBlue struct{}

func (BlanchBlue) ID() card.ID                 { return card.BlanchBlue }
func (BlanchBlue) Name() string                { return "Blanch (Blue)" }
func (BlanchBlue) Cost(*card.TurnState) int                   { return 3 }
func (BlanchBlue) Pitch() int                  { return 3 }
func (BlanchBlue) Attack() int                 { return 5 }
func (BlanchBlue) Defense() int                { return 2 }
func (BlanchBlue) Types() card.TypeSet         { return blanchTypes }
func (BlanchBlue) GoAgain() bool               { return false }
func (c BlanchBlue) Play(s *card.TurnState, _ *card.CardState) int { return blanchDamage(c.Attack()) }

// blanchDamage is a breadcrumb for the on-hit "cards they own lose all colors" rider — not
// modelled yet (see TODO.md).
func blanchDamage(attack int) int {
	if card.LikelyToHit(attack) {
		// TODO: model on-hit opponent-card color-strip rider.
	}
	return attack
}
