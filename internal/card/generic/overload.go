// Overload — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Dominate** If Overload hits, it gains **go again**."
//
// Modelling: Dominate is advertised via the card.Dominator marker so LikelyToHit credits the
// "defender capped at one blocker" bump at 5+ power. The on-hit go-again rider isn't
// modelled — Overload's printed power tops out at 3, so go-again-on-hit would only matter
// once Dominate lets it land, and we don't chain that through today.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var overloadTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type OverloadRed struct{}

func (OverloadRed) ID() card.ID                 { return card.OverloadRed }
func (OverloadRed) Name() string                { return "Overload (Red)" }
func (OverloadRed) Cost(*card.TurnState) int                   { return 0 }
func (OverloadRed) Pitch() int                  { return 1 }
func (OverloadRed) Attack() int                 { return 3 }
func (OverloadRed) Defense() int                { return 2 }
func (OverloadRed) Types() card.TypeSet         { return overloadTypes }
func (OverloadRed) GoAgain() bool               { return false }
func (OverloadRed) Dominate()                   {}
// not implemented: on-hit go-again rider (Overload's printed power tops out at 3 so it'd only matter once Dominate lets it land)
func (OverloadRed) NotImplemented()             {}
func (c OverloadRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type OverloadYellow struct{}

func (OverloadYellow) ID() card.ID                 { return card.OverloadYellow }
func (OverloadYellow) Name() string                { return "Overload (Yellow)" }
func (OverloadYellow) Cost(*card.TurnState) int                   { return 0 }
func (OverloadYellow) Pitch() int                  { return 2 }
func (OverloadYellow) Attack() int                 { return 2 }
func (OverloadYellow) Defense() int                { return 2 }
func (OverloadYellow) Types() card.TypeSet         { return overloadTypes }
func (OverloadYellow) GoAgain() bool               { return false }
func (OverloadYellow) Dominate()                   {}
// not implemented: on-hit go-again rider (Overload's printed power tops out at 3 so it'd only matter once Dominate lets it land)
func (OverloadYellow) NotImplemented()             {}
func (c OverloadYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type OverloadBlue struct{}

func (OverloadBlue) ID() card.ID                 { return card.OverloadBlue }
func (OverloadBlue) Name() string                { return "Overload (Blue)" }
func (OverloadBlue) Cost(*card.TurnState) int                   { return 0 }
func (OverloadBlue) Pitch() int                  { return 3 }
func (OverloadBlue) Attack() int                 { return 1 }
func (OverloadBlue) Defense() int                { return 2 }
func (OverloadBlue) Types() card.TypeSet         { return overloadTypes }
func (OverloadBlue) GoAgain() bool               { return false }
func (OverloadBlue) Dominate()                   {}
// not implemented: on-hit go-again rider (Overload's printed power tops out at 3 so it'd only matter once Dominate lets it land)
func (OverloadBlue) NotImplemented()             {}
func (c OverloadBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
