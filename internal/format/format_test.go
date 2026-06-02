package format

import "testing"

// fakeCard is a minimal format.Card: a name is all a Format needs to judge legality.
type fakeCard struct{ name string }

func (c fakeCard) Name() string { return c.name }

// TestSilverAgeIsCardLegal_Banned spot-checks that an implemented banned card (its exact
// Name()) is rejected. Fiddler's Green uses a straight apostrophe — the case that would break
// if the banlist were transcribed with a curly one.
func TestSilverAgeIsCardLegal_Banned(t *testing.T) {
	banned := []string{
		"Sirens of Safe Harbor",
		"Plunder Run",
		"Fiddler's Green",
		"Fate Foreseen",
		"Rosetta Thorn",
	}
	for _, name := range banned {
		if SilverAge.IsCardLegal(fakeCard{name}) {
			t.Errorf("IsCardLegal(%q) = true, want false (on the Silver Age banlist)", name)
		}
	}
}

// TestSilverAgeIsCardLegal_Legal confirms an ordinary card is legal.
func TestSilverAgeIsCardLegal_Legal(t *testing.T) {
	if !SilverAge.IsCardLegal(fakeCard{"Aether Slash"}) {
		t.Error("IsCardLegal(\"Aether Slash\") = false, want true (not banned)")
	}
}

// TestSilverAgeIsRarityLegal checks the rarity rule: Basic/Common/Rare are legal, everything
// rarer (and the empty string) is not.
func TestSilverAgeIsRarityLegal(t *testing.T) {
	for _, r := range []string{"Basic", "Common", "Rare"} {
		if !SilverAge.IsRarityLegal(r) {
			t.Errorf("IsRarityLegal(%q) = false, want true", r)
		}
	}
	for _, r := range []string{"Super Rare", "Majestic", "Legendary", "Promo", ""} {
		if SilverAge.IsRarityLegal(r) {
			t.Errorf("IsRarityLegal(%q) = true, want false", r)
		}
	}
}

// TestBannedNames checks the banlist accessor returns sorted names including a known ban and
// excluding a legal card.
func TestBannedNames(t *testing.T) {
	names := SilverAge.BannedNames()
	if len(names) == 0 {
		t.Fatal("BannedNames returned empty")
	}
	for i := 1; i < len(names); i++ {
		if names[i-1] > names[i] {
			t.Errorf("BannedNames not sorted: %q before %q", names[i-1], names[i])
		}
	}
	has := func(n string) bool {
		for _, x := range names {
			if x == n {
				return true
			}
		}
		return false
	}
	if !has("Ball Lightning") {
		t.Error("BannedNames missing a known banned card (Ball Lightning)")
	}
	if has("Aether Slash") {
		t.Error("BannedNames includes a legal card (Aether Slash)")
	}
}

// TestParse covers the known format and the self-describing error on an unknown one.
func TestParse(t *testing.T) {
	f, err := Parse("silver_age")
	if err != nil {
		t.Fatalf("Parse(\"silver_age\") returned error: %v", err)
	}
	if f != SilverAge {
		t.Errorf("Parse(\"silver_age\") = %v, want SilverAge", f)
	}
	if _, err := Parse("classic_constructed"); err == nil {
		t.Error("Parse(\"classic_constructed\") = nil error, want unknown-format error")
	}
}
