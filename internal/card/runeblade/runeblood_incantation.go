// Runeblood Incantation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. Runeblood Incantation enters the arena with N verse counters on it. At the
// beginning of your action phase, remove a verse counter. If you do, create a Runechant token.
// Otherwise, destroy Runeblood Incantation." (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling: Play flips AuraCreated and registers a start-of-turn AuraTrigger with Count=N.
// Each subsequent turn the sim fires the trigger — the handler creates one live Runechant —
// and decrements Count. After N turns Count hits zero and the sim graveyards the aura.
// Same-turn Play credits 0; every rune comes from a real future-turn fire.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runebloodIncantationTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type RunebloodIncantationRed struct{}

func (RunebloodIncantationRed) ID() card.ID              { return card.RunebloodIncantationRed }
func (RunebloodIncantationRed) Name() string             { return "Runeblood Incantation (Red)" }
func (RunebloodIncantationRed) Cost(*card.TurnState) int { return 1 }
func (RunebloodIncantationRed) Pitch() int               { return 1 }
func (RunebloodIncantationRed) Attack() int              { return 0 }
func (RunebloodIncantationRed) Defense() int             { return 2 }
func (RunebloodIncantationRed) Types() card.TypeSet      { return runebloodIncantationTypes }
func (RunebloodIncantationRed) GoAgain() bool            { return true }
func (RunebloodIncantationRed) AddsFutureValue()         {}
func (c RunebloodIncantationRed) Play(s *card.TurnState, _ *card.CardState) int {
	return runebloodPlay(s, c, 3)
}

type RunebloodIncantationYellow struct{}

func (RunebloodIncantationYellow) ID() card.ID              { return card.RunebloodIncantationYellow }
func (RunebloodIncantationYellow) Name() string             { return "Runeblood Incantation (Yellow)" }
func (RunebloodIncantationYellow) Cost(*card.TurnState) int { return 1 }
func (RunebloodIncantationYellow) Pitch() int               { return 2 }
func (RunebloodIncantationYellow) Attack() int              { return 0 }
func (RunebloodIncantationYellow) Defense() int             { return 2 }
func (RunebloodIncantationYellow) Types() card.TypeSet      { return runebloodIncantationTypes }
func (RunebloodIncantationYellow) GoAgain() bool            { return true }
func (RunebloodIncantationYellow) AddsFutureValue()         {}
func (c RunebloodIncantationYellow) Play(s *card.TurnState, _ *card.CardState) int {
	return runebloodPlay(s, c, 2)
}

type RunebloodIncantationBlue struct{}

func (RunebloodIncantationBlue) ID() card.ID              { return card.RunebloodIncantationBlue }
func (RunebloodIncantationBlue) Name() string             { return "Runeblood Incantation (Blue)" }
func (RunebloodIncantationBlue) Cost(*card.TurnState) int { return 1 }
func (RunebloodIncantationBlue) Pitch() int               { return 3 }
func (RunebloodIncantationBlue) Attack() int              { return 0 }
func (RunebloodIncantationBlue) Defense() int             { return 2 }
func (RunebloodIncantationBlue) Types() card.TypeSet      { return runebloodIncantationTypes }
func (RunebloodIncantationBlue) GoAgain() bool            { return true }
func (RunebloodIncantationBlue) AddsFutureValue()         {}
func (c RunebloodIncantationBlue) Play(s *card.TurnState, _ *card.CardState) int {
	return runebloodPlay(s, c, 1)
}

// runebloodPlay flips AuraCreated and registers the shared start-of-turn trigger with
// Count=n. Each future turn fires the handler (one Runechant per fire) and ticks Count down;
// the sim graveyards the aura when Count hits zero. Same-turn Play returns 0 — every rune
// is credited at its real future-turn fire, no flat over-credit.
func runebloodPlay(s *card.TurnState, self card.Card, n int) int {
	s.AuraCreated = true
	s.AddAuraTrigger(card.AuraTrigger{
		Self:    self,
		Type:    card.TriggerStartOfTurn,
		Count:   n,
		Handler: func(s *card.TurnState) int { return s.CreateRunechants(1) },
	})
	return 0
}
