// Shared aura-creation helper for Generic Action - Aura cards.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// setAuraCreated marks the turn state so cards that read AuraCreated (e.g. Yinti Yanti,
// Runerager Swarm) see the aura entering play. Returns 0 — the aura itself contributes
// no direct damage; its value is in the flag it leaves behind.
func setAuraCreated(s *card.TurnState) int {
	s.AuraCreated = true
	return 0
}
