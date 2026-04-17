// Reek of Corruption — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you have played or created an aura this turn, Reek of Corruption gains 'When this
// hits a hero, they discard a card.'"
//
// Simplifications:
//   - The aura condition is checked via `s.AuraCreated || s.HasPlayedType(card.TypeAura)` at
//     Play time — same pattern used by Hit the High Notes and Shrill of Skullform. When the
//     condition passes, the on-hit discard is valued at +3 (matching Consuming Volition's
//     discard valuation); otherwise Play returns only the printed attack.
//   - Assume the attack hits when the rider is active.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var reekOfCorruptionTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// reekDiscardBonus mirrors the discard-a-card rider's damage-equivalent used by Consuming Volition
// and other on-hit discard riders.
const reekDiscardBonus = 3

// reekOfCorruptionDamage returns the base attack plus the discard rider when the aura condition
// is satisfied. Extracted so all three printings share one implementation.
func reekOfCorruptionDamage(attack int, s *card.TurnState) int {
	if s != nil && (s.AuraCreated || s.HasPlayedType(card.TypeAura)) {
		return attack + reekDiscardBonus
	}
	return attack
}

type ReekOfCorruptionRed struct{}

func (ReekOfCorruptionRed) ID() card.ID                  { return card.ReekOfCorruptionRed }
func (ReekOfCorruptionRed) Name() string                 { return "Reek of Corruption (Red)" }
func (ReekOfCorruptionRed) Cost() int                    { return 2 }
func (ReekOfCorruptionRed) Pitch() int                   { return 1 }
func (ReekOfCorruptionRed) Attack() int                  { return 4 }
func (ReekOfCorruptionRed) Defense() int                 { return 3 }
func (ReekOfCorruptionRed) Types() card.TypeSet          { return reekOfCorruptionTypes }
func (ReekOfCorruptionRed) GoAgain() bool                { return false }
func (c ReekOfCorruptionRed) Play(s *card.TurnState) int { return reekOfCorruptionDamage(c.Attack(), s) }

type ReekOfCorruptionYellow struct{}

func (ReekOfCorruptionYellow) ID() card.ID                  { return card.ReekOfCorruptionYellow }
func (ReekOfCorruptionYellow) Name() string                 { return "Reek of Corruption (Yellow)" }
func (ReekOfCorruptionYellow) Cost() int                    { return 2 }
func (ReekOfCorruptionYellow) Pitch() int                   { return 2 }
func (ReekOfCorruptionYellow) Attack() int                  { return 3 }
func (ReekOfCorruptionYellow) Defense() int                 { return 3 }
func (ReekOfCorruptionYellow) Types() card.TypeSet          { return reekOfCorruptionTypes }
func (ReekOfCorruptionYellow) GoAgain() bool                { return false }
func (c ReekOfCorruptionYellow) Play(s *card.TurnState) int { return reekOfCorruptionDamage(c.Attack(), s) }

type ReekOfCorruptionBlue struct{}

func (ReekOfCorruptionBlue) ID() card.ID                  { return card.ReekOfCorruptionBlue }
func (ReekOfCorruptionBlue) Name() string                 { return "Reek of Corruption (Blue)" }
func (ReekOfCorruptionBlue) Cost() int                    { return 2 }
func (ReekOfCorruptionBlue) Pitch() int                   { return 3 }
func (ReekOfCorruptionBlue) Attack() int                  { return 2 }
func (ReekOfCorruptionBlue) Defense() int                 { return 3 }
func (ReekOfCorruptionBlue) Types() card.TypeSet          { return reekOfCorruptionTypes }
func (ReekOfCorruptionBlue) GoAgain() bool                { return false }
func (c ReekOfCorruptionBlue) Play(s *card.TurnState) int { return reekOfCorruptionDamage(c.Attack(), s) }
