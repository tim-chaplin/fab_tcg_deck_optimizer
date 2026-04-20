// Fiddler's Green — Generic Block. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 1.
//
// Text: "When this is put into your graveyard from anywhere, gain 3{h}."
//
// Modelling: using this card to defend sends it to the graveyard, so the 3{h} gain fires on
// the DR Play path — credited as +3 damage equivalent (health is valued 1-to-1 with damage).
// Pitched copies go to the bottom of the deck, not the graveyard, so they don't trigger the
// rider; the pitch-role contribution stays at the printed pitch value only.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type FiddlersGreenRed struct{}

func (FiddlersGreenRed) ID() card.ID                 { return card.FiddlersGreenRed }
func (FiddlersGreenRed) Name() string                { return "Fiddler's Green (Red)" }
func (FiddlersGreenRed) Cost(*card.TurnState) int                   { return 0 }
func (FiddlersGreenRed) Pitch() int                  { return 1 }
func (FiddlersGreenRed) Attack() int                 { return 0 }
func (FiddlersGreenRed) Defense() int                { return 1 }
func (FiddlersGreenRed) Types() card.TypeSet         { return defenseReactionTypes }
func (FiddlersGreenRed) GoAgain() bool               { return false }
func (FiddlersGreenRed) NotSilverAgeLegal()           {}
func (FiddlersGreenRed) Play(s *card.TurnState) int { return 3 }

type FiddlersGreenYellow struct{}

func (FiddlersGreenYellow) ID() card.ID                 { return card.FiddlersGreenYellow }
func (FiddlersGreenYellow) Name() string                { return "Fiddler's Green (Yellow)" }
func (FiddlersGreenYellow) Cost(*card.TurnState) int                   { return 0 }
func (FiddlersGreenYellow) Pitch() int                  { return 2 }
func (FiddlersGreenYellow) Attack() int                 { return 0 }
func (FiddlersGreenYellow) Defense() int                { return 1 }
func (FiddlersGreenYellow) Types() card.TypeSet         { return defenseReactionTypes }
func (FiddlersGreenYellow) GoAgain() bool               { return false }
func (FiddlersGreenYellow) NotSilverAgeLegal()           {}
func (FiddlersGreenYellow) Play(s *card.TurnState) int { return 3 }

type FiddlersGreenBlue struct{}

func (FiddlersGreenBlue) ID() card.ID                 { return card.FiddlersGreenBlue }
func (FiddlersGreenBlue) Name() string                { return "Fiddler's Green (Blue)" }
func (FiddlersGreenBlue) Cost(*card.TurnState) int                   { return 0 }
func (FiddlersGreenBlue) Pitch() int                  { return 3 }
func (FiddlersGreenBlue) Attack() int                 { return 0 }
func (FiddlersGreenBlue) Defense() int                { return 1 }
func (FiddlersGreenBlue) Types() card.TypeSet         { return defenseReactionTypes }
func (FiddlersGreenBlue) GoAgain() bool               { return false }
func (FiddlersGreenBlue) NotSilverAgeLegal()           {}
func (FiddlersGreenBlue) Play(s *card.TurnState) int { return 3 }
