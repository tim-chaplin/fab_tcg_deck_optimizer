// Sigil of Protection — Generic Action - Aura. Cost 1. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Defense 2.
//
// Text: "**Ward N** At the beginning of your action phase, destroy Sigil of Protection."
// (Red Ward 4, Yellow Ward 3, Blue Ward 2.)
//
// Simplification: Ward N isn't modelled (opponent's damage prevention); only the aura-created
// flag is credited. PlayNextTurn destroys the aura at the start of the next action phase so it
// moves to the graveyard — matching the printed text.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfProtectionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAura)

type SigilOfProtectionRed struct{}

func (SigilOfProtectionRed) ID() card.ID                 { return card.SigilOfProtectionRed }
func (SigilOfProtectionRed) Name() string                { return "Sigil of Protection (Red)" }
func (SigilOfProtectionRed) Cost(*card.TurnState) int    { return 1 }
func (SigilOfProtectionRed) Pitch() int                  { return 1 }
func (SigilOfProtectionRed) Attack() int                 { return 0 }
func (SigilOfProtectionRed) Defense() int                { return 2 }
func (SigilOfProtectionRed) Types() card.TypeSet         { return sigilOfProtectionTypes }
func (SigilOfProtectionRed) GoAgain() bool               { return false }
func (SigilOfProtectionRed) Play(s *card.TurnState, _ *card.PlayedCard) int  { return setAuraCreated(s) }
func (c SigilOfProtectionRed) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	s.AddToGraveyard(c)
	return card.DelayedPlayResult{}
}

type SigilOfProtectionYellow struct{}

func (SigilOfProtectionYellow) ID() card.ID                { return card.SigilOfProtectionYellow }
func (SigilOfProtectionYellow) Name() string               { return "Sigil of Protection (Yellow)" }
func (SigilOfProtectionYellow) Cost(*card.TurnState) int   { return 1 }
func (SigilOfProtectionYellow) Pitch() int                 { return 2 }
func (SigilOfProtectionYellow) Attack() int                { return 0 }
func (SigilOfProtectionYellow) Defense() int               { return 2 }
func (SigilOfProtectionYellow) Types() card.TypeSet        { return sigilOfProtectionTypes }
func (SigilOfProtectionYellow) GoAgain() bool              { return false }
func (SigilOfProtectionYellow) Play(s *card.TurnState, _ *card.PlayedCard) int { return setAuraCreated(s) }
func (c SigilOfProtectionYellow) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	s.AddToGraveyard(c)
	return card.DelayedPlayResult{}
}

type SigilOfProtectionBlue struct{}

func (SigilOfProtectionBlue) ID() card.ID                { return card.SigilOfProtectionBlue }
func (SigilOfProtectionBlue) Name() string               { return "Sigil of Protection (Blue)" }
func (SigilOfProtectionBlue) Cost(*card.TurnState) int   { return 1 }
func (SigilOfProtectionBlue) Pitch() int                 { return 3 }
func (SigilOfProtectionBlue) Attack() int                { return 0 }
func (SigilOfProtectionBlue) Defense() int               { return 2 }
func (SigilOfProtectionBlue) Types() card.TypeSet        { return sigilOfProtectionTypes }
func (SigilOfProtectionBlue) GoAgain() bool              { return false }
func (SigilOfProtectionBlue) Play(s *card.TurnState, _ *card.PlayedCard) int { return setAuraCreated(s) }
func (c SigilOfProtectionBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	s.AddToGraveyard(c)
	return card.DelayedPlayResult{}
}
