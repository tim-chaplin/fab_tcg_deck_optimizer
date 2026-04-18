// Chest Puff — Generic Action - Attack. Cost 2, Pitch 1, Power 7, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."
//
// Simplification: Pay {r} or lose 1{p} — base power is kept.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var chestPuffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type ChestPuffRed struct{}

func (ChestPuffRed) ID() card.ID                 { return card.ChestPuffRed }
func (ChestPuffRed) Name() string                { return "Chest Puff (Red)" }
func (ChestPuffRed) Cost() int                   { return 2 }
func (ChestPuffRed) Pitch() int                  { return 1 }
func (ChestPuffRed) Attack() int                 { return 7 }
func (ChestPuffRed) Defense() int                { return 3 }
func (ChestPuffRed) Types() card.TypeSet         { return chestPuffTypes }
func (ChestPuffRed) GoAgain() bool               { return false }
func (c ChestPuffRed) Play(s *card.TurnState) int { return c.Attack() }
