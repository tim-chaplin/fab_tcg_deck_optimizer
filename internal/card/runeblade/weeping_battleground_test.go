package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestWeepingBattleground_AuraInGraveyard: an aura in the graveyard gets banished for 1 arcane.
func TestWeepingBattleground_AuraInGraveyard(t *testing.T) {
	aura := SigilOfSilphidaeBlue{}
	s := card.TurnState{Graveyard: []card.Card{aura}}
	(WeepingBattlegroundRed{}).Play(&s, &card.CardState{Card: WeepingBattlegroundRed{}})
	if got := s.Value; got != 1 {
		t.Fatalf("Play() = %d, want 1", got)
	}
	if !s.ArcaneDamageDealt {
		t.Errorf("want ArcaneDamageDealt set")
	}
}

// TestWeepingBattleground_NoAuraInGraveyard: nothing banishable means the clause fizzles and
// Play returns 0. ArcaneDamageDealt stays false.
func TestWeepingBattleground_NoAuraInGraveyard(t *testing.T) {
	// Shrill of Skullform is an Action - Attack (no Aura type), so it fails the aura scan.
	nonAura := ShrillOfSkullformRed{}
	s := card.TurnState{Graveyard: []card.Card{nonAura}}
	(WeepingBattlegroundRed{}).Play(&s, &card.CardState{Card: WeepingBattlegroundRed{}})
	if got := s.Value; got != 0 {
		t.Fatalf("Play() = %d, want 0 (no aura to banish)", got)
	}
	if s.ArcaneDamageDealt {
		t.Errorf("ArcaneDamageDealt should stay false when the banish clause fizzles")
	}
}

// TestWeepingBattleground_EmptyGraveyard: no graveyard at all means no banish, no damage.
func TestWeepingBattleground_EmptyGraveyard(t *testing.T) {
	var s card.TurnState
	(WeepingBattlegroundRed{}).Play(&s, &card.CardState{Card: WeepingBattlegroundRed{}})
	if got := s.Value; got != 0 {
		t.Fatalf("Play() = %d, want 0", got)
	}
}

// TestWeepingBattleground_BanishesAura: after Play, the aura must be removed from Graveyard and
// present in Banish.
func TestWeepingBattleground_BanishesAura(t *testing.T) {
	aura := SigilOfSilphidaeBlue{}
	nonAura := ShrillOfSkullformRed{}
	s := card.TurnState{Graveyard: []card.Card{nonAura, aura}}
	(WeepingBattlegroundRed{}).Play(&s, &card.CardState{Card: WeepingBattlegroundRed{}})
	if len(s.Graveyard) != 1 || s.Graveyard[0].ID() != nonAura.ID() {
		t.Errorf("want graveyard with only non-aura left, got %+v", s.Graveyard)
	}
	if len(s.Banish) != 1 || s.Banish[0].ID() != aura.ID() {
		t.Errorf("want banish = [aura], got %+v", s.Banish)
	}
}

// TestWeepingBattleground_OnlyOneAuraBanished: two auras in the graveyard — exactly one is
// banished per Play. The other stays behind so subsequent effects (or the next Weeping
// Battleground) can use it.
func TestWeepingBattleground_OnlyOneAuraBanished(t *testing.T) {
	aura1 := SigilOfSilphidaeBlue{}
	aura2 := SigilOfSilphidaeBlue{}
	s := card.TurnState{Graveyard: []card.Card{aura1, aura2}}
	(WeepingBattlegroundRed{}).Play(&s, &card.CardState{Card: WeepingBattlegroundRed{}})
	if len(s.Graveyard) != 1 {
		t.Fatalf("want one aura left in graveyard, got %d", len(s.Graveyard))
	}
	if len(s.Banish) != 1 {
		t.Fatalf("want one aura banished, got %d", len(s.Banish))
	}
}

// TestWeepingBattleground_SecondCopyAlsoFires: two Weeping Battlegrounds against a graveyard
// with two auras each banish one. The second one must still fire because Play mutates state.
// Each fire credits +1, so cumulative s.Value = 2.
func TestWeepingBattleground_SecondCopyAlsoFires(t *testing.T) {
	aura1 := SigilOfSilphidaeBlue{}
	aura2 := SigilOfSilphidaeBlue{}
	s := card.TurnState{Graveyard: []card.Card{aura1, aura2}}
	(WeepingBattlegroundRed{}).Play(&s, &card.CardState{Card: WeepingBattlegroundRed{}})
	if got := s.Value; got != 1 {
		t.Fatalf("first Play() Value = %d, want 1", got)
	}
	(WeepingBattlegroundBlue{}).Play(&s, &card.CardState{Card: WeepingBattlegroundBlue{}})
	if got := s.Value; got != 2 {
		t.Fatalf("second Play() cumulative Value = %d, want 2", got)
	}
	if len(s.Graveyard) != 0 {
		t.Errorf("want empty graveyard after two banishes, got %d", len(s.Graveyard))
	}
	if len(s.Banish) != 2 {
		t.Errorf("want 2 banished, got %d", len(s.Banish))
	}
}

// TestWeepingBattleground_SecondCopyFizzlesWhenOutOfAuras: one aura, two Weeping Battlegrounds —
// the first banishes it, the second contributes 0 (no aura left), so cumulative Value stays 1.
func TestWeepingBattleground_SecondCopyFizzlesWhenOutOfAuras(t *testing.T) {
	aura := SigilOfSilphidaeBlue{}
	s := card.TurnState{Graveyard: []card.Card{aura}}
	(WeepingBattlegroundRed{}).Play(&s, &card.CardState{Card: WeepingBattlegroundRed{}})
	if got := s.Value; got != 1 {
		t.Fatalf("first Play() Value = %d, want 1", got)
	}
	(WeepingBattlegroundBlue{}).Play(&s, &card.CardState{Card: WeepingBattlegroundBlue{}})
	if got := s.Value; got != 1 {
		t.Fatalf("second Play() cumulative Value = %d, want 1 (no aura left, fizzle adds 0)", got)
	}
}
