// Smash Up — Generic Action - Attack. Cost 1, Pitch 1, Power 5, Defense 2. Only printed in Red.
//
// Text: "When this hits a hero, turn a card in their arsenal face-up, then banish an attack action
// card from their arsenal."
//
// Simplification: Arsenal-manipulation rider isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var smashUpTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type SmashUpRed struct{}

func (SmashUpRed) ID() card.ID                 { return card.SmashUpRed }
func (SmashUpRed) Name() string                { return "Smash Up (Red)" }
func (SmashUpRed) Cost() int                   { return 1 }
func (SmashUpRed) Pitch() int                  { return 1 }
func (SmashUpRed) Attack() int                 { return 5 }
func (SmashUpRed) Defense() int                { return 2 }
func (SmashUpRed) Types() card.TypeSet         { return smashUpTypes }
func (SmashUpRed) GoAgain() bool               { return false }
func (c SmashUpRed) Play(s *card.TurnState) int { return smashUpDamage(c.Attack()) }

// smashUpDamage is a breadcrumb for the on-hit "arsenal face-up + banish attack action" rider —
// not modelled yet (see TODO.md).
func smashUpDamage(attack int) int {
	if card.LikelyToHit(attack) {
		// TODO: model on-hit arsenal manipulation rider.
	}
	return attack
}
