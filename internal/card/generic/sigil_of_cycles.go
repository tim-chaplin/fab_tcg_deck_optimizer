// Sigil of Cycles — Generic Action - Aura. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "**Go again** At the beginning of your action phase, destroy this. When this leaves the
// arena, discard a card then draw a card."
//
// Play only flips AuraCreated — the aura stays in play. PlayNextTurn fires when the aura is
// destroyed at the start of the next action phase and hands control back to the deck loop so
// the aura is moved to the graveyard. TODO: credit the "discard a card then draw a card" rider.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfCyclesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAura)

type SigilOfCyclesBlue struct{}

func (SigilOfCyclesBlue) ID() card.ID                 { return card.SigilOfCyclesBlue }
func (SigilOfCyclesBlue) Name() string                { return "Sigil of Cycles (Blue)" }
func (SigilOfCyclesBlue) Cost(*card.TurnState) int    { return 0 }
func (SigilOfCyclesBlue) Pitch() int                  { return 3 }
func (SigilOfCyclesBlue) Attack() int                 { return 0 }
func (SigilOfCyclesBlue) Defense() int                { return 2 }
func (SigilOfCyclesBlue) Types() card.TypeSet         { return sigilOfCyclesTypes }
func (SigilOfCyclesBlue) GoAgain() bool               { return true }
func (SigilOfCyclesBlue) Play(s *card.TurnState) int  { return setAuraCreated(s) }

// PlayNextTurn destroys the aura so it moves to the graveyard next turn. The discard/draw rider
// is not modelled.
func (SigilOfCyclesBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	s.DestroyThis()
	return card.DelayedPlayResult{}
}
