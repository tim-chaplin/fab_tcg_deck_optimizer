// Sigil of Fyendal — Generic Action - Aura. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "**Go again** At the beginning of your action phase, destroy this. When this leaves the
// arena, gain 1{h}."
//
// Simplification: At-action-phase self-destroy and the 1{h} gain on leave are dropped; only the
// aura-created flag is credited.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfFyendalTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAura)

type SigilOfFyendalBlue struct{}

func (SigilOfFyendalBlue) ID() card.ID                 { return card.SigilOfFyendalBlue }
func (SigilOfFyendalBlue) Name() string                { return "Sigil of Fyendal (Blue)" }
func (SigilOfFyendalBlue) Cost() int                   { return 0 }
func (SigilOfFyendalBlue) Pitch() int                  { return 3 }
func (SigilOfFyendalBlue) Attack() int                 { return 0 }
func (SigilOfFyendalBlue) Defense() int                { return 2 }
func (SigilOfFyendalBlue) Types() card.TypeSet         { return sigilOfFyendalTypes }
func (SigilOfFyendalBlue) GoAgain() bool               { return true }
func (SigilOfFyendalBlue) Play(s *card.TurnState) int { return setAuraCreated(s) }
