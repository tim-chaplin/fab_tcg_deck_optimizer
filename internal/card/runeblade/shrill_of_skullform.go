// Shrill of Skullform — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you have played or created an aura this turn, Shrill of Skullform gains +3{p}."

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var shrillTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type ShrillOfSkullformRed struct{}

func (ShrillOfSkullformRed) ID() card.ID                 { return card.ShrillOfSkullformRed }
func (ShrillOfSkullformRed) Name() string              { return "Shrill of Skullform" }
func (ShrillOfSkullformRed) Cost(*card.TurnState) int                 { return 2 }
func (ShrillOfSkullformRed) Pitch() int                { return 1 }
func (ShrillOfSkullformRed) Attack() int               { return 4 }
func (ShrillOfSkullformRed) Defense() int              { return 3 }
func (ShrillOfSkullformRed) Types() card.TypeSet       { return shrillTypes }
func (ShrillOfSkullformRed) GoAgain() bool             { return false }
func (c ShrillOfSkullformRed) Play(s *card.TurnState, _ *card.CardState) int {
	return shrillPlay(c.Attack(), s)
}

type ShrillOfSkullformYellow struct{}

func (ShrillOfSkullformYellow) ID() card.ID                 { return card.ShrillOfSkullformYellow }
func (ShrillOfSkullformYellow) Name() string           { return "Shrill of Skullform" }
func (ShrillOfSkullformYellow) Cost(*card.TurnState) int              { return 2 }
func (ShrillOfSkullformYellow) Pitch() int             { return 2 }
func (ShrillOfSkullformYellow) Attack() int            { return 3 }
func (ShrillOfSkullformYellow) Defense() int           { return 3 }
func (ShrillOfSkullformYellow) Types() card.TypeSet    { return shrillTypes }
func (ShrillOfSkullformYellow) GoAgain() bool          { return false }
func (c ShrillOfSkullformYellow) Play(s *card.TurnState, _ *card.CardState) int {
	return shrillPlay(c.Attack(), s)
}

type ShrillOfSkullformBlue struct{}

func (ShrillOfSkullformBlue) ID() card.ID                 { return card.ShrillOfSkullformBlue }
func (ShrillOfSkullformBlue) Name() string             { return "Shrill of Skullform" }
func (ShrillOfSkullformBlue) Cost(*card.TurnState) int                { return 2 }
func (ShrillOfSkullformBlue) Pitch() int               { return 3 }
func (ShrillOfSkullformBlue) Attack() int              { return 2 }
func (ShrillOfSkullformBlue) Defense() int             { return 3 }
func (ShrillOfSkullformBlue) Types() card.TypeSet      { return shrillTypes }
func (ShrillOfSkullformBlue) GoAgain() bool            { return false }
func (c ShrillOfSkullformBlue) Play(s *card.TurnState, _ *card.CardState) int {
	return shrillPlay(c.Attack(), s)
}

func shrillPlay(base int, s *card.TurnState) int {
	if s.HasAuraInPlay() {
		return base + 3
	}
	return base
}
