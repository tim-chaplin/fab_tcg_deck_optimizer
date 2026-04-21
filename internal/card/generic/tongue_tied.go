// Tongue Tied — Generic Action - Attack. Cost 3, Pitch 1, Power 7, Defense 2. Only printed in Red.
//
// Text: "When this hits a hero, turn a card in their arsenal face-up, then banish an instant card
// from their arsenal."
//
// Simplification: Arsenal-manipulation rider isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var tongueTiedTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type TongueTiedRed struct{}

func (TongueTiedRed) ID() card.ID                 { return card.TongueTiedRed }
func (TongueTiedRed) Name() string                { return "Tongue Tied (Red)" }
func (TongueTiedRed) Cost(*card.TurnState) int                   { return 3 }
func (TongueTiedRed) Pitch() int                  { return 1 }
func (TongueTiedRed) Attack() int                 { return 7 }
func (TongueTiedRed) Defense() int                { return 2 }
func (TongueTiedRed) Types() card.TypeSet         { return tongueTiedTypes }
func (TongueTiedRed) GoAgain() bool               { return false }
func (c TongueTiedRed) Play(s *card.TurnState, _ *card.CardState) int { return tongueTiedDamage(c.Attack()) }

// tongueTiedDamage is a breadcrumb for the on-hit "arsenal face-up + banish instant" rider —
// not modelled yet (see TODO.md).
func tongueTiedDamage(attack int) int {
	if card.LikelyToHit(attack) {
		// TODO: model on-hit arsenal manipulation rider.
	}
	return attack
}
