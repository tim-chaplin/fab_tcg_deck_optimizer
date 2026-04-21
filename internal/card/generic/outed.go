// Outed — Generic Action - Attack. Cost 0, Pitch 1, Power 3, Defense 0. Only printed in Red.
//
// Text: "If you are **marked**, you can't play this. If the defending hero is **marked**, this gets
// +1{p}. **Go again**"
//
// Simplification: Marked-hero checks aren't modelled; the +1{p} rider never applies.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var outedTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type OutedRed struct{}

func (OutedRed) ID() card.ID                 { return card.OutedRed }
func (OutedRed) Name() string                { return "Outed (Red)" }
func (OutedRed) Cost(*card.TurnState) int                   { return 0 }
func (OutedRed) Pitch() int                  { return 1 }
func (OutedRed) Attack() int                 { return 3 }
func (OutedRed) Defense() int                { return 0 }
func (OutedRed) Types() card.TypeSet         { return outedTypes }
func (OutedRed) GoAgain() bool               { return true }
func (c OutedRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
