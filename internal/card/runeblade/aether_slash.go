// Aether Slash — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 1.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When Aether Slash attacks, if a 'non-attack' action card was pitched to play it, deal 1
// arcane damage to any target."
//
// Simplification: the printed 1 arcane is added to base damage unconditionally. The rider's +1
// arcane fires if any non-attack action card appears in Pitched (we don't track which pitched
// card paid for which play, so any qualifier suffices). Both land as flat damage.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var aetherSlashTypes = map[string]bool{"Runeblade": true, "Action": true, "Attack": true}

type AetherSlashRed struct{}

func (AetherSlashRed) Name() string                 { return "Aether Slash (Red)" }
func (AetherSlashRed) Cost() int                    { return 1 }
func (AetherSlashRed) Pitch() int                   { return 1 }
func (AetherSlashRed) Attack() int                  { return 4 }
func (AetherSlashRed) Defense() int                 { return 3 }
func (AetherSlashRed) Types() map[string]bool       { return aetherSlashTypes }
func (AetherSlashRed) GoAgain() bool                { return false }
func (c AetherSlashRed) Play(s *card.TurnState) int { return aetherSlashPlay(c.Attack(), s) }

type AetherSlashYellow struct{}

func (AetherSlashYellow) Name() string                 { return "Aether Slash (Yellow)" }
func (AetherSlashYellow) Cost() int                    { return 1 }
func (AetherSlashYellow) Pitch() int                   { return 2 }
func (AetherSlashYellow) Attack() int                  { return 3 }
func (AetherSlashYellow) Defense() int                 { return 3 }
func (AetherSlashYellow) Types() map[string]bool       { return aetherSlashTypes }
func (AetherSlashYellow) GoAgain() bool                { return false }
func (c AetherSlashYellow) Play(s *card.TurnState) int { return aetherSlashPlay(c.Attack(), s) }

type AetherSlashBlue struct{}

func (AetherSlashBlue) Name() string                 { return "Aether Slash (Blue)" }
func (AetherSlashBlue) Cost() int                    { return 1 }
func (AetherSlashBlue) Pitch() int                   { return 3 }
func (AetherSlashBlue) Attack() int                  { return 2 }
func (AetherSlashBlue) Defense() int                 { return 3 }
func (AetherSlashBlue) Types() map[string]bool       { return aetherSlashTypes }
func (AetherSlashBlue) GoAgain() bool                { return false }
func (c AetherSlashBlue) Play(s *card.TurnState) int { return aetherSlashPlay(c.Attack(), s) }

func aetherSlashPlay(base int, s *card.TurnState) int {
	dmg := base + 1 // printed arcane
	for _, p := range s.Pitched {
		t := p.Types()
		if t["Action"] && !t["Attack"] {
			dmg++
			break
		}
	}
	return dmg
}
