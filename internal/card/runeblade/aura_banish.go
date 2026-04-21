// Shared helper for cards that banish an aura from the graveyard for +1 arcane damage:
// Weeping Battleground, Sigil of Silphidae (both enter and leave triggers), Runic Fellingsong.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// banishAuraFromGraveyard scans s.Graveyard for the first aura-typed card that isn't skip,
// moves it to s.Banish, flips ArcaneDamageDealt, and returns 1. Returns 0 when no eligible
// aura is found. skip is nil for the usual "banish an aura" phrasing; Sigil of Silphidae's
// leave trigger passes itself so the "another aura" restriction from the printed text is
// honoured once it joins the graveyard.
func banishAuraFromGraveyard(s *card.TurnState, skip card.Card) int {
	for i, c := range s.Graveyard {
		if c == skip {
			continue
		}
		if !c.Types().Has(card.TypeAura) {
			continue
		}
		s.Banish = append(s.Banish, c)
		s.Graveyard = append(s.Graveyard[:i], s.Graveyard[i+1:]...)
		s.ArcaneDamageDealt = true
		return 1
	}
	return 0
}
