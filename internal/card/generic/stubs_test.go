package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// stubCard is a configurable Card implementation used across generic tests to build
// CardsRemaining / CardsPlayed / Pitched lists with specific type, cost, power, and pitch shapes.
// Zero-value fields mean "don't care" — tests set only what the helper under test predicates on.
type stubCard struct {
	name  string
	cost  int
	power int
	pitch int
	types card.TypeSet
}

func (s stubCard) ID() card.ID                  { return card.Invalid }
func (s stubCard) Name() string                 { return s.name }
func (s stubCard) Cost() int                    { return s.cost }
func (s stubCard) Pitch() int                   { return s.pitch }
func (s stubCard) Attack() int                  { return s.power }
func (s stubCard) Defense() int                 { return 0 }
func (s stubCard) Types() card.TypeSet          { return s.types }
func (s stubCard) GoAgain() bool                { return false }
func (s stubCard) Play(*card.TurnState) int     { return 0 }

// stubGenericAttack returns a Generic Action - Attack stub with the given cost and base power.
// Pitch defaults to 1; override via the pitch field if a test cares.
func stubGenericAttack(cost, power int) stubCard {
	return stubCard{
		name:  "stubGenericAttack",
		cost:  cost,
		power: power,
		pitch: 1,
		types: card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack),
	}
}

// stubGenericAttackPitch is stubGenericAttack with an explicit pitch value. Flying High's red
// variant rider reads pitch, so tests that exercise the +1 bonus set this.
func stubGenericAttackPitch(cost, power, pitch int) stubCard {
	s := stubGenericAttack(cost, power)
	s.pitch = pitch
	return s
}

// stubGenericAction returns a Generic Action (non-attack) stub. Used to confirm attack-typed
// lookaheads reject it.
func stubGenericAction() stubCard {
	return stubCard{
		name:  "stubGenericAction",
		types: card.NewTypeSet(card.TypeGeneric, card.TypeAction),
	}
}

// stubGenericAura returns a Generic Aura stub. Used by Yinti Yanti's HasPlayedType(TypeAura) check.
func stubGenericAura() stubCard {
	return stubCard{
		name:  "stubGenericAura",
		types: card.NewTypeSet(card.TypeGeneric, card.TypeAura),
	}
}
