// Hit the High Notes — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you've played or created an aura this turn, this gets +2{p}."
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var hitTheHighNotesTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type HitTheHighNotesRed struct{}

func (HitTheHighNotesRed) Name() string                  { return "Hit the High Notes (Red)" }
func (HitTheHighNotesRed) Cost() int                     { return 1 }
func (HitTheHighNotesRed) Pitch() int                    { return 1 }
func (HitTheHighNotesRed) Attack() int                   { return 4 }
func (HitTheHighNotesRed) Defense() int                  { return 3 }
func (HitTheHighNotesRed) Types() card.TypeSet        { return hitTheHighNotesTypes }
func (HitTheHighNotesRed) GoAgain() bool                 { return false }
func (c HitTheHighNotesRed) Play(s *card.TurnState) int  { return hitTheHighNotesPlay(c.Attack(), s) }

type HitTheHighNotesYellow struct{}

func (HitTheHighNotesYellow) Name() string                 { return "Hit the High Notes (Yellow)" }
func (HitTheHighNotesYellow) Cost() int                    { return 1 }
func (HitTheHighNotesYellow) Pitch() int                   { return 2 }
func (HitTheHighNotesYellow) Attack() int                  { return 3 }
func (HitTheHighNotesYellow) Defense() int                 { return 3 }
func (HitTheHighNotesYellow) Types() card.TypeSet       { return hitTheHighNotesTypes }
func (HitTheHighNotesYellow) GoAgain() bool                { return false }
func (c HitTheHighNotesYellow) Play(s *card.TurnState) int { return hitTheHighNotesPlay(c.Attack(), s) }

type HitTheHighNotesBlue struct{}

func (HitTheHighNotesBlue) Name() string                 { return "Hit the High Notes (Blue)" }
func (HitTheHighNotesBlue) Cost() int                    { return 1 }
func (HitTheHighNotesBlue) Pitch() int                   { return 3 }
func (HitTheHighNotesBlue) Attack() int                  { return 2 }
func (HitTheHighNotesBlue) Defense() int                 { return 3 }
func (HitTheHighNotesBlue) Types() card.TypeSet       { return hitTheHighNotesTypes }
func (HitTheHighNotesBlue) GoAgain() bool                { return false }
func (c HitTheHighNotesBlue) Play(s *card.TurnState) int { return hitTheHighNotesPlay(c.Attack(), s) }

func hitTheHighNotesPlay(base int, s *card.TurnState) int {
	if s.AuraCreated || s.HasPlayedType(card.TypeAura) {
		return base + 2
	}
	return base
}
