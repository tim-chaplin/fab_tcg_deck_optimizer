package gameengine

import (
	"sync"

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

// equipEnginePool supplies the scratch engine EquipFromCards plays weapon cards against, so a
// per-shuffle equip reuses one GameState instead of allocating a fresh ~KB struct each call. A
// weapon card's Play only ever calls CreateWeapon (which just appends to gs.weapons), so the
// pooled state needs no reset beyond clearing weapons.
var equipEnginePool = sync.Pool{New: func() any { return &GameEngine{GameState: &GameState{}} }}

// EquipFromCards builds the engine-side weapon objects for the supplied platonic weapon cards by
// playing each one: the card's Play registers its equipped object (and any trigger) via
// ge.CreateWeapon, exactly as an item card's Play calls ge.CreateItem. Weapons aren't played from
// hand — this is the game-start equip step the sim and StateBuilder drive. Returns nil for an
// empty loadout so the no-weapon case carries no allocation.
func EquipFromCards(cards []weapon.Card) []Weapon {
	if len(cards) == 0 {
		return nil
	}
	ge := equipEnginePool.Get().(*GameEngine)
	var cs card.CardState
	for _, c := range cards {
		cs = card.CardState{Card: c}
		c.Play(ge, NoopLogger{}, &cs)
	}
	out := ge.weapons
	ge.weapons = nil
	equipEnginePool.Put(ge)
	return out
}

// DestroyWeapon removes the weapon currently being fired from the arena and, when
// addToGraveyard is true, pushes its source weapon card into the graveyard. The weapon
// counterpart of DestroyAura / DestroyItem: direct splice with no cacheable flip —
// destruction is deterministic from the triggering event. Talishar's end-phase
// self-destruct routes its handler-side Destroy back here.
func (ge *GameEngine) DestroyWeapon(addToGraveyard bool) {
	i := ge.currentHookIdx
	if i < 0 || i >= len(ge.weapons) {
		return
	}
	src := ge.weapons[i].SourceCard()
	ge.weapons = append(ge.weapons[:i], ge.weapons[i+1:]...)
	ge.currentHookDestroyed = true
	if src != nil && addToGraveyard {
		ge.AppendGraveyard(src.(card.Card))
	}
}
