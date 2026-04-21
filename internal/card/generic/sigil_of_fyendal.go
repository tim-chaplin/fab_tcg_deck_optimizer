// Sigil of Fyendal — Generic Action - Aura. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "**Go again** At the beginning of your action phase, destroy this. When this leaves the
// arena, gain 1{h}."
//
// Modelling: Play only flips the aura-created flag so same-turn aura-readers see it. The
// self-destroy and 1{h} gain fire at the start of the NEXT action phase via card.DelayedPlay;
// PlayNextTurn credits +1 damage-equivalent (health valued 1-to-1 with damage).
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfFyendalTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAura)

type SigilOfFyendalBlue struct{}

func (SigilOfFyendalBlue) ID() card.ID         { return card.SigilOfFyendalBlue }
func (SigilOfFyendalBlue) Name() string        { return "Sigil of Fyendal (Blue)" }
func (SigilOfFyendalBlue) Cost(*card.TurnState) int { return 0 }
func (SigilOfFyendalBlue) Pitch() int          { return 3 }
func (SigilOfFyendalBlue) Attack() int         { return 0 }
func (SigilOfFyendalBlue) Defense() int        { return 2 }
func (SigilOfFyendalBlue) Types() card.TypeSet { return sigilOfFyendalTypes }
func (SigilOfFyendalBlue) GoAgain() bool       { return true }
func (SigilOfFyendalBlue) Play(s *card.TurnState) int { return setAuraCreated(s) }

// PlayNextTurn credits the 1{h} gain that fires when the aura leaves the arena at the start of
// the next action phase.
func (c SigilOfFyendalBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	s.AddToGraveyard(c)
	return card.DelayedPlayResult{Damage: 1}
}
