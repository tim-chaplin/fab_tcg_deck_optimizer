// Deathly Duet — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When Deathly Duet attacks, if an attack action card was pitched to play it, it gains
// +2{p}. If a 'non-attack' action card was pitched to play it, create 2 Runechant tokens."
//
// Simplifications:
//   - Both riders scan Pitched (we don't track which pitched card paid for which play; any attack
//     in Pitched satisfies the +2{p} branch, any non-attack action satisfies the runechant branch,
//     and both can fire if both kinds were pitched).
//   - The 2 Runechants are counted as +2 flat damage only if another attack (card OR weapon)
//     follows in CardsRemaining; otherwise the runechants fizzle at end of turn. AuraCreated is
//     set in the same case so following aura-conditional cards (e.g. Shrill of Skullform) see it.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var deathlyDuetTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type DeathlyDuetRed struct{}

func (DeathlyDuetRed) ID() card.ID                 { return card.DeathlyDuetRed }
func (DeathlyDuetRed) Name() string                 { return "Deathly Duet (Red)" }
func (DeathlyDuetRed) Cost() int                    { return 2 }
func (DeathlyDuetRed) Pitch() int                   { return 1 }
func (DeathlyDuetRed) Attack() int                  { return 4 }
func (DeathlyDuetRed) Defense() int                 { return 3 }
func (DeathlyDuetRed) Types() card.TypeSet       { return deathlyDuetTypes }
func (DeathlyDuetRed) GoAgain() bool                { return false }
func (c DeathlyDuetRed) Play(s *card.TurnState) int { return deathlyDuetPlay(c.Attack(), s) }

type DeathlyDuetYellow struct{}

func (DeathlyDuetYellow) ID() card.ID                 { return card.DeathlyDuetYellow }
func (DeathlyDuetYellow) Name() string                 { return "Deathly Duet (Yellow)" }
func (DeathlyDuetYellow) Cost() int                    { return 2 }
func (DeathlyDuetYellow) Pitch() int                   { return 2 }
func (DeathlyDuetYellow) Attack() int                  { return 3 }
func (DeathlyDuetYellow) Defense() int                 { return 3 }
func (DeathlyDuetYellow) Types() card.TypeSet       { return deathlyDuetTypes }
func (DeathlyDuetYellow) GoAgain() bool                { return false }
func (c DeathlyDuetYellow) Play(s *card.TurnState) int { return deathlyDuetPlay(c.Attack(), s) }

type DeathlyDuetBlue struct{}

func (DeathlyDuetBlue) ID() card.ID                 { return card.DeathlyDuetBlue }
func (DeathlyDuetBlue) Name() string                 { return "Deathly Duet (Blue)" }
func (DeathlyDuetBlue) Cost() int                    { return 2 }
func (DeathlyDuetBlue) Pitch() int                   { return 3 }
func (DeathlyDuetBlue) Attack() int                  { return 2 }
func (DeathlyDuetBlue) Defense() int                 { return 3 }
func (DeathlyDuetBlue) Types() card.TypeSet       { return deathlyDuetTypes }
func (DeathlyDuetBlue) GoAgain() bool                { return false }
func (c DeathlyDuetBlue) Play(s *card.TurnState) int { return deathlyDuetPlay(c.Attack(), s) }

func deathlyDuetPlay(base int, s *card.TurnState) int {
	var attackPitched, nonAttackActionPitched bool
	for _, p := range s.Pitched {
		t := p.Types()
		if t.Has(card.TypeAttack) {
			attackPitched = true
		}
		if t.Has(card.TypeAction) && !t.Has(card.TypeAttack) {
			nonAttackActionPitched = true
		}
	}

	dmg := base
	if attackPitched {
		dmg += 2
	}
	if nonAttackActionPitched && hasFollowingAttack(s) {
		dmg += s.CreateRunechants(2) // two Runechants, each dealing 1 when the next attack hits.
	}
	return dmg
}

// hasFollowingAttack reports whether any card in CardsRemaining is an attack — either an attack
// action card (Types["Attack"]) or a weapon (Types["Weapon"]). Used to decide whether a Runechant
// created by the current attack will land before end of turn.
func hasFollowingAttack(s *card.TurnState) bool {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if t.Has(card.TypeAttack) || t.Has(card.TypeWeapon) {
			return true
		}
	}
	return false
}
