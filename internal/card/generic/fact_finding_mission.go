// Fact-Finding Mission — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, you may look at a face-down card in their arsenal or equipment
// zones."
//
// Simplification: Peeking opponent arsenal/equipment isn't modelled.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var factFindingMissionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type FactFindingMissionRed struct{}

func (FactFindingMissionRed) ID() card.ID                 { return card.FactFindingMissionRed }
func (FactFindingMissionRed) Name() string                { return "Fact-Finding Mission (Red)" }
func (FactFindingMissionRed) Cost(*card.TurnState) int                   { return 2 }
func (FactFindingMissionRed) Pitch() int                  { return 1 }
func (FactFindingMissionRed) Attack() int                 { return 6 }
func (FactFindingMissionRed) Defense() int                { return 2 }
func (FactFindingMissionRed) Types() card.TypeSet         { return factFindingMissionTypes }
func (FactFindingMissionRed) GoAgain() bool               { return false }
func (c FactFindingMissionRed) Play(s *card.TurnState, _ *card.CardState) int { return factFindingMissionDamage(c.Attack()) }

type FactFindingMissionYellow struct{}

func (FactFindingMissionYellow) ID() card.ID                 { return card.FactFindingMissionYellow }
func (FactFindingMissionYellow) Name() string                { return "Fact-Finding Mission (Yellow)" }
func (FactFindingMissionYellow) Cost(*card.TurnState) int                   { return 2 }
func (FactFindingMissionYellow) Pitch() int                  { return 2 }
func (FactFindingMissionYellow) Attack() int                 { return 5 }
func (FactFindingMissionYellow) Defense() int                { return 2 }
func (FactFindingMissionYellow) Types() card.TypeSet         { return factFindingMissionTypes }
func (FactFindingMissionYellow) GoAgain() bool               { return false }
func (c FactFindingMissionYellow) Play(s *card.TurnState, _ *card.CardState) int { return factFindingMissionDamage(c.Attack()) }

type FactFindingMissionBlue struct{}

func (FactFindingMissionBlue) ID() card.ID                 { return card.FactFindingMissionBlue }
func (FactFindingMissionBlue) Name() string                { return "Fact-Finding Mission (Blue)" }
func (FactFindingMissionBlue) Cost(*card.TurnState) int                   { return 2 }
func (FactFindingMissionBlue) Pitch() int                  { return 3 }
func (FactFindingMissionBlue) Attack() int                 { return 4 }
func (FactFindingMissionBlue) Defense() int                { return 2 }
func (FactFindingMissionBlue) Types() card.TypeSet         { return factFindingMissionTypes }
func (FactFindingMissionBlue) GoAgain() bool               { return false }
func (c FactFindingMissionBlue) Play(s *card.TurnState, _ *card.CardState) int { return factFindingMissionDamage(c.Attack()) }

// factFindingMissionDamage is a breadcrumb for the on-hit "peek a face-down card in arsenal /
// equipment" rider — opponent-side inspection isn't modelled (see TODO.md).
func factFindingMissionDamage(attack int) int {
	if card.LikelyToHit(attack, false) {
		// TODO: model on-hit opponent-arsenal peek rider.
	}
	return attack
}
