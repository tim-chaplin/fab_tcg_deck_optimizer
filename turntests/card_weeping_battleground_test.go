package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// TestWeepingBattleground_AuraInGraveyard: an aura in the graveyard gets banished for 1 arcane.
func TestWeepingBattleground_AuraInGraveyard(t *testing.T) {
	aura := cards.SigilOfSilphidaeBlue{}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{aura}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WeepingBattlegroundRed{}})
	if got := ge.Value(); got != 1 {
		t.Fatalf("Play() = %d, want 1", got)
	}
	if !ge.ArcaneDamageDealt() {
		t.Errorf("want ArcaneDamageDealt set")
	}
}

// TestWeepingBattleground_NoAuraInGraveyard: nothing banishable means the clause fizzles and
// Play returns 0. ArcaneDamageDealt stays false.
func TestWeepingBattleground_NoAuraInGraveyard(t *testing.T) {
	// Shrill of Skullform is an Action - Attack (no Aura type), so it fails the aura scan.
	nonAura := cards.ShrillOfSkullformRed{}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{nonAura}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WeepingBattlegroundRed{}})
	if got := ge.Value(); got != 0 {
		t.Fatalf("Play() = %d, want 0 (no aura to banish)", got)
	}
	if ge.ArcaneDamageDealt() {
		t.Errorf("ArcaneDamageDealt should stay false when the banish clause fizzles")
	}
}

// TestWeepingBattleground_EmptyGraveyard: no graveyard at all means no banish, no damage.
func TestWeepingBattleground_EmptyGraveyard(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WeepingBattlegroundRed{}})
	if got := ge.Value(); got != 0 {
		t.Fatalf("Play() = %d, want 0", got)
	}
}

// TestWeepingBattleground_BanishesAura: after Play, the aura must be removed from Graveyard and
// present in Banish.
func TestWeepingBattleground_BanishesAura(t *testing.T) {
	aura := cards.SigilOfSilphidaeBlue{}
	nonAura := cards.ShrillOfSkullformRed{}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{nonAura, aura}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WeepingBattlegroundRed{}})
	g := ge.Graveyard()
	if len(g) != 1 || g[0].ID() != nonAura.ID() {
		t.Errorf("want graveyard with only non-aura left, got %+v", g)
	}
	if len(ge.Banished()) != 1 || ge.Banished()[0].ID() != aura.ID() {
		t.Errorf("want banish = [aura], got %+v", ge.Banished())
	}
}

// TestWeepingBattleground_OnlyOneAuraBanished: two auras in the graveyard — exactly one is
// banished per Play. The other stays behind so subsequent effects (or the next Weeping
// Battleground) can use it.
func TestWeepingBattleground_OnlyOneAuraBanished(t *testing.T) {
	aura1 := cards.SigilOfSilphidaeBlue{}
	aura2 := cards.SigilOfSilphidaeBlue{}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{aura1, aura2}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WeepingBattlegroundRed{}})
	if g := ge.Graveyard(); len(g) != 1 {
		t.Fatalf("want one aura left in graveyard, got %d", len(g))
	}
	if len(ge.Banished()) != 1 {
		t.Fatalf("want one aura banished, got %d", len(ge.Banished()))
	}
}

// TestWeepingBattleground_SecondCopyAlsoFires: two Weeping Battlegrounds against a graveyard
// with two auras each banish one. The second one must still fire because Play mutates state.
// Each fire credits +1, so cumulative gs.Value() = 2.
func TestWeepingBattleground_SecondCopyAlsoFires(t *testing.T) {
	aura1 := cards.SigilOfSilphidaeBlue{}
	aura2 := cards.SigilOfSilphidaeBlue{}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{aura1, aura2}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WeepingBattlegroundRed{}})
	if got := ge.Value(); got != 1 {
		t.Fatalf("first Play() Value = %d, want 1", got)
	}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WeepingBattlegroundBlue{}})
	if got := ge.Value(); got != 2 {
		t.Fatalf("second Play() cumulative Value = %d, want 2", got)
	}
	if g := ge.Graveyard(); len(g) != 0 {
		t.Errorf("want empty graveyard after two banishes, got %d", len(g))
	}
	if len(ge.Banished()) != 2 {
		t.Errorf("want 2 banished, got %d", len(ge.Banished()))
	}
}

// TestWeepingBattleground_SecondCopyFizzlesWhenOutOfAuras: one aura, two Weeping Battlegrounds —
// the first banishes it, the second contributes 0 (no aura left), so cumulative Value stays 1.
func TestWeepingBattleground_SecondCopyFizzlesWhenOutOfAuras(t *testing.T) {
	aura := cards.SigilOfSilphidaeBlue{}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{aura}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WeepingBattlegroundRed{}})
	if got := ge.Value(); got != 1 {
		t.Fatalf("first Play() Value = %d, want 1", got)
	}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WeepingBattlegroundBlue{}})
	if got := ge.Value(); got != 1 {
		t.Fatalf("second Play() cumulative Value = %d, want 1 (no aura left, fizzle adds 0)", got)
	}
}
