// Mauvrion Skies — Runeblade Action. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "The next Runeblade attack action card you play this turn gains go again and 'If this
// hits, create N Runechant tokens.'"
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling splits the two grants by when they need to resolve:
//   - Go again is a static property of the target card (not gated on how the attack plays
//     out), and must be visible before the target's chain-legality check. Play scans
//     CardsRemaining for the first matching Runeblade attack action card and flips its
//     GrantedGoAgain immediately.
//   - The "if this hits" Runechant rider depends on the target's fully-resolved attack
//     state — the target's own Play effects, hero triggers, and aura triggers can shift
//     what LikelyToHit sees (e.g. Drowning Dire gaining Dominate from an aura created
//     mid-chain). Play registers an EphemeralAttackTrigger that fires after those settle;
//     the Handler reads target.EffectiveDominate() to pick up printed + granted + conditional
//     Dominate in one probe. Damage attributes back to Mauvrion via SourceIndex.
//
// "Attack action card" excludes weapons; the Matches predicate requires both
// TypeRuneblade and the attack-action bit.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var mauvrionSkiesTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

// mauvrionTargetMatches is the shared predicate for Mauvrion's grants: the next Runeblade
// attack action card (weapons don't qualify).
func mauvrionTargetMatches(target *card.CardState) bool {
	t := target.Card.Types()
	return t.Has(card.TypeRuneblade) && t.IsAttackAction()
}

// mauvrionSkiesPlay applies the go-again grant via look-ahead and registers the ephemeral
// "if hits, create n Runechants" trigger for the same target's fully-resolved attack.
func mauvrionSkiesPlay(s *card.TurnState, source card.Card, n int) int {
	// Go again is static — flip it on the first matching target so its chain-legality check
	// sees the grant before the target's Play runs.
	for _, pc := range s.CardsRemaining {
		if mauvrionTargetMatches(pc) {
			pc.GrantedGoAgain = true
			break
		}
	}
	// The Runechant rider depends on whether the attack actually hits, so defer until the
	// target's full resolution. The trigger fires on the first matching attack action's
	// post-resolution; non-matching attacks (e.g. a Generic attack played before the
	// Runeblade one) leave it in place. Unconsumed triggers fizzle silently at end of turn.
	s.AddEphemeralAttackTrigger(card.EphemeralAttackTrigger{
		Source:  source,
		Matches: mauvrionTargetMatches,
		Handler: func(s *card.TurnState, target *card.CardState) int {
			if card.LikelyToHit(target) {
				return s.CreateRunechants(n)
			}
			return 0
		},
	})
	return 0
}

type MauvrionSkiesRed struct{}

func (MauvrionSkiesRed) ID() card.ID                                     { return card.MauvrionSkiesRed }
func (MauvrionSkiesRed) Name() string                                    { return "Mauvrion Skies (Red)" }
func (MauvrionSkiesRed) Cost(*card.TurnState) int                        { return 0 }
func (MauvrionSkiesRed) Pitch() int                                      { return 1 }
func (MauvrionSkiesRed) Attack() int                                     { return 0 }
func (MauvrionSkiesRed) Defense() int                                    { return 2 }
func (MauvrionSkiesRed) Types() card.TypeSet                             { return mauvrionSkiesTypes }
func (MauvrionSkiesRed) GoAgain() bool                                   { return true }
func (c MauvrionSkiesRed) Play(s *card.TurnState, _ *card.CardState) int { return mauvrionSkiesPlay(s, c, 3) }

type MauvrionSkiesYellow struct{}

func (MauvrionSkiesYellow) ID() card.ID                                     { return card.MauvrionSkiesYellow }
func (MauvrionSkiesYellow) Name() string                                    { return "Mauvrion Skies (Yellow)" }
func (MauvrionSkiesYellow) Cost(*card.TurnState) int                        { return 0 }
func (MauvrionSkiesYellow) Pitch() int                                      { return 2 }
func (MauvrionSkiesYellow) Attack() int                                     { return 0 }
func (MauvrionSkiesYellow) Defense() int                                    { return 2 }
func (MauvrionSkiesYellow) Types() card.TypeSet                             { return mauvrionSkiesTypes }
func (MauvrionSkiesYellow) GoAgain() bool                                   { return true }
func (c MauvrionSkiesYellow) Play(s *card.TurnState, _ *card.CardState) int { return mauvrionSkiesPlay(s, c, 2) }

type MauvrionSkiesBlue struct{}

func (MauvrionSkiesBlue) ID() card.ID                                     { return card.MauvrionSkiesBlue }
func (MauvrionSkiesBlue) Name() string                                    { return "Mauvrion Skies (Blue)" }
func (MauvrionSkiesBlue) Cost(*card.TurnState) int                        { return 0 }
func (MauvrionSkiesBlue) Pitch() int                                      { return 3 }
func (MauvrionSkiesBlue) Attack() int                                     { return 0 }
func (MauvrionSkiesBlue) Defense() int                                    { return 2 }
func (MauvrionSkiesBlue) Types() card.TypeSet                             { return mauvrionSkiesTypes }
func (MauvrionSkiesBlue) GoAgain() bool                                   { return true }
func (c MauvrionSkiesBlue) Play(s *card.TurnState, _ *card.CardState) int { return mauvrionSkiesPlay(s, c, 1) }
