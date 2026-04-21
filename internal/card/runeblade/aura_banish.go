// Shared helper for cards that banish an aura from the graveyard for +1 arcane damage:
// Weeping Battleground, Sigil of Silphidae (both enter and leave triggers), Runic Fellingsong.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// banishAuraFromGraveyard scans s.Graveyard for the first aura-typed card, moves it to
// s.Banish, flips ArcaneDamageDealt, and returns 1. Returns 0 when no aura is found.
// Callers that also destroy the source card (e.g. Sigil of Silphidae's leave trigger) should
// run this scan BEFORE adding the source to s.Graveyard so the printed "another aura"
// restriction is satisfied naturally.
func banishAuraFromGraveyard(s *card.TurnState) int {
	for i, c := range s.Graveyard {
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
