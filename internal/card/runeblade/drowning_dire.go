// Drowning Dire — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 5, Yellow 4, Blue 3.
// Text: "If you have played or created an aura this turn, Drowning Dire gains **dominate**.
//
// When Drowning Dire hits, you may put a 'non-attack' action card from your graveyard on the
// bottom of your deck."
//
// Modelling: the Dominate grant is conditional, gated on s.HasAuraInPlay(). Play flips
// self.GrantedDominate when the aura clause is live so EffectiveDominate reports the card as
// dominating this turn — downstream scanners that read pc.EffectiveDominate() see the grant.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var drowningDireTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// drowningDirePlay grants self Dominate when an aura has been played or created this turn,
// then emits the chain step.
func drowningDirePlay(s *card.TurnState, self *card.CardState) {
	if s.HasAuraInPlay() {
		self.GrantedDominate = true
	}
	s.ApplyAndLogEffectiveAttack(self)
}

type DrowningDireRed struct{}

func (DrowningDireRed) ID() card.ID              { return card.DrowningDireRed }
func (DrowningDireRed) Name() string             { return "Drowning Dire" }
func (DrowningDireRed) Cost(*card.TurnState) int { return 2 }
func (DrowningDireRed) Pitch() int               { return 1 }
func (DrowningDireRed) Attack() int              { return 5 }
func (DrowningDireRed) Defense() int             { return 3 }
func (DrowningDireRed) Types() card.TypeSet      { return drowningDireTypes }
func (DrowningDireRed) GoAgain() bool            { return false }

// not implemented: on-hit "may put a non-attack action card from your graveyard on the bottom
// of your deck" rider
func (DrowningDireRed) NotImplemented() {}
func (DrowningDireRed) Play(s *card.TurnState, self *card.CardState) {
	drowningDirePlay(s, self)
}

type DrowningDireYellow struct{}

func (DrowningDireYellow) ID() card.ID              { return card.DrowningDireYellow }
func (DrowningDireYellow) Name() string             { return "Drowning Dire" }
func (DrowningDireYellow) Cost(*card.TurnState) int { return 2 }
func (DrowningDireYellow) Pitch() int               { return 2 }
func (DrowningDireYellow) Attack() int              { return 4 }
func (DrowningDireYellow) Defense() int             { return 3 }
func (DrowningDireYellow) Types() card.TypeSet      { return drowningDireTypes }
func (DrowningDireYellow) GoAgain() bool            { return false }

// not implemented: on-hit "may put a non-attack action card from your graveyard on the bottom
// of your deck" rider
func (DrowningDireYellow) NotImplemented() {}
func (DrowningDireYellow) Play(s *card.TurnState, self *card.CardState) {
	drowningDirePlay(s, self)
}

type DrowningDireBlue struct{}

func (DrowningDireBlue) ID() card.ID              { return card.DrowningDireBlue }
func (DrowningDireBlue) Name() string             { return "Drowning Dire" }
func (DrowningDireBlue) Cost(*card.TurnState) int { return 2 }
func (DrowningDireBlue) Pitch() int               { return 3 }
func (DrowningDireBlue) Attack() int              { return 3 }
func (DrowningDireBlue) Defense() int             { return 3 }
func (DrowningDireBlue) Types() card.TypeSet      { return drowningDireTypes }
func (DrowningDireBlue) GoAgain() bool            { return false }

// not implemented: on-hit "may put a non-attack action card from your graveyard on the bottom
// of your deck" rider
func (DrowningDireBlue) NotImplemented() {}
func (DrowningDireBlue) Play(s *card.TurnState, self *card.CardState) {
	drowningDirePlay(s, self)
}
