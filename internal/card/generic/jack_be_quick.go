// Jack Be Quick — Generic Action - Attack. Cost 0, Pitch 1, Power 3, Defense 3. Only printed in
// Red.
//
// Text: "When this attacks, you may banish a Nimblism from your graveyard. If you do, this gets
// +1{p} and **go again**. When this hits a hero, {u} an ally they control, then steal it until the
// end of this action phase."
//
// Simplification: Graveyard banish for +1{p}/go-again and on-hit steal aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var jackBeQuickTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type JackBeQuickRed struct{}

func (JackBeQuickRed) ID() card.ID                 { return card.JackBeQuickRed }
func (JackBeQuickRed) Name() string                { return "Jack Be Quick (Red)" }
func (JackBeQuickRed) Cost() int                   { return 0 }
func (JackBeQuickRed) Pitch() int                  { return 1 }
func (JackBeQuickRed) Attack() int                 { return 3 }
func (JackBeQuickRed) Defense() int                { return 3 }
func (JackBeQuickRed) Types() card.TypeSet         { return jackBeQuickTypes }
func (JackBeQuickRed) GoAgain() bool               { return false }
func (c JackBeQuickRed) Play(s *card.TurnState) int { return c.Attack() }
