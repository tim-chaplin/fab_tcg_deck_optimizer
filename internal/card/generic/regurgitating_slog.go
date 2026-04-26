// Regurgitating Slog — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Regurgitating Slog, you may banish a card named Sloggism
// from your graveyard. If you do, Regurgitating Slog gains **dominate**."
//
// Modelling: the Dominate grant is gated on banishing a Sloggism — an additional cost the sim
// doesn't evaluate, so the card neither implements card.Dominator nor sets
// self.GrantedDominate. Wiring it would over-credit lines without a Sloggism to banish.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var regurgitatingSlogTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type RegurgitatingSlogRed struct{}

func (RegurgitatingSlogRed) ID() card.ID              { return card.RegurgitatingSlogRed }
func (RegurgitatingSlogRed) Name() string             { return "Regurgitating Slog" }
func (RegurgitatingSlogRed) Cost(*card.TurnState) int { return 2 }
func (RegurgitatingSlogRed) Pitch() int               { return 1 }
func (RegurgitatingSlogRed) Attack() int              { return 6 }
func (RegurgitatingSlogRed) Defense() int             { return 2 }
func (RegurgitatingSlogRed) Types() card.TypeSet      { return regurgitatingSlogTypes }
func (RegurgitatingSlogRed) GoAgain() bool            { return false }

// not implemented: Sloggism graveyard-banish Dominate grant (additional cost not evaluated,
// so the grant never fires)
func (RegurgitatingSlogRed) NotImplemented() {}
func (c RegurgitatingSlogRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type RegurgitatingSlogYellow struct{}

func (RegurgitatingSlogYellow) ID() card.ID              { return card.RegurgitatingSlogYellow }
func (RegurgitatingSlogYellow) Name() string             { return "Regurgitating Slog" }
func (RegurgitatingSlogYellow) Cost(*card.TurnState) int { return 2 }
func (RegurgitatingSlogYellow) Pitch() int               { return 2 }
func (RegurgitatingSlogYellow) Attack() int              { return 5 }
func (RegurgitatingSlogYellow) Defense() int             { return 2 }
func (RegurgitatingSlogYellow) Types() card.TypeSet      { return regurgitatingSlogTypes }
func (RegurgitatingSlogYellow) GoAgain() bool            { return false }

// not implemented: Sloggism graveyard-banish Dominate grant (additional cost not evaluated,
// so the grant never fires)
func (RegurgitatingSlogYellow) NotImplemented() {}
func (c RegurgitatingSlogYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type RegurgitatingSlogBlue struct{}

func (RegurgitatingSlogBlue) ID() card.ID              { return card.RegurgitatingSlogBlue }
func (RegurgitatingSlogBlue) Name() string             { return "Regurgitating Slog" }
func (RegurgitatingSlogBlue) Cost(*card.TurnState) int { return 2 }
func (RegurgitatingSlogBlue) Pitch() int               { return 3 }
func (RegurgitatingSlogBlue) Attack() int              { return 4 }
func (RegurgitatingSlogBlue) Defense() int             { return 2 }
func (RegurgitatingSlogBlue) Types() card.TypeSet      { return regurgitatingSlogTypes }
func (RegurgitatingSlogBlue) GoAgain() bool            { return false }

// not implemented: Sloggism graveyard-banish Dominate grant (additional cost not evaluated,
// so the grant never fires)
func (RegurgitatingSlogBlue) NotImplemented() {}
func (c RegurgitatingSlogBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
