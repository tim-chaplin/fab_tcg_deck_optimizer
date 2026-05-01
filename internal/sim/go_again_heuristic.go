package sim

// HasGoAgainHeuristic reports whether c likely resolves with Go again at play time —
// printed Go again counts (the trivial true), and so does conditional Go again that
// the card grants itself in Play when a permissive set of trigger flags is set.
//
// Implementation: when c.GoAgain() is false, run a sandboxed probe of c.Play against a
// fresh TurnState seeded with the conditional-trigger flags Runeblade-archetype cards
// gate on (AuraCreated, ArcaneDamageDealt, NonAttackActionPlayed, a positive runechant
// count). If the card flips self.GrantedGoAgain under those conditions, it has
// conditional Go again. The probe state is local — every side effect (log entries,
// runechant bumps, AuraTrigger registrations, deck reads) lands on the discarded
// TurnState, so the caller sees nothing.
//
// Caveats:
//   - The card's Play executes against the current value of package-global CurrentHero.
//     Cards that gate on hero markers (LowerHealthWanter, etc.) probe against the live
//     hero, which is the correct answer for "would this card have Go again in this
//     deck" — a Lower-Health-Wanter card in Viserai's deck never gets the conditional
//     grant.
//   - FromArsenal-gated grants are treated as "no go again" because the probe leaves
//     FromArsenal=false; for Opt heuristics this matches reality, since only one copy
//     per turn can ride the arsenal slot.
//   - The probe trips s.IsCacheable() to false on the throwaway state when the card's
//     Play touches the deck, but the throwaway is discarded so production caching
//     stays untouched.
//
// Use this for hand-shaping heuristics (Hero.Opt slot classifiers, etc.) where a
// printed-no-go-again card with a reliable conditional grant should not be treated as
// a one-per-hand finisher.
func HasGoAgainHeuristic(c Card) bool {
	if c.GoAgain() {
		return true
	}
	s := &TurnState{
		AuraCreated:           true,
		ArcaneDamageDealt:     true,
		NonAttackActionPlayed: true,
		Runechants:            1,
	}
	self := &CardState{Card: c}
	c.Play(s, self)
	return self.GrantedGoAgain
}
