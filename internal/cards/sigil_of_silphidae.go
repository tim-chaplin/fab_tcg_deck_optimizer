// Sigil of Silphidae — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 3, Go again.
// Only printed in Blue.
// Text: "When this enters or leaves the arena, you may banish another aura from your graveyard.
// If you do, deal 1 arcane damage to target hero. At the beginning of your action phase, destroy
// this."
//
// Play resolves the enter trigger directly via banishAuraFromGraveyard. The start-of-turn
// handler runs the leave trigger.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (c SigilOfSilphidaeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	enterDamage := banishAuraFromGraveyard(s)
	s.AddAura(sim.Aura{
		Trigger: sim.Trigger{TriggerType: sim.TriggerStartOfTurn, Handler: sigilOfSilphidaeAuraHandler},
		Self:    sim.CardOrTokenType{Card: c},
		Count:   1,
	})
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	if enterDamage > 0 {
		s.AddValue(enterDamage)
		s.LogRider(self, enterDamage, "Banished an aura, dealt 1 arcane damage")
	}
}

// sigilOfSilphidaeAuraHandler runs the leave trigger on the next turn: scans the graveyard
// for an aura to banish, credits 1 arcane damage on a hit, then destroys the aura.
func sigilOfSilphidaeAuraHandler(s *sim.TurnState, _ *sim.Trigger, a *sim.Aura) {
	n := banishAuraFromGraveyard(s)
	if n > 0 {
		s.AddValue(n)
		s.LogPostTrigger(a.Self.DisplayName(), "Banished an aura, dealt 1 arcane damage", n)
	}
	s.DestroyAura(a, true)
}
