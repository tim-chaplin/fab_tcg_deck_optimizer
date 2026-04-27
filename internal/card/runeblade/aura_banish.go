// Shared helper for cards that banish an aura from the graveyard for +1 arcane damage.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// banishAuraFromGraveyard finds the first aura-typed card in the graveyard, moves it to
// s.Banish, flips ArcaneDamageDealt, and returns 1. Returns 0 when no aura is found.
// Callers that also destroy the source card (e.g. Sigil of Silphidae's leave trigger) should
// run this scan BEFORE adding the source to the graveyard so the printed "another aura"
// restriction is satisfied naturally. BanishFromGraveyard flips Cacheable=false — scanning
// prior-turn graveyard contents makes the chain depend on hidden state.
func banishAuraFromGraveyard(s *card.TurnState) int {
	if _, ok := s.BanishFromGraveyard(isAura); !ok {
		return 0
	}
	return s.DealArcaneDamage(1)
}

// isAura is the predicate passed to TurnState.BanishFromGraveyard.
func isAura(c card.Card) bool { return c.Types().Has(card.TypeAura) }
