// Sigil of Cycles — Generic Action - Aura. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "**Go again** At the beginning of your action phase, destroy this. When this leaves the
// arena, discard a card then draw a card."
//
// Simplification: At-action-phase self-destroy and leaves-arena discard/draw aren't modelled; only
// the aura-created flag is credited.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfCyclesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAura)

type SigilOfCyclesBlue struct{}

func (SigilOfCyclesBlue) ID() card.ID                 { return card.SigilOfCyclesBlue }
func (SigilOfCyclesBlue) Name() string                { return "Sigil of Cycles (Blue)" }
func (SigilOfCyclesBlue) Cost(*card.TurnState) int                   { return 0 }
func (SigilOfCyclesBlue) Pitch() int                  { return 3 }
func (SigilOfCyclesBlue) Attack() int                 { return 0 }
func (SigilOfCyclesBlue) Defense() int                { return 2 }
func (SigilOfCyclesBlue) Types() card.TypeSet         { return sigilOfCyclesTypes }
func (SigilOfCyclesBlue) GoAgain() bool               { return true }
func (SigilOfCyclesBlue) Play(s *card.TurnState, _ *card.CardState) int { return setAuraCreated(s) }
