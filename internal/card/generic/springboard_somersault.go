// Springboard Somersault — Generic Defense Reaction. Cost 0, Pitch 2, Defense 2. Only printed in
// Yellow.
// Text: "If Springboard Somersault is played from arsenal, it gains +2{d}."
//
// Modelling: The +2{d} rider opts in via card.ArsenalDefenseBonus; the solver bumps the
// arsenal slot's effective defense by 2 only when this copy was the start-of-turn arsenal-in
// card.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type SpringboardSomersaultYellow struct{}

func (SpringboardSomersaultYellow) ID() card.ID                 { return card.SpringboardSomersaultYellow }
func (SpringboardSomersaultYellow) Name() string             { return "Springboard Somersault (Yellow)" }
func (SpringboardSomersaultYellow) Cost(*card.TurnState) int                { return 0 }
func (SpringboardSomersaultYellow) Pitch() int               { return 2 }
func (SpringboardSomersaultYellow) Attack() int              { return 0 }
func (SpringboardSomersaultYellow) Defense() int             { return 2 }
func (SpringboardSomersaultYellow) Types() card.TypeSet      { return defenseReactionTypes }
func (SpringboardSomersaultYellow) GoAgain() bool            { return false }
func (SpringboardSomersaultYellow) Play(*card.TurnState, *card.PlayedCard) int { return 0 }
func (SpringboardSomersaultYellow) ArsenalDefenseBonus() int { return 2 }
