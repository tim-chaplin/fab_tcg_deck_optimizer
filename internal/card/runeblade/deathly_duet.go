// Deathly Duet — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When Deathly Duet attacks, if an attack action card was pitched to play it, it gains
// +2{p}. If a 'non-attack' action card was pitched to play it, create 2 Runechant tokens."
//
// Simplification: both riders scan Pitched (we don't track which pitched card paid for which
// play; any attack in Pitched satisfies the +2{p} branch, any non-attack action satisfies the
// runechant branch, and both can fire if both kinds were pitched). The 2 Runechants enter state
// via CreateRunechants — they fire on Deathly Duet's own attack resolution downstream.
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

// deathlyDuetPlay applies two independent pitched-card riders:
//   - Attack pitched: +2{p} static buff. Credited only when the buffed total is likely to land
//     — a blockable buffed attack delivers nothing, so we don't count it.
//   - Non-attack action pitched: 2 Runechant tokens. These are arcane damage, not blocked by
//     physical blocks, so we always credit them regardless of LikelyToHit.
func deathlyDuetPlay(base int, s *card.TurnState) int {
	var attackPitched, nonAttackActionPitched bool
	for _, p := range s.Pitched {
		t := p.Types()
		if t.Has(card.TypeAttack) {
			attackPitched = true
		}
		if t.IsNonAttackAction() {
			nonAttackActionPitched = true
		}
	}

	dmg := base
	if attackPitched && (card.LikelyToHit(base+2) || card.LikelyToHit(s.Runechants)) {
		dmg += 2
	}
	if nonAttackActionPitched {
		// Two Runechants enter during Deathly Duet's own attack resolution. Deathly Duet IS
		// the attack, so no lookahead for a following attack is needed.
		dmg += s.CreateRunechants(2)
	}
	return dmg
}
