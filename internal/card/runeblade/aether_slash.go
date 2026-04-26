// Aether Slash — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 1.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When Aether Slash attacks, if a 'non-attack' action card was pitched to play it, deal 1
// arcane damage to any target."

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var aetherSlashTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type AetherSlashRed struct{}

func (AetherSlashRed) ID() card.ID                 { return card.AetherSlashRed }
func (AetherSlashRed) Name() string                 { return "Aether Slash (Red)" }
func (AetherSlashRed) Cost(*card.TurnState) int                    { return 1 }
func (AetherSlashRed) Pitch() int                   { return 1 }
func (AetherSlashRed) Attack() int                  { return 4 }
func (AetherSlashRed) Defense() int                 { return 3 }
func (AetherSlashRed) Types() card.TypeSet       { return aetherSlashTypes }
func (AetherSlashRed) GoAgain() bool                { return false }
// not implemented: Pitched scan can fire the +1 arcane rider whenever any non-attack action is in
// Pitched, regardless of which pitched card actually paid for Aether Slash (over-credits when both
// an attack and a non-attack action are pitched)
func (AetherSlashRed) NotImplemented()              {}
func (c AetherSlashRed) Play(s *card.TurnState, _ *card.CardState) int { return aetherSlashPlay(c.Attack(), s) }

type AetherSlashYellow struct{}

func (AetherSlashYellow) ID() card.ID                 { return card.AetherSlashYellow }
func (AetherSlashYellow) Name() string                 { return "Aether Slash (Yellow)" }
func (AetherSlashYellow) Cost(*card.TurnState) int                    { return 1 }
func (AetherSlashYellow) Pitch() int                   { return 2 }
func (AetherSlashYellow) Attack() int                  { return 3 }
func (AetherSlashYellow) Defense() int                 { return 3 }
func (AetherSlashYellow) Types() card.TypeSet       { return aetherSlashTypes }
func (AetherSlashYellow) GoAgain() bool                { return false }
// not implemented: Pitched scan can fire the +1 arcane rider whenever any non-attack action is in
// Pitched, regardless of which pitched card actually paid for Aether Slash (over-credits when both
// an attack and a non-attack action are pitched)
func (AetherSlashYellow) NotImplemented()              {}
func (c AetherSlashYellow) Play(s *card.TurnState, _ *card.CardState) int { return aetherSlashPlay(c.Attack(), s) }

type AetherSlashBlue struct{}

func (AetherSlashBlue) ID() card.ID                 { return card.AetherSlashBlue }
func (AetherSlashBlue) Name() string                 { return "Aether Slash (Blue)" }
func (AetherSlashBlue) Cost(*card.TurnState) int                    { return 1 }
func (AetherSlashBlue) Pitch() int                   { return 3 }
func (AetherSlashBlue) Attack() int                  { return 2 }
func (AetherSlashBlue) Defense() int                 { return 3 }
func (AetherSlashBlue) Types() card.TypeSet       { return aetherSlashTypes }
func (AetherSlashBlue) GoAgain() bool                { return false }
// not implemented: Pitched scan can fire the +1 arcane rider whenever any non-attack action is in
// Pitched, regardless of which pitched card actually paid for Aether Slash (over-credits when both
// an attack and a non-attack action are pitched)
func (AetherSlashBlue) NotImplemented()              {}
func (c AetherSlashBlue) Play(s *card.TurnState, _ *card.CardState) int { return aetherSlashPlay(c.Attack(), s) }

func aetherSlashPlay(base int, s *card.TurnState) int {
	for _, p := range s.Pitched {
		if p.Types().IsNonAttackAction() {
			return base + s.DealArcaneDamage(1)
		}
	}
	return base
}
