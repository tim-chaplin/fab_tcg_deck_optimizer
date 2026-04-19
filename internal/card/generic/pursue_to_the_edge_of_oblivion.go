// Pursue to the Edge of Oblivion — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 3.
// Only printed in Red.
//
// Text: "When this hits a hero, **mark** them."
//
// Simplification: Mark on hit isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var pursueToTheEdgeOfOblivionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type PursueToTheEdgeOfOblivionRed struct{}

func (PursueToTheEdgeOfOblivionRed) ID() card.ID                 { return card.PursueToTheEdgeOfOblivionRed }
func (PursueToTheEdgeOfOblivionRed) Name() string                { return "Pursue to the Edge of Oblivion (Red)" }
func (PursueToTheEdgeOfOblivionRed) Cost() int                   { return 0 }
func (PursueToTheEdgeOfOblivionRed) Pitch() int                  { return 1 }
func (PursueToTheEdgeOfOblivionRed) Attack() int                 { return 4 }
func (PursueToTheEdgeOfOblivionRed) Defense() int                { return 3 }
func (PursueToTheEdgeOfOblivionRed) Types() card.TypeSet         { return pursueToTheEdgeOfOblivionTypes }
func (PursueToTheEdgeOfOblivionRed) GoAgain() bool               { return false }
func (c PursueToTheEdgeOfOblivionRed) Play(s *card.TurnState) int { return pursueToTheEdgeOfOblivionDamage(c.Attack()) }

// pursueToTheEdgeOfOblivionDamage is a breadcrumb for the on-hit "mark the hero" rider — marks
// aren't tracked (see TODO.md).
func pursueToTheEdgeOfOblivionDamage(attack int) int {
	if card.LikelyToHit(attack) {
		// TODO: model on-hit mark rider.
	}
	return attack
}
