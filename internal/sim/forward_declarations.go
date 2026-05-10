package sim

// Forward declarations for the simulator: package-level vars whose real values are supplied
// by another package's init(), letting sim break would-be import cycles with packages that
// depend on sim's types. Defaults panic so a test that forgot the wiring fails loudly. See
// docs/dev-standards.md "Registry / sim split".

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// ChainStepText returns the "<DisplayName>: <VERB>[ from arsenal]" prefix the chain-step
// log line is built from. VERB picks WEAPON ATTACK for weapon activated-ability cards
// (Weapon + Attack), ATTACK for attack-action cards, DEFENSE REACTION for Defense
// Reactions, and PLAY for everything else; the "from arsenal" suffix tags entries played
// out of the arsenal slot. optimizations.init swaps in a memoised version backed by a
// per-(CardID, FromArsenal) cache.
var ChainStepText = func(self *CardState) string {
	types := self.Card.Types()
	var verb string
	switch {
	case types.IsWeaponAttack():
		verb = "WEAPON ATTACK"
	case types.IsAttackAction():
		verb = "ATTACK"
	case types.IsDefenseReaction():
		verb = "DEFENSE REACTION"
	default:
		verb = "PLAY"
	}
	if self.FromArsenal {
		verb += " from arsenal"
	}
	return self.Card.DisplayName() + ": " + verb
}

// GetCard returns the registered Card for id. Populated by simreg.init.
var GetCard = func(id ids.CardID) Card {
	panic("sim.GetCard: registry not loaded — blank-import internal/simreg to populate the hook")
}

// DeckableCards returns every CardID legal to put in a real deck. Populated by simreg.init.
var DeckableCards = func() []ids.CardID {
	panic("sim.DeckableCards: registry not loaded — blank-import internal/simreg to populate the hook")
}

// AllWeapons is the registered weapon roster. Populated by simreg.init.
var AllWeapons []Weapon
