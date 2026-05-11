// Blessing of Occult — Runeblade Action - Aura. Cost 1, Defense 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "At the start of your turn, destroy Blessing of Occult then create N Runechant tokens."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Handler creates N Runechants next turn.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// blessingOfOccultTriggerText pre-formats the trigger log text for each Runechant count
// (1 = Blue, 2 = Yellow, 3 = Red). The text is captured into a per-play closure on every
// cast, so a constant lookup avoids a Sprintf alloc per cast.
var blessingOfOccultTriggerText = [...]string{
	1: "Created a runechant",
	2: "Created 2 runechants",
	3: "Created 3 runechants",
}

func (c BlessingOfOccultRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	blessingOfOccultPlay(s, l, self, c, 3)
}

func (c BlessingOfOccultYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	blessingOfOccultPlay(s, l, self, c, 2)
}

func (c BlessingOfOccultBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	blessingOfOccultPlay(s, l, self, c, 1)
}

// blessingOfOccultPlay registers the shared next-turn trigger that creates n Runechants
// and emits the same-turn chain step (no value contribution; all credit is deferred to
// the trigger).
//
// blessingOfOccultHandler creates aura.Count runechants and destroys the aura. Count
// carries the per-variant rune count (R=3 / Y=2 / B=1) — the handler is one-shot, so
// Count's "fires remaining" interpretation collapses to "runechants to create on the
// only fire".
func blessingOfOccultHandler(s *sim.TurnState, l sim.Logger, _ *sim.Trigger, a *sim.Aura) {
	n := a.Count
	name := a.Self.DisplayName()
	s.CreateRunechants(n)
	l.AppendPostTrigger(name, blessingOfOccultTriggerText[n], n)
	s.DestroyAura(a, true)
}

func blessingOfOccultPlay(s sim.GameEngine, l sim.Logger, selfState *sim.CardState, selfCard sim.Card, n int) {
	s.AddAura(sim.Aura{
		Trigger: sim.Trigger{TriggerType: sim.TriggerStartOfTurn, Handler: blessingOfOccultHandler},
		Self:    sim.CardOrTokenType{Card: selfCard},
		Count:   n,
	})
}
