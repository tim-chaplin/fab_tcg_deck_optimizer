// Runic Fellingsong — Runeblade Action - Attack. Cost 3, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 7, Yellow 6, Blue 5.
// Text: "When this attacks, you may banish an aura from your graveyard. If you do, deal 1 arcane
// damage to target hero."
//
// Play credits Attack() plus 1 arcane when banishAuraFromGraveyard lands an aura in s.Banish.
// No aura in the graveyard → the banish rider fizzles and Play returns just Attack().

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runicFellingsongTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// runicFellingsongPlay returns printed power + (1 if the banish rider succeeds).
func runicFellingsongPlay(c card.Card, s *card.TurnState) int {
	return c.Attack() + banishAuraFromGraveyard(s)
}

type RunicFellingsongRed struct{}

func (RunicFellingsongRed) ID() card.ID              { return card.RunicFellingsongRed }
func (RunicFellingsongRed) Name() string             { return "Runic Fellingsong (Red)" }
func (RunicFellingsongRed) Cost(*card.TurnState) int { return 3 }
func (RunicFellingsongRed) Pitch() int               { return 1 }
func (RunicFellingsongRed) Attack() int              { return 7 }
func (RunicFellingsongRed) Defense() int             { return 3 }
func (RunicFellingsongRed) Types() card.TypeSet      { return runicFellingsongTypes }
func (RunicFellingsongRed) GoAgain() bool            { return false }
func (RunicFellingsongRed) NoMemo()                  {}
func (c RunicFellingsongRed) Play(s *card.TurnState, _ *card.CardState) int {
	return runicFellingsongPlay(c, s)
}

type RunicFellingsongYellow struct{}

func (RunicFellingsongYellow) ID() card.ID              { return card.RunicFellingsongYellow }
func (RunicFellingsongYellow) Name() string             { return "Runic Fellingsong (Yellow)" }
func (RunicFellingsongYellow) Cost(*card.TurnState) int { return 3 }
func (RunicFellingsongYellow) Pitch() int               { return 2 }
func (RunicFellingsongYellow) Attack() int              { return 6 }
func (RunicFellingsongYellow) Defense() int             { return 3 }
func (RunicFellingsongYellow) Types() card.TypeSet      { return runicFellingsongTypes }
func (RunicFellingsongYellow) GoAgain() bool            { return false }
func (RunicFellingsongYellow) NoMemo()                  {}
func (c RunicFellingsongYellow) Play(s *card.TurnState, _ *card.CardState) int {
	return runicFellingsongPlay(c, s)
}

type RunicFellingsongBlue struct{}

func (RunicFellingsongBlue) ID() card.ID              { return card.RunicFellingsongBlue }
func (RunicFellingsongBlue) Name() string             { return "Runic Fellingsong (Blue)" }
func (RunicFellingsongBlue) Cost(*card.TurnState) int { return 3 }
func (RunicFellingsongBlue) Pitch() int               { return 3 }
func (RunicFellingsongBlue) Attack() int              { return 5 }
func (RunicFellingsongBlue) Defense() int             { return 3 }
func (RunicFellingsongBlue) Types() card.TypeSet      { return runicFellingsongTypes }
func (RunicFellingsongBlue) GoAgain() bool            { return false }
func (RunicFellingsongBlue) NoMemo()                  {}
func (c RunicFellingsongBlue) Play(s *card.TurnState, _ *card.CardState) int {
	return runicFellingsongPlay(c, s)
}
