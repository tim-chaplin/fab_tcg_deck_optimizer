// Sigil of Silphidae — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 3, Go again.
// Only printed in Blue.
// Text: "When this enters or leaves the arena, you may banish another aura from your graveyard.
// If you do, deal 1 arcane damage to target hero. At the beginning of your action phase, destroy
// this."
//
// Modelling: Play resolves the enter trigger directly (banishAuraFromGraveyard scans
// s.Graveyard and credits 1 arcane if an aura lands in s.Banish) and registers a
// start-of-turn AuraTrigger with Count=1 for the "destroy this" clause. Next turn the
// handler runs the LEAVE trigger FIRST — the sim graveyards Self only after Count hits
// zero, so the "another aura" restriction is honoured naturally without an explicit skip.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfSilphidaeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfSilphidaeBlue struct{}

func (SigilOfSilphidaeBlue) ID() card.ID              { return card.SigilOfSilphidaeBlue }
func (SigilOfSilphidaeBlue) Name() string             { return "Sigil of Silphidae (Blue)" }
func (SigilOfSilphidaeBlue) Cost(*card.TurnState) int { return 0 }
func (SigilOfSilphidaeBlue) Pitch() int               { return 3 }
func (SigilOfSilphidaeBlue) Attack() int              { return 0 }
func (SigilOfSilphidaeBlue) Defense() int             { return 3 }
func (SigilOfSilphidaeBlue) Types() card.TypeSet      { return sigilOfSilphidaeTypes }
func (SigilOfSilphidaeBlue) GoAgain() bool            { return true }
func (SigilOfSilphidaeBlue) AddsFutureValue()         {}
func (SigilOfSilphidaeBlue) NoMemo()                  {}
func (c SigilOfSilphidaeBlue) Play(s *card.TurnState, _ *card.CardState) int {
	s.AuraCreated = true
	enterDamage := banishAuraFromGraveyard(s)
	s.AddAuraTrigger(card.AuraTrigger{
		Self:  c,
		Type:  card.TriggerStartOfTurn,
		Count: 1,
		Handler: func(s *card.TurnState) int {
			// The sim graveyards Self only AFTER this handler returns (Count hits zero),
			// so the scan here naturally can't pick up Silphidae itself — the printed
			// "another aura" restriction is satisfied without an explicit skip.
			return banishAuraFromGraveyard(s)
		},
	})
	return enterDamage
}
