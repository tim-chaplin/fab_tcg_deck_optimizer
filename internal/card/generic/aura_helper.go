// Shared aura-creation helper for Generic Action - Aura cards.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// setAuraCreated flips s.AuraCreated so cards that read it see the aura entering play.
// Returns 0 — the aura's value is in the flag, not direct damage.
func setAuraCreated(s *card.TurnState) int {
	s.AuraCreated = true
	return 0
}
