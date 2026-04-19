// Look Tuff — Generic Action - Attack. Cost 3, Pitch 1, Power 8, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."
//
// Simplification: Pay {r} or lose 1{p} — base power is kept.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var lookTuffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type LookTuffRed struct{}

func (LookTuffRed) ID() card.ID                 { return card.LookTuffRed }
func (LookTuffRed) Name() string                { return "Look Tuff (Red)" }
func (LookTuffRed) Cost(*card.TurnState) int                   { return 3 }
func (LookTuffRed) Pitch() int                  { return 1 }
func (LookTuffRed) Attack() int                 { return 8 }
func (LookTuffRed) Defense() int                { return 3 }
func (LookTuffRed) Types() card.TypeSet         { return lookTuffTypes }
func (LookTuffRed) GoAgain() bool               { return false }
func (c LookTuffRed) Play(s *card.TurnState) int { return c.Attack() }
