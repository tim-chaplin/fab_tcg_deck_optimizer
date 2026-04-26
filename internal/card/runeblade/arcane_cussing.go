// Arcane Cussing — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. When you deal or are dealt damage, destroy this. When this leaves the arena
// during your turn, create N Runechants." (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelled as a fragile aura (fragile_aura.go). The "deal damage" trigger lets any attacker
// type pop the aura, so Play passes attackActionOnly=false.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var arcaneCussingTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type ArcaneCussingRed struct{}

func (ArcaneCussingRed) ID() card.ID              { return card.ArcaneCussingRed }
func (ArcaneCussingRed) Name() string             { return "Arcane Cussing" }
func (ArcaneCussingRed) Cost(*card.TurnState) int { return 1 }
func (ArcaneCussingRed) Pitch() int               { return 1 }
func (ArcaneCussingRed) Attack() int              { return 0 }
func (ArcaneCussingRed) Defense() int             { return 2 }
func (ArcaneCussingRed) Types() card.TypeSet      { return arcaneCussingTypes }
func (ArcaneCussingRed) GoAgain() bool            { return true }
func (ArcaneCussingRed) Play(s *card.TurnState, self *card.CardState) {
	fragileAuraPlay(s, self, 3, false)
}

type ArcaneCussingYellow struct{}

func (ArcaneCussingYellow) ID() card.ID              { return card.ArcaneCussingYellow }
func (ArcaneCussingYellow) Name() string             { return "Arcane Cussing" }
func (ArcaneCussingYellow) Cost(*card.TurnState) int { return 1 }
func (ArcaneCussingYellow) Pitch() int               { return 2 }
func (ArcaneCussingYellow) Attack() int              { return 0 }
func (ArcaneCussingYellow) Defense() int             { return 2 }
func (ArcaneCussingYellow) Types() card.TypeSet      { return arcaneCussingTypes }
func (ArcaneCussingYellow) GoAgain() bool            { return true }
func (ArcaneCussingYellow) Play(s *card.TurnState, self *card.CardState) {
	fragileAuraPlay(s, self, 2, false)
}

type ArcaneCussingBlue struct{}

func (ArcaneCussingBlue) ID() card.ID              { return card.ArcaneCussingBlue }
func (ArcaneCussingBlue) Name() string             { return "Arcane Cussing" }
func (ArcaneCussingBlue) Cost(*card.TurnState) int { return 1 }
func (ArcaneCussingBlue) Pitch() int               { return 3 }
func (ArcaneCussingBlue) Attack() int              { return 0 }
func (ArcaneCussingBlue) Defense() int             { return 2 }
func (ArcaneCussingBlue) Types() card.TypeSet      { return arcaneCussingTypes }
func (ArcaneCussingBlue) GoAgain() bool            { return true }
func (ArcaneCussingBlue) Play(s *card.TurnState, self *card.CardState) {
	fragileAuraPlay(s, self, 1, false)
}
