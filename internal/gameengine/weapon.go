package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// CreateWeapon equips a card-sourced weapon — the weapon counterpart of CreateItem / CreateAura.
// A weapon card's Play registers its equipped object here; source supplies the swing Ability and
// Hands. tt / handler give a self-triggering weapon its scheduled handler (Talishar's end-phase
// self-destruct); pass (0, nil) for an untriggered weapon.
func (gs *GameState) CreateWeapon(source card.Card, tt triggertype.Type, handler func(card.GameEngine, card.Logger, card.Weapon, card.FireContext), oncePerTurn bool, filter func(card.TypeSet) bool) {
	gs.weapons = append(gs.weapons, weapon.NewFromCard(source.(weapon.Card), tt, handler, oncePerTurn, filter))
}

// EquipFromCards equips each platonic weapon card onto gs by playing it: the card's Play
// registers its object (and any trigger) via ge.CreateWeapon — the same way an item card's
// Play calls ge.CreateItem, and StateBuilder.CreateItemFromCard plays a card into the state
// being built. Weapons aren't played from hand; this is the game-start equip step the sim and
// StateBuilder drive.
func EquipFromCards(gs *GameState, cards []weapon.Card) {
	if len(cards) == 0 {
		return
	}
	ge := &GameEngine{GameState: gs}
	var cs card.CardState
	for _, c := range cards {
		cs = card.CardState{Card: c}
		c.Play(ge, NoopLogger{}, &cs)
	}
}

// DestroyWeaponObject removes the weapon equal to w from the arena, located by object identity
// and, when addToGraveyard is true, pushes its source weapon card into the graveyard. Talishar's
// end-phase self-destruct routes its handler-side self.Destroy() back here.
func (ge *GameEngine) DestroyWeaponObject(w card.Weapon, addToGraveyard bool) {
	target := any(w)
	// Fast path: a trigger handler's self.Destroy() (Talishar's end-phase) removes the firing
	// weapon — check it directly so the hot self-destruct stays index-direct instead of scanning.
	if i := ge.currentHookIdx; i >= 0 && i < len(ge.weapons) && any(ge.weapons[i]) == target {
		ge.currentHookDestroyed = true
		ge.spliceWeaponAt(i, addToGraveyard)
		return
	}
	for i := range ge.weapons {
		if any(ge.weapons[i]) == target {
			ge.spliceWeaponAt(i, addToGraveyard)
			return
		}
	}
}

// spliceWeaponAt removes ge.weapons[i] and graveyards its source weapon card when asked.
func (ge *GameEngine) spliceWeaponAt(i int, addToGraveyard bool) {
	src := ge.weapons[i].SourceCard()
	ge.weapons = append(ge.weapons[:i], ge.weapons[i+1:]...)
	if src != nil && addToGraveyard {
		ge.AppendGraveyard(src.(card.Card))
	}
}
