// Sigil of Fyendal — Generic Action - Aura. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "**Go again** At the beginning of your action phase, destroy this. When this leaves the
// arena, gain 1{h}."
//
// Modelling: the sigil self-destroys at the start of your next action phase, which always
// fires on the cadence the card prints, so the 1{h} "on leave" gain is credited as +1 damage
// equivalent on Play (health is valued 1-to-1 with damage). The aura-created flag is also
// set so same-turn "if you've created an aura" riders see it.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfFyendalTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAura)

type SigilOfFyendalBlue struct{}

func (SigilOfFyendalBlue) ID() card.ID                 { return card.SigilOfFyendalBlue }
func (SigilOfFyendalBlue) Name() string                { return "Sigil of Fyendal (Blue)" }
func (SigilOfFyendalBlue) Cost(*card.TurnState) int                   { return 0 }
func (SigilOfFyendalBlue) Pitch() int                  { return 3 }
func (SigilOfFyendalBlue) Attack() int                 { return 0 }
func (SigilOfFyendalBlue) Defense() int                { return 2 }
func (SigilOfFyendalBlue) Types() card.TypeSet         { return sigilOfFyendalTypes }
func (SigilOfFyendalBlue) GoAgain() bool               { return true }
func (SigilOfFyendalBlue) Play(s *card.TurnState) int { return setAuraCreated(s) + 1 }
