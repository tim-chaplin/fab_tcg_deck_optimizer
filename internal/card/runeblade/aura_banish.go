// Shared helper for cards that banish an aura from the graveyard for +1 arcane damage.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// banishAuraFromGraveyard scans s.Graveyard() for the first aura-typed card, moves it to
// s.Banish, flips ArcaneDamageDealt, and returns 1. Returns 0 when no aura is found.
// Callers that also destroy the source card (e.g. Sigil of Silphidae's leave trigger) should
// run this scan BEFORE adding the source to the graveyard so the printed "another aura"
// restriction is satisfied naturally. Graveyard() flips Cacheable=false — scanning prior-
// turn graveyard contents makes the chain depend on hidden state.
func banishAuraFromGraveyard(s *card.TurnState) int {
	gy := s.Graveyard()
	for i, c := range gy {
		if !c.Types().Has(card.TypeAura) {
			continue
		}
		s.Banish = append(s.Banish, c)
		s.SetGraveyard(append(gy[:i], gy[i+1:]...))
		return s.DealArcaneDamage(1)
	}
	return 0
}
