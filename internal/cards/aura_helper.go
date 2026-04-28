// Shared aura-creation helper for Generic Action - Aura cards.

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

// setAuraCreated flips s.AuraCreated so cards that read it see the aura entering play.
func setAuraCreated(s *sim.TurnState) {
	s.AuraCreated = true
}
