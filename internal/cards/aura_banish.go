// Shared helper for cards that banish an aura from the graveyard for +1 arcane damage.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// banishAuraFromGraveyard banishes the first aura-typed card in the graveyard, flips
// ArcaneDamageDealt, and returns 1. Returns 0 when no aura is found. Callers that also
// destroy the source card (e.g. Sigil of Silphidae's leave trigger) should run this scan
// BEFORE AddToGraveyard'ing the source so the printed "another aura" restriction is
// satisfied naturally. Routes through s.BanishFromGraveyard so the cacheable bit flips —
// reading the graveyard contents (which may include cards from prior turns) makes the
// chain output depend on hidden state.
func banishAuraFromGraveyard(s *sim.TurnState) int {
	if _, ok := s.BanishFromGraveyard(func(c sim.Card) bool {
		return c.Types().Has(card.TypeAura)
	}); !ok {
		return 0
	}
	return s.DealArcaneDamage(1)
}
