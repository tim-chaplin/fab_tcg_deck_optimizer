// Malefic Incantation — Runeblade Action - Aura. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "This enters the arena with N verse counters. When it has none, destroy it. Once per
// turn, when you play an attack action card, remove a verse counter from this. If you do,
// create a Runechant token." (Red N=3, Yellow N=2, Blue N=1.)
//
// AttackAction trigger with Count=N and OncePerTurn=true: each turn's first attack action
// creates 1 Runechant and burns one verse counter.

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var maleficTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type MaleficIncantationRed struct{}

func (MaleficIncantationRed) ID() card.ID              { return card.MaleficIncantationRed }
func (MaleficIncantationRed) Name() string             { return "Malefic Incantation" }
func (MaleficIncantationRed) Cost(*card.TurnState) int { return 0 }
func (MaleficIncantationRed) Pitch() int               { return 1 }
func (MaleficIncantationRed) Attack() int              { return 0 }
func (MaleficIncantationRed) Defense() int             { return 2 }
func (MaleficIncantationRed) Types() card.TypeSet      { return maleficTypes }
func (MaleficIncantationRed) GoAgain() bool            { return true }
func (MaleficIncantationRed) AddsFutureValue()         {}
func (c MaleficIncantationRed) Play(s *card.TurnState, self *card.CardState) {
	maleficPlay(s, self, c, 3)
}

type MaleficIncantationYellow struct{}

func (MaleficIncantationYellow) ID() card.ID              { return card.MaleficIncantationYellow }
func (MaleficIncantationYellow) Name() string             { return "Malefic Incantation" }
func (MaleficIncantationYellow) Cost(*card.TurnState) int { return 0 }
func (MaleficIncantationYellow) Pitch() int               { return 2 }
func (MaleficIncantationYellow) Attack() int              { return 0 }
func (MaleficIncantationYellow) Defense() int             { return 2 }
func (MaleficIncantationYellow) Types() card.TypeSet      { return maleficTypes }
func (MaleficIncantationYellow) GoAgain() bool            { return true }
func (MaleficIncantationYellow) AddsFutureValue()         {}
func (c MaleficIncantationYellow) Play(s *card.TurnState, self *card.CardState) {
	maleficPlay(s, self, c, 2)
}

type MaleficIncantationBlue struct{}

func (MaleficIncantationBlue) ID() card.ID              { return card.MaleficIncantationBlue }
func (MaleficIncantationBlue) Name() string             { return "Malefic Incantation" }
func (MaleficIncantationBlue) Cost(*card.TurnState) int { return 0 }
func (MaleficIncantationBlue) Pitch() int               { return 3 }
func (MaleficIncantationBlue) Attack() int              { return 0 }
func (MaleficIncantationBlue) Defense() int             { return 2 }
func (MaleficIncantationBlue) Types() card.TypeSet      { return maleficTypes }
func (MaleficIncantationBlue) GoAgain() bool            { return true }
func (MaleficIncantationBlue) AddsFutureValue()         {}
func (c MaleficIncantationBlue) Play(s *card.TurnState, self *card.CardState) {
	maleficPlay(s, self, c, 1)
}

// maleficPlay registers the attack-action once-per-turn trigger and emits the same-turn
// chain step. Each trigger fire creates one Runechant — the trigger handler authors a
// post-trigger log line so it groups beneath the triggering attack-action chain step.
func maleficPlay(s *card.TurnState, selfState *card.CardState, selfCard card.Card, n int) {
	s.AddAuraTrigger(card.AuraTrigger{
		Self:        selfCard,
		Type:        card.TriggerAttackAction,
		Count:       n,
		OncePerTurn: true,
		Handler: func(s *card.TurnState) int {
			return s.AddPostTriggerLogEntry(
				card.DisplayName(selfCard)+" created a runechant",
				card.DisplayName(s.TriggeringCard),
				s.CreateRunechants(1),
			)
		},
	})
	s.LogPlay(selfState)
}
